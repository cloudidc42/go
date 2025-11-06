// ===== ตัวอย่าง: Circuit Breaker Pattern =====
package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ===== Part 1: Circuit Breaker States =====

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	return []string{"CLOSED", "OPEN", "HALF-OPEN"}[s]
}

// ===== Part 2: Circuit Breaker Implementation =====

type CircuitBreaker struct {
	maxFailures   int
	resetTimeout  time.Duration
	halfOpenCalls int

	mu              sync.RWMutex
	state           State
	failures        int
	successes       int
	lastFailureTime time.Time
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:   maxFailures,
		resetTimeout:  resetTimeout,
		halfOpenCalls: 3, // Allow 3 calls in half-open state
		state:         StateClosed,
	}
}

func (cb *CircuitBreaker) Call(fn func() error) error {
	// Check if we can make the call
	if !cb.canAttempt() {
		return errors.New("circuit breaker is OPEN")
	}

	// Execute the function
	err := fn()

	// Record the result
	if err != nil {
		cb.recordFailure()
		return err
	}

	cb.recordSuccess()
	return nil
}

func (cb *CircuitBreaker) canAttempt() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Check if circuit should transition from OPEN to HALF-OPEN
	if cb.state == StateOpen {
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			fmt.Println("  🔄 Circuit breaker: OPEN → HALF-OPEN")
			cb.state = StateHalfOpen
			cb.successes = 0
			return true
		}
		return false
	}

	return true
}

func (cb *CircuitBreaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailureTime = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failures >= cb.maxFailures {
			fmt.Printf("  ❌ Circuit breaker: CLOSED → OPEN (failures: %d)\n", cb.failures)
			cb.state = StateOpen
		}

	case StateHalfOpen:
		fmt.Println("  ❌ Circuit breaker: HALF-OPEN → OPEN (test failed)")
		cb.state = StateOpen
		cb.successes = 0
	}
}

func (cb *CircuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateHalfOpen {
		cb.successes++
		if cb.successes >= cb.halfOpenCalls {
			fmt.Printf("  ✅ Circuit breaker: HALF-OPEN → CLOSED (successes: %d)\n", cb.successes)
			cb.state = StateClosed
			cb.failures = 0
			cb.successes = 0
		}
	} else if cb.state == StateClosed {
		// Reset failure count on success
		cb.failures = 0
	}
}

func (cb *CircuitBreaker) GetState() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

func (cb *CircuitBreaker) GetStats() (State, int, int) {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state, cb.failures, cb.successes
}

// ===== Part 3: Retry with Exponential Backoff =====

type RetryConfig struct {
	MaxAttempts int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

func RetryWithBackoff(config RetryConfig, fn func() error) error {
	delay := config.InitialDelay

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			if attempt > 1 {
				fmt.Printf("  ✓ Succeeded after %d attempts\n", attempt)
			}
			return nil
		}

		if attempt == config.MaxAttempts {
			return fmt.Errorf("max retries reached: %w", err)
		}

		fmt.Printf("  ⟳ Retry %d/%d after %v (error: %v)\n",
			attempt, config.MaxAttempts, delay, err)

		time.Sleep(delay)

		// Exponential backoff with jitter
		delay = time.Duration(float64(delay) * config.Multiplier)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}

		// Add jitter (randomness) to avoid thundering herd
		jitter := time.Duration(rand.Int63n(int64(delay) / 10))
		delay += jitter
	}

	return errors.New("max retries reached")
}

// ===== Part 4: Bulkhead Pattern =====

type Bulkhead struct {
	semaphore chan struct{}
	maxWait   time.Duration
}

func NewBulkhead(maxConcurrent int, maxWait time.Duration) *Bulkhead {
	return &Bulkhead{
		semaphore: make(chan struct{}, maxConcurrent),
		maxWait:   maxWait,
	}
}

func (b *Bulkhead) Execute(fn func() error) error {
	// Try to acquire permit
	select {
	case b.semaphore <- struct{}{}:
		defer func() { <-b.semaphore }()
		return fn()
	case <-time.After(b.maxWait):
		return errors.New("bulkhead: max wait time exceeded")
	}
}

// ===== Part 5: Timeout Pattern =====

func WithTimeout(timeout time.Duration, fn func() error) error {
	done := make(chan error, 1)

	go func() {
		done <- fn()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return errors.New("operation timed out")
	}
}

// ===== Part 6: Fallback Pattern =====

func WithFallback(primary func() (interface{}, error), fallback func() (interface{}, error)) (interface{}, error) {
	result, err := primary()
	if err != nil {
		fmt.Printf("  ⚠️  Primary failed: %v, using fallback\n", err)
		return fallback()
	}
	return result, nil
}

// ===== Part 7: Health Check =====

type HealthChecker struct {
	name          string
	checkFunc     func() error
	interval      time.Duration
	timeout       time.Duration
	healthy       bool
	lastCheck     time.Time
	consecutiveFails int
	mu            sync.RWMutex
}

