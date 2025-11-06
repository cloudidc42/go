// ===== ตัวอย่าง: Database Patterns =====
// NOTE: ตัวอย่างนี้ใช้ in-memory database เพื่อไม่ต้องติดตั้ง SQL database
// แต่ patterns ที่ใช้เหมือนกับการทำงานกับ database จริง

package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// ===== Part 1: Models =====

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ===== Part 2: Repository Interface Pattern =====

type UserRepository interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetByUsername(username string) (*User, error)
	GetAll() ([]*User, error)
	Update(user *User) error
	Delete(id int) error
	Count() (int, error)
}

type PostRepository interface {
	Create(post *Post) error
	GetByID(id int) (*Post, error)
	GetByUserID(userID int) ([]*Post, error)
	GetAll() ([]*Post, error)
	GetPublished() ([]*Post, error)
	Update(post *Post) error
	Delete(id int) error
	Count() (int, error)
}

type CommentRepository interface {
	Create(comment *Comment) error
	GetByID(id int) (*Comment, error)
	GetByPostID(postID int) ([]*Comment, error)
	GetByUserID(userID int) ([]*Comment, error)
	GetAll() ([]*Comment, error)
	Delete(id int) error
	Count() (int, error)
}

// ===== Part 3: In-Memory Implementation =====

type InMemoryUserRepository struct {
	mu     sync.RWMutex
	users  map[int]*User
	nextID int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

func (r *InMemoryUserRepository) Create(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if username already exists
	for _, u := range r.users {
		if u.Username == user.Username {
			return errors.New("username already exists")
		}
	}

	user.ID = r.nextID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	r.users[user.ID] = user
	r.nextID++

	return nil
}

func (r *InMemoryUserRepository) GetByID(id int) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *InMemoryUserRepository) GetByUsername(username string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *InMemoryUserRepository) GetAll() ([]*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

func (r *InMemoryUserRepository) Update(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user

	return nil
}

func (r *InMemoryUserRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(r.users, id)
	return nil
}

func (r *InMemoryUserRepository) Count() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.users), nil
}

// ===== Post Repository Implementation =====

type InMemoryPostRepository struct {
	mu     sync.RWMutex
	posts  map[int]*Post
	nextID int
}

func NewInMemoryPostRepository() *InMemoryPostRepository {
	return &InMemoryPostRepository{
		posts:  make(map[int]*Post),
		nextID: 1,
	}
}

func (r *InMemoryPostRepository) Create(post *Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	post.ID = r.nextID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	r.posts[post.ID] = post
	r.nextID++

	return nil
}

func (r *InMemoryPostRepository) GetByID(id int) (*Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}

	return post, nil
}

