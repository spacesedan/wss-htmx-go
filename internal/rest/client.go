package rest

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	ID   string
	data []byte
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan Message
	id   string
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newLine, space, -1))
		c.hub.broadcast <- msg
	}
}

type ChatMessagePayload struct {
	ID          string `json:"id"`
	ChatMessage string `json:"chat_message"`
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			var cm ChatMessagePayload

			_ = json.Unmarshal([]byte(msg.data), &cm)

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write([]byte(`<div id="chat_room" hx-swap-oob="beforeend">` + msg.ID + ":" + cm.ChatMessage + `<br></div>`))

			// n := len(c.send)
			// for i := 0; 1 < n; i++ {
			// 	_, _ = w.Write(newLine)
			// 	_, _ = w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
