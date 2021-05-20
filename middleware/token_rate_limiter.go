package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

var limiter = ratelimit.NewBucketWithQuantum(time.Second, 10, 300) //限流并发 3000 QPS

func TokenRateLimiter() gin.HandlerFunc {
	fmt.Println("token create rate:", limiter.Rate())
	fmt.Println("available token :", limiter.Available())
	return func(context *gin.Context) {
		if limiter.TakeAvailable(1) == 0 {
			log.Printf("available token :%d", limiter.Available())
			context.AbortWithStatusJSON(http.StatusTooManyRequests, "Too Many Request")
		} else {
			context.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Available()))
			context.Writer.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Capacity()))
			context.Next()
		}
	}
}