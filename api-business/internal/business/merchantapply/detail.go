package merchantapply

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type applicationDetail struct {
	ID              uint64 `json:"id"`
	MerchantName    string `json:"merchant_name"`
	ContactName     string `json:"contact_name"`
	ContactMobile   string `json:"contact_mobile"`
	CategoryName    string `json:"category_name"`
	MerchantType    string `json:"merchant_type"`
	Status          string `json:"status"`
	ReviewNote      string `json:"review_note"`
	LicenseUploaded bool   `json:"license_uploaded"`
	CreatedAt       string `json:"created_at"`
}

// Detail returns only the current user's own application and deliberately
// omits object-storage keys and URLs. A license is represented by a boolean;
// download authorization remains the responsibility of the upload service.
func (h *Handler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "申请编号错误")
		return
	}
	var row applicationRow
	err = h.db.WithContext(c.Request.Context()).Where("id = ? AND applicant_user_id = ?", id, middleware.UID(c)).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "申请记录不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询申请详情失败")
		return
	}
	response.OK(c, applicationDetail{
		ID: row.ID, MerchantName: row.MerchantName, ContactName: row.ContactName, ContactMobile: row.ContactMobile,
		CategoryName: row.CategoryName, MerchantType: row.MerchantType, Status: row.Status, ReviewNote: row.ReviewNote,
		LicenseUploaded: row.LicenseKey != "", CreatedAt: row.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}
