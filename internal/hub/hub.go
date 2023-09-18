package hub

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

type WsConnection struct {
	*websocket.Conn
}

type WsJsonResponse struct {
	Action         string       `json:"action"`
	Message        string       `json:"message"`
	MessageType    string       `json:"message_type"`
	SkipSender     bool         `json:"-"`
	IsSender       bool         `json:""`
	CurrentConn    WsConnection `json:"-"`
	ConnectedUsers []string     `json:"-"`
}

// WsPayload contains the information comming from the websocket connection
type WsPayload struct {
	// HEADERS is injected to the message by htmx
	Headers map[string]string `json:"HEADERS"`
	Action  string            `json:"action"`
	ID      string            `json:"id"`
	User    string            `json:"user"`
	Message string            `json:"message"`
	Conn    WsConnection      `json:"-"`
}

type Hub struct {
	logger  *slog.Logger
	Clients map[WsConnection]string
	WsChan  chan WsPayload
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		Clients: make(map[WsConnection]string),
		WsChan:  make(chan WsPayload),
		logger:  logger,
	}
}

func (h *Hub) ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-h.WsChan
		switch e.Action {
		case "message":
			h.logger.Info("Message recieved", slog.String("message_id", e.ID))
			h.handleChatMessage(e, response)
		case "left":

			response.SkipSender = false
			response.CurrentConn = e.Conn
			response.Action = "left"
			h.broadcastToAll(response)

			delete(h.Clients, e.Conn)
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
			h.WsChan <- payload
		}
	}
}

func (h *Hub) addToUserList(conn WsConnection, u string) []string {
	var userNames []string
	h.Clients[conn] = u
	for _, value := range h.Clients {
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
	for _, value := range h.Clients {
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
	for client := range h.Clients {
		if response.SkipSender && response.CurrentConn == client {
			continue
		}

		err := client.WriteMessage(websocket.TextMessage, []byte(response.Message))
		if err != nil {
			h.logger.Error("Error writing message", slog.String("action", response.Action), slog.String("err", err.Error()))
			_ = client.Close()
			delete(h.Clients, client)
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

	for client := range h.Clients {
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
			delete(h.Clients, client)
		}
	}
}
