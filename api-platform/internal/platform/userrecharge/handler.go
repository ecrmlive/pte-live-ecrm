// Package userrecharge exposes platform user balance recharge records
// (充值记录). Aligns with CRMEB UserRecharge list/total; refund is a
// platform admin action against qixi_crm_b_user_recharge + member balance.
package userrecharge

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/queryfilter"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	menuRead   = "accounts.recharge_record.read"
	menuRefund = "accounts.recharge_record.refund"
	tableName  = "qixi_crm_b_user_recharge"
)

type Handler struct {
	businessDB *gorm.DB
	adminDB    *gorm.DB
}

func NewHandler(businessDB, adminDB *gorm.DB) *Handler {
	return &Handler{businessDB: businessDB, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	access := middleware.RequireAdminRoles("platform", "operations")
	read := middleware.RequireAdminMenu(h.adminDB, menuRead)
	refund := middleware.RequireAdminMenu(h.adminDB, menuRefund)
	r.GET("/finance/user-recharges", access, read, h.List)
	r.GET("/finance/user-recharges/total", access, read, h.Total)
	r.POST("/finance/user-recharges/:id/refund", access, refund, h.Refund)
}

type listFilter struct {
	DateFrom     string
	DateTo       string
	Paid         *int
	RechargeType string // "" | "1"(微信系) | "2"(支付宝) | concrete type
	UserType     string
	UserKeyword  string
	OrderID      string
}

type row struct {
	RechargeID       uint64     `gorm:"column:recharge_id" json:"recharge_id"`
	UID              uint64     `gorm:"column:uid" json:"uid"`
	Nickname         string     `gorm:"column:nickname" json:"nickname"`
	Avatar           string     `gorm:"column:avatar" json:"avatar"`
	RealName         string     `gorm:"column:real_name" json:"real_name"`
	OrderID          string     `gorm:"column:order_id" json:"order_id"`
	Price            float64    `gorm:"column:price" json:"price"`
	GivePrice        float64    `gorm:"column:give_price" json:"give_price"`
	RechargeType     string     `gorm:"column:recharge_type" json:"recharge_type"`
	RechargeTypeName string     `gorm:"-" json:"recharge_type_name"`
	Paid             int        `gorm:"column:paid" json:"paid"`
	PaidName         string     `gorm:"-" json:"paid_name"`
	PayTime          *time.Time `gorm:"column:pay_time" json:"pay_time"`
	RefundPrice      float64    `gorm:"column:refund_price" json:"refund_price"`
	CanRefund        bool       `gorm:"-" json:"can_refund"`
	CreateTime       *time.Time `gorm:"column:create_time" json:"create_time"`
}

func (h *Handler) List(c *gin.Context) {
	page, limit := queryfilter.Page(c)
	f := parseFilter(c)
	var total int64
	if err := h.baseQuery(c, f).Count(&total).Error; err != nil {
		fail(c, "充值记录查询失败")
		return
	}
	rows := make([]row, 0)
	err := h.baseQuery(c, f).Select(`r.recharge_id,r.uid,COALESCE(u.nickname,'') AS nickname,
COALESCE(p.avatar_url,'') AS avatar,COALESCE(p.real_name,'') AS real_name,
r.order_id,r.price,r.give_price,r.recharge_type,r.paid,r.pay_time,r.refund_price,r.create_time`).
		Order("r.pay_time DESC, r.create_time DESC, r.recharge_id DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		fail(c, "充值记录查询失败")
		return
	}
	for i := range rows {
		decorate(&rows[i])
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Total(c *gin.Context) {
	ctx := c.Request.Context()
	var totalPay, totalRoutine, totalWx, totalRefund float64
	_ = h.businessDB.WithContext(ctx).Table(tableName).
		Where("paid = 1").Select("COALESCE(SUM(price),0)").Scan(&totalPay).Error
	_ = h.businessDB.WithContext(ctx).Table(tableName).
		Where("paid = 1 AND recharge_type = ?", "routine").
		Select("COALESCE(SUM(price),0)").Scan(&totalRoutine).Error
	_ = h.businessDB.WithContext(ctx).Table(tableName).
		Where("paid = 1 AND recharge_type IN ?", []string{"h5", "weixin", "wechat"}).
		Select("COALESCE(SUM(price),0)").Scan(&totalWx).Error
	_ = h.businessDB.WithContext(ctx).Table(tableName).
		Where("paid = 1").Select("COALESCE(SUM(refund_price),0)").Scan(&totalRefund).Error
	response.OK(c, gin.H{
		"total_pay_price":     round2(totalPay),
		"total_routine_price": round2(totalRoutine),
		"total_wx_price":      round2(totalWx),
		"total_refund_price":  round2(totalRefund),
	})
}

func (h *Handler) Refund(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req struct {
		Amount         *float64 `json:"amount"`
		IdempotencyKey string   `json:"idempotency_key"`
	}
	if c.ShouldBindJSON(&req) != nil {
		response.Fail(c, http.StatusBadRequest, "退款参数错误")
		return
	}
	idem := strings.TrimSpace(req.IdempotencyKey)
	if idem == "" || len(idem) > 128 {
		response.Fail(c, http.StatusBadRequest, "请提供幂等键")
		return
	}

	err := h.businessDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var item struct {
			RechargeID  uint64  `gorm:"column:recharge_id"`
			UID         uint64  `gorm:"column:uid"`
			Price       float64 `gorm:"column:price"`
			GivePrice   float64 `gorm:"column:give_price"`
			Paid        int     `gorm:"column:paid"`
			RefundPrice float64 `gorm:"column:refund_price"`
		}
		if err := tx.Table(tableName).Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("recharge_id,uid,price,give_price,paid,refund_price").
			Where("recharge_id = ?", id).Take(&item).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errNotFound
			}
			return err
		}
		if item.Paid != 1 {
			return errUnpaid
		}
		remain := round2(item.Price - item.RefundPrice)
		if remain <= 0 {
			return errAlreadyRefunded
		}
		amount := remain
		if req.Amount != nil {
			amount = round2(*req.Amount)
			if amount <= 0 || amount > remain+0.000001 {
				return errAmount
			}
		}

		credited := item.Price + item.GivePrice
		claw := 0.0
		if item.Price > 0 {
			claw = round2(credited * (amount / item.Price))
		}

		var previous struct {
			ID uint64 `gorm:"column:id"`
		}
		if err := tx.Table("qixi_crm_b_asset_ledger").
			Select("id").
			Where("asset_type=? AND idempotency_key=?", "balance", idem).
			Take(&previous).Error; err == nil {
			return errIdempotent
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		if err := tx.Exec(
			"INSERT INTO qixi_crm_b_member_account (user_id) VALUES (?) ON DUPLICATE KEY UPDATE user_id=VALUES(user_id)",
			item.UID,
		).Error; err != nil {
			return err
		}
		var account struct {
			Balance float64 `gorm:"column:balance"`
		}
		if err := tx.Table("qixi_crm_b_member_account").Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("balance").Where("user_id=?", item.UID).Take(&account).Error; err != nil {
			return err
		}
		if account.Balance+0.000001 < claw {
			return errBalance
		}
		after := round2(account.Balance - claw)
		if err := tx.Table("qixi_crm_b_member_account").Where("user_id=?", item.UID).
			Update("balance", after).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_asset_ledger").Create(map[string]any{
			"user_id":         item.UID,
			"asset_type":      "balance",
			"amount":          -claw,
			"reference_type":  "user_recharge_refund",
			"reference_id":    strconv.FormatUint(item.RechargeID, 10),
			"idempotency_key": idem,
		}).Error; err != nil {
			return err
		}
		newRefund := round2(item.RefundPrice + amount)
		res := tx.Table(tableName).Where("recharge_id=? AND refund_price=?", id, item.RefundPrice).
			Update("refund_price", newRefund)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errConflict
		}
		return nil
	})
	switch err {
	case nil:
		response.OK(c, gin.H{"ok": true, "recharge_id": id})
	case errIdempotent:
		response.OK(c, gin.H{"ok": true, "recharge_id": id, "idempotent": true})
	case errNotFound:
		response.Fail(c, http.StatusNotFound, "充值记录不存在")
	case errUnpaid:
		response.Fail(c, http.StatusBadRequest, "未支付记录不可退款")
	case errAlreadyRefunded:
		response.Fail(c, http.StatusBadRequest, "该记录已全额退款")
	case errAmount:
		response.Fail(c, http.StatusBadRequest, "退款金额无效")
	case errBalance:
		response.Fail(c, http.StatusBadRequest, "用户余额不足，无法回扣")
	case errConflict:
		response.Fail(c, http.StatusConflict, "充值记录已被其他操作更新，请刷新后重试")
	default:
		fail(c, "充值退款失败")
	}
}

