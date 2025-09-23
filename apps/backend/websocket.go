package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// upgrader is used to upgrade HTTP connections to WebSocket connections
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow connections from any origin (for development)
	// In production, you should check the origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// serveWebSocket handles WebSocket connections
func serveWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Get username from query parameter
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "Anonymous"
	}

	// Create a new client with simple incrementing numeric ID
	client := &Client{
		ID:       hub.NextIDString(),
		Username: username,
		Send:     make(chan Message, 256),
		Hub:      hub,
	}

	// Register the client with the hub
	client.Hub.Register <- client

	// Start goroutines for reading and writing
	go client.writePump(conn)
	go client.readPump(conn)
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump(conn *websocket.Conn) {
	defer func() {
		c.Hub.Unregister <- c
		conn.Close()
	}()

	// Set read deadline and pong handler
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		// Read message from WebSocket
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse the message
		message, err := ParseMessage(messageBytes)
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Set the username and timestamp
		message.Username = c.Username
		message.UserID = c.ID
		message.Timestamp = time.Now()
		message.Type = MessageTypeChat

		// Send the message to the hub for broadcasting
		c.Hub.Broadcast <- message
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump(conn *websocket.Conn) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The hub closed the channel
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Convert message to JSON
			messageBytes, err := message.ToJSON()
			if err != nil {
				log.Printf("Error converting message to JSON: %v", err)
				continue
			}

			// Write the message to the WebSocket connection
			if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Removed old timestamp/random ID generator in favor of simple incrementing numeric IDs