func (r *InMemoryPostRepository) GetByUserID(userID int) ([]*Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var posts []*Post
	for _, post := range r.posts {
		if post.UserID == userID {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func (r *InMemoryPostRepository) GetAll() ([]*Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	posts := make([]*Post, 0, len(r.posts))
	for _, post := range r.posts {
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *InMemoryPostRepository) GetPublished() ([]*Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var posts []*Post
	for _, post := range r.posts {
		if post.Published {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func (r *InMemoryPostRepository) Update(post *Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.posts[post.ID]; !exists {
		return errors.New("post not found")
	}

	post.UpdatedAt = time.Now()
	r.posts[post.ID] = post

	return nil
}

func (r *InMemoryPostRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.posts[id]; !exists {
		return errors.New("post not found")
	}

	delete(r.posts, id)
	return nil
}

func (r *InMemoryPostRepository) Count() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.posts), nil
}

// ===== Part 4: Service Layer =====

type BlogService struct {
	userRepo    UserRepository
	postRepo    PostRepository
	commentRepo CommentRepository
}

func NewBlogService(userRepo UserRepository, postRepo PostRepository, commentRepo CommentRepository) *BlogService {
	return &BlogService{
		userRepo:    userRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
	}
}

// User operations
func (s *BlogService) CreateUser(username, email string) (*User, error) {
	user := &User{
		Username: username,
		Email:    email,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *BlogService) GetUser(id int) (*User, error) {
	return s.userRepo.GetByID(id)
}

func (s *BlogService) GetAllUsers() ([]*User, error) {
	return s.userRepo.GetAll()
}

// Post operations
func (s *BlogService) CreatePost(userID int, title, content string) (*Post, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	post := &Post{
		UserID:    userID,
		Title:     title,
		Content:   content,
		Published: false,
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *BlogService) PublishPost(postID int) error {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	post.Published = true
	return s.postRepo.Update(post)
}

func (s *BlogService) GetUserPosts(userID int) ([]*Post, error) {
	return s.postRepo.GetByUserID(userID)
}

func (s *BlogService) GetPublishedPosts() ([]*Post, error) {
	return s.postRepo.GetPublished()
}

// Statistics
func (s *BlogService) GetStatistics() (map[string]int, error) {
	userCount, _ := s.userRepo.Count()
	postCount, _ := s.postRepo.Count()
	commentCount, _ := s.commentRepo.Count()

	return map[string]int{
		"users":    userCount,
		"posts":    postCount,
		"comments": commentCount,
	}, nil
}

// ===== Part 5: Query Builder Pattern (Simulated) =====

type QueryBuilder struct {
	table      string
	conditions []string
	orderBy    string
	limit      int
}

func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{
		table: table,
	}
}

func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
	qb.conditions = append(qb.conditions, condition)
	return qb
}

func (qb *QueryBuilder) OrderBy(field string) *QueryBuilder {
	qb.orderBy = field
	return qb
}

func (qb *QueryBuilder) Limit(n int) *QueryBuilder {
	qb.limit = n
	return qb
}

func (qb *QueryBuilder) Build() string {
	query := fmt.Sprintf("SELECT * FROM %s", qb.table)

	if len(qb.conditions) > 0 {
		query += " WHERE "
		for i, cond := range qb.conditions {
			if i > 0 {
				query += " AND "
			}
			query += cond
		}
	}

	if qb.orderBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", qb.orderBy)
	}

	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limit)
	}

	return query
}

// ===== Main Demo =====

func main() {
	fmt.Println("===== Database Patterns ใน Go =====\n")

	// ===== Initialize Repositories =====
	userRepo := NewInMemoryUserRepository()
	postRepo := NewInMemoryPostRepository()
	commentRepo := &InMemoryCommentRepository{
		comments: make(map[int]*Comment),
		nextID:   1,
	}

	// ===== Initialize Service =====
	service := NewBlogService(userRepo, postRepo, commentRepo)

	// ===== Create Users =====
	fmt.Println("--- Creating Users ---")

	user1, err := service.CreateUser("alice", "alice@example.com")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Created user: %s (ID: %d)\n", user1.Username, user1.ID)
	}

	user2, err := service.CreateUser("bob", "bob@example.com")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Created user: %s (ID: %d)\n", user2.Username, user2.ID)
	}

	user3, err := service.CreateUser("carol", "carol@example.com")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Created user: %s (ID: %d)\n", user3.Username, user3.ID)
	}

	// Try to create duplicate
	_, err = service.CreateUser("alice", "alice2@example.com")
	if err != nil {
		fmt.Printf("✓ Duplicate check worked: %v\n", err)
	}

	// ===== Create Posts =====
	fmt.Println("\n--- Creating Posts ---")

	post1, _ := service.CreatePost(user1.ID, "Introduction to Go", "Go is an amazing language...")
	fmt.Printf("✓ Created post: '%s' by user %d\n", post1.Title, post1.UserID)

	post2, _ := service.CreatePost(user1.ID, "Concurrency in Go", "Goroutines and channels...")
	fmt.Printf("✓ Created post: '%s' by user %d\n", post2.Title, post2.UserID)

	post3, _ := service.CreatePost(user2.ID, "Web Development", "Building APIs with Go...")
	fmt.Printf("✓ Created post: '%s' by user %d\n", post3.Title, post3.UserID)

	// ===== Publish Posts =====
	fmt.Println("\n--- Publishing Posts ---")

	service.PublishPost(post1.ID)
	fmt.Printf("✓ Published: %s\n", post1.Title)

	service.PublishPost(post3.ID)
	fmt.Printf("✓ Published: %s\n", post3.Title)

	// ===== Query Data =====
	fmt.Println("\n--- Querying Data ---")

	// Get all users
	users, _ := service.GetAllUsers()
	fmt.Printf("\nTotal users: %d\n", len(users))
	for _, u := range users {
		fmt.Printf("  - %s (%s)\n", u.Username, u.Email)
	}

	// Get published posts
	published, _ := service.GetPublishedPosts()
	fmt.Printf("\nPublished posts: %d\n", len(published))
	for _, p := range published {
		fmt.Printf("  - %s (by user %d)\n", p.Title, p.UserID)
	}

	// Get posts by user
	alicePosts, _ := service.GetUserPosts(user1.ID)
	fmt.Printf("\nAlice's posts: %d\n", len(alicePosts))
	for _, p := range alicePosts {
		fmt.Printf("  - %s (published: %v)\n", p.Title, p.Published)
	}

	// ===== Statistics =====
	fmt.Println("\n--- Statistics ---")

	stats, _ := service.GetStatistics()
	fmt.Printf("Total Users: %d\n", stats["users"])
	fmt.Printf("Total Posts: %d\n", stats["posts"])
	fmt.Printf("Total Comments: %d\n", stats["comments"])

	// ===== Query Builder Demo =====
	fmt.Println("\n--- Query Builder Pattern ---")

	query1 := NewQueryBuilder("users").
		Where("age > 18").
		Where("active = true").
		OrderBy("created_at DESC").
		Limit(10).
		Build()

	fmt.Println("Generated SQL:")
	fmt.Println(query1)

	query2 := NewQueryBuilder("posts").
		Where("published = true").
		OrderBy("created_at DESC").
		Limit(5).
		Build()

	fmt.Println("\n" + query2)

	// ===== Transaction Pattern (Simulated) =====
	fmt.Println("\n--- Transaction Pattern (Concept) ---")

	fmt.Println("In real database, you would:")
	fmt.Println("  1. Begin transaction")
	fmt.Println("  2. Execute multiple operations")
	fmt.Println("  3. Commit if all succeed")
	fmt.Println("  4. Rollback if any fails")

	fmt.Println("\nExample:")
	fmt.Println("  tx.Begin()")
	fmt.Println("  tx.Create(user)")
	fmt.Println("  tx.Create(profile)")
	fmt.Println("  tx.Commit()")

	// ===== Summary =====
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Database Patterns Summary:")
	fmt.Println("  ✓ Repository Pattern: Abstraction over data access")
	fmt.Println("  ✓ Service Layer: Business logic separation")
	fmt.Println("  ✓ Interfaces: Flexibility and testability")
	fmt.Println("  ✓ Thread-safe: Using sync.RWMutex")
	fmt.Println("  ✓ Query Builder: Dynamic SQL generation")
	fmt.Println(strings.Repeat("=", 60))
}

