package api

import (
	"fmt"
	"strconv"
	"testing"
)

// TestTag_List_Success 测试获取所有标签
func TestTag_List_Success(t *testing.T) {
	resp, statusCode, err := DoGet("/tags", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestTag_Articles 测试获取标签下文章
func TestTag_Articles(t *testing.T) {
	// 先创建带标签的文章
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	// 创建带标签的文章
	articleID := createTestArticleWithTags(t, token, []string{"测试标签"})
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 获取标签列表（tags API直接返回数组，不是包在list中）
	tagsResp, _, err := DoGet("/tags", "")
	if err != nil {
		t.Fatalf("获取标签失败: %v", err)
	}

	// resp.Data 是直接的数组
	var list []interface{}
	if tagsResp.Data != nil {
		if arr, ok := tagsResp.Data.([]interface{}); ok {
			list = arr
		}
	}
	if len(list) == 0 {
		t.Skip("没有可用标签，跳过测试")
		return
	}

	// 获取第一个标签的ID
	firstTag, _ := list[0].(map[string]interface{})
	var tagID int64
	switch v := firstTag["id"].(type) {
	case float64:
		tagID = int64(v)
	case int64:
		tagID = v
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			tagID = parsed
		}
	}

	if tagID == 0 {
		t.Skip("无法获取标签ID，跳过测试")
		return
	}

	resp, statusCode, err := DoGetWithQuery(fmt.Sprintf("/tags/%d/articles", tagID), "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestTag_Create_NotLogin 测试未登录创建标签
func TestTag_Create_NotLogin(t *testing.T) {
	body := map[string]string{
		"name": "测试标签",
	}

	resp, statusCode, err := DoPost("/tags", body, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回需要登录错误
	if !AssertError(t, resp, statusCode, 1006) {
		t.Logf("响应: %+v", resp)
	}
}

// TestTag_Create_Success 测试成功创建标签
func TestTag_Create_Success(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	body := map[string]string{
		"name": "测试标签_" + GenerateRandomUsername(),
	}

	resp, statusCode, err := DoPost("/tags", body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
		return
	}

	// 验证返回的标签数据
	data, err := ExtractData(resp)
	if err != nil {
		t.Fatalf("提取数据失败: %v", err)
	}

	if data["id"] == nil {
		t.Error("响应缺少 id 字段")
	}
	if data["name"] == nil {
		t.Error("响应缺少 name 字段")
	}
}

// TestTag_Update_Success 测试成功更新标签
func TestTag_Update_Success(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	// 先创建标签
	createBody := map[string]string{
		"name": "测试标签_" + GenerateRandomUsername(),
	}

	createResp, _, _ := DoPost("/tags", createBody, token)
	if createResp.Code != 1000 {
		t.Skipf("创建标签失败: %v", createResp.Msg)
		return
	}

	createData, _ := ExtractData(createResp)
	var tagID int64
	switch v := createData["id"].(type) {
	case float64:
		tagID = int64(v)
	case int64:
		tagID = v
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			tagID = parsed
		}
	}

	if tagID == 0 {
		t.Skip("无法获取标签ID，跳过测试")
		return
	}

	// 更新标签
	updateBody := map[string]string{
		"name": "更新后的标签_" + GenerateRandomUsername(),
	}

	resp, statusCode, err := DoPut(fmt.Sprintf("/tags/%d", tagID), updateBody, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestTag_Update_NotExist 测试更新不存在的标签
func TestTag_Update_NotExist(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	body := map[string]string{
		"name": "更新后的标签",
	}

	resp, statusCode, err := DoPut("/tags/999999999", body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回参数错误
	if !AssertError(t, resp, statusCode, 1001) {
		t.Logf("响应: %+v", resp)
	}
}

// TestTag_Delete_Success 测试成功删除标签
func TestTag_Delete_Success(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	// 先创建标签
	createBody := map[string]string{
		"name": "测试标签_" + GenerateRandomUsername(),
	}

	createResp, _, _ := DoPost("/tags", createBody, token)
	if createResp.Code != 1000 {
		t.Skipf("创建标签失败: %v", createResp.Msg)
		return
	}

	createData, _ := ExtractData(createResp)
	var tagID int64
	switch v := createData["id"].(type) {
	case float64:
		tagID = int64(v)
	case int64:
		tagID = v
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			tagID = parsed
		}
	}

	if tagID == 0 {
		t.Skip("无法获取标签ID，跳过测试")
		return
	}

	// 删除标签
	resp, statusCode, err := DoDelete(fmt.Sprintf("/tags/%d", tagID), token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestTag_Delete_NotExist 测试删除不存在的标签
func TestTag_Delete_NotExist(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	resp, statusCode, err := DoDelete("/tags/999999999", token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回参数错误
	if !AssertError(t, resp, statusCode, 1001) {
		t.Logf("响应: %+v", resp)
	}
}

// TestTag_AuthorTags_NotLogin 测试未登录获取作者标签
func TestTag_AuthorTags_NotLogin(t *testing.T) {
	resp, statusCode, err := DoGet("/author/tags", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回需要登录错误
	if !AssertError(t, resp, statusCode, 1006) {
		t.Logf("响应: %+v", resp)
	}
}

// TestTag_AuthorTags_Success 测试已登录获取作者标签
func TestTag_AuthorTags_Success(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	// 创建带标签的文章
	createTestArticleWithTags(t, token, []string{"标签1", "标签2"})

	resp, statusCode, err := DoGet("/author/tags", token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}