var (
	errNotFound        = fmt.Errorf("not found")
	errUnpaid          = fmt.Errorf("unpaid")
	errAlreadyRefunded = fmt.Errorf("already refunded")
	errAmount          = fmt.Errorf("amount")
	errBalance         = fmt.Errorf("balance")
	errConflict        = fmt.Errorf("conflict")
	errIdempotent      = fmt.Errorf("idempotent")
)

func (h *Handler) baseQuery(c *gin.Context, f listFilter) *gorm.DB {
	q := h.businessDB.WithContext(c.Request.Context()).Table(tableName+" AS r").
		Joins("LEFT JOIN qixi_crm_b_user u ON u.id = r.uid").
		Joins("LEFT JOIN qixi_crm_b_user_profile p ON p.user_id = r.uid")
	if f.DateFrom != "" {
		q = q.Where("COALESCE(r.pay_time, r.create_time) >= ?", f.DateFrom+" 00:00:00")
	}
	if f.DateTo != "" {
		q = q.Where("COALESCE(r.pay_time, r.create_time) <= ?", f.DateTo+" 23:59:59")
	}
	if f.Paid != nil {
		q = q.Where("r.paid = ?", *f.Paid)
	}
	switch f.RechargeType {
	case "1":
		q = q.Where("r.recharge_type IN ?", []string{"h5", "weixin", "wechat", "routine"})
	case "2":
		q = q.Where("r.recharge_type IN ?", []string{"alipay"})
	case "":
	default:
		q = q.Where("r.recharge_type = ?", f.RechargeType)
	}
	if kw := strings.TrimSpace(f.OrderID); kw != "" {
		q = q.Where("r.order_id LIKE ?", "%"+kw+"%")
	}
	if kw := strings.TrimSpace(f.UserKeyword); kw != "" {
		switch f.UserType {
		case "uid":
			if uid, err := strconv.ParseUint(kw, 10, 64); err == nil && uid > 0 {
				q = q.Where("r.uid = ?", uid)
			} else {
				q = q.Where("1 = 0")
			}
		case "phone":
			q = q.Where("u.mobile LIKE ?", "%"+kw+"%")
		case "real_name":
			q = q.Where("p.real_name LIKE ?", "%"+kw+"%")
		default:
			q = q.Where("u.nickname LIKE ?", "%"+kw+"%")
		}
	}
	return q
}

