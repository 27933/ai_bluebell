package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"errors"

	"go.uber.org/zap"
)

// GetRSSArticles 获取RSS文章列表
func GetRSSArticles(limit int) ([]*models.Article, error) {
	if limit < 1 || limit > 50 {
		limit = 20 // 默认获取20篇
	}

	// 构建查询参数 - 获取已发布的最新文章
	param := &models.ParamArticleList{
		Page:   1,
		Size:   limit,
		Sort:   "time",
		Status: "published",
	}

	articles, _, err := mysql.GetArticleList(param)
	if err != nil {
		zap.L().Error("mysql.GetArticleList() failed", zap.Error(err))
		return nil, errors.New("failed to get articles for RSS")
	}

	// 需要获取作者信息
	for range articles {
		// 这里可以填充作者信息
		// 在实际应用中需要查询作者表获取详细信息
		// 这里简化为只返回文章信息
	}

	return articles, nil
}