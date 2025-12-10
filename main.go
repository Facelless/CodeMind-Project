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
	userDatabse := repository.NewUserDatabaseRepo(client)
	userUsecase := usecase.NewUserUsecase(userDatabse)
	useController := controller.NewUserController(userUsecase)
	growth := usecase.NewGrowthUsecase(userDatabse)
	aiUsecase := usecase.NewAiUsecase(groqRepo, groqDatabaseRepo, growth)
	aiController := controller.NewAiController(aiUsecase)
	progressController := controller.NewProgressController(growth)
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
	routes.SetupRoutes(r, aiController, useController, progressController)
	r.Run(":3000")
}
