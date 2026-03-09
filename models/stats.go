package models

// SystemOverview 系统概览数据结构
type SystemOverview struct {
	UserCount              int64 `json:"user_count"`
	ArticleCount           int64 `json:"article_count"`
	CommentCount           int64 `json:"comment_count"`
	TodayNewUserCount      int64 `json:"today_new_user_count"`
	TodayNewArticleCount   int64 `json:"today_new_article_count"`
}

// DailyStats 日统计数据结构
type DailyStats struct {
	Date           string `json:"date"`
	NewUserCount   int64  `json:"new_user_count"`
	NewArticleCount int64 `json:"new_article_count"`
	NewCommentCount int64 `json:"new_comment_count"`
}

// ParamStatsDaily 日统计查询参数
type ParamStatsDaily struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}