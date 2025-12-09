package usecase

import (
	"context"
	"miservicegolang/core/domain/ai"
	"miservicegolang/core/pkg"
	"miservicegolang/infrastructure/repository"
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
			Body: map[string]any{
				"message": "Prompt cannot be empty",
			},
		}
	}
	text, log := u.repo.GenerateText(prompt)
	if log.Error {
		return ai.GroqAi{}, log
	}

	result := ai.GroqAi{
		Prompt: prompt,
		Answer: text,
	}

	id, log := u.db.Insert(context.Background(), result)
	if log.Error {
		return ai.GroqAi{}, log
	}
	result.ID = id
	return result, pkg.Log{}

}
