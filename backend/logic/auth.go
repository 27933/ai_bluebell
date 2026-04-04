package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"errors"
	"time"

	"go.uber.org/zap"
)

// LoginWithRefreshToken 用户登录并返回refresh token
func LoginWithRefreshToken(p *models.ParamUserLogin) (*models.User, *models.TokenResponse, error) {
	// 1. 根据用户名获取用户信息
	user := &models.User{
		Username: p.Username,
		Password: p.Password, // 保存密码用于验证
	}
	if err := mysql.Login(user); err != nil {
		return nil, nil, err
	}

	// 2. 生成access token和refresh token
	accessToken, err := jwt.GenToken(user.ID, user.Username)
	if err != nil {
		zap.L().Error("jwt.GenToken failed", zap.Error(err))
		return nil, nil, errors.New("failed to generate access token")
	}

	refreshToken, err := jwt.GenRefreshToken(user.ID)
	if err != nil {
		zap.L().Error("jwt.GenRefreshToken failed", zap.Error(err))
		return nil, nil, errors.New("failed to generate refresh token")
	}

	// 3. 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	if err := mysql.UpdateUserLoginInfo(user.ID, ""); err != nil {
		zap.L().Error("update last login time failed", zap.Error(err))
	}

	// 4. 返回用户信息和token
	tokenResponse := &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    28800, // 8小时，前端可以根据这个提前刷新token
	}

	return user, tokenResponse, nil
}

// RefreshAccessToken 使用refresh token获取新的access token
func RefreshAccessToken(refreshToken string) (*models.TokenResponse, error) {
	// 1. 解析refresh token
	refreshClaims, err := jwt.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 2. 检查用户是否存在且状态正常
	user, err := mysql.GetUserById(refreshClaims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Status != string(models.UserStatusActive) {
		return nil, errors.New("user account is inactive")
	}

	// 3. 生成新的access token
	accessToken, err := jwt.GenToken(user.ID, user.Username)
	if err != nil {
		zap.L().Error("jwt.GenToken failed", zap.Error(err))
		return nil, errors.New("failed to generate access token")
	}

	// 4. 返回新的token（refresh token可以继续使用，也可以生成新的）
	tokenResponse := &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // 可以复用refresh token
		ExpiresIn:    28800, // 8小时
	}

	return tokenResponse, nil
}