# MBKM Go API

Go Fiber + GORM backend untuk aplikasi MBKM (Merdeka Belajar Kampus Merdeka).

## Prerequisites

- Go 1.21+
- MySQL/MariaDB

## Setup

1. Copy environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` dengan konfigurasi database Anda

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run server:
   ```bash
   go run cmd/server/main.go
   ```

Server akan berjalan di `http://localhost:3000`

## API Endpoints

### Public
- `GET /api/v1/test` - Health check
- `POST /api/v1/register` - User registration
- `POST /api/v1/login` - User login
- `GET /api/v1/public/jobs` - List jobs
- `GET /api/v1/public/jobs/:id` - Job detail
- `GET /api/v1/articles` - List articles
- `GET /api/v1/articles/:id` - Article detail

### Protected (requires JWT token)
- `GET /api/v1/logout` - Logout
- `GET /api/v1/profile` - User profile
- `POST /api/v1/jobs` - Create job
- `PUT /api/v1/jobs/:id` - Update job
- `DELETE /api/v1/jobs/:id` - Delete job
- `POST /api/v1/jobs/:id/approve` - Approve job
- `POST /api/v1/jobs/:id/reject` - Reject job
- `POST /api/v1/jobs/:id/close` - Close job
- Companies CRUD: `/api/v1/companies`

## Authentication

Gunakan header `Authorization: Bearer <token>` untuk endpoint protected.

## Project Structure

```
go-api/
├── cmd/server/main.go       # Entry point
├── config/config.go         # Configuration
├── database/database.go     # Database connection
├── internal/
│   ├── dto/                 # Request/Response DTOs
│   ├── handlers/            # HTTP handlers
│   ├── middleware/          # Auth, CORS middleware
│   ├── models/              # GORM models
│   └── routes/              # Route definitions
└── pkg/utils/               # Utility functions
```
