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
		fmt.Printf("%+v", e)
		switch e.Action {
		case "broadcast":
			fmt.Println("Broadcast")
		case "alert":
		case "message":
			response.Action = "message"
			response.Message = fmt.Sprintf(`<div id="messages" hx-swap-oob="beforeend"><p id="message"><strong>%v says:</strong> %v</p></div>`, e.User, e.Message)
			h.broadcastToAll(response)
		case "list_users":
			fmt.Println("Listing users")
		case "connect":
		case "left":
			response.SkipSender = false
			response.CurrentConn = e.Conn
			response.Action = "left"
			response.Message = fmt.Sprintf(`<p id="leavers" hx-swap-oob="true">%v left, bye bye.</p>`, e.User)
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

		case "entered":
			userList := h.addToUserList(e.Conn, e.User)
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
		// fmt.Printf("%+v\n", payload)

		if err != nil {
			// Do nothing...
			fmt.Print(err)
		} else {
			fmt.Printf("%v\n", payload)
			payload.Conn = *conn
			h.wsChan <- payload
		}
	}
}

func (h *Hub) addToUserList(conn WsConnection, u string) []string {
	var userNames []string
	h.clients[conn] = u
	fmt.Printf("This: %v\n", h.clients[conn])
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

		fmt.Println(response.Message)
		err := client.WriteMessage(websocket.TextMessage, []byte(response.Message))
		if err != nil {
			log.Printf("Websocket error on %s: %s", response.Action, err)
			_ = client.Close()
			delete(h.clients, client)
		}
	}
}
