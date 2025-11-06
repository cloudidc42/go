// ===== ตัวอย่าง: Arrays และ Slices =====
package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	fmt.Println("===== Arrays และ Slices ใน Go =====\n")

	// ===== Part 1: Arrays (ขนาดคงที่) =====
	fmt.Println("--- 1. Arrays (Array มีขนาดคงที่) ---")

	// ประกาศ array
	var arr1 [5]int
	fmt.Printf("Array เริ่มต้น: %v\n", arr1) // [0 0 0 0 0]

	// กำหนดค่าให้ array
	arr2 := [5]int{1, 2, 3, 4, 5}
	fmt.Printf("Array ที่กำหนดค่า: %v\n", arr2)

	// ให้ compiler นับขนาดเอง
	arr3 := [...]string{"Go", "Python", "JavaScript", "Rust"}
	fmt.Printf("Array ของ strings: %v (ขนาด: %d)\n", arr3, len(arr3))

	// เข้าถึงและแก้ไขสมาชิก
	arr3[1] = "Java"
	fmt.Printf("หลังแก้ไข: %v\n\n", arr3)

	// ===== Part 2: Slices (ขนาดไม่คงที่) =====
	fmt.Println("--- 2. Slices (ขนาดเปลี่ยนได้) ---")

	// ประกาศ slice
	var slice1 []int
	fmt.Printf("Slice ว่าง: %v (len=%d, cap=%d)\n", slice1, len(slice1), cap(slice1))

	// สร้าง slice ด้วย make
	slice2 := make([]int, 5)      // length = 5
	slice3 := make([]int, 3, 10)  // length = 3, capacity = 10
	fmt.Printf("slice2: %v (len=%d, cap=%d)\n", slice2, len(slice2), cap(slice2))
	fmt.Printf("slice3: %v (len=%d, cap=%d)\n", slice3, len(slice3), cap(slice3))

	// Slice literals
	slice4 := []string{"apple", "banana", "cherry"}
	fmt.Printf("Fruit slice: %v\n\n", slice4)

	// ===== Part 3: Append (เพิ่มสมาชิก) =====
	fmt.Println("--- 3. Append - เพิ่มสมาชิก ---")

	numbers := []int{1, 2, 3}
	fmt.Printf("เริ่มต้น: %v\n", numbers)

	numbers = append(numbers, 4)
	fmt.Printf("append 4: %v\n", numbers)

	numbers = append(numbers, 5, 6, 7)
	fmt.Printf("append หลายตัว: %v\n", numbers)

	// Append slice เข้าด้วยกัน
	moreNumbers := []int{8, 9, 10}
	numbers = append(numbers, moreNumbers...)
	fmt.Printf("append slice: %v\n\n", numbers)

	// ===== Part 4: Slicing (ตัดส่วน) =====
	fmt.Println("--- 4. Slicing Operations ---")

	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("Original: %v\n", data)

	fmt.Printf("data[2:5]: %v\n", data[2:5])     // [2 3 4]
	fmt.Printf("data[:4]: %v\n", data[:4])        // [0 1 2 3]
	fmt.Printf("data[6:]: %v\n", data[6:])        // [6 7 8 9]
	fmt.Printf("data[:]: %v\n\n", data[:])        // ทั้งหมด

	// ===== Part 5: Copy =====
	fmt.Println("--- 5. Copy Slices ---")

	source := []int{1, 2, 3, 4, 5}
	destination := make([]int, len(source))

	copied := copy(destination, source)
	fmt.Printf("Source: %v\n", source)
	fmt.Printf("Destination: %v\n", destination)
	fmt.Printf("Copied %d elements\n\n", copied)

	// ===== Part 6: Iteration (วนลูป) =====
	fmt.Println("--- 6. Iteration ---")

	fruits := []string{"apple", "banana", "cherry", "date"}

	// วิธีที่ 1: ใช้ range
	fmt.Println("Using range:")
	for index, fruit := range fruits {
		fmt.Printf("  [%d] %s\n", index, fruit)
	}

	// วิธีที่ 2: ใช้ for loop ธรรมดา
	fmt.Println("\nUsing traditional for:")
	for i := 0; i < len(fruits); i++ {
		fmt.Printf("  [%d] %s\n", i, fruits[i])
	}

	// วิธีที่ 3: เอาแค่ค่า (ไม่ต้องการ index)
	fmt.Println("\nValues only:")
	for _, fruit := range fruits {
		fmt.Printf("  - %s\n", fruit)
	}
	fmt.Println()

	// ===== Part 7: Multi-dimensional Slices =====
	fmt.Println("--- 7. Multi-dimensional Slices ---")

	// 2D Slice (Matrix)
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	fmt.Println("Matrix:")
	for i, row := range matrix {
		fmt.Printf("Row %d: %v\n", i, row)
	}

	// Dynamic 2D slice
	rows, cols := 3, 4
	grid := make([][]int, rows)
	for i := range grid {
		grid[i] = make([]int, cols)
		for j := range grid[i] {
			grid[i][j] = i*cols + j + 1
		}
	}

	fmt.Println("\nDynamic Grid:")
	for _, row := range grid {
		fmt.Println(row)
	}
	fmt.Println()

	// ===== Part 8: Sorting =====
	fmt.Println("--- 8. Sorting ---")

	nums := []int{64, 34, 25, 12, 22, 11, 90}
	fmt.Printf("Before sort: %v\n", nums)

	sort.Ints(nums)
	fmt.Printf("After sort: %v\n", nums)

	words := []string{"banana", "apple", "cherry", "date"}
	fmt.Printf("Before sort: %v\n", words)

	sort.Strings(words)
	fmt.Printf("After sort: %v\n", words)

	// Reverse sort
	sort.Sort(sort.Reverse(sort.IntSlice(nums)))
	fmt.Printf("Reverse sort: %v\n\n", nums)

	// ===== Part 9: Filtering และ Mapping =====
	fmt.Println("--- 9. Filtering และ Mapping ---")

	allNumbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Filter: เอาแค่เลขคู่
	evenNumbers := filterEven(allNumbers)
	fmt.Printf("Even numbers: %v\n", evenNumbers)

	// Filter: เอาแค่เลขคี่
	oddNumbers := filterOdd(allNumbers)
	fmt.Printf("Odd numbers: %v\n", oddNumbers)

	// Map: คูณ 2
	doubled := mapDouble(allNumbers)
	fmt.Printf("Doubled: %v\n", doubled)

	// Map: ยกกำลัง 2
	squared := mapSquare(allNumbers)
	fmt.Printf("Squared: %v\n\n", squared)

	// ===== Part 10: Practical Examples =====
	fmt.Println("--- 10. Practical Examples ---")

	// Example 1: Remove duplicates
	withDuplicates := []int{1, 2, 2, 3, 4, 4, 4, 5, 1, 2}
	unique := removeDuplicates(withDuplicates)
	fmt.Printf("Original: %v\n", withDuplicates)
	fmt.Printf("Unique: %v\n", unique)

	// Example 2: Find max and min
	testNums := []int{45, 12, 78, 34, 89, 23, 67}
	max, min := findMaxMin(testNums)
	fmt.Printf("\nNumbers: %v\n", testNums)
	fmt.Printf("Max: %d, Min: %d\n", max, min)

	// Example 3: Sum and Average
	sum, avg := sumAndAverage(testNums)
	fmt.Printf("Sum: %d, Average: %.2f\n", sum, avg)

	// Example 4: Reverse slice
	reversed := reverseSlice(testNums)
	fmt.Printf("\nOriginal: %v\n", testNums)
	fmt.Printf("Reversed: %v\n", reversed)

	// Example 5: Contains
	searchNums := []int{1, 2, 3, 4, 5}
	fmt.Printf("\n%v contains 3? %v\n", searchNums, contains(searchNums, 3))
	fmt.Printf("%v contains 10? %v\n", searchNums, contains(searchNums, 10))

	// Example 6: Join strings
	tags := []string{"golang", "programming", "backend", "web"}
	joined := strings.Join(tags, ", ")
	fmt.Printf("\nTags: %s\n", joined)

	// Example 7: Split strings
	sentence := "Go is an awesome programming language"
	wordsArray := strings.Fields(sentence)
	fmt.Printf("Words: %v (count: %d)\n", wordsArray, len(wordsArray))
}

