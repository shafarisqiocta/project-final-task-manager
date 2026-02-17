package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Gagal tidak merespon:", err)
	}
	log.Println("Database berhasil terkoneksi")
	DB = db
	createTables()
}
func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id         SERIAL PRIMARY KEY,
			name       VARCHAR(100) NOT NULL,
			email      VARCHAR(100) UNIQUE NOT NULL,
			password   VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS categories (
			id         SERIAL PRIMARY KEY,
			name       VARCHAR(100) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS projects (
			id          SERIAL PRIMARY KEY,
			user_id     INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name        VARCHAR(100) NOT NULL,
			description TEXT,
			created_at  TIMESTAMP DEFAULT NOW(),
			updated_at  TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id          SERIAL PRIMARY KEY,
			project_id  INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			category_id INT REFERENCES categories(id) ON DELETE SET NULL,
			title       VARCHAR(100) NOT NULL,
			description TEXT,
			status      VARCHAR(20) DEFAULT 'todo',
			deadline    DATE,
			created_at  TIMESTAMP DEFAULT NOW(),
			updated_at  TIMESTAMP DEFAULT NOW()
		)`,
		`INSERT INTO categories (name) VALUES
			('Feature'),
			('Bug Fix'),
			('Documentation'),
			('Research'),
			('Design'),
			('Testing'),
			('Deployment'),
			('Meeting')
		ON CONFLICT (name) DO NOTHING`,
	}
	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatal("Gagal membuat table:", err)
		}
	}
	log.Println("Tabel telah terbuat / sudah ada")
}
