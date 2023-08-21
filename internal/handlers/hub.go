package handlers

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients        map[WsConnection]string
	wsChan         chan WsPayload
	connectionChan chan WsPayload
	broadcastChan  chan WsPayload
	alertChan      chan WsPayload
	whoIsThereChan chan WsPayload
	enterChan      chan WsPayload
	leaveChan      chan WsPayload
	userName       chan WsPayload
}

func newHub() *Hub {
	return &Hub{
		clients:        make(map[WsConnection]string),
		wsChan:         make(chan WsPayload),
		connectionChan: make(chan WsPayload),
		broadcastChan:  make(chan WsPayload),
		alertChan:      make(chan WsPayload),
		whoIsThereChan: make(chan WsPayload),
		enterChan:      make(chan WsPayload),
		leaveChan:      make(chan WsPayload),
		userName:       make(chan WsPayload),
	}
}

func (h *Hub) ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-h.wsChan
		switch e.Headers["action"] {
		case "broadcast":
		case "alert":
		case "message":
			fmt.Println(e.Headers)
			response.Action = "message"
			response.Message = fmt.Sprintf(`<div id="messages" hx-swap-oob="beforeend" hx-swap="scroll:bottom"><p id="message"><strong>%v says:</stong> %v</p></div>`, e.Headers["user"], e.Message)
			h.broadcastToAll(response)
		case "list_users":
		case "connect":
		case "left":
		case "username":
			response.Action = "list_users"
		}
	}
}

func (h *Hub) ListenForWS(conn *WsConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil || payload.Message == "" {
			// Do nothing...
		} else {
			payload.Conn = *conn
			h.wsChan <- payload
		}
	}
}

func (h *Hub) broadcastToAll(response WsJsonResponse) {
	for client := range h.clients {
		if response.SkipSender && response.CurrentConn == client {
			continue
		}

		err := client.WriteMessage(websocket.TextMessage, []byte(response.Message))
		if err != nil {
			log.Printf("Websocket error on %s: %s", response.Action, err)
			_ = client.Close()
			delete(h.clients, client)
		}
	}
}
