package core

import "fmt"

type Hub struct {
    clients    map[string]*Client
    register   chan *Client
    unregister chan *Client
    broadcast  chan []byte
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[string]*Client),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcast:  make(chan []byte),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client.ID] = client
            h.NotifyOnlineUsers()

        case client := <-h.unregister:
            delete(h.clients, client.ID)
            close(client.Send)
            h.NotifyOnlineUsers()

        case message := <-h.broadcast:
            for _, client := range h.clients {
                client.Send <- message
            }
        }
    }
}

func (h *Hub) NotifyOnlineUsers() {
    fmt.Printf("Notifying clients about online users update\n")
    online := []byte("ONLINE_USERS_UPDATE")
    h.broadcast <- online
}
