# ดัชนี - คู่มือการเขียนโปรแกรม Go ฉบับสมบูรณ์

## 📚 เอกสารทั้งหมด

### เอกสารหลัก
- **[README.md](README.md)** - คำแนะนำการใช้งานและภาพรวม
- **[QUICK_START.md](QUICK_START.md)** - เริ่มต้นอย่างรวดเร็ว
- **[INDEX.md](INDEX.md)** - ไฟล์นี้

### คู่มือการเรียนรู้
1. **[GO_PROGRAMMING_GUIDE_TH.md](GO_PROGRAMMING_GUIDE_TH.md)**
   - Step 1-90: พื้นฐาน
   - ติดตั้ง Go, ตัวแปร, ชนิดข้อมูล, Operators, String Operations

2. **[GO_PROGRAMMING_GUIDE_PART2.md](GO_PROGRAMMING_GUIDE_PART2.md)**
   - Step 91-250: Control Flow, Functions, Data Structures
   - If-Else, Switch, Loops, Functions, Arrays, Slices, Maps

3. **[GO_PROGRAMMING_GUIDE_PART3.md](GO_PROGRAMMING_GUIDE_PART3.md)**
   - Step 251-500: Advanced Topics
   - Pointers, Interfaces, Error Handling, Concurrency

---

## 📂 ตัวอย่างโค้ด

### Basics (พื้นฐาน)
- `examples/basics/01_hello_world.go` - Hello World
- `examples/basics/02_variables_types.go` - ตัวแปรและชนิดข้อมูล

### Intermediate (ระดับกลาง)
- `examples/intermediate/01_goroutines.go` - Goroutines และ Concurrency
- `examples/intermediate/02_channels.go` - Channels และ Communication

### Advanced (ระดับสูง)
- `examples/advanced/01_web_server.go` - Web Server
- `examples/advanced/02_rest_api.go` - REST API แบบสมบูรณ์

---

## 🗺️ Learning Path

### ระดับ 1: Beginner (2-3 สัปดาห์)
**เป้าหมาย:** เข้าใจพื้นฐาน Go และเขียนโปรแกรมง่ายๆ ได้

📖 **อ่าน:**
- GO_PROGRAMMING_GUIDE_TH.md (Step 1-90)
- GO_PROGRAMMING_GUIDE_PART2.md (Step 91-150)

💻 **ฝึกหัด:**
- รันทุกตัวอย่างใน `examples/basics/`
- สร้าง Calculator
- สร้าง Todo List (CLI)

✅ **ตรวจสอบ:**
- [ ] เขียน และรัน Go program ได้
- [ ] เข้าใจ variables, types, constants
- [ ] ใช้ if, for, switch ได้
- [ ] เขียน functions ได้
- [ ] ใช้ slices และ maps ได้

---

### ระดับ 2: Intermediate (3-4 สัปดาห์)
**เป้าหมาย:** เข้าใจ Concurrency และสร้าง Web Application ได้

📖 **อ่าน:**
- GO_PROGRAMMING_GUIDE_PART2.md (Step 151-250)
- GO_PROGRAMMING_GUIDE_PART3.md (Step 251-350)

💻 **ฝึกหัด:**
- รันทุกตัวอย่างใน `examples/intermediate/`
- สร้าง Concurrent File Processor
- สร้าง Simple Web Server

✅ **ตรวจสอบ:**
- [ ] เข้าใจ pointers
- [ ] ใช้ interfaces ได้
- [ ] Handle errors ถูกต้อง
- [ ] ใช้ goroutines และ channels ได้
- [ ] Read/Write files ได้

---

### ระดับ 3: Advanced (4-6 สัปดาห์)
**เป้าหมาย:** สร้าง Production-ready Applications

📖 **อ่าน:**
- GO_PROGRAMMING_GUIDE_PART3.md (Step 351-500)

💻 **ฝึกหัด:**
- รันทุกตัวอย่างใน `examples/advanced/`
- สร้าง REST API พร้อม Database
- เขียน Tests แบบครอบคลุม
- Deploy ขึ้น server

✅ **ตรวจสอบ:**
- [ ] สร้าง REST API ได้
- [ ] Connect database ได้
- [ ] เขียน tests ได้
- [ ] ใช้ middleware ได้
- [ ] Deploy application ได้

---

### ระดับ 4: Expert (2-3 เดือน)
**เป้าหมาย:** Microservices และ Cloud-native Development

📖 **อ่าน:**
- คู่มือเพิ่มเติมเกี่ยวกับ Microservices, gRPC, Docker, Kubernetes

💻 **ฝึกหัด:**
- สร้าง Microservices System
- Deploy ด้วย Docker + Kubernetes
- Setup Monitoring และ Logging
- Implement CI/CD

✅ **ตรวจสอบ:**
- [ ] Design microservices ได้
- [ ] ใช้ gRPC ได้
- [ ] Deploy ด้วย Docker/K8s ได้
- [ ] Setup monitoring ได้
- [ ] Handle production issues ได้

---

## 📋 Step-by-Step Index

### Part 1: พื้นฐาน (Step 1-90)

#### Installation และ Setup (Step 1-10)
- Step 1: ติดตั้ง Go
- Step 2: ตรวจสอบการติดตั้ง
- Step 3: โครงสร้างโปรเจค
- Step 4: สร้างโปรเจคแรก
- Step 5: Hello World
- Step 6: รันโปรแกรม
- Step 7: โครงสร้างพื้นฐาน
- Step 8: Comments และ Documentation
- Step 9: go fmt และ Code Style
- Step 10: Module และ Dependencies

