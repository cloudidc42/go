// ===== ตัวอย่าง: Authentication & Authorization System =====
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ===== Part 1: User Model =====

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
	RoleMod   Role = "moderator"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose in JSON
	Role         Role      `json:"role"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	LastLogin    time.Time `json:"last_login,omitempty"`
}

// ===== Part 2: Session Management =====

type Session struct {
	Token     string
	UserID    int
	CreatedAt time.Time
	ExpiresAt time.Time
	IPAddress string
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*Session),
	}
}

func (s *SessionStore) Create(userID int, ipAddress string, duration time.Duration) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate secure random token
	token, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	session := &Session{
		Token:     token,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
		IPAddress: ipAddress,
	}

	s.sessions[token] = session
	return session, nil
}

func (s *SessionStore) Get(token string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[token]
	if !exists {
		return nil, errors.New("session not found")
	}

	// Check if expired
	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session expired")
	}

	return session, nil
}

func (s *SessionStore) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, token)
}

func (s *SessionStore) DeleteByUserID(userID int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for token, session := range s.sessions {
		if session.UserID == userID {
			delete(s.sessions, token)
		}
	}
}

func (s *SessionStore) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for token, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			delete(s.sessions, token)
		}
	}
}

// ===== Part 3: User Store =====

type UserStore struct {
	mu     sync.RWMutex
	users  map[int]*User
	nextID int
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

func (s *UserStore) Create(username, email, password string, role Role) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if username exists
	for _, user := range s.users {
		if user.Username == username {
			return nil, errors.New("username already exists")
		}
	}

	// Hash password
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           s.nextID,
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		Active:       true,
		CreatedAt:    time.Now(),
	}

	s.users[user.ID] = user
	s.nextID++

	return user, nil
}

func (s *UserStore) GetByID(id int) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *UserStore) GetByUsername(username string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (s *UserStore) UpdateLastLogin(userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	user.LastLogin = time.Now()
	return nil
}

// ===== Part 4: Auth Service =====

type AuthService struct {
	userStore    *UserStore
	sessionStore *SessionStore
}

func NewAuthService(userStore *UserStore, sessionStore *SessionStore) *AuthService {
	return &AuthService{
		userStore:    userStore,
		sessionStore: sessionStore,
	}
}

func (a *AuthService) Register(username, email, password string) (*User, error) {
	// Validate input
	if username == "" || email == "" || password == "" {
		return nil, errors.New("all fields are required")
	}

	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Create user with default role
	return a.userStore.Create(username, email, password, RoleUser)
}

func (a *AuthService) Login(username, password, ipAddress string) (*Session, *User, error) {
	// Get user
	user, err := a.userStore.GetByUsername(username)
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// Check if active
	if !user.Active {
		return nil, nil, errors.New("account is inactive")
	}

	// Verify password
	if !verifyPassword(user.PasswordHash, password) {
		return nil, nil, errors.New("invalid credentials")
	}

	// Update last login
	a.userStore.UpdateLastLogin(user.ID)

	// Create session
	session, err := a.sessionStore.Create(user.ID, ipAddress, 24*time.Hour)
	if err != nil {
		return nil, nil, err
	}

	return session, user, nil
}

func (a *AuthService) Logout(token string) error {
	a.sessionStore.Delete(token)
	return nil
}

func (a *AuthService) ValidateSession(token string) (*User, error) {
	session, err := a.sessionStore.Get(token)
	if err != nil {
		return nil, err
	}

	user, err := a.userStore.GetByID(session.UserID)
	if err != nil {
		return nil, err
	}

	if !user.Active {
		return nil, errors.New("account is inactive")
	}

	return user, nil
}

// ===== Part 5: Password Utilities =====

func hashPassword(password string) (string, error) {
	// Generate salt
	salt, err := generateSecureToken(16)
	if err != nil {
		return "", err
	}

	// Hash password with salt
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	hashedPassword := hex.EncodeToString(hash.Sum(nil))

	// Return salt:hash
	return salt + ":" + hashedPassword, nil
}

func verifyPassword(storedHash, password string) bool {
	parts := strings.Split(storedHash, ":")
	if len(parts) != 2 {
		return false
	}

	salt := parts[0]
	hash := parts[1]

	// Hash the provided password with the stored salt
	testHash := sha256.New()
	testHash.Write([]byte(password + salt))
	testHashString := hex.EncodeToString(testHash.Sum(nil))

	return hash == testHashString
}

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// ===== Part 6: Middleware =====

func (a *AuthService) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract token (Bearer <token>)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate session
		user, err := a.ValidateSession(token)
		if err != nil {
			http.Error(w, "Invalid or expired session", http.StatusUnauthorized)
			return
		}

		// Store user in context (simplified - normally use context.Context)
		r.Header.Set("X-User-ID", fmt.Sprintf("%d", user.ID))
		r.Header.Set("X-User-Role", string(user.Role))

		next(w, r)
	}
}

func (a *AuthService) RequireRole(role Role) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userRole := Role(r.Header.Get("X-User-Role"))

			if userRole != role && userRole != RoleAdmin {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next(w, r)
		}
	}
}

// ===== Part 7: HTTP Handlers =====

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, AuthResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	user, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, AuthResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	sendJSON(w, http.StatusCreated, AuthResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, AuthResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	session, user, err := h.authService.Login(req.Username, req.Password, r.RemoteAddr)
	if err != nil {
		sendJSON(w, http.StatusUnauthorized, AuthResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	sendJSON(w, http.StatusOK, AuthResponse{
		Success: true,
		Message: "Login successful",
		Data: map[string]interface{}{
			"token": session.Token,
			"user":  user,
		},
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		sendJSON(w, http.StatusOK, AuthResponse{
			Success: true,
			Message: "Logged out",
		})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 {
		h.authService.Logout(parts[1])
	}

	sendJSON(w, http.StatusOK, AuthResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// This handler is protected by AuthMiddleware
	// User info is in headers
	sendJSON(w, http.StatusOK, AuthResponse{
		Success: true,
		Data: map[string]string{
			"user_id": r.Header.Get("X-User-ID"),
			"role":    r.Header.Get("X-User-Role"),
		},
	})
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ===== Main =====

func main() {
	fmt.Println("===== Authentication & Authorization System =====\n")

	// Initialize stores and services
	userStore := NewUserStore()
	sessionStore := NewSessionStore()
	authService := NewAuthService(userStore, sessionStore)
	authHandler := NewAuthHandler(authService)

	// Create admin user for demo
	adminUser, _ := userStore.Create("admin", "admin@example.com", "admin123", RoleAdmin)
	fmt.Printf("Created admin user: %s (Role: %s)\n", adminUser.Username, adminUser.Role)

	// Session cleanup goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			sessionStore.Cleanup()
			fmt.Println("Cleaned up expired sessions")
		}
	}()

	// ===== Routes =====

	// Public endpoints
	http.HandleFunc("/api/register", authHandler.Register)
	http.HandleFunc("/api/login", authHandler.Login)

	// Protected endpoints
	http.HandleFunc("/api/logout", authService.AuthMiddleware(authHandler.Logout))
	http.HandleFunc("/api/profile", authService.AuthMiddleware(authHandler.GetProfile))

	// Admin-only endpoint
	http.HandleFunc("/api/admin/users", authService.AuthMiddleware(
		authService.RequireRole(RoleAdmin)(func(w http.ResponseWriter, r *http.Request) {
			sendJSON(w, http.StatusOK, AuthResponse{
				Success: true,
				Message: "This is an admin-only endpoint",
			})
		}),
	))

	// Documentation endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>Auth System API</title>
    <style>
        body { font-family: Arial; max-width: 900px; margin: 50px auto; padding: 20px; }
        h1 { color: #00ADD8; }
        .endpoint { background: #f4f4f4; padding: 15px; margin: 10px 0; border-radius: 5px; }
        code { background: #e8e8e8; padding: 2px 8px; border-radius: 3px; }
    </style>
</head>
<body>
    <h1>🔐 Authentication & Authorization System</h1>

    <h2>Endpoints:</h2>

    <div class="endpoint">
        <strong>POST /api/register</strong><br>
        Register new user<br>
        Body: <code>{"username": "john", "email": "john@example.com", "password": "password123"}</code>
    </div>

    <div class="endpoint">
        <strong>POST /api/login</strong><br>
        Login and get session token<br>
        Body: <code>{"username": "john", "password": "password123"}</code>
    </div>

    <div class="endpoint">
        <strong>POST /api/logout</strong><br>
        Logout (requires authentication)<br>
        Header: <code>Authorization: Bearer &lt;token&gt;</code>
    </div>

    <div class="endpoint">
        <strong>GET /api/profile</strong><br>
        Get user profile (requires authentication)<br>
        Header: <code>Authorization: Bearer &lt;token&gt;</code>
    </div>

    <div class="endpoint">
        <strong>GET /api/admin/users</strong><br>
        Admin only endpoint<br>
        Header: <code>Authorization: Bearer &lt;token&gt;</code>
    </div>

    <h2>Test Commands:</h2>
    <pre>
# Register
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"password123"}'

# Get Profile (use token from login)
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Logout
curl -X POST http://localhost:8080/api/logout \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Admin endpoint (login as admin first)
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

curl -X GET http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer ADMIN_TOKEN_HERE"
    </pre>
</body>
</html>
`
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
	})

	// Start server
	fmt.Println("\nServer starting on http://localhost:8080")
	fmt.Println("\nDefault credentials:")
	fmt.Println("  Username: admin")
	fmt.Println("  Password: admin123")
	fmt.Println("\nPress Ctrl+C to stop\n")

	http.ListenAndServe(":8080", nil)
}

