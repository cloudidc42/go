// ===== ตัวอย่าง: WebSocket Real-time Chat =====
// NOTE: This example uses gorilla/websocket patterns but implements a basic version
// For production, install: go get github.com/gorilla/websocket

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ===== Part 1: Message Types =====

type MessageType string

const (
	MessageTypeJoin    MessageType = "join"
	MessageTypeLeave   MessageType = "leave"
	MessageTypeChat    MessageType = "chat"
	MessageTypePrivate MessageType = "private"
	MessageTypeSystem  MessageType = "system"
	MessageTypeTyping  MessageType = "typing"
)

type Message struct {
	Type      MessageType `json:"type"`
	From      string      `json:"from"`
	To        string      `json:"to,omitempty"`
	Content   string      `json:"content"`
	Timestamp time.Time   `json:"timestamp"`
	RoomID    string      `json:"room_id,omitempty"`
}

// ===== Part 2: Client =====

type Client struct {
	ID       string
	Username string
	Room     *Room
	Send     chan *Message
	mu       sync.Mutex
}

func NewClient(id, username string, room *Room) *Client {
	return &Client{
		ID:       id,
		Username: username,
		Room:     room,
		Send:     make(chan *Message, 256),
	}
}

// ===== Part 3: Room (Chat Room) =====

type Room struct {
	ID      string
	Name    string
	Clients map[string]*Client
	mu      sync.RWMutex

	// Message broadcast channel
	Broadcast chan *Message

	// Client registration
	Register chan *Client

	// Client unregistration
	Unregister chan *Client
}

func NewRoom(id, name string) *Room {
	return &Room{
		ID:         id,
		Name:       name,
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.registerClient(client)

		case client := <-r.Unregister:
			r.unregisterClient(client)

		case message := <-r.Broadcast:
			r.broadcastMessage(message)
		}
	}
}

func (r *Room) registerClient(client *Client) {
	r.mu.Lock()
	r.Clients[client.ID] = client
	count := len(r.Clients)
	r.mu.Unlock()

	// Send join message
	joinMsg := &Message{
		Type:      MessageTypeSystem,
		Content:   fmt.Sprintf("%s joined the room (Total users: %d)", client.Username, count),
		Timestamp: time.Now(),
		RoomID:    r.ID,
	}

	r.Broadcast <- joinMsg
	fmt.Printf("[Room %s] %s joined (Total: %d)\n", r.Name, client.Username, count)
}

func (r *Room) unregisterClient(client *Client) {
	r.mu.Lock()
	if _, exists := r.Clients[client.ID]; exists {
		delete(r.Clients, client.ID)
		close(client.Send)
	}
	count := len(r.Clients)
	r.mu.Unlock()

	// Send leave message
	leaveMsg := &Message{
		Type:      MessageTypeSystem,
		Content:   fmt.Sprintf("%s left the room (Total users: %d)", client.Username, count),
		Timestamp: time.Now(),
		RoomID:    r.ID,
	}

	r.Broadcast <- leaveMsg
	fmt.Printf("[Room %s] %s left (Total: %d)\n", r.Name, client.Username, count)
}

func (r *Room) broadcastMessage(message *Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// If it's a private message, send only to recipient
	if message.Type == MessageTypePrivate {
		for _, client := range r.Clients {
			if client.Username == message.To || client.Username == message.From {
				select {
				case client.Send <- message:
				default:
					// Client's send channel is full, skip
				}
			}
		}
		return
	}

	// Broadcast to all clients
	for _, client := range r.Clients {
		select {
		case client.Send <- message:
		default:
			// Client's send channel is full, skip
			close(client.Send)
			delete(r.Clients, client.ID)
		}
	}
}

func (r *Room) GetUserList() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]string, 0, len(r.Clients))
	for _, client := range r.Clients {
		users = append(users, client.Username)
	}
	return users
}

// ===== Part 4: Chat Server =====

type ChatServer struct {
	Rooms map[string]*Room
	mu    sync.RWMutex
}

func NewChatServer() *ChatServer {
	server := &ChatServer{
		Rooms: make(map[string]*Room),
	}

	// Create default rooms
	server.CreateRoom("general", "General Chat")
	server.CreateRoom("tech", "Tech Talk")
	server.CreateRoom("random", "Random")

	return server
}

func (s *ChatServer) CreateRoom(id, name string) *Room {
	s.mu.Lock()
	defer s.mu.Unlock()

	room := NewRoom(id, name)
	s.Rooms[id] = room

	// Start room in goroutine
	go room.Run()

	fmt.Printf("Created room: %s (%s)\n", name, id)
	return room
}

func (s *ChatServer) GetRoom(id string) (*Room, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, exists := s.Rooms[id]
	return room, exists
}

