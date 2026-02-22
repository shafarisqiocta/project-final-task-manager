TaskFlow API - Task Management REST API

TaskFlow adalah REST API untuk manajemen tugas berbasis project yang dibangun menggunakan Golang. Aplikasi ini memungkinkan pengguna untuk mengorganisir pekerjaan mereka ke dalam project-project, lalu memecah setiap project menjadi task-task kecil yang bisa dilacak statusnya.

#Fitur Utama
✅ Autentikasi JWT - Register, login, dan proteksi endpoint dengan JSON Web Token
✅ Manajemen Project - Setiap user bisa membuat dan mengelola project miliknya sendiri
✅ Manajemen Task - Task dibuat di dalam project dengan status todo, in_progress, dan done
✅ Kategori Task - Task bisa dikategorikan (Feature, Bug Fix, Documentation, dll)
✅ Keamanan Data - User hanya bisa mengakses project dan task miliknya sendiri
✅ Auto Database Setup - Tabel otomatis dibuat saat pertama kali server jalan

#Tech Stack
Go 1.25.6 - Backend programming language
Gin - HTTP web framework yang ringan dan cepat
PostgreSQL - Relational database
JWT (golang-jwt/jwt) - Untuk autentikasi yang aman
Bcrypt - Hash password
godotenv - Load environment variables

#Database Schema (ERD)
Aplikasi ini menggunakan 4 tabel yang saling berelasi:
┌─────────────┐
│   users     │
├─────────────┤
│ id (PK)     │
│ name        │
│ email       │
│ password    │
│ created_at  │
│ updated_at  │
└──────┬──────┘
       │ 1
       │
       │ N
┌──────▼──────────┐
│   projects      │
├─────────────────┤
│ id (PK)         │
│ user_id (FK)    │
│ name            │
│ description     │
│ created_at      │
│ updated_at      │
└──────┬──────────┘
       │ 1
       │
       │ N
┌──────▼──────────┐       ┌──────────────┐
│     tasks       │       │ categories   │
├─────────────────┤       ├──────────────┤
│ id (PK)         │◄──N───│ id (PK)      │
│ project_id (FK) │       │ name         │
│ category_id (FK)│       │ created_at   │
│ title           │       │ updated_at   │
│ description     │       └──────────────┘
│ status          │
│ deadline        │
│ created_at      │
│ updated_at      │
└─────────────────┘

#Relasi Antar Tabel:
users → projects : One-to-Many (1 user bisa punya banyak project)
projects → tasks : One-to-Many (1 project bisa punya banyak task)
categories → tasks : One-to-Many (1 kategori bisa dipakai banyak task)

#Instalasi & Setup
Prerequisites:
- Go 1.25.6 atau lebih baru
- PostgreSQL 12+
- Git
1. Clone Repository
git clone https://github.com/shafarisqiocta/project-final-task-manager.git
cd project-final-task-manager
2. Install Dependencies
go mod tidy
3. Setup Database PostgreSQL
Buat database baru : CREATE DATABASE taskmanager_db;
4. Setup Environment Variables
Buat file .env di root project:
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=taskmanager_db
JWT_SECRET=your_secret_key_min_32_characters
APP_PORT=8080
Catatan: Jangan lupa ganti your_password dan your_secret_key dengan nilai yang sesuai!
5. Jalankan server
go run main.go
Jika berhasil, akan menampilkan output seperti ini:
✅ Database berhasil terkoneksi!
✅ Tabel telah terbuat / sudah ada!
[GIN-debug] Listening and serving HTTP on :8080
Tabel akan otomatis terbuat dengan 8 kategori default.

# API Documentation
Base URL:
- Local : http://localhost:8080
- production: http://project-final-task-manager-production.up.railway.app/

> Authentication Endpoints
# Register // membuat akun baru
#POST /api/users/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
Response (201 Created):
{
  "message": "Register berhasil",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2026-02-20T10:00:00Z"
  }
}

# Login 
#Login & mendapatkan JWT Token.
POST /api/users/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
response (200 OK):
{
  "message": "Login berhasil",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
} // Token JWT ini harus disertakan di header Authorization: Bearer <token> untuk semua endpoint protected!
-> Protected Endpoints (Require JWT): semua endpoint di bawah memerlukan header:
Authorization: Bearer <your_jwt_token>

