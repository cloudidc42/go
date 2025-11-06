// ===== ตัวอย่าง: REST API แบบสมบูรณ์ =====
// CRUD Operations (Create, Read, Update, Delete)
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"
    "sync"
    "time"
)

// ===== Models =====

type Book struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Author      string    `json:"author"`
    ISBN        string    `json:"isbn"`
    Price       float64   `json:"price"`
    PublishedAt time.Time `json:"published_at"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}

type SuccessResponse struct {
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// ===== Database (In-Memory) =====

type BookStore struct {
    mu    sync.RWMutex
    books map[int]*Book
    nextID int
}

func NewBookStore() *BookStore {
    return &BookStore{
        books: make(map[int]*Book),
        nextID: 1,
    }
}

func (bs *BookStore) GetAll() []*Book {
    bs.mu.RLock()
    defer bs.mu.RUnlock()

    books := make([]*Book, 0, len(bs.books))
    for _, book := range bs.books {
        books = append(books, book)
    }

    return books
}

func (bs *BookStore) GetByID(id int) (*Book, bool) {
    bs.mu.RLock()
    defer bs.mu.RUnlock()

    book, exists := bs.books[id]
    return book, exists
}

func (bs *BookStore) Create(book *Book) *Book {
    bs.mu.Lock()
    defer bs.mu.Unlock()

    book.ID = bs.nextID
    bs.books[bs.nextID] = book
    bs.nextID++

    return book
}

func (bs *BookStore) Update(id int, book *Book) bool {
    bs.mu.Lock()
    defer bs.mu.Unlock()

    if _, exists := bs.books[id]; !exists {
        return false
    }

    book.ID = id
    bs.books[id] = book
    return true
}

func (bs *BookStore) Delete(id int) bool {
    bs.mu.Lock()
    defer bs.mu.Unlock()

    if _, exists := bs.books[id]; !exists {
        return false
    }

    delete(bs.books, id)
    return true
}

// ===== API Server =====

type APIServer struct {
    store *BookStore
}

func NewAPIServer() *APIServer {
    return &APIServer{
        store: NewBookStore(),
    }
}

func (s *APIServer) Start() {
    // Seed some data
    s.seedData()

    // Routes
    http.HandleFunc("/", s.homeHandler)
    http.HandleFunc("/api/books", s.booksHandler)
    http.HandleFunc("/api/books/", s.bookHandler)

    // Middleware
    handler := loggingMiddleware(http.DefaultServeMux)
    handler = corsMiddleware(handler)

    fmt.Println("=== REST API Server ===")
    fmt.Println("Server running on http://localhost:8080")
    fmt.Println("\nEndpoints:")
    fmt.Println("  GET    /api/books      - Get all books")
    fmt.Println("  GET    /api/books/{id} - Get book by ID")
    fmt.Println("  POST   /api/books      - Create new book")
    fmt.Println("  PUT    /api/books/{id} - Update book")
    fmt.Println("  DELETE /api/books/{id} - Delete book")
    fmt.Println("\nPress Ctrl+C to stop\n")

    log.Fatal(http.ListenAndServe(":8080", handler))
}

func (s *APIServer) seedData() {
    books := []*Book{
        {
            Title:       "The Go Programming Language",
            Author:      "Alan A. A. Donovan",
            ISBN:        "978-0134190440",
            Price:       45.99,
            PublishedAt: time.Date(2015, 10, 26, 0, 0, 0, 0, time.UTC),
        },
        {
            Title:       "Concurrency in Go",
            Author:      "Katherine Cox-Buday",
            ISBN:        "978-1491941195",
            Price:       39.99,
            PublishedAt: time.Date(2017, 8, 31, 0, 0, 0, 0, time.UTC),
        },
        {
            Title:       "Go in Action",
            Author:      "William Kennedy",
            ISBN:        "978-1617291784",
            Price:       44.99,
            PublishedAt: time.Date(2015, 11, 15, 0, 0, 0, 0, time.UTC),
        },
    }

    for _, book := range books {
        s.store.Create(book)
    }
}

// ===== Handlers =====

func (s *APIServer) homeHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        sendError(w, http.StatusNotFound, "Not found")
        return
    }

    html := `
