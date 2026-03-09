package controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"bluebell/logic"
	"bluebell/models"
)

// GetAuthorInfoHandler 获取作者信息
// @Summary 获取作者信息
// @Description 获取指定作者的基本信息
// @Tags 作者
// @Accept json
// @Produce json
// @Param username path string true "作者用户名"
// @Success 200 {object} controller._ResponseAuthorInfo "作者信息"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 404 {object} controller._ResponseError "用户不存在"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/authors/{username} [get]
func GetAuthorInfoHandler(c *gin.Context) {
	// 获取用户名参数
	username := c.Param("username")
	if username == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 调用逻辑层获取作者信息
	authorInfo, err := logic.GetAuthorInfo(username)
	if err != nil {
		zap.L().Error("logic.GetAuthorInfo failed", zap.String("username", username), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 用户不存在
	if authorInfo == nil {
		ResponseErrorWithMsg(c, CodeInvalidParam, "用户不存在")
		return
	}

	ResponseSuccess(c, authorInfo)
}

// GetAuthorArticlesListHandler 获取作者文章列表
// @Summary 获取作者文章列表
// @Description 获取指定作者的文章列表，支持按时间或热度排序
// @Tags 作者
// @Accept json
// @Produce json
// @Param username path string true "作者用户名"
// @Param sort query string false "排序方式" default(time) enums(time,hot)
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} controller._ResponseAuthorArticles "作者文章列表"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/authors/{username}/articles [get]
func GetAuthorArticlesListHandler(c *gin.Context) {
	// 参数获取和验证
	username := c.Param("username")
	if username == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 分页参数
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	sort := c.Query("sort") // time 或 hot

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 20
	}

	// 排序方式验证
	if sort != "time" && sort != "hot" {
		sort = "time"
	}

	// 构建参数
	params := &models.ParamAuthorArticlesList{
		Username: username,
		Page:     page,
		Size:     size,
		Sort:     sort,
	}

	// 参数校验
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取作者文章列表
	articles, total, err := logic.GetAuthorArticlesList(params)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 处理最近更新标记
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

	// 返回处理后的列表
	ResponseSuccess(c, gin.H{
		"list": articleList,
		"total": total,
		"page": page,
		"size": size,
		"pages": (total + int64(size) - 1) / int64(size),
	})
}