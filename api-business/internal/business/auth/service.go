package auth

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("账号或密码错误")
	ErrAccountDisabled    = errors.New("账号已禁用")
	ErrAccountExists      = errors.New("账号已存在")
	ErrNotFound           = errors.New("用户不存在")
	ErrBadParam           = errors.New("参数错误")
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

type Profile struct {
	ID          uint64       `json:"id"`
	UID         uint64       `json:"uid"` // 兼容既有 C 端视图；值始终等于 id。
	Nickname    string       `json:"nickname"`
	Mobile      string       `json:"mobile"`
	Channel     LoginChannel `json:"channel"`
	Subject     string       `json:"subject"`
	AuthVersion uint64       `json:"-"`
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
		return tx.Create(&Identity{UserID: created.ID, Channel: channel, Subject: subject, CredentialHash: string(hash)}).Error
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
