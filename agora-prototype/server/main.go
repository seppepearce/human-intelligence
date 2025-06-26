package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Message types for WebSocket communication
const (
	MessageTypePresence                   = "presence"
	MessageTypeSpaceUpdate                = "space_update"
	MessageTypeJoinSpace                  = "join_space"
	MessageTypeLeaveSpace                 = "leave_space"
	MessageTypeFriendUpdate               = "friend_update"
	MessageTypeMigrationSuggestion        = "migration_suggestion"
	MessageTypeCrossPollinationSuggestion = "cross_pollination_suggestion"
)

// WebSocket message structure
type Message struct {
	Type    string      `json:"type"`
	UserID  string      `json:"userId,omitempty"`
	SpaceID string      `json:"spaceId,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// User represents a connected user
type User struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Avatar       string          `json:"avatar"`
	CurrentSpace string          `json:"currentSpace,omitempty"`
	Connection   *websocket.Conn `json:"-"`
	LastSeen     time.Time       `json:"lastSeen"`
	Interests    []string        `json:"interests"`
}

// Space represents a discussion space with natural constraints
type Space struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Participants map[string]*User `json:"participants"`
	Capacity     int              `json:"capacity"`
	Topics       []string         `json:"topics"`
	Activity     string           `json:"activity"`
	LastActive   time.Time        `json:"lastActive"`
	AvgDepth     float64          `json:"avgDepth"`
	mutex        sync.RWMutex     `json:"-"`
}

// AgoraServer manages the digital agora
type AgoraServer struct {
	spaces   map[string]*Space
	users    map[string]*User
	friends  map[string][]string // userId -> []friendIds
	upgrader websocket.Upgrader
	mutex    sync.RWMutex
}

// NewAgoraServer creates a new agora server
func NewAgoraServer() *AgoraServer {
	server := &AgoraServer{
		spaces:  make(map[string]*Space),
		users:   make(map[string]*User),
		friends: make(map[string][]string),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}

	// Initialize demo spaces
	server.initializeDemoSpaces()

	// Start background processes
	go server.maintainSpaceActivity()
	go server.suggestMigrations()

	return server
}

// Initialize demo spaces for the agora
func (as *AgoraServer) initializeDemoSpaces() {
	demoSpaces := []*Space{
		{
			ID:           "philosophy",
			Name:         "Philosophy of Mind",
			Description:  "Exploring consciousness, qualia, and the hard problem",
			Participants: make(map[string]*User),
			Capacity:     12,
			Topics:       []string{"Consciousness", "Free Will", "Qualia"},
			Activity:     "high",
			LastActive:   time.Now(),
			AvgDepth:     8.5,
		},
		{
			ID:           "ai-ethics",
			Name:         "AI Ethics",
			Description:  "Moral implications of artificial intelligence",
			Participants: make(map[string]*User),
			Capacity:     12,
			Topics:       []string{"Alignment", "Bias", "Privacy"},
			Activity:     "very-high",
			LastActive:   time.Now(),
			AvgDepth:     7.2,
		},
		{
			ID:           "cognitive-science",
			Name:         "Cognitive Science",
			Description:  "How minds work across species and machines",
			Participants: make(map[string]*User),
			Capacity:     12,
			Topics:       []string{"Neuroscience", "Psychology", "AI"},
			Activity:     "medium",
			LastActive:   time.Now(),
			AvgDepth:     6.8,
		},
		{
			ID:           "ancient-wisdom",
			Name:         "Ancient Wisdom",
			Description:  "Timeless insights from classical philosophers",
			Participants: make(map[string]*User),
			Capacity:     12,
			Topics:       []string{"Stoicism", "Virtue Ethics", "Dialectics"},
			Activity:     "medium",
			LastActive:   time.Now(),
			AvgDepth:     9.1,
		},
		{
			ID:           "quantum-physics",
			Name:         "Quantum Physics",
			Description:  "The strange world of quantum mechanics",
			Participants: make(map[string]*User),
			Capacity:     12,
			Topics:       []string{"Superposition", "Entanglement", "Measurement"},
			Activity:     "low",
			LastActive:   time.Now().Add(-5 * time.Minute),
			AvgDepth:     5.4,
		},
	}

	for _, space := range demoSpaces {
		as.spaces[space.ID] = space
	}

	// Initialize demo friend networks
	as.friends["socrates"] = []string{"plato", "aristotle"}
	as.friends["plato"] = []string{"socrates", "aristotle"}
	as.friends["aristotle"] = []string{"socrates", "plato"}
}

// WebSocket handler
func (as *AgoraServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := as.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Handle WebSocket messages
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		as.handleMessage(conn, &msg)
	}
}

// Handle incoming WebSocket messages
func (as *AgoraServer) handleMessage(conn *websocket.Conn, msg *Message) {
	switch msg.Type {
	case MessageTypePresence:
		as.handlePresenceUpdate(conn, msg)
	case MessageTypeJoinSpace:
		as.handleJoinSpace(conn, msg)
	case MessageTypeLeaveSpace:
		as.handleLeaveSpace(conn, msg)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// Handle user presence updates
func (as *AgoraServer) handlePresenceUpdate(conn *websocket.Conn, msg *Message) {
	userDataBytes, _ := json.Marshal(msg.Data)
	var userData User
	json.Unmarshal(userDataBytes, &userData)

	as.mutex.Lock()
	defer as.mutex.Unlock()

	// Update or create user
	user := &User{
		ID:         userData.ID,
		Name:       userData.Name,
		Avatar:     userData.Avatar,
		Connection: conn,
		LastSeen:   time.Now(),
		Interests:  userData.Interests,
	}

	as.users[user.ID] = user

	// Broadcast presence update to friends
	as.broadcastFriendUpdate(user.ID)

	log.Printf("User %s (%s) connected", user.Name, user.ID)
}

// Handle user joining a space
func (as *AgoraServer) handleJoinSpace(conn *websocket.Conn, msg *Message) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	space, exists := as.spaces[msg.SpaceID]
	if !exists {
		return
	}

	user, exists := as.users[msg.UserID]
	if !exists {
		return
	}

	space.mutex.Lock()
	defer space.mutex.Unlock()

	// Check capacity constraints
	if len(space.Participants) >= space.Capacity {
		// Suggest alternative spaces
		as.suggestAlternativeSpaces(user, space)
		return
	}

	// Remove user from current space if any
	if user.CurrentSpace != "" {
		as.removeUserFromSpace(user.ID, user.CurrentSpace)
	}

	// Add user to new space
	space.Participants[user.ID] = user
	user.CurrentSpace = msg.SpaceID
	space.LastActive = time.Now()

	// Update activity level based on participants
	as.updateSpaceActivity(space)

	// Broadcast space update
	as.broadcastSpaceUpdate(space)

	// Notify friends about the user's new location
	as.notifyFriendsOfLocation(user)

	log.Printf("User %s joined space %s (%d/%d)", user.Name, space.Name, len(space.Participants), space.Capacity)
}

// Handle user leaving a space
func (as *AgoraServer) handleLeaveSpace(conn *websocket.Conn, msg *Message) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	as.removeUserFromSpace(msg.UserID, msg.SpaceID)
}

// Remove user from a space
func (as *AgoraServer) removeUserFromSpace(userID, spaceID string) {
	space, exists := as.spaces[spaceID]
	if !exists {
		return
	}

	user, exists := as.users[userID]
	if !exists {
		return
	}

	space.mutex.Lock()
	defer space.mutex.Unlock()

	delete(space.Participants, userID)
	user.CurrentSpace = ""

	// Update activity level
	as.updateSpaceActivity(space)

	// Broadcast space update
	as.broadcastSpaceUpdate(space)

	log.Printf("User %s left space %s (%d/%d)", user.Name, space.Name, len(space.Participants), space.Capacity)
}

// Update space activity level based on participants
func (as *AgoraServer) updateSpaceActivity(space *Space) {
	participantCount := len(space.Participants)
	capacity := space.Capacity

	percentage := float64(participantCount) / float64(capacity)

	switch {
	case percentage >= 0.9:
		space.Activity = "very-high"
	case percentage >= 0.7:
		space.Activity = "high"
	case percentage >= 0.4:
		space.Activity = "medium"
	default:
		space.Activity = "low"
	}
}

// Broadcast space update to all users
func (as *AgoraServer) broadcastSpaceUpdate(space *Space) {
	message := Message{
		Type:    MessageTypeSpaceUpdate,
		SpaceID: space.ID,
		Data:    space,
	}

	for _, user := range as.users {
		if user.Connection != nil {
			user.Connection.WriteJSON(message)
		}
	}
}

// Broadcast friend update
func (as *AgoraServer) broadcastFriendUpdate(userID string) {
	friends, exists := as.friends[userID]
	if !exists {
		return
	}

	user := as.users[userID]
	message := Message{
		Type:   MessageTypeFriendUpdate,
		UserID: userID,
		Data: map[string]interface{}{
			"user":         user,
			"currentSpace": user.CurrentSpace,
		},
	}

	for _, friendID := range friends {
		if friend, exists := as.users[friendID]; exists && friend.Connection != nil {
			friend.Connection.WriteJSON(message)
		}
	}
}

// Notify friends when a user changes location
func (as *AgoraServer) notifyFriendsOfLocation(user *User) {
	friends, exists := as.friends[user.ID]
	if !exists {
		return
	}

	space := as.spaces[user.CurrentSpace]
	message := Message{
		Type:   MessageTypeFriendUpdate,
		UserID: user.ID,
		Data: map[string]interface{}{
			"friend":  user,
			"space":   space,
			"action":  "joined",
			"message": fmt.Sprintf("%s is now exploring %s", user.Name, space.Name),
		},
	}

	for _, friendID := range friends {
		if friend, exists := as.users[friendID]; exists && friend.Connection != nil {
			friend.Connection.WriteJSON(message)
		}
	}
}

// Suggest alternative spaces when one is full
func (as *AgoraServer) suggestAlternativeSpaces(user *User, fullSpace *Space) {
	alternatives := as.findAlternativeSpaces(user, fullSpace)

	message := Message{
		Type:    MessageTypeMigrationSuggestion,
		UserID:  user.ID,
		SpaceID: fullSpace.ID,
		Data: map[string]interface{}{
			"fullSpace":    fullSpace,
			"alternatives": alternatives,
			"message":      fmt.Sprintf("%s is full. Here are some related spaces you might enjoy:", fullSpace.Name),
		},
	}

	if user.Connection != nil {
		user.Connection.WriteJSON(message)
	}
}

// Find alternative spaces based on topics and capacity
func (as *AgoraServer) findAlternativeSpaces(user *User, fullSpace *Space) []*Space {
	var alternatives []*Space

	for _, space := range as.spaces {
		if space.ID == fullSpace.ID {
			continue
		}

		// Check if space has capacity
		space.mutex.RLock()
		hasCapacity := len(space.Participants) < space.Capacity
		space.mutex.RUnlock()

		if !hasCapacity {
			continue
		}

		// Check for topic overlap or user interest alignment
		hasRelatedTopics := as.hasTopicOverlap(fullSpace.Topics, space.Topics)
		alignsWithInterests := as.alignsWithUserInterests(user.Interests, space.Topics)

		if hasRelatedTopics || alignsWithInterests {
			alternatives = append(alternatives, space)
		}
	}

	return alternatives
}

// Check if two topic lists have overlap
func (as *AgoraServer) hasTopicOverlap(topics1, topics2 []string) bool {
	for _, topic1 := range topics1 {
		for _, topic2 := range topics2 {
			if topic1 == topic2 {
				return true
			}
		}
	}
	return false
}

// Check if space topics align with user interests
func (as *AgoraServer) alignsWithUserInterests(interests, topics []string) bool {
	for _, interest := range interests {
		for _, topic := range topics {
			if interest == topic {
				return true
			}
		}
	}
	return false
}

// Background process to maintain space activity simulation
func (as *AgoraServer) maintainSpaceActivity() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		as.mutex.RLock()
		for _, space := range as.spaces {
			// Simulate small changes in participation
			change := rand.Intn(3) - 1 // -1, 0, or 1

			space.mutex.Lock()

			// Create or remove phantom participants for demo
			currentCount := len(space.Participants)

			// Don't go below actual connected users or above capacity
			realUsers := 0
			for _, participant := range space.Participants {
				if participant.Connection != nil {
					realUsers++
				}
			}

			newCount := currentCount + change
			if newCount < realUsers {
				newCount = realUsers
			}
			if newCount > space.Capacity {
				newCount = space.Capacity
			}

			// Adjust phantom participants
			if newCount > currentCount {
				// Add phantom participant
				phantomID := fmt.Sprintf("phantom_%d", rand.Intn(1000))
				space.Participants[phantomID] = &User{
					ID:     phantomID,
					Name:   "Anonymous",
					Avatar: "👤",
				}
			} else if newCount < currentCount && currentCount > realUsers {
				// Remove a phantom participant
				for id, participant := range space.Participants {
					if participant.Connection == nil {
						delete(space.Participants, id)
						break
					}
				}
			}

			as.updateSpaceActivity(space)
			space.mutex.Unlock()
		}
		as.mutex.RUnlock()
	}
}

// Background process to suggest cross-pollination opportunities
func (as *AgoraServer) suggestMigrations() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		as.mutex.RLock()
		for _, user := range as.users {
			if user.Connection == nil || user.CurrentSpace == "" {
				continue
			}

			currentSpace := as.spaces[user.CurrentSpace]
			if currentSpace == nil {
				continue
			}

			// Check if current space is getting crowded
			currentSpace.mutex.RLock()
			isCrowded := float64(len(currentSpace.Participants))/float64(currentSpace.Capacity) > 0.8
			currentSpace.mutex.RUnlock()

			if isCrowded {
				// Find cross-pollination opportunities
				crossPollinations := as.findCrossPollinationOpportunities(user)
				if len(crossPollinations) > 0 {
					message := Message{
						Type:   MessageTypeCrossPollinationSuggestion,
						UserID: user.ID,
						Data: map[string]interface{}{
							"currentSpace":    currentSpace,
							"opportunities":   crossPollinations,
							"message":         fmt.Sprintf("%s is getting crowded. Perfect time to explore new intellectual territory!", currentSpace.Name),
							"migrationReason": "crowded_space",
						},
					}

					user.Connection.WriteJSON(message)
				}
			}
		}
		as.mutex.RUnlock()
	}
}

// Find cross-pollination opportunities for a user
func (as *AgoraServer) findCrossPollinationOpportunities(user *User) []*Space {
	var opportunities []*Space

	for _, space := range as.spaces {
		if space.ID == user.CurrentSpace {
			continue
		}

		space.mutex.RLock()
		hasCapacity := len(space.Participants) < int(float64(space.Capacity)*0.7) // Not too crowded
		space.mutex.RUnlock()

		if !hasCapacity {
			continue
		}

		// Check if this would be cross-pollination (different from user's interests)
		isNewTerritory := !as.alignsWithUserInterests(user.Interests, space.Topics)

		// But still somewhat related (not completely alien)
		currentSpace := as.spaces[user.CurrentSpace]
		hasIndirectConnection := as.hasTopicOverlap(currentSpace.Topics, space.Topics) ||
			as.hasCommonFriendsInSpace(user.ID, space)

		if isNewTerritory && hasIndirectConnection {
			opportunities = append(opportunities, space)
		}
	}

	return opportunities
}

// Check if user has friends in a space
func (as *AgoraServer) hasCommonFriendsInSpace(userID string, space *Space) bool {
	friends, exists := as.friends[userID]
	if !exists {
		return false
	}

	space.mutex.RLock()
	defer space.mutex.RUnlock()

	for _, friendID := range friends {
		if _, exists := space.Participants[friendID]; exists {
			return true
		}
	}

	return false
}

// HTTP handler for space data
func (as *AgoraServer) handleSpaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	as.mutex.RLock()
	defer as.mutex.RUnlock()

	// Convert spaces to JSON-friendly format
	spaces := make([]map[string]interface{}, 0, len(as.spaces))
	for _, space := range as.spaces {
		space.mutex.RLock()
		spaceData := map[string]interface{}{
			"id":           space.ID,
			"name":         space.Name,
			"description":  space.Description,
			"participants": len(space.Participants),
			"capacity":     space.Capacity,
			"topics":       space.Topics,
			"activity":     space.Activity,
			"lastActive":   space.LastActive,
			"avgDepth":     space.AvgDepth,
		}
		space.mutex.RUnlock()
		spaces = append(spaces, spaceData)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"spaces": spaces,
	})
}

func main() {
	server := NewAgoraServer()

	r := mux.NewRouter()

	// WebSocket endpoint
	r.HandleFunc("/ws", server.handleWebSocket)

	// REST API endpoints
	r.HandleFunc("/api/spaces", server.handleSpaces).Methods("GET")

	// Serve static files (for development)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	fmt.Println("🏛️  Digital Agora Server starting on :8083")
	fmt.Println("WebSocket endpoint: ws://localhost:8083/ws")
	fmt.Println("API endpoint: http://localhost:8083/api/spaces")

	log.Fatal(http.ListenAndServe(":8083", r))
}
