# คู่มือการเขียนโปรแกรม Go - ส่วนที่ 4
## Step 501-650: Web Development และ REST API (Advanced)

---

## Step 501-530: Web Development Fundamentals

### Step 501-510: HTTP Server Advanced

```go
// ===== Step 501: Custom HTTP Server =====
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
)

func main() {
    // Create custom server
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      setupRoutes(),
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    // Start server in goroutine
    go func() {
        fmt.Println("Server starting on :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    // Graceful shutdown
    gracefulShutdown(srv)
}

func setupRoutes() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/api/", apiHandler)

    return mux
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Welcome to Go Server!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"message":"API endpoint"}`)
}

func gracefulShutdown(srv *http.Server) {
    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit

    fmt.Println("\nShutting down server...")

    // Timeout context for shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    fmt.Println("Server exited")
}

// ===== Step 502-505: Request Handling =====

func advancedRequestHandler(w http.ResponseWriter, r *http.Request) {
    // Step 502: Parse Query Parameters
    queryParams := r.URL.Query()
    name := queryParams.Get("name")
    age := queryParams.Get("age")

    // Step 503: Parse Form Data
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    // Step 504: Read Request Body
    body := make([]byte, r.ContentLength)
    r.Body.Read(body)
    defer r.Body.Close()

    // Step 505: Get Headers
    userAgent := r.Header.Get("User-Agent")
    contentType := r.Header.Get("Content-Type")

    // Response
    w.Header().Set("X-Custom-Header", "Go-Server")
    fmt.Fprintf(w, "Query: name=%s, age=%s\n", name, age)
    fmt.Fprintf(w, "User-Agent: %s\n", userAgent)
    fmt.Fprintf(w, "Content-Type: %s\n", contentType)
}

// ===== Step 506-510: Response Types =====

// Step 506: JSON Response
func jsonResponse(w http.ResponseWriter, r *http.Request) {
    type Response struct {
        Status  string `json:"status"`
        Message string `json:"message"`
        Data    map[string]interface{} `json:"data"`
    }

    resp := Response{
        Status:  "success",
        Message: "Data retrieved successfully",
        Data: map[string]interface{}{
            "id":    1,
            "name":  "John",
            "email": "john@example.com",
        },
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

// Step 507: XML Response
func xmlResponse(w http.ResponseWriter, r *http.Request) {
    type User struct {
        XMLName xml.Name `xml:"user"`
        ID      int      `xml:"id"`
        Name    string   `xml:"name"`
        Email   string   `xml:"email"`
    }

    user := User{
        ID:    1,
        Name:  "John",
        Email: "john@example.com",
    }

    w.Header().Set("Content-Type", "application/xml")
    xml.NewEncoder(w).Encode(user)
}

// Step 508: File Download
func fileDownloadHandler(w http.ResponseWriter, r *http.Request) {
    filename := "example.txt"
    content := []byte("This is a sample file content")

    w.Header().Set("Content-Disposition", "attachment; filename="+filename)
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))

    w.Write(content)
}

// Step 509: File Upload
func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse multipart form (32 MB max)
    if err := r.ParseMultipartForm(32 << 20); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    // Get file from form
    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Failed to get file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Save file
    dst, err := os.Create("./uploads/" + header.Filename)
    if err != nil {
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    io.Copy(dst, file)

    fmt.Fprintf(w, "File uploaded: %s (%d bytes)", header.Filename, header.Size)
}

