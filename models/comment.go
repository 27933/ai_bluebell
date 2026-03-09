package models

import "time"

// Comment 评论表对应结构体
type Comment struct {
	ID         int64     `json:"id,string" db:"id"`
	ArticleID  int64     `json:"article_id,string" db:"article_id"`
	UserID     int64     `json:"user_id,string" db:"user_id"`
	ParentID   *int64    `json:"parent_id,omitempty,string" db:"parent_id"`
	Content    string    `json:"content" db:"content"`
	LikeCount  int       `json:"like_count" db:"like_count"`
	Status     string    `json:"status" db:"status"`
	IPAddress  string    `json:"-" db:"ip_address"`
	UserAgent  string    `json:"-" db:"user_agent"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// CommentStatus 评论状态常量
type CommentStatus string

const (
	CommentStatusActive  CommentStatus = "active"
	CommentStatusDeleted CommentStatus = "deleted"
)

// ApiComment 评论接口返回结构体
type ApiComment struct {
	ID        int64         `json:"id,string"`
	Content   string        `json:"content"`
	LikeCount int           `json:"like_count"`
	Status    string        `json:"status"`
	Author    CommentAuthor `json:"author"`
	ParentID  *int64        `json:"parent_id,omitempty,string"`
	CreatedAt time.Time     `json:"created_at"`
}

// CommentAuthor 评论作者信息
type CommentAuthor struct {
	ID       int64  `json:"id,string"`
	Username string `json:"username"`
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

// ParamCreateComment 创建评论请求参数
type ParamCreateComment struct {
	ArticleID int64  `json:"article_id" binding:"required"`
	ParentID  *int64 `json:"parent_id,omitempty"`
	Content   string `json:"content" binding:"required,min=1,max=1000"`
	UserID    int64  `json:"-"` // 从JWT token获取，不通过JSON传递
}

// ParamCommentList 评论列表查询参数
type ParamCommentList struct {
	ArticleID int64 `form:"article_id" binding:"required"`
	Page      int   `form:"page" binding:"omitempty,min=1"`
	Size      int   `form:"size" binding:"omitempty,min=1,max=50"`
}