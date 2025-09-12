package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	grpcHandler "github.com/evrintobing17/ecommerce-system/order-service/app/delivery/grpc"
	grpcProduct "github.com/evrintobing17/ecommerce-system/shared/proto/product"
	grpcWarehouse "github.com/evrintobing17/ecommerce-system/shared/proto/warehouse"

	"github.com/evrintobing17/ecommerce-system/shared/grpc_client"
	"github.com/evrintobing17/ecommerce-system/shared/middleware"
	proto "github.com/evrintobing17/ecommerce-system/shared/proto/order"

	delivery "github.com/evrintobing17/ecommerce-system/order-service/app/delivery"
	"github.com/evrintobing17/ecommerce-system/order-service/app/models"
	"github.com/evrintobing17/ecommerce-system/order-service/app/repository"
	"github.com/evrintobing17/ecommerce-system/order-service/app/usecase"
	"google.golang.org/grpc"

	"github.com/evrintobing17/ecommerce-system/shared"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	err = shared.MigrateDB(db, &models.Order{}, &models.OrderItem{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	productServiceAddr := os.Getenv("PRODUCT_SERVICE_GRPC_ADDR")
	if productServiceAddr == "" {
		productServiceAddr = "product-service:50052"
	}

	warehouseServiceAddr := os.Getenv("WAREHOUSE_SERVICE_GRPC_ADDR")
	if warehouseServiceAddr == "" {
		warehouseServiceAddr = "warehouse-service:50055"
	}

	productConn, _ := grpc_client.NewConnection(productServiceAddr)
	defer productConn.Close()

	productClient := grpcProduct.NewProductServiceClient(productConn)

	warehouseConn, _ := grpc_client.NewConnection(warehouseServiceAddr)
	defer warehouseConn.Close()

	warehouseClient := grpcWarehouse.NewWarehouseServiceClient(warehouseConn)

	orderTimeoutMinutes := 15
	if timeoutStr := os.Getenv("ORDER_TIMEOUT_MINUTES"); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil {
			orderTimeoutMinutes = timeout
		}
	}
	orderTimeout := time.Duration(orderTimeoutMinutes) * time.Minute

	// Initialize repositories
	orderRepo := repository.NewOrderRepository(db)

	// Initialize use cases
	orderUsecase := usecase.NewOrderUsecase(orderRepo, productClient, warehouseClient, orderTimeout)
	go func() {
		ticker := time.NewTicker(5 * time.Minute) // Check every 5 minutes
		defer ticker.Stop()

		for range ticker.C {
			if err := orderUsecase.ReleaseExpiredOrders(); err != nil {
				log.Printf("Error releasing expired orders: %v", err)
			}
		}
	}()
	// Initialize HTTP server
	router := gin.Default()
	orderHandler := delivery.NewOrderHandler(orderUsecase)
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
		api.POST("/checkout", orderHandler.Checkout)
		api.POST("/orders", orderHandler.CreateOrder)
		api.GET("/orders/:id", orderHandler.GetOrder)
		api.GET("/orders", orderHandler.GetUserOrders)
		api.POST("/orders/:id/payment", orderHandler.ProcessPayment)
		api.DELETE("/orders/:id", orderHandler.CancelOrder)
	}

	// Initialize gRPC server
	orderServer := grpcHandler.NewOrderServer(orderUsecase)

	// Start gRPC server
	go func() {
		grpcPort := os.Getenv("ORDER_GRPC_PORT")
		if grpcPort == "" {
			grpcPort = "50053"
		}

		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatal("Failed to listen:", err)
		}

		grpcServer := grpc.NewServer()
		proto.RegisterOrderServiceServer(grpcServer, orderServer)

		log.Printf("Order gRPC server started on port %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	// Start HTTP server
	httpPort := os.Getenv("ORDER_SERVICE_PORT")
	if httpPort == "" {
		httpPort = "8082"
	}

	log.Printf("Order HTTP server started on port %s", httpPort)
	if err := router.Run(":" + httpPort); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
