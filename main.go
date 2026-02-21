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
	if err := godotenv.Load(); err != nil {
		log.Println(".env file tidak ditemukan, menggunakan system environment variables")
	}
	config.ConnectDB()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Task Manager API is running "})
	})

	users := r.Group("/api/users")
	{
		users.POST("/register", handlers.Register) // jadi /api/users/register
		users.POST("/login", handlers.Login)       // jadi /api/users/login
	}

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

		// Tasks
		protected.GET("/projects/:id/tasks", handlers.GetAllTasks)
		protected.POST("/projects/:id/tasks", handlers.CreateTask)
		protected.GET("/projects/:id/tasks/:taskId", handlers.GetTaskByID)
		protected.PUT("/projects/:id/tasks/:taskId", handlers.UpdateTask)
		protected.DELETE("/projects/:id/tasks/:taskId", handlers.DeleteTask)
	}

	// Railway inject PORT, fallback ke APP_PORT kalau lokal
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT")
	}
	if port == "" {
		port = "8080" // default fallback
	}
	r.Run(":" + port)

}
