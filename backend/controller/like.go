package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LikeRequest JSON请求参数结构体
type LikeRequest struct {
	TargetType string `json:"target_type" binding:"required,oneof=article comment"`
	TargetID   int64  `json:"target_id" binding:"required"`
}

// LikeHandler 点赞的处理函数（支持查询参数和JSON请求体）
// @Summary 点赞
// @Description 对文章或评论点赞
// @Tags 点赞
// @Accept json
// @Produce json
// @Param like body controller.LikeRequest true "点赞参数"
// @Success 200 {object} controller._ResponseSuccess "点赞成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/likes [post]
func LikeHandler(c *gin.Context) {
	var targetType string
	var targetID int64

	// 优先尝试JSON请求体
	req := &LikeRequest{}
	if err := c.ShouldBindJSON(req); err == nil {
		// JSON请求体成功
		targetType = req.TargetType
		targetID = req.TargetID
	} else {
		// 回退到查询参数（向后兼容）
		targetType = c.Query("target_type")
		if targetType != "article" && targetType != "comment" {
			ResponseError(c, CodeInvalidParam)
			return
		}

		targetIDStr := c.Query("target_id")
		var err error
		targetID, err = strconv.ParseInt(targetIDStr, 10, 64)
		if err != nil {
			ResponseError(c, CodeInvalidParam)
			return
		}
	}

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 2. 点赞
	if err := logic.Like(userID, models.TargetType(targetType), targetID); err != nil {
		zap.L().Error("logic.Like() failed", zap.Error(err))
		if err == mysql.ErrorArticleNotExist || err == mysql.ErrorCommentNotExist {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// UnlikeHandler 取消点赞的处理函数
// @Summary 取消点赞
// @Description 取消对文章或评论的点赞
// @Tags 点赞
// @Accept json
// @Produce json
// @Param target_type query string true "目标类型" enums(article,comment)
// @Param target_id query int true "目标ID"
// @Success 200 {object} controller._ResponseSuccess "取消成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/likes [delete]
func UnlikeHandler(c *gin.Context) {
	// 1. 获取参数
	targetType := c.Query("target_type")
	if targetType != "article" && targetType != "comment" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	targetIDStr := c.Query("target_id")
	targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
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

	// 2. 取消点赞
	if err := logic.Unlike(userID, models.TargetType(targetType), targetID); err != nil {
		zap.L().Error("logic.Unlike() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetLikeStatusHandler 获取点赞状态的处理函数
// @Summary 获取点赞状态
// @Description 获取用户对文章或评论的点赞状态
// @Tags 点赞
// @Accept json
// @Produce json
// @Param target_type query string true "目标类型" enums(article,comment)
// @Param target_id query int true "目标ID"
// @Success 200 {object} controller._ResponseLikeStatus "点赞状态"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/likes/status [get]
func GetLikeStatusHandler(c *gin.Context) {
	// 1. 获取参数
	targetType := c.Query("target_type")
	if targetType != "article" && targetType != "comment" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	targetIDStr := c.Query("target_id")
	targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户ID（可选）
	userID, _ := getCurrentUserID(c)

	// 2. 获取点赞状态
	status, err := logic.GetLikeStatus(userID, models.TargetType(targetType), targetID)
	if err != nil {
		zap.L().Error("logic.GetLikeStatus() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, status)
}

// GetUserLikesHandler 获取用户点赞列表的处理函数
// @Summary 获取用户点赞列表
// @Description 获取当前用户的点赞列表
// @Tags 点赞
// @Accept json
// @Produce json
// @Param target_type query string true "目标类型" enums(article,comment)
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} controller._ResponseLikeList "点赞列表"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/user/likes [get]
func GetUserLikesHandler(c *gin.Context) {
	// 1. 获取参数
	targetType := c.Query("target_type")
	if targetType != "article" && targetType != "comment" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 2. 获取点赞列表
	likes, total, err := logic.GetUserLikes(userID, models.TargetType(targetType), page, size)
	if err != nil {
		zap.L().Error("logic.GetUserLikes() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"list": likes,
		"total": total,
		"page": page,
		"size": size,
		"pages": (total + int64(size) - 1) / int64(size),
	})
}