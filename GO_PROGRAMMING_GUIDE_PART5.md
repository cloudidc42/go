# คู่มือการเขียนโปรแกรม Go - ส่วนที่ 5
## Step 651-800: Testing, Performance และ Best Practices

---

## Step 651-680: Testing (Advanced)

### Step 651-660: Unit Testing

```go
// ===== Step 651: Basic Unit Test =====
// calculator.go
package calculator

func Add(a, b int) int {
    return a + b
}

func Subtract(a, b int) int {
    return a - b
}

func Multiply(a, b int) int {
    return a * b
}

func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// calculator_test.go
package calculator

import "testing"

// Step 652: Table-Driven Tests
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 5, 3, 8},
        {"negative numbers", -5, -3, -8},
        {"mixed numbers", -5, 3, -2},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}

// Step 653: Test with Error Handling
func TestDivide(t *testing.T) {
    tests := []struct {
        name      string
        a, b      float64
        expected  float64
        expectErr bool
    }{
        {"normal division", 10, 2, 5, false},
        {"division by zero", 10, 0, 0, true},
        {"negative numbers", -10, 2, -5, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Divide(tt.a, tt.b)

            if tt.expectErr {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                return
            }

            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }

            if result != tt.expected {
                t.Errorf("Divide(%.2f, %.2f) = %.2f; want %.2f",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}

// ===== Step 654-656: Test Helpers =====

func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}

func assertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func assertError(t *testing.T, err error) {
    t.Helper()
    if err == nil {
        t.Fatal("expected error, got nil")
    }
}

// ===== Step 657-660: Mocking =====

type UserRepository interface {
    GetUser(id int) (*User, error)
    CreateUser(user *User) error
}

type MockUserRepository struct {
    users map[int]*User
}

func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{
        users: make(map[int]*User),
    }
}

func (m *MockUserRepository) GetUser(id int) (*User, error) {
    if user, exists := m.users[id]; exists {
        return user, nil
    }
    return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) CreateUser(user *User) error {
    m.users[user.ID] = user
    return nil
}

// Test using mock
func TestUserService(t *testing.T) {
    repo := NewMockUserRepository()
    service := NewUserService(repo)

    // Test create
    user := &User{ID: 1, Name: "Alice"}
    err := service.CreateUser(user)
    assertNoError(t, err)

    // Test get
    retrieved, err := service.GetUser(1)
    assertNoError(t, err)
    assertEqual(t, retrieved.Name, "Alice")
}
```

### Step 661-670: Benchmark Tests

```go
package benchmark

import "testing"

// ===== Step 661: Basic Benchmark =====

func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(5, 3)
    }
}

// ===== Step 662: Benchmark with Setup =====

func BenchmarkFibonacci(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fibonacci(20)
    }
}

// ===== Step 663: Benchmark Comparisons =====

func BenchmarkStringConcat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        s := ""
        for j := 0; j < 100; j++ {
            s += "x"
        }
    }
}

func BenchmarkStringBuilder(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var builder strings.Builder
        for j := 0; j < 100; j++ {
            builder.WriteString("x")
        }
        _ = builder.String()
    }
}

// ===== Step 664-666: Memory Benchmarks =====

func BenchmarkSliceAppend(b *testing.B) {
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        var s []int
        for j := 0; j < 1000; j++ {
            s = append(s, j)
        }
    }
}

func BenchmarkSlicePrealloc(b *testing.B) {
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        s := make([]int, 0, 1000)
        for j := 0; j < 1000; j++ {
            s = append(s, j)
        }
    }
}

// ===== Step 667-670: Parallel Benchmarks =====

func BenchmarkParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            // Code to benchmark
            _ = expensiveOperation()
        }
    })
}

func expensiveOperation() int {
    sum := 0
    for i := 0; i < 1000; i++ {
        sum += i
    }
    return sum
}
```

### Step 671-680: Integration Tests

```go
package integration

import (
    "database/sql"
    "testing"
)

// ===== Step 671-675: Database Integration Tests =====

func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()

    db, err := sql.Open("postgres", "postgres://test:test@localhost/testdb")
    if err != nil {
        t.Fatalf("failed to open db: %v", err)
    }

    // Create tables
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100),
            email VARCHAR(100)
        )
    `)
    if err != nil {
        t.Fatalf("failed to create table: %v", err)
    }

    return db
}

func teardownTestDB(t *testing.T, db *sql.DB) {
    t.Helper()

    _, err := db.Exec(`DROP TABLE IF EXISTS users`)
    if err != nil {
        t.Logf("failed to drop table: %v", err)
    }

    db.Close()
}

