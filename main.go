package main

import (
	"miservicegolang/core/usecase"
	"miservicegolang/infrastructure/controller"
	"miservicegolang/infrastructure/file"
	"miservicegolang/infrastructure/repository"
	"miservicegolang/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	groqRepo := repository.NewGroqAiRepo(os.Getenv("GROQ_KEY"))
	aiUsecase := usecase.NewAiUsecase(groqRepo)
	aiController := controller.NewAiController(aiUsecase)

	fileService := file.NewLocalFileService()
	unzipService := file.NewLocalUnzipService()
	uploadUc := usecase.NewUploadUsecase(fileService, unzipService)
	upload := controller.NewUploadController(uploadUc)
	r := gin.Default()
	routes.SetupRoutes(r, aiController, upload)
	r.Run(":8080")
}
