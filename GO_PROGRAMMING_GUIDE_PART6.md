# คู่มือการเขียนโปรแกรม Go - ส่วนที่ 6
## Step 801-1000: Microservices, gRPC, Docker และ Production Deployment

---

## Step 801-850: Microservices Architecture

### Step 801-810: Microservices Fundamentals

```go
// ===== Step 801: Service Structure =====
package main

/*
Microservices Project Structure:

microservices/
├── api-gateway/          # API Gateway service
│   ├── main.go
│   ├── routes.go
│   └── middleware.go
├── user-service/         # User management service
│   ├── main.go
│   ├── handler.go
│   ├── repository.go
│   └── models.go
├── product-service/      # Product management service
│   ├── main.go
│   ├── handler.go
│   └── repository.go
├── order-service/        # Order processing service
│   ├── main.go
│   ├── handler.go
│   └── repository.go
├── shared/               # Shared libraries
│   ├── logger/
│   ├── config/
│   └── middleware/
└── docker-compose.yml    # Docker orchestration
*/

// ===== Step 802-805: Service Discovery =====

type ServiceRegistry struct {
    mu       sync.RWMutex
    services map[string][]string
}

func NewServiceRegistry() *ServiceRegistry {
    return &ServiceRegistry{
        services: make(map[string][]string),
    }
}

func (sr *ServiceRegistry) Register(serviceName, address string) {
    sr.mu.Lock()
    defer sr.mu.Unlock()

    sr.services[serviceName] = append(sr.services[serviceName], address)
}

func (sr *ServiceRegistry) Deregister(serviceName, address string) {
    sr.mu.Lock()
    defer sr.mu.Unlock()

    addresses := sr.services[serviceName]
    for i, addr := range addresses {
        if addr == address {
            sr.services[serviceName] = append(addresses[:i], addresses[i+1:]...)
            break
        }
    }
}

func (sr *ServiceRegistry) GetService(serviceName string) (string, error) {
    sr.mu.RLock()
    defer sr.mu.RUnlock()

    addresses, exists := sr.services[serviceName]
    if !exists || len(addresses) == 0 {
        return "", fmt.Errorf("service not found: %s", serviceName)
    }

    // Simple round-robin
    return addresses[rand.Intn(len(addresses))], nil
}

// ===== Step 806-810: API Gateway =====

type APIGateway struct {
    registry *ServiceRegistry
    client   *http.Client
}

func NewAPIGateway(registry *ServiceRegistry) *APIGateway {
    return &APIGateway{
        registry: registry,
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (gw *APIGateway) ProxyRequest(serviceName string, w http.ResponseWriter, r *http.Request) {
    // Get service address
    serviceAddr, err := gw.registry.GetService(serviceName)
    if err != nil {
        http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
        return
    }

    // Create new request
    targetURL := fmt.Sprintf("http://%s%s", serviceAddr, r.URL.Path)
    proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
    if err != nil {
        http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
        return
    }

    // Copy headers
    proxyReq.Header = r.Header

    // Execute request
    resp, err := gw.client.Do(proxyReq)
    if err != nil {
        http.Error(w, "Service error", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    // Copy response
    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }

    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}
```

### Step 811-830: gRPC

