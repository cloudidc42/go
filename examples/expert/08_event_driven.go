// ===== ตัวอย่าง: Event-Driven Architecture =====
package main

import (
	"fmt"
	"sync"
	"time"
)

// ===== Part 1: Event Definition =====

type EventType string

const (
	EventUserCreated     EventType = "user.created"
	EventUserUpdated     EventType = "user.updated"
	EventUserDeleted     EventType = "user.deleted"
	EventOrderPlaced     EventType = "order.placed"
	EventOrderCompleted  EventType = "order.completed"
	EventPaymentReceived EventType = "payment.received"
	EventEmailSent       EventType = "email.sent"
)

type Event struct {
	ID        string
	Type      EventType
	Timestamp time.Time
	Data      map[string]interface{}
}

func NewEvent(eventType EventType, data map[string]interface{}) *Event {
	return &Event{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}
}

// ===== Part 2: Event Handler =====

type EventHandler interface {
	Handle(event *Event) error
}

type EventHandlerFunc func(event *Event) error

func (f EventHandlerFunc) Handle(event *Event) error {
	return f(event)
}

// ===== Part 3: Event Bus =====

type EventBus struct {
	handlers map[EventType][]EventHandler
	mu       sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]EventHandler),
	}
}

func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	fmt.Printf("✓ Subscribed to %s (total handlers: %d)\n", eventType, len(eb.handlers[eventType]))
}

func (eb *EventBus) Publish(event *Event) {
	eb.mu.RLock()
	handlers := eb.handlers[event.Type]
	eb.mu.RUnlock()

	fmt.Printf("\n📢 Publishing event: %s (ID: %s)\n", event.Type, event.ID)

	for _, handler := range handlers {
		// Execute handlers asynchronously
		go func(h EventHandler) {
			if err := h.Handle(event); err != nil {
				fmt.Printf("  ❌ Handler error: %v\n", err)
			}
		}(handler)
	}
}

