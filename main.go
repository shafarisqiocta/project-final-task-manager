package main

import (
	"log"
	"os"
	"project-final-task-manager/config"
	"project-final-task-manager/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal load .env:", err)
	}
	config.ConnectDB()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Task Manager API is running "})
	})

	// Public routes (tanpa JWT)
	r.POST("/api/register", handlers.Register)
	r.POST("/api/login", handlers.Login)

	// Jalankan server
	port := os.Getenv("APP_PORT")
	r.Run(":" + port)

}