func parseFilter(c *gin.Context) listFilter {
	f := listFilter{
		DateFrom:     strings.TrimSpace(c.Query("date_from")),
		DateTo:       strings.TrimSpace(c.Query("date_to")),
		RechargeType: strings.TrimSpace(c.Query("recharge_type")),
		UserType:     strings.TrimSpace(c.Query("user_type")),
		UserKeyword:  strings.TrimSpace(c.Query("user_keyword")),
		OrderID:      strings.TrimSpace(c.Query("order_id")),
	}
	if raw := strings.TrimSpace(c.Query("paid")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && (v == 0 || v == 1) {
			f.Paid = &v
		}
	}
	if f.UserType == "" {
		f.UserType = "nickname"
	}
	return f
}

func parseID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "充值记录参数错误")
		return 0, false
	}
	return id, true
}

func decorate(r *row) {
	r.RechargeTypeName = typeName(r.RechargeType)
	if r.Paid == 1 {
		r.PaidName = "已支付"
	} else {
		r.PaidName = "未支付"
	}
	r.CanRefund = r.Paid == 1 && round2(r.Price-r.RefundPrice) > 0
}

func typeName(t string) string {
	switch t {
	case "routine":
		return "小程序"
	case "h5", "weixin", "wechat":
		return "微信"
	case "alipay":
		return "支付宝"
	default:
		if t == "" {
			return "-"
		}
		return t
	}
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func fail(c *gin.Context, msg string) {
	response.Fail(c, http.StatusInternalServerError, msg)
}
