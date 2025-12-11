package server

import (
	"log"
	"miservicegolang/core/domain/match"
	"miservicegolang/core/usecase"
	"miservicegolang/infrastructure/repository"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebsocketHandler struct {
	usecase *usecase.MatchUsecase
	db      repository.UserDatabaseRepo
}

func NewWebsocketHandler(u *usecase.MatchUsecase, db repository.UserDatabaseRepo) *WebsocketHandler {
	return &WebsocketHandler{
		usecase: u,
		db:      db,
	}
}

var up = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *WebsocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := up.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "error connecting websocket", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	var eventChan chan match.MatchEvent
	var oid primitive.ObjectID

	for {
		var msg map[string]any
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("ws closed:", err)
			return
		}

		action := msg["action"].(string)

		switch action {

		case "looking_for_match":
			idStr, _ := msg["player_id"].(string)
			oid, err = primitive.ObjectIDFromHex(idStr)
			if err != nil {
				conn.WriteJSON(map[string]any{"error": "invalid_object_id"})
				continue
			}

			if eventChan == nil {
				eventChan = h.usecase.RegisterClient(oid)

				go func() {
					for evt := range eventChan {
						conn.WriteJSON(evt)
					}
				}()
			}

			userData, err := h.db.FindById(r.Context(), oid)
			if err.Error {
				conn.WriteJSON(map[string]any{"error": "user_not_found"})
				continue
			}

			h.usecase.PlayerLookingForMatch(userData)
		}
	}
}
