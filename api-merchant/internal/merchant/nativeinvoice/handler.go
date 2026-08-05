// Package nativeinvoice serves store invoice audit from qixi_crm_b_order_invoice.
package nativeinvoice

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	db, merchantDB *gorm.DB
}

func NewHandler(db, merchantDB *gorm.DB) *Handler {
	return &Handler{db: db, merchantDB: merchantDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/invoices", h.list)
	r.PUT("/invoices/:id/audit", middleware.RequireStorePermission(h.merchantDB, "invoice.audit"), h.audit)
}

func (h *Handler) list(c *gin.Context) {
	page, limit := pageLimit(c)
	q := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_order_invoice AS i").
		Select("i.id,i.order_id,o.user_id,i.title,i.tax_no,i.email,i.status,i.rejection_reason,i.invoice_no,i.requested_at").
		Joins("JOIN qixi_crm_b_order AS o ON o.id = i.order_id").
		Where("o.store_id = ?", middleware.StoreID(c))
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		if mapped := fromUIStatus(status); mapped != "" {
			q = q.Where("i.status = ?", mapped)
		} else if n, err := strconv.Atoi(status); err == nil {
			if mapped := fromLegacyStatus(n); mapped != "" {
				q = q.Where("i.status = ?", mapped)
			}
		}
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c, "查询发票失败")
		return
	}
	var rows []struct {
		ID               uint64    `gorm:"column:id"`
		OrderID          uint64    `gorm:"column:order_id"`
		UserID           uint64    `gorm:"column:user_id"`
		Title            string    `gorm:"column:title"`
		TaxNo            string    `gorm:"column:tax_no"`
		Email            string    `gorm:"column:email"`
		Status           string    `gorm:"column:status"`
		RejectionReason  string    `gorm:"column:rejection_reason"`
		InvoiceNo        string    `gorm:"column:invoice_no"`
		RequestedAt      time.Time `gorm:"column:requested_at"`
	}
	if err := q.Order("i.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c, "查询发票失败")
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		mark := row.RejectionReason
		if row.Status == "issued" && row.InvoiceNo != "" {
			mark = row.InvoiceNo
		}
		list = append(list, gin.H{
			"invoice_id": row.ID,
			"order_id":   row.OrderID,
			"uid":        row.UserID,
			"header":     row.Title,
			"tax_no":     row.TaxNo,
			"email":      row.Email,
			"status":     toLegacyStatus(row.Status),
			"status_code": row.Status,
			"mark":       mark,
			"create_time": row.RequestedAt.Format("2006-01-02 15:04:05"),
		})
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) audit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Status int    `json:"status"`
		Mark   string `json:"mark"`
	}
	if id == 0 || c.ShouldBindJSON(&req) != nil || (req.Status != 1 && req.Status != -1) {
		response.Fail(c, http.StatusBadRequest, "审核参数不合法")
		return
	}
	req.Mark = strings.TrimSpace(req.Mark)
	if utf8.RuneCountInString(req.Mark) > 200 {
		response.Fail(c, http.StatusBadRequest, "备注不能超过 200 个字符")
		return
	}
	storeID := uint64(middleware.StoreID(c))
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var row struct {
			ID     uint64 `gorm:"column:id"`
			Status string `gorm:"column:status"`
		}
		if err := tx.Table("qixi_crm_b_order_invoice AS i").
			Select("i.id,i.status").
			Joins("JOIN qixi_crm_b_order AS o ON o.id = i.order_id").
			Where("i.id = ? AND o.store_id = ?", id, storeID).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Take(&row).Error; err != nil {
			return err
		}
		if row.Status != "requested" {
			return errBadStatus
		}
		if req.Status == 1 {
			invoiceNo := req.Mark
			if invoiceNo == "" {
				invoiceNo = "INV-" + strconv.FormatUint(id, 10)
			}
			return tx.Table("qixi_crm_b_order_invoice").Where("id = ?", id).Updates(map[string]any{
				"status": "issued", "invoice_no": invoiceNo, "rejection_reason": "", "issued_at": time.Now().UTC(),
			}).Error
		}
		reason := req.Mark
		if reason == "" {
			reason = "商户驳回"
		}
		return tx.Table("qixi_crm_b_order_invoice").Where("id = ?", id).Updates(map[string]any{
			"status": "rejected", "rejection_reason": reason, "issued_at": nil,
		}).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "发票申请不存在")
		return
	}
	if errors.Is(err, errBadStatus) {
		response.Fail(c, http.StatusConflict, "当前发票状态不可审核")
		return
	}
	if err != nil {
		fail(c, "发票审核失败")
		return
	}
	response.OK(c, gin.H{"ok": true, "invoice_id": id, "status": req.Status})
}

func pageLimit(c *gin.Context) (int, int) {
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

func toLegacyStatus(v string) int {
	switch v {
	case "issued":
		return 1
	case "rejected", "voided":
		return -1
	default:
		return 0
	}
}

func fromLegacyStatus(v int) string {
	switch v {
	case 1:
		return "issued"
	case -1:
		return "rejected"
	case 0:
		return "requested"
	default:
		return ""
	}
}

func fromUIStatus(v string) string {
	switch v {
	case "requested", "issued", "rejected", "voided":
		return v
	default:
		return ""
	}
}

func fail(c *gin.Context, msg string) {
	response.Fail(c, http.StatusInternalServerError, msg)
}

var errBadStatus = statusError("bad status")

type statusError string

func (e statusError) Error() string { return string(e) }
