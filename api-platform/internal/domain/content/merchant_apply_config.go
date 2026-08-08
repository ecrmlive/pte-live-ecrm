package content

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"
)

// merchantApplyConfigKey 平台「商户管理 → 商户设置」：入驻页背景 + 入驻表单字段。
const merchantApplyConfigKey = "merchant_apply_setting"

var allowedApplyFieldTypes = map[string]struct{}{
	"checkbox":  {},
	"city":      {},
	"date":      {},
	"daterange": {},
	"radio":     {},
	"select":    {},
	"text":      {},
	"textarea":  {},
	"time":      {},
	"timerange": {},
	"image":     {},
}

var allowedApplyContentTypes = map[string]struct{}{
	"text":   {},
	"number": {},
	"mobile": {},
	"idcard": {},
	"email":  {},
}

var allowedApplyCityLevels = map[string]struct{}{
	"province_city":                 {},
	"province_city_district":        {},
	"province_city_district_street": {},
}

var allowedApplyDefaultVisible = map[string]struct{}{
	"show": {},
	"hide": {},
}

var allowedApplyDefaultMode = map[string]struct{}{
	"current": {},
	"specify": {},
}

type merchantApplyField struct {
	ID             string   `json:"id"`
	Type           string   `json:"type"`
	Title          string   `json:"title"`
	ContentType    string   `json:"content_type"`
	DefaultValue   string   `json:"default_value"`
	Placeholder    string   `json:"placeholder"`
	Required       bool     `json:"required"`
	Options        []string `json:"options,omitempty"`
	MaxUpload      int      `json:"max_upload,omitempty"`
	CityLevel      string   `json:"city_level,omitempty"`
	DefaultVisible string   `json:"default_visible,omitempty"`
	DefaultMode    string   `json:"default_mode,omitempty"`
	SpecifyValue   string   `json:"specify_value,omitempty"`
}

type merchantApplyConfig struct {
	BackgroundImage string               `json:"background_image"`
	FormFields      []merchantApplyField `json:"form_fields"`
}

func defaultMerchantApplyConfig() merchantApplyConfig {
	return merchantApplyConfig{
		BackgroundImage: "",
		FormFields:      []merchantApplyField{},
	}
}

func marshalMerchantApplyConfig(config merchantApplyConfig) string {
	if config.FormFields == nil {
		config.FormFields = []merchantApplyField{}
	}
	data, _ := json.Marshal(config)
	return string(data)
}

func parseMerchantApplyConfig(raw string) (merchantApplyConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return merchantApplyConfig{}, ErrBadParam
	}
	for key := range fields {
		switch key {
		case "background_image", "form_fields":
		default:
			return merchantApplyConfig{}, ErrBadParam
		}
	}
	var config merchantApplyConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return merchantApplyConfig{}, ErrBadParam
	}
	config.BackgroundImage = strings.TrimSpace(config.BackgroundImage)
	if utf8.RuneCountInString(config.BackgroundImage) > 1024 {
		return merchantApplyConfig{}, ErrBadParam
	}
	normalized, err := normalizeApplyFormFields(config.FormFields)
	if err != nil {
		return merchantApplyConfig{}, err
	}
	config.FormFields = normalized
	return config, nil
}

