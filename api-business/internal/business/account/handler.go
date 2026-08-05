// Package account exposes the C-end asset read model. It only reads the
// business-owned member account and immutable asset ledger.
package account

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/account/balance", h.Balance)
	r.GET("/account/points", h.Points)
	r.GET("/member/info", h.MemberInfo)
	r.GET("/member/log", h.MemberLog)
	r.GET("/sign/info", h.SignInfo)
	r.GET("/sign/lst", h.SignList)
	r.GET("/sign/month", h.SignList)
	r.POST("/sign/create", h.SignCreate)
}

type balanceSummary struct {
	Balance      float64 `json:"balance"`
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
}
type ledgerRow struct {
	ID            uint64  `json:"id"`
	Amount        float64 `json:"amount"`
	ReferenceType string  `json:"reference_type"`
	ReferenceID   string  `json:"reference_id"`
	CreatedAt     string  `json:"created_at"`
}

type pointsSummary struct {
	Points       int64 `json:"points"`
	TotalIncome  int64 `json:"total_income"`
	TotalExpense int64 `json:"total_expense"`
	FrozenPoints int64 `json:"frozen_points"`
}

type pointsLedgerRow struct {
	ID            uint64 `json:"id"`
	Amount        int64  `json:"amount"`
	ReferenceType string `json:"reference_type"`
	ReferenceID   string `json:"reference_id"`
	CreatedAt     string `json:"created_at"`
}

func ledgerPageParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

func (h *Handler) Balance(c *gin.Context) {
	uid := middleware.UID(c)
	summary := balanceSummary{}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_member_account").Select("COALESCE(balance, 0) AS balance").Where("user_id = ?", uid).Scan(&summary).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询余额失败")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_asset_ledger").Select("COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END),0) AS total_income, COALESCE(SUM(CASE WHEN amount < 0 THEN -amount ELSE 0 END),0) AS total_expense").Where("user_id = ? AND asset_type = ?", uid, "balance").Scan(&summary).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询余额流水失败")
		return
	}
	page, limit := ledgerPageParams(c)
	ledger := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_asset_ledger").Where("user_id = ? AND asset_type = ?", uid, "balance")
	var total int64
	if err := ledger.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询余额流水失败")
		return
	}
	var rows []ledgerRow
	if err := ledger.Select("id, amount, reference_type, reference_id, created_at").Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询余额流水失败")
		return
	}
	response.OK(c, gin.H{"summary": summary, "list": rows, "total": total, "page": page, "limit": limit})
}

// Points exposes the same immutable-ledger contract as balance, but keeps
// point quantities integer-valued and never rounds decimal asset records.
func (h *Handler) Points(c *gin.Context) {
	uid := middleware.UID(c)
	summary := pointsSummary{}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_member_account").
		Select("COALESCE(points, 0) AS points").Where("user_id = ?", uid).Scan(&summary).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询积分失败")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_asset_ledger").
		Select("COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END),0) AS total_income, COALESCE(SUM(CASE WHEN amount < 0 THEN -amount ELSE 0 END),0) AS total_expense").
		Where("user_id = ? AND asset_type = ?", uid, "points").Scan(&summary).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询积分流水失败")
		return
	}
	page, limit := ledgerPageParams(c)
	ledger := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_asset_ledger").Where("user_id = ? AND asset_type = ?", uid, "points")
	var total int64
	if err := ledger.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询积分流水失败")
		return
	}
	var rows []pointsLedgerRow
	if err := ledger.
		Select("id, amount, reference_type, reference_id, created_at").
		Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询积分流水失败")
		return
	}
	response.OK(c, gin.H{"summary": summary, "list": rows, "total": total, "page": page, "limit": limit})
}

type memberInfo struct {
	LevelID   uint   `json:"level_id"`
	LevelName string `json:"level_name"`
	Rank      int    `json:"rank"`
	Rules     string `json:"rules"`
	Benefits  string `json:"benefits"`
}

type memberLevel struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Rank     int    `json:"rank"`
	Rules    string `json:"rules"`
	Benefits string `json:"benefits"`
}

type memberLogRow struct {
	ID                uint64 `json:"id"`
	LevelName         string `json:"level_name"`
	PreviousLevelName string `json:"previous_level_name"`
	ChangeType        string `json:"change_type"`
	Note              string `json:"note"`
	CreatedAt         string `json:"created_at"`
}

// MemberInfo returns the current membership grade and currently enabled grades.
// Benefits and rules remain server-authored JSON strings, so the client never
// invents eligibility or purchase privileges.
func (h *Handler) MemberInfo(c *gin.Context) {
	uid := middleware.UID(c)
	current := memberInfo{LevelName: "普通会员", Rules: "[]", Benefits: "[]"}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_member_account AS a").
		Select("COALESCE(a.level_id, 0) AS level_id, COALESCE(l.name, '普通会员') AS level_name, COALESCE(l.rank, 0) AS rank, COALESCE(l.rules, '[]') AS rules, COALESCE(l.benefits, '[]') AS benefits").
		Joins("LEFT JOIN qixi_crm_b_member_level AS l ON l.id = a.level_id AND l.status = 1").
		Where("a.user_id = ?", uid).Scan(&current).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询会员信息失败")
		return
	}
	var levels []memberLevel
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_member_level").
		Select("id, name, rank, rules, benefits").Where("status = 1").Order("rank ASC, id ASC").Find(&levels).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询会员等级失败")
		return
	}
	response.OK(c, gin.H{"current": current, "levels": levels})
}