func (eb *EventBus) PublishSync(event *Event) []error {
	eb.mu.RLock()
	handlers := eb.handlers[event.Type]
	eb.mu.RUnlock()

	fmt.Printf("\n📢 Publishing event (sync): %s\n", event.Type)

	var errors []error
	for _, handler := range handlers {
		if err := handler.Handle(event); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// ===== Part 4: Event Store (Event Sourcing) =====

type EventStore struct {
	events []Event
	mu     sync.RWMutex
}

func NewEventStore() *EventStore {
	return &EventStore{
		events: make([]Event, 0),
	}
}

func (es *EventStore) Append(event *Event) {
	es.mu.Lock()
	defer es.mu.Unlock()

	es.events = append(es.events, *event)
}

func (es *EventStore) GetAll() []Event {
	es.mu.RLock()
	defer es.mu.RUnlock()

	return append([]Event{}, es.events...)
}

func (es *EventStore) GetByType(eventType EventType) []Event {
	es.mu.RLock()
	defer es.mu.RUnlock()

	var filtered []Event
	for _, event := range es.events {
		if event.Type == eventType {
			filtered = append(filtered, event)
		}
	}

	return filtered
}

func (es *EventStore) Count() int {
	es.mu.RLock()
	defer es.mu.RUnlock()

	return len(es.events)
}

// ===== Part 5: Domain Services (Event Handlers) =====

// Email Service
type EmailService struct {
	name string
}

func NewEmailService() *EmailService {
	return &EmailService{name: "EmailService"}
}

func (s *EmailService) SendWelcomeEmail(event *Event) error {
	username := event.Data["username"].(string)
	email := event.Data["email"].(string)

	time.Sleep(100 * time.Millisecond) // Simulate email sending

	fmt.Printf("  📧 [%s] Sent welcome email to %s (%s)\n", s.name, username, email)
	return nil
}

func (s *EmailService) SendOrderConfirmation(event *Event) error {
	orderID := event.Data["order_id"]
	email := event.Data["email"].(string)

	time.Sleep(100 * time.Millisecond)

	fmt.Printf("  📧 [%s] Sent order confirmation for order %v to %s\n", s.name, orderID, email)
	return nil
}

// Analytics Service
type AnalyticsService struct {
	name string
}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{name: "AnalyticsService"}
}

func (s *AnalyticsService) TrackUserCreation(event *Event) error {
	username := event.Data["username"].(string)

	fmt.Printf("  📊 [%s] Tracked user creation: %s\n", s.name, username)
	return nil
}

func (s *AnalyticsService) TrackOrder(event *Event) error {
	orderID := event.Data["order_id"]
	amount := event.Data["amount"]

	fmt.Printf("  📊 [%s] Tracked order: %v (amount: %v)\n", s.name, orderID, amount)
	return nil
}

// Notification Service
type NotificationService struct {
	name string
}

func NewNotificationService() *NotificationService {
	return &NotificationService{name: "NotificationService"}
}

func (s *NotificationService) SendPushNotification(event *Event) error {
	username := event.Data["username"].(string)
	message := fmt.Sprintf("Welcome %s!", username)

	time.Sleep(50 * time.Millisecond)

	fmt.Printf("  🔔 [%s] Sent push notification: %s\n", s.name, message)
	return nil
}

// Payment Service
type PaymentService struct {
	name string
}

func NewPaymentService() *PaymentService {
	return &PaymentService{name: "PaymentService"}
}

func (s *PaymentService) ProcessPayment(event *Event) error {
	orderID := event.Data["order_id"]
	amount := event.Data["amount"]

	time.Sleep(200 * time.Millisecond) // Simulate payment processing

	fmt.Printf("  💳 [%s] Processed payment for order %v: $%v\n", s.name, orderID, amount)
	return nil
}

// Inventory Service
type InventoryService struct {
	name string
}

func NewInventoryService() *InventoryService {
	return &InventoryService{name: "InventoryService"}
}

func (s *InventoryService) ReserveStock(event *Event) error {
	orderID := event.Data["order_id"]
	productID := event.Data["product_id"]

	time.Sleep(150 * time.Millisecond)

	fmt.Printf("  📦 [%s] Reserved stock for order %v (product: %v)\n", s.name, orderID, productID)
	return nil
}

// ===== Part 6: Saga Pattern (Distributed Transaction) =====

type SagaStep struct {
	Name        string
	Execute     func(data map[string]interface{}) error
	Compensate  func(data map[string]interface{}) error
}

type Saga struct {
	steps     []SagaStep
	completed []int
}

func NewSaga() *Saga {
	return &Saga{
		steps:     make([]SagaStep, 0),
		completed: make([]int, 0),
	}
}

func (s *Saga) AddStep(step SagaStep) {
	s.steps = append(s.steps, step)
}

func (s *Saga) Execute(data map[string]interface{}) error {
	fmt.Println("\n🎬 Executing saga...")

	for i, step := range s.steps {
		fmt.Printf("  Step %d: %s\n", i+1, step.Name)

		if err := step.Execute(data); err != nil {
			fmt.Printf("  ❌ Step %d failed: %v\n", i+1, err)
			fmt.Println("\n⏮️  Rolling back saga...")

			// Compensate completed steps in reverse order
			for j := len(s.completed) - 1; j >= 0; j-- {
				completedStep := s.completed[j]
				fmt.Printf("  Compensating step %d: %s\n", completedStep+1, s.steps[completedStep].Name)
				s.steps[completedStep].Compensate(data)
			}

			return err
		}

		s.completed = append(s.completed, i)
	}

	fmt.Println("  ✅ Saga completed successfully")
	return nil
}

// ===== Main Demo =====

func main() {
	fmt.Println("===== Event-Driven Architecture =====\n")

	// Initialize event bus and store
	eventBus := NewEventBus()
	eventStore := NewEventStore()

	// Initialize services
	emailService := NewEmailService()
	analyticsService := NewAnalyticsService()
	notificationService := NewNotificationService()
	paymentService := NewPaymentService()
	inventoryService := NewInventoryService()

	// ===== Setup Event Handlers =====
	fmt.Println("--- Setting up Event Handlers ---")

	// Subscribe to UserCreated event
	eventBus.Subscribe(EventUserCreated, EventHandlerFunc(emailService.SendWelcomeEmail))
	eventBus.Subscribe(EventUserCreated, EventHandlerFunc(analyticsService.TrackUserCreation))
	eventBus.Subscribe(EventUserCreated, EventHandlerFunc(notificationService.SendPushNotification))

	// Subscribe to OrderPlaced event
	eventBus.Subscribe(EventOrderPlaced, EventHandlerFunc(emailService.SendOrderConfirmation))
	eventBus.Subscribe(EventOrderPlaced, EventHandlerFunc(analyticsService.TrackOrder))
	eventBus.Subscribe(EventOrderPlaced, EventHandlerFunc(paymentService.ProcessPayment))
	eventBus.Subscribe(EventOrderPlaced, EventHandlerFunc(inventoryService.ReserveStock))

	// Subscribe to store all events
	eventBus.Subscribe(EventUserCreated, EventHandlerFunc(func(event *Event) error {
		eventStore.Append(event)
		return nil
	}))
	eventBus.Subscribe(EventOrderPlaced, EventHandlerFunc(func(event *Event) error {
		eventStore.Append(event)
		return nil
	}))

	fmt.Println()

	// ===== Demo 1: User Creation Event =====
	fmt.Println("--- Demo 1: User Creation Event ---")

	userCreatedEvent := NewEvent(EventUserCreated, map[string]interface{}{
		"user_id":  1,
		"username": "alice",
		"email":    "alice@example.com",
	})

	eventBus.Publish(userCreatedEvent)
	time.Sleep(500 * time.Millisecond) // Wait for async handlers

	// ===== Demo 2: Order Placement Event =====
	fmt.Println("\n--- Demo 2: Order Placement Event ---")

	orderPlacedEvent := NewEvent(EventOrderPlaced, map[string]interface{}{
		"order_id":   "ORD-001",
		"user_id":    1,
		"product_id": "PROD-123",
		"amount":     99.99,
		"email":      "alice@example.com",
	})

	eventBus.Publish(orderPlacedEvent)
	time.Sleep(500 * time.Millisecond)

	// ===== Demo 3: Multiple Events =====
	fmt.Println("\n--- Demo 3: Multiple Events ---")

	events := []*Event{
		NewEvent(EventUserCreated, map[string]interface{}{
			"user_id":  2,
			"username": "bob",
			"email":    "bob@example.com",
		}),
		NewEvent(EventOrderPlaced, map[string]interface{}{
			"order_id":   "ORD-002",
			"user_id":    2,
			"product_id": "PROD-456",
			"amount":     149.99,
			"email":      "bob@example.com",
		}),
	}

	for _, event := range events {
		eventBus.Publish(event)
		time.Sleep(300 * time.Millisecond)
	}

	time.Sleep(500 * time.Millisecond)

	// ===== Demo 4: Event Store Query =====
	fmt.Println("\n--- Demo 4: Event Store ---")

	allEvents := eventStore.GetAll()
	fmt.Printf("Total events stored: %d\n", len(allEvents))

	userEvents := eventStore.GetByType(EventUserCreated)
	fmt.Printf("User creation events: %d\n", len(userEvents))

	orderEvents := eventStore.GetByType(EventOrderPlaced)
	fmt.Printf("Order placement events: %d\n", len(orderEvents))

	fmt.Println("\nEvent History:")
	for i, event := range allEvents {
		fmt.Printf("  %d. [%s] %s at %s\n",
			i+1, event.Type, event.ID, event.Timestamp.Format("15:04:05"))
	}

	// ===== Demo 5: Saga Pattern =====
	fmt.Println("\n--- Demo 5: Saga Pattern (Distributed Transaction) ---")

	saga := NewSaga()

	orderData := map[string]interface{}{
		"order_id":   "ORD-003",
		"product_id": "PROD-789",
		"amount":     199.99,
	}

	// Add saga steps
	saga.AddStep(SagaStep{
		Name: "Reserve Inventory",
		Execute: func(data map[string]interface{}) error {
			fmt.Println("    → Reserving inventory...")
			time.Sleep(100 * time.Millisecond)
			return nil
		},
		Compensate: func(data map[string]interface{}) error {
			fmt.Println("    ← Releasing inventory...")
			return nil
		},
	})

	saga.AddStep(SagaStep{
		Name: "Process Payment",
		Execute: func(data map[string]interface{}) error {
			fmt.Println("    → Processing payment...")
			time.Sleep(100 * time.Millisecond)
			return nil
		},
		Compensate: func(data map[string]interface{}) error {
			fmt.Println("    ← Refunding payment...")
			return nil
		},
	})

	saga.AddStep(SagaStep{
		Name: "Create Shipment",
		Execute: func(data map[string]interface{}) error {
			fmt.Println("    → Creating shipment...")
			time.Sleep(100 * time.Millisecond)
			// Simulate failure
			return fmt.Errorf("shipment service unavailable")
		},
		Compensate: func(data map[string]interface{}) error {
			fmt.Println("    ← Canceling shipment...")
			return nil
		},
	})

	saga.Execute(orderData)

	// ===== Summary =====
	fmt.Println("\n" + repeatStr("=", 60))
	fmt.Println("Event-Driven Architecture Summary:")
	fmt.Println("  ✓ Event Bus: Publish-subscribe messaging")
	fmt.Println("  ✓ Event Handlers: Async event processing")
	fmt.Println("  ✓ Event Store: Event sourcing and replay")
	fmt.Println("  ✓ Decoupling: Services don't know about each other")
	fmt.Println("  ✓ Scalability: Easy to add new handlers")
	fmt.Println("  ✓ Saga Pattern: Distributed transactions")
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
// go run 08_event_driven.go
//
// Event-Driven Architecture Concepts:
//
// 1. Events:
//    - Immutable facts that happened
//    - Contains type, timestamp, and data
//    - Past tense naming (UserCreated, OrderPlaced)
//
// 2. Event Bus:
//    - Central message broker
//    - Publishers don't know subscribers
//    - Decouples services
//    - Supports multiple handlers per event
//
// 3. Event Handlers:
//    - React to events
//    - Can be async or sync
//    - Idempotent (can handle duplicate events)
//    - Independent from each other
//
// 4. Event Store (Event Sourcing):
//    - Stores all events
//    - Source of truth
//    - Can replay events
//    - Audit trail
//
// 5. Saga Pattern:
//    - Distributed transaction management
//    - Sequence of local transactions
//    - Compensation on failure
//    - Eventually consistent
//
// Benefits:
// - Loose coupling
// - Scalability
// - Flexibility (easy to add features)
// - Audit trail
// - Time travel (replay events)
// - Better testing
//
// Trade-offs:
// - Eventually consistent
// - Complex debugging
// - Event versioning challenges
// - Duplicate event handling
// - Ordering guarantees needed
//
// Best Practices:
// - Use unique event IDs
// - Include timestamp
// - Make events immutable
// - Version events (v1, v2)
// - Handle idempotency
// - Use dead letter queue
// - Monitor event flow
// - Test event handlers independently
//
// Production Considerations:
// - Distributed event bus (Kafka, RabbitMQ, NATS)
// - Event schema registry
// - Retry mechanism
// - Dead letter queue
// - Event replay capability
// - Monitoring and tracing
// - Event versioning strategy
// - Idempotency keys
// - At-least-once delivery
// - Ordering guarantees when needed
