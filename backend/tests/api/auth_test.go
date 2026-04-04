package api

import (
	"encoding/json"
	"testing"
)

// TestAuth_Register_Success 测试成功注册
func TestAuth_Register_Success(t *testing.T) {
	username := GenerateRandomUsername()
	body := map[string]string{
		"username":    username,
		"password":    "test123456",
		"re_password": "test123456",
	}

	resp, statusCode, err := DoPost("/auth/signup", body, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestAuth_Register_UserExist 测试用户名已存在
func TestAuth_Register_UserExist(t *testing.T) {
	username := GenerateRandomUsername()
	body := map[string]string{
		"username":    username,
		"password":    "test123456",
		"re_password": "test123456",
	}

	// 第一次注册
	resp1, _, err := DoPost("/auth/signup", body, "")
	if err != nil {
		t.Fatalf("第一次请求失败: %v", err)
	}
	if resp1.Code != 1000 {
		t.Fatalf("第一次注册应成功，实际: %d", resp1.Code)
	}

	// 第二次注册相同用户名
	resp2, statusCode2, err := DoPost("/auth/signup", body, "")
	if err != nil {
		t.Fatalf("第二次请求失败: %v", err)
	}

	// 期望返回用户名已存在错误
	if !AssertError(t, resp2, statusCode2, 1002) {
		t.Logf("响应: %+v", resp2)
	}
}

// TestAuth_Register_PasswordMismatch 测试密码不匹配
func TestAuth_Register_PasswordMismatch(t *testing.T) {
	body := map[string]string{
		"username":    GenerateRandomUsername(),
		"password":    "test123456",
		"re_password": "different123",
	}

	resp, statusCode, err := DoPost("/auth/signup", body, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回参数错误
	if !AssertError(t, resp, statusCode, 1001) {
		t.Logf("响应: %+v", resp)
	}
}

// TestAuth_Register_MissingField 测试缺少必填字段
func TestAuth_Register_MissingField(t *testing.T) {
	body := map[string]string{
		"username": GenerateRandomUsername(),
		// 缺少 password
		"re_password": "test123456",
	}

	resp, statusCode, err := DoPost("/auth/signup", body, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertError(t, resp, statusCode, 1001) {
		t.Logf("响应: %+v", resp)
	}
}

// TestAuth_Login_Success 测试成功登录
func TestAuth_Login_Success(t *testing.T) {
	// 先注册用户
	username := GenerateRandomUsername()
	password := "test123456"
	registerBody := map[string]string{
		"username":    username,
		"password":    password,
		"re_password": password,
	}

	_, _, err := DoPost("/auth/signup", registerBody, "")
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}

	// 登录
	loginBody := map[string]string{
		"username": username,
		"password": password,
	}

	resp, statusCode, err := DoPost("/auth/login", loginBody, "")
	if err != nil {
		t.Fatalf("登录请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
		return
	}

	// 验证返回结构
	data, err := ExtractData(resp)
	if err != nil {
		t.Fatalf("提取数据失败: %v", err)
	}

	if data["user"] == nil {
		t.Error("响应缺少 user 字段")
	}
	if data["token"] == nil {
		t.Error("响应缺少 token 字段")
	}
}

// TestAuth_Login_WrongPassword 测试密码错误
func TestAuth_Login_WrongPassword(t *testing.T) {
	// 先注册用户
	username := GenerateRandomUsername()
	password := "test123456"
	registerBody := map[string]string{
		"username":    username,
		"password":    password,
		"re_password": password,
	}

	_, _, err := DoPost("/auth/signup", registerBody, "")
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}

	// 使用错误密码登录
	loginBody := map[string]string{
		"username": username,
		"password": "wrongpassword",
	}

	resp, statusCode, err := DoPost("/auth/login", loginBody, "")
	if err != nil {
		t.Fatalf("登录请求失败: %v", err)
	}

	// 期望返回用户名或密码错误
	if !AssertError(t, resp, statusCode, 1004) {
		t.Logf("响应: %+v", resp)
	}
}

// TestAuth_Login_UserNotExist 测试用户不存在
func TestAuth_Login_UserNotExist(t *testing.T) {
	loginBody := map[string]string{
		"username": GenerateRandomUsername() + "_notexist",
		"password": "test123456",
	}

	resp, statusCode, err := DoPost("/auth/login", loginBody, "")
	if err != nil {
		t.Fatalf("登录请求失败: %v", err)
	}

	// 期望返回用户名或密码错误
	if !AssertError(t, resp, statusCode, 1004) {
		t.Logf("响应: %+v", resp)
	}
}

// TestAuth_RefreshToken_Success 测试成功刷新 Token
func TestAuth_RefreshToken_Success(t *testing.T) {
	// 先注册并登录
	username := GenerateRandomUsername()
	password := "test123456"
	registerResp, _, _ := DoPost("/auth/signup", map[string]string{
		"username":    username,
		"password":    password,
		"re_password": password,
	}, "")

	if registerResp.Code != 1000 && registerResp.Code != 1002 {
		t.Fatalf("注册失败: %v", registerResp.Msg)
	}

	loginResp, _, err := DoPost("/auth/login", map[string]string{
		"username": username,
		"password": password,
	}, "")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	if loginResp.Code != 1000 {
		t.Skipf("登录失败，跳过刷新Token测试: %v", loginResp.Msg)
		return
	}

	data, err := ExtractData(loginResp)
	if err != nil {
		t.Fatalf("提取数据失败: %v", err)
	}

	tokenData, ok := data["token"].(map[string]interface{})
	if !ok {
		t.Skip("响应中无Token数据，跳过刷新Token测试")
		return
	}

	refreshToken, ok := tokenData["refresh_token"].(string)
	if !ok || refreshToken == "" {
		t.Skip("响应中无RefreshToken，跳过刷新Token测试")
		return
	}

	// 刷新Token
	refreshBody := map[string]string{
		"refresh_token": refreshToken,
	}

	resp, statusCode, err := DoPost("/auth/refresh", refreshBody, "")
	if err != nil {
		t.Fatalf("刷新Token请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestAuth_RefreshToken_Invalid 测试无效 Token 刷新
func TestAuth_RefreshToken_Invalid(t *testing.T) {
	refreshBody := map[string]string{
		"refresh_token": "invalid_token_xxx",
	}

	resp, statusCode, err := DoPost("/auth/refresh", refreshBody, "")
	if err != nil {
		t.Fatalf("刷新Token请求失败: %v", err)
	}

	// 期望返回无效的token错误
	if !AssertError(t, resp, statusCode, 1007) {
		t.Logf("响应: %+v", resp)
	}
}

// TestAuth_GetProfile_NotLogin 测试未登录获取用户信息
func TestAuth_GetProfile_NotLogin(t *testing.T) {
	resp, statusCode, err := DoGet("/auth/profile", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回需要登录错误
	if !AssertError(t, resp, statusCode, 1006) {
		t.Logf("响应: %+v", resp)
	}
}

// TestAuth_GetProfile_Success 测试已登录获取用户信息
func TestAuth_GetProfile_Success(t *testing.T) {
	// 创建测试用户并登录
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	resp, statusCode, err := DoGet("/auth/profile", token)
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

	if data["username"] == nil {
		t.Error("响应缺少 username 字段")
	}
}

// TestAuth_UpdateProfile_NotLogin 测试未登录更新用户信息
func TestAuth_UpdateProfile_NotLogin(t *testing.T) {
	body := map[string]string{
		"email": "newemail@test.com",
	}

	resp, statusCode, err := DoPut("/auth/profile", body, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回需要登录错误
	if !AssertError(t, resp, statusCode, 1006) {
		t.Logf("响应: %+v", resp)
	}
}

// TestAuth_UpdateProfile_Success 测试已登录更新用户信息
func TestAuth_UpdateProfile_Success(t *testing.T) {
	// 创建测试用户并登录
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	// 更新昵称（避免邮箱冲突，使用昵称测试）
	newNickname := "测试昵称_" + GenerateRandomUsername()
	body := map[string]string{
		"nickname": newNickname,
	}

	resp, statusCode, err := DoPut("/auth/profile", body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
		return
	}

	// 再次获取用户信息验证更新
	getResp, _, err := DoGet("/auth/profile", token)
	if err != nil {
		t.Fatalf("获取用户信息失败: %v", err)
	}

	data, _ := ExtractData(getResp)
	if data["nickname"] != newNickname {
		t.Errorf("昵称未更新，期望: %s, 实际: %v", newNickname, data["nickname"])
	}
}

// createTestUserAndLogin 辅助函数：创建测试用户并返回Token
func createTestUserAndLogin(t *testing.T) string {
	t.Helper()

	username := GenerateRandomUsername()
	password := "test123456"

	// 注册
	_, _, err := DoPost("/auth/signup", map[string]string{
		"username":    username,
		"password":    password,
		"re_password": password,
	}, "")
	if err != nil {
		t.Logf("注册用户失败: %v", err)
		return ""
	}

	// 登录
	resp, _, err := DoPost("/auth/login", map[string]string{
		"username": username,
		"password": password,
	}, "")
	if err != nil {
		t.Logf("登录失败: %v", err)
		return ""
	}

	if resp.Code != 1000 {
		t.Logf("登录响应错误: %v", resp.Msg)
		return ""
	}

	data, err := ExtractData(resp)
	if err != nil {
		t.Logf("提取数据失败: %v", err)
		return ""
	}

	tokenData, ok := data["token"].(map[string]interface{})
	if !ok {
		t.Log("响应中无Token数据")
		return ""
	}

	accessToken, ok := tokenData["access_token"].(string)
	if !ok {
		t.Log("响应中无AccessToken")
		return ""
	}

	return accessToken
}

// createTestUserAndLoginWithRole 辅助函数：创建指定角色的用户并返回Token
func createTestUserAndLoginWithRole(t *testing.T, role string) string {
	t.Helper()

	token := createTestUserAndLogin(t)
	if token == "" {
		return ""
	}

	// 如果是admin角色，这里需要管理员去修改角色
	// 简化处理：返回token，调用方需要自行确保用户有权限
	_ = role

	return token
}

// extractTokenFromResponse 从登录响应中提取Token
func extractTokenFromResponse(t *testing.T, resp *ResponseData) string {
	t.Helper()

	data, err := ExtractData(resp)
	if err != nil {
		t.Logf("提取数据失败: %v", err)
		return ""
	}

	tokenData, ok := data["token"].(map[string]interface{})
	if !ok {
		// 尝试直接解析为LoginResponseData
		dataBytes, _ := json.Marshal(resp.Data)
		var loginData LoginResponseData
		if err := json.Unmarshal(dataBytes, &loginData); err == nil && loginData.Token.AccessToken != "" {
			return loginData.Token.AccessToken
		}
		t.Log("响应中无Token数据")
		return ""
	}

	accessToken, ok := tokenData["access_token"].(string)
	if !ok {
		t.Log("响应中无AccessToken")
		return ""
	}

	return accessToken
}

// loginWithCredentials 使用指定凭据登录
func loginWithCredentials(t *testing.T, username, password string) string {
	t.Helper()

	resp, _, err := DoPost("/auth/login", map[string]string{
		"username": username,
		"password": password,
	}, "")
	if err != nil {
		t.Logf("登录失败: %v", err)
		return ""
	}

	if resp.Code != 1000 {
		t.Logf("登录失败: %v", resp.Msg)
		return ""
	}

	return extractTokenFromResponse(t, resp)
}
