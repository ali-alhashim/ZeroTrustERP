package core

import (
    "encoding/json"
    "fmt"
)

type Hub struct {
    clients    map[string]map[*Client]bool
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

        case client := <-h.register:
            if _, ok := h.clients[client.ID]; !ok {
                h.clients[client.ID] = make(map[*Client]bool)
            }
            h.clients[client.ID][client] = true

            fmt.Println("✅ User connected:", client.ID)
            h.NotifyOnlineUsers()

        case client := <-h.unregister:
            if conns, ok := h.clients[client.ID]; ok {
                if _, exists := conns[client]; exists {
                    delete(conns, client)
                    close(client.Send)
                }
                // If no more open connections for this user → offline
                if len(conns) == 0 {
                    delete(h.clients, client.ID)
                    fmt.Println("❌ User offline:", client.ID)
                }
            }
            h.NotifyOnlineUsers()

        case message := <-h.broadcast:
            // Send to ALL connections
            for _, conns := range h.clients {
                for client := range conns {
                    select {
                    case client.Send <- message:
                    default:
                        close(client.Send)
                        delete(conns, client)
                    }
                }
            }
        }
    }
}

func (h *Hub) NotifyOnlineUsers() {
    online := h.GetOnlineUserIDs()

    msg := map[string]interface{}{
        "type":  "ONLINE_USERS_UPDATE",
        "users": online,
    }

    data, _ := json.Marshal(msg)
    h.broadcast <- data
}

func (h *Hub) GetOnlineUserIDs() []string {
    ids := []string{}
    for id := range h.clients {
        ids = append(ids, id)
    }
    return ids
}