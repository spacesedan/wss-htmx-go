package handlers

import (
	"log/slog"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/spacesedan/wss-htmx-go/internal/hub"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WssHandler struct {
	hub *hub.Hub
    logger *slog.Logger
}

func NewWssHandler(hub *hub.Hub,logger *slog.Logger) *WssHandler {
	go hub.ListenToWsChannel()
	return &WssHandler{
		hub: hub,
        logger: logger,
	}
}

func (h *WssHandler) Register(m *chi.Mux) {
	m.HandleFunc("/ws", h.Serve)
}


func (h *WssHandler) Serve(w http.ResponseWriter, r *http.Request) {

    h.logger.Info("Connected to socket")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
        h.logger.Error("Something went wrong when upgrading connection",
        slog.String("err", err.Error()))
		return
	}

	var response hub.WsJsonResponse
	response.Action = `connected`
	response.Message = `<p id="wsStatus">Welcome to the startup</p>`
	// err = ws.WriteJSON(response)
	err = ws.WriteMessage(websocket.TextMessage, []byte(response.Message))
	if err != nil {
        h.logger.Error("Something when trying to send a message to the client",
        slog.String("err", err.Error()))
	}

	conn := hub.WsConnection{Conn: ws}
	h.hub.Clients[conn] = ""

	go h.hub.ListenForWS(&conn)
}
