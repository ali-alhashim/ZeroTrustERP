package core

import (
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // You can add ZeroTrust rules later
    },
}

type Client struct {
    ID     string
    Conn   *websocket.Conn
    Send   chan []byte
    Hub    *Hub
}

func (c *Client) ReadPump() {
    defer func() {
        c.Hub.unregister <- c
        c.Conn.Close()
    }()

    for {
        _, message, err := c.Conn.ReadMessage()
        if err != nil {
            log.Println("read error:", err)
            break
        }

        // Optional: handle incoming messages
        log.Println("Received:", string(message))
    }
}

func (c *Client) WritePump() {
    defer c.Conn.Close()

    for msg := range c.Send {
        err := c.Conn.WriteMessage(websocket.TextMessage, msg)
        if err != nil {
            log.Println("write error:", err)
            break
        }
    }
}

func WebSocketHandler(hub *Hub) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Println("upgrade error:", err)
            return
        }

        userID := r.URL.Query().Get("user") // e.g. /ws?user=123
        client := &Client{
            ID:   userID,
            Conn: conn,
            Send: make(chan []byte),
            Hub:  hub,
        }

        hub.register <- client

        go client.WritePump()
        go client.ReadPump()
    }
}
