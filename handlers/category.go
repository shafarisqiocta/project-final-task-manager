package handlers

import (
	"net/http"
	"project-final-task-manager/config"
	"project-final-task-manager/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// get all category
func GetAllCAtegory(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id,name,created_at,updated_at FROM categories")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "Gagal mengambil data kategori"})
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error:": "Gagal mengambil data kategori"})
			return
		}
		categories = append(categories, category)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data kategori",
		"data":    categories,
	})
}

// post categories
func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": "Input tidak Valid"})
		return
	}
	if strings.TrimSpace(category.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kategori wajib diisi"})
		return
	}

	var newCategory models.Category
	err := config.DB.QueryRow(
		`INSERT INTO categories (name) VALUES ($1) RETURNING id, name, created_at, updated_at`,
		category.Name,
	).Scan(&newCategory.ID, &newCategory.Name, &newCategory.CreatedAt, &newCategory.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "Nama kategori sudah ada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kategori"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil membuat kategori",
		"data":    newCategory,
	})
}

// put func category
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": "inputan tidak valid"})
		return
	}
	if strings.TrimSpace(category.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kategori wajib diisi"})
		return
	}
	var updatedCategory models.Category
	err := config.DB.QueryRow(
		`UPDATE categories SET name = $1, updated_at = NOW() 
		WHERE id = $2 
		RETURNING id, name, created_at, updated_at`,
		category.Name, id,
	).Scan(&updatedCategory.ID, &updatedCategory.Name, &updatedCategory.CreatedAt, &updatedCategory.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "Nama kategori sudah ada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengupdate kategori",
		"data":    updatedCategory,
	})
}

// func delete category
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	result, err := config.DB.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil menghapus kategori"})
}
