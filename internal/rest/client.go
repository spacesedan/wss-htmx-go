package rest

import (
	"log"

	"github.com/gorilla/websocket"
)

type Message struct {
	ID   string
	data []byte
}

type WsConnection struct {
	*websocket.Conn
}

type Client struct {
	conn WsConnection
}

type ChatMessagePayload struct {
	ID          string `json:"id"`
	ChatMessage string `json:"chat_message"`
}

func (c *Client) ListenForWs() {
	defer func() {
		r := recover()
		if r != nil {
			log.Printf("%v\n", r)
		}
	}()

	var payload WsPayload

	for {
		err := c.conn.ReadJSON(&payload)
		if err != nil {
			log.Print(err)
		} else {
			c.hub.wsChan <- payload
		}

	}
}

func (c *Client) ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-c.hub.wsChKan
		switch e.Action {
		case "username":
			c.hub.clients()
		}
	}
}
