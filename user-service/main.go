package main

import (
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/evrintobing17/ecommerce-system/shared"
	userDelivery "github.com/evrintobing17/ecommerce-system/user-service/app/delivery"
	"github.com/evrintobing17/ecommerce-system/user-service/app/models"
	userRepo "github.com/evrintobing17/ecommerce-system/user-service/app/repository"
	userUsecase "github.com/evrintobing17/ecommerce-system/user-service/app/usecase"

	"github.com/evrintobing17/ecommerce-system/shared/jsonhttpresponse"
	"github.com/evrintobing17/ecommerce-system/shared/middleware"
	proto "github.com/evrintobing17/ecommerce-system/shared/proto/user"
	userGrpc "github.com/evrintobing17/ecommerce-system/user-service/app/delivery/grpc"
)

func main() {
	shared.InitLogger()
	var err error
	db, err := shared.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate models
	err = shared.MigrateDB(db, &models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(shared.GinMetricsMiddleware())
	shared.RegisterMetricsHandler(r)

	r.GET("/health", func(c *gin.Context) {
		jsonhttpresponse.OK(c, gin.H{
			"status": "OK",
		})
	})
	jwtSecret := os.Getenv("JWT_SECRET")

	userRepository := userRepo.NewUserRepository(db)
	userUseCase := userUsecase.NewUserUsecase(userRepository, string(jwtSecret))
	// Initialize HTTP server
	router := gin.Default()
	userHandler := userDelivery.NewUserHandler(userUseCase)

	api := router.Group("/api/v1")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
	}
	private := api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		private.GET("/profile", userHandler.GetProfile)
		private.PUT("/profile", userHandler.UpdateProfile)
	}

	// Initialize gRPC server
	userServer := userGrpc.NewUserServer(userUseCase)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50058")
		if err != nil {
			log.Fatal("Failed to listen:", err)
		}

		grpcServer := grpc.NewServer()
		proto.RegisterUserServiceServer(grpcServer, userServer)

		log.Println("User gRPC server started on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	// Start HTTP server
	log.Println("User HTTP server started on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
