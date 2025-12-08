package main

import (
	"miservicegolang/core/usecase"
	"miservicegolang/infrastructure/controller"
	"miservicegolang/infrastructure/repository"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	groqRepo := repository.NewGroqAiRepo(os.Getenv("GROQ_KEY"))
	aiUsecase := usecase.NewAiUsecase(groqRepo)
	aiController := controller.NewAiController(aiUsecase)
	r := gin.Default()
	r.POST("/ai", aiController.Generate)
	r.Run(":8080")
}
