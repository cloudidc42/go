// ===== ตัวอย่าง: Advanced Caching Strategies =====
package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// ===== Part 1: LRU Cache (Least Recently Used) =====

type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	lruList  *list.List
	mu       sync.RWMutex
}

type cacheEntry struct {
	key   string
	value interface{}
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		lruList:  list.New(),
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.cache[key]; exists {
		// Move to front (most recently used)
		c.lruList.MoveToFront(elem)
		return elem.Value.(*cacheEntry).value, true
	}

	return nil, false
}

func (c *LRUCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If key exists, update and move to front
	if elem, exists := c.cache[key]; exists {
		c.lruList.MoveToFront(elem)
		elem.Value.(*cacheEntry).value = value
		return
	}

	// Add new entry
	entry := &cacheEntry{key: key, value: value}
	elem := c.lruList.PushFront(entry)
	c.cache[key] = elem

	// Evict if over capacity
	if c.lruList.Len() > c.capacity {
		c.evict()
	}
}

func (c *LRUCache) evict() {
	elem := c.lruList.Back()
	if elem != nil {
		c.lruList.Remove(elem)
		entry := elem.Value.(*cacheEntry)
		delete(c.cache, entry.key)
		fmt.Printf("  [LRU] Evicted: %s\n", entry.key)
	}
}

func (c *LRUCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.cache[key]; exists {
		c.lruList.Remove(elem)
		delete(c.cache, key)
	}
}

func (c *LRUCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lruList.Len()
}

func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*list.Element)
	c.lruList = list.New()
}

// ===== Part 2: TTL Cache (Time To Live) =====

type TTLCache struct {
	cache map[string]*ttlEntry
	mu    sync.RWMutex
}

type ttlEntry struct {
	value      interface{}
	expiration time.Time
}

func NewTTLCache() *TTLCache {
	cache := &TTLCache{
		cache: make(map[string]*ttlEntry),
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

func (c *TTLCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = &ttlEntry{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

func (c *TTLCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.expiration) {
		return nil, false
	}

	return entry.value, true
}

func (c *TTLCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
}

func (c *TTLCache) cleanup() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.cache {
			if now.After(entry.expiration) {
				delete(c.cache, key)
				fmt.Printf("  [TTL] Expired: %s\n", key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *TTLCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

// ===== Part 3: Write-Through Cache =====

type WriteThroughCache struct {
	cache   map[string]interface{}
	storage map[string]interface{} // Simulated database
	mu      sync.RWMutex
}

func NewWriteThroughCache() *WriteThroughCache {
	return &WriteThroughCache{
		cache:   make(map[string]interface{}),
		storage: make(map[string]interface{}),
	}
}

func (c *WriteThroughCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Write to cache
	c.cache[key] = value

	// Write to storage (simulated database write)
	c.storage[key] = value

	fmt.Printf("  [Write-Through] Wrote to cache and storage: %s\n", key)
}

func (c *WriteThroughCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Try cache first
	if value, exists := c.cache[key]; exists {
		fmt.Printf("  [Write-Through] Cache HIT: %s\n", key)
		return value, true
	}

	// Cache miss, load from storage
	if value, exists := c.storage[key]; exists {
		fmt.Printf("  [Write-Through] Cache MISS, loaded from storage: %s\n", key)
		// Note: In production, you'd write back to cache here
		return value, true
	}

	return nil, false
}

// ===== Part 4: Write-Behind (Write-Back) Cache =====

type WriteBehindCache struct {
	cache       map[string]interface{}
	storage     map[string]interface{}
	dirtyKeys   map[string]bool
	mu          sync.RWMutex
	flushTicker *time.Ticker
}

func NewWriteBehindCache(flushInterval time.Duration) *WriteBehindCache {
	cache := &WriteBehindCache{
		cache:       make(map[string]interface{}),
		storage:     make(map[string]interface{}),
		dirtyKeys:   make(map[string]bool),
		flushTicker: time.NewTicker(flushInterval),
	}

	// Start flush goroutine
	go cache.flushLoop()

	return cache
}

func (c *WriteBehindCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Write to cache immediately
	c.cache[key] = value
	c.dirtyKeys[key] = true

	fmt.Printf("  [Write-Behind] Wrote to cache (deferred storage): %s\n", key)
}

func (c *WriteBehindCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, exists := c.cache[key]; exists {
		return value, true
	}

	if value, exists := c.storage[key]; exists {
		return value, true
	}

	return nil, false
}

func (c *WriteBehindCache) flushLoop() {
	for range c.flushTicker.C {
		c.flush()
	}
}

func (c *WriteBehindCache) flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	for key := range c.dirtyKeys {
		if value, exists := c.cache[key]; exists {
			c.storage[key] = value
			count++
		}
		delete(c.dirtyKeys, key)
	}

	if count > 0 {
		fmt.Printf("  [Write-Behind] Flushed %d items to storage\n", count)
	}
}

func (c *WriteBehindCache) Close() {
	c.flushTicker.Stop()
	c.flush() // Final flush
}

// ===== Part 5: Multi-Level Cache =====

type MultiLevelCache struct {
	l1 *LRUCache          // Fast, small cache
	l2 *TTLCache          // Slower, larger cache
	mu sync.RWMutex
}

func NewMultiLevelCache(l1Capacity int) *MultiLevelCache {
	return &MultiLevelCache{
		l1: NewLRUCache(l1Capacity),
		l2: NewTTLCache(),
	}
}

func (c *MultiLevelCache) Get(key string) (interface{}, bool) {
	// Try L1 cache first
	if value, exists := c.l1.Get(key); exists {
		fmt.Printf("  [Multi-Level] L1 HIT: %s\n", key)
		return value, true
	}

	// Try L2 cache
	if value, exists := c.l2.Get(key); exists {
		fmt.Printf("  [Multi-Level] L2 HIT: %s\n", key)
		// Promote to L1
		c.l1.Set(key, value)
		return value, true
	}

	fmt.Printf("  [Multi-Level] MISS: %s\n", key)
	return nil, false
}

func (c *MultiLevelCache) Set(key string, value interface{}, ttl time.Duration) {
	// Write to both levels
	c.l1.Set(key, value)
	c.l2.Set(key, value, ttl)
}

// ===== Part 6: Cache with Loading Function =====

type LoadingCache struct {
	cache  map[string]interface{}
	loader func(key string) (interface{}, error)
	mu     sync.RWMutex
}

func NewLoadingCache(loader func(string) (interface{}, error)) *LoadingCache {
	return &LoadingCache{
		cache:  make(map[string]interface{}),
		loader: loader,
	}
}

func (c *LoadingCache) Get(key string) (interface{}, error) {
	// Try cache
	c.mu.RLock()
	if value, exists := c.cache[key]; exists {
		c.mu.RUnlock()
		fmt.Printf("  [Loading] Cache HIT: %s\n", key)
		return value, nil
	}
	c.mu.RUnlock()

	// Cache miss, load value
	fmt.Printf("  [Loading] Cache MISS, loading: %s\n", key)
	value, err := c.loader(key)
	if err != nil {
		return nil, err
	}

	// Store in cache
	c.mu.Lock()
	c.cache[key] = value
	c.mu.Unlock()

	return value, nil
}

// ===== Part 7: Cache Statistics =====

type CacheStats struct {
	hits      uint64
	misses    uint64
	evictions uint64
	mu        sync.RWMutex
}

func NewCacheStats() *CacheStats {
	return &CacheStats{}
}

func (s *CacheStats) RecordHit() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hits++
}

func (s *CacheStats) RecordMiss() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.misses++
}

func (s *CacheStats) RecordEviction() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evictions++
}

func (s *CacheStats) HitRate() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := s.hits + s.misses
	if total == 0 {
		return 0
	}

	return float64(s.hits) / float64(total) * 100
}

func (s *CacheStats) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return fmt.Sprintf("Hits: %d, Misses: %d, Evictions: %d, Hit Rate: %.2f%%",
		s.hits, s.misses, s.evictions, s.HitRate())
}

