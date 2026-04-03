package core

import (
    "log"
    "net/http"
    "time"
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


const (
    writeWait = 10 * time.Second
    pongWait  = 30 * time.Second
    pingPeriod = (pongWait * 9) / 10 // must be < pongWait
)


func (c *Client) ReadPump() {
    defer func() {
        c.Hub.unregister <- c
        c.Conn.Close()
    }()

    c.Conn.SetReadLimit(512)
    c.Conn.SetReadDeadline(time.Now().Add(pongWait))
    c.Conn.SetPongHandler(func(string) error {
        c.Conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, _, err := c.Conn.ReadMessage()
        if err != nil {
            log.Println("read error:", err)
            break
        }
    }
}

func (c *Client) WritePump() {
    ticker := time.NewTicker(pingPeriod)

    defer func() {
        ticker.Stop()
        c.Conn.Close()
    }()

    for {
        select {
        case msg, ok := <-c.Send:
            c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
                return
            }

        case <-ticker.C:
            c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
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
