package authjwt

import (
	"strings"
)

// BearerToken 解析唯一允许的 JWT 头值。调用方必须从 Authori-zation 读取该值；
// 不接受 body、query、cookie 或 Authorization 的兼容输入。
func BearerToken(raw string) (string, error) {
	if !strings.HasPrefix(raw, "Bearer ") {
		return "", ErrInvalidToken
	}
	token := strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
	if token == "" {
		return "", ErrInvalidToken
	}
	return token, nil
}