func (s *ChatServer) GetRoomList() []map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rooms := make([]map[string]interface{}, 0, len(s.Rooms))
	for _, room := range s.Rooms {
		room.mu.RLock()
		rooms = append(rooms, map[string]interface{}{
			"id":         room.ID,
			"name":       room.Name,
			"user_count": len(room.Clients),
		})
		room.mu.RUnlock()
	}

	return rooms
}

// ===== Part 5: Simulated WebSocket Connection =====
// In production, you would use actual WebSocket library

type Connection struct {
	client   *Client
	incoming chan *Message
	outgoing chan *Message
	done     chan bool
}

func NewConnection(client *Client) *Connection {
	return &Connection{
		client:   client,
		incoming: make(chan *Message, 10),
		outgoing: make(chan *Message, 10),
		done:     make(chan bool),
	}
}

func (c *Connection) ReadPump() {
	defer func() {
		c.client.Room.Unregister <- c.client
		c.done <- true
	}()

	// Simulate reading messages
	for msg := range c.incoming {
		// Process message
		msg.Timestamp = time.Now()
		msg.RoomID = c.client.Room.ID

		// Broadcast to room
		c.client.Room.Broadcast <- msg
	}
}

func (c *Connection) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.client.Send:
			if !ok {
				return
			}
			// Send to outgoing channel (in real app, this would write to WebSocket)
			c.outgoing <- message

		case <-ticker.C:
			// Send ping (keep-alive)
			pingMsg := &Message{
				Type:      MessageTypeSystem,
				Content:   "ping",
				Timestamp: time.Now(),
			}
			c.outgoing <- pingMsg

		case <-c.done:
			return
		}
	}
}

// ===== Part 6: HTTP Handlers =====

type ChatHandler struct {
	server *ChatServer
}

func NewChatHandler(server *ChatServer) *ChatHandler {
	return &ChatHandler{
		server: server,
	}
}

func (h *ChatHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms := h.server.GetRoomList()
	sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"rooms":   rooms,
	})
}

func (h *ChatHandler) GetRoomUsers(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		roomID = "general"
	}

	room, exists := h.server.GetRoom(roomID)
	if !exists {
		sendJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Room not found",
		})
		return
	}

	users := room.GetUserList()
	sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"room":    room.Name,
		"users":   users,
		"count":   len(users),
	})
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ===== Main Demo =====

