package ai

import "go.mongodb.org/mongo-driver/bson/primitive"

type GroqAi struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"user_id" bson:"_userid, omitempty"`
	Answer    string             `json:"answer"`
	Completed bool               `json:"completed"`
	Verify    string             `json:"verify"`
}
