package main

import (
	"log"
	"net"
	"net/http"
	"os"

	delivery "github.com/evrintobing17/ecommerce-system/product-service/app/delivery"
	grpcServer "github.com/evrintobing17/ecommerce-system/product-service/app/delivery/grpc"
	"github.com/evrintobing17/ecommerce-system/product-service/app/models"
	"github.com/evrintobing17/ecommerce-system/product-service/app/repository"
	"github.com/evrintobing17/ecommerce-system/product-service/app/usecase"

	"github.com/evrintobing17/ecommerce-system/shared"
	"github.com/evrintobing17/ecommerce-system/shared/middleware"
	proto "github.com/evrintobing17/ecommerce-system/shared/proto/product"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize logger
	shared.InitLogger()

	// Initialize database
	db, err := shared.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if err := shared.CloseDB(db); err != nil {
			log.Println("Error closing database:", err)
		}
	}()

	// Auto migrate models
	err = shared.MigrateDB(db, &models.Product{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories
	productRepo := repository.NewProductRepository(db)

	// Initialize use cases
	productUsecase := usecase.NewProductUsecase(productRepo)

	// Initialize HTTP server
	router := gin.Default()
	productHandler := delivery.NewProductHandler(productUsecase)
	router.Use(gin.Recovery())
	router.Use(shared.GinMetricsMiddleware())
	shared.RegisterMetricsHandler(router)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})
	// HTTP routes
	api := router.Group("/api/v1")
	jwtSecret := os.Getenv("JWT_SECRET")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		api.GET("/products", productHandler.GetProducts)
		api.GET("/products/:id", productHandler.GetProduct)
		api.POST("/products", productHandler.CreateProduct)
		api.PUT("/products/:id", productHandler.UpdateProduct)
		api.PATCH("/products/:id/stock", productHandler.UpdateStock)
		api.DELETE("/products/:id", productHandler.DeleteProduct)
	}

	// Initialize gRPC server
	productServer := grpcServer.NewProductServer(productUsecase)

	// Start gRPC server
	go func() {
		grpcPort := os.Getenv("PRODUCT_GRPC_PORT")
		if grpcPort == "" {
			grpcPort = "50052"
		}

		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatal("Failed to listen:", err)
		}

		grpcServer := grpc.NewServer()
		proto.RegisterProductServiceServer(grpcServer, productServer)

		log.Printf("Product gRPC server started on port %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	// Start HTTP server
	httpPort := os.Getenv("PRODUCT_SERVICE_PORT")
	if httpPort == "" {
		httpPort = "8081"
	}

	log.Printf("Product HTTP server started on port %s", httpPort)
	if err := router.Run(":" + httpPort); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
