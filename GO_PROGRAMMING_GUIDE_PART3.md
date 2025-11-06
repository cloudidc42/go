# คู่มือการเขียนโปรแกรม Go - ส่วนที่ 3
## Step 181-500: Maps, Structs, Pointers, Interfaces และ Concurrency

---

## Step 181-210: Maps

### Step 181-190: Map Basics

```go
package main

import "fmt"

func main() {
    // ===== Step 181: Map Declaration =====

    // วิธีที่ 1: var (nil map)
    var map1 map[string]int
    fmt.Printf("nil map: %v, is nil: %t\n", map1, map1 == nil)

    // วิธีที่ 2: make
    map2 := make(map[string]int)
    fmt.Printf("empty map: %v, is nil: %t\n", map2, map2 == nil)

    // วิธีที่ 3: map literal
    map3 := map[string]int{
        "apple":  5,
        "banana": 3,
        "orange": 7,
    }
    fmt.Println("literal map:", map3)

    // วิธีที่ 4: make with capacity hint
    map4 := make(map[string]int, 100)
    _ = map4

    // ===== Step 182: Add and Update =====
    fruits := make(map[string]int)

    // Add
    fruits["apple"] = 5
    fruits["banana"] = 3
    fruits["orange"] = 7

    fmt.Println("After adding:", fruits)

    // Update
    fruits["apple"] = 10
    fmt.Println("After update:", fruits)

    // ===== Step 183: Read from Map =====
    // Method 1: direct access
    value := fruits["apple"]
    fmt.Println("Apples:", value)

    // Method 2: with existence check
    value2, exists := fruits["grape"]
    if exists {
        fmt.Println("Grapes:", value2)
    } else {
        fmt.Println("Grapes not found")
    }

    // ===== Step 184: Delete from Map =====
    delete(fruits, "banana")
    fmt.Println("After delete:", fruits)

    // Delete non-existent key (no error)
    delete(fruits, "watermelon")

    // ===== Step 185: Iterate over Map =====
    scores := map[string]int{
        "Alice": 95,
        "Bob":   87,
        "Carol": 92,
        "Dave":  78,
    }

    // Keys and values
    for name, score := range scores {
        fmt.Printf("%s: %d\n", name, score)
    }

    // Keys only
    fmt.Println("Names:")
    for name := range scores {
        fmt.Println("-", name)
    }

    // ===== Step 186: Map Length =====
    fmt.Printf("Map has %d entries\n", len(scores))

    // ===== Step 187: Check if Key Exists =====
    checkKeys := []string{"Alice", "Eve", "Bob"}

    for _, key := range checkKeys {
        if _, exists := scores[key]; exists {
            fmt.Printf("✓ %s exists\n", key)
        } else {
            fmt.Printf("✗ %s does not exist\n", key)
        }
    }

    // ===== Step 188: Map of Slices =====
    groups := map[string][]string{
        "fruits":     {"apple", "banana", "orange"},
        "vegetables": {"carrot", "broccoli", "spinach"},
        "meats":      {"chicken", "beef", "pork"},
    }

    for category, items := range groups {
        fmt.Printf("%s: %v\n", category, items)
    }

    // Add to slice in map
    groups["fruits"] = append(groups["fruits"], "mango")
    fmt.Println("Updated fruits:", groups["fruits"])

    // ===== Step 189: Map of Maps =====
    students := map[string]map[string]int{
        "Alice": {
            "Math":    95,
            "English": 88,
            "Science": 92,
        },
        "Bob": {
            "Math":    78,
            "English": 85,
            "Science": 80,
        },
    }

    // Access nested value
    fmt.Println("Alice's Math score:", students["Alice"]["Math"])

    // Add new student
    students["Carol"] = map[string]int{
        "Math":    90,
        "English": 87,
        "Science": 94,
    }

    // ===== Step 190: Count Occurrences =====
    words := []string{"apple", "banana", "apple", "orange", "banana", "apple"}

    count := make(map[string]int)
    for _, word := range words {
        count[word]++
    }

    fmt.Println("Word count:", count)
}
```

