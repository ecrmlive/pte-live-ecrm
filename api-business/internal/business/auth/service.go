package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrInvalidCredentials = errors.New("账号或密码错误")
	ErrAccountDisabled    = errors.New("账号已禁用")
	ErrAccountExists      = errors.New("账号已存在")
	ErrNotFound           = errors.New("用户不存在")
	ErrBadParam           = errors.New("参数错误")
	ErrCaptchaUnavailable = errors.New("安全验证已失效")
)

// LoginChannel 与 qixi_crm_b_user_identity.channel 的枚举保持一一对应。
type LoginChannel string

const (
	ChannelWechat      LoginChannel = "wechat"
	ChannelMiniProgram LoginChannel = "mini_program"
	ChannelH5          LoginChannel = "h5"
	ChannelPC          LoginChannel = "pc"
	ChannelIOS         LoginChannel = "ios"
	ChannelAndroid     LoginChannel = "android"
	ChannelHarmony     LoginChannel = "harmony"
)

func ParseChannel(raw string) (LoginChannel, error) {
	channel := LoginChannel(strings.TrimSpace(raw))
	switch channel {
	case ChannelWechat, ChannelMiniProgram, ChannelH5, ChannelPC, ChannelIOS, ChannelAndroid, ChannelHarmony:
		return channel, nil
	default:
		return "", ErrBadParam
	}
}

type User struct {
	ID       uint64 `gorm:"column:id;primaryKey" json:"id"`
	Nickname string `gorm:"column:nickname" json:"nickname"`
	// 手机号未绑定时必须写入 SQL NULL，而不是空字符串。mobile 有唯一索引，
	// 空字符串会让第二个无手机号账号注册失败。
	Mobile      *string `gorm:"column:mobile" json:"mobile"`
	Status      int8    `gorm:"column:status" json:"status"`
	AuthVersion uint64  `gorm:"column:auth_version" json:"-"`
}

func (User) TableName() string { return "qixi_crm_b_user" }

type Identity struct {
	ID             uint64       `gorm:"column:id;primaryKey"`
	UserID         uint64       `gorm:"column:user_id"`
	Channel        LoginChannel `gorm:"column:channel"`
	Subject        string       `gorm:"column:subject"`
	CredentialHash string       `gorm:"column:credential_hash"`
}

func (Identity) TableName() string { return "qixi_crm_b_user_identity" }

// CaptchaToken stores only a digest of the pte-tools-captcha token. It adds
// business-side one-time consumption; PTE itself validates signature/action.
type CaptchaToken struct {
	ID         uint64     `gorm:"column:id;primaryKey"`
	TokenHash  string     `gorm:"column:token_hash"`
	Action     string     `gorm:"column:action"`
	ExpiresAt  time.Time  `gorm:"column:expires_at"`
	ConsumedAt *time.Time `gorm:"column:consumed_at"`
}

func (CaptchaToken) TableName() string { return "qixi_crm_b_auth_captcha_token" }

type Profile struct {
	ID       uint64       `json:"id"`
	UID      uint64       `json:"uid"` // 兼容既有 C 端视图；值始终等于 id。
	Nickname string       `json:"nickname"`
	Mobile   string       `json:"mobile"`
	Channel  LoginChannel `json:"channel"`
	// Account is exposed under the stable C-end field name. The internal
	// identity subject remains the source for JWT issuance and may be a mobile
	// number, WeChat openid, or another channel-specific identifier.
	Subject     string `json:"account"`
	AuthVersion uint64 `json:"-"`
}

// StoreContext 来自业务库只读投影；api-business 不得跨库直接查询商户库。
type StoreContext struct {
	MerchantID    uint64 `gorm:"column:merchant_id" json:"merchant_id"`
	StoreID       uint64 `gorm:"column:store_id" json:"store_id"`
	MerchantAppID string `gorm:"column:merchant_app_id" json:"merchant_app_id"`
	IMSDKAppID    string `gorm:"column:im_sdk_app_id" json:"im_sdk_app_id"`
}

type Service struct{ db *gorm.DB }

func NewService(db *gorm.DB) *Service { return &Service{db: db} }

func mobileText(mobile *string) string {
	if mobile == nil {
		return ""
	}
	return *mobile
}

