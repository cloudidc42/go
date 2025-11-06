// ===== ตัวอย่าง: Web Server พื้นฐาน =====
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

// Struct สำหรับ JSON response
type Response struct {
    Message   string    `json:"message"`
    Timestamp time.Time `json:"timestamp"`
    Status    string    `json:"status"`
}

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Simple in-memory database
var users = []User{
    {ID: 1, Name: "Alice", Email: "alice@example.com"},
    {ID: 2, Name: "Bob", Email: "bob@example.com"},
    {ID: 3, Name: "Carol", Email: "carol@example.com"},
}

func main() {
    fmt.Println("=== Go Web Server ===")
    fmt.Println("Starting server on http://localhost:8080")
    fmt.Println("\nAvailable endpoints:")
    fmt.Println("  GET  /")
    fmt.Println("  GET  /hello")
    fmt.Println("  GET  /api/users")
    fmt.Println("  GET  /api/user/{id}")
    fmt.Println("  POST /api/data")
    fmt.Println("\nPress Ctrl+C to stop\n")

    // ===== Define Routes =====

    // Home page
    http.HandleFunc("/", homeHandler)

    // Hello endpoint
    http.HandleFunc("/hello", helloHandler)

    // API endpoints
    http.HandleFunc("/api/users", usersHandler)
    http.HandleFunc("/api/user/", userHandler)
    http.HandleFunc("/api/data", dataHandler)

    // Static file server
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Start server
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// ===== Handlers =====

// Home handler
func homeHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    html := `
<!DOCTYPE html>
<html>
<head>
    <title>Go Web Server</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        h1 { color: #00ADD8; }
        .endpoint { background: #f4f4f4; padding: 10px; margin: 10px 0; border-radius: 5px; }
        code { background: #e8e8e8; padding: 2px 5px; border-radius: 3px; }
    </style>
</head>
<body>
    <h1>🚀 Welcome to Go Web Server!</h1>
    <p>This is a simple web server built with Go.</p>

    <h2>Available Endpoints:</h2>

    <div class="endpoint">
        <strong>GET /hello</strong><br>
        Returns a welcome message
    </div>

    <div class="endpoint">
        <strong>GET /api/users</strong><br>
        Returns all users
    </div>

    <div class="endpoint">
        <strong>POST /api/data</strong><br>
        Accepts JSON data
    </div>

    <h2>Try it out:</h2>
    <p>
        <code>curl http://localhost:8080/hello</code><br>
        <code>curl http://localhost:8080/api/users</code>
    </p>
</body>
</html>
    `

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, html)

    logRequest(r)
}

// Hello handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
    // Get query parameter
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "World"
    }

    response := Response{
        Message:   fmt.Sprintf("Hello, %s!", name),
        Timestamp: time.Now(),
        Status:    "success",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)

    logRequest(r)
}

// Users handler - GET all users
func usersHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)

    logRequest(r)
}

// User handler - GET single user
func userHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Extract ID from URL (simple parsing)
    // In production, use a router like gorilla/mux
    id := r.URL.Path[len("/api/user/"):]

    // Find user
    for _, user := range users {
        if fmt.Sprintf("%d", user.ID) == id {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(user)
            logRequest(r)
            return
        }
    }

    http.Error(w, "User not found", http.StatusNotFound)
}

// Data handler - POST JSON data
func dataHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var data map[string]interface{}

    // Decode JSON from request body
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Process data
    response := Response{
        Message:   "Data received successfully",
        Timestamp: time.Now(),
        Status:    "success",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)

    logRequest(r)
    fmt.Printf("  Received data: %v\n", data)
}

// Helper function to log requests
func logRequest(r *http.Request) {
    fmt.Printf("[%s] %s %s %s\n",
        time.Now().Format("2006-01-02 15:04:05"),
        r.Method,
        r.URL.Path,
        r.RemoteAddr,
    )
}

// วิธีรัน:
// go run 01_web_server.go
//
// ทดสอบด้วย curl:
// curl http://localhost:8080/
// curl http://localhost:8080/hello
// curl http://localhost:8080/hello?name=Alice
// curl http://localhost:8080/api/users
// curl http://localhost:8080/api/user/1
// curl -X POST http://localhost:8080/api/data -H "Content-Type: application/json" -d '{"name":"test","value":123}'
//
// หรือเปิดในเว็บเบราว์เซอร์:
// http://localhost:8080
