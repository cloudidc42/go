# คู่มือการเขียนโปรแกรม Go - ส่วนที่ 2
## Step 91-250: Control Flow, Functions และ Data Structures

---

## Step 91-120: Control Flow

### Step 91-100: If-Else Statements

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // ===== Step 91: Basic If =====
    age := 18

    if age >= 18 {
        fmt.Println("You are an adult")
    }

    // ===== Step 92: If-Else =====
    score := 75

    if score >= 80 {
        fmt.Println("Grade: A")
    } else {
        fmt.Println("Grade: B or below")
    }

    // ===== Step 93: If-Else If-Else =====
    marks := 85

    if marks >= 90 {
        fmt.Println("Grade: A+")
    } else if marks >= 80 {
        fmt.Println("Grade: A")
    } else if marks >= 70 {
        fmt.Println("Grade: B")
    } else if marks >= 60 {
        fmt.Println("Grade: C")
    } else {
        fmt.Println("Grade: F")
    }

    // ===== Step 94: If with Short Statement =====
    // ประกาศตัวแปรได้ใน if statement
    if num := getNumber(); num > 0 {
        fmt.Printf("%d is positive\n", num)
    } else if num < 0 {
        fmt.Printf("%d is negative\n", num)
    } else {
        fmt.Printf("%d is zero\n", num)
    }
    // num ไม่สามารถใช้ได้นอก if scope

    // ===== Step 95: Nested If =====
    username := "admin"
    password := "pass123"

    if username == "admin" {
        if password == "pass123" {
            fmt.Println("Login successful")
        } else {
            fmt.Println("Wrong password")
        }
    } else {
        fmt.Println("User not found")
    }

    // ===== Step 96: Logical Operators in If =====
    age96 := 25
    hasLicense := true

    if age96 >= 18 && hasLicense {
        fmt.Println("Can drive")
    }

    isWeekend := false
    isHoliday := true

    if isWeekend || isHoliday {
        fmt.Println("No work today!")
    }

    // ===== Step 97: Checking Error =====
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Result:", result)
    }

    // ===== Step 98: Multiple Conditions =====
    temperature := 25
    humidity := 60

    if temperature >= 20 && temperature <= 30 && humidity < 70 {
        fmt.Println("Perfect weather!")
    }

    // ===== Step 99: Range Checking =====
    value := 15

    if value >= 10 && value <= 20 {
        fmt.Println("Value is in range [10-20]")
    }

    // ===== Step 100: Time-based Conditions =====
    hour := time.Now().Hour()

    if hour < 12 {
        fmt.Println("Good morning!")
    } else if hour < 17 {
        fmt.Println("Good afternoon!")
    } else {
        fmt.Println("Good evening!")
    }
}

