// ===== ตัวอย่าง: Interfaces ใน Go =====
package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// ===== Part 1: Basic Interfaces =====

// Shape interface - ทุก type ที่มี method Area() ถือว่า implement Shape
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle implements Shape
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle implements Shape
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Triangle implements Shape
type Triangle struct {
	Base, Height, Side1, Side2 float64
}

func (t Triangle) Area() float64 {
	return 0.5 * t.Base * t.Height
}

func (t Triangle) Perimeter() float64 {
	return t.Base + t.Side1 + t.Side2
}

// ฟังก์ชันที่รับ interface
func PrintShapeInfo(s Shape) {
	fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

// ===== Part 2: Multiple Interfaces =====

// Describer interface
type Describer interface {
	Describe() string
}

// Implement Describer for Rectangle
func (r Rectangle) Describe() string {
	return fmt.Sprintf("Rectangle: %.2f x %.2f", r.Width, r.Height)
}

// Implement Describer for Circle
func (c Circle) Describe() string {
	return fmt.Sprintf("Circle: radius %.2f", c.Radius)
}

// ===== Part 3: Empty Interface =====

// interface{} หรือ any รับได้ทุก type
func PrintAnything(value interface{}) {
	fmt.Printf("Type: %T, Value: %v\n", value, value)
}

// Type assertion
func DescribeValue(value interface{}) {
	switch v := value.(type) {
	case int:
		fmt.Printf("Integer: %d\n", v)
	case string:
		fmt.Printf("String: %s (length: %d)\n", v, len(v))
	case bool:
		fmt.Printf("Boolean: %v\n", v)
	case float64:
		fmt.Printf("Float: %.2f\n", v)
	case []int:
		fmt.Printf("Int slice: %v (length: %d)\n", v, len(v))
	default:
		fmt.Printf("Unknown type: %T with value: %v\n", v, v)
	}
}

// ===== Part 4: Interface Composition =====

// Reader interface
type Reader interface {
	Read() string
}

// Writer interface
type Writer interface {
	Write(data string)
}

// ReadWriter combines both
type ReadWriter interface {
	Reader
	Writer
}

// Document implements ReadWriter
type Document struct {
	Content string
}

func (d *Document) Read() string {
	return d.Content
}

func (d *Document) Write(data string) {
	d.Content = data
}

// ===== Part 5: Polymorphism Example =====

// Animal interface
type Animal interface {
	Speak() string
	Move() string
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "Woof!"
}

func (d Dog) Move() string {
	return "Running"
}

type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return "Meow!"
}

func (c Cat) Move() string {
	return "Walking gracefully"
}

type Bird struct {
	Name string
}

func (b Bird) Speak() string {
	return "Tweet!"
}

func (b Bird) Move() string {
	return "Flying"
}

// Function that works with any Animal
func AnimalDemo(a Animal, name string) {
	fmt.Printf("%s says: %s\n", name, a.Speak())
	fmt.Printf("%s is: %s\n", name, a.Move())
}

// ===== Part 6: Stringer Interface =====

// Stringer is built-in interface for custom string representation
type Person struct {
	FirstName string
	LastName  string
	Age       int
}

// Implement fmt.Stringer
func (p Person) String() string {
	return fmt.Sprintf("%s %s (Age: %d)", p.FirstName, p.LastName, p.Age)
}

// ===== Part 7: Error Interface =====

// Custom error type
type ValidationError struct {
	Field   string
	Message string
}

// Implement error interface
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on '%s': %s", e.Field, e.Message)
}

func ValidateAge(age int) error {
	if age < 0 {
		return ValidationError{Field: "age", Message: "must be positive"}
	}
	if age > 150 {
		return ValidationError{Field: "age", Message: "must be less than 150"}
	}
	return nil
}

// ===== Part 8: Sort Interface =====

type Person2 struct {
	Name string
	Age  int
}

// ByAge implements sort.Interface
type ByAge []Person2

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// ByName implements sort.Interface
type ByName []Person2

func (a ByName) Len() int           { return len(a) }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// ===== Part 9: Database Interface Example =====

// Database interface
type Database interface {
	Connect() error
	Disconnect() error
	Query(sql string) ([]map[string]interface{}, error)
}

// MySQL implementation
type MySQL struct {
	Host string
	Port int
}

func (m MySQL) Connect() error {
	fmt.Printf("Connecting to MySQL at %s:%d...\n", m.Host, m.Port)
	return nil
}

func (m MySQL) Disconnect() error {
	fmt.Println("Disconnecting from MySQL...")
	return nil
}

