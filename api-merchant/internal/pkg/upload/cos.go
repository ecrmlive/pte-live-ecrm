package upload

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// COS 腾讯云对象存储上传。
type COS struct {
	Bucket    string
	Region    string
	SecretID  string
	SecretKey string
	BaseURL   string // 对外访问前缀，如 https://cos.example.com/qixi-mergers
	KeyPrefix string // 对象键前缀，默认从 BaseURL path 推导
}

// Configured 是否已具备上传所需密钥与桶信息。
func (c COS) Configured() bool {
	return c.Bucket != "" && c.Region != "" && c.SecretID != "" && c.SecretKey != ""
}

func (c COS) keyPrefix() string {
	p := strings.Trim(strings.ReplaceAll(c.KeyPrefix, "\\", "/"), "/")
	if p != "" {
		return p
	}
	if u, err := url.Parse(strings.TrimSpace(c.BaseURL)); err == nil {
		p = strings.Trim(u.Path, "/")
		if p != "" {
			return p
		}
	}
	return "qixi-mergers"
}

func (c COS) client() (*cos.Client, error) {
	bucketURL, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", c.Bucket, c.Region))
	if err != nil {
		return nil, err
	}
	return cos.NewClient(&cos.BaseURL{BucketURL: bucketURL}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  c.SecretID,
			SecretKey: c.SecretKey,
		},
	}), nil
}

func (c COS) Save(scope string, fh *multipart.FileHeader) (publicURL, name string, err error) {
	if !c.Configured() {
		return "", "", fmt.Errorf("未配置腾讯云 COS（bucket/region/secret）")
	}
	if fh == nil {
		return "", "", fmt.Errorf("empty file")
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".mp4", ".mov", ".webm":
	default:
		return "", "", fmt.Errorf("unsupported type %s", ext)
	}
	maxSize := int64(5 << 20)
	if ext == ".mp4" || ext == ".mov" || ext == ".webm" {
		maxSize = 100 << 20
	}
	if fh.Size > maxSize {
		return "", "", fmt.Errorf("file too large")
	}
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", "", err
	}
	fname := hex.EncodeToString(b[:]) + ext
	day := time.Now().Format("20060102")
	scope = strings.Trim(strings.ReplaceAll(scope, "\\", "/"), "/")
	key := strings.Join([]string{c.keyPrefix(), scope, day, fname}, "/")

	cli, err := c.client()
	if err != nil {
		return "", "", err
	}
	src, err := fh.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	if _, err := cli.Object.Put(ctx, key, src, nil); err != nil {
		return "", "", err
	}

	base := strings.TrimRight(strings.TrimSpace(c.BaseURL), "/")
	if base == "" {
		base = fmt.Sprintf("https://%s.cos.%s.myqcloud.com", c.Bucket, c.Region)
	}
	publicURL = base + "/" + key
	name = fh.Filename
	if name == "" {
		name = fname
	}
	return publicURL, name, nil
}

var _ Store = COS{}
