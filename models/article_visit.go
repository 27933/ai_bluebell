package models

import "time"

// ArticleVisit 文章访问记录模型
type ArticleVisit struct {
	ID         int64     `db:"id" json:"id"`
	ArticleID  int64     `db:"article_id" json:"article_id"`
	UserID     *int64    `db:"user_id" json:"user_id,omitempty"` // 登录用户才有
	IPAddress  string    `db:"ip_address" json:"ip_address"`
	VisitDate  time.Time `db:"visit_date" json:"visit_date"`
	VisitHour  int       `db:"visit_hour" json:"visit_hour"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// TrendingArticle 热门文章排行模型
type TrendingArticle struct {
	ID               int64     `db:"id" json:"id"`
	ArticleID        int64     `db:"article_id" json:"article_id"`
	PeriodType       string    `db:"period_type" json:"period_type"` // daily, weekly, monthly
	PeriodDate       time.Time `db:"period_date" json:"period_date"`
	ViewCount        int       `db:"view_count" json:"view_count"`
	UniqueViewCount  int       `db:"unique_view_count" json:"unique_view_count"`
	RankPosition     int       `db:"rank_position" json:"rank_position"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

// ArticleVisitParams 创建访问记录参数
type ArticleVisitParams struct {
	ArticleID int64  `json:"article_id" binding:"required"`
	UserID    *int64 `json:"user_id,omitempty"`
	IPAddress string `json:"ip_address" binding:"required"`
	VisitDate string `json:"visit_date" binding:"required"`
	VisitHour int    `json:"visit_hour" binding:"required,min=0,max=23"`
}

// TrendingArticleQuery 热门文章查询参数
type TrendingArticleQuery struct {
	PeriodType string `form:"period" binding:"required,oneof=daily weekly monthly"`
	Page       int    `form:"page" binding:"min=1"`
	Size       int    `form:"size" binding:"min=1,max=100"`
}

// ArticleStatsQuery 文章统计查询参数
type ArticleStatsQuery struct {
	ArticleID int64  `form:"article_id" binding:"required"`
	Days      int    `form:"days" binding:"min=1,max=90"`
	GroupBy   string `form:"group_by" binding:"required,oneof=hour day week month"`
}

// ArticleStatsBatchQuery 批量统计查询参数
type ArticleStatsBatchQuery struct {
	IDs    string `form:"ids" binding:"required"` // 逗号分隔的文章ID
	Fields string `form:"fields" binding:"required"` // views,uv,trend
}

// ArticleTrendData 文章趋势数据
type ArticleTrendData struct {
	Label string `json:"label"` // 时间标签，如"2026-02-27 14:00"
	Value int    `json:"value"` // 访问量
}

// ArticleStatsData 文章统计数据
type ArticleStatsData struct {
	ArticleID       int64              `json:"article_id"`
	TotalViews      int                `json:"total_views"`
	TotalUV         int                `json:"total_uv"`
	Trend           []ArticleTrendData `json:"trend,omitempty"`
	DailyStats      map[string]int     `json:"daily_stats,omitempty"` // 日期->访问量
}

// TrendingArticleData 热门文章数据（包含文章信息）
type TrendingArticleData struct {
	ID               int64     `db:"id" json:"id"`
	ArticleID        int64     `db:"article_id" json:"article_id"`
	PeriodType       string    `db:"period_type" json:"period_type"`
	PeriodDate       time.Time `db:"period_date" json:"period_date"`
	ViewCount        int       `db:"view_count" json:"view_count"`
	UniqueViewCount  int       `db:"unique_view_count" json:"unique_view_count"`
	RankPosition     int       `db:"rank_position" json:"rank_position"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
	Title            string    `db:"title" json:"title"`
	Summary          string    `db:"summary" json:"summary"`
	AuthorID         int64     `db:"author_id" json:"author_id"`
	Status           string    `db:"status" json:"status"`
	IsFeatured       bool      `db:"is_featured" json:"is_featured"`
	LikeCount        int       `db:"like_count" json:"like_count"`
	CommentCount     int       `db:"comment_count" json:"comment_count"`
	ArticleViews     int       `db:"article_views" json:"article_views"`
	ArticleCreatedAt time.Time `db:"article_created_at" json:"article_created_at"`
	ArticleUpdatedAt time.Time `db:"article_updated_at" json:"article_updated_at"`
	AuthorUsername   string    `db:"author_username" json:"author_username"`
	AuthorNickname   string    `db:"author_nickname" json:"author_nickname"`
}

// VisitStats 访问统计数据
type VisitStats struct {
	TotalViews     int            `json:"total_views"`
	TotalUV        int            `json:"total_uv"`
	DailyViews     map[string]int `json:"daily_views"`     // 日期->浏览量
	DailyUV        map[string]int `json:"daily_uv"`        // 日期->独立访客
	HourlyDistribution [24]int    `json:"hourly_distribution"` // 24小时访问量分布
}