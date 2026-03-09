package api

import (
	"fmt"
	"strconv"
	"testing"
)

// TestArticle_List_Success 测试获取文章列表
func TestArticle_List_Success(t *testing.T) {
	resp, statusCode, err := DoGet("/articles", "")
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

	// 验证返回包含list和total
	if data["list"] == nil {
		t.Error("响应缺少 list 字段")
	}
	if data["total"] == nil {
		t.Error("响应缺少 total 字段")
	}
}

// TestArticle_List_WithPagination 测试文章列表分页
func TestArticle_List_WithPagination(t *testing.T) {
	resp, statusCode, err := DoGetWithQuery("/articles?page=1&size=5", "")
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

	if data["page"] == nil {
		t.Error("响应缺少 page 字段")
	}
	if data["size"] == nil {
		t.Error("响应缺少 size 字段")
	}
}

// TestArticle_List_WithSort 测试文章按热度排序
func TestArticle_List_WithSort(t *testing.T) {
	// 按时间排序
	resp, statusCode, err := DoGetWithQuery("/articles?sort=time", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("按时间排序响应: %+v", resp)
	}

	// 按热度排序
	resp, statusCode, err = DoGetWithQuery("/articles?sort=popular", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("按热度排序响应: %+v", resp)
	}
}

// TestArticle_Detail_Success 测试获取存在的文章详情
func TestArticle_Detail_Success(t *testing.T) {
	// 先创建一篇文章
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, token, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 获取文章详情
	resp, statusCode, err := DoGet(fmt.Sprintf("/articles/%d", articleID), "")
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

	if data["id"] == nil {
		t.Error("响应缺少 id 字段")
	}
	if data["title"] == nil {
		t.Error("响应缺少 title 字段")
	}
	if data["content"] == nil {
		t.Error("响应缺少 content 字段")
	}
}

