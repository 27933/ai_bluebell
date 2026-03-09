package logic

import (
	"fmt"
	"time"

	"bluebell/dao/mysql"
	"bluebell/models"
	"go.uber.org/zap"
)

// GetArticleStatsWithUV 获取文章统计数据（包含UV）
func GetArticleStatsWithUV(articleID int64, days int) (*models.ArticleStatsData, error) {
	// 参数校验
	if days < 1 || days > 90 {
		return nil, fmt.Errorf("days must be between 1 and 90")
	}

	// 获取文章总浏览量
	article, err := mysql.GetArticleById(articleID)
	if err != nil {
		return nil, fmt.Errorf("get article failed: %w", err)
	}

	// 获取总独立访客数（最近30天）
	totalUV, err := mysql.GetArticleUniqueViews(articleID,
		time.Now().AddDate(0, 0, -30), time.Now())
	if err != nil {
		return nil, fmt.Errorf("get article unique views failed: %w", err)
	}

	// 构建返回数据
	stats := &models.ArticleStatsData{
		ArticleID:  articleID,
		TotalViews: article.ViewCount,
		TotalUV:    totalUV,
	}

	return stats, nil
}

// GetArticleTrendData 获取文章访问趋势数据
func GetArticleTrendData(articleID int64, days int, groupBy string) ([]models.ArticleTrendData, error) {
	// 参数校验
	if days < 1 || days > 90 {
		return nil, fmt.Errorf("days must be between 1 and 90")
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	switch groupBy {
	case "hour":
		// 获取小时级别数据（只支持单日）
		if days != 1 {
			return nil, fmt.Errorf("hour groupBy only supports 1 day")
		}
		return mysql.GetArticleHourlyStats(articleID, startDate)

	case "day":
		// 获取日级别数据
		stats, err := mysql.GetArticleDailyStatsWithUV(articleID, startDate, endDate)
		if err != nil {
			return nil, err
		}

		// 转换为趋势数据格式
		var trendData []models.ArticleTrendData
		for date, data := range stats {
			trendData = append(trendData, models.ArticleTrendData{
				Label: date,
				Value: data.TotalViews,
			})
		}
		return trendData, nil

	case "week":
		// 获取周级别数据（简化实现）
		return getWeeklyTrendData(articleID, startDate, endDate)

	case "month":
		// 获取月级别数据（简化实现）
		return getMonthlyTrendData(articleID, startDate, endDate)

	default:
		return nil, fmt.Errorf("invalid group_by parameter")
	}
}

// getWeeklyTrendData 获取周级别趋势数据
func getWeeklyTrendData(articleID int64, startDate, endDate time.Time) ([]models.ArticleTrendData, error) {
	// 这里简化实现，按周聚合日数据
	stats, err := mysql.GetArticleDailyStatsWithUV(articleID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 按周聚合（简化处理）
	weekStats := make(map[string]int)
	for dateStr, data := range stats {
		date, _ := time.Parse("2006-01-02", dateStr)
		weekStart := date.AddDate(0, 0, -int(date.Weekday()))
		weekKey := weekStart.Format("2006-W02")
		weekStats[weekKey] += data.TotalViews
	}

	var trendData []models.ArticleTrendData
	for week, views := range weekStats {
		trendData = append(trendData, models.ArticleTrendData{
			Label: week,
			Value: views,
		})
	}

	return trendData, nil
}

// getMonthlyTrendData 获取月级别趋势数据
func getMonthlyTrendData(articleID int64, startDate, endDate time.Time) ([]models.ArticleTrendData, error) {
	// 获取日级别数据
	stats, err := mysql.GetArticleDailyStatsWithUV(articleID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 按月聚合
	monthStats := make(map[string]int)
	for dateStr, data := range stats {
		date, _ := time.Parse("2006-01-02", dateStr)
		monthKey := date.Format("2006-01")
		monthStats[monthKey] += data.TotalViews
	}

	var trendData []models.ArticleTrendData
	for month, views := range monthStats {
		trendData = append(trendData, models.ArticleTrendData{
			Label: month,
			Value: views,
		})
	}

	return trendData, nil
}

// GetTrendingArticles 获取热门文章排行
func GetTrendingArticles(periodType string, page, size int) ([]models.TrendingArticleData, error) {
	// 计算排行日期
	var periodDate time.Time
	switch periodType {
	case "daily":
		periodDate = time.Now().AddDate(0, 0, -1) // 昨天的数据
	case "weekly":
		periodDate = time.Now().AddDate(0, 0, -7) // 上周的数据
	case "monthly":
		periodDate = time.Now().AddDate(0, -1, 0) // 上月的数据
	default:
		return nil, fmt.Errorf("invalid period type")
	}

	// 如果没有数据，尝试更新排行
	articles, err := mysql.GetTrendingArticles(periodType, periodDate, (page-1)*size, size)
	if err != nil {
		return nil, fmt.Errorf("get trending articles failed: %w", err)
	}

	// 如果没有数据，尝试生成排行数据
	if len(articles) == 0 {
		if err := mysql.UpdateTrendingArticles(periodType, periodDate); err != nil {
			zap.L().Error("update trending articles failed", zap.Error(err))
		} else {
			// 重新获取
			articles, err = mysql.GetTrendingArticles(periodType, periodDate, (page-1)*size, size)
			if err != nil {
				return nil, fmt.Errorf("get trending articles failed: %w", err)
			}
		}
	}

	return articles, nil
}

// BatchGetArticleStats 批量获取文章统计数据
func BatchGetArticleStats(articleIDs []int64) (map[int64]models.ArticleStatsData, error) {
	if len(articleIDs) == 0 {
		return make(map[int64]models.ArticleStatsData), nil
	}

	// 调用DAO层批量获取
	stats, err := mysql.BatchGetArticleStats(articleIDs)
	if err != nil {
		return nil, fmt.Errorf("batch get article stats failed: %w", err)
	}

	return stats, nil
}