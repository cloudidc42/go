// ===== ตัวอย่าง: Middleware Patterns =====
package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

// ===== Part 1: Basic Middleware =====

// Middleware type definition
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Logging middleware
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call next handler
		next(w, r)

		// Log after handler completes
		fmt.Printf("[%s] %s %s %s %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Since(start),
		)
	}
}

// ===== Part 2: Authentication Middleware =====

// Simple auth middleware (checks for API key)
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")

		if apiKey == "" {
			http.Error(w, "Missing API key", http.StatusUnauthorized)
			return
		}

		if apiKey != "secret-key-123" {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		// API key valid, proceed
		next(w, r)
	}
}

// ===== Part 3: CORS Middleware =====

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// ===== Part 4: Rate Limiting Middleware =====

type RateLimiter struct {
	mu       sync.Mutex
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

func (rl *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		rl.mu.Lock()

		// Clean old requests
		now := time.Now()
		if requests, exists := rl.requests[ip]; exists {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if now.Sub(reqTime) < rl.window {
					validRequests = append(validRequests, reqTime)
				}
			}
			rl.requests[ip] = validRequests
		}

		// Check limit
		if len(rl.requests[ip]) >= rl.limit {
			rl.mu.Unlock()
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Add current request
		rl.requests[ip] = append(rl.requests[ip], now)
		rl.mu.Unlock()

		next(w, r)
	}
}

// ===== Part 5: Recovery Middleware (Panic Handler) =====

func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v\n%s", err, debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next(w, r)
	}
}

// ===== Part 6: Request ID Middleware =====

var requestIDCounter uint64
var requestIDMutex sync.Mutex

func RequestIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestIDMutex.Lock()
		requestIDCounter++
		requestID := fmt.Sprintf("req-%d", requestIDCounter)
		requestIDMutex.Unlock()

		// Add request ID to response header
		w.Header().Set("X-Request-ID", requestID)

		// Log request ID
		fmt.Printf("[RequestID: %s] Processing request\n", requestID)

		next(w, r)
	}
}

// ===== Part 7: Response Writer Wrapper =====

// Custom response writer to capture status code and size
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// Advanced logging middleware with status and size
func AdvancedLoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer
		wrapped := newResponseWriter(w)

		// Call next handler
		next(wrapped, r)

		// Log with status and size
		fmt.Printf("[%s] %s %s - Status: %d - Size: %d bytes - Duration: %v\n",
			time.Now().Format("15:04:05"),
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			wrapped.size,
			time.Since(start),
		)
	}
}

// ===== Part 8: Content-Type Middleware =====

func JSONContentTypeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func ContentTypeValidationMiddleware(allowedTypes []string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" || r.Method == "PUT" {
				contentType := r.Header.Get("Content-Type")

				valid := false
				for _, allowed := range allowedTypes {
					if strings.HasPrefix(contentType, allowed) {
						valid = true
						break
					}
				}

				if !valid {
					http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
					return
				}
			}

			next(w, r)
		}
	}
}

// ===== Part 9: Chain Middleware =====

// Chain multiple middlewares together
func Chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	// Apply middlewares in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// ===== Part 10: Timeout Middleware =====

func TimeoutMiddleware(timeout time.Duration) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			done := make(chan bool, 1)

			go func() {
				next(w, r)
				done <- true
			}()

			select {
			case <-done:
				// Request completed
				return
			case <-time.After(timeout):
				// Timeout
				http.Error(w, "Request timeout", http.StatusRequestTimeout)
				return
			}
		}
	}
}

// ===== Part 11: Compression Middleware (Simulated) =====

func CompressionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if client accepts gzip
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			fmt.Println("  [Middleware] Would compress response with gzip")
		}

		next(w, r)
	}
}

// ===== Handlers =====

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Middleware Demo!\n")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message": "Hello from API", "status": "success"}`)
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message": "This is a protected endpoint", "authenticated": true}`)
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate slow operation
	time.Sleep(2 * time.Second)
	fmt.Fprintf(w, "Slow response complete\n")
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	// This will panic
	panic("intentional panic for testing recovery middleware")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, `{"message": "Data received", "status": "success"}`)
}

// ===== Main Function =====