func normalizeApplyFormFields(fields []merchantApplyField) ([]merchantApplyField, error) {
	if fields == nil {
		return []merchantApplyField{}, nil
	}
	if len(fields) > 50 {
		return nil, ErrBadParam
	}
	seen := make(map[string]struct{}, len(fields))
	normalized := make([]merchantApplyField, 0, len(fields))
	for _, field := range fields {
		id := strings.TrimSpace(field.ID)
		typ := strings.TrimSpace(field.Type)
		title := strings.TrimSpace(field.Title)
		contentType := strings.TrimSpace(field.ContentType)
		if contentType == "" {
			contentType = "text"
		}
		placeholder := strings.TrimSpace(field.Placeholder)
		defaultValue := strings.TrimSpace(field.DefaultValue)
		if id == "" || title == "" {
			return nil, ErrBadParam
		}
		if _, ok := allowedApplyFieldTypes[typ]; !ok {
			return nil, ErrBadParam
		}
		if _, ok := allowedApplyContentTypes[contentType]; !ok {
			return nil, ErrBadParam
		}
		if _, dup := seen[id]; dup {
			return nil, ErrBadParam
		}
		seen[id] = struct{}{}
		if utf8.RuneCountInString(title) > 64 ||
			utf8.RuneCountInString(placeholder) > 128 ||
			utf8.RuneCountInString(defaultValue) > 256 ||
			utf8.RuneCountInString(id) > 64 {
			return nil, ErrBadParam
		}
		opts := make([]string, 0, len(field.Options))
		for _, opt := range field.Options {
			opt = strings.TrimSpace(opt)
			if opt == "" {
				continue
			}
			if utf8.RuneCountInString(opt) > 64 {
				return nil, ErrBadParam
			}
			opts = append(opts, opt)
		}
		if len(opts) > 30 {
			return nil, ErrBadParam
		}
		switch typ {
		case "radio", "checkbox", "select":
			if len(opts) == 0 {
				opts = []string{"选项一", "选项二"}
			}
		default:
			opts = nil
		}

		maxUpload := 0
		cityLevel := ""
		defaultVisible := ""
		defaultMode := ""
		specifyValue := strings.TrimSpace(field.SpecifyValue)
		if utf8.RuneCountInString(specifyValue) > 64 {
			return nil, ErrBadParam
		}

		switch typ {
		case "image":
			maxUpload = field.MaxUpload
			if maxUpload <= 0 {
				maxUpload = 8
			}
			if maxUpload > 20 {
				return nil, ErrBadParam
			}
		case "city":
			cityLevel = strings.TrimSpace(field.CityLevel)
			if cityLevel == "" {
				cityLevel = "province_city_district"
			}
			if _, ok := allowedApplyCityLevels[cityLevel]; !ok {
				return nil, ErrBadParam
			}
		case "date", "daterange", "time", "timerange":
			defaultVisible = strings.TrimSpace(field.DefaultVisible)
			if defaultVisible == "" {
				defaultVisible = "show"
			}
			if _, ok := allowedApplyDefaultVisible[defaultVisible]; !ok {
				return nil, ErrBadParam
			}
			defaultMode = strings.TrimSpace(field.DefaultMode)
			if defaultMode == "" {
				defaultMode = "current"
			}
			if _, ok := allowedApplyDefaultMode[defaultMode]; !ok {
				return nil, ErrBadParam
			}
			if defaultMode != "specify" {
				specifyValue = ""
			}
		}

		normalized = append(normalized, merchantApplyField{
			ID:             id,
			Type:           typ,
			Title:          title,
			ContentType:    contentType,
			DefaultValue:   defaultValue,
			Placeholder:    placeholder,
			Required:       field.Required,
			Options:        opts,
			MaxUpload:      maxUpload,
			CityLevel:      cityLevel,
			DefaultVisible: defaultVisible,
			DefaultMode:    defaultMode,
			SpecifyValue:   specifyValue,
		})
	}
	return normalized, nil
}

func (s *Service) GetMerchantApplyConfig(ctx context.Context) (string, error) {
	row, err := s.store.GetCache(ctx, merchantApplyConfigKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return marshalMerchantApplyConfig(defaultMerchantApplyConfig()), nil
		}
		return "", err
	}
	config, err := parseMerchantApplyConfig(row.Result)
	if err != nil {
		return marshalMerchantApplyConfig(defaultMerchantApplyConfig()), nil
	}
	return marshalMerchantApplyConfig(config), nil
}

func (s *Service) SaveMerchantApplyConfig(ctx context.Context, raw string) (string, error) {
	config, err := parseMerchantApplyConfig(raw)
	if err != nil {
		return "", err
	}
	canonical := marshalMerchantApplyConfig(config)
	if err := s.store.UpsertCache(ctx, &Cache{Key: merchantApplyConfigKey, ExpireTime: 0, Result: canonical}); err != nil {
		return "", err
	}
	return canonical, nil
}