// Step 510: Streaming Response
func streamingHandler(w http.ResponseWriter, r *http.Request) {
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming not supported", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    for i := 1; i <= 10; i++ {
        fmt.Fprintf(w, "data: Message %d\n\n", i)
        flusher.Flush()
        time.Sleep(1 * time.Second)
    }
}
```

### Step 511-520: Routing และ URL Parameters

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "regexp"
    "strings"
)

// ===== Step 511: Simple Router =====

type Route struct {
    Method  string
    Pattern *regexp.Regexp
    Handler http.HandlerFunc
}

type Router struct {
    routes []*Route
}

func NewRouter() *Router {
    return &Router{routes: []*Route{}}
}

func (router *Router) AddRoute(method, pattern string, handler http.HandlerFunc) {
    route := &Route{
        Method:  method,
        Pattern: regexp.MustCompile("^" + pattern + "$"),
        Handler: handler,
    }
    router.routes = append(router.routes, route)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    for _, route := range router.routes {
        if route.Method == r.Method && route.Pattern.MatchString(r.URL.Path) {
            route.Handler(w, r)
            return
        }
    }
    http.NotFound(w, r)
}

// ===== Step 512-515: URL Parameters =====

func extractPathParams(pattern, path string) map[string]string {
    params := make(map[string]string)

    patternParts := strings.Split(pattern, "/")
    pathParts := strings.Split(path, "/")

    if len(patternParts) != len(pathParts) {
        return params
    }

    for i, part := range patternParts {
        if strings.HasPrefix(part, ":") {
            key := strings.TrimPrefix(part, ":")
            params[key] = pathParts[i]
        }
    }

    return params
}

// ===== Step 516-520: RESTful Routes =====

func setupRESTRoutes() *Router {
    router := NewRouter()

    // Users
    router.AddRoute("GET", "/api/users", listUsersHandler)
    router.AddRoute("GET", "/api/users/[0-9]+", getUserHandler)
    router.AddRoute("POST", "/api/users", createUserHandler)
    router.AddRoute("PUT", "/api/users/[0-9]+", updateUserHandler)
    router.AddRoute("DELETE", "/api/users/[0-9]+", deleteUserHandler)

    // Posts
    router.AddRoute("GET", "/api/posts", listPostsHandler)
    router.AddRoute("GET", "/api/posts/[0-9]+", getPostHandler)
    router.AddRoute("POST", "/api/posts", createPostHandler)

    // Nested routes
    router.AddRoute("GET", "/api/users/[0-9]+/posts", getUserPostsHandler)

    return router
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
    users := []map[string]interface{}{
        {"id": 1, "name": "Alice"},
        {"id": 2, "name": "Bob"},
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
    // Extract ID from URL
    parts := strings.Split(r.URL.Path, "/")
    id := parts[len(parts)-1]

    user := map[string]interface{}{
        "id":   id,
        "name": "John Doe",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var user map[string]interface{}
    json.NewDecoder(r.Body).Decode(&user)

    user["id"] = 123 // Generate ID

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
    var user map[string]interface{}
    json.NewDecoder(r.Body).Decode(&user)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNoContent)
}

func listPostsHandler(w http.ResponseWriter, r *http.Request) {
    posts := []map[string]interface{}{
        {"id": 1, "title": "Post 1"},
        {"id": 2, "title": "Post 2"},
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

func getPostHandler(w http.ResponseWriter, r *http.Request) {
    post := map[string]interface{}{
        "id":    1,
        "title": "My Post",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
    var post map[string]interface{}
    json.NewDecoder(r.Body).Decode(&post)

    post["id"] = 123

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}

func getUserPostsHandler(w http.ResponseWriter, r *http.Request) {
    posts := []map[string]interface{}{
        {"id": 1, "title": "User's Post 1"},
        {"id": 2, "title": "User's Post 2"},
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}
```

### Step 521-530: Middleware

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

// ===== Step 521: Logging Middleware =====

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Call next handler
        next.ServeHTTP(w, r)

        // Log after handler
        log.Printf(
            "%s %s %s %v",
            r.Method,
            r.RequestURI,
            r.RemoteAddr,
            time.Since(start),
        )
    })
}

// ===== Step 522: CORS Middleware =====

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

// ===== Step 523: Authentication Middleware =====

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")

        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Verify token (simplified)
        if token != "Bearer valid-token" {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// ===== Step 524: Rate Limiting Middleware =====

type RateLimiter struct {
    requests map[string][]time.Time
    limit    int
    window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        requests: make(map[string][]time.Time),
        limit:    limit,
        window:   window,
    }
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        now := time.Now()

        // Clean old requests
        rl.requests[ip] = filterOldRequests(rl.requests[ip], now, rl.window)

        // Check limit
        if len(rl.requests[ip]) >= rl.limit {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        // Add request
        rl.requests[ip] = append(rl.requests[ip], now)

        next.ServeHTTP(w, r)
    })
}

func filterOldRequests(requests []time.Time, now time.Time, window time.Duration) []time.Time {
    var filtered []time.Time
    for _, t := range requests {
        if now.Sub(t) < window {
            filtered = append(filtered, t)
        }
    }
    return filtered
}

// ===== Step 525: Compression Middleware =====

type gzipResponseWriter struct {
    http.ResponseWriter
    Writer io.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}

func compressionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }

        w.Header().Set("Content-Encoding", "gzip")
        gz := gzip.NewWriter(w)
        defer gz.Close()

        gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
        next.ServeHTTP(gzw, r)
    })
}

