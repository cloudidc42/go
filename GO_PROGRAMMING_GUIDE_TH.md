# คู่มือการเขียนโปรแกรม Go แบบสมบูรณ์ 🚀
## จากพื้นฐานสู่มืออาชีพ (Step 1-1000)

> **คู่มือฉบับสมบูรณ์** สำหรับการเรียนรู้ภาษา Go ตั้งแต่เริ่มต้นจนถึงระดับมืออาชีพ พร้อมตัวอย่างโค้ดที่ใช้งานได้จริง 100%

---

## สารบัญ (Table of Contents)

### 📚 ส่วนที่ 1: พื้นฐาน (Basic Level) - Steps 1-250
- [Step 1-10: การติดตั้งและเริ่มต้น](#step-1-10-การติดตั้งและเริ่มต้น)
- [Step 11-30: ไวยากรณ์พื้นฐาน](#step-11-30-ไวยากรณ์พื้นฐาน)
- [Step 31-60: ตัวแปรและชนิดข้อมูล](#step-31-60-ตัวแปรและชนิดข้อมูล)
- [Step 61-90: Operators และการคำนวณ](#step-61-90-operators-และการคำนวณ)
- [Step 91-120: Control Flow](#step-91-120-control-flow)
- [Step 121-150: Functions](#step-121-150-functions)
- [Step 151-180: Arrays และ Slices](#step-151-180-arrays-และ-slices)
- [Step 181-210: Maps](#step-181-210-maps)
- [Step 211-250: Structs และ Methods](#step-211-250-structs-และ-methods)

### 🔧 ส่วนที่ 2: ระดับกลาง (Intermediate Level) - Steps 251-500
- [Step 251-280: Pointers](#step-251-280-pointers)
- [Step 281-310: Interfaces](#step-281-310-interfaces)
- [Step 311-340: Error Handling](#step-311-340-error-handling)
- [Step 341-370: Goroutines](#step-341-370-goroutines)
- [Step 371-400: Channels](#step-371-400-channels)
- [Step 401-430: File I/O](#step-401-430-file-io)
- [Step 431-460: JSON และ XML](#step-431-460-json-และ-xml)
- [Step 461-500: Testing](#step-461-500-testing)

### 🚀 ส่วนที่ 3: ระดับสูง (Advanced Level) - Steps 501-750
- [Step 501-530: Web Development](#step-501-530-web-development)
- [Step 531-560: Database Operations](#step-531-560-database-operations)
- [Step 561-590: REST API](#step-561-590-rest-api)
- [Step 591-620: Middleware](#step-591-620-middleware)
- [Step 621-650: Context](#step-621-650-context)
- [Step 651-680: Reflection](#step-651-680-reflection)
- [Step 681-710: Performance Optimization](#step-681-710-performance-optimization)
- [Step 711-750: Design Patterns](#step-711-750-design-patterns)

### 💎 ส่วนที่ 4: ระดับมืออาชีพ (Expert Level) - Steps 751-1000
- [Step 751-780: Microservices](#step-751-780-microservices)
- [Step 781-810: gRPC](#step-781-810-grpc)
- [Step 811-840: Docker และ Kubernetes](#step-811-840-docker-และ-kubernetes)
- [Step 841-870: Message Queue](#step-841-870-message-queue)
- [Step 871-900: Caching Strategies](#step-871-900-caching-strategies)
- [Step 901-930: Security Best Practices](#step-901-930-security-best-practices)
- [Step 931-960: Monitoring และ Logging](#step-931-960-monitoring-และ-logging)
- [Step 961-1000: Production Deployment](#step-961-1000-production-deployment)

---

# ส่วนที่ 1: พื้นฐาน (Basic Level)

## Step 1-10: การติดตั้งและเริ่มต้น

### Step 1: ติดตั้ง Go

```bash
# สำหรับ Ubuntu/Debian
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# เพิ่ม PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### Step 2: ตรวจสอบการติดตั้ง

```bash
# ตรวจสอบเวอร์ชัน
go version

# ตรวจสอบ environment
go env
```

### Step 3: โครงสร้างโปรเจค Go พื้นฐาน

```
my-go-project/
├── main.go           # ไฟล์หลัก
├── go.mod            # Module definition
├── go.sum            # Dependencies checksums
├── pkg/              # Public libraries
│   └── utils/
│       └── helper.go
├── internal/         # Private application code
│   └── models/
│       └── user.go
├── cmd/              # Application entry points
│   └── server/
│       └── main.go
├── tests/            # Test files
│   └── main_test.go
└── docs/             # Documentation
    └── README.md
```

### Step 4: สร้างโปรเจคแรก

```bash
# สร้างโฟลเดอร์โปรเจค
mkdir hello-go
cd hello-go

# เริ่มต้น Go module
go mod init github.com/yourusername/hello-go
```

### Step 5: Hello World - โปรแกรมแรก

```go
// main.go
package main

import "fmt"

// main เป็น entry point ของโปรแกรม
func main() {
    // แสดงข้อความบนหน้าจอ
    fmt.Println("สวัสดี Go Programming! 👋")
}
```

**คำอธิบาย:**
- `package main` - กำหนดว่าไฟล์นี้เป็น package หลัก
- `import "fmt"` - นำเข้า package สำหรับการแสดงผล
- `func main()` - ฟังก์ชันหลักที่รันเมื่อเริ่มโปรแกรม
- `fmt.Println()` - แสดงข้อความและขึ้นบรรทัดใหม่

### Step 6: รันโปรแกรม

```bash
# วิธีที่ 1: รันโดยตรง
go run main.go

# วิธีที่ 2: Build และรัน
go build -o hello
./hello

# วิธีที่ 3: Build สำหรับ OS อื่น
GOOS=windows GOARCH=amd64 go build -o hello.exe
```

### Step 7: โครงสร้างพื้นฐานของโปรแกรม

```go
// main.go
package main

// Import หลายๆ packages
import (
    "fmt"
    "time"
    "math"
)

// ตัวแปรระดับ package (global)
var globalVar = "I'm global"

// Constants
const PI = 3.14159

// init function - รันก่อน main()
func init() {
    fmt.Println("Initializing...")
}

// main function
func main() {
    fmt.Println("Program started at:", time.Now())
    fmt.Println("Square root of 16:", math.Sqrt(16))
    fmt.Println(globalVar)
}
```

**ผลลัพธ์:**
```
Initializing...
Program started at: 2024-01-15 10:30:45.123456789 +0700 +07
Square root of 16: 4
I'm global
```

### Step 8: Comments และ Documentation

```go
package main

import "fmt"

/*
Package main เป็น entry point ของโปรแกรม
รองรับ multi-line comments แบบนี้
*/

// Function สำหรับคำนวณผลรวม
// รับค่า a และ b แล้วคืนค่าผลรวม
func add(a, b int) int {
    // Single line comment
    return a + b // Inline comment
}

// ฟังก์ชันหลัก
func main() {
    result := add(5, 3)
    fmt.Printf("5 + 3 = %d\n", result)
}
```

### Step 9: การใช้ go fmt และ Code Style

```bash
# Format โค้ดอัตโนมัติ (มาตรฐาน Go)
go fmt main.go

# ตรวจสอบปัญหาในโค้ด
go vet main.go

# ติดตั้ง golint
go install golang.org/x/lint/golint@latest

# รัน linter
golint main.go
```

**Go Coding Standards:**
- ใช้ camelCase สำหรับตัวแปร (myVariable)
- ใช้ PascalCase สำหรับ exported names (MyFunction)
- ชื่อสั้น กระชับ มีความหมาย
- Tab สำหรับ indentation (ไม่ใช่ spaces)

### Step 10: Module และ Dependencies

```bash
# เริ่มต้น module
go mod init myapp

# เพิ่ม dependency
go get github.com/gorilla/mux

# อัพเดท dependencies
go get -u ./...

# ล้าง unused dependencies
go mod tidy

# ดาวน์โหลด dependencies
go mod download
```

**go.mod ตัวอย่าง:**
```go
module github.com/yourusername/myapp

go 1.21

require (
    github.com/gorilla/mux v1.8.1
    github.com/joho/godotenv v1.5.1
)
```

---

## Step 11-30: ไวยากรณ์พื้นฐาน

### Step 11: Package Organization

```go
// package declaration - ต้องอยู่บรรทัดแรก
package main

// import statements
import (
    "fmt"
    "strings"
    "myapp/internal/models"  // internal package
    "myapp/pkg/utils"         // public package
)
```

**โครงสร้างแนะนำ:**
```
myapp/
├── main.go
├── go.mod
├── internal/          # โค้ดส่วนตัว (ไม่ export)
│   └── models/
│       └── user.go
├── pkg/               # โค้ดสาธารณะ (export ได้)
│   └── utils/
│       └── string.go
└── cmd/               # Multiple entry points
    ├── server/
    │   └── main.go
    └── cli/
        └── main.go
```

### Step 12: Import Statements

```go
package main

import (
    // Standard library
    "fmt"
    "os"

    // Third-party packages
    "github.com/gin-gonic/gin"

    // Local packages
    "myapp/internal/database"

    // Aliased import
    str "strings"

    // Blank import (เรียกใช้ init() เท่านั้น)
    _ "github.com/lib/pq"
)

func main() {
    // ใช้ alias
    result := str.ToUpper("hello")
    fmt.Println(result) // HELLO
}
```

### Step 13: Visibility และ Exported Names

```go
// utils.go
package utils

// PublicFunction - exported (ขึ้นต้นด้วยตัวใหญ่)
func PublicFunction() string {
    return "Anyone can call me!"
}

// privateFunction - not exported (ขึ้นต้นด้วยตัวเล็ก)
func privateFunction() string {
    return "Only package members can call me"
}

// PublicStruct - exported struct
type PublicStruct struct {
    PublicField  string // exported field
    privateField string // not exported
}
```

```go
// main.go
package main

import (
    "fmt"
    "myapp/pkg/utils"
)

func main() {
    // ✅ ใช้ได้
    result := utils.PublicFunction()
    s := utils.PublicStruct{
        PublicField: "visible",
    }

    // ❌ Error: ใช้ไม่ได้
    // utils.privateFunction()
    // s.privateField = "hidden"
}
```

### Step 14: Zero Values

```go
package main

import "fmt"

func main() {
    // ตัวแปรที่ไม่กำหนดค่าเริ่มต้นจะมี zero value
    var i int           // 0
    var f float64       // 0.0
    var b bool          // false
    var s string        // "" (empty string)
    var p *int          // nil
    var slice []int     // nil
    var m map[string]int // nil

    fmt.Printf("int: %d\n", i)
    fmt.Printf("float64: %f\n", f)
    fmt.Printf("bool: %t\n", b)
    fmt.Printf("string: %q\n", s)
    fmt.Printf("pointer: %v\n", p)
    fmt.Printf("slice: %v\n", slice)
    fmt.Printf("map: %v\n", m)
}
```

**ผลลัพธ์:**
```
int: 0
float64: 0.000000
bool: false
string: ""
pointer: <nil>
slice: []
map: map[]
```

### Step 15: Variable Declaration - 4 วิธี

```go
package main

import "fmt"

func main() {
    // วิธีที่ 1: var keyword พร้อมชนิดข้อมูล
    var name string = "John"

    // วิธีที่ 2: var keyword (type inference)
    var age = 25

    // วิธีที่ 3: Short declaration (:=)
    // ใช้ได้เฉพาะใน function
    city := "Bangkok"

    // วิธีที่ 4: Multiple declaration
    var (
        firstName string = "Jane"
        lastName  string = "Doe"
        height    float64 = 165.5
    )

    // Multiple short declaration
    x, y, z := 1, 2, 3

    fmt.Println(name, age, city)
    fmt.Println(firstName, lastName, height)
    fmt.Println(x, y, z)
}
```

### Step 16: Constants

```go
package main

import "fmt"

// Package-level constants
const (
    // iota = auto-increment counter
    Sunday    = iota // 0
    Monday           // 1
    Tuesday          // 2
    Wednesday        // 3
    Thursday         // 4
    Friday           // 5
    Saturday         // 6
)

const (
    // iota reset ใน const block ใหม่
    _ = iota // skip 0
    KB = 1 << (10 * iota) // 1 << 10 = 1024
    MB                     // 1 << 20 = 1048576
    GB                     // 1 << 30 = 1073741824
)

func main() {
    // Typed constant
    const pi float64 = 3.14159

    // Untyped constant
    const greeting = "Hello"

    // Multiple constants
    const (
        appName    = "MyApp"
        appVersion = "1.0.0"
    )

    fmt.Println("Today is day:", Wednesday)
    fmt.Printf("1 GB = %d bytes\n", GB)
    fmt.Println(appName, appVersion)
}
```

### Step 17: Type Conversions

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    // ===== Numeric Conversions =====
    var i int = 42
    var f float64 = float64(i)
    var u uint = uint(i)

    fmt.Printf("int: %d, float64: %f, uint: %d\n", i, f, u)

    // ===== String to Number =====
    str := "123"
    num, err := strconv.Atoi(str) // string to int
    if err != nil {
        fmt.Println("Error:", err)
    }
    fmt.Printf("String '%s' to int: %d\n", str, num)

    // ParseInt with base
    num64, _ := strconv.ParseInt("1010", 2, 64) // binary to int64
    fmt.Printf("Binary 1010 to decimal: %d\n", num64)

    // ParseFloat
    floatNum, _ := strconv.ParseFloat("3.14", 64)
    fmt.Printf("String to float: %f\n", floatNum)

    // ===== Number to String =====
    numStr := strconv.Itoa(42)
    fmt.Printf("Int 42 to string: '%s'\n", numStr)

    floatStr := strconv.FormatFloat(3.14159, 'f', 2, 64)
    fmt.Printf("Float to string: '%s'\n", floatStr)

    // ===== String and Bytes =====
    text := "Hello"
    bytes := []byte(text)
    backToString := string(bytes)

    fmt.Printf("String: %s, Bytes: %v, Back: %s\n",
        text, bytes, backToString)
}
```

### Step 18: Printf Formatting

```go
package main

import "fmt"

func main() {
    name := "Alice"
    age := 30
    height := 165.5
    active := true

    // ===== General =====
    fmt.Printf("%v\n", name)      // default format
    fmt.Printf("%+v\n", struct{Name string}{name}) // with field names
    fmt.Printf("%#v\n", name)     // Go-syntax representation
    fmt.Printf("%T\n", name)      // type

    // ===== Boolean =====
    fmt.Printf("%t\n", active)    // true or false

    // ===== Integer =====
    fmt.Printf("%d\n", age)       // decimal
    fmt.Printf("%b\n", age)       // binary
    fmt.Printf("%o\n", age)       // octal
    fmt.Printf("%x\n", age)       // hexadecimal (lowercase)
    fmt.Printf("%X\n", age)       // hexadecimal (uppercase)

    // ===== Float =====
    fmt.Printf("%f\n", height)    // decimal point
    fmt.Printf("%.2f\n", height)  // 2 decimal places
    fmt.Printf("%e\n", height)    // scientific notation

    // ===== String =====
    fmt.Printf("%s\n", name)      // plain string
    fmt.Printf("%q\n", name)      // quoted string
    fmt.Printf("%10s\n", name)    // width 10, right-aligned
    fmt.Printf("%-10s\n", name)   // width 10, left-aligned

    // ===== Pointer =====
    ptr := &age
    fmt.Printf("%p\n", ptr)       // pointer address

    // ===== Complete Example =====
    fmt.Printf("Name: %-10s Age: %3d Height: %6.2f Active: %t\n",
        name, age, height, active)
}
```

### Step 19: Input/Output พื้นฐาน

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    // ===== วิธีที่ 1: fmt.Scan (อ่านทีละคำ) =====
    var name string
    var age int

    fmt.Print("Enter name and age: ")
    fmt.Scan(&name, &age)
    fmt.Printf("Name: %s, Age: %d\n", name, age)

    // ===== วิธีที่ 2: fmt.Scanln (อ่านทั้งบรรทัด) =====
    var city string
    fmt.Print("Enter city: ")
    fmt.Scanln(&city)
    fmt.Println("City:", city)

    // ===== วิธีที่ 3: bufio.Scanner (แนะนำ) =====
    scanner := bufio.NewScanner(os.Stdin)

    fmt.Print("Enter full name: ")
    scanner.Scan()
    fullName := scanner.Text()

    fmt.Print("Enter bio: ")
    scanner.Scan()
    bio := scanner.Text()

    fmt.Printf("Full Name: %s\nBio: %s\n", fullName, bio)

    // ===== วิธีที่ 4: bufio.Reader =====
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter message: ")
    message, _ := reader.ReadString('\n')
    message = strings.TrimSpace(message) // ลบ whitespace

    fmt.Println("Message:", message)
}
```

### Step 20: Working with Strings

```go
package main

import (
    "fmt"
    "strings"
    "unicode/utf8"
)

func main() {
    text := "Hello, Go Programming!"

    // ===== Length =====
    fmt.Println("Length:", len(text))
    fmt.Println("Rune count:", utf8.RuneCountInString(text))

    // ===== Case =====
    fmt.Println("Upper:", strings.ToUpper(text))
    fmt.Println("Lower:", strings.ToLower(text))
    fmt.Println("Title:", strings.Title(text))

    // ===== Contains =====
    fmt.Println("Contains 'Go':", strings.Contains(text, "Go"))
    fmt.Println("HasPrefix 'Hello':", strings.HasPrefix(text, "Hello"))
    fmt.Println("HasSuffix '!':", strings.HasSuffix(text, "!"))

    // ===== Index =====
    fmt.Println("Index of 'Go':", strings.Index(text, "Go"))
    fmt.Println("LastIndex of 'o':", strings.LastIndex(text, "o"))

    // ===== Replace =====
    replaced := strings.Replace(text, "Go", "Golang", 1) // แทนที่ 1 ครั้ง
    fmt.Println("Replace:", replaced)

    allReplaced := strings.ReplaceAll(text, "o", "0")
    fmt.Println("ReplaceAll:", allReplaced)

    // ===== Split & Join =====
    parts := strings.Split(text, " ")
    fmt.Println("Split:", parts)

    joined := strings.Join(parts, "-")
    fmt.Println("Join:", joined)

    // ===== Trim =====
    spaced := "  Hello World  "
    fmt.Println("Trim:", strings.TrimSpace(spaced))

    // ===== Repeat =====
    fmt.Println("Repeat:", strings.Repeat("Go", 3))

    // ===== String Builder (efficient) =====
    var builder strings.Builder
    builder.WriteString("Hello")
    builder.WriteString(" ")
    builder.WriteString("World")
    fmt.Println("Builder:", builder.String())
}
```

### Step 21: Runes และ UTF-8

```go
package main

import (
    "fmt"
    "unicode/utf8"
)

func main() {
    // String ใน Go เป็น UTF-8 encoded
    thai := "สวัสดี" // Thai characters
    emoji := "👋🌍"

    // ===== Length vs Rune Count =====
    fmt.Printf("Thai: len=%d, runes=%d\n",
        len(thai), utf8.RuneCountInString(thai))
    fmt.Printf("Emoji: len=%d, runes=%d\n",
        len(emoji), utf8.RuneCountInString(emoji))

    // ===== Iterate by Bytes =====
    fmt.Print("Bytes: ")
    for i := 0; i < len(thai); i++ {
        fmt.Printf("%x ", thai[i])
    }
    fmt.Println()

    // ===== Iterate by Runes (ถูกต้อง) =====
    fmt.Println("Runes:")
    for i, r := range thai {
        fmt.Printf("  [%d] %c (U+%04X)\n", i, r, r)
    }

    // ===== Rune Manipulation =====
    runes := []rune(thai)
    fmt.Printf("Rune slice: %v\n", runes)
    fmt.Printf("First rune: %c\n", runes[0])

    // ===== Emoji =====
    for _, r := range emoji {
        fmt.Printf("%c ", r)
    }
    fmt.Println()
}
```

**ผลลัพธ์:**
```
Thai: len=18, runes=6
Emoji: len=8, runes=2
Bytes: e0 b8 aa e0 b8 a7 e0 b8 b1 e0 b8 aa e0 b8 94 e0 b8 b5
Runes:
  [0] ส (U+0E2A)
  [3] ว (U+0E27)
  [6] ั (U+0E31)
  [9] ส (U+0E2A)
  [12] ด (U+0E14)
  [15] ี (U+0E35)
👋 🌍
```

### Step 22-30: Advanced String Operations

```go
package main

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
)

func main() {
    // ===== Step 22: Regular Expressions =====
    text := "Email: john@example.com, Phone: 123-456-7890"

    // Match email
    emailRegex := regexp.MustCompile(`[\w.]+@[\w.]+`)
    email := emailRegex.FindString(text)
    fmt.Println("Email:", email)

    // Match phone
    phoneRegex := regexp.MustCompile(`\d{3}-\d{3}-\d{4}`)
    phone := phoneRegex.FindString(text)
    fmt.Println("Phone:", phone)

    // Find all matches
    allNumbers := regexp.MustCompile(`\d+`).FindAllString(text, -1)
    fmt.Println("All numbers:", allNumbers)

    // ===== Step 23: String Formatting =====
    name := "Alice"
    age := 30

    // Sprintf - return formatted string
    formatted := fmt.Sprintf("Name: %s, Age: %d", name, age)
    fmt.Println(formatted)

    // Complex formatting
    result := fmt.Sprintf("%-10s | %5d | %6.2f%%", "Product", 42, 75.5)
    fmt.Println(result)

    // ===== Step 24: String Comparison =====
    s1 := "Hello"
    s2 := "hello"

    fmt.Println("Equal:", s1 == s2)
    fmt.Println("EqualFold:", strings.EqualFold(s1, s2)) // case-insensitive
    fmt.Println("Compare:", strings.Compare(s1, s2))     // -1, 0, 1

    // ===== Step 25: String Slicing =====
    str := "Hello, World!"

    fmt.Println("First 5:", str[:5])        // "Hello"
    fmt.Println("From index 7:", str[7:])   // "World!"
    fmt.Println("Middle:", str[7:12])       // "World"

    // ===== Step 26: String Contains Multiple =====
    keywords := []string{"go", "python", "java"}
    text2 := "I love go programming"

    for _, keyword := range keywords {
        if strings.Contains(strings.ToLower(text2), keyword) {
            fmt.Printf("Found keyword: %s\n", keyword)
        }
    }

    // ===== Step 27: Padding Strings =====
    padLeft := fmt.Sprintf("%10s", "Go")      // "        Go"
    padRight := fmt.Sprintf("%-10s", "Go")    // "Go        "
    padZero := fmt.Sprintf("%05d", 42)        // "00042"

    fmt.Printf("|%s|\n", padLeft)
    fmt.Printf("|%s|\n", padRight)
    fmt.Printf("|%s|\n", padZero)

    // ===== Step 28: String to Number Array =====
    numbers := "1,2,3,4,5"
    parts := strings.Split(numbers, ",")
    var intArray []int

    for _, part := range parts {
        num, _ := strconv.Atoi(part)
        intArray = append(intArray, num)
    }
    fmt.Println("Int array:", intArray)

    // ===== Step 29: Multiline Strings =====
    multiline := `
        This is a
        multiline string
        with "quotes" and 'apostrophes'
    `
    fmt.Println(multiline)

    // ===== Step 30: String Reversal =====
    original := "Hello"
    runes := []rune(original)

    // Reverse
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }

    reversed := string(runes)
    fmt.Printf("Original: %s, Reversed: %s\n", original, reversed)
}
```

---

## Step 31-60: ตัวแปรและชนิดข้อมูล

### Step 31-35: Basic Types

```go
package main

import "fmt"

func main() {
    // ===== Step 31: Integer Types =====
    var i8 int8 = 127          // -128 to 127
    var i16 int16 = 32767      // -32768 to 32767
    var i32 int32 = 2147483647
    var i64 int64 = 9223372036854775807

    var u8 uint8 = 255         // 0 to 255
    var u16 uint16 = 65535
    var u32 uint32 = 4294967295
    var u64 uint64 = 18446744073709551615

    // int และ uint (ขึ้นกับ architecture: 32 or 64 bit)
    var i int = 100
    var u uint = 100

    // byte (alias for uint8)
    var b byte = 'A'

    // rune (alias for int32) - Unicode code point
    var r rune = '中'

    fmt.Printf("int8: %d, int16: %d, int32: %d, int64: %d\n", i8, i16, i32, i64)
    fmt.Printf("uint8: %d, byte: %c, rune: %c\n", u8, b, r)

    // ===== Step 32: Float Types =====
    var f32 float32 = 3.14
    var f64 float64 = 3.141592653589793

    fmt.Printf("float32: %.2f (±1.18e-38 to ±3.4e38)\n", f32)
    fmt.Printf("float64: %.15f (±2.23e-308 to ±1.80e308)\n", f64)

    // Scientific notation
    scientific := 1.23e9 // 1.23 × 10^9
    fmt.Printf("Scientific: %e = %f\n", scientific, scientific)

    // ===== Step 33: Complex Numbers =====
    var c64 complex64 = 1 + 2i
    var c128 complex128 = 3 + 4i

    fmt.Printf("complex64: %v\n", c64)
    fmt.Printf("complex128: %v\n", c128)

    // Operations
    sum := c64 + c128
    fmt.Printf("Sum: %v\n", sum)

    // Real and Imaginary parts
    fmt.Printf("Real: %v, Imag: %v\n", real(c128), imag(c128))

    // ===== Step 34: Boolean =====
    var isActive bool = true
    var isDeleted bool = false
    var isDefault bool // zero value = false

    fmt.Printf("isActive: %t, isDeleted: %t, isDefault: %t\n",
        isActive, isDeleted, isDefault)

    // Boolean operations
    result := isActive && !isDeleted
    fmt.Printf("Active and not deleted: %t\n", result)

    // ===== Step 35: String =====
    var name string = "Go Programming"
    var empty string // zero value = ""

    // Raw string (no escape)
    var path string = `C:\Users\Documents\file.txt`

    // Concatenation
    greeting := "Hello, " + name

    fmt.Printf("Name: %s, Empty: %q, Path: %s\n", name, empty, path)
    fmt.Println(greeting)
}
```

### Step 36-40: Type Inference และ Multiple Variables

```go
package main

import "fmt"

func main() {
    // ===== Step 36: Type Inference =====
    // Go จะกำหนด type อัตโนมัติจากค่าที่กำหนด
    x := 42              // int
    y := 3.14            // float64
    z := "hello"         // string
    isValid := true      // bool
    c := 1 + 2i          // complex128

    fmt.Printf("x: %T, y: %T, z: %T, isValid: %T, c: %T\n",
        x, y, z, isValid, c)

    // ===== Step 37: Multiple Variable Declaration =====
    // วิธีที่ 1: แยกบรรทัด
    var a, b, c int = 1, 2, 3

    // วิธีที่ 2: Mixed types
    var (
        name   string  = "Alice"
        age    int     = 30
        height float64 = 165.5
    )

    // วิธีที่ 3: Short declaration
    firstName, lastName, score := "John", "Doe", 95.5

    // วิธีที่ 4: Multiple assignment
    x1, y1 := 10, 20
    x1, y1 = y1, x1 // Swap!

    fmt.Println(a, b, c)
    fmt.Println(name, age, height)
    fmt.Println(firstName, lastName, score)
    fmt.Printf("After swap: x1=%d, y1=%d\n", x1, y1)

    // ===== Step 38: Variable Scope =====
    globalScope()

    // ===== Step 39: Short vs Long Declaration =====
    // ใช้ var เมื่อ:
    // 1. ประกาศในระดับ package
    // 2. ต้องการชนิดข้อมูลเฉพาะ
    // 3. zero value
    var count int           // zero value
    var ratio float32 = 0.5 // specific type

    // ใช้ := เมื่อ:
    // 1. ภายใน function
    // 2. ต้องการ type inference
    total := 0
    message := "Hello"

    fmt.Println(count, ratio, total, message)

    // ===== Step 40: Type Alias =====
    type UserID int
    type Temperature float64
    type Callback func(int) int

    var uid UserID = 12345
    var temp Temperature = 36.5
    var cb Callback = func(x int) int { return x * 2 }

    fmt.Printf("UserID: %d, Temp: %.1f°C, Callback(5): %d\n",
        uid, temp, cb(5))
}

// Package-level variables
var packageVar = "I'm accessible everywhere in this package"

func globalScope() {
    fmt.Println("From globalScope:", packageVar)

    // Local variable
    localVar := "I only exist in this function"

    if true {
        // Block scope
        blockVar := "I only exist in this block"
        fmt.Println(localVar, blockVar)
    }

    // fmt.Println(blockVar) // ❌ Error: undefined
}
```

### Step 41-50: Arrays

```go
package main

import "fmt"

func main() {
    // ===== Step 41: Array Declaration =====
    // วิธีที่ 1: กำหนดขนาด
    var arr1 [5]int
    fmt.Println("Zero value array:", arr1) // [0 0 0 0 0]

    // วิธีที่ 2: กำหนดค่าเริ่มต้น
    arr2 := [5]int{1, 2, 3, 4, 5}
    fmt.Println("Initialized array:", arr2)

    // วิธีที่ 3: กำหนดบางตำแหน่ง
    arr3 := [5]int{0: 10, 2: 20, 4: 30}
    fmt.Println("Partial init:", arr3) // [10 0 20 0 30]

    // วิธีที่ 4: Auto size
    arr4 := [...]int{1, 2, 3, 4, 5, 6}
    fmt.Println("Auto size:", arr4, "Length:", len(arr4))

    // ===== Step 42: Array Access =====
    numbers := [5]int{10, 20, 30, 40, 50}

    // Get element
    fmt.Println("First:", numbers[0])
    fmt.Println("Last:", numbers[len(numbers)-1])

    // Set element
    numbers[2] = 999
    fmt.Println("Modified:", numbers)

    // ===== Step 43: Array Iteration =====
    fruits := [4]string{"Apple", "Banana", "Orange", "Mango"}

    // วิธีที่ 1: for with index
    for i := 0; i < len(fruits); i++ {
        fmt.Printf("[%d] %s\n", i, fruits[i])
    }

    // วิธีที่ 2: for range
    for index, fruit := range fruits {
        fmt.Printf("%d: %s\n", index, fruit)
    }

    // วิธีที่ 3: for range (value only)
    for _, fruit := range fruits {
        fmt.Println("-", fruit)
    }

    // ===== Step 44: Multidimensional Arrays =====
    // 2D Array
    var matrix [3][3]int = [3][3]int{
        {1, 2, 3},
        {4, 5, 6},
        {7, 8, 9},
    }

    fmt.Println("Matrix:")
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            fmt.Printf("%3d ", matrix[i][j])
        }
        fmt.Println()
    }

    // ===== Step 45: Array Comparison =====
    a := [3]int{1, 2, 3}
    b := [3]int{1, 2, 3}
    c := [3]int{1, 2, 4}

    fmt.Println("a == b:", a == b) // true
    fmt.Println("a == c:", a == c) // false

    // ===== Step 46: Array Copy =====
    source := [3]int{1, 2, 3}
    destination := source // Copy by value!

    destination[0] = 999

    fmt.Println("Source:", source)           // [1 2 3]
    fmt.Println("Destination:", destination) // [999 2 3]

    // ===== Step 47: Array as Function Parameter =====
    nums := [5]int{1, 2, 3, 4, 5}

    // Pass by value (copy)
    modifyArray(nums)
    fmt.Println("After modifyArray:", nums) // unchanged

    // Pass by reference (pointer)
    modifyArrayByPointer(&nums)
    fmt.Println("After modifyArrayByPointer:", nums) // changed

    // ===== Step 48: Array Sum and Average =====
    scores := [5]float64{85.5, 92.0, 78.5, 95.0, 88.5}

    sum := 0.0
    for _, score := range scores {
        sum += score
    }
    average := sum / float64(len(scores))

    fmt.Printf("Sum: %.2f, Average: %.2f\n", sum, average)

    // ===== Step 49: Finding Min/Max =====
    values := [7]int{45, 23, 67, 12, 89, 34, 56}

    min, max := values[0], values[0]
    for _, v := range values {
        if v < min {
            min = v
        }
        if v > max {
            max = v
        }
    }

    fmt.Printf("Min: %d, Max: %d\n", min, max)

    // ===== Step 50: Array Search =====
    items := [6]string{"apple", "banana", "orange", "grape", "mango", "kiwi"}
    search := "orange"

    found := false
    foundIndex := -1

    for i, item := range items {
        if item == search {
            found = true
            foundIndex = i
            break
        }
    }

    if found {
        fmt.Printf("Found '%s' at index %d\n", search, foundIndex)
    } else {
        fmt.Printf("'%s' not found\n", search)
    }
}

func modifyArray(arr [5]int) {
    arr[0] = 999 // จะไม่มีผลกับ original
}

func modifyArrayByPointer(arr *[5]int) {
    arr[0] = 999 // จะแก้ไข original
}
```

### Step 51-60: Slices (สำคัญมาก!)

```go
package main

import (
    "fmt"
    "sort"
)

func main() {
    // ===== Step 51: Slice Declaration =====
    // Slice = Dynamic array

    // วิธีที่ 1: nil slice
    var slice1 []int
    fmt.Printf("nil slice: %v, len: %d, cap: %d, is nil: %t\n",
        slice1, len(slice1), cap(slice1), slice1 == nil)

    // วิธีที่ 2: Empty slice
    slice2 := []int{}
    fmt.Printf("empty slice: %v, len: %d, cap: %d, is nil: %t\n",
        slice2, len(slice2), cap(slice2), slice2 == nil)

    // วิธีที่ 3: With values
    slice3 := []int{1, 2, 3, 4, 5}
    fmt.Printf("initialized: %v, len: %d, cap: %d\n",
        slice3, len(slice3), cap(slice3))

    // วิธีที่ 4: make() function
    slice4 := make([]int, 5)      // length = 5, capacity = 5
    slice5 := make([]int, 5, 10)  // length = 5, capacity = 10

    fmt.Printf("make(5): %v, len: %d, cap: %d\n",
        slice4, len(slice4), cap(slice4))
    fmt.Printf("make(5,10): %v, len: %d, cap: %d\n",
        slice5, len(slice5), cap(slice5))

    // ===== Step 52: Slice from Array =====
    array := [6]int{10, 20, 30, 40, 50, 60}

    slice := array[1:4]  // [20 30 40]
    slice2_52 := array[:3]   // [10 20 30]
    slice3_52 := array[3:]   // [40 50 60]
    slice4_52 := array[:]    // [10 20 30 40 50 60]

    fmt.Println("array[1:4]:", slice)
    fmt.Println("array[:3]:", slice2_52)
    fmt.Println("array[3:]:", slice3_52)
    fmt.Println("array[:]:", slice4_52)

    // ===== Step 53: Append =====
    numbers := []int{1, 2, 3}
    fmt.Printf("Original: %v, len: %d, cap: %d\n",
        numbers, len(numbers), cap(numbers))

    // Append single
    numbers = append(numbers, 4)
    fmt.Printf("After append(4): %v, len: %d, cap: %d\n",
        numbers, len(numbers), cap(numbers))

    // Append multiple
    numbers = append(numbers, 5, 6, 7)
    fmt.Printf("After append(5,6,7): %v, len: %d, cap: %d\n",
        numbers, len(numbers), cap(numbers))

    // Append slice
    more := []int{8, 9, 10}
    numbers = append(numbers, more...)
    fmt.Printf("After append slice: %v, len: %d, cap: %d\n",
        numbers, len(numbers), cap(numbers))

    // ===== Step 54: Copy =====
    source := []int{1, 2, 3, 4, 5}

    // วิธีที่ 1: copy function
    dest1 := make([]int, len(source))
    copy(dest1, source)

    // วิธีที่ 2: append (less efficient)
    dest2 := append([]int{}, source...)

    dest1[0] = 999

    fmt.Println("Source:", source)  // [1 2 3 4 5]
    fmt.Println("Dest1:", dest1)    // [999 2 3 4 5]
    fmt.Println("Dest2:", dest2)    // [1 2 3 4 5]

    // ===== Step 55: Slice Tricks =====
    s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // Delete element at index 3
    index := 3
    s = append(s[:index], s[index+1:]...)
    fmt.Println("After delete index 3:", s)

    // Insert at index 2
    s = []int{1, 2, 3, 4, 5}
    index = 2
    value := 999
    s = append(s[:index], append([]int{value}, s[index:]...)...)
    fmt.Println("After insert 999 at index 2:", s)

    // Filter (keep even numbers)
    s = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    filtered := s[:0]
    for _, v := range s {
        if v%2 == 0 {
            filtered = append(filtered, v)
        }
    }
    fmt.Println("Even numbers:", filtered)

    // ===== Step 56: 2D Slices =====
    // Create 2D slice
    rows, cols := 3, 4
    matrix := make([][]int, rows)
    for i := range matrix {
        matrix[i] = make([]int, cols)
    }

    // Fill with values
    counter := 1
    for i := 0; i < rows; i++ {
        for j := 0; j < cols; j++ {
            matrix[i][j] = counter
            counter++
        }
    }

    // Print matrix
    fmt.Println("2D Slice:")
    for _, row := range matrix {
        fmt.Println(row)
    }

    // ===== Step 57: Sorting Slices =====
    nums := []int{5, 2, 8, 1, 9, 3}

    sort.Ints(nums)
    fmt.Println("Sorted ascending:", nums)

    sort.Sort(sort.Reverse(sort.IntSlice(nums)))
    fmt.Println("Sorted descending:", nums)

    // String slice
    words := []string{"banana", "apple", "cherry", "date"}
    sort.Strings(words)
    fmt.Println("Sorted strings:", words)

    // ===== Step 58: Searching in Sorted Slice =====
    sortedNums := []int{1, 3, 5, 7, 9, 11, 13, 15}
    target := 7

    index58 := sort.SearchInts(sortedNums, target)
    if index58 < len(sortedNums) && sortedNums[index58] == target {
        fmt.Printf("Found %d at index %d\n", target, index58)
    }

    // ===== Step 59: Slice Capacity Growth =====
    s59 := []int{}

    fmt.Println("Capacity growth demonstration:")
    for i := 0; i < 20; i++ {
        s59 = append(s59, i)
        fmt.Printf("len: %2d, cap: %2d\n", len(s59), cap(s59))
    }

    // ===== Step 60: Slice as Function Parameter =====
    original := []int{1, 2, 3, 4, 5}

    // Slices are reference types!
    modifySlice(original)
    fmt.Println("After modifySlice:", original) // modified!

    // But re-slice doesn't affect original length
    original = []int{1, 2, 3, 4, 5}
    appendToSlice(original)
    fmt.Println("After appendToSlice:", original) // not modified

    // Correct way
    original = appendToSliceCorrect(original)
    fmt.Println("After appendToSliceCorrect:", original) // modified
}

func modifySlice(s []int) {
    s[0] = 999 // จะแก้ไข original
}

func appendToSlice(s []int) {
    s = append(s, 100) // ไม่มีผลกับ original
}

func appendToSliceCorrect(s []int) []int {
    return append(s, 100) // return ค่ากลับไป
}
```

---

## Step 61-90: Operators และการคำนวณ

### Step 61-70: Arithmetic Operators

```go
package main

import (
    "fmt"
    "math"
)

func main() {
    // ===== Step 61: Basic Arithmetic =====
    a, b := 10, 3

    fmt.Printf("a = %d, b = %d\n", a, b)
    fmt.Printf("Addition: %d + %d = %d\n", a, b, a+b)
    fmt.Printf("Subtraction: %d - %d = %d\n", a, b, a-b)
    fmt.Printf("Multiplication: %d * %d = %d\n", a, b, a*b)
    fmt.Printf("Division: %d / %d = %d\n", a, b, a/b)       // Integer division
    fmt.Printf("Modulus: %d %% %d = %d\n", a, b, a%b)

    // Float division
    fmt.Printf("Float division: %d / %d = %.2f\n", a, b, float64(a)/float64(b))

    // ===== Step 62: Increment/Decrement =====
    count := 5

    count++  // count = count + 1
    fmt.Println("After count++:", count)

    count--  // count = count - 1
    fmt.Println("After count--:", count)

    // Note: ++count และ count++ ต่างกัน (Go ใช้ได้แค่ postfix)
    // ❌ ++count  // Error!

    // ===== Step 63: Compound Assignment =====
    num := 10

    num += 5   // num = num + 5
    fmt.Println("After += 5:", num)

    num -= 3   // num = num - 3
    fmt.Println("After -= 3:", num)

    num *= 2   // num = num * 2
    fmt.Println("After *= 2:", num)

    num /= 4   // num = num / 4
    fmt.Println("After /= 4:", num)

    num %= 3   // num = num % 3
    fmt.Println("After %= 3:", num)

    // ===== Step 64: Comparison Operators =====
    x, y := 5, 10

    fmt.Printf("%d == %d: %t\n", x, y, x == y)  // Equal
    fmt.Printf("%d != %d: %t\n", x, y, x != y)  // Not equal
    fmt.Printf("%d < %d: %t\n", x, y, x < y)    // Less than
    fmt.Printf("%d <= %d: %t\n", x, y, x <= y)  // Less than or equal
    fmt.Printf("%d > %d: %t\n", x, y, x > y)    // Greater than
    fmt.Printf("%d >= %d: %t\n", x, y, x >= y)  // Greater than or equal

    // ===== Step 65: Logical Operators =====
    isAdult := true
    hasLicense := false

    // AND (&&)
    canDrive := isAdult && hasLicense
    fmt.Printf("Can drive (AND): %t\n", canDrive)

    // OR (||)
    hasPermission := isAdult || hasLicense
    fmt.Printf("Has permission (OR): %t\n", hasPermission)

    // NOT (!)
    isChild := !isAdult
    fmt.Printf("Is child (NOT): %t\n", isChild)

    // Complex expression
    age := 25
    result := age >= 18 && age <= 65
    fmt.Printf("Working age (18-65): %t\n", result)

    // ===== Step 66: Bitwise Operators =====
    p, q := 12, 10  // 1100 and 1010 in binary

    fmt.Printf("p = %d (%08b), q = %d (%08b)\n", p, p, q, q)
    fmt.Printf("AND: %d & %d = %d (%08b)\n", p, q, p&q, p&q)
    fmt.Printf("OR: %d | %d = %d (%08b)\n", p, q, p|q, p|q)
    fmt.Printf("XOR: %d ^ %d = %d (%08b)\n", p, q, p^q, p^q)
    fmt.Printf("AND NOT: %d &^ %d = %d (%08b)\n", p, q, p&^q, p&^q)
    fmt.Printf("NOT: ^%d = %d (%08b)\n", p, ^p, ^p)

    // ===== Step 67: Bit Shift =====
    n := 8  // 1000 in binary

    fmt.Printf("Original: %d (%08b)\n", n, n)
    fmt.Printf("Left shift << 1: %d (%08b) [multiply by 2]\n", n<<1, n<<1)
    fmt.Printf("Left shift << 2: %d (%08b) [multiply by 4]\n", n<<2, n<<2)
    fmt.Printf("Right shift >> 1: %d (%08b) [divide by 2]\n", n>>1, n>>1)
    fmt.Printf("Right shift >> 2: %d (%08b) [divide by 4]\n", n>>2, n>>2)

    // ===== Step 68: Math Package =====
    fmt.Println("\n=== Math Functions ===")
    fmt.Printf("Abs(-5.5): %.2f\n", math.Abs(-5.5))
    fmt.Printf("Ceil(4.3): %.2f\n", math.Ceil(4.3))
    fmt.Printf("Floor(4.7): %.2f\n", math.Floor(4.7))
    fmt.Printf("Round(4.5): %.2f\n", math.Round(4.5))
    fmt.Printf("Max(10, 20): %.2f\n", math.Max(10, 20))
    fmt.Printf("Min(10, 20): %.2f\n", math.Min(10, 20))
    fmt.Printf("Pow(2, 10): %.0f\n", math.Pow(2, 10))
    fmt.Printf("Sqrt(16): %.2f\n", math.Sqrt(16))

    // ===== Step 69: Trigonometric Functions =====
    angle := 45.0
    radians := angle * math.Pi / 180

    fmt.Printf("\nAngle: %.0f degrees = %.4f radians\n", angle, radians)
    fmt.Printf("Sin(%.0f°): %.4f\n", angle, math.Sin(radians))
    fmt.Printf("Cos(%.0f°): %.4f\n", angle, math.Cos(radians))
    fmt.Printf("Tan(%.0f°): %.4f\n", angle, math.Tan(radians))

    // ===== Step 70: Advanced Math =====
    fmt.Println("\n=== Advanced Math ===")
    fmt.Printf("Exp(2): %.4f\n", math.Exp(2))
    fmt.Printf("Log(100): %.4f\n", math.Log(100))
    fmt.Printf("Log10(100): %.4f\n", math.Log10(100))
    fmt.Printf("Log2(8): %.4f\n", math.Log2(8))

    // Constants
    fmt.Printf("\nPi: %.10f\n", math.Pi)
    fmt.Printf("E: %.10f\n", math.E)
}
```

### Step 71-80: More Operators และ Practical Examples

```go
package main

import (
    "fmt"
    "math"
    "math/rand"
    "time"
)

func main() {
    // ===== Step 71: Random Numbers =====
    // Seed (ควรทำครั้งเดียวตอนเริ่มโปรแกรม)
    rand.Seed(time.Now().UnixNano())

    // Random int
    fmt.Println("Random int:", rand.Int())
    fmt.Println("Random int [0-99]:", rand.Intn(100))
    fmt.Println("Random int [10-20]:", 10+rand.Intn(11))

    // Random float
    fmt.Printf("Random float [0.0-1.0]: %.4f\n", rand.Float64())
    fmt.Printf("Random float [0.0-10.0]: %.4f\n", rand.Float64()*10)

    // ===== Step 72: Temperature Conversion =====
    celsius := 25.0
    fahrenheit := celsius*9/5 + 32
    kelvin := celsius + 273.15

    fmt.Printf("\n%.1f°C = %.1f°F = %.1fK\n", celsius, fahrenheit, kelvin)

    // Reverse
    f := 77.0
    c := (f - 32) * 5 / 9
    fmt.Printf("%.1f°F = %.1f°C\n", f, c)

    // ===== Step 73: Distance Calculation =====
    // Euclidean distance
    x1, y1 := 0.0, 0.0
    x2, y2 := 3.0, 4.0

    distance := math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
    fmt.Printf("\nDistance from (%.1f,%.1f) to (%.1f,%.1f): %.2f\n",
        x1, y1, x2, y2, distance)

    // ===== Step 74: Circle Calculations =====
    radius := 5.0
    area := math.Pi * radius * radius
    circumference := 2 * math.Pi * radius

    fmt.Printf("\nCircle with radius %.1f:\n", radius)
    fmt.Printf("  Area: %.2f\n", area)
    fmt.Printf("  Circumference: %.2f\n", circumference)

    // ===== Step 75: Percentage Calculations =====
    total := 1000.0
    discount := 15.0 // 15%

    discountAmount := total * discount / 100
    finalPrice := total - discountAmount

    fmt.Printf("\nOriginal: $%.2f\n", total)
    fmt.Printf("Discount: %.0f%% ($%.2f)\n", discount, discountAmount)
    fmt.Printf("Final: $%.2f\n", finalPrice)

    // ===== Step 76: Compound Interest =====
    principal := 10000.0
    rate := 5.0        // 5% per year
    years := 10.0

    // A = P(1 + r/n)^(nt)
    amount := principal * math.Pow(1+rate/100, years)
    interest := amount - principal

    fmt.Printf("\nCompound Interest:\n")
    fmt.Printf("  Principal: $%.2f\n", principal)
    fmt.Printf("  Rate: %.1f%% per year\n", rate)
    fmt.Printf("  Time: %.0f years\n", years)
    fmt.Printf("  Final Amount: $%.2f\n", amount)
    fmt.Printf("  Interest Earned: $%.2f\n", interest)

    // ===== Step 77: BMI Calculator =====
    weightKg := 70.0
    heightM := 1.75

    bmi := weightKg / (heightM * heightM)

    fmt.Printf("\nBMI Calculator:\n")
    fmt.Printf("  Weight: %.1f kg\n", weightKg)
    fmt.Printf("  Height: %.2f m\n", heightM)
    fmt.Printf("  BMI: %.2f\n", bmi)

    var category string
    switch {
    case bmi < 18.5:
        category = "Underweight"
    case bmi < 25:
        category = "Normal"
    case bmi < 30:
        category = "Overweight"
    default:
        category = "Obese"
    }
    fmt.Printf("  Category: %s\n", category)

    // ===== Step 78: Time Calculations =====
    seconds := 3665

    hours := seconds / 3600
    minutes := (seconds % 3600) / 60
    secs := seconds % 60

    fmt.Printf("\n%d seconds = %02d:%02d:%02d\n", seconds, hours, minutes, secs)

    // ===== Step 79: Average, Variance, Standard Deviation =====
    data := []float64{10, 12, 23, 23, 16, 23, 21, 16}

    // Mean
    sum := 0.0
    for _, v := range data {
        sum += v
    }
    mean := sum / float64(len(data))

    // Variance
    variance := 0.0
    for _, v := range data {
        variance += math.Pow(v-mean, 2)
    }
    variance /= float64(len(data))

    // Standard Deviation
    stdDev := math.Sqrt(variance)

    fmt.Printf("\nStatistics:\n")
    fmt.Printf("  Data: %v\n", data)
    fmt.Printf("  Mean: %.2f\n", mean)
    fmt.Printf("  Variance: %.2f\n", variance)
    fmt.Printf("  Std Dev: %.2f\n", stdDev)

    // ===== Step 80: Prime Number Check =====
    numbers := []int{2, 17, 20, 29, 100, 97}

    fmt.Println("\nPrime Number Check:")
    for _, num := range numbers {
        if isPrime(num) {
            fmt.Printf("  %d is PRIME\n", num)
        } else {
            fmt.Printf("  %d is NOT prime\n", num)
        }
    }
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
```

### Step 81-90: Advanced Calculations

```go
package main

import (
    "fmt"
    "math"
    "math/big"
)

func main() {
    // ===== Step 81: Factorial =====
    fmt.Println("=== Factorial ===")
    for i := 0; i <= 10; i++ {
        fmt.Printf("%d! = %d\n", i, factorial(i))
    }

    // ===== Step 82: Fibonacci =====
    fmt.Println("\n=== Fibonacci ===")
    for i := 0; i < 15; i++ {
        fmt.Printf("F(%d) = %d\n", i, fibonacci(i))
    }

    // ===== Step 83: GCD และ LCM =====
    a, b := 48, 18

    gcdResult := gcd(a, b)
    lcmResult := lcm(a, b)

    fmt.Printf("\nGCD(%d, %d) = %d\n", a, b, gcdResult)
    fmt.Printf("LCM(%d, %d) = %d\n", a, b, lcmResult)

    // ===== Step 84: Power (Fast Exponentiation) =====
    base := 2
    exp := 10
    result := power(base, exp)

    fmt.Printf("\n%d^%d = %d\n", base, exp, result)

    // ===== Step 85: Big Numbers =====
    // สำหรับตัวเลขที่ใหญ่มากๆ
    bigNum1 := big.NewInt(0)
    bigNum1.SetString("123456789012345678901234567890", 10)

    bigNum2 := big.NewInt(0)
    bigNum2.SetString("987654321098765432109876543210", 10)

    sum := new(big.Int)
    sum.Add(bigNum1, bigNum2)

    fmt.Println("\n=== Big Numbers ===")
    fmt.Println("Num1:", bigNum1)
    fmt.Println("Num2:", bigNum2)
    fmt.Println("Sum:", sum)

    // ===== Step 86: Matrix Operations =====
    fmt.Println("\n=== Matrix Operations ===")

    matrixA := [][]int{
        {1, 2, 3},
        {4, 5, 6},
    }

    matrixB := [][]int{
        {7, 8},
        {9, 10},
        {11, 12},
    }

    result86 := multiplyMatrix(matrixA, matrixB)

    fmt.Println("Matrix A:")
    printMatrix(matrixA)
    fmt.Println("Matrix B:")
    printMatrix(matrixB)
    fmt.Println("A × B:")
    printMatrix(result86)

    // ===== Step 87: Quadratic Equation Solver =====
    // ax² + bx + c = 0
    a87, b87, c87 := 1.0, -5.0, 6.0 // x² - 5x + 6 = 0

    solveQuadratic(a87, b87, c87)

    // ===== Step 88: Number Base Conversion =====
    decimal := 42

    fmt.Printf("\nDecimal %d:\n", decimal)
    fmt.Printf("  Binary: %b\n", decimal)
    fmt.Printf("  Octal: %o\n", decimal)
    fmt.Printf("  Hexadecimal: %x\n", decimal)

    // Custom base
    fmt.Printf("  Base 3: %s\n", toBase(decimal, 3))
    fmt.Printf("  Base 5: %s\n", toBase(decimal, 5))

    // ===== Step 89: Combinations และ Permutations =====
    n, r := 5, 3

    comb := combinations(n, r)
    perm := permutations(n, r)

    fmt.Printf("\nC(%d,%d) = %d\n", n, r, comb)
    fmt.Printf("P(%d,%d) = %d\n", n, r, perm)

    // ===== Step 90: Number Patterns =====
    fmt.Println("\n=== Triangle Pattern ===")
    printTriangle(5)

    fmt.Println("\n=== Number Pyramid ===")
    printNumberPyramid(5)
}

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

func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

func lcm(a, b int) int {
    return (a * b) / gcd(a, b)
}

func power(base, exp int) int {
    result := 1
    for exp > 0 {
        if exp%2 == 1 {
            result *= base
        }
        base *= base
        exp /= 2
    }
    return result
}

func multiplyMatrix(a, b [][]int) [][]int {
    rowsA, colsA := len(a), len(a[0])
    rowsB, colsB := len(b), len(b[0])

    if colsA != rowsB {
        return nil
    }

    result := make([][]int, rowsA)
    for i := range result {
        result[i] = make([]int, colsB)
    }

    for i := 0; i < rowsA; i++ {
        for j := 0; j < colsB; j++ {
            for k := 0; k < colsA; k++ {
                result[i][j] += a[i][k] * b[k][j]
            }
        }
    }

    return result
}

func printMatrix(m [][]int) {
    for _, row := range m {
        for _, val := range row {
            fmt.Printf("%4d ", val)
        }
        fmt.Println()
    }
}

func solveQuadratic(a, b, c float64) {
    discriminant := b*b - 4*a*c

    fmt.Printf("\nQuadratic: %.0fx² + %.0fx + %.0f = 0\n", a, b, c)

    if discriminant > 0 {
        x1 := (-b + math.Sqrt(discriminant)) / (2 * a)
        x2 := (-b - math.Sqrt(discriminant)) / (2 * a)
        fmt.Printf("Two real roots: x₁ = %.2f, x₂ = %.2f\n", x1, x2)
    } else if discriminant == 0 {
        x := -b / (2 * a)
        fmt.Printf("One real root: x = %.2f\n", x)
    } else {
        realPart := -b / (2 * a)
        imagPart := math.Sqrt(-discriminant) / (2 * a)
        fmt.Printf("Complex roots: x = %.2f ± %.2fi\n", realPart, imagPart)
    }
}

func toBase(n, base int) string {
    if n == 0 {
        return "0"
    }

    digits := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    result := ""

    for n > 0 {
        result = string(digits[n%base]) + result
        n /= base
    }

    return result
}

func combinations(n, r int) int {
    return factorial(n) / (factorial(r) * factorial(n-r))
}

func permutations(n, r int) int {
    return factorial(n) / factorial(n-r)
}

func printTriangle(n int) {
    for i := 1; i <= n; i++ {
        for j := 1; j <= i; j++ {
            fmt.Print("* ")
        }
        fmt.Println()
    }
}

func printNumberPyramid(n int) {
    for i := 1; i <= n; i++ {
        // Spaces
        for j := 0; j < n-i; j++ {
            fmt.Print(" ")
        }
        // Numbers
        for j := 1; j <= i; j++ {
            fmt.Printf("%d ", j)
        }
        fmt.Println()
    }
}
```

---

*คู่มือนี้ยาวมาก ฉันจะสร้างไฟล์แยกสำหรับแต่ละส่วน...*

