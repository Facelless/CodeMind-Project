package routes

import (
	"miservicegolang/infrastructure/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, ai *controller.AiController) {
	aiRoute := r.Group("/ai")
	{
		aiRoute.POST("/prompt", ai.Generate)
		aiRoute.POST("/verify", ai.Verify)
	}
}
