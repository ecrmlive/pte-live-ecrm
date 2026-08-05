package integralpolicy

import (
	"net/http"

	merchantintegralpolicyevent "github.com/crmlive/pte-live-ecrm/api-merchant/internal/event/merchantintegralpolicy"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }
func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/setting/integral-policy", h.get)
	r.PUT("/setting/integral-policy", middleware.RequireStorePermission(h.db, "setting.shop"), h.put)
}

type policy struct {
	StoreID          uint64 `gorm:"column:store_id" json:"store_id"`
	Enabled          bool   `gorm:"column:enabled" json:"enabled"`
	PointsPerYuan    int64  `gorm:"column:points_per_yuan" json:"points_per_yuan"`
	MaxDeductionBps  int64  `gorm:"column:max_deduction_bps" json:"max_deduction_bps"`
	UpdatedByAccount uint64 `gorm:"column:updated_by_account_id" json:"-"`
}

func (policy) TableName() string { return "qixi_crm_m_store_integral_policy" }

func (h *Handler) get(c *gin.Context) {
	row := policy{StoreID: uint64(middleware.StoreID(c)), PointsPerYuan: 100, MaxDeductionBps: 2000}
	if err := h.db.WithContext(c.Request.Context()).Where("store_id = ?", row.StoreID).First(&row).Error; err != nil && err != gorm.ErrRecordNotFound {
		response.Fail(c, http.StatusInternalServerError, "读取积分抵扣设置失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) put(c *gin.Context) {
	var in policy
	if err := c.ShouldBindJSON(&in); err != nil || in.PointsPerYuan < 1 || in.PointsPerYuan > 100000 || in.MaxDeductionBps < 1 || in.MaxDeductionBps > 10000 {
		response.Fail(c, http.StatusBadRequest, "积分抵扣设置参数不合法")
		return
	}
	in.StoreID = uint64(middleware.StoreID(c))
	in.UpdatedByAccount = uint64(middleware.AdminID(c))
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("store_id = ?", in.StoreID).Assign(map[string]any{"enabled": in.Enabled, "points_per_yuan": in.PointsPerYuan, "max_deduction_bps": in.MaxDeductionBps, "updated_by_account_id": in.UpdatedByAccount}).FirstOrCreate(&policy{StoreID: in.StoreID}).Error; err != nil {
			return err
		}
		return merchantintegralpolicyevent.Enqueue(tx, merchantintegralpolicyevent.Payload{StoreID: in.StoreID, Enabled: in.Enabled, PointsPerYuan: in.PointsPerYuan, MaxDeductionBps: in.MaxDeductionBps})
	}); err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存积分抵扣设置失败")
		return
	}
	response.OK(c, in)
}
