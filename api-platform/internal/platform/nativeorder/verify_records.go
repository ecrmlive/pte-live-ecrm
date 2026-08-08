package nativeorder

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// verifyRecordRow is flat on purpose. GORM Scan does not populate anonymous /
// gorm:"embedded" nested structs here — only top-level fields get values.
// That bug zeroed order_sn/mer_id/pay while verify_* still worked.
type verifyRecordRow struct {
	ID                  uint64          `gorm:"column:id"`
	GroupID             uint64          `gorm:"column:group_order_id"`
	GroupOrderNo        string          `gorm:"column:group_order_no"`
	OrderNo             string          `gorm:"column:order_no"`
	MerchantID          uint64          `gorm:"column:merchant_id"`
	StoreID             uint64          `gorm:"column:store_id"`
	UserID              uint64          `gorm:"column:user_id"`
	PayAmount           float64         `gorm:"column:pay_amount"`
	TotalAmount         float64         `gorm:"column:total_amount"`
	DiscountAmount      float64         `gorm:"column:discount_amount"`
	FreightAmount       float64         `gorm:"column:freight_amount"`
	PointsAmount        int64           `gorm:"column:points_amount"`
	TotalQuantity       int             `gorm:"column:total_quantity"`
	ActivityType        int             `gorm:"column:activity_type"`
	Recipient           json.RawMessage `gorm:"column:recipient_snapshot"`
	Remark              string          `gorm:"column:remark"`
	Status              string          `gorm:"column:status"`
	CreatedAt           time.Time       `gorm:"column:created_at"`
	PaidAt              *time.Time      `gorm:"column:paid_at"`
	PayChannel          string          `gorm:"column:pay_channel"`
	UserArchivedAt      *time.Time      `gorm:"column:user_archived_at"`
	HasRefunded         int             `gorm:"column:has_refunded"`
	PendingComment      int             `gorm:"column:pending_comment"`
	VerifiedAt          *time.Time      `gorm:"column:verified_at"`
	VerifiedByAccountID uint64          `gorm:"column:verified_by_account_id"`
	VerifyStatus        string          `gorm:"column:verify_status"`
}

func (r verifyRecordRow) asOrder() order {
	return order{
		ID: r.ID, GroupID: r.GroupID, GroupOrderNo: r.GroupOrderNo, OrderNo: r.OrderNo,
		MerchantID: r.MerchantID, StoreID: r.StoreID, UserID: r.UserID,
		PayAmount: r.PayAmount, TotalAmount: r.TotalAmount, DiscountAmount: r.DiscountAmount,
		FreightAmount: r.FreightAmount, PointsAmount: r.PointsAmount, TotalQuantity: r.TotalQuantity,
		ActivityType: r.ActivityType, Recipient: r.Recipient, Remark: r.Remark, Status: r.Status,
		CreatedAt: r.CreatedAt, PaidAt: r.PaidAt, PayChannel: r.PayChannel,
		UserArchivedAt: r.UserArchivedAt, HasRefunded: r.HasRefunded, PendingComment: r.PendingComment,
	}
}

type accountNameRow struct {
	ID          uint64 `gorm:"column:id"`
	DisplayName string `gorm:"column:display_name"`
	Username    string `gorm:"column:username"`
	RoleCode    string `gorm:"column:role_code"`
}

// verifyList returns platform supervision rows for already-verified (write-off) orders.
// Read-only: no verify/cancel mutations.
func (h *Handler) verifyList(c *gin.Context) {
	page, limit := normalizePage(c)
	merchantIDs, err := h.merchantScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置订单监管数据范围")
		return
	}
	if merchantIDs != nil && len(merchantIDs) == 0 {
		response.OK(c, gin.H{"list": []gin.H{}, "total": 0, "page": page, "limit": limit})
		return
	}

	countQ := applyVerifyRecordFilters(h.verifyScopeQuery(c, merchantIDs), c, h.adminDB)
	var total int64
	if err := countQ.Distinct("o.id").Count(&total).Error; err != nil {
		fail(c, "查询核销记录失败")
		return
	}

	q := applyVerifyRecordFilters(h.verifyBase(c, merchantIDs), c, h.adminDB)
	var rows []verifyRecordRow
	if err := q.Order("v.verified_at DESC, o.id DESC").
		Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c, "查询核销记录失败")
		return
	}
	items, err := h.verifyResponses(c, rows)
	if err != nil {
		fail(c, "加载核销记录失败")
		return
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}

