package ecode

import "errors"

// 定义所有业务错误
var (
	// 通用错误
	ErrServerBusy = errors.New("服务繁忙")

	// 用户相关错误
	ErrUserNotLogin = errors.New("用户未登录")
	ErrUserExist    = errors.New("用户已存在")
	ErrUserNotExist = errors.New("用户不存在")
	ErrInvalidPassword = errors.New("用户名或密码错误")
	ErrNoPermission = errors.New("无权限操作")

	// 文章相关错误
	ErrArticleNotExist = errors.New("文章不存在")

	// 栏目相关错误
	ErrCategoryNotExist = errors.New("栏目不存在")
	ErrCategoryExist  = errors.New("栏目已存在")
	ErrNoPermissionCategory = errors.New("无权限操作该栏目")
)