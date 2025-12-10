package main

import (
	"fmt"
	"miservicegolang/core/usecase"
	"miservicegolang/infrastructure/adapter"
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
	groqRepo := adapter.NewGroqAiRepo(os.Getenv("GROQ_KEY"))
	groqDatabaseRepo := repository.NewGroqAiDatabaseRepo(client)
	aiUsecase := usecase.NewAiUsecase(groqRepo, groqDatabaseRepo)
	aiController := controller.NewAiController(aiUsecase)
	userDatabse := repository.NewUserDatabaseRepo(client)
	userUsecase := usecase.NewUserUsecase(userDatabse)
	useController := controller.NewUserController(userUsecase)
	fmt.Println(log)
	r := gin.Default()
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
	routes.SetupRoutes(r, aiController, useController)
	r.Run(":3000")
}
