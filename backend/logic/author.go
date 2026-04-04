package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

// GetAuthorInfo 获取作者信息
func GetAuthorInfo(username string) (*models.AuthorInfoResponse, error) {
	return mysql.GetAuthorInfoByUsername(username)
}

// GetAuthorArticlesList 获取作者的文章列表
func GetAuthorArticlesList(p *models.ParamAuthorArticlesList) ([]*models.Article, int64, error) {
	// 构建查询参数
	params := &models.ParamArticleList{
		Page:      p.Page,
		Size:      p.Size,
		Sort:      p.Sort,
		AuthorName: p.Username, // 使用作者用户名作为查询条件
		Status:    "published", // 只显示已发布的文章
	}

	// 调用通用的文章列表查询
	return mysql.GetArticleList(params)
}