### Step 191-200: Advanced Map Operations

```go
package main

import (
    "fmt"
    "sort"
)

func main() {
    // ===== Step 191: Copy Map =====
    original := map[string]int{
        "a": 1,
        "b": 2,
        "c": 3,
    }

    // Deep copy
    copied := make(map[string]int)
    for k, v := range original {
        copied[k] = v
    }

    copied["a"] = 999
    fmt.Println("Original:", original)
    fmt.Println("Copied:", copied)

    // ===== Step 192: Merge Maps =====
    map1 := map[string]int{"a": 1, "b": 2}
    map2 := map[string]int{"b": 3, "c": 4}

    merged := mergeMaps(map1, map2)
    fmt.Println("Merged:", merged)

    // ===== Step 193: Filter Map =====
    numbers := map[string]int{
        "one":   1,
        "two":   2,
        "three": 3,
        "four":  4,
        "five":  5,
    }

    evens := filterMap(numbers, func(k string, v int) bool {
        return v%2 == 0
    })
    fmt.Println("Even numbers:", evens)

    // ===== Step 194: Transform Map =====
    doubled := transformMap(numbers, func(k string, v int) int {
        return v * 2
    })
    fmt.Println("Doubled:", doubled)

    // ===== Step 195: Get All Keys =====
    keys := getKeys(numbers)
    fmt.Println("Keys:", keys)

    // ===== Step 196: Get All Values =====
    values := getValues(numbers)
    fmt.Println("Values:", values)

    // ===== Step 197: Sort Map by Keys =====
    data := map[string]int{
        "delta":  4,
        "alpha":  1,
        "gamma":  3,
        "beta":   2,
    }

    sortedByKeys := sortMapByKeys(data)
    fmt.Println("Sorted by keys:", sortedByKeys)

    // ===== Step 198: Sort Map by Values =====
    sortedByValues := sortMapByValues(data)
    fmt.Println("Sorted by values:", sortedByValues)

    // ===== Step 199: Invert Map =====
    original199 := map[string]string{
        "en": "English",
        "th": "Thai",
        "ja": "Japanese",
    }

    inverted := invertMap(original199)
    fmt.Println("Inverted:", inverted)

    // ===== Step 200: Group By =====
    people := []Person{
        {"Alice", "Engineering"},
        {"Bob", "Sales"},
        {"Carol", "Engineering"},
        {"Dave", "Sales"},
        {"Eve", "Marketing"},
    }

    grouped := groupBy(people, func(p Person) string {
        return p.Department
    })

    for dept, members := range grouped {
        fmt.Printf("%s: %v\n", dept, members)
    }
}

type Person struct {
    Name       string
    Department string
}

func mergeMaps(maps ...map[string]int) map[string]int {
    result := make(map[string]int)

    for _, m := range maps {
        for k, v := range m {
            result[k] = v
        }
    }

    return result
}

func filterMap(m map[string]int, predicate func(string, int) bool) map[string]int {
    result := make(map[string]int)

    for k, v := range m {
        if predicate(k, v) {
            result[k] = v
        }
    }

    return result
}

func transformMap(m map[string]int, transform func(string, int) int) map[string]int {
    result := make(map[string]int)

    for k, v := range m {
        result[k] = transform(k, v)
    }

    return result
}

func getKeys(m map[string]int) []string {
    keys := make([]string, 0, len(m))

    for k := range m {
        keys = append(keys, k)
    }

    return keys
}

func getValues(m map[string]int) []int {
    values := make([]int, 0, len(m))

    for _, v := range m {
        values = append(values, v)
    }

    return values
}

func sortMapByKeys(m map[string]int) []struct {
    Key   string
    Value int
} {
    keys := getKeys(m)
    sort.Strings(keys)

    result := make([]struct {
        Key   string
        Value int
    }, len(keys))

    for i, k := range keys {
        result[i] = struct {
            Key   string
            Value int
        }{k, m[k]}
    }

    return result
}

func sortMapByValues(m map[string]int) []struct {
    Key   string
    Value int
} {
    type kv struct {
        Key   string
        Value int
    }

    var ss []kv
    for k, v := range m {
        ss = append(ss, kv{k, v})
    }

    sort.Slice(ss, func(i, j int) bool {
        return ss[i].Value < ss[j].Value
    })

    return ss
}

func invertMap(m map[string]string) map[string]string {
    inverted := make(map[string]string)

    for k, v := range m {
        inverted[v] = k
    }

    return inverted
}

func groupBy(people []Person, keyFunc func(Person) string) map[string][]Person {
    groups := make(map[string][]Person)

    for _, person := range people {
        key := keyFunc(person)
        groups[key] = append(groups[key], person)
    }

    return groups
}
```

