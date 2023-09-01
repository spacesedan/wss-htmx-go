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
		switch e.Headers["action"] {
		case "broadcast":
		case "alert":
		case "message":
			fmt.Println(e)
			response.Action = "message"
			response.Message = fmt.Sprintf(`<div id="messages" hx-swap-oob="beforeend" hx-swap="scroll:bottom"><p id="message"><strong>%v says:</stong> %v</p></div>`, e.Headers["user"], e.Message)
			// h.broadcastToAll(response)
		case "list_users":
			fmt.Println("Listing users")
		case "connect":
		case "left":
			fmt.Printf("%v left", e.Headers["user"])
			response.SkipSender = false
			response.CurrentConn = e.Conn
			response.Action = "left"
			response.Message = fmt.Sprintf(`<p id="leavers" hx-swap-oob="true">%v left, bye bye.</p>`, e.Headers["user"])
			h.broadcastToAll(response)

			delete(h.clients, e.Conn)
			userList := h.getUserNameList()
			response.Action = "list_users"
			response.ConnectedUsers = userList
			response.SkipSender = false
			var userHtml []string
			for _, value := range response.ConnectedUsers {
				userHtml = append(userHtml, fmt.Sprintf(`<li>%v</li>`, value))
			}
			response.Message = fmt.Sprintf(`<ul id="users_list" hx-swap-oob="true">%v</ul>`, strings.Join(userHtml, ""))
			h.broadcastToAll(response)

		case "add_user":
			userList := h.addToUserList(e.Conn, e.Headers["user"])
			response.Action = "list_users"
			response.ConnectedUsers = userList
			response.SkipSender = false
			var userHtml []string
			for _, value := range response.ConnectedUsers {
				userHtml = append(userHtml, fmt.Sprintf(`<li>%v</li>`, value))
			}
			response.Message = fmt.Sprintf(`<ul id="users_list" hx-swap-oob="true">%v</ul>`, strings.Join(userHtml, ""))

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
		// sends the response in JSON requires a more hacky solution when working
		// with htmx
		// err = client.WriteJSON(response)
		// if err != nil {
		// 	log.Printf("Websocket error on %s: %s", response.Action, err)
		// 	_ = client.Close()
		// 	delete(h.clients, client)
		// }
	}
}
