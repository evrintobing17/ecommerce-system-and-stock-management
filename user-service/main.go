package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	userDelivery "user-service/app/delivery"
	userRepo "user-service/app/repository"
	userUsecase "user-service/app/usecase"
)

var db *sql.DB
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func initDB() {
	var err error
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)
	fmt.Println(connStr)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")
}

// func healthHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("OK"))
// }

func main() {
	initDB()
	defer db.Close()
	r := gin.New()
	userRepository := userRepo.NewUserRepository(db)
	userUseCase := userUsecase.NewUserUsecase(userRepository, string(jwtSecret))
	userDelivery.NewAuthHandler(r, userUseCase)

	// http.HandleFunc("/health", healthHandler)
	// http.HandleFunc("/login", loginHandler)
	// http.HandleFunc("/register", registerHandler)

	fmt.Println("User service is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
