package ai

import "go.mongodb.org/mongo-driver/bson/primitive"

type GroqAi struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Answer string             `bson:"answer" json:"answer"`
}
