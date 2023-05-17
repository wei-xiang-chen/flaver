package routers

import (
	"flaver/controllers"
	middleware "flaver/middlewares/tools"

	"github.com/gin-gonic/gin"
)

type PostRouter struct{}

func (PostRouter) InitRouter(router *gin.RouterGroup, g *gin.Engine) {
	postController := controllers.NewPostController()
	routeGroup := router.Group("posts", middleware.UserInfoCheck())
	routeGroup.GET("/", postController.GetList)
	routeGroup.GET("/:id", postController.GetDetail)
	routeGroup.POST("/", postController.Create)
	routeGroup.PATCH("/:id", postController.Update)
	routeGroup.POST("/like/:id", postController.Like)
}
