package websocket

import (
	"sync"

	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"github.com/temesgen-abebayehu/bidflow/backend/services/notification/internal/domain"
	"go.uber.org/zap"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Map userID to clients (one user can have multiple connections)
	userClients map[string]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	logger logger.Logger
	mu     sync.RWMutex
}

func NewHub(log logger.Logger) *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		userClients: make(map[string]map[*Client]bool),
		logger:      log,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if _, ok := h.userClients[client.userID]; !ok {
				h.userClients[client.userID] = make(map[*Client]bool)
			}
			h.userClients[client.userID][client] = true
			h.mu.Unlock()
			h.logger.Info("Client registered", zap.String("user_id", client.userID))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				if userMap, ok := h.userClients[client.userID]; ok {
					delete(userMap, client)
					if len(userMap) == 0 {
						delete(h.userClients, client.userID)
					}
				}
			}
			h.mu.Unlock()
			h.logger.Info("Client unregistered", zap.String("user_id", client.userID))

		case message := <-h.broadcast:
			// Broadcast to everyone (optional, maybe for system alerts)
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastToUser(userID string, message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.userClients[userID]
	if !ok {
		return // User not connected
	}

	for client := range clients {
		select {
		case client.send <- message:
		default:
			// Channel full or closed, we might want to clean up here or let the main loop handle it
			// For now, just skip
		}
	}
}

func (h *Hub) Register(client domain.Client) {
	// This method is needed to satisfy the interface but the channel is used internally
	// We can cast or change the interface.
	// For simplicity in this implementation, we use the channel directly in the handler.
	// But to satisfy domain.Hub interface:
	if c, ok := client.(*Client); ok {
		h.register <- c
	}
}

func (h *Hub) Unregister(client domain.Client) {
	if c, ok := client.(*Client); ok {
		h.unregister <- c
	}
}