// ===== Main Demo =====

func main() {
	fmt.Println("===== Advanced Caching Strategies =====\n")

	// ===== Demo 1: LRU Cache =====
	fmt.Println("--- 1. LRU Cache (Capacity: 3) ---")

	lru := NewLRUCache(3)

	lru.Set("A", "Value A")
	lru.Set("B", "Value B")
	lru.Set("C", "Value C")
	fmt.Printf("Cache size: %d\n", lru.Len())

	// Access A to make it recently used
	if val, exists := lru.Get("A"); exists {
		fmt.Printf("Get A: %v\n", val)
	}

	// Add D, should evict B (least recently used)
	lru.Set("D", "Value D")
	fmt.Printf("Cache size: %d\n", lru.Len())

	// Try to get B (should not exist)
	if _, exists := lru.Get("B"); !exists {
		fmt.Println("B was evicted ✓")
	}
	fmt.Println()

	// ===== Demo 2: TTL Cache =====
	fmt.Println("--- 2. TTL Cache ---")

	ttl := NewTTLCache()

	ttl.Set("session1", "user123", 2*time.Second)
	ttl.Set("session2", "user456", 5*time.Second)

	fmt.Printf("Cache size: %d\n", ttl.Len())

	// Get immediately
	if val, exists := ttl.Get("session1"); exists {
		fmt.Printf("Get session1: %v\n", val)
	}

	// Wait for expiration
	fmt.Println("Waiting 3 seconds...")
	time.Sleep(3 * time.Second)

	// session1 should be expired
	if _, exists := ttl.Get("session1"); !exists {
		fmt.Println("session1 expired ✓")
	}

	// session2 should still exist
	if val, exists := ttl.Get("session2"); exists {
		fmt.Printf("session2 still valid: %v\n", val)
	}
	fmt.Println()

	// ===== Demo 3: Write-Through Cache =====
	fmt.Println("--- 3. Write-Through Cache ---")

	wtCache := NewWriteThroughCache()

	wtCache.Set("user:1", map[string]string{"name": "Alice", "role": "admin"})
	wtCache.Set("user:2", map[string]string{"name": "Bob", "role": "user"})

	// Get from cache (hit)
	wtCache.Get("user:1")

	// Clear cache to simulate cache miss
	wtCache.cache = make(map[string]interface{})

	// Get again (miss, load from storage)
	wtCache.Get("user:1")
	fmt.Println()

	// ===== Demo 4: Write-Behind Cache =====
	fmt.Println("--- 4. Write-Behind Cache ---")

	wbCache := NewWriteBehindCache(2 * time.Second)
	defer wbCache.Close()

	wbCache.Set("product:1", "Laptop")
	wbCache.Set("product:2", "Mouse")
	wbCache.Set("product:3", "Keyboard")

	fmt.Println("Waiting for flush...")
	time.Sleep(3 * time.Second)
	fmt.Println()

	// ===== Demo 5: Multi-Level Cache =====
	fmt.Println("--- 5. Multi-Level Cache ---")

	mlCache := NewMultiLevelCache(2)

	mlCache.Set("data:1", "Fast Data", 10*time.Second)
	mlCache.Set("data:2", "More Data", 10*time.Second)
	mlCache.Set("data:3", "Extra Data", 10*time.Second) // Will evict from L1

	// Try to get data:1 (should be in L2, promoted to L1)
	mlCache.Get("data:1")

	// Try again (now in L1)
	mlCache.Get("data:1")
	fmt.Println()

	// ===== Demo 6: Loading Cache =====
	fmt.Println("--- 6. Loading Cache ---")

	loadingCache := NewLoadingCache(func(key string) (interface{}, error) {
		// Simulate database query
		time.Sleep(100 * time.Millisecond)
		return fmt.Sprintf("Loaded data for %s", key), nil
	})

	// First call - loads from source
	loadingCache.Get("user:100")

	// Second call - cached
	loadingCache.Get("user:100")
	fmt.Println()

	// ===== Demo 7: Cache with Statistics =====
	fmt.Println("--- 7. Cache Statistics ---")

	stats := NewCacheStats()
	simpleCache := make(map[string]string)

	// Simulate cache operations
	keys := []string{"key1", "key2", "key3", "key1", "key2", "key4"}
	simpleCache["key1"] = "value1"
	simpleCache["key2"] = "value2"

	for _, key := range keys {
		if _, exists := simpleCache[key]; exists {
			stats.RecordHit()
			fmt.Printf("Get %s: HIT\n", key)
		} else {
			stats.RecordMiss()
			fmt.Printf("Get %s: MISS\n", key)
		}
	}

	fmt.Printf("\n%s\n", stats.String())
	fmt.Println()

	// ===== Summary =====
	fmt.Println(repeatStr("=", 60))
	fmt.Println("Caching Strategies Summary:")
	fmt.Println("  ✓ LRU: Evicts least recently used items")
	fmt.Println("  ✓ TTL: Items expire after time duration")
	fmt.Println("  ✓ Write-Through: Sync write to cache and storage")
	fmt.Println("  ✓ Write-Behind: Async write to storage (better performance)")
	fmt.Println("  ✓ Multi-Level: L1 (fast) + L2 (larger) caching")
	fmt.Println("  ✓ Loading: Auto-load on cache miss")
	fmt.Println("  ✓ Statistics: Track hit rate and metrics")
	fmt.Println(repeatStr("=", 60))
}

