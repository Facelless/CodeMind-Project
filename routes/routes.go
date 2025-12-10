package routes

import (
	"miservicegolang/infrastructure/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, ai *controller.AiController, user *controller.UserController) {
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
}
