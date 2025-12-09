package adapter

import (
	"context"
	"miservicegolang/core/pkg"

	openai "github.com/sashabaranov/go-openai"
)

type GroqAiRepo interface {
	GenerateText(prompt string) (string, pkg.Log)
}

type GroqAi struct {
	client *openai.Client
}

func NewGroqAiRepo(key string) *GroqAi {
	config := openai.DefaultConfig(key)
	config.BaseURL = "https://api.groq.com/openai/v1"

	client := openai.NewClientWithConfig(config)

	return &GroqAi{
		client: client,
	}
}

func (g *GroqAi) GenerateText(prompt string) (string, pkg.Log) {

	req := openai.ChatCompletionRequest{
		Model: "llama-3.1-8b-instant",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	res, err := g.client.CreateChatCompletion(
		context.Background(),
		req,
	)

	if err != nil {
		return "", pkg.Log{
			Error: true,
			Body: map[string]any{
				"message": "Unable to connect.",
				"err":     err.Error(),
			},
		}
	}

	return res.Choices[0].Message.Content, pkg.Log{}
}
