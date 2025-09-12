package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/evrintobing17/ecommerce-system/shared"
	"github.com/evrintobing17/ecommerce-system/shared/middleware"
	proto "github.com/evrintobing17/ecommerce-system/shared/proto/shop"

	delivery "github.com/evrintobing17/ecommerce-system/shop-service/app/delivery"
	grpcServer "github.com/evrintobing17/ecommerce-system/shop-service/app/delivery/grpc"

	"github.com/evrintobing17/ecommerce-system/shop-service/app/models"
	"github.com/evrintobing17/ecommerce-system/shop-service/app/repository"
	"github.com/evrintobing17/ecommerce-system/shop-service/app/usecase"
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
	err = shared.MigrateDB(db, &models.Shop{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories
	shopRepo := repository.NewShopRepository(db)

	// Initialize use cases
	shopUsecase := usecase.NewShopUsecase(shopRepo)

	// Initialize HTTP server
	router := gin.Default()
	shopHandler := delivery.NewShopHandler(shopUsecase)
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
		api.POST("/shops", shopHandler.CreateShop)
		api.GET("/shops/:id", shopHandler.GetShop)
		api.GET("/shops", shopHandler.GetMyShops)
		api.PUT("/shops/:id", shopHandler.UpdateShop)
		api.DELETE("/shops/:id", shopHandler.DeleteShop)
	}

	// Initialize gRPC server
	shopServer := grpcServer.NewShopServer(shopUsecase)

	// Start gRPC server
	go func() {
		grpcPort := os.Getenv("SHOP_GRPC_PORT")
		if grpcPort == "" {
			grpcPort = "50054"
		}

		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatal("Failed to listen:", err)
		}

		grpcServer := grpc.NewServer()
		proto.RegisterShopServiceServer(grpcServer, shopServer)

		log.Printf("Shop gRPC server started on port %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	// Start HTTP server
	httpPort := os.Getenv("SHOP_SERVICE_PORT")
	if httpPort == "" {
		httpPort = "8083"
	}

	log.Printf("Shop HTTP server started on port %s", httpPort)
	if err := router.Run(":" + httpPort); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
