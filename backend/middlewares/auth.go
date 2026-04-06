package middlewares

import (
	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/pkg/jwt"
	"strings"

	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
)

// OptionalJWTMiddleware 可选JWT中间件：token 存在且合法则解析写入 context，否则直接放行
func OptionalJWTMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				mc, err := jwt.ParseToken(parts[1])
				if err == nil {
					c.Set(controller.CtxUserIDKey, mc.UserID)
				}
			}
		}
		c.Next()
	}
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		c.Set(controller.CtxUserIDKey, mc.UserID)

		// 查询用户信息，验证账号状态和角色
		user, err := mysql.GetUserById(mc.UserID)
		if err != nil {
			zap.L().Error("GetUserById failed", zap.Error(err))
			controller.ResponseError(c, controller.CodeServerBusy)
			c.Abort()
			return
		}
		// 账号被封禁，拒绝所有请求
		if user.Status != "active" {
			controller.ResponseError(c, controller.CodeUserBanned)
			c.Abort()
			return
		}
		c.Set(controller.CtxUserRoleKey, user.Role)

		c.Next()
	}
}
