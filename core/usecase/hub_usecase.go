package usecase

import "github.com/gorilla/websocket"

type Hub struct {
	rooms map[string]map[*websocket.Conn]bool
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]map[*websocket.Conn]bool),
	}
}

func (h *Hub) GetRoom(roomID string) map[*websocket.Conn]bool {
	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	return h.rooms[roomID]
}

func (h *Hub) Join(roomID string, conn *websocket.Conn) {
	room := h.GetRoom(roomID)
	room[conn] = true
}

func (h *Hub) Leave(roomID string, conn *websocket.Conn) {
	if room, ok := h.rooms[roomID]; ok {
		delete(room, conn)
	}
}

func (h *Hub) Broadcast(roomID string, msg []byte) {
	if room, ok := h.rooms[roomID]; ok {
		for conn := range room {
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