// ===== Step 526-530: Middleware Chain =====

func chainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}

func setupWithMiddleware() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/api/", apiHandler)

    // Chain all middleware
    handler := chainMiddleware(
        mux,
        loggingMiddleware,
        corsMiddleware,
        compressionMiddleware,
    )

    return handler
}

// Recovery Middleware
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()

        next.ServeHTTP(w, r)
    })
}

// Request ID Middleware
func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := generateRequestID()
        w.Header().Set("X-Request-ID", requestID)

        next.ServeHTTP(w, r)
    })
}

func generateRequestID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}
```

---

## Step 531-560: Database Operations

### Step 531-540: SQL Database (PostgreSQL/MySQL)

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/lib/pq" // PostgreSQL driver
)

// ===== Step 531: Database Connection =====

type Database struct {
    *sql.DB
}

func NewDatabase(connectionString string) (*Database, error) {
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, err
    }

    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    // Test connection
    if err := db.Ping(); err != nil {
        return nil, err
    }

    return &Database{db}, nil
}

// ===== Step 532-535: CRUD Operations =====

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Create
func (db *Database) CreateUser(user *User) error {
    query := `
        INSERT INTO users (name, email, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

    now := time.Now()
    err := db.QueryRow(
        query,
        user.Name,
        user.Email,
        now,
        now,
    ).Scan(&user.ID)

    if err != nil {
        return err
    }

    user.CreatedAt = now
    user.UpdatedAt = now

    return nil
}

// Read
func (db *Database) GetUser(id int) (*User, error) {
    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    user := &User{}
    err := db.QueryRow(query, id).Scan(
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

// Read All
func (db *Database) GetAllUsers() ([]*User, error) {
    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        ORDER BY id
    `

    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    users := []*User{}
    for rows.Next() {
        user := &User{}
        err := rows.Scan(
            &user.ID,
            &user.Name,
            &user.Email,
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}

// Update
func (db *Database) UpdateUser(user *User) error {
    query := `
        UPDATE users
        SET name = $1, email = $2, updated_at = $3
        WHERE id = $4
    `

    user.UpdatedAt = time.Now()

    result, err := db.Exec(
        query,
        user.Name,
        user.Email,
        user.UpdatedAt,
        user.ID,
    )

    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("user not found")
    }

    return nil
}

// Delete
func (db *Database) DeleteUser(id int) error {
    query := `DELETE FROM users WHERE id = $1`

    result, err := db.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("user not found")
    }

    return nil
}

// ===== Step 536-540: Advanced Queries =====

// Search with LIKE
func (db *Database) SearchUsers(keyword string) ([]*User, error) {
    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE name ILIKE $1 OR email ILIKE $1
        ORDER BY name
    `

    rows, err := db.Query(query, "%"+keyword+"%")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    users := []*User{}
    for rows.Next() {
        user := &User{}
        rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
        users = append(users, user)
    }

    return users, nil
}

// Pagination
func (db *Database) GetUsersPaginated(page, pageSize int) ([]*User, error) {
    offset := (page - 1) * pageSize

    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        ORDER BY id
        LIMIT $1 OFFSET $2
    `

    rows, err := db.Query(query, pageSize, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    users := []*User{}
    for rows.Next() {
        user := &User{}
        rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
        users = append(users, user)
    }

    return users, nil
}

// Count
func (db *Database) CountUsers() (int, error) {
    var count int
    query := `SELECT COUNT(*) FROM users`
    err := db.QueryRow(query).Scan(&count)
    return count, err
}

// Transaction
func (db *Database) CreateUserWithProfile(user *User, profile map[string]interface{}) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }

    defer func() {
        if err != nil {
            tx.Rollback()
            return
        }
        err = tx.Commit()
    }()

    // Insert user
    query1 := `
        INSERT INTO users (name, email, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
    now := time.Now()
    err = tx.QueryRow(query1, user.Name, user.Email, now, now).Scan(&user.ID)
    if err != nil {
        return err
    }

    // Insert profile
    query2 := `
        INSERT INTO profiles (user_id, bio, avatar)
        VALUES ($1, $2, $3)
    `
    _, err = tx.Exec(query2, user.ID, profile["bio"], profile["avatar"])
    if err != nil {
        return err
    }

    return nil
}
```

*คู่มือนี้จะมีเนื้อหามากกว่า 10,000 บรรทัด ฉันจะสร้างต่อในส่วนถัดไป...*

