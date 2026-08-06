package invoice

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/queryfilter"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ business, admin *gorm.DB }

func New(business, admin *gorm.DB) *Handler { return &Handler{business: business, admin: admin} }

func (h *Handler) Register(r gin.IRoutes) {
	read := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "accounts.invoice.read")}
	r.GET("/finance/invoices", append(read, h.List)...)
	r.GET("/finance/invoices/:id", append(read, h.Get)...)
}

type invoiceRow struct {
	ID, OrderID, MerchantID, StoreID                                     uint64
	OrderNo, MerchantName, StoreName                                     string
	ProfileType, Title, TaxNo, Email, Status, InvoiceNo, RejectionReason string
	RequestedAt                                                          time.Time
	IssuedAt                                                             *time.Time
}

func (h *Handler) List(c *gin.Context) {
	page, limit := pagination(c)
	q := h.base(c)
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		if !validStatus(status) {
			response.Fail(c, http.StatusBadRequest, "发票状态错误")
			return
		}
		q = q.Where("oi.status=?", status)
	}
	if orderNo := strings.TrimSpace(c.Query("order_no")); orderNo != "" {
		if len(orderNo) > 64 {
			response.Fail(c, http.StatusBadRequest, "订单号参数错误")
			return
		}
		q = q.Where("o.order_no LIKE ?", "%"+orderNo+"%")
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("o.order_no LIKE ? OR oi.title LIKE ? OR oi.invoice_no LIKE ?", like, like, like)
	}
	q = queryfilter.ApplyCreatedAtRange(q, c, "oi.requested_at")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "发票列表查询失败")
		return
	}
	rows := make([]invoiceRow, 0)
	if err := q.Order("oi.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "发票列表查询失败")
		return
	}
	response.OK(c, gin.H{"list": views(rows), "total": total, "page": page, "limit": limit})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "发票 ID 参数错误")
		return
	}
	var row invoiceRow
	if err := h.base(c).Where("oi.id=?", id).Take(&row).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "发票记录不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "发票详情查询失败")
		return
	}
	response.OK(c, views([]invoiceRow{row})[0])
}

func (h *Handler) base(c *gin.Context) *gorm.DB {
	return h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_order_invoice AS oi").
		Select("oi.id,oi.order_id,o.order_no,o.merchant_id,o.merchant_name_snapshot AS merchant_name,o.store_id,o.store_name_snapshot AS store_name,oi.profile_type,oi.title,oi.tax_no,oi.email,oi.status,oi.invoice_no,oi.rejection_reason,oi.requested_at,oi.issued_at").
		Joins("JOIN qixi_crm_b_order AS o ON o.id=oi.order_id")
}
func pagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}
func validStatus(v string) bool {
	return v == "requested" || v == "issued" || v == "rejected" || v == "voided"
}
func views(rows []invoiceRow) []gin.H {
	out := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		out = append(out, gin.H{"id": r.ID, "order_id": r.OrderID, "order_no": r.OrderNo, "merchant_id": r.MerchantID, "merchant_name": r.MerchantName, "store_id": r.StoreID, "store_name": r.StoreName, "profile_type": r.ProfileType, "title": r.Title, "tax_no_masked": maskTax(r.TaxNo), "email_masked": maskEmail(r.Email), "status": r.Status, "invoice_no": r.InvoiceNo, "rejection_reason": r.RejectionReason, "requested_at": r.RequestedAt, "issued_at": r.IssuedAt})
	}
	return out
}
func maskTax(v string) string {
	v = strings.TrimSpace(v)
	if len(v) <= 4 {
		return ""
	}
	return strings.Repeat("*", len(v)-4) + v[len(v)-4:]
}
func maskEmail(v string) string {
	v = strings.TrimSpace(v)
	at := strings.Index(v, "@")
	if at <= 0 {
		return ""
	}
	return string([]rune(v)[:1]) + "***" + v[at:]
}
