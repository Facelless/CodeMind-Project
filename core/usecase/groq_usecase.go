package usecase

import (
	"context"
	"fmt"
	"miservicegolang/core/domain/ai"
	"miservicegolang/core/pkg"
	"miservicegolang/infrastructure/adapter"
	"miservicegolang/infrastructure/repository"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AiUsecase struct {
	repo    adapter.GroqAiRepo
	db      repository.GroqDatabaseRepo
	usecase UserUsecase
}

func NewAiUsecase(r adapter.GroqAiRepo, db repository.GroqDatabaseRepo, usecase UserUsecase) *AiUsecase {
	return &AiUsecase{repo: r, db: db, usecase: usecase}
}

func (u *AiUsecase) Generate(prompt string, userId string) (ai.GroqAi, pkg.Log) {
	text, log := u.repo.GenerateText(prompt)
	if log.Error {
		return ai.GroqAi{}, log
	}
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return ai.GroqAi{}, pkg.Log{
			Error: true,
			Body:  map[string]any{"message": "Invalid ID", "err": err.Error()},
		}
	}

	result := ai.GroqAi{Answer: text, Completed: false, UserId: oid}

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

	prompt := "Você é um verificador de código rigoroso. Desafio original: " + original.Answer + " Código do usuário: " + code + " Avalie se o código cumpre TODOS os requisitos do desafio. Responda apenas: - sim → se todos os requisitos foram atendidos - nao: <motivo> → se algum requisito não foi atendido Não inclua nada além dessas respostas."

	generated, log2 := u.repo.GenerateText(prompt)
	if log2.Error {
		return ai.GroqAi{}, log2
	}
	answer := strings.ToLower(generated)
	answer = strings.TrimSpace(answer)
	answer = strings.ReplaceAll(answer, `"`, "")
	fmt.Println(generated)

	if answer != "sim" && answer != "não" {
		answer = "não"
	}

	original.Verify = answer
	original.Completed = (answer == "sim")
	if original.Completed == true {
		u.usecase.SetExp(oid)
		u.db.Delete(context.Background(), oid)
	}
	u.db.UpdateVerify(context.Background(), oid, answer, original.Completed)

	return original, pkg.Log{}
}
