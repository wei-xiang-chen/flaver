package routers

import (
	"flaver/controllers"
	middleware "flaver/middlewares/tools"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (UserRouter) InitRouter(router *gin.RouterGroup, g *gin.Engine) {
	userController := controllers.NewUserController()
	routeGroup := router.Group("users")
	routeGroup.POST("/login", userController.Login)
	routeGroup.POST("/refresh-token", userController.RefreshToken)
	routeGroup.POST("/register", userController.RegistUser)

	routeGroup.GET("/:userUid", middleware.UserInfoCheck(), userController.GetDetail)
	routeGroup.PATCH("/", middleware.UserInfoCheck(), userController.Update)
}