# Category Endpoints
#GET ALL CATEGORIES : GET /api/categories
Response (200 OK):
{
  "message": "Berhasil mengambil data kategori",
  "data": [
    {
      "id": 1,
      "name": "Feature",
      "created_at": "2026-02-20T10:00:00Z",
      "updated_at": "2026-02-20T10:00:00Z"
    },
    {
      "id": 2,
      "name": "Bug Fix",
      "created_at": "2026-02-20T10:00:00Z",
      "updated_at": "2026-02-20T10:00:00Z"
    }
  ]
}
#POST CREATE CATEGORY
POST /api/categories
Content-Type: application/json

{
  "name": "Urgent"
}
Response (201 Created):
{
  "message": "Berhasil membuat kategori",
  "data": {
    "id": 9,
    "name": "Urgent",
    "created_at": "2026-02-20T10:30:00Z",
    "updated_at": "2026-02-20T10:30:00Z"
  }
}
#PUT UPDATE CATEGORY
PUT /api/categories/9
Content-Type: application/json

{
  "name": "Very Urgent"
}
Response (200 OK):
{
  "message": "Berhasil mengupdate kategori",
  "data": {
    "id": 9,
    "name": "Very Urgent",
    "created_at": "2026-02-20T10:30:00Z",
    "updated_at": "2026-02-20T11:00:00Z"
  }
}
#DELETE CATEGORY
DELETE /api/categories/9
Response (200 OK):
{
  "message": "Berhasil menghapus kategori"
}

# Project Endpoints
#GET ALL PROJECTS
GET /api/projects
Response (200 OK):
{
  "message": "Berhasil mengambil data project",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "name": "Bootcamp Final Project",
      "description": "REST API Task Manager",
      "created_at": "2026-02-20T10:00:00Z",
      "updated_at": "2026-02-20T10:00:00Z"
    }
  ]
}
#POST PROJECT
POST /api/projects
Content-Type: application/json

{
  "name": "Personal Website",
  "description": "Build my portfolio website"
}
Response (201 Created):
{
  "message": "Berhasil membuat project",
  "data": {
    "id": 2,
    "user_id": 1,
    "name": "Personal Website",
    "description": "Build my portfolio website",
    "created_at": "2026-02-20T11:00:00Z",
    "updated_at": "2026-02-20T11:00:00Z"
  }
}
#GET PROJECT BY ID
GET /api/projects/1
Response (200 OK):
{
  "message": "Berhasil mengambil data project",
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "Bootcamp Final Project",
    "description": "REST API Task Manager",
    "created_at": "2026-02-20T10:00:00Z",
    "updated_at": "2026-02-20T10:00:00Z"
  }
}
#PUT UPDATE PROJECT
PUT /api/projects/1
Content-Type: application/json

{
  "name": "Bootcamp Final Project (Updated)",
  "description": "REST API with JWT Authentication"
}
Response (200 OK):
{
  "message": "Berhasil mengupdate project",
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "Bootcamp Final Project (Updated)",
    "description": "REST API with JWT Authentication",
    "created_at": "2026-02-20T10:00:00Z",
    "updated_at": "2026-02-20T12:00:00Z"
  }
}
#DELETE PROJECT
DELETE /api/projects/1
Response (200 OK):
{
  "message": "Berhasil menghapus project"
}

# Task Endpoints
#GET ALL TASKS IN PROJECT
GET /api/projects/1/tasks
Response (200 OK):
{
  "message": "Berhasil mengambil data task",
  "data": [
    {
      "id": 1,
      "project_id": 1,
      "category_id": 1,
      "title": "Buat ERD Database",
      "description": "Desain skema tabel",
      "status": "todo",
      "deadline": "2026-02-25T00:00:00Z",
      "created_at": "2026-02-20T10:00:00Z",
      "updated_at": "2026-02-20T10:00:00Z"
    }
  ]
}
#POST CREATE TASK
POST /api/projects/1/tasks
Content-Type: application/json

{
  "title": "Setup Authentication",
  "description": "Implement JWT authentication",
  "category_id": 1,
  "status": "in_progress",
  "deadline": "2026-02-28T00:00:00Z"
}
Response (201 Created):
{
  "message": "Berhasil membuat task",
  "data": {
    "id": 2,
    "project_id": 1,
    "category_id": 1,
    "title": "Setup Authentication",
    "description": "Implement JWT authentication",
    "status": "in_progress",
    "deadline": "2026-02-28T00:00:00Z",
    "created_at": "2026-02-20T11:00:00Z",
    "updated_at": "2026-02-20T11:00:00Z"
  }
}
Status yang valid: todo, in_progress, done

