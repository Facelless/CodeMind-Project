package usecase

import (
	"context"
	"miservicegolang/core/domain/ai"
	"miservicegolang/core/pkg"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroqDatabaseRepo interface {
	Insert(ctx context.Context, response ai.GroqAi) (primitive.ObjectID, pkg.Log)
	FindByID(ctx context.Context, id primitive.ObjectID) (ai.GroqAi, pkg.Log)
}

type GroqDatabase struct {
	client *mongo.Client
}

func NewGroqAiDatabaseRepo(client *mongo.Client) GroqDatabaseRepo {
	return &GroqDatabase{client: client}
}

func (g *GroqDatabase) Insert(ctx context.Context, r ai.GroqAi) (primitive.ObjectID, pkg.Log) {
	collection := g.client.Database("saas").Collection("ia")

	inserted, err := collection.InsertOne(ctx, r)
	if err != nil {
		return primitive.NilObjectID, pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "Insert failed", "err": err.Error()},
		}
	}

	id := inserted.InsertedID.(primitive.ObjectID)
	return id, pkg.Log{}
}

func (g *GroqDatabase) FindByID(ctx context.Context, id primitive.ObjectID) (ai.GroqAi, pkg.Log) {
	collection := g.client.Database("saas").Collection("ia")

	var result ai.GroqAi
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return ai.GroqAi{}, pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "Document not found", "err": err.Error()},
		}
	}

	return result, pkg.Log{}
}
