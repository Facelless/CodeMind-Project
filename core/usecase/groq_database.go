package usecase

import (
	"context"
	"miservicegolang/core/domain/ai"
	"miservicegolang/core/pkg"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroqDatabaseRepo interface {
	Insert(ctx context.Context, response ai.GroqAi) (primitive.ObjectID, pkg.Log)
}

type GroqDatabase struct {
	client *mongo.Client
}

func NewGroqAiDatabaseRepo(client *mongo.Client) GroqDatabaseRepo {
	return &GroqDatabase{client: client}
}

func (g *GroqDatabase) Insert(ctx context.Context, r ai.GroqAi) (primitive.ObjectID, pkg.Log) {

	collection := g.client.Database("saas").Collection("ia")

	insertResponse, err := collection.InsertOne(ctx, r)
	if err != nil {
		return primitive.NilObjectID, pkg.Log{
			Error: true,
			Body: map[string]any{
				"message": "failed to insert document",
				"error":   err.Error(),
			},
		}
	}

	oid, ok := insertResponse.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, pkg.Log{
			Error: true,
			Body: map[string]any{
				"message": "failed to parse inserted ID",
			},
		}
	}

	return oid, pkg.Log{}
}
