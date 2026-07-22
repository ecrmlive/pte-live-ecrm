package distribution

import "time"

// SpreadLog maps qixi_user_spread_log.
type SpreadLog struct {
	UserSpreadLogID uint      `gorm:"column:user_spread_log_id;primaryKey" json:"user_spread_log_id"`
	UID             uint      `gorm:"column:uid" json:"uid"`
	OldSpreadUID    uint      `gorm:"column:old_spread_uid" json:"old_spread_uid"`
	SpreadUID       uint      `gorm:"column:spread_uid" json:"spread_uid"`
	CreateTime      time.Time `gorm:"column:create_time" json:"create_time"`
}

func (SpreadLog) TableName() string { return "qixi_user_spread_log" }

type BindInput struct {
	SpreadUID uint `json:"spread_uid"`
}

type MeInfo struct {
	UID         uint `json:"uid"`
	SpreadUID   uint `json:"spread_uid"`
	IsPromoter  uint8 `json:"is_promoter"`
	SpreadCount int64 `json:"spread_count"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
