package handlers

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
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
		fmt.Printf("%+v\n", e)
		switch e.Action {
		case "message":
			h.handleChatMessage(e, response)
			h.broadcastToAll(response)
		case "connect":
		case "left":
			response.SkipSender = false
			response.CurrentConn = e.Conn
			response.Action = "left"
			h.broadcastToAll(response)

			fmt.Println("Before delete", h.clients)
			delete(h.clients, e.Conn)
			fmt.Println("After delete", h.clients)
			userList := h.getUserNameList()
			response.Action = "list_users"
			response.ConnectedUsers = userList
			response.SkipSender = false
			var userHtml []string
			for _, value := range response.ConnectedUsers {
				userHtml = append(userHtml, fmt.Sprintf(`<li>%v</li>`, value))
			}
			response.Message = fmt.Sprintf(`<ul id="chat_connected_users" hx-swap="innerHTML">%v</ul>`, strings.Join(userHtml, ""))
			h.broadcastToAll(response)

		case "entered":
			userList := h.addToUserList(e.Conn, e.User)
			response.Action = "list_users"
			response.ConnectedUsers = userList
			response.SkipSender = false
			var userHtml []string
			for _, value := range response.ConnectedUsers {
				userHtml = append(userHtml, fmt.Sprintf(`<li>%v</li>`, value))
			}
			response.Message = fmt.Sprintf(`<ul id="chat_connected_users" hx-swap="innerHTML">%v</ul>`, strings.Join(userHtml, ""))

			h.broadcastToAll(response)
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
		// fmt.Printf("%+v\n", payload)

		if err != nil || payload.Message == "" {
			// Do nothing...
		} else {
			payload.Conn = *conn
			h.wsChan <- payload
		}
	}
}

func (h *Hub) addToUserList(conn WsConnection, u string) []string {
	var userNames []string
	h.clients[conn] = u
	for _, value := range h.clients {
		if value != "" {
			if slices.Contains(userNames, value) {
				continue
			}
			userNames = append(userNames, value)
		}
	}
	sort.Strings(userNames)
	return userNames
}

func (h *Hub) getUserNameList() []string {
	var userNames []string
	for _, value := range h.clients {
		if value != "" {
			if slices.Contains(userNames, value) {
				continue
			}
			userNames = append(userNames, value)
		}
	}
	sort.Strings(userNames)
	return userNames
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

// handleChatMessage takes in the payload from the client and updates the
// response to be returned back to the client. It returns an HTML element that
// gets pushed onto the chat messages box.
func (h *Hub) handleChatMessage(payload WsPayload, response WsJsonResponse) {
	response.Action = "message"
	response.Message = fmt.Sprintf(`<div id="chat_messages" hx-swap-oob="beforeend"><p id="message"><strong>%v:</strong> %v</p></div>`, payload.User, payload.Message)
}
