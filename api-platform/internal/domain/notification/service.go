package notification

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"
)

type Store interface {
	List(ctx context.Context, audience Audience, page, limit int) ([]Config, int64, error)
	Get(ctx context.Context, id uint) (*Config, error)
	Save(ctx context.Context, config *Config) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) List(ctx context.Context, audience Audience, page, limit int) (*PageResult, error) {
	if !validAudience(audience) {
		return nil, ErrBadParam
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	list, total, err := s.store.List(ctx, audience, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) Get(ctx context.Context, id uint) (*Config, error) {
	if id == 0 {
		return nil, ErrBadParam
	}
	config, err := s.store.Get(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return config, err
}

func (s *Service) Save(ctx context.Context, id uint, input SaveInput) (*Config, error) {
	config, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !validFlag(input.WechatEnabled) || !validFlag(input.MiniEnabled) || !validFlag(input.SMSEnabled) {
		return nil, ErrBadParam
	}
	wechatText := strings.TrimSpace(input.WechatText)
	miniText := strings.TrimSpace(input.MiniText)
	smsText := strings.TrimSpace(input.SMSText)
	if utf8.RuneCountInString(wechatText) > 500 || utf8.RuneCountInString(miniText) > 500 || utf8.RuneCountInString(smsText) > 500 {
		return nil, ErrBadParam
	}
	if (input.WechatEnabled == 1 && wechatText == "") || (input.MiniEnabled == 1 && miniText == "") || (input.SMSEnabled == 1 && smsText == "") {
		return nil, ErrBadParam
	}
	config.WechatEnabled = input.WechatEnabled
	config.MiniEnabled = input.MiniEnabled
	config.SMSEnabled = input.SMSEnabled
	config.WechatText = wechatText
	config.MiniText = miniText
	config.SMSText = smsText
	if err := s.store.Save(ctx, config); err != nil {
		return nil, err
	}
	return config, nil
}

// Sync 目前不接入微信开放平台。保留动作端点，明确返回不可同步，而不是伪造成功。
func (s *Service) Sync(_ context.Context, audience Audience, channel Channel) error {
	if !validAudience(audience) || (channel != ChannelWechat && channel != ChannelMiniProgram) {
		return ErrBadParam
	}
	return ErrSyncUnavailable
}

func validAudience(audience Audience) bool {
	return audience == AudienceMember || audience == AudienceStore
}

func validFlag(value int8) bool { return value == 0 || value == 1 }