### Step 201-210: Map Patterns และ Use Cases

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    // ===== Step 201: Set Implementation =====
    set := NewSet()
    set.Add("apple")
    set.Add("banana")
    set.Add("apple") // duplicate

    fmt.Println("Set:", set.Values())
    fmt.Println("Contains apple:", set.Contains("apple"))
    fmt.Println("Size:", set.Size())

    // ===== Step 202: Cache Implementation =====
    cache := NewCache()
    cache.Set("user:1", "John Doe")
    cache.Set("user:2", "Jane Smith")

    if value, found := cache.Get("user:1"); found {
        fmt.Println("Found in cache:", value)
    }

    // ===== Step 203: Frequency Counter =====
    text := "hello world hello go programming"
    freq := wordFrequency(text)
    fmt.Println("Word frequency:", freq)

    // ===== Step 204: Graph Adjacency List =====
    graph := make(map[string][]string)

    graph["A"] = []string{"B", "C"}
    graph["B"] = []string{"D", "E"}
    graph["C"] = []string{"F"}
    graph["D"] = []string{}
    graph["E"] = []string{"F"}
    graph["F"] = []string{}

    fmt.Println("Graph adjacency list:")
    for node, neighbors := range graph {
        fmt.Printf("%s -> %v\n", node, neighbors)
    }

    // ===== Step 205: Memoization =====
    fib := memoizedFib()

    fmt.Println("Fib(40):", fib(40))
    fmt.Println("Fib(40) again:", fib(40)) // faster

    // ===== Step 206: Default Values =====
    config := map[string]string{
        "host": "localhost",
        "port": "8080",
    }

    host := getOrDefault(config, "host", "0.0.0.0")
    port := getOrDefault(config, "port", "80")
    timeout := getOrDefault(config, "timeout", "30")

    fmt.Printf("host=%s, port=%s, timeout=%s\n", host, port, timeout)

    // ===== Step 207: Two Sum Problem =====
    nums := []int{2, 7, 11, 15}
    target := 9

    indices := twoSum(nums, target)
    if len(indices) == 2 {
        fmt.Printf("Two sum indices: %v (nums[%d] + nums[%d] = %d)\n",
            indices, indices[0], indices[1], target)
    }

    // ===== Step 208: Anagram Grouping =====
    words := []string{"eat", "tea", "tan", "ate", "nat", "bat"}
    groups := groupAnagrams(words)

    fmt.Println("Anagram groups:")
    for _, group := range groups {
        fmt.Println(group)
    }

    // ===== Step 209: LRU Cache Simulation =====
    lru := NewLRUCache(3)
    lru.Put(1, "one")
    lru.Put(2, "two")
    lru.Put(3, "three")

    lru.Get(1)             // access 1
    lru.Put(4, "four")     // evicts 2

    fmt.Println("LRU Cache state:", lru.Keys())

    // ===== Step 210: Concurrent Safe Map =====
    safeMap := NewSafeMap()

    // Concurrent writes
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            key := fmt.Sprintf("key%d", n)
            safeMap.Set(key, n)
        }(i)
    }

    wg.Wait()

    fmt.Println("Safe map size:", safeMap.Size())
}

// Step 201: Set
type Set struct {
    data map[string]bool
}

