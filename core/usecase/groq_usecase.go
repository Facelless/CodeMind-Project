package usecase

import (
	"miservicegolang/core/pkg"
	"miservicegolang/infrastructure/repository"
)

type AiUsecase struct {
	repo repository.GroqAiRepo
}

func NewAiUsecase(r repository.GroqAiRepo) *AiUsecase {
	return &AiUsecase{repo: r}
}

func (u *AiUsecase) Generate(prompt string) (string, pkg.Log) {
	if prompt == "" {
		return "", pkg.Log{
			Error: true,
			Body: map[string]any{
				"message": "Prompt cannot be empty",
			},
		}
	}
	return u.repo.GenerateText(prompt)
}