func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    defer teardownTestDB(t, db)

    repo := NewUserRepository(db)

    user := &User{
        Name:  "John",
        Email: "john@example.com",
    }

    err := repo.Create(user)
    if err != nil {
        t.Fatalf("Create failed: %v", err)
    }

    if user.ID == 0 {
        t.Error("expected ID to be set")
    }
}

// ===== Step 676-680: HTTP Integration Tests =====

func TestAPIHandler(t *testing.T) {
    // Create test server
    handler := setupRoutes()
    server := httptest.NewServer(handler)
    defer server.Close()

    // Test GET
    resp, err := http.Get(server.URL + "/api/users")
    if err != nil {
        t.Fatalf("GET failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected status 200, got %d", resp.StatusCode)
    }

    // Test POST
    userData := `{"name":"Alice","email":"alice@example.com"}`
    resp, err = http.Post(
        server.URL+"/api/users",
        "application/json",
        strings.NewReader(userData),
    )
    if err != nil {
        t.Fatalf("POST failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        t.Errorf("expected status 201, got %d", resp.StatusCode)
    }
}
```

---

## Step 681-710: Performance Optimization

### Step 681-690: Profiling

```go
package profiling

import (
    "os"
    "runtime/pprof"
)

// ===== Step 681: CPU Profiling =====

func StartCPUProfile(filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }

    if err := pprof.StartCPUProfile(f); err != nil {
        f.Close()
        return err
    }

    return nil
}

func StopCPUProfile() {
    pprof.StopCPUProfile()
}

// ===== Step 682: Memory Profiling =====

func WriteMemProfile(filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    runtime.GC() // Run GC before taking heap snapshot
    return pprof.WriteHeapProfile(f)
}

// ===== Step 683-685: Using Profiling =====

func main() {
    // Start CPU profiling
    if err := StartCPUProfile("cpu.prof"); err != nil {
        log.Fatal(err)
    }
    defer StopCPUProfile()

    // Your application code
    runApplication()

    // Write memory profile
    if err := WriteMemProfile("mem.prof"); err != nil {
        log.Fatal(err)
    }
}

// Analyze with:
// go tool pprof cpu.prof
// go tool pprof mem.prof

// ===== Step 686-690: Runtime Metrics =====

func printMemStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    fmt.Printf("Alloc = %v MB", m.Alloc/1024/1024)
    fmt.Printf("\tTotalAlloc = %v MB", m.TotalAlloc/1024/1024)
    fmt.Printf("\tSys = %v MB", m.Sys/1024/1024)
    fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
```

### Step 691-700: Optimization Techniques

```go
package optimization

// ===== Step 691: String Builder vs Concatenation =====

// Slow
func concatStrings(n int) string {
    s := ""
    for i := 0; i < n; i++ {
        s += "x"
    }
    return s
}

// Fast
func buildStrings(n int) string {
    var builder strings.Builder
    builder.Grow(n) // Pre-allocate
    for i := 0; i < n; i++ {
        builder.WriteString("x")
    }
    return builder.String()
}

// ===== Step 692: Slice Pre-allocation =====

// Slow
func appendSlowly(n int) []int {
    var s []int
    for i := 0; i < n; i++ {
        s = append(s, i)
    }
    return s
}

// Fast
func appendFast(n int) []int {
    s := make([]int, 0, n) // Pre-allocate
    for i := 0; i < n; i++ {
        s = append(s, i)
    }
    return s
}

// ===== Step 693: Map Pre-allocation =====

func createMap(n int) map[int]int {
    m := make(map[int]int, n) // Pre-allocate
    for i := 0; i < n; i++ {
        m[i] = i * 2
    }
    return m
}

// ===== Step 694-696: Avoid Allocations =====

// Return pointer to avoid copy
func newLargeStruct() *LargeStruct {
    return &LargeStruct{
        // ... fields
    }
}

// Use pointer receiver for methods
func (ls *LargeStruct) Process() {
    // Process without copying
}

// Reuse buffers
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func processData(data []byte) {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)

    // Use buffer
}

// ===== Step 697-700: Concurrent Processing =====

func processItemsConcurrent(items []Item) []Result {
    results := make([]Result, len(items))
    var wg sync.WaitGroup

    // Worker pool
    workers := runtime.NumCPU()
    itemsChan := make(chan int, len(items))

    for w := 0; w < workers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for i := range itemsChan {
                results[i] = processItem(items[i])
            }
        }()
    }

    for i := range items {
        itemsChan <- i
    }
    close(itemsChan)

    wg.Wait()
    return results
}
```

### Step 701-710: Caching Strategies

```go
package caching