func (m MySQL) Query(sql string) ([]map[string]interface{}, error) {
	fmt.Printf("Executing MySQL query: %s\n", sql)
	return []map[string]interface{}{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}, nil
}

// PostgreSQL implementation
type PostgreSQL struct {
	Host string
	Port int
}

func (p PostgreSQL) Connect() error {
	fmt.Printf("Connecting to PostgreSQL at %s:%d...\n", p.Host, p.Port)
	return nil
}

func (p PostgreSQL) Disconnect() error {
	fmt.Println("Disconnecting from PostgreSQL...")
	return nil
}

func (p PostgreSQL) Query(sql string) ([]map[string]interface{}, error) {
	fmt.Printf("Executing PostgreSQL query: %s\n", sql)
	return []map[string]interface{}{
		{"id": 1, "name": "Charlie"},
		{"id": 2, "name": "David"},
	}, nil
}

// Function that works with any Database
func ExecuteQuery(db Database, sql string) {
	db.Connect()
	results, _ := db.Query(sql)
	fmt.Printf("Results: %v\n", results)
	db.Disconnect()
}

// ===== Part 10: Service Interface Pattern =====

// Logger interface
type Logger interface {
	Info(message string)
	Error(message string)
	Debug(message string)
}

// ConsoleLogger implements Logger
type ConsoleLogger struct {
	Prefix string
}

func (c ConsoleLogger) Info(message string) {
	fmt.Printf("[%sINFO] %s\n", c.Prefix, message)
}

func (c ConsoleLogger) Error(message string) {
	fmt.Printf("[%sERROR] %s\n", c.Prefix, message)
}

func (c ConsoleLogger) Debug(message string) {
	fmt.Printf("[%sDEBUG] %s\n", c.Prefix, message)
}

// FileLogger implements Logger (simulated)
type FileLogger struct {
	Filename string
}

func (f FileLogger) Info(message string) {
	fmt.Printf("[FileLogger:%s] INFO: %s\n", f.Filename, message)
}

func (f FileLogger) Error(message string) {
	fmt.Printf("[FileLogger:%s] ERROR: %s\n", f.Filename, message)
}

func (f FileLogger) Debug(message string) {
	fmt.Printf("[FileLogger:%s] DEBUG: %s\n", f.Filename, message)
}

// Service that uses Logger
type UserService struct {
	Logger Logger
}

func (s UserService) CreateUser(username string) {
	s.Logger.Info(fmt.Sprintf("Creating user: %s", username))
	// Create user logic here
	s.Logger.Info("User created successfully")
}

func (s UserService) DeleteUser(username string) {
	s.Logger.Info(fmt.Sprintf("Deleting user: %s", username))
	// Delete user logic here
	s.Logger.Error("User not found")
}

// ===== Main Function =====

