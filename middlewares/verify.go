package middlewares

import (
	"flaver/middlewares/tools"
	"flaver/api"

	"github.com/gin-gonic/gin"
)

// Middleware API Verify:
// Handle Request Without Add Header
func ApiVerify(ctx *gin.Context) {
	heandler := tools.GetApiVerifyHandler(ctx.Request)
	if err := heandler.ProcessRequest(ctx.Request); err != nil {
		api.SendResult(err, nil, ctx)
		ctx.Abort()
		return
	}
	ctx.Next()
}
