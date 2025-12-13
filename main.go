package main

import (
	"log"
	"miservicegolang/core/usecase"
	"miservicegolang/infrastructure/adapter"
	"miservicegolang/infrastructure/controller"
	"miservicegolang/infrastructure/database"
	"miservicegolang/infrastructure/repository"
	"miservicegolang/infrastructure/server"
	"miservicegolang/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	client, err := database.ConnectMongodb()
	if err.Error {
		log.Fatal(err)
	}

	userRepo := repository.NewUserDatabaseRepo(client)
	groqDatabaseRepo := repository.NewGroqAiDatabaseRepo(client)
	groqRepo := adapter.NewGroqAiRepo(os.Getenv("GROQ_KEY"))

	userUsecase := usecase.NewUserUsecase(userRepo)
	growthUsecase := usecase.NewGrowthUsecase(userRepo)
	aiUsecase := usecase.NewAiUsecase(groqRepo, groqDatabaseRepo, growthUsecase)

	aiController := controller.NewAiController(aiUsecase)
	userController := controller.NewUserController(userUsecase)
	progressController := controller.NewProgressController(growthUsecase)

	hubUc := usecase.NewHub()
	matchUC := usecase.NewMatchUsecase(*aiUsecase, hubUc)
	queueUC := usecase.NewQueueUsecase(userRepo, matchUC, hubUc)

	go queueUC.QueueMaker()
	go matchUC.MatchMaker()
	wsHandler := server.NewWebsocketHandler(queueUC, hubUc	)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	routes.SetupRoutes(r, aiController, userController, progressController, wsHandler)

	log.Println("Servidor rodando em :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
