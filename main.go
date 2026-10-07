package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// global db var
var DB *gorm.DB

func initDB() {
	// env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file (might be using system env vars)")
	}

	// dburl
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in .env")
	}

	// connect
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = database
	log.Println("Database connection successfully opened!")
}

func main() {
	initDB()

	r := gin.Default()

	r.GET("/api/health", func(c *gin.Context) {
		// ping db
		sqlDB, err := DB.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Database connection failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Connected to Supabase successfully!",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on http://localhost:%s\n", port)
	r.Run(":" + port)
}
