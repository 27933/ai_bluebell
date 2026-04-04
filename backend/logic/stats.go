package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"time"

	"go.uber.org/zap"
)

// GetSystemOverview 获取系统概览
func GetSystemOverview() (*models.SystemOverview, error) {
	overview := &models.SystemOverview{}

	// 获取用户总数
	userCount, err := mysql.GetUserCount()
	if err != nil {
		zap.L().Error("mysql.GetUserCount() failed", zap.Error(err))
		return nil, err
	}
	overview.UserCount = userCount

	// 获取文章总数
	articleCount, err := mysql.GetArticleCount()
	if err != nil {
		zap.L().Error("mysql.GetArticleCount() failed", zap.Error(err))
		return nil, err
	}
	overview.ArticleCount = articleCount

	// 获取评论总数
	commentCount, err := mysql.GetCommentCount()
	if err != nil {
		zap.L().Error("mysql.GetCommentCount() failed", zap.Error(err))
		return nil, err
	}
	overview.CommentCount = commentCount

	// 获取今日新增用户
	todayStart := time.Now().Format("2006-01-02") + " 00:00:00"
	todayUserCount, err := mysql.GetUserCountByDate(todayStart)
	if err != nil {
		zap.L().Error("mysql.GetUserCountByDate() failed", zap.Error(err))
		return nil, err
	}
	overview.TodayNewUserCount = todayUserCount

	// 获取今日新增文章
	todayArticleCount, err := mysql.GetArticleCountByDate(todayStart)
	if err != nil {
		zap.L().Error("mysql.GetArticleCountByDate() failed", zap.Error(err))
		return nil, err
	}
	overview.TodayNewArticleCount = todayArticleCount

	return overview, nil
}

// GetSystemDailyStats 获取系统日统计数据
func GetSystemDailyStats(param *models.ParamStatsDaily) ([]*models.DailyStats, error) {
	// 获取日统计数据
	stats, err := mysql.GetDailyStats(param.StartDate, param.EndDate)
	if err != nil {
		zap.L().Error("mysql.GetDailyStats() failed", zap.Error(err))
		return nil, err
	}

	return stats, nil
}