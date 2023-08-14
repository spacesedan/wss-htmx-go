package rest

import (
	"log"
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newLine = []byte{'\n'}
	space   = []byte{' '}
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
	go hub.run()
	return &WssHandler{
		hub: hub,
	}
}

func (h *WssHandler) Register(m *chi.Mux) {
	m.HandleFunc("/ws", h.Serve)
}

type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

type WsPayload struct {
	Action  string `json:"action"`
	ID      string `json:"id"`
	Message string `json:"message"`
	Conn    Client `json:"_-"`
}

func (h *WssHandler) Serve(w http.ResponseWriter, r *http.Request) {
	log.Println("Connected to socket")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}

	cID, err := uuid.NewRandom()
	if err != nil {
		return
	}

	var response WsJsonResponse
	response.Message = `<div id="poop" hx-swap-oob="true"><em><small>Connected to server</small></em></div>`

	err = conn.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	client := &Client{
		hub:  h.hub,
		conn: WsConnection{Conn: conn},
		send: make(chan Message),
		id:   cID.String(),
	}

	client.hub.register <- client

	go client.ListenForWs()
	go client.readPump()
}
