# Go Realtime Chat Backend

A simple real-time chat application built with Go and WebSockets, perfect for learning Go fundamentals and real-time communication patterns.

## Features

- **Real-time messaging** using WebSockets
- **In-memory storage** (no database required)
- **Multiple users** can join and chat simultaneously
- **Join/Leave notifications** when users connect or disconnect
- **REST API endpoints** for monitoring and statistics
- **CORS support** for web frontend integration

## Architecture

### Core Components

1. **Hub** (`types.go`): Central coordinator that manages all WebSocket connections
2. **Client** (`types.go`): Represents a connected user with their WebSocket connection
3. **Message** (`types.go`): Data structure for chat messages with different types
4. **WebSocket Handler** (`websocket.go`): Handles WebSocket upgrade and message processing
5. **HTTP Server** (`main.go`): Serves REST API endpoints and WebSocket connections

### Key Go Concepts Demonstrated

- **Goroutines**: Concurrent handling of multiple clients
- **Channels**: Communication between goroutines (hub, clients, messages)
- **Interfaces**: Clean separation of concerns
- **Structs and Methods**: Object-oriented patterns in Go
- **JSON Marshaling/Unmarshaling**: Data serialization
- **HTTP Handlers**: Web server patterns
- **WebSocket Protocol**: Real-time bidirectional communication

## Getting Started

### Prerequisites

- Go 1.25.0 or later
- A web browser for testing

### Installation

1. Navigate to the backend directory:
   ```bash
   cd apps/backend
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run .
   ```

The server will start on `http://localhost:8080`

### Testing the Backend

#### Option 1: Use the Test HTML Client

1. Open `test-client.html` in your web browser
2. Enter a username and click "Connect"
3. Open multiple browser tabs/windows to simulate multiple users
4. Send messages and see real-time updates

#### Option 2: Use WebSocket Tools

- **WebSocket URL**: `ws://localhost:8080/ws?username=YourName`
- **Message Format**: JSON with `type`, `username`, `content`, and `timestamp`

Example message:
```json
{
  "type": "chat",
  "username": "John",
  "content": "Hello everyone!",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## API Endpoints

### HTTP Endpoints

- `GET /` - Server status and connection info
- `GET /healthz` - Health check endpoint
- `GET /api/users` - List of currently connected users
- `GET /api/stats` - Server statistics

### WebSocket Endpoint

- `WS /ws?username=YourName` - WebSocket connection for real-time chat

## Message Types

- **`join`**: User joined the chat
- **`leave`**: User left the chat
- **`chat`**: Regular chat message
- **`error`**: Error message

## Learning Path

This backend is designed to teach Go concepts progressively:

1. **Basic HTTP Server**: Understanding Go's HTTP package
2. **WebSocket Integration**: Learning about protocol upgrades
3. **Concurrency**: Using goroutines and channels for real-time communication
4. **Data Structures**: Designing efficient message and client structures
5. **Error Handling**: Proper error management in Go
6. **JSON Processing**: Data serialization and deserialization

## Code Structure

```
apps/backend/
├── main.go          # HTTP server and main application
├── types.go         # Data structures and Hub implementation
├── websocket.go     # WebSocket connection handling
├── go.mod           # Go module dependencies
├── test-client.html # HTML test client
└── README.md        # This file
```

## Key Learning Points

### Goroutines and Channels
- The `Hub.Run()` method runs in a goroutine
- Channels (`Register`, `Unregister`, `Broadcast`) coordinate between goroutines
- Each client has its own goroutines for reading and writing

### WebSocket Lifecycle
1. HTTP connection is upgraded to WebSocket
2. Client is registered with the Hub
3. Read/Write goroutines are started
4. Messages are broadcast to all connected clients
5. Cleanup happens when client disconnects

### Error Handling
- WebSocket errors are logged but don't crash the server
- JSON parsing errors are handled gracefully
- Connection timeouts prevent resource leaks

## Extending the Application

Some ideas for further learning:

1. **Add user authentication**
2. **Implement multiple chat rooms**
3. **Add message persistence**
4. **Implement private messaging**
5. **Add file sharing capabilities**
6. **Scale with Redis pub/sub**

## Troubleshooting

### Common Issues

1. **Port already in use**: Change the port in `main.go` (line 101)
2. **CORS errors**: The server includes CORS headers, but check browser console
3. **WebSocket connection fails**: Ensure the server is running and accessible

### Debugging

- Check server logs for error messages
- Use browser developer tools to inspect WebSocket connections
- Test API endpoints with curl or Postman

## Contributing

This is a learning project! Feel free to:
- Add new features
- Improve error handling
- Add more comprehensive tests
- Optimize performance
- Enhance the documentation
