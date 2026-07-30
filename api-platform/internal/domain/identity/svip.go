package identity

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// UserSvipActive is_svip: 1体验 2有效期 3永久；2 需未过期。
func UserSvipActive(u *User) bool {
	if u == nil {
		return false
	}
	switch u.IsSvip {
	case 1, 3:
		return true
	case 2:
		if u.SvipEndtime == nil {
			return false
		}
		return u.SvipEndtime.After(time.Now())
	default:
		return false
	}
}

type SvipInput struct {
	IsSvip      int8       `json:"is_svip"`
	SvipEndtime *time.Time `json:"svip_endtime"`
}

type UserPage struct {
	List  []User `json:"list"`
	Total int64  `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

func (s *Service) ListUsers(ctx context.Context, page, limit int) (*UserPage, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	list, total, err := s.store.ListUsers(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return &UserPage{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) SetUserSvip(ctx context.Context, uid uint, in SvipInput) (*User, error) {
	u, err := s.store.FindUserByID(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if in.IsSvip < -1 || in.IsSvip > 3 {
		return nil, ErrBadParam
	}
	u.IsSvip = in.IsSvip
	if in.IsSvip == 2 {
		if in.SvipEndtime == nil {
			t := time.Now().AddDate(0, 1, 0)
			u.SvipEndtime = &t
		} else {
			u.SvipEndtime = in.SvipEndtime
		}
	} else {
		u.SvipEndtime = nil
	}
	if err := s.store.UpdateUserSvip(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}
