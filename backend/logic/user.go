package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// SignUp 用户注册
func SignUp(p *models.ParamSignUp) error {
	// 1. 检查用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	// 2. 生成用户ID
	userID := snowflake.GenID()

	// 3. 构造用户实例
	user := &models.User{
		ID:       userID,
		Username: p.Username,
		Password: p.Password,
		Role:     string(models.UserRoleReader), // 默认角色为读者
		Status:   string(models.UserStatusActive),
		Extra:    "{}", // 设置为空JSON对象，避免数据库JSON字段错误
		Avatar:   sql.NullString{Valid: false}, // 设置avatar为NULL
		Bio:      sql.NullString{Valid: false}, // 设置bio为NULL
		Email:    sql.NullString{Valid: false}, // 设置email为NULL
		Nickname: sql.NullString{Valid: false}, // 设置nickname为NULL
	}

	// 4. 保存到数据库
	return mysql.CreateUser(user)
}

// Login 用户登录
func Login(p *models.ParamLogin) (user *models.User, err error) {
	// 1. 根据用户名获取用户信息
	user = &models.User{
		Username: p.Username,
		Password: p.Password, // 保存密码用于验证
	}
	if err := mysql.Login(user); err != nil {
		return nil, err
	}

	// 2. 生成JWT Token
	token, err := jwt.GenToken(user.ID, user.Username)
	if err != nil {
		zap.L().Error("jwt.GenToken failed", zap.Error(err))
		return nil, errors.New("failed to generate token")
	}

	// 3. 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	if err := mysql.UpdateUserLoginInfo(user.ID, ""); err != nil {
		zap.L().Error("update last login time failed", zap.Error(err))
	}

	// 4. 返回用户信息（包含token）
	user.Token = token
	return user, nil
}

// GetUserByID 根据ID获取用户信息
func GetUserByID(userID int64) (*models.User, error) {
	return mysql.GetUserById(userID)
}

// GetUserByUsername 根据用户名获取用户信息
func GetUserByUsername(username string) (*models.User, error) {
	return mysql.GetUserByUsername(username)
}

