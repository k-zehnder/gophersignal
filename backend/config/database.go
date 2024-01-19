package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func NewConnection() *gorm.DB {
	sqlconnection := os.Getenv("MYSQL_DSN")
	if sqlconnection == "" {
		log.Fatal("MYSQL_DSN environment variable is not set.")
	}

	// Establish a connection to the MySQL server
	db, err := gorm.Open(mysql.Open(sqlconnection), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the MySQL server:", err)
	}

	// Create the database if it doesn't exist
	err = db.Exec("CREATE DATABASE IF NOT EXISTS gophersignal").Error
	if err != nil {
		log.Fatal("Failed to create the database:", err)
	}

	// Close the current connection and reconnect to the new database
	db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s/gophersignal", sqlconnection)), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the gophersignal database:", err)
	}

	// Create the 'articles' table if it doesn't exist
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS articles (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			link VARCHAR(512) NOT NULL,
			content TEXT,
			summary VARCHAR(2000),
			source VARCHAR(100) NOT NULL,
			is_on_homepage BOOLEAN,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);
	`).Error
	if err != nil {
		log.Fatal("Failed to create the 'articles' table:", err)
	}

	return db
}
