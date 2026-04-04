package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/ecode"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

// CreateCategory 创建栏目
func CreateCategory(userID int64, p *models.ParamCreateCategory) error {
	// 1. 检查栏目名称是否已存在
	if _, err := mysql.GetCategoryByName(p.CategoryName); err == nil {
		return ecode.ErrCategoryExist
	}

	// 2. 生成栏目ID
	categoryID := snowflake.GenID()

	// 3. 构造栏目实例
	category := &models.Category{
		ID:           categoryID,
		CategoryName: p.CategoryName,
		Introduction: p.Introduction,
		CreatedBy:    userID,
	}

	// 4. 保存到数据库
	return mysql.CreateCategory(category)
}

// GetCategoryList 获取栏目列表
func GetCategoryList() ([]*models.ApiCategoryDetail, error) {
	// 获取栏目列表（包含文章数量）
	return mysql.GetCategoryListWithArticleCount()
}

// GetCategoryById 根据ID获取栏目详情
func GetCategoryById(categoryID int64) (*models.Category, error) {
	return mysql.GetCategoryById(categoryID)
}

// UpdateCategory 更新栏目
func UpdateCategory(userID int64, userRole string, categoryID int64, p *models.ParamUpdateCategory) error {
	// 1. 获取栏目信息
	category, err := mysql.GetCategoryById(categoryID)
	if err != nil {
		return err
	}

	// 2. 权限检查：只有创建者和管理员可以更新
	if category.CreatedBy != userID && userRole != "admin" {
		return ecode.ErrNoPermissionCategory
	}

	// 3. 如果有新的栏目名称，检查是否已存在
	if p.CategoryName != "" && p.CategoryName != category.CategoryName {
		if _, err := mysql.GetCategoryByName(p.CategoryName); err == nil {
			return ecode.ErrCategoryExist
		}
		category.CategoryName = p.CategoryName
	}

	// 4. 更新字段
	if p.Introduction != "" {
		category.Introduction = p.Introduction
	}

	// 5. 保存更新
	return mysql.UpdateCategory(category)
}

// DeleteCategory 删除栏目
func DeleteCategory(userID int64, userRole string, categoryID int64) error {
	// 1. 获取栏目信息
	category, err := mysql.GetCategoryById(categoryID)
	if err != nil {
		return err
	}

	// 2. 权限检查：只有创建者和管理员可以删除
	if category.CreatedBy != userID && userRole != "admin" {
		return ecode.ErrNoPermissionCategory
	}

	// 3. 删除栏目的所有文章关联
	if err := mysql.DeleteCategoryArticles(categoryID); err != nil {
		zap.L().Error("DeleteCategoryArticles failed", zap.Error(err))
		return err
	}

	// 4. 删除栏目
	return mysql.DeleteCategory(categoryID)
}

// AddArticleToCategories 添加文章到栏目
func AddArticleToCategories(userID int64, userRole string, articleID int64, categoryIDs []int64) error {
	// 1. 检查文章是否存在
	article, err := mysql.GetArticleById(articleID)
	if err != nil {
		return err
	}

	// 2. 权限检查：作者只能关联自己的文章，管理员可以关联任何文章
	if article.AuthorID != userID && userRole != "admin" {
		return ecode.ErrNoPermission
	}

	// 3. 检查所有栏目是否存在
	for _, categoryID := range categoryIDs {
		if _, err := mysql.GetCategoryById(categoryID); err != nil {
			return ecode.ErrCategoryNotExist
		}
	}

	// 4. 添加关联（使用事务，先删除旧关联再添加新关联）
	return mysql.AddArticleToCategories(articleID, categoryIDs)
}

// GetArticlesByCategory 获取栏目下的文章ID列表
func GetArticlesByCategory(categoryID int64, page, size int) ([]int64, int64, error) {
	// 1. 检查栏目是否存在
	if _, err := mysql.GetCategoryById(categoryID); err != nil {
		return nil, 0, err
	}

	// 2. 获取文章总数
	total, err := mysql.GetCategoryArticleCount(categoryID)
	if err != nil {
		return nil, 0, err
	}

	// 3. 计算偏移量
	offset := (page - 1) * size
	if offset < 0 {
		offset = 0
	}

	// 4. 获取文章ID列表
	articleIDs, err := mysql.GetArticlesByCategoryID(categoryID, offset, size)
	if err != nil {
		return nil, 0, err
	}

	return articleIDs, total, nil
}

// GetCategoriesByArticle 获取文章所属的栏目列表
func GetCategoriesByArticle(articleID int64) ([]*models.Category, error) {
	// 1. 检查文章是否存在
	if _, err := mysql.GetArticleById(articleID); err != nil {
		return nil, err
	}

	// 2. 获取栏目列表
	return mysql.GetCategoriesByArticleID(articleID)
}