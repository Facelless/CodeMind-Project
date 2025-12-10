package usecase

import (
	"context"
	"miservicegolang/core/domain/user"
	"miservicegolang/core/pkg"
	"miservicegolang/infrastructure/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GrowthUsecase interface {
	SetExp(id primitive.ObjectID) pkg.Log
	Rank() pkg.Log
}

type Growth struct {
	db repository.UserDatabaseRepo
}

func NewGrowthUsecase(db repository.UserDatabaseRepo) GrowthUsecase {
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

func (u *Growth) Rank() pkg.Log {
	cursor, err := u.db.Find(context.Background(), bson.M{"level": bson.M{"$exists": true}}, bson.M{"level": -1})
	if err.Error {
		return err
	}
	var results []user.User
	if cursor.All(context.Background(), &results); err.Error {
		return err
	}

	return pkg.Log{
		Error: false,
		Body: map[string]any{
			"Results": results,
		},
	}
}
