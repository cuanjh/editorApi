package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var RequestId = uuid.NewV4().String()
		c.Writer.Header().Set("X-Request-Id", RequestId)
		c.Request.Header.Set("X-Request-Id", RequestId)
		c.Next()
	}
}
