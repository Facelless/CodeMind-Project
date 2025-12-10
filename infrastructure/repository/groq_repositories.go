package repository

import (
	"context"
	"errors"
	"miservicegolang/core/domain/ai"
	"miservicegolang/core/pkg"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroqDatabaseRepo interface {
	Insert(ctx context.Context, response ai.GroqAi) (primitive.ObjectID, pkg.Log)
	FindByID(ctx context.Context, id primitive.ObjectID) (ai.GroqAi, pkg.Log)
	UpdateVerify(ctx context.Context, id primitive.ObjectID, verify string, completed bool) pkg.Log
	Delete(ctx context.Context, id primitive.ObjectID) pkg.Log
}

type GroqDatabase struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewGroqAiDatabaseRepo(client *mongo.Client) GroqDatabaseRepo {
	return &GroqDatabase{
		client:     client,
		collection: client.Database("saas").Collection("ia"),
	}
}

func (g *GroqDatabase) Insert(ctx context.Context, r ai.GroqAi) (primitive.ObjectID, pkg.Log) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := g.collection.InsertOne(c, r)
	if err != nil {
		return primitive.NilObjectID, pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "Insert failed", "err": err.Error()},
		}
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "Invalid ObjectID returned"},
		}
	}

	return oid, pkg.Log{}
}

func (g *GroqDatabase) FindByID(ctx context.Context, id primitive.ObjectID) (ai.GroqAi, pkg.Log) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var result ai.GroqAi
	err := g.collection.FindOne(c, bson.M{"_userid": id}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ai.GroqAi{}, pkg.Log{
				Error: true,
				Body:  map[string]any{"message": "Document not found"},
			}
		}

		return ai.GroqAi{}, pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "Find failed", "err": err.Error()},
		}
	}

	return result, pkg.Log{}
}

func (g *GroqDatabase) UpdateVerify(ctx context.Context, id primitive.ObjectID, verify string, completed bool) pkg.Log {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"verify":    verify,
			"completed": completed,
			"updatedAt": time.Now(),
		},
	}

	result, err := g.collection.UpdateOne(c, bson.M{"_id": id}, update)
	if err != nil {
		return pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "Update failed", "err": err.Error()},
		}
	}

	if result.MatchedCount == 0 {
		return pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "Document not found"},
		}
	}

	return pkg.Log{}
}

func (g *GroqDatabase) Delete(ctx context.Context, id primitive.ObjectID) pkg.Log {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := g.collection.DeleteOne(c, bson.M{"_id": id})
	if err != nil {
		return pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "Update failed", "err": err.Error()},
		}
	}

	return pkg.Log{}
}