```go
// ===== Step 811-815: Protocol Buffers Definition =====

/*
// user.proto
syntax = "proto3";

package user;

option go_package = "./proto";

service UserService {
  rpc CreateUser(CreateUserRequest) returns (User);
  rpc GetUser(GetUserRequest) returns (User);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  rpc UpdateUser(UpdateUserRequest) returns (User);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
}

message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
  string created_at = 4;
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
}

message GetUserRequest {
  int32 id = 1;
}

message ListUsersRequest {
  int32 page = 1;
  int32 page_size = 2;
}

message ListUsersResponse {
  repeated User users = 1;
  int32 total = 2;
}

message UpdateUserRequest {
  int32 id = 1;
  string name = 2;
  string email = 3;
}

message DeleteUserRequest {
  int32 id = 1;
}

message DeleteUserResponse {
  bool success = 1;
}

Generate code:
protoc --go_out=. --go-grpc_out=. user.proto
*/

// ===== Step 816-820: gRPC Server Implementation =====

package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "myapp/proto"
)

type UserServer struct {
    pb.UnimplementedUserServiceServer
    users map[int32]*pb.User
    mu    sync.RWMutex
}

func NewUserServer() *UserServer {
    return &UserServer{
        users: make(map[int32]*pb.User),
    }
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    id := int32(len(s.users) + 1)
    user := &pb.User{
        Id:        id,
        Name:      req.Name,
        Email:     req.Email,
        CreatedAt: time.Now().Format(time.RFC3339),
    }

    s.users[id] = user
    return user, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    user, exists := s.users[req.Id]
    if !exists {
        return nil, fmt.Errorf("user not found")
    }

    return user, nil
}

func (s *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var users []*pb.User
    for _, user := range s.users {
        users = append(users, user)
    }

    return &pb.ListUsersResponse{
        Users: users,
        Total: int32(len(users)),
    }, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterUserServiceServer(grpcServer, NewUserServer())

    log.Println("gRPC server listening on :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

// ===== Step 821-825: gRPC Client =====

func createGRPCClient() {
    // Connect to server
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewUserServiceClient(conn)

    // Create user
    ctx := context.Background()
    user, err := client.CreateUser(ctx, &pb.CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    })
    if err != nil {
        log.Fatalf("CreateUser failed: %v", err)
    }

    log.Printf("Created user: %v", user)

    // Get user
    getUser, err := client.GetUser(ctx, &pb.GetUserRequest{
        Id: user.Id,
    })
    if err != nil {
        log.Fatalf("GetUser failed: %v", err)
    }

    log.Printf("Retrieved user: %v", getUser)
}

// ===== Step 826-830: gRPC Streaming =====

// Server streaming
func (s *UserServer) StreamUsers(req *pb.StreamUsersRequest, stream pb.UserService_StreamUsersServer) error {
    s.mu.RLock()
    defer s.mu.RUnlock()

    for _, user := range s.users {
        if err := stream.Send(user); err != nil {
            return err
        }
        time.Sleep(100 * time.Millisecond) // Simulate delay
    }

    return nil
}

// Client streaming
func (s *UserServer) CreateMultipleUsers(stream pb.UserService_CreateMultipleUsersServer) error {
    count := 0

    for {
        req, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&pb.CreateMultipleUsersResponse{
                Count: int32(count),
            })
        }
        if err != nil {
            return err
        }

        // Create user
        user := &pb.User{
            Id:    int32(len(s.users) + 1),
            Name:  req.Name,
            Email: req.Email,
        }
        s.users[user.Id] = user
        count++
    }
}

// Bidirectional streaming
func (s *UserServer) Chat(stream pb.UserService_ChatServer) error {
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }

        // Process and send response
        response := &pb.ChatMessage{
            Content: "Echo: " + msg.Content,
        }

        if err := stream.Send(response); err != nil {
            return err
        }
    }
}
```

### Step 831-850: Docker และ Kubernetes

```dockerfile
# ===== Step 831-835: Dockerfile =====

# Multi-stage build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run
CMD ["./main"]

# ===== Step 836-840: Docker Compose =====

# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: myuser
      POSTGRES_PASSWORD: mypassword
      POSTGRES_DB: mydb
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  api-gateway:
    build:
      context: ./api-gateway
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    depends_on:
      - user-service
      - product-service
    environment:
      - USER_SERVICE_URL=user-service:8081
      - PRODUCT_SERVICE_URL=product-service:8082

  user-service:
    build:
      context: ./user-service
      dockerfile: Dockerfile
    ports:
      - "8081:8081"
    depends_on:
      - postgres
      - redis
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=myuser
      - DB_PASSWORD=mypassword
      - DB_NAME=mydb
      - REDIS_URL=redis:6379

  product-service:
    build:
      context: ./product-service
      dockerfile: Dockerfile
    ports:
      - "8082:8082"
    depends_on:
      - postgres
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432

volumes:
  postgres_data:

# ===== Step 841-845: Kubernetes Deployment =====

# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  labels:
    app: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: myregistry/user-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db.host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
        resources:
          requests:
            memory: "128Mi"
            cpu: "250m"
          limits:
            memory: "256Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5

---
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: user-service
spec:
  selector:
    app: user-service
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer

---
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  db.host: "postgres.default.svc.cluster.local"
  db.port: "5432"
  db.name: "mydb"

---
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-secret
type: Opaque
data:
  password: cGFzc3dvcmQxMjM=  # base64 encoded

# ===== Step 846-850: Horizontal Pod Autoscaler =====

# hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: user-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: user-service
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

---

## Step 851-900: Monitoring และ Logging

### Step 851-870: Prometheus Metrics

```go
// ===== Step 851-855: Metrics Collection =====

