package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func OpenDB() (*gorm.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	time.Sleep(30 * time.Second)
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable host=postgres-database", dbUser, dbName, dbPassword)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// sqlDB, err := db.DB()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Set the maximum number of idle connections in the pool.
	// sqlDB.SetMaxIdleConns(10)

	// // Set the maximum number of open connections to the database.
	// sqlDB.SetMaxOpenConns(100)

	// // Check if the connection is successful
	// if err := sqlDB.Ping(); err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("Connected to the database")

	return db, err
}

func PingDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")
}
