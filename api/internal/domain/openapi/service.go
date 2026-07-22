package openapi

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/authjwt"
	"gorm.io/gorm"
)

type Store interface {
	GetByAccessKey(ctx context.Context, accessKey string) (*OpenAuth, error)
	TouchLogin(ctx context.Context, id uint, ip string, at time.Time) error
}

type Service struct {
	store  Store
	jwtMgr *authjwt.Manager
}

func NewService(store Store, jwtMgr *authjwt.Manager) *Service {
	return &Service{store: store, jwtMgr: jwtMgr}
}

func (s *Service) Authenticate(ctx context.Context, in AuthInput, clientIP string) (*AuthResult, error) {
	ak := strings.TrimSpace(in.AccessKey)
	if ak == "" {
		return nil, ErrBadParam
	}
	row, err := s.store.GetByAccessKey(ctx, ak)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}
	if row.IsDel == 1 || row.Status != 1 {
		return nil, ErrDisabled
	}
	if row.MerID == 0 {
		return nil, ErrUnauthorized
	}

	sig := strings.TrimSpace(in.Signature)
	skDemo := strings.TrimSpace(in.SecretKey)
	switch {
	case sig != "":
		if err := VerifySignature(ak, row.SecretKey, strings.TrimSpace(in.Unique), sig, in.Expiration, time.Now()); err != nil {
			return nil, err
		}
	case skDemo != "":
		// 本地演示捷径：直传 secret_key（生产应走 signature）
		if skDemo != row.SecretKey {
			return nil, ErrUnauthorized
		}
	default:
		return nil, ErrBadParam
	}

	pair, err := s.jwtMgr.Issue(authjwt.PortalOpen, 0, row.MerID, 0, row.AccessKey)
	if err != nil {
		return nil, err
	}
	_ = s.store.TouchLogin(ctx, row.ID, clientIP, time.Now())
	return &AuthResult{
		Token:     pair.AccessToken,
		Exp:       pair.ExpiresIn,
		MerID:     row.MerID,
		AccessKey: row.AccessKey,
	}, nil
}

func (s *Service) LoadAuth(ctx context.Context, accessKey string) (*OpenAuth, error) {
	row, err := s.store.GetByAccessKey(ctx, accessKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}
	if row.IsDel == 1 || row.Status != 1 {
		return nil, ErrDisabled
	}
	return row, nil
}

func (a *OpenAuth) AllowProduct() bool {
	return hasRule(a.Auth, AuthRuleProduct)
}

func (a *OpenAuth) AllowOrder() bool {
	return hasRule(a.Auth, AuthRuleOrder)
}

func hasRule(auth, rule string) bool {
	if strings.TrimSpace(auth) == "" {
		return true
	}
	for _, p := range strings.Split(auth, ",") {
		if strings.TrimSpace(p) == rule {
			return true
		}
	}
	return false
}