func NewSet() *Set {
    return &Set{data: make(map[string]bool)}
}

func (s *Set) Add(item string) {
    s.data[item] = true
}

func (s *Set) Remove(item string) {
    delete(s.data, item)
}

func (s *Set) Contains(item string) bool {
    return s.data[item]
}

func (s *Set) Values() []string {
    values := make([]string, 0, len(s.data))
    for k := range s.data {
        values = append(values, k)
    }
    return values
}

func (s *Set) Size() int {
    return len(s.data)
}

// Step 202: Cache
type Cache struct {
    data map[string]string
}

func NewCache() *Cache {
    return &Cache{data: make(map[string]string)}
}

func (c *Cache) Set(key, value string) {
    c.data[key] = value
}

func (c *Cache) Get(key string) (string, bool) {
    value, found := c.data[key]
    return value, found
}

// Step 203
func wordFrequency(text string) map[string]int {
    freq := make(map[string]int)
    words := splitWords(text)

    for _, word := range words {
        freq[word]++
    }

    return freq
}

func splitWords(text string) []string {
    var words []string
    var current string

    for _, ch := range text {
        if ch == ' ' {
            if current != "" {
                words = append(words, current)
                current = ""
            }
        } else {
            current += string(ch)
        }
    }

    if current != "" {
        words = append(words, current)
    }

    return words
}

// Step 205
func memoizedFib() func(int) int {
    cache := make(map[int]int)

    var fib func(int) int
    fib = func(n int) int {
        if v, found := cache[n]; found {
            return v
        }

        if n <= 1 {
            return n
        }

        result := fib(n-1) + fib(n-2)
        cache[n] = result
        return result
    }

    return fib
}

// Step 206
func getOrDefault(m map[string]string, key, defaultValue string) string {
    if value, found := m[key]; found {
        return value
    }
    return defaultValue
}

// Step 207
func twoSum(nums []int, target int) []int {
    seen := make(map[int]int)

    for i, num := range nums {
        complement := target - num
        if j, found := seen[complement]; found {
            return []int{j, i}
        }
        seen[num] = i
    }

    return []int{}
}

// Step 208
func groupAnagrams(words []string) [][]string {
    groups := make(map[string][]string)

    for _, word := range words {
        key := sortString(word)
        groups[key] = append(groups[key], word)
    }

    result := make([][]string, 0, len(groups))
    for _, group := range groups {
        result = append(result, group)
    }

    return result
}

func sortString(s string) string {
    runes := []rune(s)
    sort.Slice(runes, func(i, j int) bool {
        return runes[i] < runes[j]
    })
    return string(runes)
}

// Step 209: Simple LRU
type LRUCache struct {
    capacity int
    data     map[int]string
    order    []int
}

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity: capacity,
        data:     make(map[int]string),
        order:    make([]int, 0),
    }
}

func (lru *LRUCache) Get(key int) (string, bool) {
    value, found := lru.data[key]
    if found {
        lru.moveToFront(key)
    }
    return value, found
}

func (lru *LRUCache) Put(key int, value string) {
    if _, found := lru.data[key]; found {
        lru.data[key] = value
        lru.moveToFront(key)
        return
    }

    if len(lru.data) >= lru.capacity {
        // Evict least recently used
        oldest := lru.order[0]
        delete(lru.data, oldest)
        lru.order = lru.order[1:]
    }

    lru.data[key] = value
    lru.order = append(lru.order, key)
}

func (lru *LRUCache) moveToFront(key int) {
    // Remove from current position
    for i, k := range lru.order {
        if k == key {
            lru.order = append(lru.order[:i], lru.order[i+1:]...)
            break
        }
    }
    // Add to end
    lru.order = append(lru.order, key)
}

func (lru *LRUCache) Keys() []int {
    return lru.order
}

// Step 210: Thread-safe Map
type SafeMap struct {
    mu   sync.RWMutex
    data map[string]int
}

func NewSafeMap() *SafeMap {
    return &SafeMap{
        data: make(map[string]int),
    }
}

