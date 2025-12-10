package repository

import (
	"context"
	"miservicegolang/core/pkg"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GrowthRepo interface {
	SetExp(id primitive.ObjectID) pkg.Log
}

type Growth struct {
	db UserDatabaseRepo
}

func NewGrowthRepo(db UserDatabaseRepo) GrowthRepo {
	return &Growth{db: db}
}

func (u *Growth) SetExp(id primitive.ObjectID) pkg.Log {
	user, log := u.db.FindById(context.Background(), id)
	if log.Error {
		return log
	}

	if user.Exp >= 100 {
		return u.db.Update(context.Background(), id, bson.M{
			"$set": bson.M{"exp": 0},
			"$inc": bson.M{"level": 1},
		})
	}

	return u.db.Update(context.Background(), id, bson.M{
		"$inc": bson.M{"exp": 10},
	})
}