func getNumber() int {
    return -5
}

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```

### Step 101-110: Switch Statements

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // ===== Step 101: Basic Switch =====
    day := 3

    switch day {
    case 1:
        fmt.Println("Monday")
    case 2:
        fmt.Println("Tuesday")
    case 3:
        fmt.Println("Wednesday")
    case 4:
        fmt.Println("Thursday")
    case 5:
        fmt.Println("Friday")
    case 6:
        fmt.Println("Saturday")
    case 7:
        fmt.Println("Sunday")
    default:
        fmt.Println("Invalid day")
    }

    // ===== Step 102: Multiple Cases =====
    month := 4

    switch month {
    case 1, 3, 5, 7, 8, 10, 12:
        fmt.Println("31 days")
    case 4, 6, 9, 11:
        fmt.Println("30 days")
    case 2:
        fmt.Println("28 or 29 days")
    default:
        fmt.Println("Invalid month")
    }

    // ===== Step 103: Switch with Expressions =====
    score := 85

    switch {
    case score >= 90:
        fmt.Println("Excellent!")
    case score >= 80:
        fmt.Println("Very Good!")
    case score >= 70:
        fmt.Println("Good")
    case score >= 60:
        fmt.Println("Pass")
    default:
        fmt.Println("Fail")
    }

    // ===== Step 104: Switch with Short Statement =====
    switch num := getRandomNumber(); {
    case num < 0:
        fmt.Printf("%d is negative\n", num)
    case num == 0:
        fmt.Printf("%d is zero\n", num)
    case num > 0:
        fmt.Printf("%d is positive\n", num)
    }

    // ===== Step 105: Type Switch =====
    var i interface{} = "hello"

    switch v := i.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    case bool:
        fmt.Printf("Boolean: %t\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }

    // ===== Step 106: Fallthrough =====
    number := 1

    switch number {
    case 1:
        fmt.Println("One")
        fallthrough // จะทำ case ถัดไปด้วย
    case 2:
        fmt.Println("Two")
    case 3:
        fmt.Println("Three")
    }

    // ===== Step 107: Switch on Type =====
    checkType(42)
    checkType("Hello")
    checkType(3.14)
    checkType(true)
    checkType([]int{1, 2, 3})

    // ===== Step 108: OS Detection =====
    switch os := runtime.GOOS; os {
    case "darwin":
        fmt.Println("macOS")
    case "linux":
        fmt.Println("Linux")
    case "windows":
        fmt.Println("Windows")
    default:
        fmt.Printf("Unknown OS: %s\n", os)
    }

    // ===== Step 109: Weekday/Weekend =====
    today := time.Now().Weekday()

    switch today {
    case time.Saturday, time.Sunday:
        fmt.Println("It's the weekend!")
    default:
        fmt.Println("It's a weekday")
    }

    // ===== Step 110: Advanced Switch Patterns =====
    processValue := func(val interface{}) {
        switch v := val.(type) {
        case nil:
            fmt.Println("nil value")
        case int:
            if v > 0 {
                fmt.Printf("Positive int: %d\n", v)
            } else {
                fmt.Printf("Non-positive int: %d\n", v)
            }
        case string:
            if len(v) > 0 {
                fmt.Printf("Non-empty string: %s\n", v)
            } else {
                fmt.Println("Empty string")
            }
        default:
            fmt.Printf("Unknown type: %T\n", v)
        }
    }

    processValue(42)
    processValue("Go")
    processValue(nil)
}

func getRandomNumber() int {
    return 5
}

func checkType(val interface{}) {
    switch val.(type) {
    case int:
        fmt.Printf("%v is an integer\n", val)
    case string:
        fmt.Printf("%v is a string\n", val)
    case float64:
        fmt.Printf("%v is a float64\n", val)
    case bool:
        fmt.Printf("%v is a boolean\n", val)
    default:
        fmt.Printf("%v is of type %T\n", val, val)
    }
}
```