func main() {
	fmt.Println("===== WebSocket Chat Server =====\n")

	// Initialize chat server
	chatServer := NewChatServer()
	chatHandler := NewChatHandler(chatServer)

	// Simulate some users joining
	simulateChat(chatServer)

	// ===== HTTP Routes =====

	// API endpoints
	http.HandleFunc("/api/rooms", chatHandler.GetRooms)
	http.HandleFunc("/api/room/users", chatHandler.GetRoomUsers)

	// Documentation
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>WebSocket Chat</title>
    <style>
        body { font-family: Arial; max-width: 900px; margin: 50px auto; padding: 20px; }
        h1 { color: #00ADD8; }
        .info { background: #f0f0f0; padding: 15px; border-radius: 5px; margin: 10px 0; }
        .endpoint { background: #e8f4f8; padding: 15px; margin: 10px 0; border-radius: 5px; }
        code { background: #e8e8e8; padding: 2px 8px; border-radius: 3px; }
        .feature { margin: 5px 0; }
    </style>
</head>
<body>
    <h1>💬 WebSocket Real-time Chat System</h1>

    <div class="info">
        <h2>Features</h2>
        <div class="feature">✓ Multiple chat rooms</div>
        <div class="feature">✓ Real-time messaging</div>
        <div class="feature">✓ Private messages</div>
        <div class="feature">✓ User presence (join/leave notifications)</div>
        <div class="feature">✓ Typing indicators</div>
        <div class="feature">✓ Thread-safe concurrent access</div>
        <div class="feature">✓ Message broadcasting</div>
    </div>

    <h2>API Endpoints:</h2>

    <div class="endpoint">
        <strong>GET /api/rooms</strong><br>
        Get list of all chat rooms
    </div>

    <div class="endpoint">
        <strong>GET /api/room/users?room=general</strong><br>
        Get list of users in a room
    </div>

    <h2>WebSocket Connection (Simulated):</h2>
    <div class="info">
        <p><strong>Connect:</strong> <code>ws://localhost:8080/ws?room=general&username=alice</code></p>
        <p><strong>Message Format:</strong></p>
        <pre>{
  "type": "chat",
  "from": "alice",
  "content": "Hello, everyone!",
  "timestamp": "2024-01-01T12:00:00Z"
}</pre>
    </div>

    <h2>Message Types:</h2>
    <ul>
        <li><strong>chat</strong> - Regular chat message</li>
        <li><strong>private</strong> - Private message to specific user</li>
        <li><strong>system</strong> - System notification</li>
        <li><strong>join</strong> - User joined</li>
        <li><strong>leave</strong> - User left</li>
        <li><strong>typing</strong> - User is typing</li>
    </ul>

    <h2>Test Commands:</h2>
    <pre>
# Get all rooms
curl http://localhost:8080/api/rooms

# Get users in general room
curl http://localhost:8080/api/room/users?room=general

# Get users in tech room
curl http://localhost:8080/api/room/users?room=tech
    </pre>

    <div class="info">
        <h3>Architecture:</h3>
        <p><strong>ChatServer</strong> → Manages multiple rooms</p>
        <p><strong>Room</strong> → Manages clients and message broadcasting</p>
        <p><strong>Client</strong> → Represents connected user</p>
        <p><strong>Message</strong> → Chat message with type, content, timestamp</p>
    </div>
</body>
</html>
`
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
	})

	// Start server
	fmt.Println("Chat server running on http://localhost:8080")
	fmt.Println("\nActive rooms:")
	for id, room := range chatServer.Rooms {
		fmt.Printf("  - %s (%s)\n", room.Name, id)
	}
	fmt.Println("\nPress Ctrl+C to stop\n")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ===== Simulation =====

func simulateChat(server *ChatServer) {
	// Get general room
	generalRoom, _ := server.GetRoom("general")
	techRoom, _ := server.GetRoom("tech")

	// Simulate users joining
	go func() {
		time.Sleep(2 * time.Second)

		// Alice joins general
		alice := NewClient("alice-123", "Alice", generalRoom)
		generalRoom.Register <- alice
		go simulateUser(alice, []string{
			"Hello everyone!",
			"How's everyone doing?",
		})

		time.Sleep(1 * time.Second)

		// Bob joins general
		bob := NewClient("bob-456", "Bob", generalRoom)
		generalRoom.Register <- bob
		go simulateUser(bob, []string{
			"Hi Alice!",
			"I'm doing great, thanks!",
		})

		time.Sleep(1 * time.Second)

		// Carol joins tech
		carol := NewClient("carol-789", "Carol", techRoom)
		techRoom.Register <- carol
		go simulateUser(carol, []string{
			"Anyone here interested in Go?",
			"I'm learning about goroutines",
		})
	}()
}

func simulateUser(client *Client, messages []string) {
	conn := NewConnection(client)

	// Start read/write pumps
	go conn.ReadPump()
	go conn.WritePump()

	// Send messages
	for i, content := range messages {
		time.Sleep(time.Duration(2+i) * time.Second)

		msg := &Message{
			Type:    MessageTypeChat,
			From:    client.Username,
			Content: content,
		}

		conn.incoming <- msg
	}

	// Monitor outgoing messages
	go func() {
		for msg := range conn.outgoing {
			if msg.Type != MessageTypeSystem || msg.Content != "ping" {
				fmt.Printf("[%s] %s: %s\n", client.Room.Name, msg.From, msg.Content)
			}
		}
	}()
}

// วิธีรัน:
// go run 04_websocket_chat.go
//
// ทดสอบ API:
// curl http://localhost:8080/api/rooms
// curl http://localhost:8080/api/room/users?room=general
//
// Key Concepts:
//
// 1. WebSocket Communication:
//    - Bidirectional real-time communication
//    - Full-duplex connection
//    - Low latency
//
// 2. Chat Architecture:
//    - Server manages multiple rooms
//    - Each room has multiple clients
//    - Messages are broadcast to all clients in room
//    - Private messages to specific users
//
// 3. Concurrency Patterns:
//    - Each room runs in its own goroutine
//    - Read pump for incoming messages
//    - Write pump for outgoing messages
//    - Channels for message passing
//
// 4. Message Types:
//    - Chat: Regular messages
//    - Private: Direct messages
//    - System: Join/leave notifications
//    - Typing: Typing indicators
//
// Production Considerations:
// - Use gorilla/websocket for real WebSocket
// - Implement authentication
// - Add message persistence (database)
// - Implement rate limiting
// - Add reconnection logic
// - Scale with Redis pub/sub
// - Monitor connections
// - Handle disconnections gracefully
// - Implement heartbeat/ping
// - Add message history
// - Implement file sharing
// - Add emoji/reactions
// - Room permissions
// - Ban/kick functionality
//
// For Real WebSocket:
// go get github.com/gorilla/websocket
//
// upgrader := websocket.Upgrader{
//     CheckOrigin: func(r *http.Request) bool { return true },
// }
//
// conn, err := upgrader.Upgrade(w, r, nil)
// if err != nil {
//     log.Println(err)
//     return
// }
