package distribution

import (
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type commissionOrderRow struct {
	ID        uint64  `gorm:"column:id"`
	OrderID   uint64  `gorm:"column:order_id"`
	Amount    float64 `gorm:"column:amount"`
	Status    string  `gorm:"column:status"`
	CreatedAt string  `gorm:"column:created_at"`
}

func spreadOrderKeyword(raw string) (uint64, bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return 0, false, nil
	}
	if len(value) > 20 {
		return 0, false, errInvalidSpreadOrderKeyword
	}
	for _, char := range value {
		if !unicode.IsDigit(char) {
			return 0, false, errInvalidSpreadOrderKeyword
		}
	}
	orderID, err := strconv.ParseUint(value, 10, 64)
	if err != nil || orderID == 0 {
		return 0, false, errInvalidSpreadOrderKeyword
	}
	return orderID, true, nil
}

func commissionOrderStatus(status string) string {
	switch status {
	case "pending":
		return "待结算"
	case "available":
		return "可提现"
	case "settled":
		return "已结算"
	default:
		return "佣金记录"
	}
}

func (h *Handler) orders(c *gin.Context) {
	orderID, filtered, err := spreadOrderKeyword(c.Query("keyword"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "订单号格式错误")
		return
	}
	page, limit := pageParams(c)
	uid := uint64(middleware.UID(c))
	q := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_commission_ledger").Where("user_id = ? AND status <> ?", uid, "voided")
	if filtered {
		q = q.Where("order_id = ?", orderID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	var totalAmount float64
	if err := q.Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount).Error; err != nil {
		fail(c)
		return
	}
	rows := make([]commissionOrderRow, 0)
	if err := q.Select("id, order_id, amount, status, DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s') AS created_at").Order("id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, gin.H{"record_id": row.ID, "order_id": row.OrderID, "commission_amount": row.Amount, "status": row.Status, "status_label": commissionOrderStatus(row.Status), "created_at": row.CreatedAt})
	}
	response.OK(c, gin.H{"list": list, "total": total, "sum_commission": totalAmount, "page": page, "limit": limit})
}

var errInvalidSpreadOrderKeyword = &spreadOrderKeywordError{}

type spreadOrderKeywordError struct{}

func (*spreadOrderKeywordError) Error() string { return "invalid spread order keyword" }