#GET TASK BY ID
GET /api/projects/1/tasks/2
Response (200 OK):
{
  "message": "Berhasil mengambil data task",
  "data": {
    "id": 2,
    "project_id": 1,
    "category_id": 1,
    "title": "Setup Authentication",
    "description": "Implement JWT authentication",
    "status": "in_progress",
    "deadline": "2026-02-28T00:00:00Z",
    "created_at": "2026-02-20T11:00:00Z",
    "updated_at": "2026-02-20T11:00:00Z"
  }
}

#PUT UPDATE TASK
PUT /api/projects/1/tasks/2
Content-Type: application/json

{
  "title": "Setup Authentication",
  "description": "JWT authentication completed",
  "category_id": 1,
  "status": "done",
  "deadline": "2026-02-28T00:00:00Z"
}
Response (200 OK):
{
  "message": "Berhasil mengupdate task",
  "data": {
    "id": 2,
    "project_id": 1,
    "category_id": 1,
    "title": "Setup Authentication",
    "description": "JWT authentication completed",
    "status": "done",
    "deadline": "2026-02-28T00:00:00Z",
    "created_at": "2026-02-20T11:00:00Z",
    "updated_at": "2026-02-20T13:00:00Z"
  }
}

#DELETE TASK
DELETE /api/projects/1/tasks/2
Response (200 OK):
{
  "message": "Berhasil menghapus task"
}

# Security Features
1. Password Hashing

Password di-hash menggunakan bcrypt dengan cost 14
Password asli tidak pernah disimpan di database
Password hash tidak pernah muncul di response JSON

2. JWT Authentication

Token JWT expired dalam 24 jam
Setiap token berisi user_id dan email
Token di-sign dengan secret key yang aman
Middleware memvalidasi token di setiap request protected

3. Authorization & Access Control

User hanya bisa melihat dan mengelola project miliknya sendiri
User hanya bisa membuat/edit/hapus task di project miliknya
Validasi kepemilikan dilakukan di setiap endpoint protected
Response 403 Forbidden jika user coba akses data milik user lain

4. Input Validation

Validasi field wajib (name, email, password, title, dll)
Validasi format email
Validasi status task (hanya todo, in_progress, done)
Validasi foreign key (category_id, project_id harus ada)

# Error Handling
API ini menggunakan HTTP status code standar:
200 OK : Request Berhasil
201 Created : Resource Berhasil Terbuat
400 Bad Request : Input Tidak Valid
401 Unauthorized : Token Tidak Ada Atau Invalid
403 Forbidden : User Tidak Punya Akses Ke Resource
404 Not Found : Resource Tidak Ditemukan
409 COnflict: Data Duplikat (seperti email sudah terdaftar)
500 Internal Server Error : Error Di Server

# Deployment
Aplikasi ini sudah di-deploy ke Railway dan bisa diakses di:
http://project-final-task-manager-production.up.railway.app/ 

#nvironment Variables (Production)
Di Railway, set environment variables berikut:
DB_HOST=<railway_postgres_host>
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=<railway_postgres_password>
DB_NAME=railway
JWT_SECRET=<strong_random_secret>
PORT=8080
Railway akan otomatis inject PORT, dan database credentials akan auto-generate saat membuat PostgreSQL service.

# Project Structure
project-final-task-manager/
├── config/
│   └── database.go          # Koneksi database & auto-create tables
├── handlers/
│   ├── auth_handler.go      # Register & Login
│   ├── category_handler.go  # CRUD Categories
│   ├── project_handler.go   # CRUD Projects
│   └── task_handler.go      # CRUD Tasks
├── helpers/
│   ├── bcrypt.go            # Hash & verify password
│   └── jwt.go               # Generate & validate JWT
├── middleware/
│   └── auth_middleware.go   # JWT authentication middleware
├── models/
│   ├── user.go
│   ├── category.go
│   ├── project.go
│   └── task.go
├── .env                     # Environment variables (local)
├── .gitignore
├── go.mod
├── go.sum
├── main.go                  # Entry point & routing
└── README.md

# Notes
- Database schema akan otomatis dibuat saat pertama kali server dijalankan
- Data 8 kategori akan otomatis dimasukkan ke database
- JWT token berlaku selama 24 jam, setelah itu user harus login ulang
- Format deadline menggunakan ISO 8601: YYYY-MM-DDTHH:MM:SSZ
- Password minimal 6 karakter (bisa disesuaikan di handler)




