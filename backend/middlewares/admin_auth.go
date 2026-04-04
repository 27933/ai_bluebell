package middlewares

import (
	"bluebell/controller"
	"bluebell/models"
	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware 管理员权限验证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取当前用户角色
		role, exists := c.Get("userRole")
		if !exists {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}

		// 验证是否为管理员
		if role != string(models.UserRoleAdmin) {
			controller.ResponseError(c, controller.CodeNoPermission)
			c.Abort()
			return
		}

		c.Next()
	}
}