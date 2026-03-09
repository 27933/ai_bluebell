package mysql

import (
	"fmt"
	"time"

	"bluebell/models"
	"bluebell/pkg/snowflake"

	"github.com/jmoiron/sqlx"
)

// RecordArticleVisit 记录文章访问（带防刷机制）
func RecordArticleVisit(articleID int64, userID *int64, ipAddress string, visitTime time.Time) error {
	// 构建访问日期和小时
	visitDate := visitTime.Format("2006-01-02")
	visitHour := visitTime.Hour()

	// 检查是否已存在访问记录（用于去重）
	var exists bool
	if userID != nil {
		// 登录用户：检查user_id + article_id + visit_date
		query := `SELECT EXISTS(SELECT 1 FROM article_visits WHERE user_id = ? AND article_id = ? AND visit_date = ?)`
		err := db.Get(&exists, query, *userID, articleID, visitDate)
		if err != nil {
			return fmt.Errorf("check user visit exists failed: %w", err)
		}
	} else {
		// 未登录用户：检查ip_address + article_id + visit_date
		query := `SELECT EXISTS(SELECT 1 FROM article_visits WHERE ip_address = ? AND article_id = ? AND visit_date = ?)`
		err := db.Get(&exists, query, ipAddress, articleID, visitDate)
		if err != nil {
			return fmt.Errorf("check ip visit exists failed: %w", err)
		}
	}

	// 如果已存在，不重复记录
	if exists {
		return nil
	}

	// 生成雪花ID
	visitID := snowflake.GenID()

	// 插入访问记录
	query := `INSERT INTO article_visits (id, article_id, user_id, ip_address, visit_date, visit_hour, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, visitID, articleID, userID, ipAddress, visitDate, visitHour, time.Now())
	if err != nil {
		return fmt.Errorf("insert article visit failed: %w", err)
	}

	return nil
}

// GetArticleUniqueViews 获取文章独立访客数
func GetArticleUniqueViews(articleID int64, startDate, endDate time.Time) (int, error) {
	var uv int
	query := `SELECT COUNT(DISTINCT ip_address) FROM article_visits
			  WHERE article_id = ? AND visit_date >= ? AND visit_date <= ?`
	err := db.Get(&uv, query, articleID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return 0, fmt.Errorf("get article unique views failed: %w", err)
	}
	return uv, nil
}

// GetArticleHourlyStats 获取文章小时级别统计数据
func GetArticleHourlyStats(articleID int64, date time.Time) ([]models.ArticleTrendData, error) {
	var stats []models.ArticleTrendData

	query := `SELECT
			  CONCAT(visit_date, ' ', LPAD(visit_hour, 2, '0'), ':00') as label,
			  COUNT(*) as value
			  FROM article_visits
			  WHERE article_id = ? AND visit_date = ?
			  GROUP BY visit_date, visit_hour
			  ORDER BY visit_date, visit_hour`

	err := db.Select(&stats, query, articleID, date.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("get article hourly stats failed: %w", err)
	}

	return stats, nil
}

// GetArticleDailyStatsWithUV 获取文章每日统计数据（包含UV）
func GetArticleDailyStatsWithUV(articleID int64, startDate, endDate time.Time) (map[string]models.ArticleStatsData, error) {
	stats := make(map[string]models.ArticleStatsData)

	// 获取每日PV和UV
	query := `SELECT
			  visit_date,
			  COUNT(*) as pv,
			  COUNT(DISTINCT ip_address) as uv
			  FROM article_visits
			  WHERE article_id = ? AND visit_date >= ? AND visit_date <= ?
			  GROUP BY visit_date
			  ORDER BY visit_date`

	type dailyStat struct {
		VisitDate string `db:"visit_date"`
		PV        int    `db:"pv"`
		UV        int    `db:"uv"`
	}

	var rows []dailyStat
	err := db.Select(&rows, query, articleID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("get article daily stats failed: %w", err)
	}

	for _, row := range rows {
		stats[row.VisitDate] = models.ArticleStatsData{
			ArticleID:  articleID,
			TotalViews: row.PV,
			TotalUV:    row.UV,
		}
	}

	return stats, nil
}

// GetTrendingArticles 获取热门文章排行
func GetTrendingArticles(periodType string, periodDate time.Time, offset, limit int) ([]models.TrendingArticleData, error) {
	var articles []models.TrendingArticleData

	query := `SELECT
			  ta.id, ta.article_id, ta.period_type, ta.period_date,
			  ta.view_count, ta.unique_view_count, ta.rank_position,
			  ta.created_at, ta.updated_at,
			  a.title, a.summary, a.author_id, a.status, a.is_featured,
			  a.like_count, a.comment_count, a.view_count as article_views,
			  a.created_at as article_created_at, a.updated_at as article_updated_at,
			  u.username as author_username, u.nickname as author_nickname
			  FROM trending_articles ta
			  INNER JOIN articles a ON ta.article_id = a.id
			  INNER JOIN users u ON a.author_id = u.id
			  WHERE ta.period_type = ? AND ta.period_date = ?
			  ORDER BY ta.rank_position
			  LIMIT ? OFFSET ?`

	err := db.Select(&articles, query, periodType, periodDate.Format("2006-01-02"), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get trending articles failed: %w", err)
	}

	return articles, nil
}

// UpdateTrendingArticles 更新热门文章排行（调用存储过程）
func UpdateTrendingArticles(periodType string, periodDate time.Time) error {
	query := `CALL InitTrendingArticles(?, ?)`
	_, err := db.Exec(query, periodType, periodDate.Format("2006-01-02"))
	if err != nil {
		return fmt.Errorf("update trending articles failed: %w", err)
	}
	return nil
}

// CheckIPVisitLimit 检查IP访问限制（防刷机制）
func CheckIPVisitLimit(articleID int64, ipAddress string, visitDate time.Time, limit int) (bool, int, error) {
	// 获取当日该IP对该文章的访问次数
	var count int
	query := `SELECT COUNT(*) FROM article_visits
			  WHERE article_id = ? AND ip_address = ? AND visit_date = ?`
	err := db.Get(&count, query, articleID, ipAddress, visitDate.Format("2006-01-02"))
	if err != nil {
		return false, 0, fmt.Errorf("get ip visit count failed: %w", err)
	}

	// 如果已达到限制，返回false和当前次数
	if count >= limit {
		return false, count, nil
	}

	// 否则返回true和当前次数
	return true, count, nil
}

// BatchGetArticleStats 批量获取文章统计数据
func BatchGetArticleStats(articleIDs []int64) (map[int64]models.ArticleStatsData, error) {
	stats := make(map[int64]models.ArticleStatsData)

	if len(articleIDs) == 0 {
		return stats, nil
	}

	// 获取总浏览量
	query := `SELECT id, view_count FROM articles WHERE id IN (?)`
	query, args, err := sqlx.In(query, articleIDs)
	if err != nil {
		return nil, fmt.Errorf("build in query failed: %w", err)
	}

	type articleView struct {
		ID        int64 `db:"id"`
		ViewCount int   `db:"view_count"`
	}

	var articles []articleView
	err = db.Select(&articles, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get article views failed: %w", err)
	}

	for _, article := range articles {
		stats[article.ID] = models.ArticleStatsData{
			ArticleID:  article.ID,
			TotalViews: article.ViewCount,
		}
	}

	// 获取独立访客数（最近30天）
	startDate := time.Now().AddDate(0, 0, -30)
	query = `SELECT article_id, COUNT(DISTINCT ip_address) as uv
			 FROM article_visits
			 WHERE article_id IN (?) AND visit_date >= ?
			 GROUP BY article_id`
	query, args, err = sqlx.In(query, articleIDs, startDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("build in query for uv failed: %w", err)
	}

	type uvStat struct {
		ArticleID int64 `db:"article_id"`
		UV        int   `db:"uv"`
	}

	var uvStats []uvStat
	err = db.Select(&uvStats, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get article uv failed: %w", err)
	}

	for _, stat := range uvStats {
		if data, exists := stats[stat.ArticleID]; exists {
			data.TotalUV = stat.UV
			stats[stat.ArticleID] = data
		}
	}

	return stats, nil
}