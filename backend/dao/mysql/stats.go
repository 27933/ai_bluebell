package mysql

import (
	"bluebell/models"
	"fmt"
	"time"
)

// GetUserCount 获取用户总数
func GetUserCount() (int64, error) {
	var count int64
	sqlStr := `SELECT COUNT(*) FROM users WHERE status = 'active'`
	err := db.Get(&count, sqlStr)
	return count, err
}

// GetArticleCount 获取文章总数
func GetArticleCount() (int64, error) {
	var count int64
	sqlStr := `SELECT COUNT(*) FROM articles WHERE status = 'published'`
	err := db.Get(&count, sqlStr)
	return count, err
}

// GetCommentCount 获取评论总数
func GetCommentCount() (int64, error) {
	var count int64
	sqlStr := `SELECT COUNT(*) FROM comments WHERE status = 'active'`
	err := db.Get(&count, sqlStr)
	return count, err
}

// GetUserCountByDate 获取指定日期后的新增用户数
func GetUserCountByDate(startDate string) (int64, error) {
	var count int64
	sqlStr := `SELECT COUNT(*) FROM users WHERE created_at >= ? AND status = 'active'`
	err := db.Get(&count, sqlStr, startDate)
	return count, err
}

// GetArticleCountByDate 获取指定日期后的新增文章数
func GetArticleCountByDate(startDate string) (int64, error) {
	var count int64
	sqlStr := `SELECT COUNT(*) FROM articles WHERE created_at >= ? AND status = 'published'`
	err := db.Get(&count, sqlStr, startDate)
	return count, err
}

// GetDailyStats 获取日统计数据
func GetDailyStats(startDate, endDate string) ([]*models.DailyStats, error) {
	var stats []*models.DailyStats

	// 如果没有提供日期范围，默认获取最近30天的数据
	if startDate == "" || endDate == "" {
		endDate = time.Now().Format("2006-01-02")
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}

	sqlStr := fmt.Sprintf(`
		SELECT
			DATE(created_at) as date,
			SUM(CASE WHEN type = 'user' THEN 1 ELSE 0 END) as new_user_count,
			SUM(CASE WHEN type = 'article' THEN 1 ELSE 0 END) as new_article_count,
			SUM(CASE WHEN type = 'comment' THEN 1 ELSE 0 END) as new_comment_count
		FROM (
			SELECT created_at, 'user' as type FROM users WHERE created_at >= '%s' AND created_at <= '%s 23:59:59'
			UNION ALL
			SELECT created_at, 'article' as type FROM articles WHERE created_at >= '%s' AND created_at <= '%s 23:59:59'
			UNION ALL
			SELECT created_at, 'comment' as type FROM comments WHERE created_at >= '%s' AND created_at <= '%s 23:59:59'
		) combined
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`, startDate, endDate, startDate, endDate, startDate, endDate)

	err := db.Select(&stats, sqlStr)
	return stats, err
}