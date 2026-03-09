package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"bluebell/logic"
	"bluebell/models"
)

// ExportArticleHandler 导出单篇文章
// @Summary 导出单篇文章
// @Description 导出单篇文章为Markdown格式
// @Tags 导出
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} controller._ResponseExportArticle "导出数据"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 404 {object} controller._ResponseError "文章不存在"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/author/articles/{id}/export [get]
func ExportArticleHandler(c *gin.Context) {
	// 获取文章ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户信息
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	userRole, _ := getCurrentUserRole(c)

	// 导出文章
	exportData, err := logic.ExportArticle(userID, userRole, articleID)
	if err != nil {
		if err == models.ErrorNoPermission {
			ResponseErrorWithMsg(c, CodeNoPermission, "只能导出自己的文章")
			return
		}
		if err == models.ErrorArticleNotExist {
			ResponseErrorWithMsg(c, CodeArticleNotExist, "文章不存在")
			return
		}
		zap.L().Error("logic.ExportArticle failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, exportData)
}

// ExportArticlesBatchHandler 批量导出文章
// @Summary 批量导出文章
// @Description 批量导出多篇文章为Markdown格式
// @Tags 导出
// @Accept json
// @Produce json
// @Param request body models.ExportArticlesBatchRequest true "导出请求"
// @Success 200 {object} controller._ResponseExportBatch "批量导出数据"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 404 {object} controller._ResponseError "文章不存在"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/author/articles/export [post]
func ExportArticlesBatchHandler(c *gin.Context) {
	// 获取请求参数
	var req models.ExportArticlesBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 验证文章ID数量
	if len(req.ArticleIDs) == 0 || len(req.ArticleIDs) > 50 {
		ResponseErrorWithMsg(c, CodeInvalidParam, "文章数量必须在1-50之间")
		return
	}

	// 获取当前用户信息
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	userRole, _ := getCurrentUserRole(c)

	// 批量导出文章
	exportData, err := logic.ExportArticlesBatch(userID, userRole, req.ArticleIDs)
	if err != nil {
		if err == models.ErrorNoPermission {
			ResponseErrorWithMsg(c, CodeNoPermission, "只能导出自己的文章")
			return
		}
		if err == models.ErrorArticleNotExist {
			ResponseErrorWithMsg(c, CodeArticleNotExist, "部分文章不存在")
			return
		}
		zap.L().Error("logic.ExportArticlesBatch failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, exportData)
}