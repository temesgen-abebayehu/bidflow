package websocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"go.uber.org/zap"
)

// MockLogger for Hub tests
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, fields ...zap.Field)  {}
func (m *MockLogger) Info(msg string, fields ...zap.Field)   {}
func (m *MockLogger) Warn(msg string, fields ...zap.Field)   {}
func (m *MockLogger) Error(msg string, fields ...zap.Field)  {}
func (m *MockLogger) Fatal(msg string, fields ...zap.Field)  {}
func (m *MockLogger) With(fields ...zap.Field) logger.Logger { return m }
func (m *MockLogger) Sync() error                            { return nil }

func TestNewHub(t *testing.T) {
	log := &MockLogger{}
	hub := NewHub(log)

	assert.NotNil(t, hub)
	assert.NotNil(t, hub.clients)
	assert.NotNil(t, hub.userClients)
	assert.NotNil(t, hub.broadcast)
	assert.NotNil(t, hub.register)
	assert.NotNil(t, hub.unregister)
}

func TestHub_Register(t *testing.T) {
	log := &MockLogger{}
	hub := NewHub(log)

	// Start the hub in a goroutine
	go hub.Run()

	client := &Client{
		hub:    hub,
		userID: "user-1",
		send:   make(chan interface{}, 1), // Buffered to prevent blocking if logic changes
	}

	// Register
	hub.register <- client

	// Wait for registration to complete (simple sleep for unit test)
	// In a real scenario, we might check hub.clients size with a mutex, but sleep is enough here.
	// 10ms is usually plenty for a goroutine context switch.
	// A better way would be to use a WaitGroup or a callback if the Hub supported it.
	// Or check the map in a loop.
	for i := 0; i < 10; i++ {
		hub.mu.RLock()
		_, ok := hub.clients[client]
		hub.mu.RUnlock()
		if ok {
			break
		}
		// yield
	}

	msg := "hello"
	hub.BroadcastToUser("user-1", msg)
	
	received := <-client.send
	assert.Equal(t, msg, received)
}

func TestHub_Unregister(t *testing.T) {
	log := &MockLogger{}
	hub := NewHub(log)
	go hub.Run()

	client := &Client{
		hub:    hub,
		userID: "user-1",
		send:   make(chan interface{}, 1),
	}

	hub.register <- client
	
	// Wait for registration
	for i := 0; i < 10; i++ {
		hub.mu.RLock()
		_, ok := hub.clients[client]
		hub.mu.RUnlock()
		if ok {
			break
		}
	}

	// Unregister
	hub.unregister <- client
	
	// Verify channel is closed
	// We need to wait for the unregister to process
	// The hub closes the channel.
	// We can just read from it.
	_, ok := <-client.send
	assert.False(t, ok, "Channel should be closed after unregister")
}
