package openapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// SignPolicy 对齐 CRMEB OpenApiAuthMiddleware：
// policy = {access_key, conditions: ak/mer/openapi, expiration: YmdHis, unique}
// signature = hex(hmac_sha256(base64(json(policy)), secret_key))
func SignPolicy(accessKey, secretKey, unique string, expirationUnix int64) (string, error) {
	if accessKey == "" || secretKey == "" || unique == "" || expirationUnix <= 0 {
		return "", ErrBadParam
	}
	loc := time.FixedZone("CST", 8*3600)
	policy := map[string]string{
		"access_key": accessKey,
		"conditions": fmt.Sprintf("%s/mer/openapi", accessKey),
		"expiration": time.Unix(expirationUnix, 0).In(loc).Format("20060102150405"),
		"unique":     unique,
	}
	raw, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, _ = mac.Write([]byte(base64.StdEncoding.EncodeToString(raw)))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func VerifySignature(accessKey, secretKey, unique, signature string, expirationUnix int64, now time.Time) error {
	if accessKey == "" || unique == "" || signature == "" || expirationUnix <= 0 {
		return ErrBadParam
	}
	// CRMEB: (time() - expiration) > 300 → 过期（允许未来时间戳）
	if now.Unix()-expirationUnix > 300 {
		return ErrUnauthorized
	}
	expect, err := SignPolicy(accessKey, secretKey, unique, expirationUnix)
	if err != nil {
		return err
	}
	if !hmac.Equal([]byte(expect), []byte(signature)) {
		return ErrUnauthorized
	}
	return nil
}