func NewHealthChecker(name string, checkFunc func() error, interval, timeout time.Duration) *HealthChecker {
	hc := &HealthChecker{
		name:      name,
		checkFunc: checkFunc,
		interval:  interval,
		timeout:   timeout,
		healthy:   true,
	}

	go hc.start()

	return hc
}

func (hc *HealthChecker) start() {
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	for range ticker.C {
		hc.check()
	}
}

func (hc *HealthChecker) check() {
	err := WithTimeout(hc.timeout, hc.checkFunc)

	hc.mu.Lock()
	defer hc.mu.Unlock()

	hc.lastCheck = time.Now()

	if err != nil {
		hc.consecutiveFails++
		if hc.consecutiveFails >= 3 && hc.healthy {
			hc.healthy = false
			fmt.Printf("  ❌ %s is now UNHEALTHY (failures: %d)\n", hc.name, hc.consecutiveFails)
		}
	} else {
		if !hc.healthy {
			fmt.Printf("  ✅ %s is now HEALTHY\n", hc.name)
		}
		hc.healthy = true
		hc.consecutiveFails = 0
	}
}

func (hc *HealthChecker) IsHealthy() bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	return hc.healthy
}

// ===== Simulated External Service =====

type ExternalService struct {
	failureRate float64
	latency     time.Duration
	mu          sync.RWMutex
}

func NewExternalService() *ExternalService {
	return &ExternalService{
		failureRate: 0.0,
		latency:     100 * time.Millisecond,
	}
}

func (s *ExternalService) SetFailureRate(rate float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failureRate = rate
}

func (s *ExternalService) SetLatency(latency time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.latency = latency
}

func (s *ExternalService) Call() error {
	s.mu.RLock()
	failureRate := s.failureRate
	latency := s.latency
	s.mu.RUnlock()

	time.Sleep(latency)

	if rand.Float64() < failureRate {
		return errors.New("service unavailable")
	}

	return nil
}

func (s *ExternalService) GetData() (string, error) {
	err := s.Call()
	if err != nil {
		return "", err
	}
	return "Success data", nil
}

// ===== Main Demo =====

