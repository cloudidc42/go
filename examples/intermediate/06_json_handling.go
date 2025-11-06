// ===== ตัวอย่าง: JSON Handling ใน Go =====
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// ===== Part 1: Basic Structs =====

type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	Country string `json:"country"`
	ZipCode string `json:"zip_code"`
}

type Employee struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Address   Address   `json:"address"`
	Position  string    `json:"position"`
	Salary    float64   `json:"salary"`
	HireDate  time.Time `json:"hire_date"`
	Active    bool      `json:"active"`
}

// ===== Part 2: Nested Structs =====

type Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	InStock     bool     `json:"in_stock"`
	Tags        []string `json:"tags"`
	CreatedAt   string   `json:"created_at"`
}

type Order struct {
	OrderID    string    `json:"order_id"`
	CustomerID int       `json:"customer_id"`
	Products   []Product `json:"products"`
	Total      float64   `json:"total"`
	Status     string    `json:"status"`
	OrderDate  string    `json:"order_date"`
}

// ===== Part 3: JSON Tags และ Options =====

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // ไม่ export ไป JSON
	Email     string    `json:"email,omitempty"` // ไม่แสดงถ้าเป็น empty
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	IsAdmin   bool      `json:"is_admin"`
}

// ===== Part 4: Custom JSON Marshal =====

type CustomDate struct {
	time.Time
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	formatted := cd.Format("2006-01-02")
	return json.Marshal(formatted)
}

func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	cd.Time = t
	return nil
}

type Event struct {
	Name string     `json:"name"`
	Date CustomDate `json:"date"`
}

// ===== Main Function =====

