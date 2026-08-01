package upload

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
)

// ConfigResolver 由云配置中心实现；返回值只能留在服务端运行时使用。
type ConfigResolver interface {
	Values(ctx context.Context, group string) (map[string]string, error)
}

// DatabaseCOS 每次上传只读取后台数据库配置。后台保存后无需重启即可在下一次上传生效。
type DatabaseCOS struct {
	Resolver ConfigResolver
}

func (d DatabaseCOS) Save(scope string, fh *multipart.FileHeader) (string, string, error) {
	if d.Resolver == nil {
		return "", "", fmt.Errorf("后台 COS 配置服务不可用")
	}
	cosValues, err := d.Resolver.Values(context.Background(), "cos")
	if err != nil {
		return "", "", fmt.Errorf("读取 COS 配置: %w", err)
	}
	if !enabled(cosValues["enabled"]) {
		return "", "", fmt.Errorf("后台 COS 未启用")
	}
	account, err := d.Resolver.Values(context.Background(), "tencent_account")
	if err != nil {
		return "", "", fmt.Errorf("读取腾讯云账号配置: %w", err)
	}
	region := cosValues["region"]
	if strings.TrimSpace(region) == "" {
		region = account["region"]
	}
	store := COS{
		Bucket: cosValues["bucket"], Region: region, BaseURL: cosValues["base_url"], KeyPrefix: cosValues["key_prefix"],
		SecretID: account["secret_id"], SecretKey: account["secret_key"],
	}
	if !store.Configured() {
		return "", "", fmt.Errorf("数据库 COS 已启用但 bucket/region/腾讯云密钥未配置完整")
	}
	return store.Save(scope, fh)
}

func enabled(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

var _ Store = DatabaseCOS{}
