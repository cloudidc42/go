# โครงสร้างโปรเจค Go แบบสมบูรณ์ (Production-Ready)

## Overview

โปรเจคตัวอย่างนี้แสดงโครงสร้างที่ดีที่สุดสำหรับ Production Application ใน Go

---

## Project Structure

```
e-commerce-api/
├── cmd/
│   ├── api/                    # API Server entry point
│   │   └── main.go
│   ├── worker/                 # Background worker
│   │   └── main.go
│   └── migrate/                # Database migration tool
│       └── main.go
│
├── internal/                   # Private application code
│   ├── config/                 # Configuration
│   │   └── config.go
│   │
│   ├── domain/                 # Business domain models
│   │   ├── user/
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── product/
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   └── order/
│   │       ├── model.go
│   │       ├── repository.go
│   │       └── service.go
│   │
│   ├── handler/                # HTTP handlers
│   │   ├── user_handler.go
│   │   ├── product_handler.go
│   │   └── order_handler.go
│   │
│   ├── middleware/             # HTTP middleware
│   │   ├── auth.go
│   │   ├── logging.go
│   │   ├── cors.go
│   │   └── rate_limit.go
│   │
│   ├── repository/             # Database implementations
│   │   ├── postgres/
│   │   │   ├── user_repo.go
│   │   │   ├── product_repo.go
│   │   │   └── order_repo.go
│   │   └── redis/
│   │       └── cache_repo.go
│   │
│   └── util/                   # Utility functions
│       ├── logger/
│       │   └── logger.go
│       ├── validator/
│       │   └── validator.go
│       └── password/
│           └── password.go
│
├── pkg/                        # Public libraries
│   ├── http/                   # HTTP helpers
│   │   └── response.go
│   ├── jwt/                    # JWT utilities
│   │   └── jwt.go
│   └── pagination/             # Pagination helper
│       └── pagination.go
│
├── api/                        # API Specifications
│   ├── openapi/
│   │   └── swagger.yaml
│   └── proto/                  # gRPC proto files
│       └── service.proto
│
├── migrations/                 # Database migrations
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_products_table.up.sql
│   └── 000002_create_products_table.down.sql
│
├── scripts/                    # Build and deployment scripts
│   ├── build.sh
│   ├── deploy.sh
│   └── test.sh
│
├── tests/                      # Test files
│   ├── integration/
│   │   ├── user_test.go
│   │   └── product_test.go
│   └── e2e/
│       └── api_test.go
│
├── configs/                    # Configuration files
│   ├── config.yaml
│   ├── config.dev.yaml
│   └── config.prod.yaml
│
├── deployments/                # Deployment configurations
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── kubernetes/
│       ├── deployment.yaml
│       ├── service.yaml
│       ├── ingress.yaml
│       └── configmap.yaml
│
├── docs/                       # Documentation
│   ├── architecture.md
│   ├── api.md
│   └── deployment.md
│
├── .github/                    # GitHub specific files
│   └── workflows/
│       ├── ci.yml
│       └── cd.yml
│
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── Makefile                    # Build automation
├── README.md                   # Project README
├── .gitignore                  # Git ignore file
├── .env.example                # Environment variables example
└── LICENSE                     # License file
```

---

## File Examples

### 1. cmd/api/main.go

```go
package main

import (
    "log"
    "os"

    "e-commerce-api/internal/config"
    "e-commerce-api/internal/handler"
    "e-commerce-api/internal/repository/postgres"
    "e-commerce-api/internal/domain/user"
    "e-commerce-api/internal/util/logger"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize logger
    if err := logger.Init(cfg.LogLevel); err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }

    // Connect to database
    db, err := postgres.Connect(cfg.Database)
    if err != nil {
        logger.Fatal("Failed to connect to database", "error", err)
    }
    defer db.Close()

    // Initialize repositories
    userRepo := postgres.NewUserRepository(db)

    // Initialize services
    userService := user.NewService(userRepo)

    // Initialize handlers
    userHandler := handler.NewUserHandler(userService)

    // Setup router
    router := setupRouter(cfg, userHandler)

    // Start server
    logger.Info("Starting server", "port", cfg.Server.Port)
    if err := router.Run(cfg.Server.Address()); err != nil {
        logger.Fatal("Server failed", "error", err)
    }
}
```

### 2. internal/domain/user/model.go

```go
package user

import "time"

type User struct {
    ID        int       `json:"id" db:"id"`
    Email     string    `json:"email" db:"email"`
    Name      string    `json:"name" db:"name"`
    Password  string    `json:"-" db:"password"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Email string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
    Token string `json:"token"`
    User  *User  `json:"user"`
}
```

### 3. internal/domain/user/repository.go

```go
package user

import "context"

type Repository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id int) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int) error
    List(ctx context.Context, limit, offset int) ([]*User, error)
}
```

### 4. internal/domain/user/service.go

```go
package user

