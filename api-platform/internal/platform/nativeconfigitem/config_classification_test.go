package nativeconfigitem

import "testing"

func TestNormalizeConfigClassificationInput(t *testing.T) {
	name, key, description, icon, status, sort, err := normalizeConfigClassificationInput(configClassificationInput{
		Name: "短信配置", ClassifyKey: "tencent.sms", Description: "平台短信配置", Icon: "lucide:mail", Status: groupDataIntPtr(1), Sort: groupDataIntPtr(10),
	}, 0, 0)
	if err != nil || name != "短信配置" || key != "tencent.sms" || description != "平台短信配置" || icon != "lucide:mail" || status != 1 || sort != 10 {
		t.Fatalf("unexpected normalized classification: %q %q %q %q %d %d %v", name, key, description, icon, status, sort, err)
	}
	if _, _, _, _, _, _, err = normalizeConfigClassificationInput(configClassificationInput{Name: "测试", ClassifyKey: "中文"}, 1, 0); err == nil {
		t.Fatal("expected invalid key error")
	}
}

func TestNormalizeConfigClassificationItemInput(t *testing.T) {
	name, key, fieldType, backendType, content, description, status, sort, err := normalizeConfigClassificationItemInput(configClassificationItemInput{
		Name: "是否启用", ConfigKey: "enabled", Content: "true", Description: "开启默认行为", Status: groupDataIntPtr(1), Sort: groupDataIntPtr(2),
	}, "input", 0, 0, 0)
	if err != nil || name != "是否启用" || key != "enabled" || fieldType != "input" || backendType != 0 || content != "true" || description != "开启默认行为" || status != 1 || sort != 2 {
		t.Fatalf("unexpected normalized item: %q %q %q %d %q %q %d %d %v", name, key, fieldType, backendType, content, description, status, sort, err)
	}
}
