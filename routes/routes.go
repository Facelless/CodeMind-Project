package routes

import (
	"miservicegolang/infrastructure/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, ai *controller.AiController, upload *controller.UploadController) {
	aiRoute := r.Group("/ai")
	{
		aiRoute.POST("/prompt", ai.Generate)
	}
	fileRoute := r.Group("/file")
	{
		fileRoute.POST("/upload", upload.UploadFile)
	}
}