func main() {
	fmt.Println("===== Circuit Breaker & Resilience Patterns =====\n")
	rand.Seed(time.Now().UnixNano())

	service := NewExternalService()

	// ===== Demo 1: Circuit Breaker =====
	fmt.Println("--- 1. Circuit Breaker Pattern ---")

	cb := NewCircuitBreaker(3, 3*time.Second)

	// Phase 1: Service is healthy
	fmt.Println("\nPhase 1: Service healthy (failure rate: 0%)")
	for i := 1; i <= 5; i++ {
		err := cb.Call(func() error {
			return service.Call()
		})
		if err != nil {
			fmt.Printf("Call %d: ❌ %v\n", i, err)
		} else {
			fmt.Printf("Call %d: ✅ Success\n", i)
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Phase 2: Service starts failing
	fmt.Println("\nPhase 2: Service degraded (failure rate: 80%)")
	service.SetFailureRate(0.8)

	for i := 1; i <= 8; i++ {
		err := cb.Call(func() error {
			return service.Call()
		})
		state, failures, _ := cb.GetStats()
		if err != nil {
			fmt.Printf("Call %d: ❌ %v (State: %s, Failures: %d)\n", i, err, state, failures)
		} else {
			fmt.Printf("Call %d: ✅ Success (State: %s)\n", i, state)
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Phase 3: Wait for circuit to become half-open
	fmt.Println("\nPhase 3: Waiting for circuit breaker reset...")
	time.Sleep(4 * time.Second)

	// Service recovers
	fmt.Println("\nPhase 4: Service recovered (failure rate: 0%)")
	service.SetFailureRate(0.0)

	for i := 1; i <= 5; i++ {
		err := cb.Call(func() error {
			return service.Call()
		})
		state, _, successes := cb.GetStats()
		if err != nil {
			fmt.Printf("Call %d: ❌ %v (State: %s)\n", i, err, state)
		} else {
			fmt.Printf("Call %d: ✅ Success (State: %s, Successes: %d)\n", i, state, successes)
		}
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println()

	// ===== Demo 2: Retry with Exponential Backoff =====
	fmt.Println("--- 2. Retry with Exponential Backoff ---")

	service.SetFailureRate(0.7) // 70% failure rate

	retryConfig := RetryConfig{
		MaxAttempts:  5,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     2 * time.Second,
		Multiplier:   2.0,
	}

	err := RetryWithBackoff(retryConfig, func() error {
		return service.Call()
	})

	if err != nil {
		fmt.Printf("Final result: ❌ %v\n", err)
	} else {
		fmt.Println("Final result: ✅ Success")
	}
	fmt.Println()

	// ===== Demo 3: Bulkhead Pattern =====
	fmt.Println("--- 3. Bulkhead Pattern (Max 2 concurrent) ---")

	service.SetFailureRate(0.0)
	service.SetLatency(500 * time.Millisecond)

	bulkhead := NewBulkhead(2, 1*time.Second)
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			start := time.Now()
			err := bulkhead.Execute(func() error {
				fmt.Printf("  Task %d: Started\n", id)
				err := service.Call()
				fmt.Printf("  Task %d: Completed in %v\n", id, time.Since(start))
				return err
			})

			if err != nil {
				fmt.Printf("  Task %d: ❌ %v\n", id, err)
			}
		}(i)
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println()

	// ===== Demo 4: Timeout Pattern =====
	fmt.Println("--- 4. Timeout Pattern ---")

	service.SetLatency(2 * time.Second)

	err = WithTimeout(1*time.Second, func() error {
		fmt.Println("  Executing slow operation...")
		return service.Call()
	})

	if err != nil {
		fmt.Printf("Result: ❌ %v\n", err)
	} else {
		fmt.Println("Result: ✅ Success")
	}
	fmt.Println()

	// ===== Demo 5: Fallback Pattern =====
	fmt.Println("--- 5. Fallback Pattern ---")

	service.SetFailureRate(1.0) // Always fail
	service.SetLatency(100 * time.Millisecond)

	result, err := WithFallback(
		func() (interface{}, error) {
			return service.GetData()
		},
		func() (interface{}, error) {
			return "Cached data (fallback)", nil
		},
	)

	if err != nil {
		fmt.Printf("Result: ❌ %v\n", err)
	} else {
		fmt.Printf("Result: ✅ %v\n", result)
	}
	fmt.Println()

	// ===== Demo 6: Health Check =====
	fmt.Println("--- 6. Health Check ---")

	service.SetFailureRate(0.0)

	healthChecker := NewHealthChecker("ExternalService",
		func() error {
			return service.Call()
		},
		1*time.Second,
		500*time.Millisecond,
	)

	fmt.Println("Health check started (checking every 1s)")
	time.Sleep(2 * time.Second)

	// Simulate service failure
	fmt.Println("\nSimulating service failure...")
	service.SetFailureRate(1.0)
	time.Sleep(4 * time.Second)

	// Service recovers
	fmt.Println("\nService recovered...")
	service.SetFailureRate(0.0)
	time.Sleep(2 * time.Second)

	fmt.Printf("\nFinal health status: %v\n", healthChecker.IsHealthy())

	// ===== Summary =====
	fmt.Println("\n" + repeatString("=", 60))
	fmt.Println("Resilience Patterns Summary:")
	fmt.Println("  ✓ Circuit Breaker: Prevent cascading failures")
	fmt.Println("  ✓ Retry + Backoff: Handle transient failures")
	fmt.Println("  ✓ Bulkhead: Isolate resources")
	fmt.Println("  ✓ Timeout: Prevent hanging operations")
	fmt.Println("  ✓ Fallback: Provide alternative responses")
	fmt.Println("  ✓ Health Check: Monitor service health")
	fmt.Println(repeatString("=", 60))
}

func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

// วิธีรัน:
// go run 07_circuit_breaker.go
//
// Resilience Patterns Explained:
//
// 1. Circuit Breaker:
//    States:
//    - CLOSED: Normal operation, requests pass through
//    - OPEN: Failures exceeded threshold, requests blocked
//    - HALF-OPEN: Testing if service recovered
//
//    Benefits:
//    - Prevents cascading failures
//    - Fast failure (no waiting for timeout)
//    - Automatic recovery testing
//
// 2. Retry with Exponential Backoff:
//    - Retry failed operations
//    - Increase delay between retries
//    - Add jitter to prevent thundering herd
//    - Max attempts and max delay
//
// 3. Bulkhead Pattern:
//    - Isolate resources
//    - Limit concurrent requests
//    - Prevent resource exhaustion
//    - Like watertight compartments in a ship
//
// 4. Timeout Pattern:
//    - Set maximum time for operation
//    - Prevent indefinite waiting
//    - Free up resources quickly
//
// 5. Fallback Pattern:
//    - Provide alternative response on failure
//    - Use cached data
//    - Return default values
//    - Graceful degradation
//
// 6. Health Check:
//    - Periodic service monitoring
//    - Detect failures early
//    - Automatic recovery detection
//    - Used by load balancers
//
// Best Practices:
// - Combine multiple patterns
// - Circuit breaker + Retry + Fallback
// - Set appropriate thresholds
// - Monitor metrics (failure rate, latency)
// - Log state transitions
// - Use distributed tracing
// - Test failure scenarios
// - Implement graceful degradation
//
// Production Considerations:
// - Distributed circuit breaker (Redis)
// - Metrics and alerting
// - Dashboard for monitoring
// - Configuration management
// - Rolling window for failure counting
// - Different timeouts for different operations
// - Bulkhead per service
// - Health check endpoints
// - SLA monitoring
