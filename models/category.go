package models

import "time"

// Category 栏目模型
type Category struct {
	ID           int64     `db:"id" json:"id"`
	CategoryName string    `db:"category_name" json:"category_name"`
	Introduction string    `db:"introduction" json:"introduction"`
	CreatedBy    int64     `db:"created_by" json:"created_by"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// ArticleCategory 文章-栏目关联模型
type ArticleCategory struct {
	ID         int64     `db:"id" json:"id"`
	ArticleID  int64     `db:"article_id" json:"article_id"`
	CategoryID int64     `db:"category_id" json:"category_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// ParamCreateCategory 创建栏目参数
type ParamCreateCategory struct {
	CategoryName string `json:"category_name" binding:"required,min=2,max=128"`
	Introduction string `json:"introduction" binding:"required,min=2,max=256"`
}

// ParamUpdateCategory 更新栏目参数
type ParamUpdateCategory struct {
	CategoryName string `json:"category_name" binding:"min=2,max=128"`
	Introduction string `json:"introduction" binding:"min=2,max=256"`
}

// ParamAddArticleCategories 添加文章到栏目参数
type ParamAddArticleCategories struct {
	CategoryIDs []int64 `json:"category_ids" binding:"required,min=1"`
}

// ApiCategoryDetail 栏目详情（包含文章数量）
type ApiCategoryDetail struct {
	ID             int64     `db:"id" json:"id"`
	CategoryName   string    `db:"category_name" json:"category_name"`
	Introduction   string    `db:"introduction" json:"introduction"`
	CreatedBy      int64     `db:"created_by" json:"created_by"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	ArticleCount   int64     `db:"article_count" json:"article_count"`   // 栏目下的文章数量
}