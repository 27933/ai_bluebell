package mysql

import (
	"time"
)

// UpdateArticleView 更新文章访问量
func UpdateArticleView(articleId int64, date time.Time) error {
	sqlStr := `INSERT INTO article_stats (article_id, date, views)
		VALUES (?, ?, 1)
		ON DUPLICATE KEY UPDATE views = views + 1`

	_, err := db.Exec(sqlStr, articleId, date.Format("2006-01-02"))
	if err != nil {
		return err
	}

	// 同时更新文章总访问量
	_, err = db.Exec("UPDATE articles SET view_count = view_count + 1 WHERE id = ?", articleId)
	return err
}

// GetArticleDailyStats 获取文章的日访问量统计
func GetArticleDailyStats(articleId int64, days int) ([]map[string]interface{}, error) {
	var stats []map[string]interface{}

	sqlStr := `SELECT date, views
		FROM article_stats
		WHERE article_id = ? AND date >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
		ORDER BY date DESC`

	err := db.Select(&stats, sqlStr, articleId, days)
	return stats, err
}

// GetArticleTotalViews 获取文章总访问量
func GetArticleTotalViews(articleId int64) (int64, error) {
	var total int64
	sqlStr := `SELECT COALESCE(SUM(views), 0) FROM article_stats WHERE article_id = ?`
	err := db.Get(&total, sqlStr, articleId)
	return total, err
}