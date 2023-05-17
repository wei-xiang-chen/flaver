package routers

import (
	"flaver/controllers"
	middleware "flaver/middlewares/tools"

	"github.com/gin-gonic/gin"
)

type FileRouter struct{}

func (FileRouter) InitRouter(router *gin.RouterGroup, g *gin.Engine) {
	fileController := controllers.NewFileController()
	routeGroup := router.Group("files", middleware.UserInfoCheck())
	routeGroup.POST("/", fileController.UploadImage)
}
