package handlers

import (
	"log"
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
}

func NewWssHandler() *WssHandler {
	hub := newHub()
	go hub.ListenToWsChannel()
	return &WssHandler{
		hub: hub,
	}
}

func (h *WssHandler) Register(m *chi.Mux) {
	m.HandleFunc("/ws", h.Serve)
}

type WsConnection struct {
	*websocket.Conn
}

type WsJsonResponse struct {
	Action        string       `json:"action"`
	Message       string       `json:"message"`
	MessageType   string       `json:"message_type"`
	SkipSender    bool         `json:"-"`
	CurrentConn   WsConnection `json:"-"`
	ConnectedUser []string     `json:"-"`
}

type WsPayload struct {
	Action  string       `json:"action"`
	ID      string       `json:"id"`
	Message string       `json:"message"`
	Conn    WsConnection `json:"-"`
}

func (h *WssHandler) Serve(w http.ResponseWriter, r *http.Request) {
	log.Println("Connected to socket")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}

	var response WsJsonResponse
	response.Message = `<span id="status" hx-swap-oob="true"> Welcome to the startup...</span>`
	err = ws.WriteMessage(websocket.TextMessage, []byte(response.Message))
	if err != nil {
		log.Println(err)
	}

	conn := WsConnection{Conn: ws}
	h.hub.clients[conn] = ""

	go h.hub.ListenForWS(&conn)
}
