package models

import "time"

// Tag 标签表对应结构体
type Tag struct {
	ID          int64     `json:"id,string" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	Slug        string    `json:"slug,omitempty" db:"slug"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ApiTag 标签接口返回结构体
type ApiTag struct {
	ID           int64     `json:"id,string"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Slug         string    `json:"slug,omitempty"`
	ArticleCount int       `json:"article_count"`
	CreatedAt    time.Time `json:"created_at"`
}

// ParamCreateTag 创建标签请求参数
type ParamCreateTag struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"omitempty,max=200"`
	Slug        string `json:"slug" binding:"omitempty,max=50"`
}

// ParamUpdateTag 更新标签请求参数
type ParamUpdateTag struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=50"`
	Description string `json:"description" binding:"omitempty,max=200"`
	Slug        string `json:"slug" binding:"omitempty,max=50"`
}