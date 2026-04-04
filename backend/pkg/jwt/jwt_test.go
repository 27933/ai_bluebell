package jwt

import (
	"testing"
	"time"
	"github.com/spf13/viper"
)

func TestGenToken(t *testing.T) {
	// 设置测试配置
	viper.Set("auth.jwt_expire", 8760)

	// 生成token
	token, err := GenToken(123, "testuser")
	if err != nil {
		t.Fatalf("GenToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("Token should not be empty")
	}

	// 解析token
	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	// 验证claims
	if claims.UserID != 123 {
		t.Errorf("Expected UserID 123, got %d", claims.UserID)
	}

	if claims.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got %s", claims.Username)
	}

	// 检查过期时间
	expectedExpire := time.Now().Add(time.Duration(8760) * time.Hour).Unix()
	if claims.ExpiresAt < expectedExpire-10 || claims.ExpiresAt > expectedExpire+10 {
		t.Errorf("ExpiresAt not within expected range")
	}
}