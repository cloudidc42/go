// ===== Complete CRUD API with In-Memory Database =====
// Production-ready REST API example
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

type Product struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Stock       int       `json:"stock"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
}

type UpdateProductRequest struct {
    Name        *string  `json:"name,omitempty"`
    Description *string  `json:"description,omitempty"`
    Price       *float64 `json:"price,omitempty"`
    Stock       *int     `json:"stock,omitempty"`
}

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

// ===== In-Memory Database =====

type ProductStore struct {
    mu       sync.RWMutex
    products map[int]*Product
    nextID   int
}

func NewProductStore() *ProductStore {
    return &ProductStore{
        products: make(map[int]*Product),
        nextID:   1,
    }
}

func (s *ProductStore) Create(req *CreateProductRequest) *Product {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now()
    product := &Product{
        ID:          s.nextID,
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    s.products[s.nextID] = product
    s.nextID++

    return product
}

func (s *ProductStore) GetAll() []*Product {
    s.mu.RLock()
    defer s.mu.RUnlock()

    products := make([]*Product, 0, len(s.products))
    for _, product := range s.products {
        products = append(products, product)
    }

    return products
}

func (s *ProductStore) GetByID(id int) (*Product, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    product, exists := s.products[id]
    return product, exists
}

func (s *ProductStore) Update(id int, req *UpdateProductRequest) (*Product, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()

    product, exists := s.products[id]
    if !exists {
        return nil, false
    }

    if req.Name != nil {
        product.Name = *req.Name
    }
    if req.Description != nil {
        product.Description = *req.Description
    }
    if req.Price != nil {
        product.Price = *req.Price
    }
    if req.Stock != nil {
        product.Stock = *req.Stock
    }

    product.UpdatedAt = time.Now()

    return product, true
}

func (s *ProductStore) Delete(id int) bool {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, exists := s.products[id]; !exists {
        return false
    }

    delete(s.products, id)
    return true
}

func (s *ProductStore) Search(query string) []*Product {
    s.mu.RLock()
    defer s.mu.RUnlock()

    results := []*Product{}
    query = strings.ToLower(query)

    for _, product := range s.products {
        if strings.Contains(strings.ToLower(product.Name), query) ||
            strings.Contains(strings.ToLower(product.Description), query) {
            results = append(results, product)
        }
    }

    return results
}

// ===== API Server =====

type APIServer struct {
    store *ProductStore
}

func NewAPIServer() *APIServer {
    server := &APIServer{
        store: NewProductStore(),
    }

    // Seed some data
    server.seedData()

    return server
}

func (s *APIServer) seedData() {
    products := []CreateProductRequest{
        {
            Name:        "Laptop",
            Description: "High-performance laptop for developers",
            Price:       1299.99,
            Stock:       15,
        },
        {
            Name:        "Wireless Mouse",
            Description: "Ergonomic wireless mouse",
            Price:       29.99,
            Stock:       50,
        },
        {
            Name:        "Mechanical Keyboard",
            Description: "RGB mechanical gaming keyboard",
            Price:       89.99,
            Stock:       30,
        },
        {
            Name:        "USB-C Hub",
            Description: "7-in-1 USB-C hub adapter",
            Price:       49.99,
            Stock:       25,
        },
        {
            Name:        "Monitor Stand",
            Description: "Adjustable monitor stand with storage",
            Price:       39.99,
            Stock:       20,
        },
    }

    for _, p := range products {
        s.store.Create(&p)
    }
}

func (s *APIServer) setupRoutes() http.Handler {
    mux := http.NewServeMux()

    // API Routes
    mux.HandleFunc("/api/products", s.handleProducts)
    mux.HandleFunc("/api/products/", s.handleProduct)
    mux.HandleFunc("/api/products/search", s.handleSearch)

    // Root
    mux.HandleFunc("/", s.handleRoot)

    // Apply middleware
    handler := loggingMiddleware(mux)
    handler = corsMiddleware(handler)

    return handler
}

// ===== Handlers =====

func (s *APIServer) handleRoot(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    html := `
<!DOCTYPE html>
<html>
<head>
    <title>Product API</title>
    <style>
        body {
            font-family: 'Segoe UI', Arial, sans-serif;
            max-width: 1000px;
            margin: 50px auto;
            padding: 20px;
            background: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 { color: #00ADD8; }
        .endpoint {
            background: #f8f9fa;
            padding: 15px;
            margin: 10px 0;
            border-left: 4px solid #00ADD8;
            border-radius: 4px;
        }
        .method {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 3px;
            font-weight: bold;
            font-size: 12px;
            margin-right: 10px;
        }
        .get { background: #61affe; color: white; }
        .post { background: #49cc90; color: white; }
        .put { background: #fca130; color: white; }
        .delete { background: #f93e3e; color: white; }
        code {
            background: #e8e8e8;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
        }
        .example {
            background: #2d2d2d;
            color: #f8f8f2;
            padding: 15px;
            border-radius: 4px;
            overflow-x: auto;
            margin: 10px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🛍️ Product Management API</h1>
        <p>Complete CRUD REST API with in-memory database</p>

        <h2>📚 API Endpoints</h2>

        <div class="endpoint">
            <span class="method get">GET</span>
            <code>/api/products</code><br>
            <small>Get all products</small>
        </div>

        <div class="endpoint">
            <span class="method get">GET</span>
            <code>/api/products/{id}</code><br>
            <small>Get product by ID</small>
        </div>

        <div class="endpoint">
            <span class="method get">GET</span>
            <code>/api/products/search?q={query}</code><br>
            <small>Search products</small>
        </div>

        <div class="endpoint">
            <span class="method post">POST</span>
            <code>/api/products</code><br>
            <small>Create new product</small>
        </div>

        <div class="endpoint">
            <span class="method put">PUT</span>
            <code>/api/products/{id}</code><br>
            <small>Update product</small>
        </div>

        <div class="endpoint">
            <span class="method delete">DELETE</span>
            <code>/api/products/{id}</code><br>
            <small>Delete product</small>
        </div>

        <h2>📝 Example Requests</h2>

        <h3>Get all products:</h3>
        <div class="example">curl http://localhost:8080/api/products</div>

        <h3>Get specific product:</h3>
        <div class="example">curl http://localhost:8080/api/products/1</div>

        <h3>Search products:</h3>
        <div class="example">curl http://localhost:8080/api/products/search?q=laptop</div>

        <h3>Create product:</h3>
        <div class="example">curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"New Product","description":"Amazing product","price":99.99,"stock":10}'</div>

        <h3>Update product:</h3>
        <div class="example">curl -X PUT http://localhost:8080/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{"price":1199.99,"stock":20}'</div>

        <h3>Delete product:</h3>
        <div class="example">curl -X DELETE http://localhost:8080/api/products/1</div>

        <h2>✨ Features</h2>
        <ul>
            <li>✅ Complete CRUD operations</li>
            <li>✅ In-memory database (thread-safe)</li>
            <li>✅ Search functionality</li>
            <li>✅ JSON API</li>
            <li>✅ CORS enabled</li>
            <li>✅ Request logging</li>
            <li>✅ Error handling</li>
            <li>✅ Validation</li>
        </ul>
    </div>
</body>
</html>
    `

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, html)
}

func (s *APIServer) handleProducts(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        s.getAllProducts(w, r)
    case http.MethodPost:
        s.createProduct(w, r)
    default:
        respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
    }
}

func (s *APIServer) handleProduct(w http.ResponseWriter, r *http.Request) {
    // Extract ID from URL
    path := strings.TrimPrefix(r.URL.Path, "/api/products/")
    if path == "search" {
        s.handleSearch(w, r)
        return
    }

    id, err := strconv.Atoi(path)
    if err != nil {
        respondError(w, http.StatusBadRequest, "Invalid product ID")
        return
    }

    switch r.Method {
    case http.MethodGet:
        s.getProduct(w, r, id)
    case http.MethodPut:
        s.updateProduct(w, r, id)
    case http.MethodDelete:
        s.deleteProduct(w, r, id)
    default:
        respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
    }
}

func (s *APIServer) handleSearch(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
        return
    }

    query := r.URL.Query().Get("q")
    if query == "" {
        respondError(w, http.StatusBadRequest, "Query parameter 'q' is required")
        return
    }

    products := s.store.Search(query)
    respondSuccess(w, http.StatusOK, "Products found", products)
}

func (s *APIServer) getAllProducts(w http.ResponseWriter, r *http.Request) {
    products := s.store.GetAll()
    respondSuccess(w, http.StatusOK, "Products retrieved successfully", products)
}

func (s *APIServer) getProduct(w http.ResponseWriter, r *http.Request, id int) {
    product, exists := s.store.GetByID(id)
    if !exists {
        respondError(w, http.StatusNotFound, "Product not found")
        return
    }

    respondSuccess(w, http.StatusOK, "Product retrieved successfully", product)
}

func (s *APIServer) createProduct(w http.ResponseWriter, r *http.Request) {
    var req CreateProductRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }

    // Validation
    if req.Name == "" {
        respondError(w, http.StatusBadRequest, "Name is required")
        return
    }
    if req.Price < 0 {
        respondError(w, http.StatusBadRequest, "Price must be positive")
        return
    }
    if req.Stock < 0 {
        respondError(w, http.StatusBadRequest, "Stock must be non-negative")
        return
    }

    product := s.store.Create(&req)
    respondSuccess(w, http.StatusCreated, "Product created successfully", product)
}

func (s *APIServer) updateProduct(w http.ResponseWriter, r *http.Request, id int) {
    var req UpdateProductRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }

    // Validation
    if req.Price != nil && *req.Price < 0 {
        respondError(w, http.StatusBadRequest, "Price must be positive")
        return
    }
    if req.Stock != nil && *req.Stock < 0 {
        respondError(w, http.StatusBadRequest, "Stock must be non-negative")
        return
    }

    product, exists := s.store.Update(id, &req)
    if !exists {
        respondError(w, http.StatusNotFound, "Product not found")
        return
    }

    respondSuccess(w, http.StatusOK, "Product updated successfully", product)
}

func (s *APIServer) deleteProduct(w http.ResponseWriter, r *http.Request, id int) {
    if !s.store.Delete(id) {
        respondError(w, http.StatusNotFound, "Product not found")
        return
    }

    respondSuccess(w, http.StatusOK, "Product deleted successfully", nil)
}

// ===== Helper Functions =====

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

// ===== Main =====

func main() {
    server := NewAPIServer()
    handler := server.setupRoutes()

    fmt.Println("=== Product Management API ===")
    fmt.Println("Server starting on http://localhost:8080")
    fmt.Println("\nEndpoints:")
    fmt.Println("  GET    /api/products         - List all products")
    fmt.Println("  GET    /api/products/{id}    - Get product by ID")
    fmt.Println("  GET    /api/products/search  - Search products")
    fmt.Println("  POST   /api/products         - Create product")
    fmt.Println("  PUT    /api/products/{id}    - Update product")
    fmt.Println("  DELETE /api/products/{id}    - Delete product")
    fmt.Println("\nOpen http://localhost:8080 in your browser for documentation")
    fmt.Println("Press Ctrl+C to stop\n")

    log.Fatal(http.ListenAndServe(":8080", handler))
}

// วิธีรัน:
// go run 03_complete_crud_api.go
//
// ทดสอบ API:
// curl http://localhost:8080/api/products
// curl http://localhost:8080/api/products/1
// curl http://localhost:8080/api/products/search?q=laptop
// curl -X POST http://localhost:8080/api/products -H "Content-Type: application/json" -d '{"name":"Test","description":"Test product","price":19.99,"stock":5}'
// curl -X PUT http://localhost:8080/api/products/1 -H "Content-Type: application/json" -d '{"price":999.99}'
// curl -X DELETE http://localhost:8080/api/products/1
//
// Features:
// ✅ Complete CRUD operations
// ✅ Thread-safe in-memory database
// ✅ Search functionality
// ✅ JSON API
// ✅ CORS middleware
// ✅ Request logging
// ✅ Error handling
// ✅ Input validation
// ✅ Beautiful HTML documentation
// ✅ Production-ready code
//
// ใช้งานได้จริง 100%!
