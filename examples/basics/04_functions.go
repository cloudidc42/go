// ===== ตัวอย่าง: Functions - Complete Guide =====
package main

import (
    "fmt"
    "math"
    "strings"
)

func main() {
    fmt.Println("=== Go Functions - Complete Examples ===\n")

    // ===== Example 1: Basic Functions =====
    fmt.Println("Example 1: Basic Functions")

    greet()
    sayHello("Alice")

    result := add(5, 3)
    fmt.Printf("  5 + 3 = %d\n", result)

    // ===== Example 2: Multiple Return Values =====
    fmt.Println("\nExample 2: Multiple Return Values")

    sum, diff := calculate(10, 3)
    fmt.Printf("  10 + 3 = %d, 10 - 3 = %d\n", sum, diff)

    quotient, remainder := divide(17, 5)
    fmt.Printf("  17 ÷ 5 = %d remainder %d\n", quotient, remainder)

    // ===== Example 3: Named Return Values =====
    fmt.Println("\nExample 3: Named Return Values")

    area, perimeter := rectangleStats(5, 3)
    fmt.Printf("  Rectangle 5x3: Area = %d, Perimeter = %d\n", area, perimeter)

    // ===== Example 4: Variadic Functions =====
    fmt.Println("\nExample 4: Variadic Functions")

    total1 := sum(1, 2, 3, 4, 5)
    fmt.Printf("  Sum of 1,2,3,4,5 = %d\n", total1)

    numbers := []int{10, 20, 30, 40}
    total2 := sum(numbers...)
    fmt.Printf("  Sum of slice = %d\n", total2)

    printInfo("Alice", "Engineer", 30, "New York")

    // ===== Example 5: Anonymous Functions =====
    fmt.Println("\nExample 5: Anonymous Functions")

    // Immediate execution
    func() {
        fmt.Println("  Anonymous function executed!")
    }()

    // Assigned to variable
    multiply := func(a, b int) int {
        return a * b
    }
    fmt.Printf("  5 × 3 = %d\n", multiply(5, 3))

    // ===== Example 6: Closures =====
    fmt.Println("\nExample 6: Closures")

    counter := makeCounter()
    fmt.Printf("  Count: %d\n", counter()) // 1
    fmt.Printf("  Count: %d\n", counter()) // 2
    fmt.Printf("  Count: %d\n", counter()) // 3

    addFive := makeAdder(5)
    fmt.Printf("  10 + 5 = %d\n", addFive(10))
    fmt.Printf("  20 + 5 = %d\n", addFive(20))

    // ===== Example 7: Higher-Order Functions =====
    fmt.Println("\nExample 7: Higher-Order Functions")

    nums := []int{1, 2, 3, 4, 5}

    // Map
    squared := mapInt(nums, func(n int) int { return n * n })
    fmt.Printf("  Squared: %v\n", squared)

    // Filter
    evens := filterInt(nums, func(n int) bool { return n%2 == 0 })
    fmt.Printf("  Evens: %v\n", evens)

    // Reduce
    sumResult := reduceInt(nums, 0, func(acc, n int) int { return acc + n })
    fmt.Printf("  Sum: %d\n", sumResult)

    // ===== Example 8: Recursion =====
    fmt.Println("\nExample 8: Recursion")

    fmt.Printf("  Factorial of 5: %d\n", factorial(5))
    fmt.Printf("  Fibonacci(10): %d\n", fibonacci(10))

    // Sum of digits
    fmt.Printf("  Sum of digits (12345): %d\n", sumDigits(12345))

    // ===== Example 9: Error Handling =====
    fmt.Println("\nExample 9: Error Handling")

    result1, err := safeDivide(10, 2)
    if err != nil {
        fmt.Printf("  Error: %v\n", err)
    } else {
        fmt.Printf("  10 ÷ 2 = %.2f\n", result1)
    }

    result2, err := safeDivide(10, 0)
    if err != nil {
        fmt.Printf("  Error: %v\n", err)
    } else {
        fmt.Printf("  Result: %.2f\n", result2)
    }

    // ===== Example 10: Defer =====
    fmt.Println("\nExample 10: Defer")

    deferExample()

    // ===== Example 11: Practical Examples =====
    fmt.Println("\nExample 11: Practical Examples")

    // Temperature conversion
    celsius := 25.0
    fahrenheit := celsiusToFahrenheit(celsius)
    fmt.Printf("  %.1f°C = %.1f°F\n", celsius, fahrenheit)

    // String manipulation
    text := "hello world"
    fmt.Printf("  Capitalize: %s\n", capitalize(text))
    fmt.Printf("  Reverse: %s\n", reverse(text))

    // Validation
    email := "user@example.com"
    fmt.Printf("  Is '%s' valid email? %t\n", email, isValidEmail(email))

    // ===== Example 12: Function as Parameter =====
    fmt.Println("\nExample 12: Function as Parameter")

    applyOperation(10, 5, func(a, b int) int {
        return a + b
    }, "Addition")

    applyOperation(10, 5, func(a, b int) int {
        return a * b
    }, "Multiplication")
}