// verifySummary aggregates paid counts / amounts for the same filter set as verifyList.
func (h *Handler) verifySummary(c *gin.Context) {
	merchantIDs, err := h.merchantScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置订单监管数据范围")
		return
	}
	empty := gin.H{
		"paid_count": 0, "pay_amount": 0.0, "refund_amount": 0.0,
		"wechat_amount": 0.0, "balance_amount": 0.0, "alipay_amount": 0.0,
	}
	if merchantIDs != nil && len(merchantIDs) == 0 {
		response.OK(c, empty)
		return
	}

	var paidCount int64
	if err := applyVerifyRecordFilters(h.verifyScopeQuery(c, merchantIDs), c, h.adminDB).
		Distinct("o.id").Count(&paidCount).Error; err != nil {
		fail(c, "统计核销记录失败")
		return
	}

	type sumRow struct {
		PayAmount     float64 `gorm:"column:pay_amount"`
		WechatAmount  float64 `gorm:"column:wechat_amount"`
		BalanceAmount float64 `gorm:"column:balance_amount"`
		AlipayAmount  float64 `gorm:"column:alipay_amount"`
	}
	var sums sumRow
	if err := applyVerifyRecordFilters(h.verifyScopeQuery(c, merchantIDs), c, h.adminDB).Select(`
		COALESCE(SUM(o.pay_amount),0) AS pay_amount,
		COALESCE(SUM(CASE WHEN g.pay_channel='wechat' THEN o.pay_amount ELSE 0 END),0) AS wechat_amount,
		COALESCE(SUM(CASE WHEN g.pay_channel='balance' THEN o.pay_amount ELSE 0 END),0) AS balance_amount,
		COALESCE(SUM(CASE WHEN g.pay_channel='alipay' THEN o.pay_amount ELSE 0 END),0) AS alipay_amount
	`).Scan(&sums).Error; err != nil {
		fail(c, "统计核销记录失败")
		return
	}

	var refundAmount float64
	if err := applyVerifyRecordFilters(h.verifyScopeQuery(c, merchantIDs), c, h.adminDB).
		Joins("JOIN qixi_crm_b_refund r ON r.order_id = o.id AND r.status = 'refunded'").
		Select("COALESCE(SUM(r.amount),0)").Scan(&refundAmount).Error; err != nil {
		fail(c, "统计核销退款失败")
		return
	}

	response.OK(c, gin.H{
		"paid_count":     paidCount,
		"pay_amount":     sums.PayAmount,
		"refund_amount":  refundAmount,
		"wechat_amount":  sums.WechatAmount,
		"balance_amount": sums.BalanceAmount,
		"alipay_amount":  sums.AlipayAmount,
	})
}

func (h *Handler) verifyScopeQuery(c *gin.Context, merchantIDs []uint64) *gorm.DB {
	// One used verification per order (latest id) so list/sum never double-count.
	q := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order AS o").
		Joins("JOIN qixi_crm_b_group_order AS g ON g.id = o.group_order_id").
		Joins(`JOIN (
			SELECT order_id, MAX(id) AS id
			FROM qixi_crm_b_order_verification
			WHERE status = 'used'
			GROUP BY order_id
		) AS vu ON vu.order_id = o.id`).
		Joins("JOIN qixi_crm_b_order_verification AS v ON v.id = vu.id")
	if merchantIDs != nil {
		q = q.Where("o.merchant_id IN ?", merchantIDs)
	}
	return q
}

func (h *Handler) verifyBase(c *gin.Context, merchantIDs []uint64) *gorm.DB {
	return h.verifyScopeQuery(c, merchantIDs).Select(`o.id,o.group_order_id,g.order_no AS group_order_no,o.order_no,o.merchant_id,o.store_id,o.user_id,
			o.pay_amount,o.total_amount,o.discount_amount,o.freight_amount,o.points_amount,o.total_quantity,o.activity_type,
			o.recipient_snapshot,o.remark,o.status,o.created_at,o.paid_at,g.pay_channel,g.user_archived_at,
			CASE WHEN EXISTS(SELECT 1 FROM qixi_crm_b_refund r WHERE r.order_id=o.id AND r.status='refunded') THEN 1 ELSE 0 END AS has_refunded,
			0 AS pending_comment,
			v.verified_at, COALESCE(v.verified_by_account_id,0) AS verified_by_account_id, v.status AS verify_status`)
}

