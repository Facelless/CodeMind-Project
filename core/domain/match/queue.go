package match

import (
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Queue struct {
	Conn *websocket.Conn
	Type string             `json:"type"`
	Id   primitive.ObjectID `json:"id"`
}