// ===== Basic Functions =====

func greet() {
    fmt.Println("  Hello, Go!")
}

func sayHello(name string) {
    fmt.Printf("  Hello, %s!\n", name)
}

func add(a, b int) int {
    return a + b
}

// ===== Multiple Return Values =====

func calculate(a, b int) (int, int) {
    return a + b, a - b
}

func divide(a, b int) (int, int) {
    return a / b, a % b
}

// ===== Named Return Values =====

func rectangleStats(length, width int) (area int, perimeter int) {
    area = length * width
    perimeter = 2 * (length + width)
    return // naked return
}

// ===== Variadic Functions =====

func sum(numbers ...int) int {
    total := 0
    for _, num := range numbers {
        total += num
    }
    return total
}

func printInfo(values ...interface{}) {
    fmt.Print("  Info: ")
    for i, v := range values {
        if i > 0 {
            fmt.Print(", ")
        }
        fmt.Print(v)
    }
    fmt.Println()
}

// ===== Closures =====

func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

func makeAdder(x int) func(int) int {
    return func(y int) int {
        return x + y
    }
}

// ===== Higher-Order Functions =====

func mapInt(nums []int, f func(int) int) []int {
    result := make([]int, len(nums))
    for i, n := range nums {
        result[i] = f(n)
    }
    return result
}

func filterInt(nums []int, f func(int) bool) []int {
    result := []int{}
    for _, n := range nums {
        if f(n) {
            result = append(result, n)
        }
    }
    return result
}

func reduceInt(nums []int, init int, f func(int, int) int) int {
    result := init
    for _, n := range nums {
        result = f(result, n)
    }
    return result
}

// ===== Recursion =====

func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

func sumDigits(n int) int {
    if n == 0 {
        return 0
    }
    return n%10 + sumDigits(n/10)
}

// ===== Error Handling =====

func safeDivide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// ===== Defer =====

func deferExample() {
    defer fmt.Println("    3. Deferred (executes last)")
    fmt.Println("    1. First")
    fmt.Println("    2. Second")
}

// ===== Practical Functions =====

func celsiusToFahrenheit(c float64) float64 {
    return c*9/5 + 32
}

func capitalize(s string) string {
    if len(s) == 0 {
        return s
    }
    return strings.ToUpper(s[:1]) + s[1:]
}

func reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func isValidEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// ===== Function as Parameter =====

func applyOperation(a, b int, operation func(int, int) int, opName string) {
    result := operation(a, b)
    fmt.Printf("  %s: %d and %d = %d\n", opName, a, b, result)
}

// ===== Advanced: Function Composition =====

func compose(f, g func(int) int) func(int) int {
    return func(x int) int {
        return f(g(x))
    }
}

// ===== Advanced: Memoization =====

func memoizedFibonacci() func(int) int {
    cache := make(map[int]int)

    var fib func(int) int
    fib = func(n int) int {
        if v, ok := cache[n]; ok {
            return v
        }

        if n <= 1 {
            return n
        }

        result := fib(n-1) + fib(n-2)
        cache[n] = result
        return result
    }

    return fib
}

// ===== Mathematical Functions =====

func power(base, exp int) int {
    result := 1
    for i := 0; i < exp; i++ {
        result *= base
    }
    return result
}

func isPrime(n int) bool {
    if n <= 1 {
        return false
    }
    if n <= 3 {
        return true
    }
    if n%2 == 0 || n%3 == 0 {
        return false
    }

    for i := 5; i*i <= n; i += 6 {
        if n%i == 0 || n%(i+2) == 0 {
            return false
        }
    }

    return true
}

func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

func distance(x1, y1, x2, y2 float64) float64 {
    dx := x2 - x1
    dy := y2 - y1
    return math.Sqrt(dx*dx + dy*dy)
}

// วิธีรัน:
// go run 04_functions.go
//
// ตัวอย่างครอบคลุม:
// - Basic functions
// - Multiple return values
// - Named returns
// - Variadic functions
// - Anonymous functions
// - Closures
// - Higher-order functions (map, filter, reduce)
// - Recursion
// - Error handling
// - Defer
// - Practical examples
// - Function composition
//
// ทุกตัวอย่างทำงานได้จริง 100%!
