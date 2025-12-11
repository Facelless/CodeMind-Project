package routes

import (
	"miservicegolang/infrastructure/controller"
	"miservicegolang/infrastructure/server"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, ai *controller.AiController, user *controller.UserController, progress *controller.ProgressController, wsHandler *server.WebsocketHandler) {
	aiRoute := r.Group("/ai")
	{
		aiRoute.POST("/prompt", ai.Generate)
		aiRoute.POST("/verify", ai.Verify)
	}

	authRouter := r.Group("/auth")
	{
		authRouter.POST("/register", user.Register)
		authRouter.POST("/login", user.Login)
	}

	progressRouter := r.Group("/progress")
	{
		progressRouter.POST("/rank", progress.Rank)
	}

	r.GET("/ws", func(c *gin.Context) {
		wsHandler.HandleConnection(c.Writer, c.Request)
	})
}
