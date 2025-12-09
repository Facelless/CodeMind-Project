package usecase

import (
	"context"
	"miservicegolang/core/domain/ai"
	"miservicegolang/core/pkg"
	"miservicegolang/infrastructure/adapter"
	"miservicegolang/infrastructure/repository"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AiUsecase struct {
	repo adapter.GroqAiRepo
	db   repository.GroqDatabaseRepo
}

func NewAiUsecase(r adapter.GroqAiRepo, db repository.GroqDatabaseRepo) *AiUsecase {
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

	result := ai.GroqAi{Answer: text, Completed: false}

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

	prompt := `
	Você é um verificador de código extremamente preciso.

	Desafio original:
	"` + original.Answer + `"

	Código enviado pelo usuário:
	"` + code + `"

	Analise se o código atende COMPLETAMENTE todos os requisitos.
	Ignore estilo, comentários e pequenas diferenças irrelevantes.

	Responda APENAS com:
	sim
	ou
	não
	`
	generated, log2 := u.Generate(prompt)
	if log2.Error {
		return ai.GroqAi{}, log2
	}
	answer := strings.ToLower(generated.Answer)
	answer = strings.TrimSpace(answer)
	answer = strings.ReplaceAll(answer, `"`, "")

	if answer != "sim" && answer != "não" {
		answer = "não"
	}

	original.Verify = answer
	original.Completed = (answer == "sim")

	u.db.UpdateVerify(context.Background(), oid, answer, original.Completed)

	return original, pkg.Log{}
}
