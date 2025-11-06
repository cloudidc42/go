// ===== ตัวอย่าง: Advanced Concurrency Patterns =====
package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ===== Part 1: Worker Pool Pattern =====

type Job struct {
	ID      int
	Data    string
	Process func(string) string
}

type Result struct {
	JobID  int
	Output string
	Error  error
}

type WorkerPool struct {
	numWorkers int
	jobs       chan Job
	results    chan Result
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, 100),
		results:    make(chan Result, 100),
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 1; i <= wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(ctx, i)
	}
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: Shutting down\n", id)
			return
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}

			// Process job
			start := time.Now()
			output := job.Process(job.Data)
			duration := time.Since(start)

			fmt.Printf("Worker %d: Processed job %d in %v\n", id, job.ID, duration)

			wp.results <- Result{
				JobID:  job.ID,
				Output: output,
				Error:  nil,
			}
		}
	}
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.jobs <- job
}

func (wp *WorkerPool) Close() {
	close(wp.jobs)
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.results)
}

func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}

// ===== Part 2: Pipeline Pattern =====

func pipeline(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func double(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * 2:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// ===== Part 3: Fan-Out, Fan-In Pattern =====

func fanOut(ctx context.Context, in <-chan int, numWorkers int) []<-chan int {
	channels := make([]<-chan int, numWorkers)

	for i := 0; i < numWorkers; i++ {
		channels[i] = processWithDelay(ctx, in)
	}

	return channels
}

func processWithDelay(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for n := range in {
			// Simulate processing
			time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

			select {
			case out <- n * 2:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func fanIn(ctx context.Context, channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	// Start a goroutine for each input channel
	for _, c := range channels {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for n := range ch {
				select {
				case out <- n:
				case <-ctx.Done():
					return
				}
			}
		}(c)
	}

	// Close out channel when all inputs are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ===== Part 4: Rate Limiter with Token Bucket =====

type TokenBucket struct {
	capacity int
	tokens   int
	refill   time.Duration
	mu       sync.Mutex
}

func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
	tb := &TokenBucket{
		capacity: capacity,
		tokens:   capacity,
		refill:   refillRate,
	}

	// Start refill goroutine
	go tb.refillTokens()

	return tb
}

func (tb *TokenBucket) refillTokens() {
	ticker := time.NewTicker(tb.refill)
	defer ticker.Stop()

	for range ticker.C {
		tb.mu.Lock()
		if tb.tokens < tb.capacity {
			tb.tokens++
		}
		tb.mu.Unlock()
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

func (tb *TokenBucket) Wait(ctx context.Context) error {
	for {
		if tb.Allow() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Millisecond):
			continue
		}
	}
}

// ===== Part 5: Semaphore Pattern =====

type Semaphore struct {
	permits chan struct{}
}

func NewSemaphore(maxPermits int) *Semaphore {
	return &Semaphore{
		permits: make(chan struct{}, maxPermits),
	}
}

func (s *Semaphore) Acquire() {
	s.permits <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.permits
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case s.permits <- struct{}{}:
		return true
	default:
		return false
	}
}

// ===== Part 6: Future/Promise Pattern =====

type Future struct {
	result chan interface{}
	err    chan error
}

func NewFuture(fn func() (interface{}, error)) *Future {
	f := &Future{
		result: make(chan interface{}, 1),
		err:    make(chan error, 1),
	}

	go func() {
		result, err := fn()
		if err != nil {
			f.err <- err
		} else {
			f.result <- result
		}
	}()

	return f
}

func (f *Future) Get() (interface{}, error) {
	select {
	case result := <-f.result:
		return result, nil
	case err := <-f.err:
		return nil, err
	}
}

func (f *Future) GetWithTimeout(timeout time.Duration) (interface{}, error) {
	select {
	case result := <-f.result:
		return result, nil
	case err := <-f.err:
		return nil, err
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout after %v", timeout)
	}
}

// ===== Part 7: Broadcast Pattern =====

type Broadcaster struct {
	listeners []chan interface{}
	mu        sync.RWMutex
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		listeners: make([]chan interface{}, 0),
	}
}

func (b *Broadcaster) Subscribe() <-chan interface{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan interface{}, 10)
	b.listeners = append(b.listeners, ch)
	return ch
}

func (b *Broadcaster) Broadcast(data interface{}) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.listeners {
		select {
		case ch <- data:
		default:
			// Skip if channel is full
		}
	}
}

