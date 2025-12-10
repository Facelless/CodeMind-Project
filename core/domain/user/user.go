package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email    *string            `json:"email"`
	Password string             `json:"password"`
	Username string             `json:"username"`
	Avatar   string             `json:"avatar"`
	Exp      int                `json:"exp"`
	Level    int                `json:"level"`
}
