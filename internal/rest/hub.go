package rest

import (
	"fmt"
)

type Hub struct {
	clients    map[*Client]bool
	wsChan     chan WsPayload
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		wsChan:     make(chan WsPayload),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case msg := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- Message{
					ID:   client.id,
					data: msg,
				}:
				default:
					fmt.Println(h.clients[client])
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
