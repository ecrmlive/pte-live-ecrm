package distribution

import (
	"net/http"
	"strconv"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler serves the C-end distribution projection from business-owned tables.
type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/spread/me", h.me)
	r.GET("/spread/rank", h.rank)
	r.GET("/spread/users", h.users)
	r.GET("/spread/orders", h.orders)
	r.POST("/spread/bind", h.bind)
	r.GET("/spread/bills", h.bills)
	r.GET("/spread/withdrawals", h.withdrawals)
	r.POST("/spread/withdrawals", h.applyWithdrawal)
}

type bindInput struct {
	SpreadUID uint64 `json:"spread_uid"`
}

type commissionRow struct {
	ID        uint64    `gorm:"column:id"`
	Amount    float64   `gorm:"column:amount"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (h *Handler) me(c *gin.Context) {
	uid := uint64(middleware.UID(c))
	var parent uint64
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_relation").Select("COALESCE(parent_user_id,0)").Where("user_id = ?", uid).Scan(&parent).Error; err != nil {
		fail(c)
		return
	}
	var promoterCount, childCount int64
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_promoter").Where("user_id = ? AND status = 1", uid).Count(&promoterCount).Error; err != nil {
		fail(c)
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_relation").Where("parent_user_id = ?", uid).Count(&childCount).Error; err != nil {
		fail(c)
		return
	}
	response.OK(c, gin.H{"uid": uid, "spread_uid": parent, "is_promoter": promoterCount == 1, "spread_count": childCount})
}

func (h *Handler) bind(c *gin.Context) {
	var input bindInput
	if err := c.ShouldBindJSON(&input); err != nil || input.SpreadUID == 0 {
		response.Fail(c, http.StatusBadRequest, "推广员 UID 错误")
		return
	}
	uid := uint64(middleware.UID(c))
	if uid == input.SpreadUID {
		response.Fail(c, http.StatusBadRequest, "不能绑定自己为推广员")
		return
	}
	var promoters int64
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_promoter").Where("user_id = ? AND status = 1", input.SpreadUID).Count(&promoters).Error; err != nil {
		fail(c)
		return
	}
	if promoters != 1 {
		response.Fail(c, http.StatusBadRequest, "推广员不存在或未启用")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("INSERT INTO qixi_crm_b_distribution_relation (user_id,parent_user_id,bound_at) VALUES (?,?,?) ON DUPLICATE KEY UPDATE parent_user_id=IF(parent_user_id IS NULL,VALUES(parent_user_id),parent_user_id),bound_at=IF(parent_user_id IS NULL,VALUES(bound_at),bound_at)", uid, input.SpreadUID, time.Now()).Error; err != nil {
			return err
		}
		var parent uint64
		if err := tx.Table("qixi_crm_b_distribution_relation").Select("COALESCE(parent_user_id,0)").Where("user_id = ?", uid).Scan(&parent).Error; err != nil {
			return err
		}
		if parent != input.SpreadUID {
			return errAlreadyBound
		}
		return nil
	}); err != nil {
		if err == errAlreadyBound {
			response.Fail(c, http.StatusBadRequest, "已绑定其他推广员，不能更换")
			return
		}
		fail(c)
		return
	}
	response.OK(c, gin.H{"spread_uid": input.SpreadUID, "bound": true})
}

func (h *Handler) bills(c *gin.Context) {
	page, limit := pageParams(c)
	uid := uint64(middleware.UID(c))
	var total int64
	q := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_commission_ledger").Where("user_id = ? AND status <> ?", uid, "voided")
	if err := q.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	rows := make([]commissionRow, 0)
	if err := q.Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		fail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, gin.H{"bill_id": row.ID, "title": "推广佣金", "number": row.Amount, "balance": 0, "mark": commissionMark(row.Status), "create_time": row.CreatedAt.Format("2006-01-02 15:04:05"), "status": row.Status})
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func commissionMark(status string) string {
	switch status {
	case "pending":
		return "待结算"
	case "available":
		return "可结算"
	case "settled":
		return "已结算"
	default:
		return "佣金记录"
	}
}

func pageParams(c *gin.Context) (int, int) {
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

func fail(c *gin.Context) { response.Fail(c, http.StatusInternalServerError, "分销服务异常") }

var errAlreadyBound = &bindingError{}

type bindingError struct{}

func (*bindingError) Error() string { return "already bound" }
