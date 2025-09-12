package main

import (
	"log"
	"net"
	"os"

	"github.com/evrintobing17/ecommerce-system/shared"
	"github.com/evrintobing17/ecommerce-system/shared/jsonhttpresponse"
	"github.com/evrintobing17/ecommerce-system/shared/middleware"
	proto "github.com/evrintobing17/ecommerce-system/shared/proto/warehouse"
	http "github.com/evrintobing17/ecommerce-system/warehouse-service/app/delivery"
	grpcServer "github.com/evrintobing17/ecommerce-system/warehouse-service/app/delivery/grpc"
	"github.com/evrintobing17/ecommerce-system/warehouse-service/app/models"
	"github.com/evrintobing17/ecommerce-system/warehouse-service/app/repository"
	"github.com/evrintobing17/ecommerce-system/warehouse-service/app/usecase"
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
	err = shared.MigrateDB(db, &models.Warehouse{}, &models.Stock{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories
	warehouseRepo := repository.NewWarehouseRepository(db)
	stockRepo := repository.NewStockRepository(db)

	// Initialize use cases
	warehouseUsecase := usecase.NewWarehouseUsecase(warehouseRepo, stockRepo)

	// Initialize HTTP server
	router := gin.Default()
	warehouseHandler := http.NewWarehouseHandler(warehouseUsecase)

	router.Use(gin.Recovery())
	router.Use(shared.GinMetricsMiddleware())
	shared.RegisterMetricsHandler(router)

	router.GET("/health", func(c *gin.Context) {
		jsonhttpresponse.OK(c, gin.H{
			"status": "OK",
		})
	})
	// HTTP routes
	api := router.Group("/api/v1")
	jwtSecret := os.Getenv("JWT_SECRET")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		api.GET("/warehouses/:id", warehouseHandler.GetWarehouse)
		api.GET("/warehouses", warehouseHandler.GetWarehouses)
		api.POST("/warehouses", warehouseHandler.CreateWarehouse)
		api.PUT("/warehouses/:id", warehouseHandler.UpdateWarehouse)
		api.POST("/warehouses/transfer", warehouseHandler.TransferStock)
		api.GET("/warehouses/stock", warehouseHandler.GetStock)
		api.PATCH("/warehouses/stock", warehouseHandler.UpdateStock)
	}

	// Initialize gRPC server
	warehouseServer := grpcServer.NewWarehouseServer(warehouseUsecase)

	// Start gRPC server
	go func() {
		grpcPort := os.Getenv("WAREHOUSE_GRPC_PORT")
		if grpcPort == "" {
			grpcPort = "50055"
		}

		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatal("Failed to listen:", err)
		}

		grpcServer := grpc.NewServer()
		proto.RegisterWarehouseServiceServer(grpcServer, warehouseServer)

		log.Printf("Warehouse gRPC server started on port %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	// Start HTTP server
	httpPort := os.Getenv("WAREHOUSE_SERVICE_PORT")
	if httpPort == "" {
		httpPort = "8084"
	}

	log.Printf("Warehouse HTTP server started on port %s", httpPort)
	if err := router.Run(":" + httpPort); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
