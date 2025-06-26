package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// LiveViewMessage represents a WebSocket message in LiveView style
type LiveViewMessage struct {
	Type      string                 `json:"type"`
	Topic     string                 `json:"topic,omitempty"`
	Event     string                 `json:"event,omitempty"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	Timestamp int64                  `json:"timestamp"`
}

// Client represents a WebSocket client connection
type Client struct {
	ID       string
	Conn     *websocket.Conn
	Hub      *Hub
	Send     chan LiveViewMessage
	Topics   map[string]bool
	UserID   string
	mu       sync.RWMutex
	IsActive bool
	LastPing time.Time
	Metadata map[string]interface{}
}

// Hub manages all WebSocket connections and topics
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Topic subscriptions: topic -> set of clients
	topics map[string]map[*Client]bool

	// Presence tracking: topic -> user data
	presence map[string]map[string]interface{}

	// Channel for registering clients
	register chan *Client

	// Channel for unregistering clients
	unregister chan *Client

	// Channel for broadcasting messages
	broadcast chan BroadcastMessage

	// Channel for topic-specific messages
	topicMessages chan TopicMessage

	// Mutex for thread safety
	mu sync.RWMutex

	// Context for cleanup
	ctx    context.Context
	cancel context.CancelFunc
}

// BroadcastMessage represents a message to broadcast to all clients
type BroadcastMessage struct {
	Message LiveViewMessage
	Exclude *Client // Optional client to exclude from broadcast
}

// TopicMessage represents a message for a specific topic
type TopicMessage struct {
	Topic   string
	Message LiveViewMessage
	Exclude *Client
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	ctx, cancel := context.WithCancel(context.Background())

	return &Hub{
		clients:       make(map[*Client]bool),
		topics:        make(map[string]map[*Client]bool),
		presence:      make(map[string]map[string]interface{}),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		broadcast:     make(chan BroadcastMessage),
		topicMessages: make(chan TopicMessage),
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	defer h.cancel()

	// Start heartbeat ticker
	heartbeat := time.NewTicker(30 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-h.ctx.Done():
			return

		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)

		case topicMsg := <-h.topicMessages:
			h.sendToTopic(topicMsg)

		case <-heartbeat.C:
			h.sendHeartbeat()
		}
	}
}

// registerClient adds a new client to the hub
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	client.IsActive = true
	client.LastPing = time.Now()

	log.Printf("Client %s connected. Total clients: %d", client.ID, len(h.clients))

	// Send welcome message
	welcome := LiveViewMessage{
		Type:      "welcome",
		Payload:   map[string]interface{}{"client_id": client.ID},
		Timestamp: time.Now().Unix(),
	}

	select {
	case client.Send <- welcome:
	default:
		close(client.Send)
		delete(h.clients, client)
	}
}

// unregisterClient removes a client from the hub
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		// Remove from all topics
		for topic := range client.Topics {
			h.unsubscribeFromTopic(client, topic)
		}

		// Remove from presence
		h.removeFromPresence(client)

		delete(h.clients, client)
		close(client.Send)

		log.Printf("Client %s disconnected. Total clients: %d", client.ID, len(h.clients))
	}
}

// broadcastMessage sends a message to all connected clients
func (h *Hub) broadcastMessage(broadcastMsg BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client == broadcastMsg.Exclude {
			continue
		}

		select {
		case client.Send <- broadcastMsg.Message:
		default:
			// Client's channel is full, close it
			close(client.Send)
			delete(h.clients, client)
		}
	}
}

// sendToTopic sends a message to all clients subscribed to a topic
func (h *Hub) sendToTopic(topicMsg TopicMessage) {
	h.mu.RLock()
	topicClients, exists := h.topics[topicMsg.Topic]
	h.mu.RUnlock()

	if !exists {
		return
	}

	for client := range topicClients {
		if client == topicMsg.Exclude {
			continue
		}

		select {
		case client.Send <- topicMsg.Message:
		default:
			// Client's channel is full, remove from topic
			h.unsubscribeFromTopic(client, topicMsg.Topic)
		}
	}
}

// subscribeToTopic adds a client to a topic
func (h *Hub) subscribeToTopic(client *Client, topic string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.topics[topic] == nil {
		h.topics[topic] = make(map[*Client]bool)
	}

	h.topics[topic][client] = true
	client.Topics[topic] = true

	log.Printf("Client %s subscribed to topic %s", client.ID, topic)

	// Send subscription confirmation
	confirm := LiveViewMessage{
		Type:      "subscribed",
		Topic:     topic,
		Payload:   map[string]interface{}{"status": "ok"},
		Timestamp: time.Now().Unix(),
	}

	select {
	case client.Send <- confirm:
	default:
		close(client.Send)
		delete(h.clients, client)
	}
}

// unsubscribeFromTopic removes a client from a topic
func (h *Hub) unsubscribeFromTopic(client *Client, topic string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if topicClients, exists := h.topics[topic]; exists {
		delete(topicClients, client)
		if len(topicClients) == 0 {
			delete(h.topics, topic)
		}
	}

	delete(client.Topics, topic)
}

// joinPresence adds a client to presence tracking for a topic
func (h *Hub) joinPresence(client *Client, topic string, data map[string]interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.presence[topic] == nil {
		h.presence[topic] = make(map[string]interface{})
	}

	userKey := fmt.Sprintf("%s:%s", client.UserID, client.ID)
	h.presence[topic][userKey] = data

	// Broadcast presence update
	go func() {
		presenceData := h.getPresenceData(topic)
		message := LiveViewMessage{
			Type:      "presence_update",
			Topic:     topic,
			Payload:   map[string]interface{}{"presence": presenceData},
			Timestamp: time.Now().Unix(),
		}

		h.topicMessages <- TopicMessage{
			Topic:   topic,
			Message: message,
		}
	}()
}

// removeFromPresence removes a client from all presence tracking
func (h *Hub) removeFromPresence(client *Client) {
	for topic := range h.presence {
		userKey := fmt.Sprintf("%s:%s", client.UserID, client.ID)
		delete(h.presence[topic], userKey)

		// Clean up empty presence topics
		if len(h.presence[topic]) == 0 {
			delete(h.presence, topic)
		}
	}
}

// getPresenceData returns presence data for a topic
func (h *Hub) getPresenceData(topic string) map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if presence, exists := h.presence[topic]; exists {
		return map[string]interface{}{
			"users": presence,
			"count": len(presence),
		}
	}

	return map[string]interface{}{
		"users": make(map[string]interface{}),
		"count": 0,
	}
}

// sendHeartbeat sends heartbeat to all clients
func (h *Hub) sendHeartbeat() {
	heartbeat := LiveViewMessage{
		Type:      "heartbeat",
		Timestamp: time.Now().Unix(),
	}

	h.broadcast <- BroadcastMessage{Message: heartbeat}
}

// BroadcastToTopic sends a message to all clients in a topic
func (h *Hub) BroadcastToTopic(topic string, event string, payload map[string]interface{}) {
	message := LiveViewMessage{
		Type:      "broadcast",
		Topic:     topic,
		Event:     event,
		Payload:   payload,
		Timestamp: time.Now().Unix(),
	}

	h.topicMessages <- TopicMessage{
		Topic:   topic,
		Message: message,
	}
}

// GetTopicCount returns the number of clients subscribed to a topic
func (h *Hub) GetTopicCount(topic string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, exists := h.topics[topic]; exists {
		return len(clients)
	}
	return 0
}

// GetOnlineCount returns the total number of connected clients
func (h *Hub) GetOnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// HandleMessage processes incoming messages from clients
func (h *Hub) HandleMessage(client *Client, messageType int, data []byte) {
	if messageType != websocket.TextMessage {
		return
	}

	var message LiveViewMessage
	if err := json.Unmarshal(data, &message); err != nil {
		log.Printf("Error parsing message from client %s: %v", client.ID, err)
		return
	}

	message.Timestamp = time.Now().Unix()

	switch message.Type {
	case "subscribe":
		if message.Topic != "" {
			h.subscribeToTopic(client, message.Topic)
		}

	case "unsubscribe":
		if message.Topic != "" {
			h.unsubscribeFromTopic(client, message.Topic)
		}

	case "join":
		if message.Topic != "" {
			h.subscribeToTopic(client, message.Topic)
			h.joinPresence(client, message.Topic, message.Payload)
		}

	case "leave":
		if message.Topic != "" {
			h.unsubscribeFromTopic(client, message.Topic)
		}

	case "push":
		// Handle client events (like voting, commenting, etc.)
		h.handleClientEvent(client, message)

	case "heartbeat_ack":
		client.LastPing = time.Now()

	default:
		log.Printf("Unknown message type from client %s: %s", client.ID, message.Type)
	}
}

// handleClientEvent processes client-initiated events
func (h *Hub) handleClientEvent(client *Client, message LiveViewMessage) {
	// This is where you'd handle business logic events
	// For example: voting, creating nodes, etc.

	switch message.Event {
	case "vote_node":
		// Handle node voting
		log.Printf("Client %s voted on node: %v", client.ID, message.Payload)

		// Broadcast the vote to other clients in the topic
		broadcastMsg := LiveViewMessage{
			Type:      "update",
			Topic:     message.Topic,
			Event:     "node_voted",
			Payload:   message.Payload,
			Timestamp: time.Now().Unix(),
		}

		h.topicMessages <- TopicMessage{
			Topic:   message.Topic,
			Message: broadcastMsg,
			Exclude: client, // Don't send back to the sender
		}

	case "create_node":
		// Handle node creation
		log.Printf("Client %s created node: %v", client.ID, message.Payload)

		// Broadcast to activity feed
		activityMsg := LiveViewMessage{
			Type:      "broadcast",
			Topic:     "activity_feed",
			Event:     "node_created",
			Payload:   message.Payload,
			Timestamp: time.Now().Unix(),
		}

		h.topicMessages <- TopicMessage{
			Topic:   "activity_feed",
			Message: activityMsg,
		}

	default:
		log.Printf("Unknown client event: %s", message.Event)
	}
}

// NewClient creates a new WebSocket client
func NewClient(id string, conn *websocket.Conn, hub *Hub, userID string) *Client {
	return &Client{
		ID:       id,
		Conn:     conn,
		Hub:      hub,
		Send:     make(chan LiveViewMessage, 256),
		Topics:   make(map[string]bool),
		UserID:   userID,
		IsActive: false,
		LastPing: time.Now(),
		Metadata: make(map[string]interface{}),
	}
}

// WritePump handles writing messages to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("Error writing to client %s: %v", c.ID, err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ReadPump handles reading messages from the WebSocket connection
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for client %s: %v", c.ID, err)
			}
			break
		}

		c.Hub.HandleMessage(c, messageType, message)
	}
}

// Upgrader configures the WebSocket upgrader
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
}
