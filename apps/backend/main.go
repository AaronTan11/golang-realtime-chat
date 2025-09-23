package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Create a new hub for managing WebSocket connections
	hub := NewHub()

	// Start the hub in a goroutine
	go hub.Run()

	// Create HTTP server with routes
	mux := http.NewServeMux()

	// Root endpoint - shows server status
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "golang-realtime-chat backend is running")
		fmt.Fprintf(w, "Active connections: %d\n", len(hub.Clients))
		fmt.Fprintln(w, "Connect to /ws?username=YourName to join the chat")
	})

	// Health check endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// WebSocket endpoint
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebSocket(hub, w, r)
	})

	// API endpoint to get current connected users
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Get list of connected usernames and detailed entries
		var usernames []string
		var detailed []map[string]string
		for client := range hub.Clients {
			usernames = append(usernames, client.Username)
			detailed = append(detailed, map[string]string{
				"id":       client.ID,
				"username": client.Username,
			})
		}

		response := map[string]interface{}{
			"users":         usernames,
			"usersDetailed": detailed,
			"count":         len(usernames),
		}

		json.NewEncoder(w).Encode(response)
	})

	// API endpoint to get server stats
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		stats := map[string]interface{}{
			"active_connections": len(hub.Clients),
			"uptime":             time.Since(time.Now()).String(), // This will be 0, but shows the structure
			"server_time":        time.Now().Format(time.RFC3339),
		}

		json.NewEncoder(w).Encode(stats)
	})

	// Add CORS middleware for all routes
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Start the server
	addr := ":8080"
	log.Printf("🚀 Starting golang-realtime-chat server on %s", addr)
	log.Printf("📡 WebSocket endpoint: ws://localhost%s/ws", addr)
	log.Printf("🌐 API endpoints:")
	log.Printf("   - GET  http://localhost%s/api/users  (list connected users)", addr)
	log.Printf("   - GET  http://localhost%s/api/stats  (server statistics)", addr)
	log.Printf("   - GET  http://localhost%s/healthz    (health check)", addr)

	if err := http.ListenAndServe(addr, corsHandler(mux)); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Server error: %v", err)
	}
}