func main() {
	fmt.Println("===== JSON Handling ใน Go =====\n")

	// ===== Part 1: Marshal (Struct → JSON) =====
	fmt.Println("--- 1. Marshal: Struct to JSON ---")

	person := Person{
		Name:  "Alice Smith",
		Age:   30,
		Email: "alice@example.com",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(person)
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		return
	}

	fmt.Printf("JSON: %s\n", string(jsonData))

	// Marshal with indentation (pretty print)
	jsonPretty, err := json.MarshalIndent(person, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		return
	}

	fmt.Println("\nPretty JSON:")
	fmt.Println(string(jsonPretty))
	fmt.Println()

	// ===== Part 2: Unmarshal (JSON → Struct) =====
	fmt.Println("--- 2. Unmarshal: JSON to Struct ---")

	jsonStr := `{
		"name": "Bob Johnson",
		"age": 25,
		"email": "bob@example.com"
	}`

	var person2 Person
	err = json.Unmarshal([]byte(jsonStr), &person2)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v\n", err)
		return
	}

	fmt.Printf("Parsed Person: %+v\n", person2)
	fmt.Printf("Name: %s, Age: %d\n\n", person2.Name, person2.Age)

	// ===== Part 3: Complex Nested Struct =====
	fmt.Println("--- 3. Complex Nested Struct ---")

	employee := Employee{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@company.com",
		Age:       35,
		Address: Address{
			Street:  "123 Main St",
			City:    "Bangkok",
			Country: "Thailand",
			ZipCode: "10110",
		},
		Position: "Senior Developer",
		Salary:   75000.00,
		HireDate: time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
		Active:   true,
	}

	empJSON, _ := json.MarshalIndent(employee, "", "  ")
	fmt.Println("Employee JSON:")
	fmt.Println(string(empJSON))
	fmt.Println()

	// ===== Part 4: Slice of Structs =====
	fmt.Println("--- 4. Slice of Structs ---")

	people := []Person{
		{Name: "Alice", Age: 30, Email: "alice@example.com"},
		{Name: "Bob", Age: 25, Email: "bob@example.com"},
		{Name: "Carol", Age: 28, Email: "carol@example.com"},
	}

	peopleJSON, _ := json.MarshalIndent(people, "", "  ")
	fmt.Println("People JSON:")
	fmt.Println(string(peopleJSON))
	fmt.Println()

	// ===== Part 5: Map to JSON =====
	fmt.Println("--- 5. Map to JSON ---")

	data := map[string]interface{}{
		"name":    "Test User",
		"age":     30,
		"active":  true,
		"balance": 1234.56,
		"tags":    []string{"vip", "premium"},
	}

	mapJSON, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println("Map JSON:")
	fmt.Println(string(mapJSON))
	fmt.Println()

	// ===== Part 6: JSON to Map =====
	fmt.Println("--- 6. JSON to Map ---")

	jsonData2 := `{
		"status": "success",
		"code": 200,
		"data": {
			"user_id": 123,
			"username": "testuser"
		}
	}`

	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsonData2), &result)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v\n", err)
		return
	}

	fmt.Printf("Parsed Map: %+v\n", result)
	fmt.Printf("Status: %v\n", result["status"])
	fmt.Printf("Code: %v\n", result["code"])

	// Type assertion for nested map
	if dataMap, ok := result["data"].(map[string]interface{}); ok {
		fmt.Printf("User ID: %v\n", dataMap["user_id"])
		fmt.Printf("Username: %v\n", dataMap["username"])
	}
	fmt.Println()

	// ===== Part 7: Omitempty Tag =====
	fmt.Println("--- 7. Omitempty Tag ---")

	user1 := User{
		ID:        1,
		Username:  "alice",
		Password:  "secret123", // จะไม่ปรากฏใน JSON
		Email:     "alice@example.com",
		FullName:  "Alice Smith",
		CreatedAt: time.Now(),
		// UpdatedAt ไม่กำหนด (omitempty จะทำให้ไม่แสดง)
		IsAdmin: false,
	}

	user2 := User{
		ID:        2,
		Username:  "bob",
		Password:  "secret456",
		Email:     "", // empty จะไม่แสดง (omitempty)
		FullName:  "Bob Johnson",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(), // มีค่า จะแสดง
		IsAdmin:   true,
	}

	user1JSON, _ := json.MarshalIndent(user1, "", "  ")
	fmt.Println("User 1 JSON (Email shown, UpdatedAt hidden):")
	fmt.Println(string(user1JSON))

	user2JSON, _ := json.MarshalIndent(user2, "", "  ")
	fmt.Println("\nUser 2 JSON (Email hidden, UpdatedAt shown):")
	fmt.Println(string(user2JSON))
	fmt.Println()

	// ===== Part 8: Working with Files =====
	fmt.Println("--- 8. Write JSON to File ---")

	products := []Product{
		{
			ID:          1,
			Name:        "Laptop",
			Description: "High-performance laptop",
			Price:       999.99,
			InStock:     true,
			Tags:        []string{"electronics", "computers"},
			CreatedAt:   time.Now().Format(time.RFC3339),
		},
		{
			ID:          2,
			Name:        "Mouse",
			Description: "Wireless mouse",
			Price:       29.99,
			InStock:     true,
			Tags:        []string{"electronics", "accessories"},
			CreatedAt:   time.Now().Format(time.RFC3339),
		},
	}

	// Write to file
	productsJSON, _ := json.MarshalIndent(products, "", "  ")
	err = ioutil.WriteFile("products.json", productsJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
	} else {
		fmt.Println("✓ Created products.json")
	}

	// Read from file
	fileData, err := ioutil.ReadFile("products.json")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	} else {
		var loadedProducts []Product
		err = json.Unmarshal(fileData, &loadedProducts)
		if err != nil {
			fmt.Printf("Error unmarshaling: %v\n", err)
		} else {
			fmt.Printf("✓ Loaded %d products from file\n", len(loadedProducts))
			for _, p := range loadedProducts {
				fmt.Printf("  - %s: $%.2f\n", p.Name, p.Price)
			}
		}
	}
	fmt.Println()

	// ===== Part 9: JSON Encoder/Decoder (Streaming) =====
	fmt.Println("--- 9. JSON Encoder/Decoder ---")

	// Encoder (write directly to file)
	file, err := os.Create("employees.json")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	employees := []Employee{
		{
			ID:        1,
			FirstName: "Alice",
			LastName:  "Smith",
			Email:     "alice@company.com",
			Age:       30,
			Address:   Address{City: "Bangkok", Country: "Thailand"},
			Position:  "Developer",
			Salary:    60000,
			HireDate:  time.Now(),
			Active:    true,
		},
		{
			ID:        2,
			FirstName: "Bob",
			LastName:  "Johnson",
			Email:     "bob@company.com",
			Age:       35,
			Address:   Address{City: "Chiang Mai", Country: "Thailand"},
			Position:  "Manager",
			Salary:    80000,
			HireDate:  time.Now(),
			Active:    true,
		},
	}

	err = encoder.Encode(employees)
	if err != nil {
		fmt.Printf("Error encoding: %v\n", err)
	} else {
		fmt.Println("✓ Encoded employees to file")
	}
	file.Close()

	// Decoder (read from file)
	file2, err := os.Open("employees.json")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file2.Close()

	decoder := json.NewDecoder(file2)

	var loadedEmployees []Employee
	err = decoder.Decode(&loadedEmployees)
	if err != nil {
		fmt.Printf("Error decoding: %v\n", err)
	} else {
		fmt.Printf("✓ Decoded %d employees\n", len(loadedEmployees))
		for _, emp := range loadedEmployees {
			fmt.Printf("  - %s %s (%s)\n", emp.FirstName, emp.LastName, emp.Position)
		}
	}
	fmt.Println()

	// ===== Part 10: Custom JSON Marshal/Unmarshal =====
	fmt.Println("--- 10. Custom JSON Marshal/Unmarshal ---")

	event := Event{
		Name: "Go Conference 2024",
		Date: CustomDate{time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC)},
	}

	eventJSON, _ := json.MarshalIndent(event, "", "  ")
	fmt.Println("Event JSON (custom date format):")
	fmt.Println(string(eventJSON))

	// Unmarshal
	jsonEvent := `{
		"name": "Golang Meetup",
		"date": "2024-11-20"
	}`

	var parsedEvent Event
	err = json.Unmarshal([]byte(jsonEvent), &parsedEvent)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v\n", err)
	} else {
		fmt.Printf("\nParsed Event: %s on %s\n", parsedEvent.Name, parsedEvent.Date.Format("2006-01-02"))
	}

	// ===== Part 11: Error Handling =====
	fmt.Println("\n--- 11. Error Handling ---")

	invalidJSON := `{"name": "Test", "age": "invalid"}`
	var testPerson Person

	err = json.Unmarshal([]byte(invalidJSON), &testPerson)
	if err != nil {
		fmt.Printf("✓ Caught error: %v\n", err)
	}

	// ===== Part 12: Practical Example - API Response =====
	fmt.Println("\n--- 12. API Response Example ---")

	type APIResponse struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
		Error   string      `json:"error,omitempty"`
	}

	// Success response
	successResp := APIResponse{
		Success: true,
		Message: "User created successfully",
		Data: map[string]interface{}{
			"user_id": 12345,
			"username": "newuser",
		},
	}

	successJSON, _ := json.MarshalIndent(successResp, "", "  ")
	fmt.Println("Success Response:")
	fmt.Println(string(successJSON))

	// Error response
	errorResp := APIResponse{
		Success: false,
		Message: "Failed to create user",
		Error:   "Username already exists",
	}

	errorJSON, _ := json.MarshalIndent(errorResp, "", "  ")
	fmt.Println("\nError Response:")
	fmt.Println(string(errorJSON))

	// ===== Cleanup =====
	fmt.Println("\n--- Cleanup ---")
	os.Remove("products.json")
	os.Remove("employees.json")
	fmt.Println("✓ Cleaned up JSON files")
}

// วิธีรัน:
// go run 06_json_handling.go
//
// Key Points:
// - json.Marshal() - Struct → JSON
// - json.Unmarshal() - JSON → Struct
// - json.MarshalIndent() - Pretty print
// - json.NewEncoder/Decoder() - Streaming
//
// JSON Tags:
// - `json:"field_name"` - กำหนดชื่อ field ใน JSON
// - `json:"-"` - ไม่ export field นี้
// - `json:",omitempty"` - ไม่แสดงถ้าค่าเป็น zero value
// - `json:"name,omitempty"` - ทั้งสองอย่าง
//
// Best Practices:
// - ใช้ struct tags เพื่อควบคุม JSON output
// - ตรวจสอบ errors จาก Marshal/Unmarshal เสมอ
// - ใช้ interface{} สำหรับ dynamic JSON
// - ใช้ Encoder/Decoder สำหรับ streaming
// - Implement custom MarshalJSON สำหรับ custom formatting
//
// Common Patterns:
// - API responses: struct with Success, Message, Data, Error
// - Config files: map[string]interface{} or specific struct
// - Nested data: embedded structs
// - Arrays: []Struct
