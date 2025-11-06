// ===== ตัวอย่าง: Channels - Communication between Goroutines =====
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("=== Channels Examples ===\n")

    // ===== Example 1: Basic Channel =====
    fmt.Println("Example 1: Basic Channel")

    messages := make(chan string)

    // Send value to channel (in goroutine)
    go func() {
        messages <- "Hello from goroutine!"
    }()

    // Receive value from channel
    msg := <-messages
    fmt.Printf("  Received: %s\n", msg)

    // ===== Example 2: Channel Synchronization =====
    fmt.Println("\nExample 2: Channel Synchronization")

    done := make(chan bool)

    go func() {
        fmt.Println("  Working...")
        time.Sleep(1 * time.Second)
        fmt.Println("  Done!")

        done <- true
    }()

    <-done // รอจนได้รับ signal

    // ===== Example 3: Buffered Channel =====
    fmt.Println("\nExample 3: Buffered Channel")

    // สร้าง buffered channel (capacity = 3)
    buffered := make(chan string, 3)

    // ส่งค่าได้โดยไม่ต้องมี receiver รอ (จนกว่า buffer เต็ม)
    buffered <- "first"
    buffered <- "second"
    buffered <- "third"

    // รับค่า
    fmt.Printf("  %s\n", <-buffered)
    fmt.Printf("  %s\n", <-buffered)
    fmt.Printf("  %s\n", <-buffered)

    // ===== Example 4: Channel Direction =====
    fmt.Println("\nExample 4: Channel Direction")

    pings := make(chan string, 1)
    pongs := make(chan string, 1)

    ping(pings, "message")
    pong(pings, pongs)

    fmt.Printf("  %s\n", <-pongs)

    // ===== Example 5: Select Statement =====
    fmt.Println("\nExample 5: Select Statement")

    c1 := make(chan string)
    c2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        c1 <- "one"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        c2 <- "two"
    }()

    // Select รอ channel ไหนก็ได้ที่พร้อม
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-c1:
            fmt.Printf("  Received: %s\n", msg1)
        case msg2 := <-c2:
            fmt.Printf("  Received: %s\n", msg2)
        }
    }

    // ===== Example 6: Timeout =====
    fmt.Println("\nExample 6: Timeout with Select")

    ch := make(chan string)

    go func() {
        time.Sleep(2 * time.Second)
        ch <- "result"
    }()

    select {
    case res := <-ch:
        fmt.Printf("  Got: %s\n", res)
    case <-time.After(1 * time.Second):
        fmt.Println("  Timeout!")
    }

    // ===== Example 7: Non-blocking Select =====
    fmt.Println("\nExample 7: Non-blocking Select")

    messages2 := make(chan string)
    signals := make(chan bool)

    // Non-blocking receive
    select {
    case msg := <-messages2:
        fmt.Printf("  Received message: %s\n", msg)
    default:
        fmt.Println("  No message received")
    }

    // Non-blocking send
    select {
    case messages2 <- "hi":
        fmt.Println("  Sent message")
    default:
        fmt.Println("  No message sent")
    }

    // Multi-way non-blocking select
    select {
    case msg := <-messages2:
        fmt.Printf("  Received: %s\n", msg)
    case sig := <-signals:
        fmt.Printf("  Received signal: %t\n", sig)
    default:
        fmt.Println("  No activity")
    }

    // ===== Example 8: Closing Channels =====
    fmt.Println("\nExample 8: Closing Channels")

    jobs := make(chan int, 5)

    // Send jobs
    go func() {
        for j := 1; j <= 3; j++ {
            jobs <- j
            fmt.Printf("  Sent job %d\n", j)
        }
        close(jobs) // ปิด channel
    }()

    // Receive jobs
    for j := range jobs {
        fmt.Printf("  Received job %d\n", j)
    }

    // ===== Example 9: Range over Channel =====
    fmt.Println("\nExample 9: Range over Channel")

    queue := make(chan string, 2)
    queue <- "one"
    queue <- "two"
    close(queue)

    for elem := range queue {
        fmt.Printf("  %s\n", elem)
    }

    // ===== Example 10: Pipeline Pattern =====
    fmt.Println("\nExample 10: Pipeline Pattern")

    // Stage 1: Generate numbers
    nums := gen(2, 3, 4, 5)

    // Stage 2: Square numbers
    squared := square(nums)

    // Stage 3: Print results
    for n := range squared {
        fmt.Printf("  %d\n", n)
    }

    // ===== Example 11: Fan-out, Fan-in Pattern =====
    fmt.Println("\nExample 11: Fan-out, Fan-in Pattern")

    input := make(chan int)

    // Fan-out: Start multiple workers
    c1_11 := worker11(input)
    c2_11 := worker11(input)
    c3_11 := worker11(input)

    // Fan-in: Merge results
    output := merge(c1_11, c2_11, c3_11)

    // Send work
    go func() {
        for i := 1; i <= 9; i++ {
            input <- i
        }
        close(input)
    }()

    // Collect results
    for result := range output {
        fmt.Printf("  Result: %d\n", result)
    }
}

// Channel direction: send-only
func ping(pings chan<- string, msg string) {
    pings <- msg
}

// Channel directions: receive-only and send-only
func pong(pings <-chan string, pongs chan<- string) {
    msg := <-pings
    pongs <- msg
}

// Pipeline: Generate numbers
func gen(nums ...int) <-chan int {
    out := make(chan int)

    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()

    return out
}

// Pipeline: Square numbers
func square(in <-chan int) <-chan int {
    out := make(chan int)

    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()

    return out
}

// Worker for fan-out/fan-in
func worker11(input <-chan int) <-chan int {
    output := make(chan int)

    go func() {
        for n := range input {
            output <- n * n
        }
        close(output)
    }()

    return output
}

// Merge multiple channels
func merge(cs ...<-chan int) <-chan int {
    out := make(chan int)

    for _, c := range cs {
        go func(ch <-chan int) {
            for n := range ch {
                out <- n
            }
        }(c)
    }

    // Close output when all inputs are done
    go func() {
        // Wait for all goroutines (simplified)
        time.Sleep(1 * time.Second)
        close(out)
    }()

    return out
}

// วิธีรัน:
// go run 02_channels.go
//
// สิ่งที่เรียนรู้:
// 1. Unbuffered vs Buffered channels
// 2. Channel direction (send-only, receive-only)
// 3. Select statement สำหรับ multiple channels
// 4. Timeout และ non-blocking operations
// 5. Closing channels และ range
// 6. Pipeline pattern
// 7. Fan-out, Fan-in pattern