#### Syntax Basics (Step 11-30)
- Step 11-13: Package Organization
- Step 14: Zero Values
- Step 15: Variable Declaration (4 วิธี)
- Step 16: Constants และ iota
- Step 17: Type Conversions
- Step 18: Printf Formatting
- Step 19: Input/Output
- Step 20-21: Working with Strings และ Runes
- Step 22-30: Advanced String Operations

#### Variables และ Types (Step 31-60)
- Step 31-35: Basic Types (int, float, bool, string)
- Step 36-40: Type Inference, Multiple Variables
- Step 41-50: Arrays
- Step 51-60: Slices

#### Operators (Step 61-90)
- Step 61-70: Arithmetic Operators, Math Functions
- Step 71-80: Random Numbers, Calculations
- Step 81-90: Advanced Algorithms

---

### Part 2: Control Flow และ Functions (Step 91-250)

#### Control Flow (Step 91-120)
- Step 91-100: If-Else Statements
- Step 101-110: Switch Statements
- Step 111-120: Loops (for, range)

#### Functions (Step 121-150)
- Step 121-130: Basic Functions
- Step 131-140: Advanced Functions (Error Handling, Defer, Panic, Recover)
- Step 141-150: Function Patterns

#### Data Structures (Step 151-210)
- Step 151-180: Arrays และ Slices (Advanced)
- Step 181-210: Maps

#### Structs (Step 211-250)
- Step 211-230: Struct Basics
- Step 231-250: Methods และ Patterns

---

### Part 3: Advanced Topics (Step 251-500)

#### Pointers (Step 251-280)
- Pointer basics, dereferencing, use cases

#### Interfaces (Step 281-310)
- Interface basics, type assertion, polymorphism

#### Error Handling (Step 311-340)
- Error types, custom errors, wrapping

#### Concurrency (Step 341-400)
- Step 341-370: Goroutines
- Step 371-400: Channels

#### File I/O (Step 401-430)
- Reading, writing, working with files

#### JSON และ XML (Step 431-460)
- Encoding, decoding, struct tags

#### Testing (Step 461-500)
- Unit tests, table-driven tests, benchmarks

---

## 🎯 หัวข้อเฉพาะทาง

### Web Development
- HTTP Server
- Routing
- Middleware
- Templates
- REST API
- WebSockets

### Database
- SQL (database/sql)
- PostgreSQL
- MySQL
- MongoDB
- ORM (GORM)

### Concurrency Patterns
- Worker Pools
- Pipeline
- Fan-out/Fan-in
- Rate Limiting
- Context

### Testing
- Unit Testing
- Integration Testing
- Table-Driven Tests
- Mocking
- Benchmarking

### DevOps
- Docker
- Kubernetes
- CI/CD
- Monitoring
- Logging

---

## 🔍 ค้นหาหัวข้อ

### ต้องการเรียนรู้...

**Variables และ Types?**
→ GO_PROGRAMMING_GUIDE_TH.md (Step 31-60)

**Control Flow (if, for, switch)?**
→ GO_PROGRAMMING_GUIDE_PART2.md (Step 91-120)

**Functions?**
→ GO_PROGRAMMING_GUIDE_PART2.md (Step 121-150)

**Slices และ Maps?**
→ GO_PROGRAMMING_GUIDE_PART2.md (Step 151-210)

**Structs?**
→ GO_PROGRAMMING_GUIDE_PART2.md (Step 211-250)

**Pointers?**
→ GO_PROGRAMMING_GUIDE_PART3.md (Step 251-280)

**Interfaces?**
→ GO_PROGRAMMING_GUIDE_PART3.md (Step 281-310)

**Goroutines และ Channels?**
→ GO_PROGRAMMING_GUIDE_PART3.md (Step 341-400)
→ examples/intermediate/01_goroutines.go
→ examples/intermediate/02_channels.go

**Web Server?**
→ examples/advanced/01_web_server.go

**REST API?**
→ examples/advanced/02_rest_api.go

---

## 📊 ความยาวของเนื้อหา

| Part | Steps | Pages | Time to Complete |
|------|-------|-------|------------------|
| Part 1 | 1-90 | ~50 | 1-2 weeks |
| Part 2 | 91-250 | ~80 | 2-3 weeks |
| Part 3 | 251-500 | ~120 | 4-6 weeks |
| **Total** | **1-500** | **~250** | **2-3 months** |

---

## 💡 เคล็ดลับการใช้คู่มือ

1. **อ่านตามลำดับ** - อย่าข้ามพื้นฐาน
2. **ทำโค้ดทุก step** - อย่าแค่อ่าน
3. **ทำแบบฝึกหัด** - ฝึกจนชำนาญ
4. **สร้างโปรเจค** - Apply สิ่งที่เรียน
5. **Review บ่อยๆ** - กลับมาอ่านซ้ำ

---

## 🚀 เริ่มต้นเลย!

```bash
# เริ่มจาก Quick Start
cat QUICK_START.md

# หรือเริ่มจากคู่มือหลัก
cat GO_PROGRAMMING_GUIDE_TH.md

# หรือลองรันตัวอย่างเลย
cd examples/basics
go run 01_hello_world.go
```

---

**Happy Learning! 🎉**