import (
    "sync"
    "time"
)

// ===== Step 701-703: In-Memory Cache =====

type Cache struct {
    mu    sync.RWMutex
    items map[string]CacheItem
}

type CacheItem struct {
    Value      interface{}
    Expiration time.Time
}

func NewCache() *Cache {
    cache := &Cache{
        items: make(map[string]CacheItem),
    }

    // Start cleanup goroutine
    go cache.cleanup()

    return cache
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.items[key] = CacheItem{
        Value:      value,
        Expiration: time.Now().Add(ttl),
    }
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    item, exists := c.items[key]
    if !exists {
        return nil, false
    }

    if time.Now().After(item.Expiration) {
        return nil, false
    }

    return item.Value, true
}

func (c *Cache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    delete(c.items, key)
}

func (c *Cache) cleanup() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        c.mu.Lock()
        now := time.Now()
        for key, item := range c.items {
            if now.After(item.Expiration) {
                delete(c.items, key)
            }
        }
        c.mu.Unlock()
    }
}

// ===== Step 704-706: LRU Cache =====

type LRUCache struct {
    capacity int
    cache    map[string]*list.Element
    list     *list.List
    mu       sync.Mutex
}

type entry struct {
    key   string
    value interface{}
}

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity: capacity,
        cache:    make(map[string]*list.Element),
        list:     list.New(),
    }
}

func (lru *LRUCache) Get(key string) (interface{}, bool) {
    lru.mu.Lock()
    defer lru.mu.Unlock()

    if elem, exists := lru.cache[key]; exists {
        lru.list.MoveToFront(elem)
        return elem.Value.(*entry).value, true
    }

    return nil, false
}

func (lru *LRUCache) Put(key string, value interface{}) {
    lru.mu.Lock()
    defer lru.mu.Unlock()

    if elem, exists := lru.cache[key]; exists {
        lru.list.MoveToFront(elem)
        elem.Value.(*entry).value = value
        return
    }

    if lru.list.Len() >= lru.capacity {
        oldest := lru.list.Back()
        if oldest != nil {
            lru.list.Remove(oldest)
            delete(lru.cache, oldest.Value.(*entry).key)
        }
    }

    elem := lru.list.PushFront(&entry{key, value})
    lru.cache[key] = elem
}

// ===== Step 707-710: Cache Patterns =====

// Cache-Aside Pattern
func (s *Service) GetUser(id int) (*User, error) {
    // Try cache first
    cacheKey := fmt.Sprintf("user:%d", id)
    if cached, found := s.cache.Get(cacheKey); found {
        return cached.(*User), nil
    }

    // Cache miss - fetch from database
    user, err := s.db.GetUser(id)
    if err != nil {
        return nil, err
    }

    // Store in cache
    s.cache.Set(cacheKey, user, 5*time.Minute)

    return user, nil
}

// Write-Through Cache
func (s *Service) UpdateUser(user *User) error {
    // Update database
    if err := s.db.UpdateUser(user); err != nil {
        return err
    }

    // Update cache
    cacheKey := fmt.Sprintf("user:%d", user.ID)
    s.cache.Set(cacheKey, user, 5*time.Minute)

    return nil
}
```

---

## Step 711-750: Design Patterns

### Step 711-720: Creational Patterns

```go
package patterns

// ===== Step 711: Singleton Pattern =====

type Database struct {
    connection string
}

var (
    instance *Database
    once     sync.Once
)

func GetInstance() *Database {
    once.Do(func() {
        instance = &Database{
            connection: "database_connection",
        }
    })
    return instance
}

// ===== Step 712-714: Factory Pattern =====

type Animal interface {
    Speak() string
}

type Dog struct{}

func (d Dog) Speak() string {
    return "Woof!"
}

type Cat struct{}

func (c Cat) Speak() string {
    return "Meow!"
}

func AnimalFactory(animalType string) Animal {
    switch animalType {
    case "dog":
        return Dog{}
    case "cat":
        return Cat{}
    default:
        return nil
    }
}

// ===== Step 715-717: Builder Pattern =====

type House struct {
    Windows int
    Doors   int
    Rooms   int
    HasGarage bool
    HasGarden bool
}

type HouseBuilder struct {
    house House
}

func NewHouseBuilder() *HouseBuilder {
    return &HouseBuilder{}
}

