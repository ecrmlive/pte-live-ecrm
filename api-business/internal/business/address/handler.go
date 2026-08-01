// Package address owns user address data in qixi_crm_business.
package address

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/response"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/address", h.List)
	r.POST("/address", h.Create)
	r.PUT("/address/:id", h.Update)
	r.DELETE("/address/:id", h.Delete)
}

type row struct {
	ID         uint64 `gorm:"column:id"`
	UserID     uint64 `gorm:"column:user_id"`
	Recipient  string `gorm:"column:recipient"`
	Mobile     string `gorm:"column:mobile"`
	Province   string `gorm:"column:province"`
	City       string `gorm:"column:city"`
	District   string `gorm:"column:district"`
	RegionCode string `gorm:"column:region_code"`
	Detail     string `gorm:"column:detail"`
	PostCode   int    `gorm:"column:post_code"`
	IsDefault  int8   `gorm:"column:is_default"`
}

func (row) TableName() string { return "qixi_crm_b_user_address" }

type request struct {
	RealName   *string `json:"real_name"`
	Phone      *string `json:"phone"`
	Province   *string `json:"province"`
	City       *string `json:"city"`
	District   *string `json:"district"`
	RegionCode *string `json:"region_code"`
	Detail     *string `json:"detail"`
	PostCode   *int    `json:"post_code"`
	IsDefault  *int8   `json:"is_default"`
}

func (h *Handler) List(c *gin.Context) {
	var rows []row
	err := h.db.WithContext(c.Request.Context()).Where("user_id = ?", middleware.UID(c)).
		Order("is_default DESC, id DESC").Find(&rows).Error
	if err != nil {
		serverError(c)
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, item := range rows {
		items = append(items, toResponse(item))
	}
	response.OK(c, gin.H{"list": items})
}

func (h *Handler) Create(c *gin.Context) {
	var input request
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	address := row{UserID: uint64(middleware.UID(c))}
	if err := merge(&address, input, true); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if address.IsDefault == 1 {
			if err := tx.Model(&row{}).Where("user_id = ?", address.UserID).Update("is_default", 0).Error; err != nil {
				return err
			}
		}
		return tx.Create(&address).Error
	}); err != nil {
		serverError(c)
		return
	}
	response.OK(c, toResponse(address))
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "地址 ID 错误")
		return
	}
	var input request
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	uid := uint64(middleware.UID(c))
	var address row
	err = h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND user_id = ?", id, uid).First(&address).Error; err != nil {
			return err
		}
		if err := merge(&address, input, false); err != nil {
			return err
		}
		if address.IsDefault == 1 {
			if err := tx.Model(&row{}).Where("user_id = ? AND id <> ?", uid, id).Update("is_default", 0).Error; err != nil {
				return err
			}
		}
		return tx.Save(&address).Error
	})
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, toResponse(address))
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "地址 ID 错误")
		return
	}
	result := h.db.WithContext(c.Request.Context()).Where("id = ? AND user_id = ?", id, middleware.UID(c)).Delete(&row{})
	if result.Error != nil {
		serverError(c)
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "地址不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func merge(address *row, input request, creating bool) error {
	if input.RealName != nil {
		address.Recipient = strings.TrimSpace(*input.RealName)
	}
	if input.Phone != nil {
		address.Mobile = strings.TrimSpace(*input.Phone)
	}
	if input.Province != nil {
		address.Province = strings.TrimSpace(*input.Province)
	}
	if input.City != nil {
		address.City = strings.TrimSpace(*input.City)
	}
	if input.District != nil {
		address.District = strings.TrimSpace(*input.District)
	}
	if input.RegionCode != nil {
		address.RegionCode = strings.TrimSpace(*input.RegionCode)
	}
	if input.Detail != nil {
		address.Detail = strings.TrimSpace(*input.Detail)
	}
	if input.PostCode != nil {
		address.PostCode = *input.PostCode
	}
	if input.IsDefault != nil {
		address.IsDefault = *input.IsDefault
	}
	if address.Recipient == "" || address.Mobile == "" || address.Detail == "" {
		return errors.New("收货人、手机号和详细地址不能为空")
	}
	if address.IsDefault != 0 && address.IsDefault != 1 {
		return errors.New("默认地址参数错误")
	}
	if creating && (input.RealName == nil || input.Phone == nil || input.Detail == nil) {
		return errors.New("收货人、手机号和详细地址不能为空")
	}
	return nil
}

func toResponse(address row) gin.H {
	return gin.H{
		"address_id": address.ID, "real_name": address.Recipient, "phone": address.Mobile,
		"province": address.Province, "city": address.City, "district": address.District,
		"region_code": address.RegionCode, "detail": address.Detail, "post_code": address.PostCode,
		"is_default": address.IsDefault,
	}
}

func writeError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "地址不存在")
		return
	}
	if strings.Contains(err.Error(), "不能为空") || strings.Contains(err.Error(), "参数错误") {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	serverError(c)
}

func serverError(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "地址服务异常")
}