### Step 111-120: Loops

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // ===== Step 111: Basic For Loop =====
    fmt.Println("Counting 1 to 5:")
    for i := 1; i <= 5; i++ {
        fmt.Println(i)
    }

    // ===== Step 112: For as While =====
    count := 0
    for count < 5 {
        fmt.Println("Count:", count)
        count++
    }

    // ===== Step 113: Infinite Loop =====
    /*
    for {
        fmt.Println("Infinite loop - use Ctrl+C to stop")
        time.Sleep(1 * time.Second)
    }
    */

    // ===== Step 114: Break =====
    for i := 1; i <= 10; i++ {
        if i == 5 {
            break // หยุด loop
        }
        fmt.Print(i, " ")
    }
    fmt.Println()

    // ===== Step 115: Continue =====
    for i := 1; i <= 10; i++ {
        if i%2 == 0 {
            continue // ข้ามไปรอบถัดไป
        }
        fmt.Print(i, " ") // แสดงเฉพาะเลขคี่
    }
    fmt.Println()

    // ===== Step 116: Nested Loops =====
    fmt.Println("\nMultiplication Table:")
    for i := 1; i <= 5; i++ {
        for j := 1; j <= 5; j++ {
            fmt.Printf("%3d ", i*j)
        }
        fmt.Println()
    }

    // ===== Step 117: Range over Slice =====
    fruits := []string{"Apple", "Banana", "Cherry", "Date"}

    for index, fruit := range fruits {
        fmt.Printf("%d: %s\n", index, fruit)
    }

    // Only values
    for _, fruit := range fruits {
        fmt.Println("-", fruit)
    }

    // Only indices
    for index := range fruits {
        fmt.Println("Index:", index)
    }

    // ===== Step 118: Range over Map =====
    ages := map[string]int{
        "Alice": 25,
        "Bob":   30,
        "Carol": 28,
    }

    for name, age := range ages {
        fmt.Printf("%s is %d years old\n", name, age)
    }

    // ===== Step 119: Range over String =====
    text := "Hello, 世界"

    // Iterate by runes
    for index, runeValue := range text {
        fmt.Printf("Index %d: %c (Unicode: U+%04X)\n",
            index, runeValue, runeValue)
    }

    // ===== Step 120: Loop Patterns =====

    // Pattern 1: Sum
    numbers := []int{1, 2, 3, 4, 5}
    sum := 0
    for _, num := range numbers {
        sum += num
    }
    fmt.Printf("Sum: %d\n", sum)

    // Pattern 2: Find
    target := 3
    found := false
    for _, num := range numbers {
        if num == target {
            found = true
            break
        }
    }
    fmt.Printf("Found %d: %t\n", target, found)

    // Pattern 3: Filter
    evens := []int{}
    for _, num := range numbers {
        if num%2 == 0 {
            evens = append(evens, num)
        }
    }
    fmt.Printf("Even numbers: %v\n", evens)

    // Pattern 4: Transform
    doubled := make([]int, len(numbers))
    for i, num := range numbers {
        doubled[i] = num * 2
    }
    fmt.Printf("Doubled: %v\n", doubled)

    // Pattern 5: Count with timeout
    counter := 0
    timeout := time.After(2 * time.Second)

    for {
        select {
        case <-timeout:
            fmt.Printf("\nTimeout! Counted to %d\n", counter)
            goto done
        default:
            counter++
            time.Sleep(100 * time.Millisecond)
        }
    }

done:
    fmt.Println("Loop completed")
}
```

---

## Step 121-150: Functions

### Step 121-130: Basic Functions

```go
package main

import "fmt"

func main() {
    // ===== Step 121: Simple Function =====
    greet()

    // ===== Step 122: Function with Parameters =====
    sayHello("Alice")
    sayHello("Bob")

    // ===== Step 123: Function with Return Value =====
    result := add(5, 3)
    fmt.Printf("5 + 3 = %d\n", result)

    // ===== Step 124: Multiple Parameters =====
    fullName := getFullName("John", "Doe")
    fmt.Println("Full name:", fullName)

    // ===== Step 125: Multiple Return Values =====
    sum, diff := calculate(10, 3)
    fmt.Printf("Sum: %d, Difference: %d\n", sum, diff)

    // ===== Step 126: Named Return Values =====
    area, perimeter := rectangleStats(5, 3)
    fmt.Printf("Rectangle: Area=%d, Perimeter=%d\n", area, perimeter)

    // ===== Step 127: Variadic Functions =====
    total := sum(1, 2, 3, 4, 5)
    fmt.Println("Sum:", total)

    numbers := []int{10, 20, 30}
    total2 := sum(numbers...) // Spread operator
    fmt.Println("Sum of slice:", total2)

    // ===== Step 128: Anonymous Functions =====
    // ประกาศและเรียกใช้ทันที
    func() {
        fmt.Println("Anonymous function executed!")
    }()

    // เก็บไว้ในตัวแปร
    multiply := func(a, b int) int {
        return a * b
    }
    fmt.Println("5 × 3 =", multiply(5, 3))

    // ===== Step 129: Closures =====
    nextNumber := counter()
    fmt.Println(nextNumber()) // 1
    fmt.Println(nextNumber()) // 2
    fmt.Println(nextNumber()) // 3

    // ===== Step 130: Recursion =====
    fmt.Printf("Factorial of 5: %d\n", factorial(5))
    fmt.Printf("Fibonacci(10): %d\n", fibonacci(10))
}

// Step 121
func greet() {
    fmt.Println("Hello, World!")
}