func applyVerifyRecordFilters(q *gorm.DB, c *gin.Context, adminDB *gorm.DB) *gorm.DB {
	if payType := strings.TrimSpace(c.Query("pay_type")); payType != "" {
		if channel := payChannelFromType(payType); channel != "" {
			q = q.Where("g.pay_channel = ?", channel)
		}
	}

	if isTrader := strings.TrimSpace(c.Query("is_trader")); isTrader != "" && adminDB != nil {
		if trader, err := strconv.Atoi(isTrader); err == nil && (trader == 0 || trader == 1) {
			var ids []uint64
			_ = adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").
				Select("merchant_id").Where("is_trader = ?", trader).Pluck("merchant_id", &ids)
			if len(ids) == 0 {
				q = q.Where("1 = 0")
			} else {
				q = q.Where("o.merchant_id IN ?", ids)
			}
		}
	}

	orderKW := strings.TrimSpace(c.Query("order_keyword"))
	if orderKW == "" {
		orderKW = strings.TrimSpace(c.Query("order_sn"))
	}
	if orderKW != "" {
		like := "%" + orderKW + "%"
		q = q.Where(
			`o.order_no LIKE ? OR JSON_UNQUOTE(JSON_EXTRACT(o.recipient_snapshot,'$.recipient')) LIKE ?
			 OR JSON_UNQUOTE(JSON_EXTRACT(o.recipient_snapshot,'$.mobile')) LIKE ?`,
			like, like, like,
		)
	}

	userKW := strings.TrimSpace(c.Query("user_keyword"))
	if userKW == "" {
		userKW = strings.TrimSpace(c.Query("username"))
	}
	if userKW != "" {
		like := "%" + userKW + "%"
		q = q.Where(`EXISTS (
			SELECT 1 FROM qixi_crm_b_user u
			WHERE u.id = o.user_id AND (u.nickname LIKE ? OR u.mobile LIKE ?)
		)`, like, like)
	}

	if from := strings.TrimSpace(c.Query("date_from")); from != "" {
		if t, err := time.ParseInLocation("2006-01-02", from, time.Local); err == nil {
			q = q.Where("v.verified_at >= ?", t)
		}
	}
	if to := strings.TrimSpace(c.Query("date_to")); to != "" {
		if t, err := time.ParseInLocation("2006-01-02", to, time.Local); err == nil {
			q = q.Where("v.verified_at < ?", t.AddDate(0, 0, 1))
		}
	}
	return q
}

func (h *Handler) verifyResponses(c *gin.Context, rows []verifyRecordRow) ([]gin.H, error) {
	if len(rows) == 0 {
		return []gin.H{}, nil
	}
	baseOrders := make([]order, 0, len(rows))
	accountIDs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		baseOrders = append(baseOrders, row.asOrder())
		if row.VerifiedByAccountID > 0 {
			accountIDs = append(accountIDs, row.VerifiedByAccountID)
		}
	}
	items, err := h.responses(c, baseOrders, true)
	if err != nil {
		return nil, err
	}

	accounts := map[uint64]accountNameRow{}
	if len(accountIDs) > 0 {
		var accountRows []accountNameRow
		_ = h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_account").
			Select("id,display_name,username,role_code").Where("id IN ?", accountIDs).Find(&accountRows)
		for _, a := range accountRows {
			accounts[a.ID] = a
		}
	}

	for i := range items {
		row := rows[i]
		items[i]["order_type_label"] = "核销订单"
		items[i]["verify_status"] = "used"
		items[i]["verify_status_label"] = "已核销"
		items[i]["verify_time"] = ""
		if row.VerifiedAt != nil {
			items[i]["verify_time"] = row.VerifiedAt.Format("2006-01-02 15:04:05")
		}
		items[i]["verifier_account_id"] = row.VerifiedByAccountID
		items[i]["verifier_name"] = verifierDisplayName(row.VerifiedByAccountID, accounts[row.VerifiedByAccountID])
		items[i]["pay_type_label"] = verifyPayTypeLabel(row.PayChannel, paid(row.Status))
	}
	return items, nil
}

func verifierDisplayName(accountID uint64, account accountNameRow) string {
	if accountID == 0 {
		return "管理员核销"
	}
	if name := strings.TrimSpace(account.DisplayName); name != "" {
		return name
	}
	if name := strings.TrimSpace(account.Username); name != "" {
		return name
	}
	if account.RoleCode == "owner" || account.RoleCode == "manager" {
		return "管理员核销"
	}
	return strconv.FormatUint(accountID, 10)
}

func verifyPayTypeLabel(channel string, paidFlag int) string {
	if paidFlag == 0 || channel == "" {
		return "未支付"
	}
	switch channel {
	case "wechat":
		return "微信支付"
	case "alipay":
		return "支付宝支付"
	case "balance":
		return "余额支付"
	case "mock":
		return "模拟支付"
	default:
		return channel
	}
}