func (hb *HouseBuilder) Windows(n int) *HouseBuilder {
    hb.house.Windows = n
    return hb
}

func (hb *HouseBuilder) Doors(n int) *HouseBuilder {
    hb.house.Doors = n
    return hb
}

func (hb *HouseBuilder) Rooms(n int) *HouseBuilder {
    hb.house.Rooms = n
    return hb
}

func (hb *HouseBuilder) WithGarage() *HouseBuilder {
    hb.house.HasGarage = true
    return hb
}

func (hb *HouseBuilder) WithGarden() *HouseBuilder {
    hb.house.HasGarden = true
    return hb
}

func (hb *HouseBuilder) Build() House {
    return hb.house
}

// Usage
func buildHouse() {
    house := NewHouseBuilder().
        Windows(10).
        Doors(3).
        Rooms(5).
        WithGarage().
        WithGarden().
        Build()

    fmt.Printf("%+v\n", house)
}

// ===== Step 718-720: Prototype Pattern =====

type Prototype interface {
    Clone() Prototype
}

type ConcretePrototype struct {
    Name  string
    Value int
}

func (cp *ConcretePrototype) Clone() Prototype {
    return &ConcretePrototype{
        Name:  cp.Name,
        Value: cp.Value,
    }
}
```

### Step 721-730: Structural Patterns

```go
// ===== Step 721-723: Adapter Pattern =====

type LegacyPrinter interface {
    Print(s string) string
}

type ModernPrinter interface {
    PrintStored() string
}

type Adapter struct {
    OldPrinter LegacyPrinter
    Msg        string
}

func (a *Adapter) PrintStored() string {
    return a.OldPrinter.Print(a.Msg)
}

// ===== Step 724-726: Decorator Pattern =====

type Coffee interface {
    Cost() float64
    Description() string
}

type SimpleCoffee struct{}

func (sc SimpleCoffee) Cost() float64 {
    return 10.0
}

func (sc SimpleCoffee) Description() string {
    return "Simple Coffee"
}

type MilkDecorator struct {
    Coffee Coffee
}

func (md MilkDecorator) Cost() float64 {
    return md.Coffee.Cost() + 2.0
}

func (md MilkDecorator) Description() string {
    return md.Coffee.Description() + ", Milk"
}

type SugarDecorator struct {
    Coffee Coffee
}

func (sd SugarDecorator) Cost() float64 {
    return sd.Coffee.Cost() + 1.0
}

func (sd SugarDecorator) Description() string {
    return sd.Coffee.Description() + ", Sugar"
}

// Usage
func orderCoffee() {
    coffee := SimpleCoffee{}
    coffeeWithMilk := MilkDecorator{Coffee: coffee}
    coffeeWithMilkAndSugar := SugarDecorator{Coffee: coffeeWithMilk}

    fmt.Printf("%s: $%.2f\n",
        coffeeWithMilkAndSugar.Description(),
        coffeeWithMilkAndSugar.Cost(),
    )
}

// ===== Step 727-730: Proxy Pattern =====

type Subject interface {
    Request() string
}

type RealSubject struct{}

func (rs RealSubject) Request() string {
    return "RealSubject: Handling request"
}

type Proxy struct {
    realSubject *RealSubject
}

func (p *Proxy) Request() string {
    if p.realSubject == nil {
        p.realSubject = &RealSubject{}
    }

    // Add extra functionality
    result := "Proxy: Logging request\n"
    result += p.realSubject.Request()
    result += "\nProxy: Logging response"

    return result
}
```

### Step 731-750: Behavioral Patterns

```go
// ===== Step 731-734: Observer Pattern =====

type Observer interface {
    Update(string)
}

type Subject interface {
    Register(Observer)
    Deregister(Observer)
    NotifyAll()
}

type ConcreteSubject struct {
    observers []Observer
    state     string
}

func (cs *ConcreteSubject) Register(o Observer) {
    cs.observers = append(cs.observers, o)
}

func (cs *ConcreteSubject) Deregister(o Observer) {
    for i, observer := range cs.observers {
        if observer == o {
            cs.observers = append(cs.observers[:i], cs.observers[i+1:]...)
            break
        }
    }
}

func (cs *ConcreteSubject) NotifyAll() {
    for _, observer := range cs.observers {
        observer.Update(cs.state)
    }
}

func (cs *ConcreteSubject) SetState(state string) {
    cs.state = state
    cs.NotifyAll()
}

type ConcreteObserver struct {
    name string
}