// Step 122
func sayHello(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

// Step 123
func add(a, b int) int {
    return a + b
}

// Step 124
func getFullName(firstName, lastName string) string {
    return firstName + " " + lastName
}

// Step 125
func calculate(a, b int) (int, int) {
    return a + b, a - b
}

// Step 126
func rectangleStats(length, width int) (area int, perimeter int) {
    area = length * width
    perimeter = 2 * (length + width)
    return // naked return
}

// Step 127
func sum(numbers ...int) int {
    total := 0
    for _, num := range numbers {
        total += num
    }
    return total
}

// Step 129
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

// Step 130
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
```

### Step 131-140: Advanced Functions

```go
package main

import (
    "errors"
    "fmt"
)

func main() {
    // ===== Step 131: Error Handling =====
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("10 ÷ 2 = %.2f\n", result)
    }

    result2, err := divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("Result: %.2f\n", result2)
    }

    // ===== Step 132: Defer =====
    deferExample()

    // ===== Step 133: Defer Multiple =====
    multipleDeferExample()

    // ===== Step 134: Panic and Recover =====
    safeDivide(10, 0)
    fmt.Println("Program continues...")

    // ===== Step 135: Function as Parameter =====
    numbers := []int{1, 2, 3, 4, 5}

    applyFunc(numbers, func(n int) int {
        return n * n
    })

    // ===== Step 136: Function as Return Value =====
    addFunc := makeAdder(10)
    fmt.Println("10 + 5 =", addFunc(5))
    fmt.Println("10 + 15 =", addFunc(15))

    // ===== Step 137: Method Chaining =====
    calc := NewCalculator(10)
    result137 := calc.Add(5).Multiply(2).Subtract(10).GetValue()
    fmt.Printf("((10 + 5) × 2) - 10 = %d\n", result137)

    // ===== Step 138: Callback Functions =====
    processData([]int{1, 2, 3}, func(n int) {
        fmt.Printf("Processing: %d\n", n)
    })

    // ===== Step 139: Generic-like Functions =====
    printSlice([]int{1, 2, 3})
    printSlice([]string{"a", "b", "c"})

    // ===== Step 140: Memoization =====
    fib := memoizedFibonacci()
    fmt.Println("Fib(40):", fib(40))
    fmt.Println("Fib(40) again (cached):", fib(40))
}

// Step 131
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Step 132
func deferExample() {
    defer fmt.Println("3. This executes last (deferred)")
    fmt.Println("1. This executes first")
    fmt.Println("2. This executes second")
}

// Step 133
func multipleDeferExample() {
    fmt.Println("Start")
    defer fmt.Println("1st defer")
    defer fmt.Println("2nd defer")
    defer fmt.Println("3rd defer")
    fmt.Println("End")
    // Output order: Start, End, 3rd defer, 2nd defer, 1st defer
}

// Step 134
func safeDivide(a, b int) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()

    if b == 0 {
        panic("cannot divide by zero")
    }

    fmt.Printf("%d ÷ %d = %d\n", a, b, a/b)
}

// Step 135
func applyFunc(numbers []int, f func(int) int) {
    for i, n := range numbers {
        numbers[i] = f(n)
    }
    fmt.Println("After apply:", numbers)
}

// Step 136
func makeAdder(x int) func(int) int {
    return func(y int) int {
        return x + y
    }
}

// Step 137
type Calculator struct {
    value int
}

func NewCalculator(initial int) *Calculator {
    return &Calculator{value: initial}
}

func (c *Calculator) Add(n int) *Calculator {
    c.value += n
    return c
}

func (c *Calculator) Subtract(n int) *Calculator {
    c.value -= n
    return c
}

func (c *Calculator) Multiply(n int) *Calculator {
    c.value *= n
    return c
}

func (c *Calculator) GetValue() int {
    return c.value
}

// Step 138
func processData(data []int, callback func(int)) {
    for _, item := range data {
        callback(item)
    }
}

// Step 139
func printSlice(s interface{}) {
    fmt.Printf("Slice: %v, Type: %T\n", s, s)
}

