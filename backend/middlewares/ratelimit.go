package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// IPRateLimitMiddleware 按 IP 的令牌桶限速
// fillInterval: 每隔多久补充一个令牌
// cap: 桶容量（最大突发请求数）
func IPRateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	var limiters sync.Map
	return func(c *gin.Context) {
		ip := c.ClientIP()
		v, _ := limiters.LoadOrStore(ip, ratelimit.NewBucket(fillInterval, cap))
		bucket := v.(*ratelimit.Bucket)
		if bucket.TakeAvailable(1) != 1 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