func main() {
	fmt.Println("===== Interfaces ใน Go =====\n")

	// ===== Example 1: Basic Shapes =====
	fmt.Println("--- 1. Basic Shapes Interface ---")

	shapes := []Shape{
		Rectangle{Width: 10, Height: 5},
		Circle{Radius: 7},
		Triangle{Base: 6, Height: 8, Side1: 5, Side2: 5},
	}

	for i, shape := range shapes {
		fmt.Printf("Shape %d: ", i+1)
		PrintShapeInfo(shape)
	}
	fmt.Println()

	// ===== Example 2: Describer Interface =====
	fmt.Println("--- 2. Multiple Interfaces ---")

	rect := Rectangle{Width: 12, Height: 8}
	fmt.Println(rect.Describe())
	PrintShapeInfo(rect)

	circle := Circle{Radius: 5}
	fmt.Println(circle.Describe())
	PrintShapeInfo(circle)
	fmt.Println()

	// ===== Example 3: Empty Interface =====
	fmt.Println("--- 3. Empty Interface (interface{}) ---")

	PrintAnything(42)
	PrintAnything("Hello, Go!")
	PrintAnything(3.14159)
	PrintAnything(true)
	PrintAnything([]int{1, 2, 3})

	fmt.Println("\nType switching:")
	DescribeValue(100)
	DescribeValue("Hello")
	DescribeValue(true)
	DescribeValue(3.14)
	DescribeValue([]int{1, 2, 3, 4, 5})
	fmt.Println()

	// ===== Example 4: ReadWriter Interface =====
	fmt.Println("--- 4. Interface Composition ---")

	doc := &Document{}
	doc.Write("Hello, Go Interfaces!")
	fmt.Printf("Document content: %s\n\n", doc.Read())

	// ===== Example 5: Polymorphism =====
	fmt.Println("--- 5. Polymorphism with Animals ---")

	animals := []Animal{
		Dog{Name: "Buddy"},
		Cat{Name: "Whiskers"},
		Bird{Name: "Tweety"},
	}

	names := []string{"Buddy the Dog", "Whiskers the Cat", "Tweety the Bird"}

	for i, animal := range animals {
		AnimalDemo(animal, names[i])
		fmt.Println()
	}

	// ===== Example 6: Stringer Interface =====
	fmt.Println("--- 6. Stringer Interface ---")

	person := Person{FirstName: "Alice", LastName: "Smith", Age: 30}
	fmt.Println(person) // Automatically calls String() method

	people := []Person{
		{FirstName: "Bob", LastName: "Johnson", Age: 25},
		{FirstName: "Carol", LastName: "Williams", Age: 35},
	}

	for _, p := range people {
		fmt.Println(p)
	}
	fmt.Println()

	// ===== Example 7: Error Interface =====
	fmt.Println("--- 7. Custom Error Interface ---")

	ages := []int{25, -5, 30, 200, 45}

	for _, age := range ages {
		if err := ValidateAge(age); err != nil {
			fmt.Printf("Age %d: %v\n", age, err)
		} else {
			fmt.Printf("Age %d: Valid\n", age)
		}
	}
	fmt.Println()

	// ===== Example 8: Sort Interface =====
	fmt.Println("--- 8. Sort Interface ---")

	people2 := []Person2{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
		{Name: "David", Age: 28},
	}

	fmt.Println("Original:")
	for _, p := range people2 {
		fmt.Printf("  %s: %d years\n", p.Name, p.Age)
	}

	sort.Sort(ByAge(people2))
	fmt.Println("\nSorted by Age:")
	for _, p := range people2 {
		fmt.Printf("  %s: %d years\n", p.Name, p.Age)
	}

	sort.Sort(ByName(people2))
	fmt.Println("\nSorted by Name:")
	for _, p := range people2 {
		fmt.Printf("  %s: %d years\n", p.Name, p.Age)
	}
	fmt.Println()

	// ===== Example 9: Database Interface =====
	fmt.Println("--- 9. Database Interface (Polymorphism) ---")

	mysql := MySQL{Host: "localhost", Port: 3306}
	postgres := PostgreSQL{Host: "localhost", Port: 5432}

	fmt.Println("Using MySQL:")
	ExecuteQuery(mysql, "SELECT * FROM users")

	fmt.Println("\nUsing PostgreSQL:")
	ExecuteQuery(postgres, "SELECT * FROM users")
	fmt.Println()

	// ===== Example 10: Logger Interface =====
	fmt.Println("--- 10. Service with Logger Interface ---")

	consoleLogger := ConsoleLogger{Prefix: "APP "}
	fileLogger := FileLogger{Filename: "app.log"}

	fmt.Println("Using Console Logger:")
	service1 := UserService{Logger: consoleLogger}
	service1.CreateUser("alice")
	service1.DeleteUser("bob")

	fmt.Println("\nUsing File Logger:")
	service2 := UserService{Logger: fileLogger}
	service2.CreateUser("charlie")
	service2.DeleteUser("david")

	// ===== Bonus: Type Assertion =====
	fmt.Println("\n--- Bonus: Type Assertion ---")

	var i interface{} = "Hello, Go!"

	// Type assertion
	s, ok := i.(string)
	if ok {
		fmt.Printf("String value: %s (uppercase: %s)\n", s, strings.ToUpper(s))
	}

	// This would panic if not checked
	// n := i.(int) // panic!

	// Safe type assertion
	if n, ok := i.(int); ok {
		fmt.Printf("Int value: %d\n", n)
	} else {
		fmt.Println("Not an integer")
	}
}

// วิธีรัน:
// go run 04_interfaces.go
//
// Key Points:
// - Interface คือ collection ของ method signatures
// - Type ไม่ต้องประกาศว่า implement interface (implicit)
// - Empty interface (interface{} หรือ any) รับได้ทุก type
// - Type assertion: value.(Type)
// - Type switch: switch v := value.(type)
// - Interface composition: รวม interfaces เข้าด้วยกัน
//
// Common Patterns:
// - Polymorphism: หลาย types implement interface เดียวกัน
// - Dependency Injection: รับ interface แทน concrete type
// - Testing: ใช้ mock interfaces
// - Strategy Pattern: เปลี่ยน implementation ได้
//
// Built-in Interfaces:
// - error: Error() string
// - fmt.Stringer: String() string
// - io.Reader: Read([]byte) (int, error)
// - io.Writer: Write([]byte) (int, error)
// - sort.Interface: Len(), Less(), Swap()
