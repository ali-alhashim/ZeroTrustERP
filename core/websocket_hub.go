package core

import "fmt"

type Hub struct {
    clients    map[string]map[*Client]bool // userID → set of connections
    register   chan *Client
    unregister chan *Client
    broadcast  chan []byte
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[string]map[*Client]bool),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcast:  make(chan []byte),
    }
}

func (h *Hub) Run() {
    for {
        select {

        // ✅ New client connection
        case client := <-h.register:
            if _, exists := h.clients[client.ID]; !exists {
                h.clients[client.ID] = make(map[*Client]bool)
            }
            h.clients[client.ID][client] = true

            fmt.Println("User connected:", client.ID)
            h.NotifyOnlineUsers()

        // ✅ Client disconnected
        case client := <-h.unregister:
            if conns, exists := h.clients[client.ID]; exists {
                // Remove client connection
                if _, ok := conns[client]; ok {
                    delete(conns, client)
                }

                // ✅ If user has zero websocket sessions → user is offline
                if len(conns) == 0 {
                    delete(h.clients, client.ID)
                    fmt.Println("User offline:", client.ID)
                }
            }

            // Prevent panic: close only user's channel
            close(client.Send)

            h.NotifyOnlineUsers()

        // ✅ Broadcast a message to ALL websocket connections
        case message := <-h.broadcast:
            for _, conns := range h.clients {
                for client := range conns {
                    select {
                    case client.Send <- message:
                    default:
                        // Drop stuck client
                        close(client.Send)
                        delete(conns, client)
                    }
                }
            }
        }
    }
}

func (h *Hub) NotifyOnlineUsers() {
    fmt.Println("Notifying clients about online users update")
    h.broadcast <- []byte("ONLINE_USERS_UPDATE")
}

// ✅ Helper: return all online user IDs
func (h *Hub) GetOnlineUserIDs() []string {
    ids := []string{}
    for userID := range h.clients {
        ids = append(ids, userID)
    }
    return ids
}