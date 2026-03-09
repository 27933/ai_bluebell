package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
	ErrorInvalidID       = errors.New("无效的ID")
	ErrorArticleNotExist = errors.New("文章不存在")
	ErrorTagNotExist     = errors.New("标签不存在")
	ErrorTagInUse        = errors.New("标签正在被使用")
	ErrorCommentNotExist = errors.New("评论不存在")
	ErrorLikeNotExist    = errors.New("点赞记录不存在")

	// 栏目相关错误
	ErrorCategoryNotExist = errors.New("栏目不存在")
	ErrorCategoryExist    = errors.New("栏目已存在")
	ErrorCategoryInUse    = errors.New("栏目正在被使用")

	// 导出相关错误
	ErrorNoPermission = errors.New("没有权限")
	ErrorInvalidParam = errors.New("参数错误")
)