package monitoring

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // Counter
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    // Histogram
    httpDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_duration_seconds",
            Help:    "Duration of HTTP requests",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )

    // Gauge
    activeConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
            Help: "Number of active connections",
        },
    )

    // Summary
    requestSize = promauto.NewSummary(
        prometheus.SummaryOpts{
            Name: "request_size_bytes",
            Help: "Request size in bytes",
        },
    )
)

// ===== Step 856-860: Metrics Middleware =====

func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Wrap response writer to capture status code
        ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

        // Process request
        next.ServeHTTP(ww, r)

        // Record metrics
        duration := time.Since(start).Seconds()
        status := fmt.Sprintf("%d", ww.statusCode)

        httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
        httpDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
        requestSize.Observe(float64(r.ContentLength))
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// ===== Step 861-865: Custom Metrics =====

type BusinessMetrics struct {
    ordersProcessed prometheus.Counter
    orderValue      prometheus.Histogram
    activeUsers     prometheus.Gauge
}

func NewBusinessMetrics() *BusinessMetrics {
    return &BusinessMetrics{
        ordersProcessed: promauto.NewCounter(prometheus.CounterOpts{
            Name: "orders_processed_total",
            Help: "Total number of orders processed",
        }),
        orderValue: promauto.NewHistogram(prometheus.HistogramOpts{
            Name:    "order_value_dollars",
            Help:    "Order value in dollars",
            Buckets: []float64{10, 50, 100, 500, 1000, 5000},
        }),
        activeUsers: promauto.NewGauge(prometheus.GaugeOpts{
            Name: "active_users",
            Help: "Number of active users",
        }),
    }
}

func (bm *BusinessMetrics) RecordOrder(value float64) {
    bm.ordersProcessed.Inc()
    bm.orderValue.Observe(value)
}

func (bm *BusinessMetrics) SetActiveUsers(count float64) {
    bm.activeUsers.Set(count)
}

// ===== Step 866-870: Prometheus Server =====

