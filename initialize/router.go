package initialize

import (
	"flaver/api"
	"flaver/globals"
	"flaver/middlewares"
	"flaver/routers"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func Routers() *gin.Engine {

	var rootRouter = gin.Default()
	// 跨域
	rootRouter.Use(middlewares.Cors())
	globals.GetLogger().Info("use middleware cors")
	// globals.GetLogger().Info("register swagger handler")

	rootRouter.Use(middlewares.Secure())
	rootRouter.Use(middlewares.ApiVerify)

	m := ginmetrics.GetMonitor()

	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	// set middleware for gin, path /metrics
	m.Use(rootRouter)

	rootRouter.GET("healthy", healthyRequest)

	routeGroup := routers.RouterGroupApp
	v1Group := rootRouter.Group("v1")
	{
		routeGroup.UserGroup.InitRouter(v1Group, rootRouter)
		routeGroup.PostGroup.InitRouter(v1Group, rootRouter)
		routeGroup.FileGroup.InitRouter(v1Group, rootRouter)
		routeGroup.TopicGroup.InitRouter(v1Group, rootRouter)
	}

	globals.GetLogger().Info("router register success")
	return rootRouter
}

func healthyRequest(c *gin.Context) {
	data := map[string]interface{}{
		"healthy": "ok",
		"code":    200,
		"content": "ok",
	}
	api.SendResult(nil, data, c)
}
