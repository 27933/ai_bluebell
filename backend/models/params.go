package models

import (
	"database/sql"
	"errors"
)

var (
	ErrorNoPermission    = errors.New("没有权限")
	ErrorArticleNotExist = errors.New("文章不存在")
)

// 定义请求的参数结构体

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required,min=3,max=50"`
	Password   string `json:"password" binding:"required,min=8,max=100"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	// UserID 从请求中获取当前的用户
	PostID    string `json:"post_id" binding:"required"`               // 贴子id
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1" ` // 赞成票(1)还是反对票(-1)取消投票(0)
}

// ParamPostList 获取帖子列表query string参数
type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"`   // 可以为空
	Page        int64  `json:"page" form:"page" example:"1"`       // 页码
	Size        int64  `json:"size" form:"size" example:"10"`      // 每页数据量
	Order       string `json:"order" form:"order" example:"score"` // 排序依据
}

// TargetInfo 批量点赞查询目标信息
type TargetInfo struct {
	TargetType string `json:"target_type" binding:"required,oneof=article comment"`
	TargetID   int64  `json:"target_id,string" binding:"required"`
}

// ParamBatchLikeStatus 批量获取点赞状态请求参数
type ParamBatchLikeStatus struct {
	Targets []TargetInfo `json:"targets" binding:"required,min=1,max=50,dive"`
}

// ParamAuthorArticlesList 作者文章列表请求参数
type ParamAuthorArticlesList struct {
	Username string `json:"username" binding:"required" form:"username"`
	Page     int    `json:"page" form:"page" example:"1"`   // 页码
	Size     int    `json:"size" form:"size" example:"20"`  // 每页数量
	Sort     string `json:"sort" form:"sort" example:"time"` // 排序方式：time(时间)/hot(热度)
}

// AuthorInfoResponse 作者信息响应
type AuthorInfoResponse struct {
	Username     string         `json:"username" db:"username"`
	Nickname     sql.NullString `json:"nickname" db:"nickname"`
	Bio          sql.NullString `json:"bio" db:"bio"`
	JoinDate     string         `json:"join_date" db:"join_date"`
	ArticleCount int64          `json:"article_count" db:"article_count"`
	TotalViews   int64          `json:"total_views" db:"total_views"`
	TotalLikes   int64          `json:"total_likes" db:"total_likes"`
}

// ExportArticlesBatchRequest 批量导出文章请求
type ExportArticlesBatchRequest struct {
	ArticleIDs []int64 `json:"article_ids" binding:"required,min=1,max=50"`
}

// ExportArticleResponse 单篇文章导出响应
type ExportArticleResponse struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
	Size     int64  `json:"size"`
}

// ExportFileInfo 导出文件信息
type ExportFileInfo struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
	Size     int64  `json:"size"`
}

// ExportBatchResponse 批量导出响应
type ExportBatchResponse struct {
	BatchID   string            `json:"batch_id"`
	Files     []*ExportFileInfo `json:"files"`
	FileCount int64             `json:"file_count"`
	TotalSize int64             `json:"total_size"`
}

// BatchLikeStatusResponse 批量点赞状态响应结果
type BatchLikeStatusResponse struct {
	TargetType string `json:"target_type"`
	TargetID   int64  `json:"target_id,string"`
	IsLiked    *bool  `json:"is_liked"` // 使用指针类型，null表示未查询（目标不存在等）
	LikeCount  int    `json:"like_count"`
}

// ParamLikeHistory 点赞历史查询参数
type ParamLikeHistory struct {
	TargetType   *string `json:"target_type" form:"target_type" binding:"omitempty,oneof=article comment"` // 目标类型，可选
	StartTime    *string `json:"start_time" form:"start_time"`                                             // 开始时间，可选
	EndTime      *string `json:"end_time" form:"end_time"`                                                 // 结束时间，可选
	Page         int64   `json:"page" form:"page" example:"1"`                                             // 页码
	Size         int64   `json:"size" form:"size" example:"20"`                                            // 每页大小
}
