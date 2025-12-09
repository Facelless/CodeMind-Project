package ai

import "go.mongodb.org/mongo-driver/bson/primitive"

type GroqAi struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Prompt string             `bson:"prompt" json:"prompt"`
	Answer string             `bson:"answer" json:"answer"`
}