func (s *Service) Register(ctx context.Context, subject, password, nickname string, channel LoginChannel) (*Profile, error) {
	subject = strings.TrimSpace(subject)
	nickname = strings.TrimSpace(nickname)
	if subject == "" || len(password) < 6 {
		return nil, ErrBadParam
	}
	if nickname == "" {
		nickname = subject
	}
	runes := []rune(nickname)
	if len(runes) > 64 {
		nickname = string(runes[:64])
	}

	var created User
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing Identity
		err := tx.Where("channel = ? AND subject = ?", channel, subject).First(&existing).Error
		if err == nil {
			return ErrAccountExists
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		created = User{Nickname: nickname, Status: 1, AuthVersion: 1}
		if err := tx.Create(&created).Error; err != nil {
			return err
		}
		if err := tx.Create(&Identity{UserID: created.ID, Channel: channel, Subject: subject, CredentialHash: string(hash)}).Error; err != nil {
			return err
		}
		return grantOnboardingCoupons(tx, created.ID)
	})
	if err != nil {
		return nil, err
	}
	return &Profile{ID: created.ID, UID: created.ID, Nickname: created.Nickname, Mobile: mobileText(created.Mobile), Channel: channel, Subject: subject, AuthVersion: created.AuthVersion}, nil
}

func (s *Service) Login(ctx context.Context, subject, password string, channel LoginChannel) (*Profile, error) {
	subject = strings.TrimSpace(subject)
	var identity Identity
	if err := s.db.WithContext(ctx).Where("channel = ? AND subject = ?", channel, subject).First(&identity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	var user User
	if err := s.db.WithContext(ctx).First(&user, identity.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if user.Status != 1 {
		return nil, ErrAccountDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(identity.CredentialHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return &Profile{ID: user.ID, UID: user.ID, Nickname: user.Nickname, Mobile: mobileText(user.Mobile), Channel: channel, Subject: subject, AuthVersion: user.AuthVersion}, nil
}

// LoginOrRegisterExternal binds a provider-issued immutable subject (such as
// WeChat openid) to one C-end user. The caller cannot provide a user ID,
// password, or mobile number.
func (s *Service) LoginOrRegisterExternal(ctx context.Context, subject, nickname string, channel LoginChannel) (*Profile, error) {
	subject = strings.TrimSpace(subject)
	if subject == "" || (channel != ChannelMiniProgram && channel != ChannelWechat) {
		return nil, ErrBadParam
	}
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		nickname = "微信用户"
	}
	runes := []rune(nickname)
	if len(runes) > 64 {
		nickname = string(runes[:64])
	}

	var user User
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var identity Identity
		err := tx.Where("channel = ? AND subject = ?", channel, subject).First(&identity).Error
		if err == nil {
			if err := tx.First(&user, identity.UserID).Error; err != nil {
				return err
			}
			if user.Status != 1 {
				return ErrAccountDisabled
			}
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		user = User{Nickname: nickname, Status: 1, AuthVersion: 1}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if err := tx.Create(&Identity{UserID: user.ID, Channel: channel, Subject: subject}).Error; err != nil {
			return err
		}
		return grantOnboardingCoupons(tx, user.ID)
	})
	if err != nil {
		return nil, err
	}
	return &Profile{ID: user.ID, UID: user.ID, Nickname: user.Nickname, Mobile: mobileText(user.Mobile), Channel: channel, Subject: subject, AuthVersion: user.AuthVersion}, nil
}

func captchaTokenHash(token string) string {
	digest := sha256.Sum256([]byte(strings.TrimSpace(token)))
	return hex.EncodeToString(digest[:])
}

// RecordCaptchaToken records a pte-issued proof before it is returned to the
// client. The source token never enters MySQL; a later action can consume it once.
func (s *Service) RecordCaptchaToken(ctx context.Context, token, action string) error {
	if strings.TrimSpace(token) == "" || strings.TrimSpace(action) == "" {
		return ErrBadParam
	}
	row := CaptchaToken{TokenHash: captchaTokenHash(token), Action: action, ExpiresAt: time.Now().Add(10 * time.Minute)}
	return s.db.WithContext(ctx).Create(&row).Error
}

func (s *Service) ConsumeCaptchaToken(ctx context.Context, token, action string) error {
	if strings.TrimSpace(token) == "" || strings.TrimSpace(action) == "" {
		return ErrCaptchaUnavailable
	}
	now := time.Now()
	result := s.db.WithContext(ctx).Model(&CaptchaToken{}).
		Where("token_hash = ? AND action = ? AND consumed_at IS NULL AND expires_at > ?", captchaTokenHash(token), action, now).
		Update("consumed_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return ErrCaptchaUnavailable
	}
	return nil
}

func normalizeNickname(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if len([]rune(value)) == 0 || len([]rune(value)) > 64 {
		return "", ErrBadParam
	}
	return value, nil
}

func (s *Service) UpdateNickname(ctx context.Context, userID uint64, nickname string) error {
	value, err := normalizeNickname(nickname)
	if err != nil {
		return err
	}
	result := s.db.WithContext(ctx).Model(&User{}).Where("id=? AND status=1", userID).Update("nickname", value)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return ErrNotFound
	}
	return nil
}

func (s *Service) ResetPasswordByMobile(ctx context.Context, mobile, newPassword string) error {
	mobile = strings.TrimSpace(mobile)
	if !validMobile(mobile) || !validNewPassword(newPassword) {
		return ErrBadParam
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("mobile = ?", mobile).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if user.Status != 1 {
			return ErrAccountDisabled
		}
		var identity Identity
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND channel = ?", user.ID, ChannelH5).First(&identity).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			identity = Identity{UserID: user.ID, Channel: ChannelH5, Subject: "mobile:" + mobile}
			if err = tx.Create(&identity).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if err = tx.Model(&Identity{}).Where("id=?", identity.ID).Update("credential_hash", string(hash)).Error; err != nil {
			return err
		}
		result := tx.Model(&User{}).Where("id=? AND auth_version=?", user.ID, user.AuthVersion).Update("auth_version", user.AuthVersion+1)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrNotFound
		}
		return nil
	})
}

