package handler

import (
	"BankKibikov/internal/chat"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type ChatHandler struct {
	hub *chat.Hub
}

func NewChatHandler(hub *chat.Hub) *ChatHandler {
	return &ChatHandler{hub: hub}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *ChatHandler) ServeWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := &chat.Client{Conn: conn, Send: make(chan []byte, 256)}
	h.hub.Register <- client

	go func() {
		defer func() {
			h.hub.Unregister <- client
			conn.Close()
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			h.hub.Broadcast <- message
		}
	}()

	go func() {
		for msg := range client.Send {
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}()
}
