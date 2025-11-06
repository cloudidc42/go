// ===== ตัวอย่าง: Goroutines - Concurrent Programming =====
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    fmt.Println("=== Goroutines Examples ===\n")

    // ===== Example 1: Simple Goroutine =====
    fmt.Println("Example 1: Simple Goroutine")

    // Sequential execution
    sayHello("Sequential")

    // Concurrent execution
    go sayHello("Goroutine")

    // รอให้ goroutine ทำงานเสร็จ
    time.Sleep(1 * time.Second)

    // ===== Example 2: Multiple Goroutines =====
    fmt.Println("\nExample 2: Multiple Goroutines")

    for i := 1; i <= 5; i++ {
        go func(n int) {
            fmt.Printf("  Goroutine %d is running\n", n)
        }(i)
    }

    time.Sleep(1 * time.Second)

    // ===== Example 3: WaitGroup =====
    fmt.Println("\nExample 3: Using WaitGroup")

    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1) // เพิ่ม counter

        go func(n int) {
            defer wg.Done() // ลด counter เมื่อเสร็จ

            fmt.Printf("  Worker %d: Starting\n", n)
            time.Sleep(time.Duration(n*100) * time.Millisecond)
            fmt.Printf("  Worker %d: Done\n", n)
        }(i)
    }

    wg.Wait() // รอให้ทุก goroutine เสร็จ
    fmt.Println("  All workers completed!")

    // ===== Example 4: Concurrent Calculations =====
    fmt.Println("\nExample 4: Concurrent Calculations")

    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    results := make([]int, len(numbers))

    var wg2 sync.WaitGroup

    for i, num := range numbers {
        wg2.Add(1)

        go func(index, value int) {
            defer wg2.Done()

            // คำนวณกำลังสอง
            results[index] = value * value
        }(i, num)
    }

    wg2.Wait()

    fmt.Printf("  Numbers: %v\n", numbers)
    fmt.Printf("  Squares: %v\n", results)

    // ===== Example 5: Race Condition (Problem) =====
    fmt.Println("\nExample 5: Race Condition (Unsafe)")

    counter := 0
    var wg3 sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg3.Add(1)

        go func() {
            defer wg3.Done()
            counter++ // ไม่ thread-safe!
        }()
    }

    wg3.Wait()
    fmt.Printf("  Counter (unsafe): %d (expected 1000)\n", counter)

    // ===== Example 6: Mutex (Solution) =====
    fmt.Println("\nExample 6: Using Mutex (Safe)")

    safeCounter := 0
    var mu sync.Mutex
    var wg4 sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg4.Add(1)

        go func() {
            defer wg4.Done()

            mu.Lock()
            safeCounter++
            mu.Unlock()
        }()
    }

    wg4.Wait()
    fmt.Printf("  Counter (safe): %d\n", safeCounter)

    // ===== Example 7: Real-world Use Case - Web Scraper =====
    fmt.Println("\nExample 7: Concurrent Web Scraper Simulation")

    urls := []string{
        "https://example.com/page1",
        "https://example.com/page2",
        "https://example.com/page3",
        "https://example.com/page4",
        "https://example.com/page5",
    }

    var wg5 sync.WaitGroup

    for _, url := range urls {
        wg5.Add(1)

        go func(u string) {
            defer wg5.Done()
            fetchURL(u)
        }(url)
    }

    wg5.Wait()
    fmt.Println("  All URLs fetched!")

    // ===== Example 8: Worker Pool Pattern =====
    fmt.Println("\nExample 8: Worker Pool Pattern")

    jobs := make(chan int, 10)
    var wg6 sync.WaitGroup

    // Start 3 workers
    for w := 1; w <= 3; w++ {
        wg6.Add(1)

        go worker(w, jobs, &wg6)
    }

    // Send jobs
    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs)

    wg6.Wait()
    fmt.Println("  All jobs processed!")
}

func sayHello(from string) {
    for i := 0; i < 3; i++ {
        fmt.Printf("  %s: Hello %d\n", from, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func fetchURL(url string) {
    fmt.Printf("  Fetching: %s\n", url)
    time.Sleep(time.Duration(100+len(url)) * time.Millisecond)
    fmt.Printf("  Completed: %s\n", url)
}

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
    defer wg.Done()

    for j := range jobs {
        fmt.Printf("  Worker %d: Processing job %d\n", id, j)
        time.Sleep(200 * time.Millisecond)
        fmt.Printf("  Worker %d: Finished job %d\n", id, j)
    }
}

// วิธีรัน:
// go run 01_goroutines.go
//
// สิ่งที่เรียนรู้:
// 1. go keyword สำหรับสร้าง goroutine
// 2. sync.WaitGroup สำหรับรอ goroutines
// 3. sync.Mutex สำหรับป้องกัน race condition
// 4. Worker Pool pattern สำหรับจัดการงานหลายๆ อัน
