package content

import (
	"context"
	"encoding/json"
	"strings"
)

// PageVisualizationConfigKey 存储商城首页、签到及开屏广告的可视化配置。
const PageVisualizationConfigKey = "page_visualization_config"

type pageVisualizationImage struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Link string `json:"link"`
}

type pageVisualizationReward struct {
	Day     int  `json:"day"`
	Points  int  `json:"points"`
	Sort    int  `json:"sort"`
	Enabled bool `json:"enabled"`
}

type pageVisualizationRange struct {
	Name    string `json:"name"`
	Min     int    `json:"min"`
	Max     int    `json:"max"`
	Sort    int    `json:"sort"`
	Enabled bool   `json:"enabled"`
}

type pageVisualizationSplash struct {
	Enabled         bool                     `json:"enabled"`
	DisplaySeconds  int                      `json:"display_seconds"`
	IntervalHours   int                      `json:"interval_hours"`
	Images          []pageVisualizationImage `json:"images"`
}

type pageVisualizationConfig struct {
	Carousels    map[string][]pageVisualizationImage `json:"carousels"`
	SignRewards  []pageVisualizationReward           `json:"sign_rewards"`
	PointRanges  []pageVisualizationRange            `json:"point_ranges"`
	Splash       pageVisualizationSplash             `json:"splash"`
}

func defaultPageVisualizationConfig() pageVisualizationConfig {
	return pageVisualizationConfig{
		Carousels: map[string][]pageVisualizationImage{
			"home": {}, "hot": {}, "featured": {}, "recommend": {}, "points": {},
		},
		SignRewards: []pageVisualizationReward{},
		PointRanges: []pageVisualizationRange{},
		Splash: pageVisualizationSplash{Enabled: false, DisplaySeconds: 3, IntervalHours: 0, Images: []pageVisualizationImage{}},
	}
}

func parsePageVisualizationConfig(raw string) (pageVisualizationConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return pageVisualizationConfig{}, ErrBadParam
	}
	for key := range fields {
		if isSensitiveConfigKey(key) || (key != "carousels" && key != "sign_rewards" && key != "point_ranges" && key != "splash") {
			return pageVisualizationConfig{}, ErrBadParam
		}
	}
	config := defaultPageVisualizationConfig()
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return pageVisualizationConfig{}, ErrBadParam
	}
	allowedCarouselKeys := map[string]struct{}{"home": {}, "hot": {}, "featured": {}, "recommend": {}, "points": {}}
	for key := range config.Carousels {
		if _, ok := allowedCarouselKeys[key]; !ok {
			return pageVisualizationConfig{}, ErrBadParam
		}
	}
	for _, key := range []string{"home", "hot", "featured", "recommend", "points"} {
		if _, ok := config.Carousels[key]; !ok {
			config.Carousels[key] = []pageVisualizationImage{}
		}
		if err := validatePageVisualizationImages(config.Carousels[key]); err != nil {
			return pageVisualizationConfig{}, err
		}
	}
	if err := validatePageVisualizationImages(config.Splash.Images); err != nil || config.Splash.DisplaySeconds < 1 || config.Splash.DisplaySeconds > 30 || config.Splash.IntervalHours < 0 || config.Splash.IntervalHours > 8760 {
		return pageVisualizationConfig{}, ErrBadParam
	}
	if len(config.SignRewards) > 31 || len(config.PointRanges) > 50 {
		return pageVisualizationConfig{}, ErrBadParam
	}
	for _, reward := range config.SignRewards {
		if reward.Day < 1 || reward.Day > 366 || reward.Points < 0 || reward.Sort < 0 {
			return pageVisualizationConfig{}, ErrBadParam
		}
	}
	for _, item := range config.PointRanges {
		item.Name = strings.TrimSpace(item.Name)
		if item.Name == "" || len([]rune(item.Name)) > 64 || item.Min < 0 || item.Max < item.Min || item.Sort < 0 {
			return pageVisualizationConfig{}, ErrBadParam
		}
	}
	return config, nil
}

func validatePageVisualizationImages(images []pageVisualizationImage) error {
	if len(images) > 5 {
		return ErrBadParam
	}
	for _, image := range images {
		if len([]rune(strings.TrimSpace(image.Name))) > 64 || len([]rune(strings.TrimSpace(image.URL))) > 1024 || len([]rune(strings.TrimSpace(image.Link))) > 1024 {
			return ErrBadParam
		}
	}
	return nil
}

func (s *Service) GetPageVisualizationConfig(ctx context.Context) (string, error) {
	return getJSONConfig(s, ctx, PageVisualizationConfigKey, defaultPageVisualizationConfig(), parsePageVisualizationConfig)
}

func (s *Service) SavePageVisualizationConfig(ctx context.Context, raw string) (string, error) {
	return saveJSONConfig(s, ctx, PageVisualizationConfigKey, raw, parsePageVisualizationConfig)
}