// ===== Helper Functions =====

// Filter even numbers
func filterEven(numbers []int) []int {
	var result []int
	for _, num := range numbers {
		if num%2 == 0 {
			result = append(result, num)
		}
	}
	return result
}

// Filter odd numbers
func filterOdd(numbers []int) []int {
	var result []int
	for _, num := range numbers {
		if num%2 != 0 {
			result = append(result, num)
		}
	}
	return result
}

// Map: double each number
func mapDouble(numbers []int) []int {
	result := make([]int, len(numbers))
	for i, num := range numbers {
		result[i] = num * 2
	}
	return result
}

// Map: square each number
func mapSquare(numbers []int) []int {
	result := make([]int, len(numbers))
	for i, num := range numbers {
		result[i] = num * num
	}
	return result
}

// Remove duplicates
func removeDuplicates(numbers []int) []int {
	seen := make(map[int]bool)
	var result []int

	for _, num := range numbers {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}

	return result
}

// Find max and min
func findMaxMin(numbers []int) (max, min int) {
	if len(numbers) == 0 {
		return 0, 0
	}

	max, min = numbers[0], numbers[0]

	for _, num := range numbers {
		if num > max {
			max = num
		}
		if num < min {
			min = num
		}
	}

	return max, min
}

// Sum and average
func sumAndAverage(numbers []int) (sum int, avg float64) {
	if len(numbers) == 0 {
		return 0, 0
	}

	for _, num := range numbers {
		sum += num
	}

	avg = float64(sum) / float64(len(numbers))
	return sum, avg
}

// Reverse slice
func reverseSlice(numbers []int) []int {
	result := make([]int, len(numbers))
	for i, num := range numbers {
		result[len(numbers)-1-i] = num
	}
	return result
}

// Contains checks if slice contains a value
func contains(numbers []int, target int) bool {
	for _, num := range numbers {
		if num == target {
			return true
		}
	}
	return false
}

// วิธีรัน:
// go run 05_slices_arrays.go
//
// Key Points:
// - Arrays มีขนาดคงที่ (fixed size)
// - Slices มีขนาดเปลี่ยนได้ (dynamic)
// - ใช้ append() เพื่อเพิ่มสมาชิก
// - Slices คือ reference type (ระวังเรื่อง sharing)
// - ใช้ copy() เพื่อคัดลอกอย่างปลอดภัย
// - len() = จำนวนสมาชิก, cap() = ความจุ
//
// Common Operations:
// - Append: slice = append(slice, element)
// - Slicing: slice[start:end]
// - Copy: copy(dst, src)
// - Iteration: for i, v := range slice
// - Sort: sort.Ints(), sort.Strings()
