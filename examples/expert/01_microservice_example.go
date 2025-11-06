// ===== ตัวอย่าง: Microservice แบบสมบูรณ์ =====
// User Service - Production Ready
package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"

    _ "github.com/lib/pq"
)

// ===== Models =====

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

// ===== Configuration =====

type Config struct {
    ServerPort      string
    DBConnectionStr string
    ShutdownTimeout time.Duration
}

func LoadConfig() *Config {
    return &Config{
        ServerPort:      getEnv("PORT", "8080"),
        DBConnectionStr: getEnv("DATABASE_URL", "postgres://user:password@localhost/mydb?sslmode=disable"),
        ShutdownTimeout: 30 * time.Second,
    }
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}

// ===== Database Repository =====

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *User) error {
    query := `
        INSERT INTO users (name, email, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
    now := time.Now()
    return r.db.QueryRow(query, user.Name, user.Email, now, now).Scan(&user.ID)
}

func (r *UserRepository) GetByID(id int) (*User, error) {
    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE id = $1
    `
    user := &User{}
    err := r.db.QueryRow(query, id).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }
    return user, nil
}

func (r *UserRepository) GetAll() ([]*User, error) {
    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        ORDER BY id
    `
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    users := []*User{}
    for rows.Next() {
        user := &User{}
        err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}

func (r *UserRepository) Update(user *User) error {
    query := `
        UPDATE users
        SET name = $1, email = $2, updated_at = $3
        WHERE id = $4
    `
    user.UpdatedAt = time.Now()
    _, err := r.db.Exec(query, user.Name, user.Email, user.UpdatedAt, user.ID)
    return err
}

func (r *UserRepository) Delete(id int) error {
    query := `DELETE FROM users WHERE id = $1`
    _, err := r.db.Exec(query, id)
    return err
}

// ===== Service Layer =====

type UserService struct {
    repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *User) error {
    // Validation
    if user.Name == "" || user.Email == "" {
        return fmt.Errorf("name and email are required")
    }

    return s.repo.Create(user)
}

func (s *UserService) GetUser(id int) (*User, error) {
    return s.repo.GetByID(id)
}

func (s *UserService) GetAllUsers() ([]*User, error) {
    return s.repo.GetAll()
}

func (s *UserService) UpdateUser(user *User) error {
    // Check if user exists
    _, err := s.repo.GetByID(user.ID)
    if err != nil {
        return err
    }

    return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id int) error {
    // Check if user exists
    _, err := s.repo.GetByID(id)
    if err != nil {
        return err
    }

    return s.repo.Delete(id)
}

// ===== HTTP Handlers =====

type UserHandler struct {
    service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
        return
    }

    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }

    if err := h.service.CreateUser(&user); err != nil {
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    respondSuccess(w, http.StatusCreated, "User created successfully", user)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
        return
    }

    users, err := h.service.GetAllUsers()
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondSuccess(w, http.StatusOK, "Users retrieved successfully", users)
}

// Helper functions
func respondSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func respondError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(Response{
        Success: false,
        Error:   message,
    })
}

// ===== Middleware =====

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
        next.ServeHTTP(w, r)
        log.Printf("Completed in %v", time.Since(start))
    })
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// ===== Health Check =====

func healthHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check database
        if err := db.Ping(); err != nil {
            respondError(w, http.StatusServiceUnavailable, "Database unhealthy")
            return
        }

        respondSuccess(w, http.StatusOK, "Service healthy", map[string]string{
            "status":    "healthy",
            "timestamp": time.Now().Format(time.RFC3339),
        })
    }
}

// ===== Application =====

type Application struct {
    config  *Config
    db      *sql.DB
    server  *http.Server
    handler *UserHandler
}

func NewApplication() (*Application, error) {
    config := LoadConfig()

    // Connect to database
    db, err := sql.Open("postgres", config.DBConnectionStr)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %v", err)
    }

    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %v", err)
    }

    // Initialize layers
    repo := NewUserRepository(db)
    service := NewUserService(repo)
    handler := NewUserHandler(service)

    // Setup routes
    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthHandler(db))
    mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            handler.Create(w, r)
        } else {
            handler.GetAll(w, r)
        }
    })

    // Apply middleware
    var h http.Handler = mux
    h = corsMiddleware(h)
    h = loggingMiddleware(h)

    // Create server
    server := &http.Server{
        Addr:         ":" + config.ServerPort,
        Handler:      h,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    return &Application{
        config:  config,
        db:      db,
        server:  server,
        handler: handler,
    }, nil
}

func (app *Application) Start() error {
    // Start server in goroutine
    go func() {
        log.Printf("Server starting on %s", app.server.Addr)
        if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), app.config.ShutdownTimeout)
    defer cancel()

    if err := app.server.Shutdown(ctx); err != nil {
        return fmt.Errorf("server forced to shutdown: %v", err)
    }

    // Close database
    if err := app.db.Close(); err != nil {
        return fmt.Errorf("failed to close database: %v", err)
    }

    log.Println("Server exited gracefully")
    return nil
}

// ===== Main =====

func main() {
    app, err := NewApplication()
    if err != nil {
        log.Fatalf("Failed to create application: %v", err)
    }

    if err := app.Start(); err != nil {
        log.Fatalf("Application error: %v", err)
    }
}

/*
===== Database Schema =====

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_users_email ON users(email);

===== Docker Compose =====

version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: mydb
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  user-service:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      DATABASE_URL: "postgres://user:password@postgres:5432/mydb?sslmode=disable"
      PORT: "8080"

volumes:
  postgres_data:

===== Dockerfile =====

FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o user-service

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/user-service .
EXPOSE 8080
CMD ["./user-service"]

===== Usage =====

# Run with Docker Compose
docker-compose up -d

# Test API
curl http://localhost:8080/health
curl http://localhost:8080/api/users
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

Features:
✅ Clean Architecture (Repository → Service → Handler)
✅ Database Connection Pool
✅ Graceful Shutdown
✅ Health Check Endpoint
✅ CORS Middleware
✅ Logging Middleware
✅ Error Handling
✅ JSON API
✅ Docker Support
✅ Production Ready
*/
