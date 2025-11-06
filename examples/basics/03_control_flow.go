// ===== ตัวอย่าง: Control Flow - All Patterns =====
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    fmt.Println("=== Go Control Flow - Complete Examples ===\n")

    // ===== Example 1: If-Else =====
    fmt.Println("Example 1: If-Else Statements")

    age := 18
    if age >= 18 {
        fmt.Printf("  Age %d: You can vote! ✓\n", age)
    } else {
        fmt.Printf("  Age %d: Too young to vote\n", age)
    }

    // If with short statement
    if num := rand.Intn(100); num > 50 {
        fmt.Printf("  Random number %d is greater than 50\n", num)
    } else {
        fmt.Printf("  Random number %d is less than or equal to 50\n", num)
    }

    // ===== Example 2: Switch =====
    fmt.Println("\nExample 2: Switch Statements")

    day := time.Now().Weekday()
    switch day {
    case time.Saturday, time.Sunday:
        fmt.Printf("  It's %s - Weekend! 🎉\n", day)
    case time.Monday:
        fmt.Printf("  It's %s - Start of the week\n", day)
    case time.Friday:
        fmt.Printf("  It's %s - TGIF!\n", day)
    default:
        fmt.Printf("  It's %s - Regular weekday\n", day)
    }

    // Switch without expression
    score := 85
    switch {
    case score >= 90:
        fmt.Printf("  Score %d: Grade A (Excellent!)\n", score)
    case score >= 80:
        fmt.Printf("  Score %d: Grade B (Very Good!)\n", score)
    case score >= 70:
        fmt.Printf("  Score %d: Grade C (Good)\n", score)
    case score >= 60:
        fmt.Printf("  Score %d: Grade D (Pass)\n", score)
    default:
        fmt.Printf("  Score %d: Grade F (Fail)\n", score)
    }

    // ===== Example 3: For Loops =====
    fmt.Println("\nExample 3: For Loops")

    // Standard for loop
    fmt.Println("  Counting 1 to 5:")
    for i := 1; i <= 5; i++ {
        fmt.Printf("    %d ", i)
    }
    fmt.Println()

    // While-style loop
    fmt.Println("  Countdown from 5:")
    count := 5
    for count > 0 {
        fmt.Printf("    %d... ", count)
        count--
    }
    fmt.Println("Blast off! 🚀")

    // For with range (slice)
    fmt.Println("  Iterating fruits:")
    fruits := []string{"Apple", "Banana", "Orange", "Mango"}
    for i, fruit := range fruits {
        fmt.Printf("    [%d] %s\n", i, fruit)
    }

    // For with range (map)
    fmt.Println("  Student scores:")
    scores := map[string]int{
        "Alice": 95,
        "Bob":   87,
        "Carol": 92,
    }
    for name, score := range scores {
        fmt.Printf("    %s: %d points\n", name, score)
    }

    // ===== Example 4: Break and Continue =====
    fmt.Println("\nExample 4: Break and Continue")

    fmt.Println("  Finding first number divisible by 7:")
    for i := 1; i <= 50; i++ {
        if i%7 == 0 {
            fmt.Printf("    Found: %d\n", i)
            break // Exit loop
        }
    }

    fmt.Println("  Odd numbers only (1-10):")
    fmt.Print("    ")
    for i := 1; i <= 10; i++ {
        if i%2 == 0 {
            continue // Skip even numbers
        }
        fmt.Printf("%d ", i)
    }
    fmt.Println()

    // ===== Example 5: Nested Loops =====
    fmt.Println("\nExample 5: Nested Loops - Multiplication Table")
    for i := 1; i <= 5; i++ {
        fmt.Print("    ")
        for j := 1; j <= 5; j++ {
            fmt.Printf("%3d ", i*j)
        }
        fmt.Println()
    }

    // ===== Example 6: Pattern Printing =====
    fmt.Println("\nExample 6: Pattern Printing")

    // Triangle
    fmt.Println("  Right Triangle:")
    for i := 1; i <= 5; i++ {
        fmt.Print("    ")
        for j := 1; j <= i; j++ {
            fmt.Print("* ")
        }
        fmt.Println()
    }

    // Pyramid
    fmt.Println("  Pyramid:")
    n := 5
    for i := 1; i <= n; i++ {
        // Spaces
        for j := 1; j <= n-i; j++ {
            fmt.Print(" ")
        }
        // Stars
        for j := 1; j <= 2*i-1; j++ {
            fmt.Print("*")
        }
        fmt.Println()
    }

    // ===== Example 7: Practical Examples =====
    fmt.Println("\nExample 7: Practical Examples")

    // FizzBuzz
    fmt.Println("  FizzBuzz (1-15):")
    fmt.Print("    ")
    for i := 1; i <= 15; i++ {
        switch {
        case i%15 == 0:
            fmt.Print("FizzBuzz ")
        case i%3 == 0:
            fmt.Print("Fizz ")
        case i%5 == 0:
            fmt.Print("Buzz ")
        default:
            fmt.Printf("%d ", i)
        }
    }
    fmt.Println()

    // Prime numbers
    fmt.Println("  Prime numbers (1-30):")
    fmt.Print("    ")
    for num := 2; num <= 30; num++ {
        isPrime := true
        for i := 2; i*i <= num; i++ {
            if num%i == 0 {
                isPrime = false
                break
            }
        }
        if isPrime {
            fmt.Printf("%d ", num)
        }
    }
    fmt.Println()

    // ===== Example 8: Menu System =====
    fmt.Println("\nExample 8: Simple Menu System")
    displayMenu()
}

func displayMenu() {
    choices := []string{
        "1. View Profile",
        "2. Edit Settings",
        "3. View Reports",
        "4. Exit",
    }

    fmt.Println("  Menu Options:")
    for _, choice := range choices {
        fmt.Printf("    %s\n", choice)
    }

    // Simulate user choice
    userChoice := 2
    fmt.Printf("  Selected: %s\n", choices[userChoice-1])

    switch userChoice {
    case 1:
        fmt.Println("    → Opening Profile...")
    case 2:
        fmt.Println("    → Opening Settings...")
    case 3:
        fmt.Println("    → Loading Reports...")
    case 4:
        fmt.Println("    → Goodbye!")
    default:
        fmt.Println("    → Invalid choice")
    }
}

// วิธีรัน:
// go run 03_control_flow.go
//
// Output จะแสดง:
// - If-Else examples
// - Switch examples (with/without expression)
// - For loop variations (standard, while-style, range)
// - Break and Continue
// - Nested loops
// - Pattern printing (triangle, pyramid)
// - Practical examples (FizzBuzz, Prime numbers)
// - Menu system
//
// ทุกตัวอย่างทำงานได้จริง 100%!
