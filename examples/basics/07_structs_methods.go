// ===== ตัวอย่าง: Structs และ Methods =====
package main

import (
	"fmt"
	"math"
	"time"
)

// ===== Part 1: Basic Structs =====

type Person struct {
	FirstName string
	LastName  string
	Age       int
	Email     string
}

type Point struct {
	X, Y float64
}

type Rectangle struct {
	Width, Height float64
}

type Circle struct {
	Radius float64
}

// ===== Part 2: Embedded Structs =====

type Address struct {
	Street  string
	City    string
	Country string
	ZipCode string
}

type Employee struct {
	Person  // Embedded struct (composition)
	Address // Embedded struct
	ID      int
	Salary  float64
	HireDate time.Time
}

// ===== Part 3: Methods =====

// Method สำหรับ Person
func (p Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

func (p Person) Greet() string {
	return fmt.Sprintf("Hello, I'm %s and I'm %d years old.", p.FullName(), p.Age)
}

// Method สำหรับ Rectangle
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Method สำหรับ Circle
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Circumference() float64 {
	return 2 * math.Pi * c.Radius
}

// Method สำหรับ Point
func (p Point) Distance(other Point) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// ===== Part 4: Pointer Receivers (แก้ไขได้) =====

// Pointer receiver - สามารถแก้ไข struct ได้
func (p *Person) Birthday() {
	p.Age++
}

func (p *Person) UpdateEmail(newEmail string) {
	p.Email = newEmail
}

// Pointer receiver สำหรับ Rectangle
func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

// ===== Part 5: Bank Account Example =====

type BankAccount struct {
	AccountNumber string
	Owner         string
	Balance       float64
}

func NewBankAccount(accountNumber, owner string, initialBalance float64) *BankAccount {
	return &BankAccount{
		AccountNumber: accountNumber,
		Owner:         owner,
		Balance:       initialBalance,
	}
}

func (ba *BankAccount) Deposit(amount float64) {
	if amount > 0 {
		ba.Balance += amount
		fmt.Printf("Deposited $%.2f. New balance: $%.2f\n", amount, ba.Balance)
	} else {
		fmt.Println("Invalid deposit amount")
	}
}

func (ba *BankAccount) Withdraw(amount float64) bool {
	if amount <= 0 {
		fmt.Println("Invalid withdrawal amount")
		return false
	}

	if amount > ba.Balance {
		fmt.Println("Insufficient funds")
		return false
	}

	ba.Balance -= amount
	fmt.Printf("Withdrawn $%.2f. New balance: $%.2f\n", amount, ba.Balance)
	return true
}

func (ba BankAccount) GetBalance() float64 {
	return ba.Balance
}

func (ba BankAccount) String() string {
	return fmt.Sprintf("Account[%s] Owner: %s, Balance: $%.2f",
		ba.AccountNumber, ba.Owner, ba.Balance)
}

// ===== Part 6: Student Management =====

type Grade struct {
	Subject string
	Score   float64
}

type Student struct {
	ID     string
	Name   string
	Grades []Grade
}

func NewStudent(id, name string) *Student {
	return &Student{
		ID:     id,
		Name:   name,
		Grades: make([]Grade, 0),
	}
}

func (s *Student) AddGrade(subject string, score float64) {
	s.Grades = append(s.Grades, Grade{
		Subject: subject,
		Score:   score,
	})
}

func (s Student) AverageGrade() float64 {
	if len(s.Grades) == 0 {
		return 0
	}

	total := 0.0
	for _, grade := range s.Grades {
		total += grade.Score
	}

	return total / float64(len(s.Grades))
}

func (s Student) PrintReport() {
	fmt.Printf("\n=== Report Card for %s (%s) ===\n", s.Name, s.ID)
	fmt.Println("Grades:")

	for _, grade := range s.Grades {
		fmt.Printf("  %-15s: %.2f\n", grade.Subject, grade.Score)
	}

	avg := s.AverageGrade()
	fmt.Printf("\nAverage: %.2f\n", avg)

	// Determine grade
	var letterGrade string
	switch {
	case avg >= 90:
		letterGrade = "A"
	case avg >= 80:
		letterGrade = "B"
	case avg >= 70:
		letterGrade = "C"
	case avg >= 60:
		letterGrade = "D"
	default:
		letterGrade = "F"
	}

	fmt.Printf("Letter Grade: %s\n", letterGrade)
}

// ===== Part 7: Task Manager =====

type TaskStatus int

const (
	Pending TaskStatus = iota
	InProgress
	Completed
)

func (ts TaskStatus) String() string {
	return [...]string{"Pending", "In Progress", "Completed"}[ts]
}

type Task struct {
	ID          int
	Title       string
	Description string
	Status      TaskStatus
	CreatedAt   time.Time
	CompletedAt *time.Time
}

func NewTask(id int, title, description string) *Task {
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      Pending,
		CreatedAt:   time.Now(),
	}
}

