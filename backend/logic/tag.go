package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// GetAllTags 获取所有标签
func GetAllTags() ([]*models.Tag, error) {
	tags, err := mysql.GetAllTags()
	if err != nil {
		zap.L().Error("mysql.GetAllTags() failed", zap.Error(err))
		return nil, err
	}
	return tags, nil
}

// GetArticlesByTag 根据标签获取文章
func GetArticlesByTag(tagID int64, param *models.ParamArticleList) ([]*models.Article, int64, error) {
	// 参数验证
	if param.Page < 1 {
		param.Page = 1
	}
	if param.Size < 1 || param.Size > 50 {
		param.Size = 20
	}
	if param.Sort == "" {
		param.Sort = "time"
	}
	if param.Status == "" {
		param.Status = "published"
	}

	// 调用DAO层
	articles, total, err := mysql.GetArticlesByTag(tagID, param)
	if err != nil {
		zap.L().Error("mysql.GetArticlesByTag() failed", zap.Error(err))
		return nil, 0, err
	}

	return articles, total, nil
}

// CreateTag 创建标签
func CreateTag(req *models.ParamCreateTag) (*models.Tag, error) {
	// 参数验证
	if req.Name == "" {
		return nil, errors.New("tag name cannot be empty")
	}

	// 检查标签是否已存在
	existingTag, err := mysql.GetTagByName(req.Name)
	if err == nil && existingTag != nil {
		return nil, errors.New("tag already exists")
	}

	// 创建标签
	tag := &models.Tag{
		ID:          snowflake.GenID(),
		Name:        req.Name,
		Slug:        generateSlug(req.Name),
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存标签
	if err := mysql.CreateTag(tag); err != nil {
		zap.L().Error("mysql.CreateTag() failed", zap.Error(err))
		return nil, err
	}

	return tag, nil
}

// UpdateTag 更新标签
func UpdateTag(tagID int64, req *models.ParamUpdateTag) error {
	// 检查标签是否存在
	tag, err := mysql.GetTagById(tagID)
	if err != nil {
		return mysql.ErrorTagNotExist
	}

	// 检查新名称是否已被使用（如果提供了新名称）
	if req.Name != "" && req.Name != tag.Name {
		existingTag, err := mysql.GetTagByName(req.Name)
		if err == nil && existingTag != nil && existingTag.ID != tagID {
			return errors.New("tag name already exists")
		}
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
		updates["slug"] = generateSlug(req.Name)
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if len(updates) == 0 {
		return nil // 没有需要更新的字段
	}

	// 更新标签
	if err := mysql.UpdateTag(tagID, updates); err != nil {
		zap.L().Error("mysql.UpdateTag() failed", zap.Error(err))
		return err
	}

	return nil
}

// DeleteTag 删除标签
func DeleteTag(tagID int64) error {
	// 检查标签是否存在
	if _, err := mysql.GetTagById(tagID); err != nil {
		return mysql.ErrorTagNotExist
	}

	// 检查标签是否正在被使用
	count, err := mysql.CountArticlesByTag(tagID)
	if err != nil {
		zap.L().Error("mysql.CountArticlesByTag() failed", zap.Error(err))
		return err
	}
	if count > 0 {
		return errors.New(fmt.Sprintf("tag is used by %d articles", count))
	}

	// 删除标签
	if err := mysql.DeleteTag(tagID); err != nil {
		zap.L().Error("mysql.DeleteTag() failed", zap.Error(err))
		return err
	}

	return nil
}

// GetAuthorTags 获取作者标签列表
func GetAuthorTags(authorID int64) ([]*models.Tag, error) {
	tags, err := mysql.GetAuthorTags(authorID)
	if err != nil {
		zap.L().Error("mysql.GetAuthorTags() failed", zap.Error(err))
		return nil, err
	}
	return tags, nil
}

// GetArticlesByTagSlug 根据标签slug获取文章
func GetArticlesByTagSlug(tagSlug string, page, size int) ([]*models.Article, int64, error) {
	// 先获取标签
	tag, err := mysql.GetTagBySlug(tagSlug)
	if err != nil {
		return nil, 0, errors.New("tag not found")
	}

	param := &models.ParamArticleList{
		Page:   page,
		Size:   size,
		Sort:   "time",
		Status: "published",
	}

	return GetArticlesByTag(tag.ID, param)
}

