# Quick Start Guide - Go Programming 🚀

## การเริ่มต้นอย่างรวดเร็ว

### 1. ติดตั้ง Go (5 นาที)

```bash
# Linux/macOS
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# ตรวจสอบ
go version
```

### 2. โปรแกรมแรก (2 นาที)

```bash
# สร้างไฟล์ hello.go
cat > hello.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go! 🎉")
}
EOF

# รัน
go run hello.go
```

### 3. เรียนรู้ตามลำดับ (แนะนำ)

#### สัปดาห์ที่ 1-2: พื้นฐาน
```bash
# อ่านคู่มือ
cat GO_PROGRAMMING_GUIDE_TH.md

# รันตัวอย่าง
cd examples/basics
go run 01_hello_world.go
go run 02_variables_types.go
```

**เรียนรู้:**
- ตัวแปร และ ชนิดข้อมูล
- Control Flow (if, for, switch)
- Functions
- Arrays, Slices, Maps
- Structs

#### สัปดาห์ที่ 3-4: ระดับกลาง
```bash
cd examples/intermediate
go run 01_goroutines.go
go run 02_channels.go
```

**เรียนรู้:**
- Pointers
- Interfaces
- Error Handling
- Goroutines & Channels (Concurrency)
- File I/O

#### สัปดาห์ที่ 5-6: ระดับสูง
```bash
cd examples/advanced
go run 01_web_server.go
go run 02_rest_api.go
```

**เรียนรู้:**
- Web Server
- REST API
- Database Operations
- Testing
- Performance

### 4. โปรเจคแนะนำ

#### Beginner Projects
1. **Calculator** - เครื่องคิดเลขพื้นฐาน
2. **Todo List** - CLI todo app
3. **Number Guessing Game** - เกมทายตัวเลข

#### Intermediate Projects
1. **Weather CLI** - ดึงข้อมูลสภาพอากาศจาก API
2. **URL Shortener** - บริการย่อ URL
3. **Chat Server** - Real-time chat

#### Advanced Projects
1. **Blog API** - REST API สำหรับ blog
2. **Task Manager** - จัดการ tasks พร้อม database
3. **File Storage Service** - อัพโหลดและดาวน์โหลดไฟล์

### 5. คำสั่งที่ใช้บ่อย

```bash
# รันโปรแกรม
go run main.go

# Build executable
go build -o myapp

# Format โค้ด
go fmt ./...

# ตรวจสอบปัญหา
go vet ./...

# รัน tests
go test ./...

# ติดตั้ง dependencies
go get github.com/package/name

# จัดการ dependencies
go mod init myproject
go mod tidy
```

### 6. โครงสร้างโปรเจค Go มาตรฐาน

```
myproject/
├── main.go           # Entry point
├── go.mod            # Dependencies
├── go.sum            # Checksums
├── README.md         # Documentation
│
├── cmd/              # Main applications
│   └── server/
│       └── main.go
│
├── internal/         # Private code
│   ├── models/
│   ├── handlers/
│   └── services/
│
├── pkg/              # Public libraries
│   └── utils/
│
├── api/              # API definitions
├── web/              # Web assets
├── configs/          # Configuration files
├── scripts/          # Scripts
├── test/             # Additional tests
└── docs/             # Documentation
```

### 7. Resources

- **คู่มือหลัก:** `GO_PROGRAMMING_GUIDE_TH.md`
- **README:** `README.md`
- **ตัวอย่าง:** `examples/`
- **Official:** https://go.dev
- **Tour:** https://go.dev/tour

### 8. Tips สำหรับผู้เริ่มต้น

✅ **DO:**
- เขียนโค้ดทุกวัน (ถึงแค่ 30 นาที)
- อ่าน error messages ให้เข้าใจ
- ใช้ `go fmt` เสมอ
- เขียน tests
- อ่าน standard library source code

❌ **DON'T:**
- อย่าข้าม fundamentals
- อย่าคัดลอกโค้ดโดยไม่เข้าใจ
- อย่าเขียนโค้ดที่ซับซ้อนเกินไป (Keep it simple)
- อย่าลืม handle errors

### 9. เริ่มต้นเลย!

```bash
# Clone/Download คู่มือนี้
cd /path/to/go

# เริ่มจากตัวอย่างแรก
cd examples/basics
go run 01_hello_world.go

# อ่านคู่มือ
cat ../GO_PROGRAMMING_GUIDE_TH.md

# ทำโปรเจคเล็กๆ
mkdir my-first-project
cd my-first-project
go mod init my-first-project

# สร้าง main.go
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    name := "Gopher"
    fmt.Printf("Hello, %s! Let's learn Go! 🚀\n", name)
}
EOF

# รัน
go run main.go
```

### 10. Next Steps

1. **Week 1:** Complete basics (Step 1-100)
2. **Week 2:** Practice with mini projects
3. **Week 3:** Learn concurrency (Goroutines & Channels)
4. **Week 4:** Build a web server
5. **Week 5+:** Advanced topics และ production deployment

---

## 📌 สรุป Learning Path

```
Day 1-7:    Basics (variables, types, control flow)
Day 8-14:   Functions, Slices, Maps, Structs
Day 15-21:  Pointers, Interfaces, Error Handling
Day 22-28:  Goroutines & Channels
Day 29-35:  File I/O, JSON, Testing
Day 36-42:  Web Server & REST API
Day 43-49:  Database & Advanced Patterns
Day 50+:    Production & Real Projects
```

---

## 🎯 Goal Setting

**2 Weeks:** ทำโปรเจคพื้นฐานได้
**1 Month:** สร้าง REST API ได้
**2 Months:** Deploy production app ได้
**3 Months:** Microservices และ Cloud

---

## 🚀 Happy Learning!

```go
package main

import "fmt"

func main() {
    fmt.Println("Your Go journey starts NOW! 🚀")
}
```
