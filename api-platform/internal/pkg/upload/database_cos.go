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

// DatabaseCOS 每次上传时读取数据库配置。未启用数据库 COS 时，回退到 YAML 的本地/COS 存储。
// 这样后台保存成功后无需重启 api-platform 即可在下一次上传生效。
type DatabaseCOS struct {
	Resolver ConfigResolver
	Fallback Store
}

func (d DatabaseCOS) Save(scope string, fh *multipart.FileHeader) (string, string, error) {
	if d.Resolver == nil {
		return d.fallback().Save(scope, fh)
	}
	cosValues, err := d.Resolver.Values(context.Background(), "cos")
	if err != nil {
		return "", "", fmt.Errorf("读取 COS 配置: %w", err)
	}
	if !enabled(cosValues["enabled"]) {
		return d.fallback().Save(scope, fh)
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

func (d DatabaseCOS) fallback() Store {
	if d.Fallback != nil {
		return d.Fallback
	}
	return Local{}
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
