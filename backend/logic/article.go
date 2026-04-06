package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

// GetFeaturedArticles 获取精选文章
func GetFeaturedArticles(limit int) ([]*models.Article, error) {
	articles, err := mysql.GetFeaturedArticles(limit)
	if err != nil {
		zap.L().Error("mysql.GetFeaturedArticles() failed", zap.Error(err))
		return nil, err
	}
	return articles, nil
}

// GetArticleList 获取文章列表（支持排序和标签过滤）
func GetArticleList(param *models.ParamArticleList) ([]*models.Article, int64, error) {
	// 参数验证
	if param.Page < 1 {
		param.Page = 1
	}
	if param.Size < 1 || param.Size > 50 {
		param.Size = 20
	}

	articles, total, err := mysql.GetArticleList(param)
	if err != nil {
		zap.L().Error("mysql.GetArticleList() failed", zap.Error(err))
		return nil, 0, err
	}

	return articles, total, nil
}

// GetArticleDetail 获取文章详情
func GetArticleDetail(articleID int64) (*models.Article, error) {
	// 获取文章详情
	article, err := mysql.GetArticleById(articleID)
	if err != nil {
		zap.L().Error("mysql.GetArticleById() failed", zap.Int64("articleID", articleID), zap.Error(err))
		return nil, err
	}

	// 检查文章状态
	if article.Status != string(models.ArticleStatusPublished) {
		return nil, fmt.Errorf("article not published")
	}

	// 注意：view_count 由 RecordArticleViewWithAntiCheat 负责更新（带去重）
	// 此处不再直接调用 UpdateArticleView，避免每次 GET 都无限累加

	return article, nil
}

