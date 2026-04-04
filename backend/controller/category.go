package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"bluebell/pkg/ecode"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateCategoryHandler 创建栏目
// @Summary 创建栏目
// @Description 创建新的文章栏目
// @Tags 栏目
// @Accept json
// @Produce json
// @Param category body models.ParamCreateCategory true "栏目参数"
// @Success 200 {object} controller._ResponseSuccess "创建成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/categories [post]
func CreateCategoryHandler(c *gin.Context) {
	// 1. 获取参数
	p := new(models.ParamCreateCategory)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取当前用户ID和角色
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 3. 创建栏目
	if err := logic.CreateCategory(userID, p); err != nil {
		zap.L().Error("logic.CreateCategory failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 4. 返回响应
	ResponseSuccess(c, nil)
}

// GetCategoryListHandler 获取栏目列表
// @Summary 获取栏目列表
// @Description 获取所有栏目列表
// @Tags 栏目
// @Accept json
// @Produce json
// @Success 200 {object} controller._ResponseCategoryList "栏目列表"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/categories [get]
func GetCategoryListHandler(c *gin.Context) {
	// 1. 获取栏目列表
	categories, err := logic.GetCategoryList()
	if err != nil {
		zap.L().Error("logic.GetCategoryList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 2. 返回响应
	ResponseSuccess(c, gin.H{
		"list": categories,
	})
}

// GetCategoryByIdHandler 根据ID获取栏目详情
// @Summary 获取栏目详情
// @Description 根据ID获取栏目详细信息
// @Tags 栏目
// @Accept json
// @Produce json
// @Param id path int true "栏目ID"
// @Success 200 {object} controller._ResponseCategory "栏目详情"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 404 {object} controller._ResponseError "栏目不存在"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/categories/{id} [get]
func GetCategoryByIdHandler(c *gin.Context) {
	// 1. 获取参数
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取栏目详情
	category, err := logic.GetCategoryById(categoryID)
	if err != nil {
		if err == mysql.ErrorCategoryNotExist {
			ResponseError(c, CodeCategoryNotExist)
			return
		}
		zap.L().Error("logic.GetCategoryById failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, category)
}

// UpdateCategoryHandler 更新栏目
// @Summary 更新栏目
// @Description 更新栏目信息
// @Tags 栏目
// @Accept json
// @Produce json
// @Param id path int true "栏目ID"
// @Param category body models.ParamUpdateCategory true "栏目参数"
// @Success 200 {object} controller._ResponseSuccess "更新成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 404 {object} controller._ResponseError "栏目不存在"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/categories/{id} [put]
func UpdateCategoryHandler(c *gin.Context) {
	// 1. 获取参数
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	p := new(models.ParamUpdateCategory)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取当前用户ID和角色
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	userRole, _ := getCurrentUserRole(c)

	// 3. 更新栏目
	if err := logic.UpdateCategory(userID, userRole, categoryID, p); err != nil {
		zap.L().Error("logic.UpdateCategory failed", zap.Error(err))
		// 权限错误统一返回服务器繁忙，不暴露具体原因
		if err == ecode.ErrNoPermissionCategory || err == ecode.ErrNoPermission {
			ResponseError(c, CodeServerBusy)
			return
		}
		if err == mysql.ErrorCategoryExist {
			ResponseError(c, CodeCategoryExist)
			return
		}
		if err == mysql.ErrorCategoryNotExist {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 4. 返回响应
	ResponseSuccess(c, nil)
}

// DeleteCategoryHandler 删除栏目
// @Summary 删除栏目
// @Description 删除指定栏目
// @Tags 栏目
// @Accept json
// @Produce json
// @Param id path int true "栏目ID"
// @Success 200 {object} controller._ResponseSuccess "删除成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/categories/{id} [delete]
func DeleteCategoryHandler(c *gin.Context) {
	// 1. 获取参数
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取当前用户ID和角色
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	userRole, _ := getCurrentUserRole(c)

	// 3. 删除栏目
	if err := logic.DeleteCategory(userID, userRole, categoryID); err != nil {
		zap.L().Error("logic.DeleteCategory failed", zap.Error(err))
		// 权限错误统一返回服务器繁忙，不暴露具体原因
		if err == ecode.ErrNoPermissionCategory || err == ecode.ErrNoPermission {
			ResponseError(c, CodeServerBusy)
			return
		}
		if err == mysql.ErrorCategoryNotExist {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 4. 返回响应
	ResponseSuccess(c, nil)
}

// AddArticleToCategoriesHandler 添加文章到栏目
// @Summary 添加文章到栏目
// @Description 将文章添加到指定栏目
// @Tags 栏目
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param categories body models.ParamAddArticleCategories true "栏目ID列表"
// @Success 200 {object} controller._ResponseSuccess "添加成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/articles/{id}/categories [post]
func AddArticleToCategoriesHandler(c *gin.Context) {
	// 1. 获取文章ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取参数
	p := new(models.ParamAddArticleCategories)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 3. 获取当前用户ID和角色
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	userRole, _ := getCurrentUserRole(c)

	// 4. 添加文章到栏目
	if err := logic.AddArticleToCategories(userID, userRole, articleID, p.CategoryIDs); err != nil {
		zap.L().Error("logic.AddArticleToCategories failed", zap.Error(err))
		if err == mysql.ErrorCategoryNotExist {
			ResponseError(c, CodeCategoryNotExist)
			return
		}
		// 权限错误统一返回服务器繁忙，不暴露具体原因
		if err == ecode.ErrNoPermissionCategory || err == ecode.ErrNoPermission {
			ResponseError(c, CodeServerBusy)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 5. 返回响应
	ResponseSuccess(c, nil)
}

// GetArticlesByCategoryHandler 获取栏目下的文章
// @Summary 获取栏目下的文章
// @Description 获取指定栏目下的文章列表
// @Tags 栏目
// @Accept json
// @Produce json
// @Param id path int true "栏目ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} controller._ResponseArticleList "文章列表"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 404 {object} controller._ResponseError "栏目不存在"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/categories/{id}/articles [get]
func GetArticlesByCategoryHandler(c *gin.Context) {
	// 1. 获取栏目ID
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取分页参数
	page, size := getPageInfo(c)

	// 3. 获取文章ID列表
	articleIDs, total, err := logic.GetArticlesByCategory(categoryID, int(page), int(size))
	if err != nil {
		if err == mysql.ErrorCategoryNotExist {
			ResponseError(c, CodeCategoryNotExist)
			return
		}
		zap.L().Error("logic.GetArticlesByCategory failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 4. 获取文章详情列表
	articles := make([]*models.ApiArticleListItem, 0, len(articleIDs))
	now := time.Now()
	for _, articleID := range articleIDs {
		article, err := mysql.GetArticleById(articleID)
		if err != nil {
			zap.L().Error("mysql.GetArticleById failed", zap.Error(err))
			continue
		}
		// 只返回公开文章
		if article.Status == string(models.ArticleStatusPublished) {
			// 转换为ApiArticleListItem格式
			apiArticle := &models.ApiArticleListItem{
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
			articles = append(articles, apiArticle)
		}
	}

	// 5. 返回响应
	ResponseSuccess(c, gin.H{
		"list":  articles,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetCategoriesByArticleHandler 获取文章的栏目列表
// @Summary 获取文章的栏目列表
// @Description 获取指定文章所属的栏目列表
// @Tags 栏目
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} controller._ResponseCategoryList "栏目列表"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/articles/{id}/categories [get]
func GetCategoriesByArticleHandler(c *gin.Context) {
	// 1. 获取文章ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取栏目列表
	categories, err := logic.GetCategoriesByArticle(articleID)
	if err != nil {
		if err == mysql.ErrorArticleNotExist {
			ResponseError(c, CodeInvalidParam)
			return
		}
		zap.L().Error("logic.GetCategoriesByArticle failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"list": categories,
	})
}