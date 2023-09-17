package handlers

import (
	"log/slog"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WssHandler struct {
	hub *Hub
    logger *slog.Logger
}

func NewWssHandler(logger *slog.Logger) *WssHandler {
	hub := newHub(logger)
	go hub.ListenToWsChannel()
	return &WssHandler{
		hub: hub,
        logger: logger,
	}
}

func (h *WssHandler) Register(m *chi.Mux) {
	m.HandleFunc("/ws", h.Serve)
}

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

func (h *WssHandler) Serve(w http.ResponseWriter, r *http.Request) {

    h.logger.Info("Connected to socket")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
        h.logger.Error("Something went wrong when upgrading connection",
        slog.String("err", err.Error()))
		return
	}

	var response WsJsonResponse
	response.Action = `connected`
	response.Message = `<p id="wsStatus">Welcome to the startup</p>`
	// err = ws.WriteJSON(response)
	err = ws.WriteMessage(websocket.TextMessage, []byte(response.Message))
	if err != nil {
        h.logger.Error("Something when trying to send a message to the client",
        slog.String("err", err.Error()))
	}

	conn := WsConnection{Conn: ws}
	h.hub.clients[conn] = ""

	go h.hub.ListenForWS(&conn)
}
