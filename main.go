package main

import (
	"fmt"
	"miservicegolang/core/usecase"
	"miservicegolang/infrastructure/controller"
	"miservicegolang/infrastructure/database"
	"miservicegolang/infrastructure/repository"
	"miservicegolang/routes"
	"net/http"
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

	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	routes.SetupRoutes(r, aiController)
	r.Run(":3000")
}
