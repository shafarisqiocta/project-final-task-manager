package handlers

import (
	"database/sql"
	"net/http"
	"project-final-task-manager/config"
	"project-final-task-manager/helpers"
	"project-final-task-manager/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// Register
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "Input tidak valid"})
		return
	}
	//validasi input tidak boleh kosong
	if strings.TrimSpace(user.Name) == "" || strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email, password wajib diisi"})
		return
	}
	//validasi email sudah terdaftar atau belum
	var existingID int
	err := config.DB.QueryRow("SELECT id FROM users WHERE email =$1", user.Email).Scan(&existingID)
	if err != sql.ErrNoRows {
		c.JSON(http.StatusConflict, gin.H{"error:": "email sudah terdaftar"})
		return
	}
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "gagal memproses password"})
		return
	}
	//simpan ke database
	var newUser models.User
	err = config.DB.QueryRow(
		`INSERT INTO users (name, email, password) 
		VALUES ($1, $2, $3) 
		RETURNING id, name, email, created_at`,
		user.Name, user.Email, hashedPassword,
	).Scan(&newUser.ID, &newUser.Name, &newUser.Email, &newUser.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "Gagal menyimpan user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Register berhasil",
		"user":    newUser,
	})
}

// Login
func Login(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "Input tidak valid"})
		return
	}
	//validasi inputan tidak boleh kosong
	if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Password) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email dan password wajib diisi"})
		return
	}
	// Cari user berdasarkan email
	var user models.User
	err := config.DB.QueryRow(
		"SELECT id, name, email, password FROM users WHERE email = $1",
		input.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}

	// Cek password
	if !helpers.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// Generate JWT token
	token, err := helpers.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
