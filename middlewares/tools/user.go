package tools

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func UserInfoCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		c.Set("userUid", strings.Replace(authHeader, "Bearer ", "", 1))
	}
}
