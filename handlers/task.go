package handlers

import (
	"database/sql"
	"net/http"
	"project-final-task-manager/config"
	"project-final-task-manager/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// get all task
func GetAllTasks(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID := c.Param("id")

	//validasi kepemilikan project
	var ownerID int
	err := config.DB.QueryRow("SELECT user_id FROM projects WHERE ID =$1", projectID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error:": "Project tidak ditemukan"})
		return
	}
	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error:": "Anda tidak memiliki akses ke project ini!"})
		return
	}
	rows, err := config.DB.Query(
		`SELECT id, project_id, category_id, title, description, status, deadline, created_at, updated_at 
		FROM tasks WHERE project_id = $1`,
		projectID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "Gagal mengambil data task"})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.ProjectID, &task.CategoryID, &task.Title, &task.Description, &task.Status, &task.Deadline, &task.CreatedAt, &task.UpdatedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error:": "Gagal membaca data task"})
			return
		}
		tasks = append(tasks, task)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data task",
		"data":    tasks,
	})
}

// req post tambah task
func CreateTask(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID := c.Param("id")

	//validasi kepemilikan project
	var ownerID int
	err := config.DB.QueryRow("SELECT user_id FROM projects WHERE ID =$1", projectID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error:": "Project tidak ditemukan"})
		return
	}
	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error:": "Anda tidak memiliki akses ke project ini!"})
		return
	}
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "Input tidak valid"})
		return
	}
	if strings.TrimSpace(task.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title task wajib diisi"})
		return
	}
	//validasi status
	validStatuses := []string{"todo", "in_progress", "done"}
	if task.Status == "" {
		task.Status = "todo"
	} else {
		isValid := false
		for _, status := range validStatuses {
			if task.Status == status {
				isValid = true
				break
			}
		}
		if !isValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus salah satu dari: todo, in_progress, done"})
			return
		}
	}
	// Validasi category_id ada atau tidak
	if task.CategoryID != 0 {
		var categoryExists int
		err = config.DB.QueryRow("SELECT id FROM categories WHERE id = $1", task.CategoryID).Scan(&categoryExists)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category tidak ditemukan"})
			return
		}
	}

	var newTask models.Task
	err = config.DB.QueryRow(
		`INSERT INTO tasks (project_id, category_id, title, description, status, deadline) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, project_id, category_id, title, description, status, deadline, created_at, updated_at`,
		projectID, task.CategoryID, task.Title, task.Description, task.Status, task.Deadline,
	).Scan(&newTask.ID, &newTask.ProjectID, &newTask.CategoryID, &newTask.Title, &newTask.Description, &newTask.Status, &newTask.Deadline, &newTask.CreatedAt, &newTask.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil membuat task",
		"data":    newTask,
	})
}

// get detail 1 task
func GetTaskByID(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID := c.Param("id")
	taskID := c.Param("taskId")

	// Validasi kepemilikan project
	var ownerID int
	err := config.DB.QueryRow("SELECT user_id FROM projects WHERE id = $1", projectID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project tidak ditemukan"})
		return
	}
	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke project ini"})
		return
	}

	var task models.Task
	err = config.DB.QueryRow(
		`SELECT id, project_id, category_id, title, description, status, deadline, created_at, updated_at 
		FROM tasks WHERE id = $1 AND project_id = $2`,
		taskID, projectID,
	).Scan(&task.ID, &task.ProjectID, &task.CategoryID, &task.Title, &task.Description, &task.Status, &task.Deadline, &task.CreatedAt, &task.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data task",
		"data":    task,
	})
}

// PUT /api/projects/:id/tasks/:taskId - update task
func UpdateTask(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID := c.Param("id")
	taskID := c.Param("taskId")

	// Validasi kepemilikan project
	var ownerID int
	err := config.DB.QueryRow("SELECT user_id FROM projects WHERE id = $1", projectID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project tidak ditemukan"})
		return
	}
	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke project ini"})
		return
	}

	// Cek task ada atau tidak
	var existingTaskID int
	err = config.DB.QueryRow("SELECT id FROM tasks WHERE id = $1 AND project_id = $2", taskID, projectID).Scan(&existingTaskID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title task wajib diisi"})
		return
	}

	// Validasi status
	validStatuses := []string{"todo", "in_progress", "done"}
	if task.Status != "" {
		isValid := false
		for _, status := range validStatuses {
			if task.Status == status {
				isValid = true
				break
			}
		}
		if !isValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus salah satu dari: todo, in_progress, done"})
			return
		}
	}

	// Validasi category_id ada atau tidak
	if task.CategoryID != 0 {
		var categoryExists int
		err = config.DB.QueryRow("SELECT id FROM categories WHERE id = $1", task.CategoryID).Scan(&categoryExists)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category tidak ditemukan"})
			return
		}
	}

	var updatedTask models.Task
	err = config.DB.QueryRow(
		`UPDATE tasks SET category_id = $1, title = $2, description = $3, status = $4, deadline = $5, updated_at = NOW() 
		WHERE id = $6 AND project_id = $7 
		RETURNING id, project_id, category_id, title, description, status, deadline, created_at, updated_at`,
		task.CategoryID, task.Title, task.Description, task.Status, task.Deadline, taskID, projectID,
	).Scan(&updatedTask.ID, &updatedTask.ProjectID, &updatedTask.CategoryID, &updatedTask.Title, &updatedTask.Description, &updatedTask.Status, &updatedTask.Deadline, &updatedTask.CreatedAt, &updatedTask.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengupdate task",
		"data":    updatedTask,
	})
}

// DELETE /api/projects/:id/tasks/:taskId - hapus task
func DeleteTask(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID := c.Param("id")
	taskID := c.Param("taskId")

	// Validasi kepemilikan project
	var ownerID int
	err := config.DB.QueryRow("SELECT user_id FROM projects WHERE id = $1", projectID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project tidak ditemukan"})
		return
	}
	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke project ini"})
		return
	}

	result, err := config.DB.Exec("DELETE FROM tasks WHERE id = $1 AND project_id = $2", taskID, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus task"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil menghapus task"})
}