import (
    "context"
    "fmt"

    "e-commerce-api/internal/util/password"
)

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Check if user exists
    existing, _ := s.repo.GetByEmail(ctx, req.Email)
    if existing != nil {
        return nil, fmt.Errorf("email already exists")
    }

    // Hash password
    hashedPassword, err := password.Hash(req.Password)
    if err != nil {
        return nil, err
    }

    // Create user
    user := &User{
        Email:    req.Email,
        Name:     req.Name,
        Password: hashedPassword,
    }

    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*User, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int, req *UpdateUserRequest) (*User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    user.Name = req.Name
    user.Email = req.Email

    if err := s.repo.Update(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
    return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, page, pageSize int) ([]*User, error) {
    offset := (page - 1) * pageSize
    return s.repo.List(ctx, pageSize, offset)
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
    user, err := s.repo.GetByEmail(ctx, req.Email)
    if err != nil {
        return nil, fmt.Errorf("invalid credentials")
    }

    if !password.Check(req.Password, user.Password) {
        return nil, fmt.Errorf("invalid credentials")
    }

    // Generate JWT token
    token, err := generateJWT(user.ID)
    if err != nil {
        return nil, err
    }

    return &LoginResponse{
        Token: token,
        User:  user,
    }, nil
}
```

### 5. internal/handler/user_handler.go

```go
package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "e-commerce-api/internal/domain/user"
    "e-commerce-api/pkg/http/response"
)

type UserHandler struct {
    service *user.Service
}

func NewUserHandler(service *user.Service) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *gin.Context) {
    var req user.CreateUserRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request", err)
        return
    }

    user, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Failed to create user", err)
        return
    }

    response.Success(c, http.StatusCreated, "User created successfully", user)
}

func (h *UserHandler) Get(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
        return
    }

    user, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        response.Error(c, http.StatusNotFound, "User not found", err)
        return
    }

    response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

func (h *UserHandler) Update(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
        return
    }

    var req user.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request", err)
        return
    }

    user, err := h.service.Update(c.Request.Context(), id, &req)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Failed to update user", err)
        return
    }

    response.Success(c, http.StatusOK, "User updated successfully", user)
}

func (h *UserHandler) Delete(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
        return
    }

    if err := h.service.Delete(c.Request.Context(), id); err != nil {
        response.Error(c, http.StatusBadRequest, "Failed to delete user", err)
        return
    }

    response.Success(c, http.StatusOK, "User deleted successfully", nil)
}

func (h *UserHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

    users, err := h.service.List(c.Request.Context(), page, pageSize)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to list users", err)
        return
    }

    response.Success(c, http.StatusOK, "Users retrieved successfully", users)
}

func (h *UserHandler) Login(c *gin.Context) {
    var req user.LoginRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request", err)
        return
    }

    loginResp, err := h.service.Login(c.Request.Context(), &req)
    if err != nil {
        response.Error(c, http.StatusUnauthorized, "Login failed", err)
        return
    }

    response.Success(c, http.StatusOK, "Login successful", loginResp)
}
```

### 6. Makefile

```makefile
.PHONY: build test run clean docker-build docker-run migrate-up migrate-down

# Variables
APP_NAME=e-commerce-api
BUILD_DIR=bin
DOCKER_IMAGE=$(APP_NAME):latest

# Build
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) cmd/api/main.go

# Test
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

test-coverage:
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run
run:
	@echo "Running $(APP_NAME)..."
	@go run cmd/api/main.go

# Clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Docker
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) -f deployments/docker/Dockerfile .

docker-run:
	@echo "Running with Docker Compose..."
	@docker-compose -f deployments/docker/docker-compose.yml up -d

docker-stop:
	@docker-compose -f deployments/docker/docker-compose.yml down

# Database Migrations
migrate-up:
	@echo "Running migrations up..."
	@migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	@echo "Running migrations down..."
	@migrate -path migrations -database "$(DATABASE_URL)" down

migrate-create:
	@echo "Creating migration: $(name)"
	@migrate create -ext sql -dir migrations -seq $(name)

# Linting
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .

# Dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Generate
generate:
	@echo "Generating code..."
	@go generate ./...

# Install tools
install-tools:
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Swagger
swagger:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/api/main.go -o api/swagger

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  run            - Run the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  docker-stop    - Stop Docker Compose"
	@echo "  migrate-up     - Run database migrations up"
	@echo "  migrate-down   - Run database migrations down"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  deps           - Download dependencies"
	@echo "  swagger        - Generate Swagger documentation"
```

---

## Key Features

✅ **Clean Architecture** - Domain-driven design
✅ **Layered Structure** - Separation of concerns
✅ **Dependency Injection** - Testable code
✅ **Configuration Management** - YAML config files
✅ **Database Migrations** - Version control for schema
✅ **Docker Support** - Containerization ready
✅ **Kubernetes Ready** - K8s manifests included
✅ **CI/CD Pipeline** - GitHub Actions workflows
✅ **API Documentation** - OpenAPI/Swagger specs
✅ **Comprehensive Testing** - Unit, integration, e2e
✅ **Logging** - Structured logging
✅ **Monitoring** - Prometheus metrics
✅ **Health Checks** - Liveness and readiness probes
✅ **Graceful Shutdown** - Signal handling
✅ **Security** - JWT authentication, rate limiting

---

## Next Steps

1. Clone this structure for your project
2. Customize domain models for your business
3. Implement repository interfaces
4. Add business logic in services
5. Create HTTP handlers
6. Write tests
7. Deploy to production

Happy Coding! 🚀
