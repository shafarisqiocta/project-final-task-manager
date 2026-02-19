package handlers

import (
	"database/sql"
	"net/http"
	"project-final-task-manager/config"
	"project-final-task-manager/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// get all data project
func GetAllProjects(c *gin.Context) {
	userID := c.GetInt("user_id")

	rows, err := config.DB.Query(
		"SELECT id, user_id, name, description, created_at, updated_at FROM projects where user_id = $1", userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "gagal mengambil data project"})
		return
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(&project.ID, &project.UserID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data project"})
			return
		}
		projects = append(projects, project)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data project",
		"data":    projects,
	})
}

// POST /api/projects - buat project baru
func CreateProject(c *gin.Context) {
	userID := c.GetInt("user_id")

	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	if strings.TrimSpace(project.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama project wajib diisi"})
		return
	}

	var newProject models.Project
	err := config.DB.QueryRow(
		`INSERT INTO projects (user_id, name, description) 
		VALUES ($1, $2, $3) 
		RETURNING id, user_id, name, description, created_at, updated_at`,
		userID, project.Name, project.Description,
	).Scan(&newProject.ID, &newProject.UserID, &newProject.Name, &newProject.Description, &newProject.CreatedAt, &newProject.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil membuat project",
		"data":    newProject,
	})
}

// GET /api/projects/:id - ambil detail 1 project
func GetProjectByID(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID := c.Param("id")

	var project models.Project
	err := config.DB.QueryRow(
		"SELECT id, user_id, name, description, created_at, updated_at FROM projects WHERE id = $1",
		projectID,
	).Scan(&project.ID, &project.UserID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project tidak ditemukan"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data project"})
		return
	}

	// Cek kepemilikan project
	if project.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke project ini"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data project",
		"data":    project,
	})
}

// PUT /api/projects/:id - update project
func UpdateProject(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID := c.Param("id")

	// Cek kepemilikan project dulu
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

	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	if strings.TrimSpace(project.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama project wajib diisi"})
		return
	}

	var updatedProject models.Project
	err = config.DB.QueryRow(
		`UPDATE projects SET name = $1, description = $2, updated_at = NOW() 
		WHERE id = $3 
		RETURNING id, user_id, name, description, created_at, updated_at`,
		project.Name, project.Description, projectID,
	).Scan(&updatedProject.ID, &updatedProject.UserID, &updatedProject.Name, &updatedProject.Description, &updatedProject.CreatedAt, &updatedProject.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengupdate project",
		"data":    updatedProject,
	})
}

// DELETE /api/projects/:id - hapus project
func DeleteProject(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID := c.Param("id")

	// Cek kepemilikan project dulu
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

	_, err = config.DB.Exec("DELETE FROM projects WHERE id = $1", projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil menghapus project"})
}
