package mysql

import (
	"bluebell/models"
	"database/sql"
)

// CreateCategory 创建栏目
func CreateCategory(category *models.Category) error {
	sqlStr := `INSERT INTO categories (category_name, introduction, created_by) VALUES (?, ?, ?)`
	result, err := db.Exec(sqlStr, category.CategoryName, category.Introduction, category.CreatedBy)
	if err != nil {
		return err
	}
	categoryID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	category.ID = categoryID
	return nil
}

// GetCategoryById 根据ID获取栏目
func GetCategoryById(id int64) (*models.Category, error) {
	category := new(models.Category)
	sqlStr := `SELECT id, category_name, introduction, created_by, created_at, updated_at FROM categories WHERE id = ?`
	err := db.Get(category, sqlStr, id)
	if err == sql.ErrNoRows {
		return nil, ErrorCategoryNotExist
	}
	return category, err
}

// GetCategoryByName 根据名称获取栏目
func GetCategoryByName(name string) (*models.Category, error) {
	category := new(models.Category)
	sqlStr := `SELECT id, category_name, introduction, created_by, created_at, updated_at FROM categories WHERE category_name = ?`
	err := db.Get(category, sqlStr, name)
	if err == sql.ErrNoRows {
		return nil, ErrorCategoryNotExist
	}
	return category, err
}

// GetCategoryList 获取栏目列表
func GetCategoryList() ([]*models.Category, error) {
	var categories []*models.Category
	sqlStr := `SELECT id, category_name, introduction, created_by, created_at, updated_at FROM categories ORDER BY created_at DESC`
	err := db.Select(&categories, sqlStr)
	return categories, err
}

// UpdateCategory 更新栏目
func UpdateCategory(category *models.Category) error {
	sqlStr := `UPDATE categories SET category_name = ?, introduction = ? WHERE id = ?`
	_, err := db.Exec(sqlStr, category.CategoryName, category.Introduction, category.ID)
	return err
}

// DeleteCategory 删除栏目
func DeleteCategory(id int64) error {
	sqlStr := `DELETE FROM categories WHERE id = ?`
	_, err := db.Exec(sqlStr, id)
	return err
}

// GetCategoryArticleCount 获取栏目下的文章数量
func GetCategoryArticleCount(categoryID int64) (int64, error) {
	var count int64
	sqlStr := `SELECT COUNT(*) FROM article_categories WHERE category_id = ?`
	err := db.Get(&count, sqlStr, categoryID)
	return count, err
}

// AddArticleToCategories 添加文章到多个栏目
func AddArticleToCategories(articleID int64, categoryIDs []int64) error {
	// 开启事务
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 先删除文章已有的栏目关联
	delSql := `DELETE FROM article_categories WHERE article_id = ?`
	if _, err = tx.Exec(delSql, articleID); err != nil {
		return err
	}

	// 批量插入新的关联
	if len(categoryIDs) > 0 {
		insertSql := `INSERT INTO article_categories (article_id, category_id) VALUES (?, ?)`
		for _, categoryID := range categoryIDs {
			if _, err = tx.Exec(insertSql, articleID, categoryID); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// GetArticlesByCategoryID 获取栏目下的文章ID列表
func GetArticlesByCategoryID(categoryID int64, offset, limit int) ([]int64, error) {
	var articleIDs []int64
	sqlStr := `SELECT article_id FROM article_categories WHERE category_id = ? ORDER BY id DESC LIMIT ? OFFSET ?`
	err := db.Select(&articleIDs, sqlStr, categoryID, limit, offset)
	return articleIDs, err
}

// GetCategoriesByArticleID 获取文章所属的栏目列表
func GetCategoriesByArticleID(articleID int64) ([]*models.Category, error) {
	var categories []*models.Category
	sqlStr := `
		SELECT c.id, c.category_name, c.introduction, c.created_by, c.created_at, c.updated_at
		FROM categories c
		INNER JOIN article_categories ac ON c.id = ac.category_id
		WHERE ac.article_id = ?
		ORDER BY c.created_at DESC`
	err := db.Select(&categories, sqlStr, articleID)
	return categories, err
}

// CheckArticleInCategory 检查文章是否已在栏目中
func CheckArticleInCategory(articleID, categoryID int64) (bool, error) {
	var count int
	sqlStr := `SELECT COUNT(*) FROM article_categories WHERE article_id = ? AND category_id = ?`
	err := db.Get(&count, sqlStr, articleID, categoryID)
	return count > 0, err
}

// DeleteArticleCategories 删除文章的所有栏目关联
func DeleteArticleCategories(articleID int64) error {
	sqlStr := `DELETE FROM article_categories WHERE article_id = ?`
	_, err := db.Exec(sqlStr, articleID)
	return err
}

// DeleteCategoryArticles 删除栏目的所有文章关联
func DeleteCategoryArticles(categoryID int64) error {
	sqlStr := `DELETE FROM article_categories WHERE category_id = ?`
	_, err := db.Exec(sqlStr, categoryID)
	return err
}

// GetCategoryListWithArticleCount 获取栏目列表（包含文章数量）
func GetCategoryListWithArticleCount() ([]*models.ApiCategoryDetail, error) {
	var categories []*models.ApiCategoryDetail
	sqlStr := `
		SELECT c.id, c.category_name, c.introduction, c.created_by, c.created_at,
			(SELECT COUNT(*) FROM article_categories ac WHERE ac.category_id = c.id) as article_count
		FROM categories c
		ORDER BY c.created_at DESC`
	err := db.Select(&categories, sqlStr)
	return categories, err
}