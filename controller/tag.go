package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetAllTagsHandler 获取所有标签的处理函数
// @Summary 获取所有标签
// @Description 获取系统中所有标签列表
// @Tags 标签
// @Accept json
// @Produce json
// @Success 200 {object} controller._ResponseTagList "标签列表"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/tags [get]
func GetAllTagsHandler(c *gin.Context) {
	// 1. 获取数据
	tags, err := logic.GetAllTags()
	if err != nil {
		zap.L().Error("logic.GetAllTags() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 2. 返回响应
	ResponseSuccess(c, tags)
}

// GetArticlesByTagHandler 获取标签下文章的处理函数
// @Summary 获取标签下文章
// @Description 获取指定标签下的文章列表
// @Tags 标签
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} controller._ResponseArticleList "文章列表"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/tags/{id}/articles [get]
func GetArticlesByTagHandler(c *gin.Context) {
	// 1. 获取参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	p := &models.ParamArticleList{
		Page: 1,
		Size: 20,
		Status: "published",
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetArticlesByTagHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取数据
	articles, total, err := logic.GetArticlesByTag(id, p)
	if err != nil {
		zap.L().Error("logic.GetArticlesByTag() failed", zap.Error(err))
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

// CreateTagHandler 创建标签的处理函数（需要作者权限）
// @Summary 创建标签
// @Description 创建新的文章标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param tag body models.ParamCreateTag true "标签参数"
// @Success 200 {object} controller._ResponseTag "创建的标签"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/tags [post]
func CreateTagHandler(c *gin.Context) {
	// 1. 获取参数
	p := new(models.ParamCreateTag)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreateTagHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 创建标签
	tag, err := logic.CreateTag(p)
	if err != nil {
		zap.L().Error("logic.CreateTag() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, tag)
}

// UpdateTagHandler 更新标签的处理函数（需要作者权限）
// @Summary 更新标签
// @Description 更新标签信息
// @Tags 标签
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param tag body models.ParamUpdateTag true "标签参数"
// @Success 200 {object} controller._ResponseSuccess "更新成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/tags/{id} [put]
func UpdateTagHandler(c *gin.Context) {
	// 1. 获取参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	p := new(models.ParamUpdateTag)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdateTagHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 更新标签
	if err := logic.UpdateTag(id, p); err != nil {
		zap.L().Error("logic.UpdateTag() failed", zap.Error(err))
		if err == mysql.ErrorTagNotExist {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// DeleteTagHandler 删除标签的处理函数（需要作者权限）
// @Summary 删除标签
// @Description 删除指定标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} controller._ResponseSuccess "删除成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/tags/{id} [delete]
func DeleteTagHandler(c *gin.Context) {
	// 1. 获取参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 删除标签
	if err := logic.DeleteTag(id); err != nil {
		zap.L().Error("logic.DeleteTag() failed", zap.Error(err))
		if err == mysql.ErrorTagNotExist {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetAuthorTagsHandler 获取作者标签的处理函数
// @Summary 获取作者标签
// @Description 获取当前登录作者的所有标签
// @Tags 标签
// @Accept json
// @Produce json
// @Success 200 {object} controller._ResponseTagList "标签列表"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/author/tags [get]
func GetAuthorTagsHandler(c *gin.Context) {
	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 获取作者标签
	tags, err := logic.GetAuthorTags(userID)
	if err != nil {
		zap.L().Error("logic.GetAuthorTags() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, tags)
}