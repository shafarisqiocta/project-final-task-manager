package main

import (
	"log"
	"os"
	"project-final-task-manager/config"
	"project-final-task-manager/handlers"
	"project-final-task-manager/middleware"

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

	//protected routes (wajib jwt)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Categories
		protected.GET("/categories", handlers.GetAllCAtegory)
		protected.POST("/categories", handlers.CreateCategory)
		protected.PUT("/categories/:id", handlers.UpdateCategory)
		protected.DELETE("/categories/:id", handlers.DeleteCategory)

		// Projects
		protected.GET("/projects", handlers.GetAllProjects)
		protected.POST("/projects", handlers.CreateProject)
		protected.GET("/projects/:id", handlers.GetProjectByID)
		protected.PUT("/projects/:id", handlers.UpdateProject)
		protected.DELETE("/projects/:id", handlers.DeleteProject)
	}

	// Jalankan server
	port := os.Getenv("APP_PORT")
	r.Run(":" + port)

}
