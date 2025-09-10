package shared

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func InitDB() (*sql.DB, error) {

	config := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		config.User,
		config.Password,
		config.DBName,
		config.Host,
		config.Port,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to the database")
	return db, nil
}