func (b *Broadcaster) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, ch := range b.listeners {
		close(ch)
	}
}

// ===== Part 8: Bounded Parallelism =====

func boundedParallel(ctx context.Context, maxGoroutines int, tasks []func() error) []error {
	sem := NewSemaphore(maxGoroutines)
	errors := make([]error, len(tasks))
	var wg sync.WaitGroup

	for i, task := range tasks {
		wg.Add(1)
		go func(index int, t func() error) {
			defer wg.Done()

			sem.Acquire()
			defer sem.Release()

			select {
			case <-ctx.Done():
				errors[index] = ctx.Err()
			default:
				errors[index] = t()
			}
		}(i, task)
	}

	wg.Wait()
	return errors
}

// ===== Main Demo =====

func main() {
	fmt.Println("===== Advanced Concurrency Patterns =====\n")

	// ===== Demo 1: Worker Pool =====
	fmt.Println("--- 1. Worker Pool Pattern ---")

	ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel1()

	pool := NewWorkerPool(3)
	pool.Start(ctx1)

	// Add jobs
	for i := 1; i <= 10; i++ {
		job := Job{
			ID:   i,
			Data: fmt.Sprintf("Task-%d", i),
			Process: func(data string) string {
				time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
				return fmt.Sprintf("Processed: %s", data)
			},
		}
		pool.AddJob(job)
	}

	pool.Close()

	// Collect results
	go func() {
		for result := range pool.Results() {
			fmt.Printf("Result %d: %s\n", result.JobID, result.Output)
		}
	}()

	pool.Wait()
	fmt.Println()

	// ===== Demo 2: Pipeline Pattern =====
	fmt.Println("--- 2. Pipeline Pattern ---")

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	// Create pipeline: generate -> square -> double
	nums := pipeline(ctx2, 1, 2, 3, 4, 5)
	squared := square(ctx2, nums)
	doubled := double(ctx2, squared)

	fmt.Print("Pipeline (n -> n² -> n*2): ")
	for result := range doubled {
		fmt.Printf("%d ", result)
	}
	fmt.Println("\n")

	// ===== Demo 3: Fan-Out, Fan-In =====
	fmt.Println("--- 3. Fan-Out, Fan-In Pattern ---")

	ctx3, cancel3 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel3()

	input := pipeline(ctx3, 1, 2, 3, 4, 5, 6, 7, 8)
	workers := fanOut(ctx3, input, 3)
	output := fanIn(ctx3, workers...)

	fmt.Print("Fan-Out/Fan-In results: ")
	for result := range output {
		fmt.Printf("%d ", result)
	}
	fmt.Println("\n")

	// ===== Demo 4: Token Bucket Rate Limiter =====
	fmt.Println("--- 4. Token Bucket Rate Limiter ---")

	limiter := NewTokenBucket(3, 500*time.Millisecond)

	for i := 1; i <= 10; i++ {
		if limiter.Allow() {
			fmt.Printf("Request %d: Allowed ✓\n", i)
		} else {
			fmt.Printf("Request %d: Rate limited ✗\n", i)
		}
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Println()

	// ===== Demo 5: Semaphore Pattern =====
	fmt.Println("--- 5. Semaphore Pattern (Max 2 concurrent) ---")

	sem := NewSemaphore(2)
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			sem.Acquire()
			defer sem.Release()

			fmt.Printf("Task %d: Started\n", id)
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Task %d: Completed\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println()

	// ===== Demo 6: Future/Promise Pattern =====
	fmt.Println("--- 6. Future/Promise Pattern ---")

	future1 := NewFuture(func() (interface{}, error) {
		time.Sleep(500 * time.Millisecond)
		return "Task 1 completed", nil
	})

	future2 := NewFuture(func() (interface{}, error) {
		time.Sleep(300 * time.Millisecond)
		return 42, nil
	})

	result1, _ := future1.Get()
	fmt.Printf("Future 1: %v\n", result1)

	result2, _ := future2.GetWithTimeout(1 * time.Second)
	fmt.Printf("Future 2: %v\n", result2)
	fmt.Println()

	// ===== Demo 7: Broadcast Pattern =====
	fmt.Println("--- 7. Broadcast Pattern ---")

	broadcaster := NewBroadcaster()

	// Create 3 subscribers
	listener1 := broadcaster.Subscribe()
	listener2 := broadcaster.Subscribe()
	listener3 := broadcaster.Subscribe()

	// Listen for broadcasts
	var wgBroadcast sync.WaitGroup
	wgBroadcast.Add(3)

	go func() {
		defer wgBroadcast.Done()
		for msg := range listener1 {
			fmt.Printf("Listener 1 received: %v\n", msg)
		}
	}()

	go func() {
		defer wgBroadcast.Done()
		for msg := range listener2 {
			fmt.Printf("Listener 2 received: %v\n", msg)
		}
	}()

	go func() {
		defer wgBroadcast.Done()
		for msg := range listener3 {
			fmt.Printf("Listener 3 received: %v\n", msg)
		}
	}()

	// Broadcast messages
	time.Sleep(100 * time.Millisecond)
	broadcaster.Broadcast("Hello")
	time.Sleep(100 * time.Millisecond)
	broadcaster.Broadcast("World")
	time.Sleep(100 * time.Millisecond)

	broadcaster.Close()
	wgBroadcast.Wait()
	fmt.Println()

	// ===== Demo 8: Bounded Parallelism =====
	fmt.Println("--- 8. Bounded Parallelism (Max 2 concurrent) ---")

	ctx8, cancel8 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel8()

	tasks := make([]func() error, 6)
	for i := range tasks {
		taskID := i + 1
		tasks[i] = func() error {
			fmt.Printf("Task %d: Started\n", taskID)
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Task %d: Completed\n", taskID)
			return nil
		}
	}

	errors := boundedParallel(ctx8, 2, tasks)

	fmt.Printf("\nCompleted with %d errors\n", countErrors(errors))

	// ===== Summary =====
	fmt.Println("\n" + repeatString("=", 60))
	fmt.Println("Advanced Concurrency Patterns Summary:")
	fmt.Println("  ✓ Worker Pool: Distribute work across fixed workers")
	fmt.Println("  ✓ Pipeline: Chain processing stages")
	fmt.Println("  ✓ Fan-Out/Fan-In: Parallel processing and merge")
	fmt.Println("  ✓ Token Bucket: Rate limiting with refill")
	fmt.Println("  ✓ Semaphore: Limit concurrent access")
	fmt.Println("  ✓ Future/Promise: Async result handling")
	fmt.Println("  ✓ Broadcast: Multi-subscriber messaging")
	fmt.Println("  ✓ Bounded Parallelism: Control max goroutines")
	fmt.Println(repeatString("=", 60))
}

func countErrors(errors []error) int {
	count := 0
	for _, err := range errors {
		if err != nil {
			count++
		}
	}
	return count
}

func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

// วิธีรัน:
// go run 05_advanced_concurrency.go
//
// Key Concepts:
//
// 1. Worker Pool:
//    - Fixed number of workers
//    - Job queue distribution
//    - Result collection
//    - Graceful shutdown
//
// 2. Pipeline:
//    - Sequential processing stages
//    - Each stage is a goroutine
//    - Data flows through channels
//    - Context cancellation
//
// 3. Fan-Out/Fan-In:
//    - Fan-Out: Distribute to multiple workers
//    - Fan-In: Merge results from workers
//    - Parallel processing
//    - Result aggregation
//
// 4. Token Bucket:
//    - Rate limiting algorithm
//    - Token refill over time
//    - Burst handling
//    - Non-blocking and blocking modes
//
// 5. Semaphore:
//    - Limit concurrent access
//    - Resource pool management
//    - Acquire/Release pattern
//    - Try-acquire for non-blocking
//
// 6. Future/Promise:
//    - Async computation
//    - Deferred result
//    - Timeout support
//    - Error handling
//
// 7. Broadcast:
//    - One-to-many messaging
//    - Multiple subscribers
//    - Non-blocking send
//    - Clean shutdown
//
// 8. Bounded Parallelism:
//    - Control max goroutines
//    - Resource management
//    - Error collection
//    - Context support
//
// Best Practices:
// - Always use context for cancellation
// - Close channels properly
// - Use WaitGroups for synchronization
// - Handle channel full scenarios
// - Implement graceful shutdown
// - Use buffered channels appropriately
// - Avoid goroutine leaks
// - Use select with default for non-blocking
//
// Production Tips:
// - Monitor goroutine count
// - Implement metrics
// - Add logging
// - Handle errors properly
// - Test with race detector: go run -race
// - Use pprof for profiling
// - Set proper timeouts
// - Implement circuit breakers
