package content

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"
)

// WechatMenusCacheKey 对齐 CRMEB CacheRepository wechat_menus（菜单 button JSON 数组）。
const WechatMenusCacheKey = "wechat_menus"

// WechatMenuButton 微信自定义菜单按钮（官方 button / sub_button 结构）。
type WechatMenuButton struct {
	Type      string             `json:"type,omitempty"`
	Name      string             `json:"name"`
	Key       string             `json:"key,omitempty"`
	URL       string             `json:"url,omitempty"`
	AppID     string             `json:"appid,omitempty"`
	PagePath  string             `json:"pagepath,omitempty"`
	SubButton []WechatMenuButton `json:"sub_button,omitempty"`
}

func defaultWechatMenus() []WechatMenuButton {
	return []WechatMenuButton{
		{
			Name: "商城",
			Type: "view",
			URL:  "https://mer.crmeb.net",
		},
		{
			Name: "一级菜单",
			Type: "click",
			Key:  "MENU_LEVEL_1",
		},
	}
}

func (s *Service) GetWechatMenus(ctx context.Context) ([]WechatMenuButton, error) {
	row, err := s.store.GetCache(ctx, WechatMenusCacheKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return defaultWechatMenus(), nil
		}
		return nil, err
	}
	buttons, err := parseWechatMenus(row.Result)
	if err != nil {
		return defaultWechatMenus(), nil
	}
	return buttons, nil
}

func (s *Service) SaveWechatMenus(ctx context.Context, buttons []WechatMenuButton) ([]WechatMenuButton, error) {
	cleaned, err := normalizeWechatMenus(buttons)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(cleaned)
	if err != nil {
		return nil, ErrBadParam
	}
	if err := s.store.UpsertCache(ctx, &Cache{Key: WechatMenusCacheKey, ExpireTime: 0, Result: string(data)}); err != nil {
		return nil, err
	}
	return cleaned, nil
}

func parseWechatMenus(raw string) ([]WechatMenuButton, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "null" {
		return nil, ErrBadParam
	}
	var buttons []WechatMenuButton
	if err := json.Unmarshal([]byte(raw), &buttons); err != nil {
		return nil, ErrBadParam
	}
	return normalizeWechatMenus(buttons)
}

func normalizeWechatMenus(buttons []WechatMenuButton) ([]WechatMenuButton, error) {
	if len(buttons) == 0 {
		return nil, ErrBadParam
	}
	if len(buttons) > 3 {
		return nil, ErrBadParam
	}
	out := make([]WechatMenuButton, 0, len(buttons))
	for _, top := range buttons {
		item, err := normalizeWechatMenuButton(top, true)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func normalizeWechatMenuButton(btn WechatMenuButton, topLevel bool) (WechatMenuButton, error) {
	name := strings.TrimSpace(btn.Name)
	if name == "" || utf8.RuneCountInString(name) > 16 {
		return WechatMenuButton{}, ErrBadParam
	}
	subs := btn.SubButton
	if len(subs) > 5 {
		return WechatMenuButton{}, ErrBadParam
	}
	if len(subs) > 0 {
		if !topLevel {
			return WechatMenuButton{}, ErrBadParam
		}
		children := make([]WechatMenuButton, 0, len(subs))
		for _, child := range subs {
			item, err := normalizeWechatMenuButton(child, false)
			if err != nil {
				return WechatMenuButton{}, err
			}
			children = append(children, item)
		}
		return WechatMenuButton{Name: name, SubButton: children}, nil
	}

	typ := strings.TrimSpace(btn.Type)
	switch typ {
	case "click":
		key := strings.TrimSpace(btn.Key)
		if key == "" || utf8.RuneCountInString(key) > 128 {
			return WechatMenuButton{}, ErrBadParam
		}
		return WechatMenuButton{Type: typ, Name: name, Key: key}, nil
	case "view":
		u := strings.TrimSpace(btn.URL)
		if !validWechatMenuURL(u) {
			return WechatMenuButton{}, ErrBadParam
		}
		return WechatMenuButton{Type: typ, Name: name, URL: u}, nil
	case "miniprogram":
		appID := strings.TrimSpace(btn.AppID)
		pagePath := strings.TrimSpace(btn.PagePath)
		u := strings.TrimSpace(btn.URL)
		if appID == "" || pagePath == "" || !validWechatMenuURL(u) ||
			utf8.RuneCountInString(appID) > 64 || utf8.RuneCountInString(pagePath) > 256 {
			return WechatMenuButton{}, ErrBadParam
		}
		return WechatMenuButton{
			Type: typ, Name: name, AppID: appID, PagePath: pagePath, URL: u,
		}, nil
	default:
		return WechatMenuButton{}, ErrBadParam
	}
}

func validWechatMenuURL(raw string) bool {
	if raw == "" || utf8.RuneCountInString(raw) > 1024 {
		return false
	}
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return u.Host != ""
}
