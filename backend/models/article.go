package models

import "time"

// Article 文章表对应结构体
type Article struct {
	ID               int64     `json:"id,string" db:"id"`
	Title            string    `json:"title" db:"title" binding:"required"`
	Content          string    `json:"content" db:"content" binding:"required"`
	Summary          string    `json:"summary" db:"summary"`
	WordCount        int       `json:"word_count" db:"word_count"`
	AuthorID         int64     `json:"author_id,string" db:"author_id"`
	Status           string    `json:"status" db:"status"`
	IsFeatured       bool      `json:"is_featured" db:"is_featured"`
	FeaturedAt       *time.Time `json:"featured_at,omitempty" db:"featured_at"`
	AllowComment     bool      `json:"allow_comment" db:"allow_comment"`
	LikeCount        int       `json:"like_count" db:"like_count"`
	CommentCount     int       `json:"comment_count" db:"comment_count"`
	Slug             string    `json:"slug" db:"slug"`
	MetaKeywords     string    `json:"meta_keywords,omitempty" db:"meta_keywords"`
	MetaDescription  string    `json:"meta_description,omitempty" db:"meta_description"`
	ViewCount        int       `json:"view_count" db:"view_count"`
	IPAddress        string    `json:"-" db:"ip_address"`
	UserAgent        string    `json:"-" db:"user_agent"`
	Extra            string    `json:"extra,omitempty" db:"extra"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// ArticleStatus 文章状态常量
type ArticleStatus string

const (
	ArticleStatusDraft    ArticleStatus = "draft"
	ArticleStatusPublished ArticleStatus = "published"
	ArticleStatusOffline   ArticleStatus = "offline"
)

// ApiArticleDetail 文章详情接口返回结构体
type ApiArticleDetail struct {
	ID               int64      `json:"id,string"`
	Title            string     `json:"title"`
	Content          string     `json:"content"`
	Summary          string     `json:"summary"`
	WordCount        int        `json:"word_count"`
	Author           *AuthorInfo `json:"author"`
	Tags             []Tag      `json:"tags"`
	Status           string     `json:"status"`
	IsFeatured       bool       `json:"is_featured"`
	FeaturedAt       *time.Time `json:"featured_at,omitempty"`
	AllowComment     bool       `json:"allow_comment"`
	LikeCount        int        `json:"like_count"`
	CommentCount     int        `json:"comment_count"`
	Slug             string     `json:"slug"`
	MetaKeywords     string     `json:"meta_keywords,omitempty"`
	MetaDescription  string     `json:"meta_description,omitempty"`
	ViewCount        int        `json:"view_count"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ApiArticleListItem 文章列表项接口返回结构体
type ApiArticleListItem struct {
	ID           int64      `json:"id,string"`
	Title        string     `json:"title"`
	Summary      string     `json:"summary"`
	Author       AuthorInfo `json:"author"`
	Tags         []Tag      `json:"tags"`
	ViewCount    int        `json:"view_count"`
	LikeCount    int        `json:"like_count"`
	CommentCount int        `json:"comment_count"`
	IsFeatured   bool       `json:"is_featured"`
	IsRecent     bool       `json:"is_recent"`     // 最近更新标记（24小时内）
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// AuthorInfo 作者信息结构体
type AuthorInfo struct {
	ID       int64  `json:"id,string"`
	Username string `json:"username"`
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Bio      string `json:"bio,omitempty"`
}

// AdminArticleListItem 管理员文章列表项（含作者信息，JOIN users 表）
type AdminArticleListItem struct {
	ID             int64      `json:"id,string" db:"id"`
	Title          string     `json:"title" db:"title"`
	Summary        string     `json:"summary" db:"summary"`
	WordCount      int        `json:"word_count" db:"word_count"`
	AuthorID       int64      `json:"author_id,string" db:"author_id"`
	AuthorUsername string     `json:"author_username" db:"author_username"`
	AuthorNickname string     `json:"author_nickname,omitempty" db:"author_nickname"`
	Status         string     `json:"status" db:"status"`
	IsFeatured     bool       `json:"is_featured" db:"is_featured"`
	FeaturedAt     *time.Time `json:"featured_at,omitempty" db:"featured_at"`
	AllowComment   bool       `json:"allow_comment" db:"allow_comment"`
	LikeCount      int        `json:"like_count" db:"like_count"`
	CommentCount   int        `json:"comment_count" db:"comment_count"`
	ViewCount      int        `json:"view_count" db:"view_count"`
	Slug           string     `json:"slug" db:"slug"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// ParamCreateArticle 创建文章请求参数
type ParamCreateArticle struct {
	Title           string   `json:"title" binding:"required,min=1,max=200"`
	Content         string   `json:"content" binding:"required,min=1"`
	Tags            []string `json:"tags" binding:"max=10"`
	Status          string   `json:"status" binding:"omitempty,oneof=draft published offline"`
	AllowComment    bool     `json:"allow_comment"`
	Slug            string   `json:"slug" binding:"omitempty,max=200"`
	MetaKeywords    string   `json:"meta_keywords" binding:"omitempty,max=500"`
	MetaDescription string   `json:"meta_description" binding:"omitempty,max=500"`
}

// ParamUpdateArticle 更新文章请求参数
type ParamUpdateArticle struct {
	Title           string   `json:"title" binding:"omitempty,min=1,max=200"`
	Content         string   `json:"content" binding:"omitempty,min=1"`
	Tags            []string `json:"tags" binding:"max=10"`
	Summary         string   `json:"summary" binding:"omitempty,max=500"`
	Status          string   `json:"status" binding:"omitempty,oneof=draft published offline"`
	IsFeatured      *bool    `json:"is_featured"`
	AllowComment    *bool    `json:"allow_comment"`
	Slug            string   `json:"slug" binding:"omitempty,max=200"`
	MetaKeywords    string   `json:"meta_keywords" binding:"omitempty,max=500"`
	MetaDescription string   `json:"meta_description" binding:"omitempty,max=500"`
}

// ParamArticleList 文章列表查询参数
type ParamArticleList struct {
	Sort       string `form:"sort" binding:"omitempty,oneof=time popular"`
	Page       int    `form:"page" binding:"omitempty,min=1"`
	Size       int    `form:"size" binding:"omitempty,min=1,max=50"`
	Tag        string `form:"tag"`
	Status     string `form:"status" binding:"omitempty,oneof=draft published offline all"`
	AuthorID   int64  `form:"author_id"`
	AuthorName string `form:"author_name"` // 按作者用户名搜索
	Days       int    `form:"days" binding:"omitempty,min=1,max=365"`
	Keyword    string `form:"keyword"`
	IsFeatured *bool  `form:"is_featured"` // 精选筛选，nil=不过滤
}

// ParamArticleStatus 更新文章状态参数
type ParamArticleStatus struct {
	Status string `json:"status" binding:"required,oneof=draft published offline"`
}

// ParamArticleFeatured 更新文章精选状态参数
type ParamArticleFeatured struct {
	IsFeatured bool `json:"is_featured"`
}