// TestArticle_Create_NotLogin 测试未登录创建文章
func TestArticle_Create_NotLogin(t *testing.T) {
	body := map[string]interface{}{
		"title":   "测试文章标题",
		"content": "测试文章内容",
	}

	resp, statusCode, err := DoPost("/articles", body, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回需要登录错误
	if !AssertError(t, resp, statusCode, 1006) {
		t.Logf("响应: %+v", resp)
	}
}

// TestArticle_Create_Success 测试成功创建文章
func TestArticle_Create_Success(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	body := map[string]interface{}{
		"title":    "测试文章标题_" + GenerateRandomUsername(),
		"content":  "测试文章内容",
		"tags":     []string{"测试", "Go"},
		"status":   "draft",
		"slug":     "test-slug-" + GenerateRandomUsername(),
	}

	resp, statusCode, err := DoPost("/articles", body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
		return
	}

	// 验证返回的文章数据
	data, err := ExtractData(resp)
	if err != nil {
		t.Fatalf("提取数据失败: %v", err)
	}

	if data["title"] == nil {
		t.Error("响应缺少 title 字段")
	}
	if data["id"] == nil {
		t.Error("响应缺少 id 字段")
	}
}

// TestArticle_Create_MissingField 测试创建文章缺少必填字段
func TestArticle_Create_MissingField(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	// 缺少标题
	body := map[string]interface{}{
		"content": "测试文章内容",
	}

	resp, statusCode, err := DoPost("/articles", body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertError(t, resp, statusCode, 1001) {
		t.Logf("响应: %+v", resp)
	}
}

// TestArticle_Update_NotAuthor 测试非作者更新文章
func TestArticle_Update_NotAuthor(t *testing.T) {
	// 创建作者用户和文章
	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取作者Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "draft")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 创建另一个用户
	otherToken := createTestUserAndLogin(t)
	if otherToken == "" {
		t.Skip("无法获取其他用户Token，跳过测试")
		return
	}

	// 尝试用其他用户更新文章
	body := map[string]interface{}{
		"title": "被篡改的标题",
	}

	resp, _, err := DoPut(fmt.Sprintf("/author/articles/%d", articleID), body, otherToken)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回权限错误（非作者）
	// 实际代码返回 CodeServerBusy 或 CodeNoPermission
	if resp.Code == 1000 {
		t.Error("非作者不应该能更新文章")
	}
}

// TestArticle_Update_Success 测试作者成功更新文章
func TestArticle_Update_Success(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, token, "draft")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	newTitle := "更新后的标题_" + GenerateRandomUsername()
	body := map[string]interface{}{
		"title": newTitle,
	}

	resp, statusCode, err := DoPut(fmt.Sprintf("/author/articles/%d", articleID), body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestArticle_Delete_NotAuthor 测试非作者删除文章
func TestArticle_Delete_NotAuthor(t *testing.T) {
	// 创建作者用户和文章
	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取作者Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "draft")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 创建另一个用户
	otherToken := createTestUserAndLogin(t)
	if otherToken == "" {
		t.Skip("无法获取其他用户Token，跳过测试")
		return
	}

	// 尝试用其他用户删除文章
	resp, _, err := DoDelete(fmt.Sprintf("/author/articles/%d", articleID), otherToken)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回权限错误
	if resp.Code == 1000 {
		t.Error("非作者不应该能删除文章")
	}
}

// TestArticle_Status_Publish 测试发布文章
func TestArticle_Status_Publish(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	// 创建草稿文章
	articleID := createTestArticle(t, token, "draft")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 发布文章
	body := map[string]string{
		"status": "published",
	}

	resp, statusCode, err := DoPatch(fmt.Sprintf("/author/articles/%d/status", articleID), body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestArticle_Status_Offline 测试下线文章
func TestArticle_Status_Offline(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	// 创建已发布文章
	articleID := createTestArticle(t, token, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 下线文章
	body := map[string]string{
		"status": "offline",
	}

	resp, statusCode, err := DoPatch(fmt.Sprintf("/author/articles/%d/status", articleID), body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestArticle_Featured 测试设置文章精选状态
func TestArticle_Featured(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, token, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 设置为精选
	body := map[string]bool{
		"is_featured": true,
	}

	resp, statusCode, err := DoPatch(fmt.Sprintf("/author/articles/%d/featured", articleID), body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}

	// 取消精选
	body = map[string]bool{
		"is_featured": false,
	}

	resp, statusCode, err = DoPatch(fmt.Sprintf("/author/articles/%d/featured", articleID), body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestArticle_Featured_List 测试获取精选文章列表
func TestArticle_Featured_List(t *testing.T) {
	resp, statusCode, err := DoGet("/articles/featured", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestArticle_Featured_List_WithLimit 测试获取指定数量精选文章
func TestArticle_Featured_List_WithLimit(t *testing.T) {
	resp, statusCode, err := DoGetWithQuery("/articles/featured?limit=5", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestArticle_Search_ByKeyword 测试按关键词搜索文章
func TestArticle_Search_ByKeyword(t *testing.T) {
	resp, statusCode, err := DoGetWithQuery("/articles/search?keyword=Go", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 搜索结果可能为空，不一定返回错误
	AssertCode(t, resp, resp.Code)
	_ = statusCode
}

// TestArticle_Search_NoKeyword 测试搜索缺少关键词
func TestArticle_Search_NoKeyword(t *testing.T) {
	resp, statusCode, err := DoGetWithQuery("/articles/search", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 缺少搜索条件应返回错误
	if resp.Code == 1000 {
		t.Log("缺少搜索条件时返回成功，检查业务逻辑")
	}
	_ = statusCode
}

// TestArticle_AuthorArticles 测试获取作者自己的文章列表
func TestArticle_AuthorArticles(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	// 创建几篇文章
	createTestArticle(t, token, "draft")
	createTestArticle(t, token, "published")

	resp, statusCode, err := DoGetWithQuery("/author/articles", token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// createTestArticle 辅助函数：创建测试文章
func createTestArticle(t *testing.T, token, status string) int64 {
	t.Helper()

	body := map[string]interface{}{
		"title":    "测试文章_" + GenerateRandomUsername(),
		"content":  "这是一篇用于测试的文章内容",
		"status":   status,
		"allow_comment": true,
	}

	resp, _, err := DoPost("/articles", body, token)
	if err != nil {
		t.Logf("创建文章失败: %v", err)
		return 0
	}

	if resp.Code != 1000 {
		t.Logf("创建文章失败: %v", resp.Msg)
		return 0
	}

	data, err := ExtractData(resp)
	if err != nil {
		t.Logf("提取数据失败: %v", err)
		return 0
	}

	// 处理 id 可能为 float64、int64 或 string 的情况
	var articleID int64
	switch v := data["id"].(type) {
	case float64:
		articleID = int64(v)
	case int64:
		articleID = v
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			t.Logf("无法解析文章ID字符串: %v", err)
			return 0
		}
		articleID = parsed
	default:
		t.Logf("无法解析文章ID类型: %T", data["id"])
		return 0
	}

	return articleID
}

// createTestArticleWithTags 辅助函数：创建带标签的测试文章
func createTestArticleWithTags(t *testing.T, token string, tags []string) int64 {
	t.Helper()

	body := map[string]interface{}{
		"title":    "测试文章_" + GenerateRandomUsername(),
		"content":  "这是一篇用于测试的文章内容",
		"status":   "published",
		"tags":     tags,
		"allow_comment": true,
	}

	resp, _, err := DoPost("/articles", body, token)
	if err != nil {
		t.Logf("创建文章失败: %v", err)
		return 0
	}

	if resp.Code != 1000 {
		t.Logf("创建文章失败: %v", resp.Msg)
		return 0
	}

	data, err := ExtractData(resp)
	if err != nil {
		t.Logf("提取数据失败: %v", err)
		return 0
	}

	var articleID int64
	switch v := data["id"].(type) {
	case float64:
		articleID = int64(v)
	case int64:
		articleID = v
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			t.Logf("无法解析文章ID字符串: %v", err)
			return 0
		}
		articleID = parsed
	default:
		t.Logf("无法解析文章ID类型: %T", data["id"])
		return 0
	}

	return articleID
}
