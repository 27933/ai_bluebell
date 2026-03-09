package models

import "time"

// Like 点赞表对应结构体
type Like struct {
	ID         int64     `json:"id,string" db:"id"`
	UserID     int64     `json:"user_id,string" db:"user_id"`
	TargetType string    `json:"target_type" db:"target_type"`
	TargetID   int64     `json:"target_id,string" db:"target_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// TargetType 点赞目标类型常量
type TargetType string

const (
	TargetTypeArticle TargetType = "article"
	TargetTypeComment TargetType = "comment"
)

// ApiLikeStatus 点赞状态接口返回
type ApiLikeStatus struct {
	IsLiked   bool      `json:"is_liked"`
	LikeCount int       `json:"like_count"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}