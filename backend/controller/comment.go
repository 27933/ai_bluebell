package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateCommentHandler 创建评论的处理函数
// @Summary 创建评论
// @Description 创建新的评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param comment body models.ParamCreateComment true "评论参数"
// @Success 200 {object} controller._ResponseComment "创建的评论"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/comments [post]
func CreateCommentHandler(c *gin.Context) {
	// 1. 获取参数
	p := new(models.ParamCreateComment)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreateCommentHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.UserID = userID

	// 2. 创建评论
	comment, err := logic.CreateComment(p)
	if err != nil {
		zap.L().Error("logic.CreateComment() failed", zap.Error(err))
		if err == mysql.ErrorArticleNotExist {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, comment)
}

// GetCommentListHandler 获取评论列表的处理函数
// @Summary 获取评论列表
// @Description 获取指定文章的评论列表（访客可访问）
// @Tags 评论
// @Accept json
// @Produce json
// @Param article_id query int true "文章ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} controller._ResponseCommentList "评论列表"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/comments [get]
func GetCommentListHandler(c *gin.Context) {
	// 1. 获取参数
	articleIDStr := c.Query("article_id")
	if articleIDStr == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	p := &models.ParamCommentList{
		ArticleID: articleID,
		Page: 1,
		Size: 20,
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommentListHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取评论列表
	comments, total, err := logic.GetCommentList(p)
	if err != nil {
		zap.L().Error("logic.GetCommentList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"list": comments,
		"total": total,
		"page": p.Page,
		"size": p.Size,
		"pages": (total + int64(p.Size) - 1) / int64(p.Size),
	})
}

// UpdateCommentHandler 更新评论的处理函数
// @Summary 更新评论
// @Description 更新评论内容
// @Tags 评论
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Param comment body models.Comment true "评论内容"
// @Success 200 {object} controller._ResponseSuccess "更新成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/comments/{id} [put]
func UpdateCommentHandler(c *gin.Context) {
	// 1. 获取参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	p := new(models.Comment)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdateCommentHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 2. 更新评论
	if err := logic.UpdateComment(id, userID, p.Content); err != nil {
		zap.L().Error("logic.UpdateComment() failed", zap.Error(err))
		if err == mysql.ErrorCommentNotExist {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// DeleteCommentHandler 删除评论的处理函数
// @Summary 删除评论
// @Description 删除指定评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} controller._ResponseSuccess "删除成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/comments/{id} [delete]
func DeleteCommentHandler(c *gin.Context) {
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

	// 2. 删除评论
	if err := logic.DeleteComment(id, userID); err != nil {
		zap.L().Error("logic.DeleteComment() failed", zap.Error(err))
		if err == mysql.ErrorCommentNotExist {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}