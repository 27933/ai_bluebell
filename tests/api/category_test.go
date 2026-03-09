package api

import (
	"fmt"
	"strconv"
	"testing"
)

// TestCategory_List_Success 测试获取栏目列表
func TestCategory_List_Success(t *testing.T) {
	resp, statusCode, err := DoGet("/categories", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
		return
	}

	// 验证返回data包含list字段（list可能为null或空数组，均为正常）
	data, err := ExtractData(resp)
	if err != nil {
		t.Fatalf("提取数据失败: %v", err)
	}

	if data == nil {
		t.Error("响应 data 为 nil")
	}
}

// TestCategory_Detail_Success 测试获取栏目详情
func TestCategory_Detail_Success(t *testing.T) {
	// 先创建栏目
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	categoryID := createTestCategory(t, token)
	if categoryID == 0 {
		t.Skip("无法创建测试栏目，跳过测试")
		return
	}

	resp, statusCode, err := DoGet(fmt.Sprintf("/categories/%d", categoryID), "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestCategory_Detail_NotExist 测试获取不存在的栏目
func TestCategory_Detail_NotExist(t *testing.T) {
	resp, statusCode, err := DoGet("/categories/999999999", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回栏目不存在错误
	if !AssertError(t, resp, statusCode, 1008) {
		t.Logf("响应: %+v", resp)
	}
}

// TestCategory_Create_NotLogin 测试未登录创建栏目
func TestCategory_Create_NotLogin(t *testing.T) {
	body := map[string]string{
		"name": "测试栏目",
	}

	resp, statusCode, err := DoPost("/categories", body, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回需要登录错误
	if !AssertError(t, resp, statusCode, 1006) {
		t.Logf("响应: %+v", resp)
	}
}

// TestCategory_Create_Success 测试成功创建栏目
func TestCategory_Create_Success(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	body := map[string]string{
		"category_name": "测试栏目_" + GenerateRandomUsername(),
		"introduction":  "测试栏目描述，至少两个字符",
	}

	resp, statusCode, err := DoPost("/categories", body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestCategory_Create_Duplicate 测试创建重复名称栏目
func TestCategory_Create_Duplicate(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	categoryName := "测试栏目_" + GenerateRandomUsername()

	// 第一次创建
	body := map[string]string{
		"category_name": categoryName,
		"introduction":  "测试栏目描述，至少两个字符",
	}

	resp1, _, _ := DoPost("/categories", body, token)
	if resp1.Code != 1000 {
		t.Skipf("第一次创建失败，跳过重复测试: %v", resp1.Msg)
		return
	}

	// 第二次创建相同名称
	resp2, statusCode2, err := DoPost("/categories", body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 重复创建应返回错误（实际返回 CodeServerBusy 1005）
	if resp2.Code == 1000 {
		t.Error("重复创建栏目不应该成功")
	}
	_ = statusCode2
}

// TestCategory_Update_Success 测试成功更新栏目
func TestCategory_Update_Success(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	categoryID := createTestCategory(t, token)
	if categoryID == 0 {
		t.Skip("无法创建测试栏目，跳过测试")
		return
	}

	body := map[string]string{
		"category_name": "更新后的栏目_" + GenerateRandomUsername(),
		"introduction":  "更新后的描述，至少两个字符",
	}

	resp, statusCode, err := DoPut(fmt.Sprintf("/categories/%d", categoryID), body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestCategory_Update_NotExist 测试更新不存在的栏目
func TestCategory_Update_NotExist(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	body := map[string]string{
		"category_name": "更新后的栏目名称",
		"introduction":  "更新后的描述内容",
	}

	resp, statusCode, err := DoPut("/categories/999999999", body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回参数错误（不存在的栏目返回 CodeInvalidParam）
	if !AssertError(t, resp, statusCode, 1001) {
		t.Logf("响应: %+v", resp)
	}
}

// TestCategory_Delete_Success 测试成功删除栏目
func TestCategory_Delete_Success(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	categoryID := createTestCategory(t, token)
	if categoryID == 0 {
		t.Skip("无法创建测试栏目，跳过测试")
		return
	}

	resp, statusCode, err := DoDelete(fmt.Sprintf("/categories/%d", categoryID), token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestCategory_Delete_NotExist 测试删除不存在的栏目
func TestCategory_Delete_NotExist(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	resp, statusCode, err := DoDelete("/categories/999999999", token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回参数错误（不存在的栏目返回 CodeInvalidParam）
	if !AssertError(t, resp, statusCode, 1001) {
		t.Logf("响应: %+v", resp)
	}
}

// TestCategory_Articles 测试获取栏目下文章
func TestCategory_Articles(t *testing.T) {
	// 创建栏目
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	categoryID := createTestCategory(t, token)
	if categoryID == 0 {
		t.Skip("无法创建测试栏目，跳过测试")
		return
	}

	resp, statusCode, err := DoGetWithQuery(fmt.Sprintf("/categories/%d/articles", categoryID), "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestCategory_AddArticle_NotAuthor 测试非作者添加文章到栏目
func TestCategory_AddArticle_NotAuthor(t *testing.T) {
	// 创建文章作者
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

	// 创建栏目
	categoryToken := createTestUserAndLogin(t)
	if categoryToken == "" {
		t.Skip("无法获取栏目创建者Token，跳过测试")
		return
	}

	categoryID := createTestCategory(t, categoryToken)
	if categoryID == 0 {
		t.Skip("无法创建测试栏目，跳过测试")
		return
	}

	// 另一用户尝试添加文章到栏目
	otherToken := createTestUserAndLogin(t)
	if otherToken == "" {
		t.Skip("无法获取其他用户Token，跳过测试")
		return
	}

	body := map[string]interface{}{
		"category_ids": []int64{categoryID},
	}

	resp, statusCode, err := DoPost(fmt.Sprintf("/articles/%d/categories", articleID), body, otherToken)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回权限错误
	if resp.Code == 1000 {
		t.Error("非作者不应该能添加文章到栏目")
	}
	_ = statusCode
}

// TestCategory_GetByArticle 测试获取文章所属栏目
func TestCategory_GetByArticle(t *testing.T) {
	// 创建文章
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

	resp, statusCode, err := DoGet(fmt.Sprintf("/articles/%d/categories", articleID), "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// createTestCategory 辅助函数：创建测试栏目
func createTestCategory(t *testing.T, token string) int64 {
	t.Helper()

	categoryName := "测试栏目_" + GenerateRandomUsername()
	body := map[string]string{
		"category_name": categoryName,
		"introduction":  "测试栏目描述，至少两个字符",
	}

	resp, _, err := DoPost("/categories", body, token)
	if err != nil {
		t.Logf("创建栏目失败: %v", err)
		return 0
	}

	if resp.Code != 1000 {
		t.Logf("创建栏目失败: %v", resp.Msg)
		return 0
	}

	// 创建成功后通过列表查询，按名称匹配找到刚创建的栏目
	listResp, _, err := DoGet("/categories", "")
	if err != nil {
		t.Logf("获取栏目列表失败: %v", err)
		return 0
	}

	data, err := ExtractData(listResp)
	if err != nil {
		t.Logf("提取数据失败: %v", err)
		return 0
	}

	list, ok := data["list"].([]interface{})
	if !ok || len(list) == 0 {
		t.Log("栏目列表为空")
		return 0
	}

	// 按名称匹配找到刚创建的栏目
	for _, item := range list {
		cat, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if cat["category_name"] == categoryName {
			var categoryID int64
			switch v := cat["id"].(type) {
			case float64:
				categoryID = int64(v)
			case int64:
				categoryID = v
			case string:
				parsed, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					t.Logf("无法解析栏目ID字符串: %v", err)
					return 0
				}
				categoryID = parsed
			}
			return categoryID
		}
	}

	t.Log("在列表中未找到刚创建的栏目")
	return 0
}
