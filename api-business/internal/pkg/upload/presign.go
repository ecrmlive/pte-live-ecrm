// Package upload issues narrowly scoped COS PUT URLs. The API never receives
// file bytes and the browser persists only the verified object key.
package upload

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type ConfigResolver interface {
	Values(ctx context.Context, group string) (map[string]string, error)
}

type PresignInput struct {
	Scope       string
	Filename    string
	ContentType string
	Size        int64
}

type PresignResult struct {
	ObjectKey string            `json:"object_key"`
	UploadURL string            `json:"upload_url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	ExpiresAt time.Time         `json:"expires_at"`
}

// DatabaseCOS reads COS credentials only from the platform-owned config table.
type DatabaseCOS struct{ Resolver ConfigResolver }

func (d DatabaseCOS) PresignPut(ctx context.Context, in PresignInput) (*PresignResult, error) {
	if d.Resolver == nil {
		return nil, fmt.Errorf("后台 COS 配置服务不可用")
	}
	if err := validate(in); err != nil {
		return nil, err
	}
	cosValues, err := d.Resolver.Values(ctx, "cos")
	if err != nil {
		return nil, fmt.Errorf("读取 COS 配置: %w", err)
	}
	if !enabled(cosValues["enabled"]) {
		return nil, fmt.Errorf("后台 COS 未启用")
	}
	account, err := d.Resolver.Values(ctx, "tencent_account")
	if err != nil {
		return nil, fmt.Errorf("读取腾讯云账号配置: %w", err)
	}
	region := strings.TrimSpace(cosValues["region"])
	if region == "" {
		region = strings.TrimSpace(account["region"])
	}
	c := COS{
		Bucket: cosValues["bucket"], Region: region, BaseURL: cosValues["base_url"], KeyPrefix: cosValues["key_prefix"],
		SecretID: account["secret_id"], SecretKey: account["secret_key"],
	}
	return c.PresignPut(ctx, in)
}

func (d DatabaseCOS) Exists(ctx context.Context, key string) error {
	if d.Resolver == nil {
		return fmt.Errorf("后台 COS 配置服务不可用")
	}
	cosValues, err := d.Resolver.Values(ctx, "cos")
	if err != nil {
		return fmt.Errorf("读取 COS 配置: %w", err)
	}
	account, err := d.Resolver.Values(ctx, "tencent_account")
	if err != nil {
		return fmt.Errorf("读取腾讯云账号配置: %w", err)
	}
	region := strings.TrimSpace(cosValues["region"])
	if region == "" {
		region = strings.TrimSpace(account["region"])
	}
	c := COS{Bucket: cosValues["bucket"], Region: region, SecretID: account["secret_id"], SecretKey: account["secret_key"]}
	if !c.Configured() {
		return fmt.Errorf("数据库 COS 配置不完整")
	}
	return c.Exists(ctx, key)
}

type COS struct {
	Bucket    string
	Region    string
	SecretID  string
	SecretKey string
	BaseURL   string
	KeyPrefix string
}

func (c COS) Configured() bool {
	return strings.TrimSpace(c.Bucket) != "" && strings.TrimSpace(c.Region) != "" && strings.TrimSpace(c.SecretID) != "" && strings.TrimSpace(c.SecretKey) != ""
}

func (c COS) PresignPut(ctx context.Context, in PresignInput) (*PresignResult, error) {
	if !c.Configured() {
		return nil, fmt.Errorf("数据库 COS 已启用但 bucket/region/腾讯云密钥未配置完整")
	}
	if err := validate(in); err != nil {
		return nil, err
	}
	key, err := c.newKey(in.Scope, in.Filename)
	if err != nil {
		return nil, err
	}
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	headers := http.Header{}
	headers.Set("Content-Type", normalizedContentType(in.ContentType, in.Filename))
	expiresAt := time.Now().Add(10 * time.Minute)
	presigned, err := client.Object.GetPresignedURL(ctx, http.MethodPut, key, c.SecretID, c.SecretKey, time.Until(expiresAt), &cos.PresignedURLOptions{Header: &headers}, false)
	if err != nil {
		return nil, fmt.Errorf("签发 COS PUT 地址: %w", err)
	}
	return &PresignResult{ObjectKey: key, UploadURL: presigned.String(), Method: http.MethodPut, Headers: map[string]string{"Content-Type": headers.Get("Content-Type")}, ExpiresAt: expiresAt}, nil
}

func (c COS) Exists(ctx context.Context, key string) error {
	if !c.Configured() {
		return fmt.Errorf("数据库 COS 配置不完整")
	}
	if strings.TrimSpace(key) == "" || strings.Contains(key, "..") {
		return fmt.Errorf("对象键不合法")
	}
	client, err := c.client()
	if err != nil {
		return err
	}
	_, err = client.Object.Head(ctx, key, nil)
	if err != nil {
		return fmt.Errorf("COS 对象不存在或不可访问: %w", err)
	}
	return nil
}

func (c COS) client() (*cos.Client, error) {
	bucketURL, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", c.Bucket, c.Region))
	if err != nil {
		return nil, err
	}
	return cos.NewClient(&cos.BaseURL{BucketURL: bucketURL}, &http.Client{Transport: &cos.AuthorizationTransport{SecretID: c.SecretID, SecretKey: c.SecretKey}}), nil
}

func (c COS) newKey(scope, filename string) (string, error) {
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		return "", err
	}
	ext := strings.ToLower(filepath.Ext(filename))
	prefix := strings.Trim(strings.ReplaceAll(c.KeyPrefix, "\\", "/"), "/")
	if prefix == "" {
		prefix = "pte-live-ecrm"
	}
	scope = strings.Trim(strings.ReplaceAll(scope, "\\", "/"), "/")
	return strings.Join([]string{prefix, scope, time.Now().Format("20060102"), hex.EncodeToString(random[:]) + ext}, "/"), nil
}

func validate(in PresignInput) error {
	if strings.TrimSpace(in.Scope) == "" || strings.Contains(in.Scope, "..") {
		return fmt.Errorf("上传范围不合法")
	}
	if strings.TrimSpace(in.Filename) == "" || in.Size <= 0 || in.Size > 10<<20 {
		return fmt.Errorf("文件大小必须在 1B 到 10MB 之间")
	}
	switch strings.ToLower(filepath.Ext(in.Filename)) {
	case ".jpg", ".jpeg", ".png", ".webp":
		return nil
	default:
		return fmt.Errorf("仅支持 JPG、PNG、WebP 图片")
	}
}

func normalizedContentType(contentType, filename string) string {
	contentType = strings.TrimSpace(strings.ToLower(strings.Split(contentType, ";")[0]))
	if strings.HasPrefix(contentType, "image/") {
		return contentType
	}
	if guessed := mime.TypeByExtension(filepath.Ext(filename)); guessed != "" {
		return guessed
	}
	return "application/octet-stream"
}

func enabled(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