func (sm *SafeMap) Set(key string, value int) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.data[key] = value
}

func (sm *SafeMap) Get(key string) (int, bool) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    value, found := sm.data[key]
    return value, found
}

func (sm *SafeMap) Size() int {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    return len(sm.data)
}
```

---

## Step 211-250: Structs และ Methods

### Step 211-230: Struct Basics

```go
package main

import "fmt"

func main() {
    // ===== Step 211: Basic Struct =====
    type Person struct {
        Name string
        Age  int
    }

    // Declaration
    var p1 Person
    fmt.Printf("Zero value: %+v\n", p1)

    // With values
    p2 := Person{
        Name: "Alice",
        Age:  30,
    }
    fmt.Printf("p2: %+v\n", p2)

    // Short form (must be in order)
    p3 := Person{"Bob", 25}
    fmt.Printf("p3: %+v\n", p3)

    // ===== Step 212: Access Fields =====
    fmt.Println("Name:", p2.Name)
    fmt.Println("Age:", p2.Age)

    // Modify
    p2.Age = 31
    fmt.Printf("After update: %+v\n", p2)

    // ===== Step 213: Nested Structs =====
    type Address struct {
        Street  string
        City    string
        ZipCode string
    }

    type Employee struct {
        Name    string
        Age     int
        Address Address
    }

    emp := Employee{
        Name: "John Doe",
        Age:  35,
        Address: Address{
            Street:  "123 Main St",
            City:    "Bangkok",
            ZipCode: "10100",
        },
    }

    fmt.Printf("Employee: %+v\n", emp)
    fmt.Println("City:", emp.Address.City)

    // ===== Step 214: Anonymous Fields =====
    type User struct {
        Name  string
        Email string
        Address // embedded
    }

    user := User{
        Name:  "Alice",
        Email: "alice@example.com",
        Address: Address{
            Street:  "456 Oak Ave",
            City:    "Chiang Mai",
            ZipCode: "50000",
        },
    }

    // Can access directly
    fmt.Println("User city:", user.City)

    // ===== Step 215: Pointer to Struct =====
    p := &Person{Name: "Charlie", Age: 28}

    // Auto-dereferencing
    fmt.Println("Name:", p.Name) // same as (*p).Name
    p.Age = 29

    fmt.Printf("Person: %+v\n", p)

    // ===== Step 216: Constructor Pattern =====
    alice := NewPerson("Alice", 30)
    bob := NewPerson("Bob", 25)

    fmt.Printf("Alice: %+v\n", alice)
    fmt.Printf("Bob: %+v\n", bob)

    // ===== Step 217: Methods =====
    rect := Rectangle{Width: 10, Height: 5}

    area := rect.Area()
    perimeter := rect.Perimeter()

    fmt.Printf("Rectangle: %+v\n", rect)
    fmt.Printf("Area: %d, Perimeter: %d\n", area, perimeter)

    // ===== Step 218: Pointer Receivers =====
    circle := Circle{Radius: 5}

    fmt.Println("Before scale:", circle.Radius)
    circle.Scale(2)
    fmt.Println("After scale:", circle.Radius)

    // ===== Step 219: Value vs Pointer Receiver =====
    // Value receiver - doesn't modify original
    p219 := Point{X: 1, Y: 2}
    p219.AddValue(5)
    fmt.Printf("After AddValue: %+v\n", p219) // unchanged

    // Pointer receiver - modifies original
    p219.AddPointer(5)
    fmt.Printf("After AddPointer: %+v\n", p219) // changed

    // ===== Step 220: Method Chaining =====
    builder := NewStringBuilder()
    result := builder.
        Append("Hello").
        Append(" ").
        Append("World").
        String()

    fmt.Println("Built string:", result)

    // ===== Step 221-225: More Complex Structs =====
    // Bank account
    account := NewBankAccount("ACC001", 1000)
    account.Deposit(500)
    account.Withdraw(200)
    account.PrintStatement()

    // ===== Step 226-230: Struct Tags =====
    type Product struct {
        ID    int     `json:"id"`
        Name  string  `json:"name"`
        Price float64 `json:"price"`
    }

    // We'll use these tags in JSON encoding later
}

