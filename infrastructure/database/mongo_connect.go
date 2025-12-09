package database

import (
	"context"
	"miservicegolang/core/pkg"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongodb() (*mongo.Client, pkg.Log) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, pkg.Log{
			Error: true,
			Body: map[string]any{
				"message": "Error loading .env",
			},
		}
	}
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, pkg.Log{
			Error: true,
			Body: map[string]any{
				"message": "Error connecting to mongodb.",
				"err":     err,
			},
		}
	}
	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, pkg.Log{
			Error: true,
			Body: map[string]any{
				"message": "Error sending connection ping",
				"err":     err,
			},
		}
	}
	return client, pkg.Log{
		Error: false,
		Body: map[string]any{
			"message": "MongoClient connected",
		},
	}
}