func (s *Service) LoginOrRegisterMobile(ctx context.Context, mobile string, channel LoginChannel) (*Profile, error) {
	mobile = strings.TrimSpace(mobile)
	if !validMobile(mobile) || channel != ChannelH5 {
		return nil, ErrBadParam
	}
	var user User
	var identity Identity
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("mobile = ?", mobile).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			value := mobile
			user = User{Nickname: "用户" + mobile[len(mobile)-4:], Mobile: &value, Status: 1, AuthVersion: 1}
			if err = tx.Create(&user).Error; err != nil {
				return err
			}
			if err = grantOnboardingCoupons(tx, user.ID); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else if user.Status != 1 {
			return ErrAccountDisabled
		}
		err = tx.Where("user_id = ? AND channel = ?", user.ID, channel).First(&identity).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			identity = Identity{UserID: user.ID, Channel: channel, Subject: "mobile:" + mobile}
			return tx.Create(&identity).Error
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return &Profile{ID: user.ID, UID: user.ID, Nickname: user.Nickname, Mobile: mobileText(user.Mobile), Channel: identity.Channel, Subject: identity.Subject, AuthVersion: user.AuthVersion}, nil
}

func (s *Service) BindMobile(ctx context.Context, userID uint64, mobile string) error {
	mobile = strings.TrimSpace(mobile)
	if userID == 0 || !validMobile(mobile) {
		return ErrBadParam
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var other User
		err := tx.Where("mobile = ?", mobile).First(&other).Error
		if err == nil && other.ID != userID {
			return ErrAccountExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		result := tx.Model(&User{}).Where("id = ? AND status = ?", userID, 1).Update("mobile", mobile)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrNotFound
		}
		return nil
	})
}

func (s *Service) Profile(ctx context.Context, userID uint64, channel LoginChannel) (*Profile, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	var identity Identity
	if err := s.db.WithContext(ctx).Where("user_id = ? AND channel = ?", userID, channel).First(&identity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &Profile{ID: user.ID, UID: user.ID, Nickname: user.Nickname, Mobile: mobileText(user.Mobile), Channel: identity.Channel, Subject: identity.Subject, AuthVersion: user.AuthVersion}, nil
}

func (s *Service) ResolveStoreContext(ctx context.Context, merchantAppID string) (*StoreContext, error) {
	merchantAppID = strings.TrimSpace(merchantAppID)
	if merchantAppID == "" {
		return nil, ErrBadParam
	}
	var row StoreContext
	err := s.db.WithContext(ctx).Table("qixi_crm_b_store_view AS s").
		Select("s.merchant_id,s.store_id,s.store_app_id AS merchant_app_id,COALESCE(im.sdk_app_id, '') AS im_sdk_app_id").
		Joins("LEFT JOIN qixi_crm_b_merchant_im_sdk_app_view AS im ON im.merchant_id = s.merchant_id").
		Where("s.store_app_id = ? AND s.status = 1", merchantAppID).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &row, err
}