// Step 216
func NewPerson(name string, age int) *Person {
    return &Person{
        Name: name,
        Age:  age,
    }
}

// Step 217
type Rectangle struct {
    Width  int
    Height int
}

func (r Rectangle) Area() int {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() int {
    return 2 * (r.Width + r.Height)
}

// Step 218
type Circle struct {
    Radius float64
}

func (c *Circle) Scale(factor float64) {
    c.Radius *= factor
}

// Step 219
type Point struct {
    X, Y int
}

func (p Point) AddValue(n int) {
    p.X += n
    p.Y += n
}

func (p *Point) AddPointer(n int) {
    p.X += n
    p.Y += n
}

// Step 220
type StringBuilder struct {
    data string
}

func NewStringBuilder() *StringBuilder {
    return &StringBuilder{}
}

func (sb *StringBuilder) Append(s string) *StringBuilder {
    sb.data += s
    return sb
}

func (sb *StringBuilder) String() string {
    return sb.data
}

// Step 221-225
type BankAccount struct {
    AccountNumber string
    Balance       float64
    Transactions  []Transaction
}

type Transaction struct {
    Type   string
    Amount float64
}

func NewBankAccount(accountNumber string, initialBalance float64) *BankAccount {
    return &BankAccount{
        AccountNumber: accountNumber,
        Balance:       initialBalance,
        Transactions:  []Transaction{},
    }
}

func (ba *BankAccount) Deposit(amount float64) {
    ba.Balance += amount
    ba.Transactions = append(ba.Transactions, Transaction{
        Type:   "DEPOSIT",
        Amount: amount,
    })
}

func (ba *BankAccount) Withdraw(amount float64) bool {
    if amount > ba.Balance {
        return false
    }

    ba.Balance -= amount
    ba.Transactions = append(ba.Transactions, Transaction{
        Type:   "WITHDRAW",
        Amount: amount,
    })

    return true
}

func (ba *BankAccount) PrintStatement() {
    fmt.Printf("\n=== Account Statement ===\n")
    fmt.Printf("Account: %s\n", ba.AccountNumber)
    fmt.Printf("Balance: $%.2f\n", ba.Balance)
    fmt.Println("Transactions:")

    for _, tx := range ba.Transactions {
        fmt.Printf("  %s: $%.2f\n", tx.Type, tx.Amount)
    }
}
```

*ต่อด้วย Step 231-250 และ sections อื่นๆ...*

---

## Step 251-280: Pointers

```go
package main

import "fmt"

func main() {
    // ===== Step 251: Pointer Basics =====
    x := 42

    // Declare pointer
    var p *int
    fmt.Printf("Nil pointer: %v\n", p)

    // Address-of operator (&)
    p = &x
    fmt.Printf("Pointer value: %v\n", p)
    fmt.Printf("Pointed value: %d\n", *p)

    // ===== Step 252: Dereferencing =====
    *p = 100 // modify through pointer
    fmt.Println("x after *p = 100:", x)

    // ===== Step 253: Pointer to Pointer =====
    pp := &p
    fmt.Printf("**pp = %d\n", **pp)

    // ===== Step 254: new() Function =====
    ptr := new(int)
    *ptr = 42
    fmt.Println("Value via new():", *ptr)

    // ===== Step 255-260: Pointer Use Cases =====
    // Swap function
    a, b := 10, 20
    fmt.Printf("Before swap: a=%d, b=%d\n", a, b)
    swap(&a, &b)
    fmt.Printf("After swap: a=%d, b=%d\n", a, b)

    // Modify struct
    person := Person{Name: "Alice", Age: 30}
    updateAge(&person, 31)
    fmt.Printf("Updated person: %+v\n", person)

    // ===== Step 261-270: Pointer Safety =====
    // Nil pointer check
    var nilPtr *int
    if nilPtr != nil {
        fmt.Println(*nilPtr)
    } else {
        fmt.Println("Pointer is nil")
    }

    // ===== Step 271-280: Pointer Patterns =====
    // Optional values
    result := findUser(1)
    if result != nil {
        fmt.Printf("Found: %+v\n", *result)
    } else {
        fmt.Println("User not found")
    }
}

type Person struct {
    Name string
    Age  int
}

func swap(a, b *int) {
    *a, *b = *b, *a
}

func updateAge(p *Person, newAge int) {
    p.Age = newAge
}

func findUser(id int) *Person {
    if id == 1 {
        return &Person{Name: "Alice", Age: 30}
    }
    return nil
}
```

---

## Step 281-310: Interfaces

```go
package main

import (
    "fmt"
    "math"
)

func main() {
    // ===== Step 281: Interface Basics =====
    var s Shape

    s = Circle{Radius: 5}
    fmt.Printf("Circle area: %.2f\n", s.Area())

    s = Rectangle{Width: 10, Height: 5}
    fmt.Printf("Rectangle area: %.2f\n", s.Area())

    // ===== Step 282: Interface with Multiple Methods =====
    var g Geometry = Square{Side: 5}
    fmt.Printf("Square area: %.2f\n", g.Area())
    fmt.Printf("Square perimeter: %.2f\n", g.Perimeter())

    // ===== Step 283: Empty Interface =====
    var i interface{}

    i = 42
    fmt.Printf("i = %v (type: %T)\n", i, i)

    i = "hello"
    fmt.Printf("i = %v (type: %T)\n", i, i)

    i = Circle{Radius: 3}
    fmt.Printf("i = %v (type: %T)\n", i, i)

    // ===== Step 284: Type Assertion =====
    var val interface{} = "hello"

    // Safe type assertion
    if str, ok := val.(string); ok {
        fmt.Println("String value:", str)
    }

    // Type switch
    checkType(42)
    checkType("hello")
    checkType(3.14)
    checkType(true)

    // ===== Step 285-290: Interface Patterns =====
    // Polymorphism
    shapes := []Shape{
        Circle{Radius: 5},
        Rectangle{Width: 10, Height: 5},
        Square{Side: 7},
    }

    totalArea := 0.0
    for _, shape := range shapes {
        totalArea += shape.Area()
    }
    fmt.Printf("Total area: %.2f\n", totalArea)

    // ===== Step 291-300: Common Interfaces =====
    // Stringer interface
    p := Person{Name: "Alice", Age: 30}
    fmt.Println(p) // calls String() method

    // ===== Step 301-310: Advanced Interfaces =====
    // Interface composition
    var rw ReadWriter = &File{Content: "Hello"}
    rw.Write("World")
    content := rw.Read()
    fmt.Println("File content:", content)
}

// Step 281
type Shape interface {
    Area() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Step 282
type Geometry interface {
    Area() float64
    Perimeter() float64
}

type Square struct {
    Side float64
}

func (s Square) Area() float64 {
    return s.Side * s.Side
}

func (s Square) Perimeter() float64 {
    return 4 * s.Side
}

// Step 284
func checkType(val interface{}) {
    switch v := val.(type) {
    case int:
        fmt.Printf("%v is int\n", v)
    case string:
        fmt.Printf("%v is string\n", v)
    case float64:
        fmt.Printf("%v is float64\n", v)
    default:
        fmt.Printf("%v is %T\n", v, v)
    }
}

// Step 291
type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
}

// Step 301
type Reader interface {
    Read() string
}

type Writer interface {
    Write(string)
}

type ReadWriter interface {
    Reader
    Writer
}

type File struct {
    Content string
}

func (f *File) Read() string {
    return f.Content
}

func (f *File) Write(content string) {
    f.Content += content
}
```

เนื่องจากคู่มือนี้ยาวมาก ฉันจะสร้างไฟล์ที่เหลือและตัวอย่างโค้ดที่ใช้งานได้จริง...

