// ===== ตัวอย่าง: Variables และ Types =====
package main

import "fmt"

func main() {
    // ===== การประกาศตัวแปร =====

    // วิธีที่ 1: var keyword
    var name string = "John Doe"
    var age int = 30
    var height float64 = 175.5

    // วิธีที่ 2: type inference
    var city = "Bangkok"
    var temperature = 35.5

    // วิธีที่ 3: short declaration
    country := "Thailand"
    population := 70000000

    // วิธีที่ 4: multiple variables
    var (
        firstName string = "Jane"
        lastName  string = "Smith"
        score     int    = 95
    )

    // แสดงผล
    fmt.Println("=== Personal Information ===")
    fmt.Printf("Name: %s\n", name)
    fmt.Printf("Age: %d\n", age)
    fmt.Printf("Height: %.1f cm\n", height)
    fmt.Printf("City: %s\n", city)
    fmt.Printf("Temperature: %.1f°C\n", temperature)

    fmt.Println("\n=== Location ===")
    fmt.Printf("Country: %s\n", country)
    fmt.Printf("Population: %d\n", population)

    fmt.Println("\n=== Student Info ===")
    fmt.Printf("Name: %s %s\n", firstName, lastName)
    fmt.Printf("Score: %d\n", score)

    // ===== ชนิดข้อมูลพื้นฐาน =====

    // Integer types
    var i8 int8 = 127
    var i16 int16 = 32767
    var i32 int32 = 2147483647
    var i64 int64 = 9223372036854775807

    // Unsigned integer types
    var u8 uint8 = 255
    var u16 uint16 = 65535

    // Float types
    var f32 float32 = 3.14
    var f64 float64 = 3.14159265359

    // Boolean
    var isActive bool = true
    var isDeleted bool = false

    // String
    var message string = "Hello, Go!"

    fmt.Println("\n=== Data Types ===")
    fmt.Printf("int8: %d\n", i8)
    fmt.Printf("int16: %d\n", i16)
    fmt.Printf("int32: %d\n", i32)
    fmt.Printf("int64: %d\n", i64)
    fmt.Printf("uint8: %d\n", u8)
    fmt.Printf("uint16: %d\n", u16)
    fmt.Printf("float32: %.2f\n", f32)
    fmt.Printf("float64: %.11f\n", f64)
    fmt.Printf("bool: isActive=%t, isDeleted=%t\n", isActive, isDeleted)
    fmt.Printf("string: %s\n", message)

    // ===== Constants =====
    const Pi = 3.14159
    const AppName = "My Go App"
    const MaxUsers = 100

    fmt.Println("\n=== Constants ===")
    fmt.Printf("Pi: %.5f\n", Pi)
    fmt.Printf("App Name: %s\n", AppName)
    fmt.Printf("Max Users: %d\n", MaxUsers)

    // ===== Type Conversion =====
    var x int = 42
    var y float64 = float64(x)
    var z uint = uint(x)

    fmt.Println("\n=== Type Conversion ===")
    fmt.Printf("int: %d\n", x)
    fmt.Printf("float64: %.2f\n", y)
    fmt.Printf("uint: %d\n", z)
}

// วิธีรัน:
// go run 02_variables_types.go
