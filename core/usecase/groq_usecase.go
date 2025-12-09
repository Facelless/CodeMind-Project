package usecase

import (
	"context"
	"miservicegolang/core/domain/ai"
	"miservicegolang/core/pkg"
	"miservicegolang/infrastructure/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AiUsecase struct {
	repo repository.GroqAiRepo
	db   GroqDatabaseRepo
}

func NewAiUsecase(r repository.GroqAiRepo, db GroqDatabaseRepo) *AiUsecase {
	return &AiUsecase{repo: r, db: db}
}

func (u *AiUsecase) Generate(prompt string) (ai.GroqAi, pkg.Log) {
	if prompt == "" {
		return ai.GroqAi{}, pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "Prompt cannot be empty"},
		}
	}

	text, log := u.repo.GenerateText(prompt)
	if log.Error {
		return ai.GroqAi{}, log
	}

	result := ai.GroqAi{Answer: text}

	id, log := u.db.Insert(context.Background(), result)
	if log.Error {
		return ai.GroqAi{}, log
	}

	result.ID = id
	return result, pkg.Log{}
}

func (u *AiUsecase) Verify(id string, code string) (ai.GroqAi, pkg.Log) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ai.GroqAi{}, pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "Invalid ID", "err": err.Error()},
		}
	}

	original, log := u.db.FindByID(context.Background(), oid)
	if log.Error {
		return ai.GroqAi{}, log
	}

	prompt := "Verifique se o código `" + code + "` está de acordo com `" + original.Answer + "`. Responda apenas se está de acordo ou não."

	verifyResult, log := u.repo.GenerateText(prompt)
	if log.Error {
		return ai.GroqAi{}, log
	}

	return ai.GroqAi{Answer: verifyResult}, pkg.Log{}
}
