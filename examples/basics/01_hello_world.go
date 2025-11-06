// ===== ตัวอย่าง: Hello World =====
// โปรแกรม Go แรกของคุณ
package main

import "fmt"

func main() {
    // แสดงข้อความ Hello World
    fmt.Println("สวัสดี Go Programming! 👋")
    fmt.Println("Hello World!")

    // ใช้ Printf สำหรับ formatted output
    name := "Gopher"
    fmt.Printf("Hello, %s!\n", name)

    // ตัวแปรหลายตัว
    language := "Go"
    version := 1.21
    fmt.Printf("Learning %s version %.2f\n", language, version)
}

// วิธีรัน:
// go run 01_hello_world.go
