package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"testing"
)

// getArticleIDStr 从响应数据中提取文章ID字符串
func getArticleIDStr(t *testing.T, data map[string]interface{}) string {
	t.Helper()
	switch v := data["id"].(type) {
	case float64:
		return fmt.Sprintf("%.0f", v)
	case int64:
		return strconv.FormatInt(v, 10)
	case string:
		return v
	default:
		t.Logf("无法解析文章ID类型: %T", data["id"])
		return ""
	}
}

// getArticleIDInt64 从响应数据中提取文章ID为int64
func getArticleIDInt64(t *testing.T, data map[string]interface{}) int64 {
	t.Helper()
	switch v := data["id"].(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			t.Logf("无法解析文章ID字符串: %v", err)
			return 0
		}
		return parsed
	default:
		t.Logf("无法解析文章ID类型: %T", data["id"])
		return 0
	}
}

// ============================================
// 作者主页测试
// ============================================

// TestAuthor_Info_Success 测试获取作者信息
func TestAuthor_Info_Success(t *testing.T) {
	// 先创建一个用户
	username := GenerateRandomUsername()
	password := "test123456"

	_, _, err := DoPost("/auth/signup", map[string]string{
		"username":    username,
		"password":    password,
		"re_password": password,
	}, "")
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}

	// 获取作者信息
	resp, statusCode, err := DoGet("/authors/"+username, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestAuthor_Info_NotExist 测试获取不存在的作者信息
func TestAuthor_Info_NotExist(t *testing.T) {
	resp, statusCode, err := DoGet("/authors/nonexistent_user_12345", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 应该返回参数错误（用户不存在）
	if !AssertError(t, resp, statusCode, 1001) {
		t.Logf("响应: %+v", resp)
	}
}

// TestAuthor_Articles_Success 测试获取作者文章列表
func TestAuthor_Articles_Success(t *testing.T) {
	// 创建用户并发布文章
	username := GenerateRandomUsername()
	password := "test123456"

	_, _, err := DoPost("/auth/signup", map[string]string{
		"username":    username,
		"password":    password,
		"re_password": password,
	}, "")
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}

	// 登录获取token
	loginResp, _, _ := DoPost("/auth/login", map[string]string{
		"username": username,
		"password": password,
	}, "")

	if loginResp.Code != 1000 {
		t.Skip("登录失败，跳过测试")
		return
	}

	token := extractTokenFromResponse(t, loginResp)

	// 创建并发布文章
	_, _, _ = DoPost("/articles", map[string]interface{}{
		"title":   "作者文章测试",
		"content": "这是作者的文章内容",
		"status":  "published",
	}, token)

	// 获取作者文章列表
	resp, statusCode, err := DoGet("/authors/"+username+"/articles", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestAuthor_Articles_WithPagination 测试作者文章列表分页
func TestAuthor_Articles_WithPagination(t *testing.T) {
	username := GenerateRandomUsername()
	password := "test123456"

	_, _, _ = DoPost("/auth/signup", map[string]string{
		"username":    username,
		"password":    password,
		"re_password": password,
	}, "")

	resp, statusCode, err := DoGetWithQuery("/authors/"+username+"/articles?page=1&size=10&sort=time", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestAuthor_Articles_SortByHot 测试作者文章按热度排序
func TestAuthor_Articles_SortByHot(t *testing.T) {
	username := GenerateRandomUsername()
	password := "test123456"

	_, _, _ = DoPost("/auth/signup", map[string]string{
		"username":    username,
		"password":    password,
		"re_password": password,
	}, "")

	resp, statusCode, err := DoGetWithQuery("/authors/"+username+"/articles?sort=hot", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// ============================================
// RSS订阅测试
// ============================================

// TestRSS_Success 测试获取RSS订阅
func TestRSS_Success(t *testing.T) {
	resp, statusCode, err := DoGet("/rss", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
		return
	}

	// 验证返回数据结构
	data, err := ExtractData(resp)
	if err != nil {
		t.Fatalf("提取数据失败: %v", err)
	}

	if data["title"] == nil {
		t.Error("响应缺少 title 字段")
	}
	if data["items"] == nil {
		t.Error("响应缺少 items 字段")
	}
}

// ============================================
// 文章导出测试
// ============================================

// TestExport_Single_NotLogin 测试未登录导出文章
func TestExport_Single_NotLogin(t *testing.T) {
	resp, statusCode, err := DoGet("/author/articles/1/export", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回需要登录错误
	if !AssertError(t, resp, statusCode, 1006) {
		t.Logf("响应: %+v", resp)
	}
}

// TestExport_Single_Success 测试导出自己的文章
func TestExport_Single_Success(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	// 创建文章
	articleResp, _, _ := DoPost("/articles", map[string]interface{}{
		"title":   "导出测试文章",
		"content": "# 标题\n\n这是导出测试的文章内容，包含Markdown格式。",
	}, token)

	if articleResp.Code != 1000 {
		t.Skipf("创建文章失败: %v", articleResp.Msg)
		return
	}

	data, _ := ExtractData(articleResp)
	articleID := getArticleIDStr(t, data)
	if articleID == "" {
		t.Skip("无法获取文章ID，跳过测试")
		return
	}

	// 导出文章
	resp, statusCode, err := DoGet("/author/articles/"+articleID+"/export", token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestExport_Single_NotAuthor 测试导出他人文章
func TestExport_Single_NotAuthor(t *testing.T) {
	// 用户1创建文章
	token1 := createTestUserAndLogin(t)
	if token1 == "" {
		t.Skip("无法获取Token1，跳过测试")
		return
	}

	articleResp, _, _ := DoPost("/articles", map[string]interface{}{
		"title":   "他人文章导出测试",
		"content": "这是他人的文章内容",
	}, token1)

	if articleResp.Code != 1000 {
		t.Skip("创建文章失败，跳过测试")
		return
	}

	data, _ := ExtractData(articleResp)
	articleID := getArticleIDStr(t, data)
	if articleID == "" {
		t.Skip("无法获取文章ID，跳过测试")
		return
	}

	// 用户2尝试导出
	token2 := createTestUserAndLogin(t)
	if token2 == "" {
		t.Skip("无法获取Token2，跳过测试")
		return
	}

	resp, statusCode, err := DoGet("/author/articles/"+articleID+"/export", token2)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回无权限错误
	if !AssertError(t, resp, statusCode, 1013) {
		t.Logf("响应: %+v", resp)
	}
}

// TestExport_Single_NotExist 测试导出不存在的文章
func TestExport_Single_NotExist(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	resp, statusCode, err := DoGet("/author/articles/999999999/export", token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回错误（1014文章不存在 或 1005服务繁忙）
	if resp.Code == 1000 {
		t.Errorf("期望返回错误，实际返回成功")
	}
	_ = statusCode
}

// TestExport_Batch_NotLogin 测试未登录批量导出
func TestExport_Batch_NotLogin(t *testing.T) {
	resp, statusCode, err := DoPost("/author/articles/export", map[string]interface{}{
		"article_ids": []int64{1, 2, 3},
	}, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回需要登录错误
	if !AssertError(t, resp, statusCode, 1006) {
		t.Logf("响应: %+v", resp)
	}
}

// TestExport_Batch_Success 测试批量导出自己的文章
func TestExport_Batch_Success(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	// 创建两篇文章
	article1Resp, _, _ := DoPost("/articles", map[string]interface{}{
		"title":   "批量导出测试文章1",
		"content": "文章内容1",
	}, token)

	article2Resp, _, _ := DoPost("/articles", map[string]interface{}{
		"title":   "批量导出测试文章2",
		"content": "文章内容2",
	}, token)

	if article1Resp.Code != 1000 || article2Resp.Code != 1000 {
		t.Skip("创建文章失败，跳过测试")
		return
	}

	data1, _ := ExtractData(article1Resp)
	data2, _ := ExtractData(article2Resp)
	id1 := getArticleIDInt64(t, data1)
	id2 := getArticleIDInt64(t, data2)
	if id1 == 0 || id2 == 0 {
		t.Skip("无法获取文章ID，跳过测试")
		return
	}

	// 批量导出
	resp, statusCode, err := DoPost("/author/articles/export", map[string]interface{}{
		"article_ids": []int64{id1, id2},
	}, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestExport_Batch_EmptyIDs 测试批量导出空ID列表
func TestExport_Batch_EmptyIDs(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	resp, statusCode, err := DoPost("/author/articles/export", map[string]interface{}{
		"article_ids": []int64{},
	}, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回参数错误
	if !AssertError(t, resp, statusCode, 1001) {
		t.Logf("响应: %+v", resp)
	}
}

// ============================================
// 文件上传测试
// ============================================

// TestUpload_Image_NotLogin 测试未登录上传图片
func TestUpload_Image_NotLogin(t *testing.T) {
	// 创建一个假的文件上传请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatalf("创建表单失败: %v", err)
	}
	part.Write([]byte("fake image content"))
	writer.Close()

	url := BaseURL + APIVersion + "/upload/image"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result ResponseData
	json.Unmarshal(respBody, &result)

	// 期望返回需要登录错误
	if result.Code != 1006 {
		t.Errorf("期望业务码1006，实际得到: %d", result.Code)
	}
}

// TestUpload_Image_Success 测试成功上传图片
func TestUpload_Image_Success(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	// 创建一个假的文件上传请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatalf("创建表单失败: %v", err)
	}
	part.Write([]byte("fake image content for upload test"))
	writer.Close()

	url := BaseURL + APIVersion + "/upload/image"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result ResponseData
	json.Unmarshal(respBody, &result)

	if result.Code != 1000 {
		t.Errorf("期望业务码1000，实际得到: %d, msg: %v", result.Code, result.Msg)
	}
}

// TestUpload_Image_NoFile 测试上传图片但没有文件
func TestUpload_Image_NoFile(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	// 发送空的请求
	url := BaseURL + APIVersion + "/upload/image"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result ResponseData
	json.Unmarshal(respBody, &result)

	// 期望返回参数错误
	if result.Code != 1001 {
		t.Errorf("期望业务码1001，实际得到: %d", result.Code)
	}
}

// TestUpload_Attachment_NotLogin 测试未登录上传附件
func TestUpload_Attachment_NotLogin(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.pdf")
	if err != nil {
		t.Fatalf("创建表单失败: %v", err)
	}
	part.Write([]byte("fake pdf content"))
	writer.Close()

	url := BaseURL + APIVersion + "/upload/attachment"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result ResponseData
	json.Unmarshal(respBody, &result)

	// 期望返回需要登录错误
	if result.Code != 1006 {
		t.Errorf("期望业务码1006，实际得到: %d", result.Code)
	}
}

// TestUpload_Attachment_Success 测试成功上传附件
func TestUpload_Attachment_Success(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "document.pdf")
	if err != nil {
		t.Fatalf("创建表单失败: %v", err)
	}
	part.Write([]byte("fake pdf content for upload test"))
	writer.Close()

	url := BaseURL + APIVersion + "/upload/attachment"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result ResponseData
	json.Unmarshal(respBody, &result)

	if result.Code != 1000 {
		t.Errorf("期望业务码1000，实际得到: %d, msg: %v", result.Code, result.Msg)
	}
}
