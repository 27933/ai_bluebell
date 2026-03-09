package models

import (
	"database/sql"
	"time"
)

// User 用户表对应结构体
type User struct {
	ID            int64          `json:"id,string" db:"id"`
	Username      string         `json:"username" db:"username"`
	Email         sql.NullString `json:"email,omitempty" db:"email"`
	Password      string         `json:"-" db:"password"`
	Role          string         `json:"role" db:"role"`
	Status        string         `json:"status" db:"status"`
	Nickname      sql.NullString `json:"nickname,omitempty" db:"nickname"`
	Avatar        sql.NullString `json:"avatar,omitempty" db:"avatar"`
	Bio           sql.NullString `json:"bio,omitempty" db:"bio"`
	TotalWords    int64          `json:"total_words" db:"total_words"`
	TotalLikes    int64          `json:"total_likes" db:"total_likes"`
	WechatOpenid  string         `json:"-" db:"wechat_openid"`
	GithubID      string         `json:"-" db:"github_id"`
	IPAddress     string         `json:"-" db:"ip_address"`
	UserAgent     string         `json:"-" db:"user_agent"`
	LastLoginAt   *time.Time     `json:"last_login_at,omitempty" db:"last_login_at"`
	Extra         string         `json:"extra,omitempty" db:"extra"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
	Token         string         `json:"token,omitempty" db:"-"` // JWT token，不存储在数据库
}

// UserRole 用户角色常量
type UserRole string

const (
	UserRoleVisitor UserRole = "visitor"
	UserRoleReader  UserRole = "reader"
	UserRoleAuthor  UserRole = "author"
	UserRoleAdmin   UserRole = "admin"
)

// UserStatus 用户状态常量
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

// ApiUserInfo 用户信息接口返回结构体
type ApiUserInfo struct {
	ID           int64      `json:"id,string"`
	Username     string     `json:"username"`
	Email        string     `json:"email,omitempty"`
	Role         string     `json:"role"`
	Status       string     `json:"status"`
	Nickname     string     `json:"nickname,omitempty"`
	Avatar       string     `json:"avatar,omitempty"`
	Bio          string     `json:"bio,omitempty"`
	TotalWords   int64      `json:"total_words"`
	TotalLikes   int64      `json:"total_likes"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ApiAuthorInfo 作者信息接口返回结构体
type ApiAuthorInfo struct {
	ID         int64  `json:"id,string"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname,omitempty"`
	Avatar     string `json:"avatar,omitempty"`
	Bio        string `json:"bio,omitempty"`
	TotalWords int64  `json:"total_words"`
	TotalLikes int64  `json:"total_likes"`
}

// ParamUserLogin 用户登录请求参数
type ParamUserLogin struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// ParamUserUpdate 更新用户信息请求参数
type ParamUserUpdate struct {
	Nickname string `json:"nickname" binding:"omitempty,min=1,max=100"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Bio      string `json:"bio" binding:"omitempty,max=500"`
	Avatar   string `json:"avatar" binding:"omitempty,max=500"`
}

// TokenResponse 登录响应结构体
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // access token过期时间（秒）
}

// RefreshTokenRequest 刷新token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// 兼容旧代码
// User 旧版本用户结构体（用于兼容）
type OldUser struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Token    string
}
