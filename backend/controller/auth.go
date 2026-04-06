package controller

import (
	"errors"
	"fmt"

	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoginHandler 登录处理函数（支持refresh token）
// @Summary 用户登录
// @Description 用户通过用户名和密码登录，返回用户信息和token
// @Tags 认证
// @Accept json
// @Produce json
// @Param login body models.ParamLogin true "登录参数"
// @Success 200 {object} controller._ResponseLogin "登录成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "用户名或密码错误"
// @Router /api/v1/auth/login [post]
func LoginHandler(c *gin.Context) {
	// 1. 获取参数
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 业务逻辑处理（使用新的带refresh token的登录逻辑）
	user, tokenResponse, err := logic.LoginWithRefreshToken(&models.ParamUserLogin{
		Username: p.Username,
		Password: p.Password,
	})
	if err != nil {
		zap.L().Error("logic.LoginWithRefreshToken failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserBanned) {
			ResponseError(c, CodeUserBanned)
		} else {
			ResponseError(c, CodeInvalidPassword)
		}
		return
	}

	// 3. 构建响应数据
	response := gin.H{
		"user": gin.H{
			"id":       fmt.Sprintf("%d", user.ID), // 与文章 author_id 保持一致（string 类型）
			"username": user.Username,
			"email":    user.Email.String,
			"role":     user.Role,
			"status":   user.Status,
		},
		"token": tokenResponse,
	}

	// 4. 返回响应
	ResponseSuccess(c, response)
}

// GetUserProfileHandler 获取用户信息处理函数
// @Summary 获取用户信息
// @Description 获取当前登录用户的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} controller._ResponseUserProfile "用户信息"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/auth/profile [get]
func GetUserProfileHandler(c *gin.Context) {
	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 获取用户信息
	user, err := logic.GetUserProfile(userID)
	if err != nil {
		zap.L().Error("logic.GetUserProfile failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, user)
}

// UpdateUserProfileHandler 更新用户信息处理函数
// @Summary 更新用户信息
// @Description 更新当前登录用户的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param profile body models.ParamUserUpdate true "用户信息"
// @Success 200 {object} controller._ResponseSuccess "更新成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/auth/profile [put]
func UpdateUserProfileHandler(c *gin.Context) {
	// 1. 获取参数
	p := new(models.ParamUserUpdate)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdateUserProfile with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 2. 更新用户信息
	if err := logic.UpdateUserProfile(userID, p); err != nil {
		zap.L().Error("logic.UpdateUserProfile failed", zap.Error(err))
		if err.Error() == "email already exists" {
			ResponseErrorWithMsg(c, CodeInvalidParam, "该邮箱已被其他账号使用")
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// RefreshTokenHandler 刷新access token处理函数
// @Summary 刷新Access Token
// @Description 使用Refresh Token获取新的Access Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param refresh body models.RefreshTokenRequest true "刷新token请求"
// @Success 200 {object} controller._ResponseRefreshToken "新的token信息"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "token无效"
// @Router /api/v1/auth/refresh [post]
func RefreshTokenHandler(c *gin.Context) {
	// 1. 获取参数
	p := new(models.RefreshTokenRequest)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("RefreshToken with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 业务逻辑处理
	tokenResponse, err := logic.RefreshAccessToken(p.RefreshToken)
	if err != nil {
		zap.L().Error("logic.RefreshAccessToken failed", zap.Error(err))
		ResponseError(c, CodeInvalidToken)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, tokenResponse)
}