// CreateArticle 创建文章
func CreateArticle(authorID int64, req *models.ParamCreateArticle) (*models.Article, error) {
	// 生成文章ID
	articleID := snowflake.GenID()

	// 创建文章
	now := time.Now()
	article := &models.Article{
		ID:          articleID,
		Title:       req.Title,
		Content:     req.Content,
		Summary:     generateSummary(req.Content),
		WordCount:   countWords(req.Content),
		AuthorID:    authorID,
		Status:      req.Status,
		IsFeatured:  false, // 创建文章时不能设置为精选
		AllowComment: req.AllowComment,
		Slug:        generateSlug(req.Title),
		MetaKeywords: generateMetaKeywords(req.Tags),
		MetaDescription: generateMetaDescription(req.Content),
		Extra:       "{}", // 设置默认的JSON字符串
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 如果状态为空，默认为草稿
	if article.Status == "" {
		article.Status = string(models.ArticleStatusDraft)
	}

	// 保存文章
	if err := mysql.CreateArticle(article); err != nil {
		zap.L().Error("mysql.CreateArticle() failed", zap.Error(err))
		return nil, err
	}

	// 处理标签
	if len(req.Tags) > 0 {
		// 先获取或创建标签
		var tagIds []int64
		for _, tagName := range req.Tags {
			tag, err := mysql.GetTagByName(tagName)
			if err != nil {
				// 标签不存在，创建新标签
				newTag := &models.Tag{
					Name: tagName,
					Slug: generateSlug(tagName),
				}
				if err := mysql.CreateTag(newTag); err != nil {
					zap.L().Error("mysql.CreateTag() failed", zap.Error(err))
					continue
				}
				tagIds = append(tagIds, newTag.ID)
			} else {
				tagIds = append(tagIds, tag.ID)
			}
		}

		// 添加文章标签关联
		if len(tagIds) > 0 {
			if err := mysql.AddArticleTags(article.ID, tagIds); err != nil {
				zap.L().Error("mysql.AddArticleTags() failed", zap.Error(err))
			}
		}
	}

	return article, nil
}

// GetAuthorArticles 获取作者文章列表
func GetAuthorArticles(authorID int64, param *models.ParamArticleList) ([]*models.Article, int64, error) {
	// 参数验证
	if param.Page < 1 {
		param.Page = 1
	}
	if param.Size < 1 || param.Size > 50 {
		param.Size = 20
	}

	// 设置作者ID
	param.AuthorID = authorID

	articles, total, err := mysql.GetArticleList(param)
	if err != nil {
		zap.L().Error("mysql.GetArticleList() failed", zap.Error(err))
		return nil, 0, err
	}

	return articles, total, nil
}

// GetAuthorArticleDetail 获取作者文章详情
func GetAuthorArticleDetail(authorID, articleID int64) (*models.Article, error) {
	// 获取文章详情
	article, err := mysql.GetArticleById(articleID)
	if err != nil {
		return nil, err
	}

	// 检查是否属于该作者
	if article.AuthorID != authorID {
		return nil, fmt.Errorf("article not found or access denied")
	}

	return article, nil
}

// UpdateArticle 更新文章
func UpdateArticle(authorID, articleID int64, req *models.ParamUpdateArticle) error {
	// 检查文章是否存在且属于该作者
	if _, err := GetAuthorArticleDetail(authorID, articleID); err != nil {
		return err
	}

	// 构建更新字段
	updates := make(map[string]interface{})

	if req.Title != "" {
		updates["title"] = req.Title
		updates["slug"] = generateSlug(req.Title)
	}
	if req.Content != "" {
		updates["content"] = req.Content
		updates["summary"] = generateSummary(req.Content)
		updates["word_count"] = countWords(req.Content)
		updates["meta_description"] = generateMetaDescription(req.Content)
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.AllowComment != nil {
		updates["allow_comment"] = *req.AllowComment
	}

	// 更新文章
	if err := mysql.UpdateArticle(articleID, updates); err != nil {
		zap.L().Error("mysql.UpdateArticle() failed", zap.Error(err))
		return err
	}

	// 更新标签
	if len(req.Tags) > 0 {
		// 将标签名转换为ID
		tagIds := make([]int64, 0, len(req.Tags))
		for _, tagName := range req.Tags {
			tag, err := mysql.GetTagByName(tagName)
			if err != nil {
				zap.L().Warn("tag not found", zap.String("tag", tagName))
				continue
			}
			tagIds = append(tagIds, tag.ID)
		}

		if len(tagIds) > 0 {
			if err := mysql.UpdateArticleTags(articleID, tagIds); err != nil {
				zap.L().Error("mysql.UpdateArticleTags() failed", zap.Error(err))
			}
		}
	}

	return nil
}

// DeleteArticle 删除文章
func DeleteArticle(authorID, articleID int64) error {
	// 检查文章是否存在且属于该作者
	if _, err := GetAuthorArticleDetail(authorID, articleID); err != nil {
		return err
	}

	// 删除文章（软删除）
	if err := mysql.DeleteArticle(articleID); err != nil {
		zap.L().Error("mysql.DeleteArticle() failed", zap.Error(err))
		return err
	}

	return nil
}

// UpdateArticleStatus 更新文章状态
func UpdateArticleStatus(authorID, articleID int64, status string) error {
	// 检查文章是否存在且属于该作者
	if _, err := GetAuthorArticleDetail(authorID, articleID); err != nil {
		return err
	}

	if err := mysql.UpdateArticleStatus(articleID, status); err != nil {
		zap.L().Error("mysql.UpdateArticleStatus() failed", zap.Error(err))
		return err
	}

	return nil
}

// SearchArticles 搜索文章
func SearchArticles(param *models.ParamArticleList) ([]*models.Article, int64, error) {
	// 验证搜索条件：关键词、作者ID、作者名、标签至少有一个
	if param.Keyword == "" && param.AuthorID == 0 && param.AuthorName == "" && param.Tag == "" {
		return nil, 0, fmt.Errorf("at least one search condition is required")
	}

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

	articles, total, err := mysql.GetArticleList(param)
	if err != nil {
		zap.L().Error("mysql.GetArticleList() failed", zap.Error(err))
		return nil, 0, err
	}

	return articles, total, nil
}

// GetAdminArticles 管理员获取所有文章
func GetAdminArticles(param *models.ParamArticleList) ([]*models.Article, int64, error) {
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

	// 管理员可以查看所有状态的文章
	articles, total, err := mysql.GetArticleList(param)
	if err != nil {
		zap.L().Error("mysql.GetArticleList() failed", zap.Error(err))
		return nil, 0, err
	}

	return articles, total, nil
}

// AdminSetArticleFeatured 管理员设置文章精选状态
func AdminSetArticleFeatured(id int64, isFeatured bool) error {
	return mysql.UpdateArticleFeatured(id, isFeatured)
}

// GetArticlesByTagSlug 根据标签获取文章（旧版本，接受tagSlug）
func GetArticlesByTagSlugOld(tagSlug string, page, size int) ([]*models.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}

	return GetArticlesByTagSlug(tagSlug, page, size)
}

// GetArticlesByTagWithID 根据标签ID获取文章（适应controller调用）
func GetArticlesByTagWithID(tagID int64, param *models.ParamArticleList) ([]*models.Article, int64, error) {
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

	articles, total, err := mysql.GetArticlesByTag(tagID, param)
	if err != nil {
		zap.L().Error("mysql.GetArticlesByTag() failed", zap.Error(err))
		return nil, 0, err
	}

	return articles, total, nil
}

// SetFeaturedArticle 设置精选文章（管理员功能）
func SetFeaturedArticle(articleID int64, isFeatured bool) error {
	if err := mysql.UpdateArticleFeatured(articleID, isFeatured); err != nil {
		zap.L().Error("mysql.UpdateArticleFeatured() failed", zap.Error(err))
		return err
	}

	return nil
}

// UpdateArticleFeatured 更新文章精选状态（与controller匹配的函数名）
func UpdateArticleFeatured(userID, articleID int64, isFeatured bool) error {
	// 检查文章是否存在且属于该作者
	if _, err := GetAuthorArticleDetail(userID, articleID); err != nil {
		return err
	}

	// 检查用户是否有权限设置精选（应该是管理员权限）
	// 这里简化处理，实际应该检查用户角色
	return SetFeaturedArticle(articleID, isFeatured)
}

// RecordArticleView 记录文章访问
func RecordArticleView(articleID int64) error {
	// 异步记录访问量
	go func() {
		date := time.Now()
		if err := mysql.UpdateArticleView(articleID, date); err != nil {
			zap.L().Error("mysql.UpdateArticleView() failed", zap.Error(err))
		}
	}()
	return nil
}

// GetArticleDailyStats 获取文章日访问量统计
func GetArticleDailyStats(articleID int64, days int) ([]map[string]interface{}, error) {
	if days < 1 {
		days = 30
	}
	if days > 90 {
		days = 90
	}

	stats, err := mysql.GetArticleDailyStats(articleID, days)
	if err != nil {
		zap.L().Error("mysql.GetArticleDailyStats() failed", zap.Error(err))
		return nil, err
	}

	return stats, nil
}

// 辅助函数

func generateSummary(content string) string {
	// 移除markdown标记，截取前200字符
	plainText := removeMarkdown(content)
	if len(plainText) > 200 {
		return plainText[:200] + "..."
	}
	return plainText
}

func countWords(content string) int {
	// 简单统计中文字符和英文单词
	words := 0
	inWord := false
	for _, r := range content {
		if r >= 0x4E00 && r <= 0x9FA5 {
			words++ // 中文字符
		} else if r == ' ' || r == '\n' || r == '\t' {
			inWord = false
		} else if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			if !inWord {
				words++
				inWord = true
			}
		} else {
			inWord = false
		}
	}
	return words
}

func generateSlug(title string) string {
	// 简单的slug生成，实际需要更复杂的逻辑
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, ".", "")
	slug = strings.ReplaceAll(slug, ",", "")
	slug = strings.ReplaceAll(slug, "?", "")
	slug = strings.ReplaceAll(slug, "!", "")
	return slug
}

func generateMetaKeywords(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	return strings.Join(tags, ",")
}

func generateMetaDescription(content string) string {
	summary := generateSummary(content)
	if len(summary) > 160 {
		return summary[:160]
	}
	return summary
}

func removeMarkdown(content string) string {
	// 简单的markdown移除，实际需要更完善的解析
	result := content
	// 移除标题
	result = strings.ReplaceAll(result, "#", "")
	// 移除粗体
	result = strings.ReplaceAll(result, "**", "")
	result = strings.ReplaceAll(result, "__", "")
	// 移除斜体
	result = strings.ReplaceAll(result, "*", "")
	result = strings.ReplaceAll(result, "_", "")
	// 移除链接
	result = strings.ReplaceAll(result, "[", "")
	result = strings.ReplaceAll(result, "]", "")
	result = strings.ReplaceAll(result, "(", "")
	result = strings.ReplaceAll(result, ")", "")
	return strings.TrimSpace(result)
}