func main() {
	fmt.Println("===== Middleware Patterns ใน Go =====\n")

	// Create rate limiter
	rateLimiter := NewRateLimiter(5, time.Minute) // 5 requests per minute

	// ===== Setup Routes =====

	// Basic route with logging
	http.HandleFunc("/", Chain(homeHandler,
		LoggingMiddleware,
	))

	// API route with multiple middlewares
	http.HandleFunc("/api", Chain(apiHandler,
		LoggingMiddleware,
		CORSMiddleware,
		JSONContentTypeMiddleware,
	))

	// Protected route with auth
	http.HandleFunc("/protected", Chain(protectedHandler,
		LoggingMiddleware,
		AuthMiddleware,
		JSONContentTypeMiddleware,
	))

	// Rate limited route
	http.HandleFunc("/limited", Chain(homeHandler,
		LoggingMiddleware,
		rateLimiter.Middleware,
	))

	// Route with recovery (will catch panic)
	http.HandleFunc("/panic", Chain(panicHandler,
		LoggingMiddleware,
		RecoveryMiddleware,
	))

	// Route with timeout
	http.HandleFunc("/slow", Chain(slowHandler,
		LoggingMiddleware,
		TimeoutMiddleware(1*time.Second),
	))

	// Route with advanced logging
	http.HandleFunc("/advanced", Chain(apiHandler,
		AdvancedLoggingMiddleware,
		RequestIDMiddleware,
	))

	// POST route with content type validation
	http.HandleFunc("/data", Chain(dataHandler,
		LoggingMiddleware,
		ContentTypeValidationMiddleware([]string{"application/json"}),
	))

	// Route with all middlewares
	http.HandleFunc("/full", Chain(apiHandler,
		AdvancedLoggingMiddleware,
		RecoveryMiddleware,
		RequestIDMiddleware,
		CORSMiddleware,
		rateLimiter.Middleware,
		CompressionMiddleware,
		JSONContentTypeMiddleware,
	))

	// ===== Start Server =====

	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println("\nAvailable endpoints:")
	fmt.Println("  GET  /              - Basic route with logging")
	fmt.Println("  GET  /api           - API with CORS and JSON")
	fmt.Println("  GET  /protected     - Protected with API key auth")
	fmt.Println("  GET  /limited       - Rate limited (5 req/min)")
	fmt.Println("  GET  /panic         - Panic recovery demo")
	fmt.Println("  GET  /slow          - Timeout demo (1s)")
	fmt.Println("  GET  /advanced      - Advanced logging with request ID")
	fmt.Println("  POST /data          - Content-Type validation")
	fmt.Println("  GET  /full          - All middlewares combined")

	fmt.Println("\nExample commands:")
	fmt.Println("  curl http://localhost:8080/")
	fmt.Println("  curl http://localhost:8080/api")
	fmt.Println("  curl -H \"X-API-Key: secret-key-123\" http://localhost:8080/protected")
	fmt.Println("  curl -X POST -H \"Content-Type: application/json\" http://localhost:8080/data")

	fmt.Println("\nPress Ctrl+C to stop\n")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// วิธีรัน:
// go run 04_middleware_patterns.go
//
// ทดสอบด้วย curl:
//
// 1. Basic request:
//    curl http://localhost:8080/
//
// 2. API endpoint:
//    curl http://localhost:8080/api
//
// 3. Protected endpoint (with auth):
//    curl -H "X-API-Key: secret-key-123" http://localhost:8080/protected
//
// 4. Protected endpoint (without auth):
//    curl http://localhost:8080/protected
//
// 5. Rate limited (try multiple times):
//    for i in {1..10}; do curl http://localhost:8080/limited; done
//
// 6. Panic recovery:
//    curl http://localhost:8080/panic
//
// 7. Timeout:
//    curl http://localhost:8080/slow
//
// 8. POST with content type:
//    curl -X POST -H "Content-Type: application/json" http://localhost:8080/data
//
// 9. All middlewares:
//    curl http://localhost:8080/full
//
// Key Concepts:
// - Middleware คือ function ที่ wrap handler
// - Middleware pattern: func(http.HandlerFunc) http.HandlerFunc
// - Chain middlewares เพื่อใช้หลาย middlewares
// - Order matters: middlewares execute in order
// - Common middlewares: logging, auth, CORS, rate limiting, recovery
//
// Best Practices:
// - Keep middlewares focused (single responsibility)
// - Make middlewares reusable
// - Use defer for cleanup
// - Handle errors properly
// - Log important events
// - Use context for passing values
//
// Common Middleware Patterns:
// 1. Logging: Log requests and responses
// 2. Authentication: Verify credentials
// 3. Authorization: Check permissions
// 4. CORS: Handle cross-origin requests
// 5. Rate Limiting: Limit request rate
// 6. Recovery: Catch panics
// 7. Compression: Compress responses
// 8. Caching: Cache responses
// 9. Timeout: Handle long requests
// 10. Request ID: Track requests