// Step 140
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
```

### Step 141-150: Function Patterns

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // ===== Step 141: Decorator Pattern =====
    slow := slowFunction
    timed := timeFunction(slow)
    timed("test")

    // ===== Step 142: Pipeline Pattern =====
    numbers := []int{1, 2, 3, 4, 5}
    result := pipeline(
        numbers,
        double,
        addOne,
        square,
    )
    fmt.Println("Pipeline result:", result)

    // ===== Step 143: Option Pattern =====
    server := NewServer(
        WithPort(8080),
        WithHost("localhost"),
        WithTimeout(30),
    )
    server.Print()

    // ===== Step 144: Builder Pattern =====
    user := NewUserBuilder().
        SetName("John Doe").
        SetEmail("john@example.com").
        SetAge(30).
        Build()
    fmt.Printf("User: %+v\n", user)

    // ===== Step 145: Retry Pattern =====
    err := retry(3, time.Second, func() error {
        fmt.Println("Attempting operation...")
        return fmt.Errorf("simulated error")
    })
    if err != nil {
        fmt.Println("Failed after retries:", err)
    }

    // ===== Step 146: Curry Pattern =====
    add := func(a int) func(int) int {
        return func(b int) int {
            return a + b
        }
    }
    add5 := add(5)
    fmt.Println("5 + 10 =", add5(10))
    fmt.Println("5 + 20 =", add5(20))

    // ===== Step 147: Map, Filter, Reduce =====
    nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // Map
    squared := mapInt(nums, func(n int) int { return n * n })
    fmt.Println("Squared:", squared)

    // Filter
    evens := filterInt(nums, func(n int) bool { return n%2 == 0 })
    fmt.Println("Evens:", evens)

    // Reduce
    sum := reduceInt(nums, 0, func(acc, n int) int { return acc + n })
    fmt.Println("Sum:", sum)

    // ===== Step 148: Composition =====
    addTwo := func(x int) int { return x + 2 }
    multiplyThree := func(x int) int { return x * 3 }

    composed := compose(multiplyThree, addTwo)
    fmt.Println("Composed (5):", composed(5)) // (5 + 2) * 3 = 21

    // ===== Step 149: Partial Application =====
    multiply := func(a, b, c int) int { return a * b * c }
    partial := partial3to2(multiply, 2)
    result149 := partial(3, 4) // 2 * 3 * 4 = 24
    fmt.Println("Partial result:", result149)

    // ===== Step 150: Function Factory =====
    validators := map[string]func(string) bool{
        "email":    emailValidator(),
        "phone":    phoneValidator(),
        "username": usernameValidator(),
    }

    testEmail := "user@example.com"
    testPhone := "123-456-7890"

    for name, validator := range validators {
        fmt.Printf("%s validator: %t\n", name, validator(testEmail))
    }
}

// Step 141
func timeFunction(f func(string)) func(string) {
    return func(s string) {
        start := time.Now()
        f(s)
        fmt.Printf("Execution time: %v\n", time.Since(start))
    }
}

func slowFunction(s string) {
    time.Sleep(100 * time.Millisecond)
    fmt.Println("Processing:", s)
}

// Step 142
func pipeline(nums []int, funcs ...func([]int) []int) []int {
    result := nums
    for _, f := range funcs {
        result = f(result)
    }
    return result
}

func double(nums []int) []int {
    result := make([]int, len(nums))
    for i, n := range nums {
        result[i] = n * 2
    }
    return result
}

func addOne(nums []int) []int {
    result := make([]int, len(nums))
    for i, n := range nums {
        result[i] = n + 1
    }
    return result
}

func square(nums []int) []int {
    result := make([]int, len(nums))
    for i, n := range nums {
        result[i] = n * n
    }
    return result
}

// Step 143
type Server struct {
    host    string
    port    int
    timeout int
}

type ServerOption func(*Server)

func NewServer(opts ...ServerOption) *Server {
    s := &Server{
        host:    "0.0.0.0",
        port:    80,
        timeout: 10,
    }

    for _, opt := range opts {
        opt(s)
    }

    return s
}

func WithHost(host string) ServerOption {
    return func(s *Server) {
        s.host = host
    }
}

func WithPort(port int) ServerOption {
    return func(s *Server) {
        s.port = port
    }
}

func WithTimeout(timeout int) ServerOption {
    return func(s *Server) {
        s.timeout = timeout
    }
}

func (s *Server) Print() {
    fmt.Printf("Server{host: %s, port: %d, timeout: %d}\n",
        s.host, s.port, s.timeout)
}

// Step 144
type User struct {
    name  string
    email string
    age   int
}

type UserBuilder struct {
    user User
}

func NewUserBuilder() *UserBuilder {
    return &UserBuilder{}
}

func (b *UserBuilder) SetName(name string) *UserBuilder {
    b.user.name = name
    return b
}

func (b *UserBuilder) SetEmail(email string) *UserBuilder {
    b.user.email = email
    return b
}

func (b *UserBuilder) SetAge(age int) *UserBuilder {
    b.user.age = age
    return b
}

func (b *UserBuilder) Build() User {
    return b.user
}

// Step 145
func retry(attempts int, sleep time.Duration, f func() error) error {
    var err error
    for i := 0; i < attempts; i++ {
        if err = f(); err == nil {
            return nil
        }
        fmt.Printf("Attempt %d failed, retrying in %v...\n", i+1, sleep)
        time.Sleep(sleep)
    }
    return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

// Step 147
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

// Step 148
func compose(f, g func(int) int) func(int) int {
    return func(x int) int {
        return f(g(x))
    }
}

// Step 149
func partial3to2(f func(int, int, int) int, a int) func(int, int) int {
    return func(b, c int) int {
        return f(a, b, c)
    }
}

// Step 150
func emailValidator() func(string) bool {
    return func(s string) bool {
        // Simple validation
        return len(s) > 0 && contains(s, "@")
    }
}

func phoneValidator() func(string) bool {
    return func(s string) bool {
        return len(s) >= 10
    }
}

func usernameValidator() func(string) bool {
    return func(s string) bool {
        return len(s) >= 3 && len(s) <= 20
    }
}

func contains(s, substr string) bool {
    for i := 0; i <= len(s)-len(substr); i++ {
        if s[i:i+len(substr)] == substr {
            return true
        }
    }
    return false
}
```