// วิธีรัน:
// go run 03_auth_system.go
//
// ทดสอบด้วย curl:
//
// 1. Register new user:
//    curl -X POST http://localhost:8080/api/register \
//      -H "Content-Type: application/json" \
//      -d '{"username":"john","email":"john@example.com","password":"password123"}'
//
// 2. Login:
//    curl -X POST http://localhost:8080/api/login \
//      -H "Content-Type: application/json" \
//      -d '{"username":"john","password":"password123"}'
//
// 3. Get profile (copy token from login response):
//    curl http://localhost:8080/api/profile \
//      -H "Authorization: Bearer YOUR_TOKEN"
//
// 4. Admin login:
//    curl -X POST http://localhost:8080/api/login \
//      -H "Content-Type: application/json" \
//      -d '{"username":"admin","password":"admin123"}'
//
// Key Features:
// - User registration and login
// - Password hashing with salt
// - Session-based authentication
// - Role-based authorization (RBAC)
// - Secure token generation
// - Session expiration
// - Middleware for auth checks
// - Thread-safe stores
//
// Security Best Practices:
// - Never store plain text passwords
// - Use secure random tokens
// - Implement session expiration
// - Use HTTPS in production
// - Implement rate limiting
// - Use strong password requirements
// - Log authentication attempts
// - Implement account lockout
// - Use JWT for stateless auth (alternative)