func main() {
    // Regular HTTP routes
    http.HandleFunc("/api/users", usersHandler)

    // Prometheus metrics endpoint
    http.Handle("/metrics", promhttp.Handler())

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Step 871-890: Structured Logging

```go
// ===== Step 871-875: Logger Implementation =====

package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(environment string) error {
    var config zap.Config

    if environment == "production" {
        config = zap.NewProductionConfig()
    } else {
        config = zap.NewDevelopmentConfig()
    }

    // Customize config
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    var err error
    log, err = config.Build()
    if err != nil {
        return err
    }

    return nil
}

func Info(msg string, fields ...zap.Field) {
    log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
    log.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
    log.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
    log.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
    log.Fatal(msg, fields...)
}

// ===== Step 876-880: Request Logging Middleware =====

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Generate request ID
        requestID := generateRequestID()

        // Create logger with request context
        reqLogger := log.With(
            zap.String("request_id", requestID),
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
            zap.String("remote_addr", r.RemoteAddr),
        )

        // Add request ID to header
        w.Header().Set("X-Request-ID", requestID)

        // Log request
        reqLogger.Info("Request started")

        // Wrap response writer
        ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

        // Process request
        next.ServeHTTP(ww, r)

        // Log response
        duration := time.Since(start)
        reqLogger.Info("Request completed",
            zap.Int("status", ww.statusCode),
            zap.Duration("duration", duration),
        )
    })
}

// ===== Step 881-885: Application Logging =====

func processOrder(orderID string, amount float64) error {
    logger.Info("Processing order",
        zap.String("order_id", orderID),
        zap.Float64("amount", amount),
    )

    // Business logic
    if amount < 0 {
        logger.Error("Invalid order amount",
            zap.String("order_id", orderID),
            zap.Float64("amount", amount),
        )
        return fmt.Errorf("invalid amount")
    }

    // Success
    logger.Info("Order processed successfully",
        zap.String("order_id", orderID),
        zap.Float64("amount", amount),
    )

    return nil
}

// ===== Step 886-890: Log Aggregation =====

/*
Configure log shipping to ELK Stack or similar:

1. Filebeat configuration (filebeat.yml):

filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/myapp/*.log
  json.keys_under_root: true
  json.add_error_key: true

output.elasticsearch:
  hosts: ["localhost:9200"]
  index: "myapp-logs-%{+yyyy.MM.dd}"

setup.kibana:
  host: "localhost:5601"

2. Fluentd configuration (fluent.conf):

<source>
  @type tail
  path /var/log/myapp/*.log
  pos_file /var/log/td-agent/myapp.log.pos
  tag myapp.logs
  format json
</source>

<match myapp.logs>
  @type elasticsearch
  host localhost
  port 9200
  index_name myapp-logs
  type_name log
</match>
*/
```

---

## Step 891-950: Production Best Practices

### Step 891-910: Configuration Management

```go
// ===== Step 891-895: Configuration Structure =====

package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
    Logging  LoggingConfig
}

type ServerConfig struct {
    Host            string
    Port            int
    ReadTimeout     time.Duration
    WriteTimeout    time.Duration
    ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
    Host            string
    Port            int
    User            string
    Password        string
    Database        string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
}

type RedisConfig struct {
    Host     string
    Port     int
    Password string
    DB       int
}

type JWTConfig struct {
    Secret     string
    Expiration time.Duration
}

type LoggingConfig struct {
    Level       string
    Format      string
    OutputPaths []string
}

// ===== Step 896-900: Load Configuration =====

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")
    viper.AddConfigPath("/etc/myapp/")

    // Environment variables
    viper.AutomaticEnv()
    viper.SetEnvPrefix("MYAPP")

    // Defaults
    setDefaults()

    // Read config file
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            // Config file not found; using defaults and env vars
        } else {
            return nil, err
        }
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

func setDefaults() {
    viper.SetDefault("server.host", "0.0.0.0")
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.read_timeout", "10s")
    viper.SetDefault("server.write_timeout", "10s")
    viper.SetDefault("server.shutdown_timeout", "30s")

    viper.SetDefault("database.host", "localhost")
    viper.SetDefault("database.port", 5432)
    viper.SetDefault("database.max_open_conns", 25)
    viper.SetDefault("database.max_idle_conns", 5)
    viper.SetDefault("database.conn_max_lifetime", "5m")

    viper.SetDefault("logging.level", "info")
    viper.SetDefault("logging.format", "json")
}

// config.yaml example:
/*
server:
  host: 0.0.0.0
  port: 8080
  read_timeout: 10s
  write_timeout: 10s
  shutdown_timeout: 30s

database:
  host: localhost
  port: 5432
  user: myuser
  password: mypassword
  database: mydb
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key"
  expiration: 24h

logging:
  level: info
  format: json
  output_paths:
    - stdout
    - /var/log/myapp/app.log
*/

// ===== Step 901-905: Health Checks =====

type HealthChecker struct {
    db    *sql.DB
    redis *redis.Client
}

func NewHealthChecker(db *sql.DB, redis *redis.Client) *HealthChecker {
    return &HealthChecker{
        db:    db,
        redis: redis,
    }
}

func (hc *HealthChecker) Check() map[string]string {
    health := make(map[string]string)

    // Check database
    if err := hc.db.Ping(); err != nil {
        health["database"] = "unhealthy: " + err.Error()
    } else {
        health["database"] = "healthy"
    }

    // Check Redis
    if err := hc.redis.Ping(context.Background()).Err(); err != nil {
        health["redis"] = "unhealthy: " + err.Error()
    } else {
        health["redis"] = "healthy"
    }

    return health
}

func healthHandler(hc *HealthChecker) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        health := hc.Check()

        allHealthy := true
        for _, status := range health {
            if !strings.HasPrefix(status, "healthy") {
                allHealthy = false
                break
            }
        }

        statusCode := http.StatusOK
        if !allHealthy {
            statusCode = http.StatusServiceUnavailable
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(statusCode)
        json.NewEncoder(w).Encode(health)
    }
}

// ===== Step 906-910: Graceful Shutdown =====

func RunServer(config *Config, handler http.Handler) error {
    srv := &http.Server{
        Addr:         fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
        Handler:      handler,
        ReadTimeout:  config.Server.ReadTimeout,
        WriteTimeout: config.Server.WriteTimeout,
    }

    // Start server
    go func() {
        logger.Info("Server starting",
            zap.String("addr", srv.Addr),
        )

        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal("Server failed to start", zap.Error(err))
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit

    logger.Info("Server shutting down...")

    // Shutdown with timeout
    ctx, cancel := context.WithTimeout(
        context.Background(),
        config.Server.ShutdownTimeout,
    )
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Error("Server forced to shutdown", zap.Error(err))
        return err
    }

    logger.Info("Server exited gracefully")
    return nil
}
```

---

## Step 951-1000: Advanced Topics

### Step 951-970: Message Queues (RabbitMQ/Kafka)

```go
// ===== Step 951-955: RabbitMQ Producer =====

package messaging

import (
    "github.com/streadway/amqp"
)

type RabbitMQProducer struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

func NewRabbitMQProducer(url string) (*RabbitMQProducer, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }

    return &RabbitMQProducer{
        conn:    conn,
        channel: ch,
    }, nil
}

func (p *RabbitMQProducer) Publish(queueName string, message []byte) error {
    // Declare queue
    _, err := p.channel.QueueDeclare(
        queueName,
        true,  // durable
        false, // delete when unused
        false, // exclusive
        false, // no-wait
        nil,   // arguments
    )
    if err != nil {
        return err
    }

    // Publish message
    err = p.channel.Publish(
        "",        // exchange
        queueName, // routing key
        false,     // mandatory
        false,     // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        message,
        },
    )

    return err
}

func (p *RabbitMQProducer) Close() {
    p.channel.Close()
    p.conn.Close()
}

// ===== Step 956-960: RabbitMQ Consumer =====

type RabbitMQConsumer struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

func NewRabbitMQConsumer(url string) (*RabbitMQConsumer, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }

    return &RabbitMQConsumer{
        conn:    conn,
        channel: ch,
    }, nil
}

func (c *RabbitMQConsumer) Consume(queueName string, handler func([]byte) error) error {
    // Declare queue
    _, err := c.channel.QueueDeclare(
        queueName,
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

    // Get messages
    msgs, err := c.channel.Consume(
        queueName,
        "",    // consumer
        false, // auto-ack
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil,   // args
    )
    if err != nil {
        return err
    }

    // Process messages
    for msg := range msgs {
        if err := handler(msg.Body); err != nil {
            msg.Nack(false, true) // Requeue
        } else {
            msg.Ack(false)
        }
    }

    return nil
}

// ===== Step 961-970: Event-Driven Architecture =====

type Event struct {
    Type      string                 `json:"type"`
    Timestamp time.Time              `json:"timestamp"`
    Data      map[string]interface{} `json:"data"`
}

type EventBus struct {
    subscribers map[string][]chan Event
    mu          sync.RWMutex
}

func NewEventBus() *EventBus {
    return &EventBus{
        subscribers: make(map[string][]chan Event),
    }
}

func (eb *EventBus) Subscribe(eventType string) <-chan Event {
    eb.mu.Lock()
    defer eb.mu.Unlock()

    ch := make(chan Event, 10)
    eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)

    return ch
}

func (eb *EventBus) Publish(event Event) {
    eb.mu.RLock()
    defer eb.mu.RUnlock()

    if subscribers, exists := eb.subscribers[event.Type]; exists {
        for _, ch := range subscribers {
            select {
            case ch <- event:
            default:
                // Channel full, skip
            }
        }
    }
}
```

### Step 971-990: CI/CD Pipeline

```yaml
# ===== Step 971-975: GitHub Actions =====

# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-

    - name: Install dependencies
      run: go mod download

    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...

    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.out

    - name: Run linter
      uses: golangci/golangci-lint-action@v3

  build:
    needs: test
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Build
      run: go build -v ./...

    - name: Build Docker image
      run: docker build -t myapp:${{ github.sha }} .

    - name: Push to Docker Hub
      if: github.ref == 'refs/heads/main'
      run: |
        echo ${{ secrets.DOCKER_PASSWORD }} | docker login -u ${{ secrets.DOCKER_USERNAME }} --password-stdin
        docker tag myapp:${{ github.sha }} myregistry/myapp:latest
        docker push myregistry/myapp:latest

# ===== Step 976-980: GitLab CI =====

# .gitlab-ci.yml
stages:
  - test
  - build
  - deploy

test:
  stage: test
  image: golang:1.21
  script:
    - go test -v -race -coverprofile=coverage.out ./...
    - go tool cover -html=coverage.out -o coverage.html
  artifacts:
    paths:
      - coverage.html
    expire_in: 1 week

build:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
  only:
    - main

deploy:
  stage: deploy
  script:
    - kubectl set image deployment/myapp myapp=$CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
    - kubectl rollout status deployment/myapp
  only:
    - main
  environment:
    name: production

# ===== Step 981-990: Makefile =====

# Makefile
.PHONY: build test run clean docker-build docker-run

build:
	go build -o bin/myapp cmd/main.go

test:
	go test -v -race -coverprofile=coverage.out ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

run:
	go run cmd/main.go

clean:
	rm -rf bin/
	go clean

docker-build:
	docker build -t myapp:latest .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

lint:
	golangci-lint run

fmt:
	go fmt ./...
	gofmt -s -w .

vet:
	go vet ./...

deps:
	go mod download
	go mod tidy

migrate-up:
	migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" down

proto:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  run           - Run the application"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
```

### Step 991-1000: Final Production Checklist

```go
/*
===== Step 991-1000: Production Readiness Checklist =====

✅ CODE QUALITY
- [ ] All tests passing (unit, integration, e2e)
- [ ] Code coverage > 80%
- [ ] Linter passing (golangci-lint)
- [ ] No race conditions (go test -race)
- [ ] Documentation complete
- [ ] Code reviewed

✅ PERFORMANCE
- [ ] Load testing completed
- [ ] Benchmark tests passing
- [ ] Memory leaks checked
- [ ] CPU profiling done
- [ ] Database queries optimized
- [ ] Caching implemented

✅ SECURITY
- [ ] Authentication implemented
- [ ] Authorization working
- [ ] Input validation
- [ ] SQL injection prevention
- [ ] XSS prevention
- [ ] CSRF protection
- [ ] Rate limiting
- [ ] HTTPS enabled
- [ ] Secrets management (vault/secrets manager)
- [ ] Security headers configured

✅ MONITORING & LOGGING
- [ ] Structured logging (JSON)
- [ ] Log aggregation (ELK/Splunk)
- [ ] Metrics collection (Prometheus)
- [ ] Alerting configured
- [ ] Dashboards created (Grafana)
- [ ] Tracing implemented (Jaeger/Zipkin)
- [ ] Error tracking (Sentry)

✅ INFRASTRUCTURE
- [ ] Docker images optimized
- [ ] Kubernetes manifests ready
- [ ] Health checks implemented
- [ ] Readiness probes configured
- [ ] Resource limits set
- [ ] Horizontal Pod Autoscaler configured
- [ ] Load balancer configured
- [ ] CDN configured (if needed)

✅ DATABASE
- [ ] Migrations tested
- [ ] Backup strategy implemented
- [ ] Connection pooling configured
- [ ] Indexes optimized
- [ ] Replication configured
- [ ] Disaster recovery plan

✅ CI/CD
- [ ] Automated tests in pipeline
- [ ] Docker build automated
- [ ] Deployment automated
- [ ] Rollback strategy defined
- [ ] Blue-green/canary deployment

✅ DOCUMENTATION
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Architecture diagrams
- [ ] Deployment guide
- [ ] Troubleshooting guide
- [ ] Runbooks
- [ ] README.md complete

✅ COMPLIANCE
- [ ] GDPR compliance (if applicable)
- [ ] Data encryption
- [ ] Audit logging
- [ ] Privacy policy
- [ ] Terms of service

🎉 CONGRATULATIONS! 🎉

คุณได้เรียนรู้ Go ครบทั้ง 1000 steps แล้ว!

ตอนนี้คุณพร้อมที่จะ:
- สร้าง Production-ready Applications
- Design Microservices Architecture
- Deploy to Cloud (AWS, GCP, Azure)
- Implement CI/CD Pipelines
- Monitor และ Scale Applications
- Lead Development Teams

Next Steps:
1. สร้างโปรเจคขนาดใหญ่
2. Contribute to Open Source
3. เขียน Technical Blog
4. Mentor ผู้อื่น
5. Build Products!

Happy Coding with Go! 🚀
*/
```

---

**จบ Step 1-1000 สมบูรณ์!** 🎉