func repeatStr(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

// วิธีรัน:
// go run 06_caching_strategies.go
//
// Caching Strategies Explained:
//
// 1. LRU Cache (Least Recently Used):
//    - Evicts least recently accessed items
//    - O(1) get and set operations
//    - Uses doubly-linked list + hashmap
//    - Good for: General purpose caching
//
// 2. TTL Cache (Time To Live):
//    - Items expire after duration
//    - Background cleanup
//    - Good for: Sessions, temporary data
//
// 3. Write-Through Cache:
//    - Writes to cache AND storage synchronously
//    - Consistency guaranteed
//    - Higher latency
//    - Good for: Strong consistency requirements
//
// 4. Write-Behind Cache (Write-Back):
//    - Writes to cache immediately, storage later
//    - Better performance
//    - Risk of data loss
//    - Good for: High-write workloads
//
// 5. Multi-Level Cache:
//    - L1: Small, fast cache
//    - L2: Larger, slower cache
//    - Promotion on access
//    - Good for: Hierarchical data access
//
// 6. Loading Cache:
//    - Auto-loads on cache miss
//    - Transparent caching
//    - Good for: Read-heavy workloads
//
// 7. Cache Statistics:
//    - Track hit/miss rate
//    - Monitor performance
//    - Good for: Optimization
//
// Best Practices:
// - Choose strategy based on access pattern
// - Monitor hit rates
// - Set appropriate capacity
// - Use TTL for temporary data
// - Implement eviction policies
// - Handle cache stampede
// - Consider cache warming
// - Use compression for large values
// - Implement cache-aside pattern
// - Handle concurrent access
//
// Production Considerations:
// - Distributed caching (Redis, Memcached)
// - Cache invalidation strategy
// - Serialization overhead
// - Memory limits
// - Network latency
// - Consistency requirements
// - Monitoring and alerting
// - Cache penetration protection
// - Circuit breaker integration