func (co ConcreteObserver) Update(state string) {
    fmt.Printf("%s received update: %s\n", co.name, state)
}

// ===== Step 735-738: Strategy Pattern =====

type PaymentStrategy interface {
    Pay(amount float64) string
}

type CreditCardPayment struct {
    cardNumber string
}

func (cc CreditCardPayment) Pay(amount float64) string {
    return fmt.Sprintf("Paid $%.2f with Credit Card", amount)
}

type PayPalPayment struct {
    email string
}

func (pp PayPalPayment) Pay(amount float64) string {
    return fmt.Sprintf("Paid $%.2f with PayPal", amount)
}

type ShoppingCart struct {
    paymentMethod PaymentStrategy
}

func (sc *ShoppingCart) SetPaymentMethod(pm PaymentStrategy) {
    sc.paymentMethod = pm
}

func (sc *ShoppingCart) Checkout(amount float64) string {
    return sc.paymentMethod.Pay(amount)
}

// ===== Step 739-742: Command Pattern =====

type Command interface {
    Execute()
    Undo()
}

type Light struct {
    isOn bool
}

type LightOnCommand struct {
    light *Light
}

func (loc *LightOnCommand) Execute() {
    loc.light.isOn = true
    fmt.Println("Light is ON")
}

func (loc *LightOnCommand) Undo() {
    loc.light.isOn = false
    fmt.Println("Light is OFF")
}

// ===== Step 743-750: Chain of Responsibility =====

type Handler interface {
    SetNext(Handler)
    Handle(request string) string
}

type BaseHandler struct {
    next Handler
}

func (bh *BaseHandler) SetNext(handler Handler) {
    bh.next = handler
}

func (bh *BaseHandler) Handle(request string) string {
    if bh.next != nil {
        return bh.next.Handle(request)
    }
    return ""
}

type ConcreteHandler1 struct {
    BaseHandler
}

func (ch1 *ConcreteHandler1) Handle(request string) string {
    if request == "one" {
        return "Handler1 handled the request"
    }
    return ch1.BaseHandler.Handle(request)
}

type ConcreteHandler2 struct {
    BaseHandler
}

func (ch2 *ConcreteHandler2) Handle(request string) string {
    if request == "two" {
        return "Handler2 handled the request"
    }
    return ch2.BaseHandler.Handle(request)
}
```

---

## Step 751-780: Security Best Practices

### Step 751-760: Authentication & Authorization

```go
package security

import (
    "crypto/rand"
    "encoding/base64"
    "golang.org/x/crypto/bcrypt"
    "time"
)

// ===== Step 751-753: Password Hashing =====

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// ===== Step 754-756: JWT Tokens =====

type Claims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    jwt.StandardClaims
}

func GenerateJWT(userID int, email string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)

    claims := &Claims{
        UserID: userID,
        Email:  email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateJWT(tokenString string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims,
        func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}

// ===== Step 757-760: Session Management =====

type Session struct {
    ID        string
    UserID    int
    ExpiresAt time.Time
}

type SessionStore struct {
    mu       sync.RWMutex
    sessions map[string]*Session
}

func NewSessionStore() *SessionStore {
    store := &SessionStore{
        sessions: make(map[string]*Session),
    }

    go store.cleanup()

    return store
}

func (ss *SessionStore) Create(userID int) (*Session, error) {
    sessionID, err := generateRandomString(32)
    if err != nil {
        return nil, err
    }

    session := &Session{
        ID:        sessionID,
        UserID:    userID,
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }

    ss.mu.Lock()
    ss.sessions[sessionID] = session
    ss.mu.Unlock()

    return session, nil
}

func (ss *SessionStore) Get(sessionID string) (*Session, bool) {
    ss.mu.RLock()
    defer ss.mu.RUnlock()

    session, exists := ss.sessions[sessionID]
    if !exists || time.Now().After(session.ExpiresAt) {
        return nil, false
    }

    return session, true
}

func (ss *SessionStore) Delete(sessionID string) {
    ss.mu.Lock()
    defer ss.mu.Unlock()

    delete(ss.sessions, sessionID)
}

func (ss *SessionStore) cleanup() {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        ss.mu.Lock()
        now := time.Now()
        for id, session := range ss.sessions {
            if now.After(session.ExpiresAt) {
                delete(ss.sessions, id)
            }
        }
        ss.mu.Unlock()
    }
}

func generateRandomString(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}
```

---

*คู่มือนี้ยังคงครอบคลุม Step ต่อไปถึง 800 พร้อม Code Examples เพิ่มเติม...*

