package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"shared"
	"shared/jsonhttpresponse"
	productDelivery "product-service/app/delivery"
	productRepo "product-service/app/repository"
	productUsecase "product-service/app/usecase"
)

var (
	db        *sql.DB
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)

func main() {
	logger := shared.NewLogger("PRODUCT-SERVICE")
	var err error
	db, err = shared.InitDB()
	if err != nil {
		logger.ErrorLog(err)
		log.Fatal(err)
	}

	defer db.Close()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(shared.GinMetricsMiddleware())
	shared.RegisterMetricsHandler(r)

	r.GET("/health", func(c *gin.Context) {
		jsonhttpresponse.OK(c, gin.H{
			"status": "OK",
		})
	})

	productRepository := productRepo.NewproductRepository(db)
	productUseCase := productUsecase.NewproductUsecase(productRepository, string(jwtSecret), logger)
	productDelivery.NewAuthHandler(r, logger, productUseCase)

	fmt.Println("product service is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