// GetUserList 获取用户列表（管理员功能）
func GetUserList(role, status string, page, size int) ([]*models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	users, total, err := mysql.GetUserList(role, status, page, size)
	if err != nil {
		zap.L().Error("mysql.GetUserList() failed", zap.Error(err))
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUserStatus 更新用户状态（管理员功能）
func UpdateUserStatus(userID int64, status string) error {
	if status != string(models.UserStatusActive) && status != string(models.UserStatusInactive) {
		return errors.New("invalid status")
	}

	if err := mysql.UpdateUserStatus(userID, status); err != nil {
		zap.L().Error("mysql.UpdateUserStatus() failed", zap.Error(err))
		return err
	}

	return nil
}

// GetUserDetail 获取用户详情（管理员功能）
func GetUserDetail(userID int64) (*models.User, error) {
	user, err := mysql.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		zap.L().Error("mysql.GetUserByID() failed", zap.Error(err))
		return nil, err
	}
	return user, nil
}

// UpdateUserRole 更新用户角色（管理员功能）
func UpdateUserRole(userID int64, role string) error {
	// 验证角色是否有效
	validRoles := []string{
		string(models.UserRoleAdmin),
		string(models.UserRoleAuthor),
		string(models.UserRoleReader),
	}

	isValid := false
	for _, r := range validRoles {
		if role == r {
			isValid = true
			break
		}
	}

	if !isValid {
		return errors.New("invalid role")
	}

	if err := mysql.UpdateUserRole(userID, role); err != nil {
		zap.L().Error("mysql.UpdateUserRole() failed", zap.Error(err))
		return err
	}

	return nil
}

// BatchUpdateUserStatus 批量更新用户状态（管理员功能）
func BatchUpdateUserStatus(userIDs []int64, status string) error {
	if status != string(models.UserStatusActive) && status != string(models.UserStatusInactive) {
		return errors.New("invalid status")
	}

	if len(userIDs) == 0 {
		return errors.New("user_ids cannot be empty")
	}

	if len(userIDs) > 100 {
		return errors.New("too many users")
	}

	if err := mysql.BatchUpdateUserStatus(userIDs, status); err != nil {
		zap.L().Error("mysql.BatchUpdateUserStatus() failed", zap.Error(err))
		return err
	}

	return nil
}

// UpdateUserProfile 更新用户资料
func UpdateUserProfile(userID int64, req *models.ParamUserUpdate) error {
	updates := make(map[string]interface{})

	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Email != "" {
		// 检查邮箱是否已被使用
		if user, _ := mysql.GetUserByEmail(req.Email); user != nil && user.ID != userID {
			return errors.New("email already exists")
		}
		updates["email"] = req.Email
	}

	if len(updates) == 0 {
		return nil
	}

	if err := mysql.UpdateUser(userID, updates); err != nil {
		zap.L().Error("mysql.UpdateUser() failed", zap.Error(err))
		return err
	}

	return nil
}

// ChangePassword 修改密码
func ChangePassword(userID int64, oldPassword, newPassword string) error {
	// 1. 获取用户信息
	user, err := mysql.GetUserById(userID)
	if err != nil {
		return err
	}

	// 2. 验证旧密码
	if user.Password != oldPassword {
		return errors.New("old password is incorrect")
	}

	// 3. 更新密码
	if err := mysql.UpdateUser(userID, map[string]interface{}{
		"password": newPassword,
	}); err != nil {
		zap.L().Error("mysql.UpdateUser() failed", zap.Error(err))
		return err
	}

	return nil
}

// ResetPassword 重置密码（管理员功能）
func ResetPassword(userID int64, newPassword string) error {
	// 更新密码
	if err := mysql.UpdateUser(userID, map[string]interface{}{
		"password": newPassword,
	}); err != nil {
		zap.L().Error("mysql.UpdateUser() failed", zap.Error(err))
		return err
	}

	return nil
}

// GetUserByEmail 根据邮箱获取用户信息
func GetUserByEmail(email string) (*models.User, error) {
	return mysql.GetUserByEmail(email)
}

// GetUserByWechatOpenID 根据微信OpenID获取用户信息
func GetUserByWechatOpenID(openID string) (*models.User, error) {
	return mysql.GetUserByWechatOpenid(openID)
}

// GetUserByGithubID 根据Github ID获取用户信息
func GetUserByGithubID(githubID string) (*models.User, error) {
	return mysql.GetUserByGithubID(githubID)
}

// UpdateUserStats 更新用户统计信息
func UpdateUserStats(userID, totalWords, totalLikes int64) error {
	if err := mysql.UpdateUserStats(userID, totalWords, totalLikes); err != nil {
		zap.L().Error("mysql.UpdateUserStats() failed", zap.Error(err))
		return err
	}
	return nil
}

// BindSocialAccount 绑定社交账号
func BindSocialAccount(userID int64, provider, openID string) error {
	var field string
	switch provider {
	case "wechat":
		field = "wechat_openid"
	case "github":
		field = "github_id"
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}

	// 检查该社交账号是否已被绑定
	var existingUser *models.User
	var err error

	if provider == "wechat" {
		existingUser, err = mysql.GetUserByWechatOpenid(openID)
	} else {
		existingUser, err = mysql.GetUserByGithubID(openID)
	}

	if err == nil && existingUser != nil && existingUser.ID != userID {
		return errors.New("social account already bound to another user")
	}

	// 绑定账号
	if err := mysql.UpdateUser(userID, map[string]interface{}{
		field: openID,
	}); err != nil {
		zap.L().Error("mysql.UpdateUser() failed", zap.Error(err))
		return err
	}

	return nil
}

// UnbindSocialAccount 解绑社交账号
func UnbindSocialAccount(userID int64, provider string) error {
	var field string
	switch provider {
	case "wechat":
		field = "wechat_openid"
	case "github":
		field = "github_id"
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}

	// 解绑账号
	if err := mysql.UpdateUser(userID, map[string]interface{}{
		field: "",
	}); err != nil {
		zap.L().Error("mysql.UpdateUser() failed", zap.Error(err))
		return err
	}

	return nil
}

// GetUserProfile 获取用户资料（返回ApiUserInfo结构体）
func GetUserProfile(userID int64) (*models.ApiUserInfo, error) {
	// 获取用户信息
	user, err := GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// 转换为ApiUserInfo
	profile := &models.ApiUserInfo{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email.String,
		Role:        user.Role,
		Status:      user.Status,
		Nickname:    user.Nickname.String,
		Avatar:      user.Avatar.String,
		Bio:         user.Bio.String,
		TotalWords:  user.TotalWords,
		TotalLikes:  user.TotalLikes,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
	}

	return profile, nil
}