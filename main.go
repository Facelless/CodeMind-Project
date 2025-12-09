package main

import (
	"fmt"
	"miservicegolang/core/usecase"
	"miservicegolang/infrastructure/controller"
	"miservicegolang/infrastructure/database"
	"miservicegolang/infrastructure/file"
	"miservicegolang/infrastructure/repository"
	"miservicegolang/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	client, log := database.ConnectMongodb()
	groqRepo := repository.NewGroqAiRepo(os.Getenv("GROQ_KEY"))
	groqDatabaseRepo := usecase.NewGroqAiDatabaseRepo(client)
	aiUsecase := usecase.NewAiUsecase(groqRepo, groqDatabaseRepo)
	aiController := controller.NewAiController(aiUsecase)
	fmt.Println(log)
	fileService := file.NewLocalFileService()
	unzipService := file.NewLocalUnzipService()
	uploadUc := usecase.NewUploadUsecase(fileService, unzipService)
	upload := controller.NewUploadController(uploadUc)
	r := gin.Default()
	routes.SetupRoutes(r, aiController, upload)
	r.Run(":8080")
}