// ===== Comment Repository Implementation =====

type InMemoryCommentRepository struct {
	mu       sync.RWMutex
	comments map[int]*Comment
	nextID   int
}

func (r *InMemoryCommentRepository) Create(comment *Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	comment.ID = r.nextID
	comment.CreatedAt = time.Now()

	r.comments[comment.ID] = comment
	r.nextID++

	return nil
}

func (r *InMemoryCommentRepository) GetByID(id int) (*Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, exists := r.comments[id]
	if !exists {
		return nil, errors.New("comment not found")
	}

	return comment, nil
}

func (r *InMemoryCommentRepository) GetByPostID(postID int) ([]*Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var comments []*Comment
	for _, comment := range r.comments {
		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}

	return comments, nil
}

func (r *InMemoryCommentRepository) GetByUserID(userID int) ([]*Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var comments []*Comment
	for _, comment := range r.comments {
		if comment.UserID == userID {
			comments = append(comments, comment)
		}
	}

	return comments, nil
}

func (r *InMemoryCommentRepository) GetAll() ([]*Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comments := make([]*Comment, 0, len(r.comments))
	for _, comment := range r.comments {
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *InMemoryCommentRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.comments[id]; !exists {
		return errors.New("comment not found")
	}

	delete(r.comments, id)
	return nil
}

func (r *InMemoryCommentRepository) Count() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.comments), nil
}

// วิธีรัน:
// go run 06_database_patterns.go
//
// Key Patterns:
//
// 1. Repository Pattern:
//    - Abstraction layer over data access
//    - Interface-based for flexibility
//    - Easy to test with mocks
//    - Can swap implementations (in-memory, SQL, NoSQL)
//
// 2. Service Layer:
//    - Business logic separation
//    - Orchestrates multiple repositories
//    - Transaction management
//    - Validation and error handling
//
// 3. Models/Entities:
//    - Domain objects
//    - Data structure definition
//    - Validation methods
//
// 4. Query Builder:
//    - Dynamic SQL generation
//    - Type-safe queries
//    - Composable query parts
//
// Best Practices:
// - Use interfaces for repositories
// - Separate business logic from data access
// - Use transactions for multiple operations
// - Handle errors properly
// - Use connection pooling
// - Implement caching where appropriate
// - Use prepared statements (SQL injection prevention)
// - Implement proper indexing
// - Use migrations for schema changes
//
// Real Database Example (PostgreSQL/MySQL):
// - Use database/sql package
// - Use sqlx for easier querying
// - Use GORM or similar ORM
// - Implement connection pooling
// - Use transactions
// - Handle NULL values
// - Use context for timeouts