// MemberLog reads immutable membership grade changes for the current user only.
func (h *Handler) MemberLog(c *gin.Context) {
	uid := middleware.UID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	const limit = 20
	base := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_member_level_log AS log").Where("log.user_id = ?", uid)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询会员成长记录失败")
		return
	}
	var rows []memberLogRow
	if err := base.
		Select("log.id, COALESCE(current_level.name, '普通会员') AS level_name, COALESCE(previous_level.name, '') AS previous_level_name, log.change_type, log.note, log.created_at").
		Joins("LEFT JOIN qixi_crm_b_member_level AS current_level ON current_level.id = log.level_id").
		Joins("LEFT JOIN qixi_crm_b_member_level AS previous_level ON previous_level.id = log.previous_level_id").
		Order("log.id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询会员成长记录失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

var errSignAlready = errors.New("daily sign already exists")

const dailySignReward int64 = 5

type signRow struct {
	ID             uint64 `json:"id"`
	SignDate       string `json:"sign_date"`
	Points         int64  `json:"points"`
	ContinuousDays int    `json:"continuous_days"`
	CreatedAt      string `json:"created_at"`
}

func shanghaiToday() string {
	return time.Now().In(time.FixedZone("CST", 8*60*60)).Format("2006-01-02")
}

func (h *Handler) SignInfo(c *gin.Context) {
	uid, today := middleware.UID(c), shanghaiToday()
	row := signRow{}
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_user_sign").
		Select("id, DATE_FORMAT(sign_date, '%Y-%m-%d') AS sign_date, points, continuous_days, created_at").
		Where("user_id = ? AND sign_date = ?", uid, today).Take(&row).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusInternalServerError, "查询签到状态失败")
		return
	}
	response.OK(c, gin.H{"server_date": today, "signed_today": err == nil, "today_points": row.Points, "daily_reward": dailySignReward, "continuous_days": row.ContinuousDays})
}

func (h *Handler) SignList(c *gin.Context) {
	uid := middleware.UID(c)
	month := c.DefaultQuery("month", shanghaiToday()[:7])
	if _, err := time.Parse("2006-01", month); err != nil {
		response.Fail(c, http.StatusBadRequest, "月份格式应为 YYYY-MM")
		return
	}
	var rows []signRow
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_user_sign").
		Select("id, DATE_FORMAT(sign_date, '%Y-%m-%d') AS sign_date, points, continuous_days, created_at").
		Where("user_id = ? AND DATE_FORMAT(sign_date, '%Y-%m') = ?", uid, month).Order("sign_date DESC").Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询签到记录失败")
		return
	}
	response.OK(c, gin.H{"month": month, "list": rows})
}

// SignCreate grants the fixed daily sign-in reward exactly once. The user-day
// unique key and asset-ledger idempotency key protect concurrent requests.
func (h *Handler) SignCreate(c *gin.Context) {
	uid, today := middleware.UID(c), shanghaiToday()
	var row signRow
	already := false
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		err := tx.Table("qixi_crm_b_user_sign").Where("user_id = ? AND sign_date = ?", uid, today).Take(&row).Error
		if err == nil {
			already = true
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		previous := signRow{}
		if err := tx.Table("qixi_crm_b_user_sign").Select("continuous_days").Where("user_id = ? AND sign_date = DATE_SUB(?, INTERVAL 1 DAY)", uid, today).Take(&previous).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		row = signRow{SignDate: today, Points: dailySignReward, ContinuousDays: previous.ContinuousDays + 1}
		if err := tx.Table("qixi_crm_b_user_sign").Create(map[string]any{"user_id": uid, "sign_date": today, "points": dailySignReward, "continuous_days": row.ContinuousDays}).Error; err != nil {
			return errSignAlready
		}
		if err := tx.Exec("INSERT INTO qixi_crm_b_member_account (user_id, points) VALUES (?, ?) ON DUPLICATE KEY UPDATE points = points + VALUES(points)", uid, dailySignReward).Error; err != nil {
			return err
		}
		key := "daily-sign:" + strconv.FormatUint(uint64(uid), 10) + ":" + today
		return tx.Exec("INSERT INTO qixi_crm_b_asset_ledger (user_id, asset_type, amount, reference_type, reference_id, idempotency_key) VALUES (?, 'points', ?, 'sign_in', ?, ?)", uid, dailySignReward, today, key).Error
	})
	if err != nil {
		if errors.Is(err, errSignAlready) {
			if findErr := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_user_sign").Select("id, DATE_FORMAT(sign_date, '%Y-%m-%d') AS sign_date, points, continuous_days, created_at").Where("user_id = ? AND sign_date = ?", uid, today).Take(&row).Error; findErr == nil {
				response.OK(c, gin.H{"already_signed": true, "signed_today": true, "points": row.Points, "continuous_days": row.ContinuousDays, "server_date": today})
				return
			}
		}
		response.Fail(c, http.StatusInternalServerError, "签到失败，请稍后重试")
		return
	}
	response.OK(c, gin.H{"already_signed": already, "signed_today": true, "points": row.Points, "continuous_days": row.ContinuousDays, "server_date": today})
}
