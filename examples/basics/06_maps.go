// ===== ตัวอย่าง: Maps (Hash Tables / Dictionaries) =====
package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println("===== Maps ใน Go =====\n")

	// ===== Part 1: การสร้าง Maps =====
	fmt.Println("--- 1. การสร้าง Maps ---")

	// วิธีที่ 1: Map literal
	ages := map[string]int{
		"Alice":   25,
		"Bob":     30,
		"Charlie": 35,
	}
	fmt.Printf("Ages map: %v\n", ages)

	// วิธีที่ 2: ใช้ make()
	scores := make(map[string]int)
	scores["Math"] = 95
	scores["English"] = 88
	scores["Science"] = 92
	fmt.Printf("Scores map: %v\n", scores)

	// วิธีที่ 3: Nil map (ยังใช้งานไม่ได้)
	var nilMap map[string]int
	fmt.Printf("Nil map: %v (is nil: %v)\n\n", nilMap, nilMap == nil)

	// ===== Part 2: การเข้าถึงและเพิ่มข้อมูล =====
	fmt.Println("--- 2. การเข้าถึงและเพิ่มข้อมูล ---")

	// เพิ่มข้อมูล
	ages["David"] = 28
	ages["Eve"] = 32
	fmt.Printf("After adding: %v\n", ages)

	// อ่านค่า
	aliceAge := ages["Alice"]
	fmt.Printf("Alice's age: %d\n", aliceAge)

	// อ่านค่าที่ไม่มี (จะได้ zero value)
	unknownAge := ages["Unknown"]
	fmt.Printf("Unknown person's age: %d\n", unknownAge)

	// ตรวจสอบว่ามี key หรือไม่ (comma ok idiom)
	if age, exists := ages["Alice"]; exists {
		fmt.Printf("Alice exists with age: %d\n", age)
	}

	if age, exists := ages["Unknown"]; !exists {
		fmt.Printf("Unknown does not exist (got: %d)\n\n", age)
	}

	// ===== Part 3: การลบข้อมูล =====
	fmt.Println("--- 3. การลบข้อมูล ---")

	fmt.Printf("Before delete: %v\n", ages)
	delete(ages, "Bob")
	fmt.Printf("After delete Bob: %v\n", ages)

	// ลบ key ที่ไม่มี (ไม่เกิด error)
	delete(ages, "NonExistent")
	fmt.Printf("After delete non-existent: %v\n\n", ages)

	// ===== Part 4: การวนลูป =====
	fmt.Println("--- 4. การวนลูป ---")

	products := map[string]float64{
		"Laptop":  899.99,
		"Mouse":   25.50,
		"Keyboard": 75.00,
		"Monitor": 249.99,
	}

	fmt.Println("Products and prices:")
	for name, price := range products {
		fmt.Printf("  %-10s: $%.2f\n", name, price)
	}

	// วนเอาแค่ keys
	fmt.Println("\nProduct names:")
	for name := range products {
		fmt.Printf("  - %s\n", name)
	}
	fmt.Println()

	// ===== Part 5: ความยาว (Length) =====
	fmt.Println("--- 5. Map Length ---")

	fmt.Printf("Number of ages: %d\n", len(ages))
	fmt.Printf("Number of products: %d\n", len(products))
	fmt.Printf("Number of items in nil map: %d\n\n", len(nilMap))

	// ===== Part 6: Nested Maps =====
	fmt.Println("--- 6. Nested Maps ---")

	// Map of maps (nested structure)
	users := map[string]map[string]interface{}{
		"user1": {
			"name":  "Alice",
			"age":   25,
			"email": "alice@example.com",
		},
		"user2": {
			"name":  "Bob",
			"age":   30,
			"email": "bob@example.com",
		},
	}

	fmt.Println("Users:")
	for userID, userData := range users {
		fmt.Printf("\n%s:\n", userID)
		for key, value := range userData {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
	fmt.Println()

	// ===== Part 7: Maps กับ Structs =====
	fmt.Println("--- 7. Maps with Structs ---")

	type Person struct {
		Name  string
		Age   int
		City  string
	}

	people := map[string]Person{
		"emp001": {Name: "Alice", Age: 25, City: "Bangkok"},
		"emp002": {Name: "Bob", Age: 30, City: "Chiang Mai"},
		"emp003": {Name: "Carol", Age: 28, City: "Phuket"},
	}

	fmt.Println("Employees:")
	for empID, person := range people {
		fmt.Printf("  %s: %s, %d years old, from %s\n",
			empID, person.Name, person.Age, person.City)
	}
	fmt.Println()

	// ===== Part 8: Sorted Iteration =====
	fmt.Println("--- 8. Sorted Iteration ---")

	// Maps ไม่มี order รับประกัน ต้อง sort keys เอง
	countries := map[string]string{
		"TH": "Thailand",
		"US": "United States",
		"JP": "Japan",
		"UK": "United Kingdom",
		"FR": "France",
	}

	// Sort keys
	var keys []string
	for k := range countries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("Countries (sorted by code):")
	for _, code := range keys {
		fmt.Printf("  %s: %s\n", code, countries[code])
	}
	fmt.Println()

	// ===== Part 9: Counting with Maps =====
	fmt.Println("--- 9. Counting with Maps ---")

	text := "hello world hello go world"
	wordCount := countWords(text)

	fmt.Printf("Text: %s\n", text)
	fmt.Println("Word counts:")
	for word, count := range wordCount {
		fmt.Printf("  '%s': %d\n", word, count)
	}
	fmt.Println()

	// Count characters
	charCount := countChars("hello")
	fmt.Println("Character counts in 'hello':")
	for char, count := range charCount {
		fmt.Printf("  '%c': %d\n", char, count)
	}
	fmt.Println()

	// ===== Part 10: Practical Examples =====
	fmt.Println("--- 10. Practical Examples ---")

	// Example 1: Phone book
	phoneBook := map[string]string{
		"Alice":   "081-234-5678",
		"Bob":     "082-345-6789",
		"Charlie": "083-456-7890",
	}

	fmt.Println("Phone Book:")
	displayPhoneBook(phoneBook)

	// Example 2: Inventory management
	inventory := map[string]int{
		"Apples":  50,
		"Bananas": 30,
		"Oranges": 25,
	}

	fmt.Println("\nInventory System:")
	displayInventory(inventory)

	// Add items
	addInventory(inventory, "Apples", 20)
	addInventory(inventory, "Grapes", 40)
	fmt.Println("\nAfter adding items:")
	displayInventory(inventory)

	// Remove items
	removeInventory(inventory, "Bananas", 10)
	fmt.Println("\nAfter removing 10 bananas:")
	displayInventory(inventory)

	// Example 3: Grouping data
	students := []string{"Alice", "Bob", "Alice", "Charlie", "Bob", "Alice", "David"}
	attendance := groupAndCount(students)
	fmt.Println("\nStudent Attendance:")
	for student, count := range attendance {
		fmt.Printf("  %s attended %d time(s)\n", student, count)
	}

	// Example 4: Cache/Memoization
	fmt.Println("\nFibonacci with Memoization:")
	cache := make(map[int]int)
	for i := 0; i <= 10; i++ {
		result := fibMemo(i, cache)
		fmt.Printf("  fib(%d) = %d\n", i, result)
	}
	fmt.Printf("Cache size: %d entries\n", len(cache))

	// Example 5: Set operations (using map)
	fmt.Println("\nSet Operations:")
	set1 := makeSet([]int{1, 2, 3, 4, 5})
	set2 := makeSet([]int{4, 5, 6, 7, 8})

	fmt.Printf("Set 1: %v\n", getSetElements(set1))
	fmt.Printf("Set 2: %v\n", getSetElements(set2))
	fmt.Printf("Intersection: %v\n", intersection(set1, set2))
	fmt.Printf("Union: %v\n", union(set1, set2))
}

// ===== Helper Functions =====

// Count words in text
func countWords(text string) map[string]int {
	counts := make(map[string]int)
	word := ""

	for _, char := range text + " " {
		if char == ' ' {
			if word != "" {
				counts[word]++
				word = ""
			}
		} else {
			word += string(char)
		}
	}

	return counts
}

// Count characters
func countChars(text string) map[rune]int {
	counts := make(map[rune]int)
	for _, char := range text {
		counts[char]++
	}
	return counts
}

// Display phone book
func displayPhoneBook(pb map[string]string) {
	for name, phone := range pb {
		fmt.Printf("  %-10s: %s\n", name, phone)
	}
}

// Display inventory
func displayInventory(inv map[string]int) {
	for item, quantity := range inv {
		fmt.Printf("  %-10s: %d units\n", item, quantity)
	}
}

// Add to inventory
func addInventory(inv map[string]int, item string, quantity int) {
	inv[item] += quantity
}

// Remove from inventory
func removeInventory(inv map[string]int, item string, quantity int) {
	if inv[item] >= quantity {
		inv[item] -= quantity
	} else {
		fmt.Printf("  Warning: Not enough %s in inventory\n", item)
	}
}

// Group and count occurrences
func groupAndCount(items []string) map[string]int {
	counts := make(map[string]int)
	for _, item := range items {
		counts[item]++
	}
	return counts
}

// Fibonacci with memoization
func fibMemo(n int, cache map[int]int) int {
	// Check cache
	if val, exists := cache[n]; exists {
		return val
	}

	// Base cases
	if n <= 1 {
		return n
	}

	// Calculate and cache
	result := fibMemo(n-1, cache) + fibMemo(n-2, cache)
	cache[n] = result

	return result
}

// Set operations using maps
func makeSet(numbers []int) map[int]bool {
	set := make(map[int]bool)
	for _, num := range numbers {
		set[num] = true
	}
	return set
}

func getSetElements(set map[int]bool) []int {
	var elements []int
	for num := range set {
		elements = append(elements, num)
	}
	sort.Ints(elements)
	return elements
}

func intersection(set1, set2 map[int]bool) []int {
	result := make(map[int]bool)
	for num := range set1 {
		if set2[num] {
			result[num] = true
		}
	}
	return getSetElements(result)
}

func union(set1, set2 map[int]bool) []int {
	result := make(map[int]bool)
	for num := range set1 {
		result[num] = true
	}
	for num := range set2 {
		result[num] = true
	}
	return getSetElements(result)
}

// วิธีรัน:
// go run 06_maps.go
//
// Key Points:
// - Maps เป็น reference type
// - Zero value ของ map คือ nil (ใช้งานไม่ได้)
// - ต้องสร้างด้วย make() หรือ map literal
// - การเข้าถึง key ที่ไม่มีจะได้ zero value
// - ใช้ comma ok idiom เพื่อตรวจสอบว่ามี key
// - Maps ไม่รับประกัน order ของ keys
// - ใช้ delete() เพื่อลบ key
//
// Common Patterns:
// - Counter/Frequency: map[string]int
// - Set: map[T]bool
// - Cache/Memoization: map[key]value
// - Index/Lookup: map[id]struct
// - Grouping: map[category][]items
