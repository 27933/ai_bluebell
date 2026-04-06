package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"errors"

	"go.uber.org/zap"
)

// GetRSSArticles 获取RSS文章列表（含作者用户名）
func GetRSSArticles(limit int) ([]*models.AdminArticleListItem, error) {
	if limit < 1 || limit > 50 {
		limit = 20
	}

	param := &models.ParamArticleList{
		Page:   1,
		Size:   limit,
		Sort:   "time",
		Status: "published",
	}

	articles, _, err := mysql.GetAdminArticleList(param)
	if err != nil {
		zap.L().Error("mysql.GetAdminArticleList() failed", zap.Error(err))
		return nil, errors.New("failed to get articles for RSS")
	}

	return articles, nil
}