---

## Step 151-180: Arrays และ Slices (Advanced)

```go
package main

import (
    "fmt"
    "sort"
)

func main() {
    // ===== Step 151: Slice Initialization Patterns =====
    var nil_slice []int                    // nil slice
    empty_slice := []int{}                 // empty slice
    make_slice := make([]int, 5)          // length 5
    make_cap := make([]int, 5, 10)        // length 5, capacity 10
    literal := []int{1, 2, 3, 4, 5}       // literal

    fmt.Printf("nil: %v, len=%d, cap=%d, is nil=%t\n",
        nil_slice, len(nil_slice), cap(nil_slice), nil_slice == nil)
    fmt.Printf("empty: %v, len=%d, cap=%d, is nil=%t\n",
        empty_slice, len(empty_slice), cap(empty_slice), empty_slice == nil)

    // ===== Step 152: Slice Internals =====
    arr := [5]int{1, 2, 3, 4, 5}
    slice1 := arr[1:4]  // [2, 3, 4]
    slice2 := arr[1:4]  // same backing array

    slice1[0] = 999     // affects both slices!

    fmt.Println("Array:", arr)
    fmt.Println("Slice1:", slice1)
    fmt.Println("Slice2:", slice2)

    // ===== Step 153: Full Slice Expression =====
    s := []int{1, 2, 3, 4, 5, 6}
    //     s[low:high:max]
    s1 := s[1:4:4]  // [2 3 4], capacity = 4-1 = 3

    fmt.Printf("s1: %v, len=%d, cap=%d\n", s1, len(s1), cap(s1))

    // ===== Step 154: Append and Capacity Growth =====
    var s154 []int
    printSliceInfo("Initial", s154)

    for i := 0; i < 10; i++ {
        s154 = append(s154, i)
        printSliceInfo(fmt.Sprintf("After append %d", i), s154)
    }

    // ===== Step 155: Pre-allocate for Performance =====
    // Inefficient
    var slow []int
    for i := 0; i < 1000; i++ {
        slow = append(slow, i)
    }

    // Efficient
    fast := make([]int, 0, 1000)
    for i := 0; i < 1000; i++ {
        fast = append(fast, i)
    }

    // ===== Step 156: Copy Slices =====
    source := []int{1, 2, 3, 4, 5}

    // Method 1: copy()
    dest1 := make([]int, len(source))
    n := copy(dest1, source)
    fmt.Printf("Copied %d elements: %v\n", n, dest1)

    // Method 2: append
    dest2 := append([]int{}, source...)
    fmt.Println("Append copy:", dest2)

    // Partial copy
    dest3 := make([]int, 3)
    copy(dest3, source)
    fmt.Println("Partial copy:", dest3)

    // ===== Step 157: Insert into Slice =====
    s157 := []int{1, 2, 4, 5}
    index := 2
    value := 3

    s157 = append(s157[:index], append([]int{value}, s157[index:]...)...)
    fmt.Println("After insert:", s157)

    // ===== Step 158: Delete from Slice =====
    s158 := []int{1, 2, 3, 4, 5}

    // Delete at index 2
    i := 2
    s158 = append(s158[:i], s158[i+1:]...)
    fmt.Println("After delete:", s158)

    // Delete range [1:3]
    s158 = []int{1, 2, 3, 4, 5}
    s158 = append(s158[:1], s158[3:]...)
    fmt.Println("After range delete:", s158)

    // ===== Step 159: Remove Duplicates =====
    withDupes := []int{1, 2, 2, 3, 3, 3, 4, 5, 5}
    unique := removeDuplicates(withDupes)
    fmt.Println("Unique:", unique)

    // ===== Step 160: Reverse Slice =====
    nums := []int{1, 2, 3, 4, 5}
    reverse(nums)
    fmt.Println("Reversed:", nums)

    // ===== Step 161: Rotate Slice =====
    nums161 := []int{1, 2, 3, 4, 5}
    rotate(nums161, 2)
    fmt.Println("Rotated by 2:", nums161)

    // ===== Step 162: Find Min/Max =====
    values := []int{45, 23, 67, 12, 89, 34}
    min, max := findMinMax(values)
    fmt.Printf("Min: %d, Max: %d\n", min, max)

    // ===== Step 163: Binary Search =====
    sorted := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
    target := 11

    index163 := binarySearch(sorted, target)
    if index163 != -1 {
        fmt.Printf("Found %d at index %d\n", target, index163)
    }

    // ===== Step 164: Merge Sorted Slices =====
    a := []int{1, 3, 5, 7}
    b := []int{2, 4, 6, 8}
    merged := mergeSorted(a, b)
    fmt.Println("Merged:", merged)

    // ===== Step 165: Partition Slice =====
    nums165 := []int{1, 5, 3, 8, 2, 9, 4, 7, 6}
    even, odd := partition(nums165, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println("Even:", even)
    fmt.Println("Odd:", odd)

    // ===== Step 166: Chunk Slice =====
    data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    chunks := chunk(data, 3)
    fmt.Println("Chunks of 3:", chunks)

    // ===== Step 167: Flatten Nested Slices =====
    nested := [][]int{{1, 2}, {3, 4, 5}, {6, 7, 8, 9}}
    flat := flatten(nested)
    fmt.Println("Flattened:", flat)

    // ===== Step 168: Intersection of Slices =====
    s1 := []int{1, 2, 3, 4, 5}
    s2 := []int{3, 4, 5, 6, 7}
    intersection := intersect(s1, s2)
    fmt.Println("Intersection:", intersection)

    // ===== Step 169: Union of Slices =====
    union := union(s1, s2)
    fmt.Println("Union:", union)

    // ===== Step 170: Difference of Slices =====
    diff := difference(s1, s2)
    fmt.Println("Difference s1-s2:", diff)

    // ===== Step 171-175: Sorting Variations =====

    // Custom sort
    people := []Person{
        {"Alice", 30},
        {"Bob", 25},
        {"Carol", 35},
    }

    sort.Slice(people, func(i, j int) bool {
        return people[i].Age < people[j].Age
    })
    fmt.Println("Sorted by age:", people)

    sort.Slice(people, func(i, j int) bool {
        return people[i].Name < people[j].Name
    })
    fmt.Println("Sorted by name:", people)

    // ===== Step 176-180: Advanced Slice Algorithms =====

    // Sliding window
    arr180 := []int{1, 2, 3, 4, 5, 6, 7, 8}
    windowSize := 3
    sums := slidingWindow(arr180, windowSize)
    fmt.Println("Sliding window sums:", sums)

    // Two pointers
    arr2ptr := []int{1, 2, 3, 4, 5, 6}
    pairs := findPairs(arr2ptr, 7)
    fmt.Println("Pairs that sum to 7:", pairs)
}

// Helper functions

func printSliceInfo(label string, s []int) {
    fmt.Printf("%s: len=%d, cap=%d, %v\n", label, len(s), cap(s), s)
}

func removeDuplicates(nums []int) []int {
    seen := make(map[int]bool)
    result := []int{}

    for _, num := range nums {
        if !seen[num] {
            seen[num] = true
            result = append(result, num)
        }
    }

    return result
}

func reverse(nums []int) {
    for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
        nums[i], nums[j] = nums[j], nums[i]
    }
}

func rotate(nums []int, k int) {
    k = k % len(nums)
    reverse(nums)
    reverse(nums[:k])
    reverse(nums[k:])
}

func findMinMax(nums []int) (int, int) {
    if len(nums) == 0 {
        return 0, 0
    }

    min, max := nums[0], nums[0]
    for _, num := range nums {
        if num < min {
            min = num
        }
        if num > max {
            max = num
        }
    }

    return min, max
}

func binarySearch(nums []int, target int) int {
    left, right := 0, len(nums)-1

    for left <= right {
        mid := left + (right-left)/2

        if nums[mid] == target {
            return mid
        } else if nums[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }

    return -1
}

func mergeSorted(a, b []int) []int {
    result := make([]int, 0, len(a)+len(b))
    i, j := 0, 0

    for i < len(a) && j < len(b) {
        if a[i] < b[j] {
            result = append(result, a[i])
            i++
        } else {
            result = append(result, b[j])
            j++
        }
    }

    result = append(result, a[i:]...)
    result = append(result, b[j:]...)

    return result
}

func partition(nums []int, predicate func(int) bool) ([]int, []int) {
    var pass, fail []int

    for _, num := range nums {
        if predicate(num) {
            pass = append(pass, num)
        } else {
            fail = append(fail, num)
        }
    }

    return pass, fail
}

func chunk(nums []int, size int) [][]int {
    var chunks [][]int

    for i := 0; i < len(nums); i += size {
        end := i + size
        if end > len(nums) {
            end = len(nums)
        }
        chunks = append(chunks, nums[i:end])
    }

    return chunks
}

func flatten(nested [][]int) []int {
    var result []int

    for _, slice := range nested {
        result = append(result, slice...)
    }

    return result
}

func intersect(a, b []int) []int {
    set := make(map[int]bool)
    for _, num := range a {
        set[num] = true
    }

    var result []int
    for _, num := range b {
        if set[num] {
            result = append(result, num)
            delete(set, num) // avoid duplicates
        }
    }

    return result
}

func union(a, b []int) []int {
    set := make(map[int]bool)

    for _, num := range a {
        set[num] = true
    }
    for _, num := range b {
        set[num] = true
    }

    var result []int
    for num := range set {
        result = append(result, num)
    }

    return result
}

func difference(a, b []int) []int {
    set := make(map[int]bool)
    for _, num := range b {
        set[num] = true
    }

    var result []int
    for _, num := range a {
        if !set[num] {
            result = append(result, num)
        }
    }

    return result
}

type Person struct {
    Name string
    Age  int
}

func slidingWindow(nums []int, k int) []int {
    if k > len(nums) {
        return nil
    }

    sums := make([]int, 0, len(nums)-k+1)

    for i := 0; i <= len(nums)-k; i++ {
        sum := 0
        for j := i; j < i+k; j++ {
            sum += nums[j]
        }
        sums = append(sums, sum)
    }

    return sums
}

func findPairs(nums []int, target int) [][2]int {
    var pairs [][2]int
    left, right := 0, len(nums)-1

    for left < right {
        sum := nums[left] + nums[right]
        if sum == target {
            pairs = append(pairs, [2]int{nums[left], nums[right]})
            left++
            right--
        } else if sum < target {
            left++
        } else {
            right--
        }
    }

    return pairs
}
```

ฉันจะสร้างส่วนที่เหลือต่อไป...

