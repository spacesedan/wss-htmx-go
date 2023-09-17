package handlers

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

type Hub struct {
	logger  *slog.Logger
	clients map[WsConnection]string
	wsChan  chan WsPayload
}

func newHub(logger *slog.Logger) *Hub {
	return &Hub{
		clients: make(map[WsConnection]string),
		wsChan:  make(chan WsPayload),
		logger:  logger,
	}
}

func (h *Hub) ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-h.wsChan
		switch e.Action {
		case "message":
			h.logger.Info("Message recieved", slog.String("message_id", e.ID))
			h.handleChatMessage(e, response)
		case "left":

			response.SkipSender = false
			response.CurrentConn = e.Conn
			response.Action = "left"
			h.broadcastToAll(response)

			delete(h.clients, e.Conn)
			userList := h.getUserNameList()
			response.Action = "list_users"
			response.ConnectedUsers = userList
			response.SkipSender = false
			var userHtml []string
			for _, value := range response.ConnectedUsers {
				userHtml = append(userHtml, fmt.Sprintf(`<li class="font-bold">%v</li>`, value))
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
				userHtml = append(userHtml, fmt.Sprintf(`<li class="font-bold">%v</li>`, value))
			}
			response.Message = fmt.Sprintf(`<ul id="chat_connected_users" hx-swap="innerHTML">%v</ul>`, strings.Join(userHtml, ""))

			h.broadcastToAll(response)
		}
	}
}

func (h *Hub) ListenForWS(conn *WsConnection) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Error("Error: Attempting to recover", slog.Any("err", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)

		if err != nil {
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
			h.logger.Error("Error writing message", slog.String("action", response.Action), slog.String("err", err.Error()))
			_ = client.Close()
			delete(h.clients, client)
		}
	}
}

func (h *Hub) handleChatMessage(payload WsPayload, response WsJsonResponse) {
	response.Action = "message"
	response.CurrentConn = payload.Conn

	// prevents empty messages from being sent to connected clients
	if payload.Message == "" {
		return
	}

	for client := range h.clients {
		if response.CurrentConn == client {
			response.Message = fmt.Sprintf(`
          <div id="chat_messages" hx-swap-oob="beforeend">
            <div id="message" class="flex gap-3 justify-end items-start p-3 font-mono">
              <p class="bg-green-400 px-3 py-2 rounded-md">%v</p>
              <img src="https://ui-avatars.com/api/?name=%v&size=32&rounded=true" alt="profile image for user: %v"></img>
            </div>
          </div>
      `, payload.Message, payload.User, payload.User)
		} else {
			response.Message = fmt.Sprintf(`
          <div id="chat_messages" hx-swap-oob="beforeend">
            <div id="message" class="flex gap-3 justify-start items-start p-3 font-mono">
              <img src="https://ui-avatars.com/api/?name=%v&size=32&rounded=true" alt="profile image for user: %v"></img>
              <p class="bg-indigo-400 px-3 py-2 rounded-md">%v</p>
            </div>
          </div>
      `, payload.User, payload.User, payload.Message)
		}

		err := client.WriteMessage(websocket.TextMessage, []byte(response.Message))
		if err != nil {
			h.logger.Error("Error writing message", slog.String("action", response.Action), slog.String("err", err.Error()))
			_ = client.Close()
			delete(h.clients, client)
		}
	}
}
