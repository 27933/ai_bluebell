package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UpdateArticleHandler 更新文章的处理函数
// @Summary 更新文章
// @Description 作者更新自己的文章内容
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param article body models.ParamUpdateArticle true "文章更新参数"
// @Success 200 {object} controller._ResponseSuccess "更新成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/author/articles/{id} [put]
func UpdateArticleHandler(c *gin.Context) {
	// 1. 获取参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	p := new(models.ParamUpdateArticle)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdateArticleHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 更新文章
	if err := logic.UpdateArticle(userID, id, p); err != nil {
		zap.L().Error("logic.UpdateArticle() failed", zap.Error(err))
		if err == mysql.ErrorArticleNotExist {
			ResponseError(c, CodeArticleNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// DeleteArticleHandler 删除文章的处理函数
// @Summary 删除文章
// @Description 作者删除自己的文章
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} controller._ResponseSuccess "删除成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/author/articles/{id} [delete]
func DeleteArticleHandler(c *gin.Context) {
	// 1. 获取参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 2. 删除文章
	if err := logic.DeleteArticle(userID, id); err != nil {
		zap.L().Error("logic.DeleteArticle() failed", zap.Error(err))
		if err == mysql.ErrorArticleNotExist {
			ResponseError(c, CodeArticleNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// UpdateArticleStatusHandler 更新文章状态的处理函数
// @Summary 更新文章状态
// @Description 作者更新自己文章的状态（草稿/发布/下线）
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param status body models.ParamArticleStatus true "文章状态"
// @Success 200 {object} controller._ResponseSuccess "更新成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/author/articles/{id}/status [patch]
func UpdateArticleStatusHandler(c *gin.Context) {
	// 1. 获取参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	p := new(models.ParamArticleStatus)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdateArticleStatusHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 更新文章状态
	if err := logic.UpdateArticleStatus(userID, id, p.Status); err != nil {
		zap.L().Error("logic.UpdateArticleStatus() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// UpdateArticleFeaturedHandler 更新文章精选状态的处理函数
// @Summary 更新文章精选状态
// @Description 作者更新自己文章的精选状态
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param featured body models.ParamArticleFeatured true "精选状态"
// @Success 200 {object} controller._ResponseSuccess "更新成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/author/articles/{id}/featured [patch]
func UpdateArticleFeaturedHandler(c *gin.Context) {
	// 1. 获取参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	p := new(models.ParamArticleFeatured)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdateArticleFeaturedHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 更新文章精选状态
	if err := logic.UpdateArticleFeatured(userID, id, p.IsFeatured); err != nil {
		zap.L().Error("logic.UpdateArticleFeatured() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetAuthorArticlesHandler 获取作者文章列表的处理函数
// @Summary 获取作者文章列表
// @Description 获取当前登录作者的文章列表
// @Tags 文章
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Param status query string false "文章状态" enums(all,published,draft,offline)
// @Success 200 {object} controller._ResponseArticleList "文章列表"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/author/articles [get]
func GetAuthorArticlesHandler(c *gin.Context) {
	// 1. 获取参数
	p := &models.ParamArticleList{
		Page: 1,
		Size: 20,
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetAuthorArticlesHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 2. 获取作者文章列表
	articles, total, err := logic.GetAuthorArticles(userID, p)
	if err != nil {
		zap.L().Error("logic.GetAuthorArticles() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"list": articles,
		"total": total,
		"page": p.Page,
		"size": p.Size,
		"pages": (total + int64(p.Size) - 1) / int64(p.Size),
	})
}

// GetAdminArticlesHandler 管理员获取文章列表的处理函数
// @Summary 管理员获取文章列表
// @Description 管理员获取所有文章列表，支持状态筛选
// @Tags 管理员
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Param status query string false "文章状态" enums(all,published,draft,offline)
// @Param author query string false "作者用户名"
// @Success 200 {object} controller._ResponseArticleList "文章列表"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/admin/articles [get]
func GetAdminArticlesHandler(c *gin.Context) {
	// 1. 获取参数
	p := &models.ParamArticleList{
		Page: 1,
		Size: 20,
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetAdminArticlesHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取所有文章列表（管理员权限）
	articles, total, err := logic.GetAdminArticles(p)
	if err != nil {
		zap.L().Error("logic.GetAdminArticles() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"list": articles,
		"total": total,
		"page": p.Page,
		"size": p.Size,
		"pages": (total + int64(p.Size) - 1) / int64(p.Size),
	})
}

// AdminSetArticleFeaturedHandler 管理员设置文章精选状态的处理函数
// @Summary 管理员设置文章精选
// @Description 管理员设置文章的精选状态
// @Tags 管理员
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param featured body models.ParamArticleFeatured true "精选状态"
// @Success 200 {object} controller._ResponseSuccess "设置成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/admin/articles/{id}/featured [patch]
func AdminSetArticleFeaturedHandler(c *gin.Context) {
	// 1. 获取参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	p := new(models.ParamArticleFeatured)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("AdminSetArticleFeaturedHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 设置文章精选状态
	if err := logic.AdminSetArticleFeatured(id, p.IsFeatured); err != nil {
		zap.L().Error("logic.AdminSetArticleFeatured() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// CreateArticleHandler 创建文章的处理函数
// @Summary 创建文章
// @Description 作者创建新文章
// @Tags 文章
// @Accept json
// @Produce json
// @Param article body models.ParamCreateArticle true "文章参数"
// @Success 200 {object} controller._ResponseCreateArticle "创建成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/articles [post]
func CreateArticleHandler(c *gin.Context) {
	// 1. 获取参数及参数校验
	p := new(models.ParamCreateArticle)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create article with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 从 c 取到当前发请求的用户的ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 检查权限：只有 author 和 admin 角色可以创建文章
	userRole, err := getCurrentUserRole(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	if userRole != "author" && userRole != "admin" {
		ResponseError(c, CodeNoPermission)
		return
	}

	// 2. 创建文章
	article, err := logic.CreateArticle(userID, p)
	if err != nil {
		zap.L().Error("logic.CreateArticle failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, article)
}

// GetArticleDetailHandler 获取文章详情的处理函数
// @Summary 获取文章详情
// @Description 获取指定文章的详细信息
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} controller._ResponseArticleDetail "文章详情"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 404 {object} controller._ResponseError "文章不存在"
// @Router /api/v1/articles/{id} [get]
func GetArticleDetailHandler(c *gin.Context) {
	// 1. 获取参数（从URL中获取文章的id）
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("get article detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取用户ID（如果有）
	var userID *int64
	if currentUserID, err := getCurrentUserID(c); err == nil {
		userID = &currentUserID
	}

	// 获取客户端IP
	ip := c.ClientIP()

	// 2. 根据id取出文章数据
	article, err := logic.GetArticleDetail(id)
	if err != nil {
		zap.L().Error("logic.GetArticleDetail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 异步记录访问（带防刷机制）
	go func() {
		if err := logic.RecordArticleViewWithAntiCheat(id, userID, ip); err != nil {
			zap.L().Error("RecordArticleViewWithAntiCheat failed", zap.Error(err))
		}
	}()

	// 4. 返回响应
	ResponseSuccess(c, article)
}

// GetArticleListHandler 获取文章列表的处理函数
// @Summary 获取文章列表
// @Description 获取公开的文章列表，支持分页和状态筛选
// @Tags 文章
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Param status query string false "文章状态" default(published) enums(published,draft,offline)
// @Param sort query string false "排序方式" default(time) enums(time,hot)
// @Param tag query string false "标签"
// @Param category_id query int false "栏目ID"
// @Success 200 {object} controller._ResponseArticleList "文章列表"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/articles [get]
func GetArticleListHandler(c *gin.Context) {
	// 1. 获取参数
	p := &models.ParamArticleList{
		Page: 1,
		Size: 20,
		Sort: "time",
		Status: "published", // 默认只显示已发布的文章
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetArticleListHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取数据
	articles, total, err := logic.GetArticleList(p)
	if err != nil {
		zap.L().Error("logic.GetArticleList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 处理最近更新标记
	now := time.Now()
	articleList := make([]*models.ApiArticleListItem, 0, len(articles))
	for _, article := range articles {
		item := &models.ApiArticleListItem{
			ID:           article.ID,
			Title:        article.Title,
			Summary:      article.Summary,
			ViewCount:    article.ViewCount,
			LikeCount:    article.LikeCount,
			CommentCount: article.CommentCount,
			IsFeatured:   article.IsFeatured,
			IsRecent:     article.UpdatedAt.After(now.Add(-24 * time.Hour)), // 24小时内更新
			CreatedAt:    article.CreatedAt,
			UpdatedAt:    article.UpdatedAt,
		}
		articleList = append(articleList, item)
	}

	// 4. 返回响应
	ResponseSuccess(c, gin.H{
		"list": articleList,
		"total": total,
		"page": p.Page,
		"size": p.Size,
		"pages": (total + int64(p.Size) - 1) / int64(p.Size),
	})
}

// GetFeaturedArticlesHandler 获取精选文章的处理函数
// @Summary 获取精选文章
// @Description 获取精选文章列表
// @Tags 文章
// @Accept json
// @Produce json
// @Param limit query int false "数量限制" default(3) minimum(1) maximum(10)
// @Success 200 {object} controller._ResponseFeaturedArticles "精选文章列表"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/articles/featured [get]
func GetFeaturedArticlesHandler(c *gin.Context) {
	// 1. 获取参数
	limitStr := c.DefaultQuery("limit", "3")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 3
	}
	if limit > 10 {
		limit = 10
	}

	// 2. 获取精选文章
	articles, err := logic.GetFeaturedArticles(limit)
	if err != nil {
		zap.L().Error("logic.GetFeaturedArticles() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, articles)
}

// SearchArticlesHandler 搜索文章的处理函数
// @Summary 搜索文章
// @Description 根据关键词、作者、标签搜索文章
// @Tags 文章
// @Accept json
// @Produce json
// @Param keyword query string false "搜索关键词"
// @Param author_name query string false "作者用户名"
// @Param tag query string false "标签"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} controller._ResponseArticleList "搜索结果"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/articles/search [get]
func SearchArticlesHandler(c *gin.Context) {
	// 1. 获取参数
	p := &models.ParamArticleList{
		Page: 1,
		Size: 20,
		Status: "published",
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("SearchArticlesHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 验证搜索条件：关键词、作者ID、作者名、标签至少有一个
	if p.Keyword == "" && p.AuthorID == 0 && p.AuthorName == "" && p.Tag == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 搜索文章
	articles, total, err := logic.SearchArticles(p)
	if err != nil {
		zap.L().Error("logic.SearchArticles() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"list": articles,
		"total": total,
		"page": p.Page,
		"size": p.Size,
		"pages": (total + int64(p.Size) - 1) / int64(p.Size),
	})
}

// RecordArticleViewHandler 记录文章访问的处理函数
// @Summary 记录文章访问
// @Description 记录文章访问（已废弃，使用防刷版本）
// @Tags 文章
// @Accept json
// @Produce json
// @Param article_id query int true "文章ID"
// @Success 200 {object} controller._ResponseSuccess "记录成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/articles/view [post]
func RecordArticleViewHandler(c *gin.Context) {
	// 1. 获取参数（现在使用查询参数article_id）
	idStr := c.Query("article_id")
	if idStr == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 记录访问
	if err := logic.RecordArticleView(id); err != nil {
		zap.L().Error("logic.RecordArticleView failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetArticleDailyStatsHandler 获取文章日访问量统计的处理函数
// @Summary 获取文章日访问量统计
// @Description 获取指定文章的日访问量统计数据
// @Tags 统计
// @Accept json
// @Produce json
// @Param article_id query int true "文章ID"
// @Param days query int false "统计天数" default(7) minimum(1) maximum(90)
// @Success 200 {object} controller._ResponseArticleDailyStats "日访问统计数据"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/article-stats/daily [get]
func GetArticleDailyStatsHandler(c *gin.Context) {
	// 1. 获取参数（现在使用查询参数article_id）
	idStr := c.Query("article_id")
	if idStr == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}
	if days > 90 {
		days = 90
	}

	// 2. 获取统计数据
	stats, err := logic.GetArticleDailyStats(id, days)
	if err != nil {
		zap.L().Error("logic.GetArticleDailyStats failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, stats)
}