package server

import (
	"encoding/json"
	"log"
	"miservicegolang/core/domain/match"
	"miservicegolang/core/usecase"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebsocketHandler struct {
	Queue *usecase.QueueUsecase
	Hub   *usecase.Hub
}

func NewWebsocketHandler(queue *usecase.QueueUsecase, hub *usecase.Hub) *WebsocketHandler {
	return &WebsocketHandler{
		Queue: queue,
		Hub:   hub,
	}
}

var up = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *WebsocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := up.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	h.Hub.Join("lobby", conn)

	defer func() {
		h.Hub.Leave("lobby", conn)
		conn.Close()
	}()

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var msg match.Queue
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			continue
		}

		if msg.Id == primitive.NilObjectID {
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid id"))
			continue
		}

		if msg.Type == "looking_for_match" {
			h.Queue.Queue <- &match.Queue{
				Id:   msg.Id,
				Conn: conn,
			}
		}

	}
}
