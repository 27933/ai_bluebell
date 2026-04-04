package api

import (
	"fmt"
	"strconv"
	"testing"
)

// TestComment_Create_NotLogin 测试未登录创建评论
func TestComment_Create_NotLogin(t *testing.T) {
	body := map[string]interface{}{
		"article_id": 1,
		"content":    "测试评论内容",
	}

	resp, statusCode, err := DoPost("/comments", body, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回需要登录错误
	if !AssertError(t, resp, statusCode, 1006) {
		t.Logf("响应: %+v", resp)
	}
}

// TestComment_Create_Success 测试成功创建评论
func TestComment_Create_Success(t *testing.T) {
	// 创建评论者用户
	commenterToken := createTestUserAndLogin(t)
	if commenterToken == "" {
		t.Skip("无法获取评论者Token，跳过测试")
		return
	}

	// 创建文章作者和文章
	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取作者Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 创建评论
	body := map[string]interface{}{
		"article_id": articleID,
		"content":    "这是一条测试评论",
	}

	resp, statusCode, err := DoPost("/comments", body, commenterToken)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
		return
	}

	// 验证返回的评论数据
	data, err := ExtractData(resp)
	if err != nil {
		t.Fatalf("提取数据失败: %v", err)
	}

	if data["id"] == nil {
		t.Error("响应缺少 id 字段")
	}
	if data["content"] == nil {
		t.Error("响应缺少 content 字段")
	}
	if data["article_id"] == nil {
		t.Error("响应缺少 article_id 字段")
	}
}

// TestComment_Create_ArticleNotExist 测试评论不存在的文章
func TestComment_Create_ArticleNotExist(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	body := map[string]interface{}{
		"article_id": 999999999,
		"content":    "测试评论内容",
	}

	resp, statusCode, err := DoPost("/comments", body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回参数错误
	if !AssertError(t, resp, statusCode, 1001) {
		t.Logf("响应: %+v", resp)
	}
}

// TestComment_List_VisitorAccess 测试访客可以访问评论列表
func TestComment_List_VisitorAccess(t *testing.T) {
	// 创建作者和文章
	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 添加几条评论
	createTestComment(t, authorToken, articleID, "评论1")
	createTestComment(t, authorToken, articleID, "评论2")

	// 获取评论列表（访客可访问，不需要登录）
	resp, statusCode, err := DoGetWithQuery(fmt.Sprintf("/comments?article_id=%d", articleID), "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
		return
	}

	data, err := ExtractData(resp)
	if err != nil {
		t.Fatalf("提取数据失败: %v", err)
	}

	if data["list"] == nil {
		t.Error("响应缺少 list 字段")
	}
	if data["total"] == nil {
		t.Error("响应缺少 total 字段")
	}
}

// TestComment_List_Success 测试获取评论列表（已登录用户）
func TestComment_List_Success(t *testing.T) {
	// 创建作者和文章
	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 添加几条评论
	createTestComment(t, authorToken, articleID, "评论1")
	createTestComment(t, authorToken, articleID, "评论2")

	// 获取评论列表（已登录用户）
	resp, statusCode, err := DoGetWithQuery(fmt.Sprintf("/comments?article_id=%d", articleID), authorToken)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
		return
	}

	data, err := ExtractData(resp)
	if err != nil {
		t.Fatalf("提取数据失败: %v", err)
	}

	if data["list"] == nil {
		t.Error("响应缺少 list 字段")
	}
	if data["total"] == nil {
		t.Error("响应缺少 total 字段")
	}
}

// TestComment_List_MissingArticleID 测试评论列表缺少文章ID
func TestComment_List_MissingArticleID(t *testing.T) {
	// GET /comments 现在是公开接口，不需要登录也可以访问
	resp, statusCode, err := DoGet("/comments", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回参数错误
	if !AssertError(t, resp, statusCode, 1001) {
		t.Logf("响应: %+v", resp)
	}
}

// TestComment_Update_NotAuthor 测试非作者更新评论
func TestComment_Update_NotAuthor(t *testing.T) {
	// 创建文章和评论者1
	user1Token := createTestUserAndLogin(t)
	if user1Token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取作者Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	commentID := createTestComment(t, user1Token, articleID, "原评论内容")
	if commentID == 0 {
		t.Skip("无法创建测试评论，跳过测试")
		return
	}

	// 用户2尝试更新评论
	user2Token := createTestUserAndLogin(t)
	if user2Token == "" {
		t.Skip("无法获取用户2Token，跳过测试")
		return
	}

	body := map[string]string{
		"content": "被篡改的评论",
	}

	resp, statusCode, err := DoPut(fmt.Sprintf("/comments/%d", commentID), body, user2Token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回权限错误
	if resp.Code == 1000 {
		t.Error("非作者不应该能更新评论")
	}
	_ = statusCode
}

// TestComment_Update_Success 测试作者成功更新评论
func TestComment_Update_Success(t *testing.T) {
	// 创建文章和评论
	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	commentID := createTestComment(t, authorToken, articleID, "原评论内容")
	if commentID == 0 {
		t.Skip("无法创建测试评论，跳过测试")
		return
	}

	newContent := "更新后的评论内容"
	body := map[string]string{
		"content": newContent,
	}

	resp, statusCode, err := DoPut(fmt.Sprintf("/comments/%d", commentID), body, authorToken)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestComment_Delete_NotAuthor 测试非作者删除评论
func TestComment_Delete_NotAuthor(t *testing.T) {
	// 创建文章和评论者1
	user1Token := createTestUserAndLogin(t)
	if user1Token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取作者Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	commentID := createTestComment(t, user1Token, articleID, "评论内容")
	if commentID == 0 {
		t.Skip("无法创建测试评论，跳过测试")
		return
	}

	// 用户2尝试删除评论
	user2Token := createTestUserAndLogin(t)
	if user2Token == "" {
		t.Skip("无法获取用户2Token，跳过测试")
		return
	}

	resp, statusCode, err := DoDelete(fmt.Sprintf("/comments/%d", commentID), user2Token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回权限错误
	if resp.Code == 1000 {
		t.Error("非作者不应该能删除评论")
	}
	_ = statusCode
}

// TestComment_Delete_Success 测试作者成功删除评论
func TestComment_Delete_Success(t *testing.T) {
	// 创建文章和评论
	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	commentID := createTestComment(t, authorToken, articleID, "将被删除的评论")
	if commentID == 0 {
		t.Skip("无法创建测试评论，跳过测试")
		return
	}

	resp, statusCode, err := DoDelete(fmt.Sprintf("/comments/%d", commentID), authorToken)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// createTestComment 辅助函数：创建测试评论
func createTestComment(t *testing.T, token string, articleID int64, content string) int64 {
	t.Helper()

	body := map[string]interface{}{
		"article_id": articleID,
		"content":    content,
	}

	resp, _, err := DoPost("/comments", body, token)
	if err != nil {
		t.Logf("创建评论失败: %v", err)
		return 0
	}

	if resp.Code != 1000 {
		t.Logf("创建评论失败: %v", resp.Msg)
		return 0
	}

	data, err := ExtractData(resp)
	if err != nil {
		t.Logf("提取数据失败: %v", err)
		return 0
	}

	var commentID int64
	switch v := data["id"].(type) {
	case float64:
		commentID = int64(v)
	case int64:
		commentID = v
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			t.Logf("无法解析评论ID字符串: %v", err)
			return 0
		}
		commentID = parsed
	default:
		t.Logf("无法解析评论ID类型: %T", data["id"])
		return 0
	}

	return commentID
}
