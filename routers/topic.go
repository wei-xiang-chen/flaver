package routers

import (
	"flaver/controllers"

	"github.com/gin-gonic/gin"
)

type TopicRouter struct{}

func (TopicRouter) InitRouter(router *gin.RouterGroup, g *gin.Engine) {
	topicController := controllers.NewTopicController()
	routeGroup := router.Group("topics")
	routeGroup.GET("/", topicController.GetList)
}
