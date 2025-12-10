package usecase

import (
	"context"
	"miservicegolang/core/domain/user"
	"miservicegolang/core/pkg"
	"miservicegolang/infrastructure/repository"
	"miservicegolang/infrastructure/security"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	db repository.UserDatabaseRepo
}

func NewUserUsecase(db repository.UserDatabaseRepo) *UserUsecase {
	return &UserUsecase{db: db}
}

func (u *UserUsecase) Register(dates user.User) (user.User, pkg.Log) {
	exists, _ := u.db.FindByEmail(context.Background(), *dates.Email)

	if exists.Email != nil {
		return user.User{}, pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "email already registered"},
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(dates.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.User{}, pkg.Log{
			Error: true,
			Body:  map[string]any{"Err": err.Error()},
		}
	}

	dates.Password = string(hash)
	result := user.User{
		Email:    dates.Email,
		Password: dates.Password,
		Username: dates.Username,
		Avatar:   "",
		Exp:      0,
		Level:    0,
	}
	id, log := u.db.Insert(context.Background(), result)
	if log.Error {
		return user.User{}, log
	}

	result.ID = id
	return result, pkg.Log{}

}

func (u *UserUsecase) Login(email, password string) (string, pkg.Log) {
	dates, log := u.db.FindByEmail(context.Background(), email)
	if log.Error {
		return "", log
	}
	if bcrypt.CompareHashAndPassword([]byte(dates.Password), []byte(password)) != nil {
		return "", pkg.Log{
			Error: true,
			Body: map[string]any{
				"message": "Incorrect password",
			},
		}
	}

	token, err := security.GenereateToken(dates.ID)
	if err != nil {
		return "", pkg.Log{
			Error: true,
			Body:  map[string]any{"err": err.Error()},
		}
	}
	return token, pkg.Log{}
}
