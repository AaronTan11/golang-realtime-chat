package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)

// MessageType represents the type of message being sent
type MessageType string

const (
	// MessageTypeJoin indicates a user joined the chat
	MessageTypeJoin MessageType = "join"
	// MessageTypeLeave indicates a user left the chat
	MessageTypeLeave MessageType = "leave"
	// MessageTypeChat indicates a regular chat message
	MessageTypeChat MessageType = "chat"
	// MessageTypeError indicates an error occurred
	MessageTypeError MessageType = "error"
	// MessageTypeWelcome indicates a server welcome with assigned user ID
	MessageTypeWelcome MessageType = "welcome"
)

// Message represents a chat message
type Message struct {
	Type      MessageType `json:"type"`
	Username  string      `json:"username"`
	UserID    string      `json:"userId,omitempty"`
	Content   string      `json:"content,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Client represents a connected WebSocket client
type Client struct {
	ID       string
	Username string
	Send     chan Message
	Hub      *Hub
}

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered clients
	Clients map[*Client]bool
	// Inbound messages from clients
	Broadcast chan Message
	// Register requests from clients
	Register chan *Client
	// Unregister requests from clients
	Unregister chan *Client
	// nextID is an atomically incremented counter for assigning client IDs
	nextID uint64
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// NextIDString increments and returns the next client ID as a string
func (h *Hub) NextIDString() string {
	id := atomic.AddUint64(&h.nextID, 1)
	return strconv.FormatUint(id, 10)
}

// Run starts the hub and handles client registration/unregistration and message broadcasting
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			// Send a welcome message to the newly connected client with their assigned ID
			select {
			case client.Send <- Message{
				Type:      MessageTypeWelcome,
				Username:  client.Username,
				UserID:    client.ID,
				Content:   "Welcome to the chat!",
				Timestamp: time.Now(),
			}:
			default:
				// if client's send buffer is full, drop welcome
			}
			// Notify all clients that someone joined
			joinMessage := Message{
				Type:      MessageTypeJoin,
				Username:  client.Username,
				UserID:    client.ID,
				Content:   fmt.Sprintf("%s (%s) joined the chat", client.Username, client.ID),
				Timestamp: time.Now(),
			}
			h.broadcastMessage(joinMessage)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				// Notify all clients that someone left
				leaveMessage := Message{
					Type:      MessageTypeLeave,
					Username:  client.Username,
					UserID:    client.ID,
					Content:   fmt.Sprintf("%s (%s) left the chat", client.Username, client.ID),
					Timestamp: time.Now(),
				}
				h.broadcastMessage(leaveMessage)
			}

		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

// broadcastMessage sends a message to all connected clients
func (h *Hub) broadcastMessage(message Message) {
	for client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.Clients, client)
		}
	}
}

// ToJSON converts a Message to JSON string
func (m Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// ParseMessage parses JSON data into a Message
func ParseMessage(data []byte) (Message, error) {
	var message Message
	err := json.Unmarshal(data, &message)
	return message, err
}