<!DOCTYPE html>
<html>
<head>
    <title>REST API - Book Store</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 900px; margin: 50px auto; padding: 20px; }
        h1 { color: #00ADD8; }
        .endpoint { background: #f4f4f4; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .method { display: inline-block; width: 80px; font-weight: bold; }
        .get { color: #61affe; }
        .post { color: #49cc90; }
        .put { color: #fca130; }
        .delete { color: #f93e3e; }
        code { background: #e8e8e8; padding: 2px 8px; border-radius: 3px; font-family: monospace; }
    </style>
</head>
<body>
    <h1>📚 Book Store REST API</h1>
    <p>Complete CRUD operations for managing books</p>

    <h2>Endpoints:</h2>

    <div class="endpoint">
        <span class="method get">GET</span>
        <code>/api/books</code><br>
        Get all books
    </div>

    <div class="endpoint">
        <span class="method get">GET</span>
        <code>/api/books/{id}</code><br>
        Get book by ID
    </div>

    <div class="endpoint">
        <span class="method post">POST</span>
        <code>/api/books</code><br>
        Create new book (send JSON in body)
    </div>

    <div class="endpoint">
        <span class="method put">PUT</span>
        <code>/api/books/{id}</code><br>
        Update existing book
    </div>

    <div class="endpoint">
        <span class="method delete">DELETE</span>
        <code>/api/books/{id}</code><br>
        Delete book
    </div>

    <h2>Examples:</h2>
    <h3>Get all books:</h3>
    <code>curl http://localhost:8080/api/books</code>

    <h3>Get specific book:</h3>
    <code>curl http://localhost:8080/api/books/1</code>

    <h3>Create new book:</h3>
    <code>curl -X POST http://localhost:8080/api/books \<br>
    &nbsp;&nbsp;-H "Content-Type: application/json" \<br>
    &nbsp;&nbsp;-d '{"title":"New Book","author":"John Doe","isbn":"123","price":29.99}'</code>

    <h3>Update book:</h3>
    <code>curl -X PUT http://localhost:8080/api/books/1 \<br>
    &nbsp;&nbsp;-H "Content-Type: application/json" \<br>
    &nbsp;&nbsp;-d '{"title":"Updated Title","author":"Jane Doe","isbn":"456","price":34.99}'</code>

    <h3>Delete book:</h3>
    <code>curl -X DELETE http://localhost:8080/api/books/1</code>
</body>
</html>
    `

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, html)
}

func (s *APIServer) booksHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        s.getAllBooks(w, r)
    case http.MethodPost:
        s.createBook(w, r)
    default:
        sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
    }
}

func (s *APIServer) bookHandler(w http.ResponseWriter, r *http.Request) {
    // Extract ID
    path := strings.TrimPrefix(r.URL.Path, "/api/books/")
    id, err := strconv.Atoi(path)
    if err != nil {
        sendError(w, http.StatusBadRequest, "Invalid book ID")
        return
    }

    switch r.Method {
    case http.MethodGet:
        s.getBook(w, r, id)
    case http.MethodPut:
        s.updateBook(w, r, id)
    case http.MethodDelete:
        s.deleteBook(w, r, id)
    default:
        sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
    }
}

// GET /api/books
func (s *APIServer) getAllBooks(w http.ResponseWriter, r *http.Request) {
    books := s.store.GetAll()
    sendJSON(w, http.StatusOK, books)
}

// GET /api/books/{id}
func (s *APIServer) getBook(w http.ResponseWriter, r *http.Request, id int) {
    book, exists := s.store.GetByID(id)
    if !exists {
        sendError(w, http.StatusNotFound, "Book not found")
        return
    }

    sendJSON(w, http.StatusOK, book)
}

// POST /api/books
func (s *APIServer) createBook(w http.ResponseWriter, r *http.Request) {
    var book Book

    if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
        sendError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }

    // Validation
    if book.Title == "" || book.Author == "" {
        sendError(w, http.StatusBadRequest, "Title and Author are required")
        return
    }

    created := s.store.Create(&book)

    sendJSON(w, http.StatusCreated, SuccessResponse{
        Message: "Book created successfully",
        Data:    created,
    })
}

// PUT /api/books/{id}
func (s *APIServer) updateBook(w http.ResponseWriter, r *http.Request, id int) {
    var book Book

    if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
        sendError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }

    if !s.store.Update(id, &book) {
        sendError(w, http.StatusNotFound, "Book not found")
        return
    }

    sendJSON(w, http.StatusOK, SuccessResponse{
        Message: "Book updated successfully",
        Data:    book,
    })
}

// DELETE /api/books/{id}
func (s *APIServer) deleteBook(w http.ResponseWriter, r *http.Request, id int) {
    if !s.store.Delete(id) {
        sendError(w, http.StatusNotFound, "Book not found")
        return
    }

    sendJSON(w, http.StatusOK, SuccessResponse{
        Message: "Book deleted successfully",
    })
}

// ===== Middleware =====

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        next.ServeHTTP(w, r)

        fmt.Printf("[%s] %s %s %s %v\n",
            start.Format("2006-01-02 15:04:05"),
            r.Method,
            r.URL.Path,
            r.RemoteAddr,
            time.Since(start),
        )
    })
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// ===== Helper Functions =====

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, message string) {
    sendJSON(w, status, ErrorResponse{
        Error:   http.StatusText(status),
        Message: message,
    })
}

// ===== Main =====

func main() {
    server := NewAPIServer()
    server.Start()
}

// วิธีรัน:
// go run 02_rest_api.go
//
// ทดสอบด้วย curl:
//
// 1. Get all books:
//    curl http://localhost:8080/api/books
//
// 2. Get specific book:
//    curl http://localhost:8080/api/books/1
//
// 3. Create new book:
//    curl -X POST http://localhost:8080/api/books \
//      -H "Content-Type: application/json" \
//      -d '{"title":"Learning Go","author":"Jon Bodner","isbn":"978-1492077213","price":49.99}'
//
// 4. Update book:
//    curl -X PUT http://localhost:8080/api/books/1 \
//      -H "Content-Type: application/json" \
//      -d '{"title":"Updated Title","author":"Updated Author","isbn":"111","price":39.99}'
//
// 5. Delete book:
//    curl -X DELETE http://localhost:8080/api/books/1
//
// หรือเปิดในเว็บเบราว์เซอร์:
// http://localhost:8080