func (t *Task) Start() {
	if t.Status == Pending {
		t.Status = InProgress
		fmt.Printf("Task '%s' started\n", t.Title)
	}
}

func (t *Task) Complete() {
	if t.Status != Completed {
		t.Status = Completed
		now := time.Now()
		t.CompletedAt = &now
		fmt.Printf("Task '%s' completed\n", t.Title)
	}
}

func (t Task) String() string {
	status := fmt.Sprintf("[%s] %s", t.Status, t.Title)
	if t.CompletedAt != nil {
		duration := t.CompletedAt.Sub(t.CreatedAt)
		status += fmt.Sprintf(" (Completed in %v)", duration.Round(time.Second))
	}
	return status
}

// ===== Main Function =====

func main() {
	fmt.Println("===== Structs และ Methods ใน Go =====\n")

	// ===== Example 1: Basic Person =====
	fmt.Println("--- 1. Basic Person ---")

	person1 := Person{
		FirstName: "Alice",
		LastName:  "Smith",
		Age:       25,
		Email:     "alice@example.com",
	}

	fmt.Println(person1.Greet())
	fmt.Printf("Email: %s\n", person1.Email)

	// Birthday (ต้องใช้ pointer)
	person1.Birthday()
	fmt.Printf("After birthday: Age = %d\n", person1.Age)

	person1.UpdateEmail("alice.smith@example.com")
	fmt.Printf("Updated email: %s\n\n", person1.Email)

	// ===== Example 2: Geometric Shapes =====
	fmt.Println("--- 2. Geometric Shapes ---")

	rect := Rectangle{Width: 10, Height: 5}
	fmt.Printf("Rectangle: %.0f x %.0f\n", rect.Width, rect.Height)
	fmt.Printf("Area: %.2f\n", rect.Area())
	fmt.Printf("Perimeter: %.2f\n", rect.Perimeter())

	rect.Scale(2)
	fmt.Printf("\nAfter scaling by 2:\n")
	fmt.Printf("Rectangle: %.0f x %.0f\n", rect.Width, rect.Height)
	fmt.Printf("Area: %.2f\n", rect.Area())

	circle := Circle{Radius: 5}
	fmt.Printf("\nCircle with radius %.0f\n", circle.Radius)
	fmt.Printf("Area: %.2f\n", circle.Area())
	fmt.Printf("Circumference: %.2f\n\n", circle.Circumference())

	// ===== Example 3: Points and Distance =====
	fmt.Println("--- 3. Points and Distance ---")

	p1 := Point{X: 0, Y: 0}
	p2 := Point{X: 3, Y: 4}
	p3 := Point{X: 10, Y: 10}

	fmt.Printf("Point 1: (%.0f, %.0f)\n", p1.X, p1.Y)
	fmt.Printf("Point 2: (%.0f, %.0f)\n", p2.X, p2.Y)
	fmt.Printf("Distance: %.2f\n", p1.Distance(p2))

	fmt.Printf("\nPoint 2: (%.0f, %.0f)\n", p2.X, p2.Y)
	fmt.Printf("Point 3: (%.0f, %.0f)\n", p3.X, p3.Y)
	fmt.Printf("Distance: %.2f\n\n", p2.Distance(p3))

	// ===== Example 4: Embedded Structs =====
	fmt.Println("--- 4. Embedded Structs (Employee) ---")

	employee := Employee{
		Person: Person{
			FirstName: "Bob",
			LastName:  "Johnson",
			Age:       30,
			Email:     "bob@company.com",
		},
		Address: Address{
			Street:  "123 Main St",
			City:    "Bangkok",
			Country: "Thailand",
			ZipCode: "10110",
		},
		ID:       12345,
		Salary:   50000,
		HireDate: time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	// สามารถเข้าถึง fields จาก embedded struct ได้โดยตรง
	fmt.Printf("Employee: %s %s\n", employee.FirstName, employee.LastName)
	fmt.Printf("ID: %d\n", employee.ID)
	fmt.Printf("Email: %s\n", employee.Email)
	fmt.Printf("City: %s\n", employee.City)
	fmt.Printf("Salary: $%.2f\n", employee.Salary)
	fmt.Println(employee.Greet()) // Method จาก embedded Person
	fmt.Println()

	// ===== Example 5: Bank Account =====
	fmt.Println("--- 5. Bank Account ---")

	account := NewBankAccount("ACC001", "Alice Smith", 1000.0)
	fmt.Println(account)

	account.Deposit(500)
	account.Withdraw(200)
	account.Withdraw(2000) // Should fail

	fmt.Printf("\nFinal balance: $%.2f\n\n", account.GetBalance())

	// ===== Example 6: Student Management =====
	fmt.Println("--- 6. Student Management ---")

	student := NewStudent("STU001", "Charlie Brown")
	student.AddGrade("Mathematics", 85)
	student.AddGrade("Physics", 90)
	student.AddGrade("Chemistry", 88)
	student.AddGrade("English", 92)
	student.AddGrade("History", 87)

	student.PrintReport()

	// ===== Example 7: Task Manager =====
	fmt.Println("\n--- 7. Task Manager ---")

	tasks := []*Task{
		NewTask(1, "Design database schema", "Create ERD and tables"),
		NewTask(2, "Implement API endpoints", "Build REST API"),
		NewTask(3, "Write unit tests", "Test all functions"),
	}

	fmt.Println("Initial tasks:")
	for _, task := range tasks {
		fmt.Printf("  %s\n", task)
	}

	// Simulate task workflow
	fmt.Println("\nWorking on tasks...")
	tasks[0].Start()
	time.Sleep(100 * time.Millisecond)
	tasks[0].Complete()

	tasks[1].Start()
	time.Sleep(50 * time.Millisecond)
	tasks[1].Complete()

	fmt.Println("\nUpdated tasks:")
	for _, task := range tasks {
		fmt.Printf("  %s\n", task)
	}

	// ===== Example 8: Anonymous Structs =====
	fmt.Println("\n--- 8. Anonymous Structs ---")

	config := struct {
		Host string
		Port int
		SSL  bool
	}{
		Host: "localhost",
		Port: 8080,
		SSL:  false,
	}

	fmt.Printf("Config: %+v\n", config)
	fmt.Printf("Server: %s:%d (SSL: %v)\n", config.Host, config.Port, config.SSL)

	// ===== Example 9: Struct Comparison =====
	fmt.Println("\n--- 9. Struct Comparison ---")

	rect1 := Rectangle{Width: 10, Height: 5}
	rect2 := Rectangle{Width: 10, Height: 5}
	rect3 := Rectangle{Width: 8, Height: 6}

	fmt.Printf("rect1 == rect2: %v\n", rect1 == rect2)
	fmt.Printf("rect1 == rect3: %v\n", rect1 == rect3)

	// ===== Example 10: Struct Tags (Preview) =====
	fmt.Println("\n--- 10. Struct Tags (for JSON, etc.) ---")

	type Product struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	product := Product{
		ID:    1,
		Name:  "Laptop",
		Price: 999.99,
	}

	fmt.Printf("Product: %+v\n", product)
	fmt.Println("(Struct tags are used with json.Marshal, etc.)")
}

// วิธีรัน:
// go run 07_structs_methods.go
//
// Key Points:
// - Structs คือ collection ของ fields
// - Methods เป็น functions ที่ผูกกับ struct
// - Value receiver: func (s Struct) Method()
// - Pointer receiver: func (s *Struct) Method() - แก้ไขได้
// - Embedded structs = composition (inheritance-like)
// - Constructor pattern: NewType() *Type
//
// Best Practices:
// - ใช้ pointer receivers เมื่อต้องแก้ไข struct
// - ใช้ pointer receivers สำหรับ structs ขนาดใหญ่
// - ใช้ value receivers สำหรับ immutable operations
// - Export structs และ fields ด้วย capital letter
// - ใช้ constructor functions (New...) สำหรับ initialization
