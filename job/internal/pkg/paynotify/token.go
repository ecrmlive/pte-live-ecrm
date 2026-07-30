package paynotify

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// MakeToken 签发沙箱回调令牌（仅下发给已登录支付发起方）。
func MakeToken(secret, channel, outTradeNo string, uid, groupOrderID uint, amount float64) string {
	raw := payload(channel, outTradeNo, uid, groupOrderID, amount)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(raw))
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifyToken 校验回调令牌。
func VerifyToken(secret, channel, outTradeNo string, uid, groupOrderID uint, amount float64, token string) bool {
	want := MakeToken(secret, channel, outTradeNo, uid, groupOrderID, amount)
	return hmac.Equal([]byte(strings.ToLower(want)), []byte(strings.ToLower(strings.TrimSpace(token))))
}

func payload(channel, outTradeNo string, uid, groupOrderID uint, amount float64) string {
	return fmt.Sprintf("%s|%s|%d|%d|%s",
		strings.TrimSpace(channel),
		strings.TrimSpace(outTradeNo),
		uid, groupOrderID,
		strconv.FormatFloat(amount, 'f', 2, 64),
	)
}
