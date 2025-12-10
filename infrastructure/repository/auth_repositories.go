package repository

import (
	"context"
	"errors"
	"miservicegolang/core/domain/user"
	"miservicegolang/core/pkg"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserDatabaseRepo interface {
	Insert(ctx context.Context, dates user.User) (primitive.ObjectID, pkg.Log)
	FindById(ctx context.Context, id primitive.ObjectID) (user.User, pkg.Log)
	FindByEmail(ctx context.Context, email string) (user.User, pkg.Log)
	Delete(ctx context.Context, id primitive.ObjectID) pkg.Log
	Update(ctx context.Context, id primitive.ObjectID, dates bson.M) pkg.Log
	Find(ctx context.Context, filter bson.M, filterOp bson.M) (*mongo.Cursor, pkg.Log)
}

type UserDatabase struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewUserDatabaseRepo(client *mongo.Client) UserDatabaseRepo {
	return &UserDatabase{
		client:     client,
		collection: client.Database("saas").Collection("users"),
	}
}

func (u *UserDatabase) Insert(ctx context.Context, dates user.User) (primitive.ObjectID, pkg.Log) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := u.collection.InsertOne(c, dates)
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
			Body:  map[string]any{"message": "Invalid ObjectId returned"},
		}
	}

	return oid, pkg.Log{}
}

func (u *UserDatabase) FindById(ctx context.Context, id primitive.ObjectID) (user.User, pkg.Log) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var result user.User
	err := u.collection.FindOne(c, bson.M{"_id": id}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.User{}, pkg.Log{
				Error: true,
				Body:  map[string]any{"message": "Document not found"},
			}
		}
		return user.User{}, pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "Find failed", "err": err.Error()},
		}
	}
	return result, pkg.Log{}
}

func (u *UserDatabase) FindByEmail(ctx context.Context, email string) (user.User, pkg.Log) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var result user.User
	err := u.collection.FindOne(c, bson.M{"email": email}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.User{}, pkg.Log{
				Error: true,
				Body:  map[string]any{"message": "Document not found"},
			}
		}
		return user.User{}, pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "Find failed", "err": err.Error()},
		}
	}
	return result, pkg.Log{}
}

func (u *UserDatabase) Delete(ctx context.Context, id primitive.ObjectID) pkg.Log {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := u.collection.DeleteOne(c, bson.M{"_id": id})
	if err != nil {
		return pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "Update failed", "err": err.Error()},
		}
	}

	return pkg.Log{}
}

func (u *UserDatabase) Update(ctx context.Context, id primitive.ObjectID, dates bson.M) pkg.Log {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result, err := u.collection.UpdateOne(c, bson.M{"_id": id}, dates)
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

func (u *UserDatabase) Find(ctx context.Context, filter bson.M, filterOp bson.M) (*mongo.Cursor, pkg.Log) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := u.collection.Find(c, filter, options.Find().SetSort(filterOp))
	if err != nil {
		return nil, pkg.Log{Error: true, Body: map[string]any{"err": err.Error()}}
	}

	return result, pkg.Log{}
}
