// ===== ตัวอย่าง: Testing ใน Go =====
// File: calculator.go (โค้ดที่จะทดสอบ)

package main

import (
	"errors"
	"fmt"
	"strings"
)

// ===== Calculator Functions =====

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
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func IsEven(n int) bool {
	return n%2 == 0
}

func Factorial(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("factorial of negative number")
	}
	if n == 0 || n == 1 {
		return 1, nil
	}

	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result, nil
}

// ===== String Functions =====

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func IsPalindrome(s string) bool {
	return s == Reverse(s)
}

func CountVowels(s string) int {
	vowels := "aeiouAEIOU"
	count := 0
	for _, char := range s {
		for _, vowel := range vowels {
			if char == vowel {
				count++
				break
			}
		}
	}
	return count
}

// ===== User Management =====

type User struct {
	ID       int
	Username string
	Email    string
	Age      int
}

func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Age < 0 {
		return errors.New("age must be positive")
	}
	if u.Age > 150 {
		return errors.New("age must be less than 150")
	}
	return nil
}

func (u *User) IsAdult() bool {
	return u.Age >= 18
}

type UserService struct {
	users map[int]*User
	nextID int
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[int]*User),
		nextID: 1,
	}
}

func (s *UserService) Create(username, email string, age int) (*User, error) {
	user := &User{
		ID:       s.nextID,
		Username: username,
		Email:    email,
		Age:      age,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	s.users[user.ID] = user
	s.nextID++

	return user, nil
}

func (s *UserService) GetByID(id int) (*User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) GetAll() []*User {
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

func (s *UserService) Update(id int, username, email string, age int) (*User, error) {
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.Username = username
	user.Email = email
	user.Age = age

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(id int) error {
	if _, exists := s.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(s.users, id)
	return nil
}

func (s *UserService) Count() int {
	return len(s.users)
}

// ===== Main Demo =====

func main() {
	fmt.Println("===== Testing Examples Demo =====\n")

	// ===== Calculator Demo =====
	fmt.Println("--- Calculator ---")
	fmt.Printf("Add(5, 3) = %d\n", Add(5, 3))
	fmt.Printf("Subtract(10, 4) = %d\n", Subtract(10, 4))
	fmt.Printf("Multiply(6, 7) = %d\n", Multiply(6, 7))

	if result, err := Divide(10, 2); err == nil {
		fmt.Printf("Divide(10, 2) = %.2f\n", result)
	}

	if _, err := Divide(10, 0); err != nil {
		fmt.Printf("Divide(10, 0) = Error: %v\n", err)
	}

	fmt.Printf("IsEven(4) = %v\n", IsEven(4))
	fmt.Printf("IsEven(7) = %v\n", IsEven(7))

	if result, err := Factorial(5); err == nil {
		fmt.Printf("Factorial(5) = %d\n", result)
	}

	// ===== String Functions Demo =====
	fmt.Println("\n--- String Functions ---")
	text := "hello"
	fmt.Printf("Reverse('%s') = '%s'\n", text, Reverse(text))
	fmt.Printf("IsPalindrome('radar') = %v\n", IsPalindrome("radar"))
	fmt.Printf("IsPalindrome('hello') = %v\n", IsPalindrome("hello"))
	fmt.Printf("CountVowels('hello world') = %d\n", CountVowels("hello world"))

	// ===== User Service Demo =====
	fmt.Println("\n--- User Service ---")

	service := NewUserService()

	// Create users
	user1, err := service.Create("alice", "alice@example.com", 25)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
	} else {
		fmt.Printf("Created user: %+v\n", user1)
		fmt.Printf("Is adult: %v\n", user1.IsAdult())
	}

	user2, err := service.Create("bob", "bob@example.com", 16)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
	} else {
		fmt.Printf("Created user: %+v\n", user2)
		fmt.Printf("Is adult: %v\n", user2.IsAdult())
	}

	// Try to create invalid user
	_, err = service.Create("", "invalid@example.com", 30)
	if err != nil {
		fmt.Printf("Validation error: %v\n", err)
	}

	// Get all users
	fmt.Printf("\nTotal users: %d\n", service.Count())
	fmt.Println("All users:")
	for _, user := range service.GetAll() {
		fmt.Printf("  - %s (%s), Age: %d\n", user.Username, user.Email, user.Age)
	}

	// Update user
	updatedUser, err := service.Update(1, "alice_updated", "alice.new@example.com", 26)
	if err != nil {
		fmt.Printf("Error updating user: %v\n", err)
	} else {
		fmt.Printf("\nUpdated user: %+v\n", updatedUser)
	}

	// Delete user
	err = service.Delete(2)
	if err != nil {
		fmt.Printf("Error deleting user: %v\n", err)
	} else {
		fmt.Println("Deleted user ID 2")
		fmt.Printf("Total users after delete: %d\n", service.Count())
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("NOTE: To run tests, create a file named 'calculator_test.go'")
	fmt.Println("Then run: go test -v")
	fmt.Println(strings.Repeat("=", 60))
}

// วิธีรัน:
// go run 05_testing_example.go
//
// วิธีรัน tests (สร้างไฟล์ *_test.go ก่อน):
// go test -v
// go test -cover
// go test -bench .
// go test -race
//
// NOTE: ตัวอย่าง test file ด้านล่างนี้ควรสร้างเป็นไฟล์แยก
// ชื่อ calculator_test.go

/*
===== calculator_test.go =====
package main

import "testing"

// ===== Basic Tests =====

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 10, 3, 7},
		{"negative result", 5, 10, -5},
		{"zero", 5, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Subtract(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Subtract(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	// Test successful division
	result, err := Divide(10, 2)
	if err != nil {
		t.Errorf("Divide(10, 2) unexpected error: %v", err)
	}
	if result != 5.0 {
		t.Errorf("Divide(10, 2) = %.2f; want 5.00", result)
	}

	// Test division by zero
	_, err = Divide(10, 0)
	if err == nil {
		t.Error("Divide(10, 0) expected error, got nil")
	}
}

func TestIsEven(t *testing.T) {
	tests := []struct {
		input    int
		expected bool
	}{
		{2, true},
		{3, false},
		{0, true},
		{-4, true},
		{-5, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("IsEven(%d)", tt.input), func(t *testing.T) {
			result := IsEven(tt.input)
			if result != tt.expected {
				t.Errorf("IsEven(%d) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFactorial(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
		wantErr  bool
	}{
		{"factorial of 0", 0, 1, false},
		{"factorial of 1", 1, 1, false},
		{"factorial of 5", 5, 120, false},
		{"negative number", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Factorial(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Factorial(%d) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Factorial(%d) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("Factorial(%d) = %d; want %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// ===== String Tests =====

func TestReverse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "olleh"},
		{"Go", "oG"},
		{"", ""},
		{"a", "a"},
		{"12345", "54321"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Reverse(tt.input)
			if result != tt.expected {
				t.Errorf("Reverse(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"radar", true},
		{"hello", false},
		{"level", true},
		{"Go", false},
		{"", true},
		{"a", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := IsPalindrome(tt.input)
			if result != tt.expected {
				t.Errorf("IsPalindrome(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// ===== User Tests =====

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		wantErr  bool
	}{
		{
			"valid user",
			&User{Username: "alice", Email: "alice@example.com", Age: 25},
			false,
		},
		{
			"missing username",
			&User{Username: "", Email: "test@example.com", Age: 25},
			true,
		},
		{
			"missing email",
			&User{Username: "bob", Email: "", Age: 25},
			true,
		},
		{
			"negative age",
			&User{Username: "charlie", Email: "charlie@example.com", Age: -5},
			true,
		},
		{
			"age too high",
			&User{Username: "david", Email: "david@example.com", Age: 200},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService(t *testing.T) {
	service := NewUserService()

	// Test Create
	user, err := service.Create("alice", "alice@example.com", 25)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("User ID = %d; want 1", user.ID)
	}

	// Test GetByID
	retrieved, err := service.GetByID(user.ID)
	if err != nil {
		t.Errorf("GetByID(%d) error = %v", user.ID, err)
	}
	if retrieved.Username != "alice" {
		t.Errorf("Username = %s; want alice", retrieved.Username)
	}

	// Test Count
	if count := service.Count(); count != 1 {
		t.Errorf("Count() = %d; want 1", count)
	}

	// Test Update
	updated, err := service.Update(user.ID, "alice_new", "alice.new@example.com", 26)
	if err != nil {
		t.Errorf("Update error = %v", err)
	}
	if updated.Age != 26 {
		t.Errorf("Updated age = %d; want 26", updated.Age)
	}

	// Test Delete
	err = service.Delete(user.ID)
	if err != nil {
		t.Errorf("Delete error = %v", err)
	}

	if count := service.Count(); count != 0 {
		t.Errorf("Count after delete = %d; want 0", count)
	}
}

// ===== Benchmark Tests =====

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(100, 200)
	}
}

func BenchmarkFactorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Factorial(10)
	}
}

func BenchmarkReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Reverse("hello world this is a test string")
	}
}

// ===== Example Tests (for documentation) =====

func ExampleAdd() {
	result := Add(2, 3)
	fmt.Println(result)
	// Output: 5
}

func ExampleReverse() {
	result := Reverse("hello")
	fmt.Println(result)
	// Output: olleh
}
*/

// สรุป Testing Best Practices:
//
// 1. ตั้งชื่อ test function: Test + ชื่อฟังก์ชัน
// 2. ใช้ table-driven tests สำหรับหลาย test cases
// 3. ใช้ t.Run() สำหรับ subtests
// 4. Test ทั้ง happy path และ error cases
// 5. ใช้ t.Errorf() สำหรับ failures
// 6. ใช้ t.Fatalf() เมื่อต้องการหยุด test ทันที
// 7. เขียน benchmarks เพื่อวัด performance
// 8. ใช้ -cover flag เพื่อดู test coverage
// 9. ใช้ -race flag เพื่อหา race conditions
// 10. เขียน example tests สำหรับ documentation
