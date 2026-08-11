package invoice

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const menuRead = "accounts.invoice.read"

type Handler struct{ business, admin *gorm.DB }

func New(business, admin *gorm.DB) *Handler { return &Handler{business: business, admin: admin} }

func (h *Handler) Register(r gin.IRoutes) {
	read := []gin.HandlerFunc{
		middleware.RequireAdminRoles("platform", "operations"),
		middleware.RequireAdminMenu(h.admin, menuRead),
	}
	r.GET("/finance/invoices", append(read, h.List)...)
	r.GET("/finance/invoices/:id", append(read, h.Get)...)
}

type invoiceRow struct {
	ID              uint64     `gorm:"column:id"`
	OrderID         uint64     `gorm:"column:order_id"`
	MerchantID      uint64     `gorm:"column:merchant_id"`
	StoreID         uint64     `gorm:"column:store_id"`
	UserID          uint64     `gorm:"column:user_id"`
	OrderNo         string     `gorm:"column:order_no"`
	MerchantName    string     `gorm:"column:merchant_name"`
	StoreName       string     `gorm:"column:store_name"`
	Nickname        string     `gorm:"column:nickname"`
	UserMobile      string     `gorm:"column:user_mobile"`
	OrderStatus     string     `gorm:"column:order_status"`
	PayAmount       float64    `gorm:"column:pay_amount"`
	OrderCreatedAt  time.Time  `gorm:"column:order_created_at"`
	IsSystemDel     int        `gorm:"column:is_system_del"`
	HasRefunded     int        `gorm:"column:has_refunded"`
	RecipientJSON   []byte     `gorm:"column:recipient_snapshot"`
	ProfileType     string     `gorm:"column:profile_type"`
	InvoiceType     int        `gorm:"column:invoice_type"`
	ReceiptSN       string     `gorm:"column:receipt_sn"`
	InvoiceAmount   float64    `gorm:"column:invoice_amount"`
	Title           string     `gorm:"column:title"`
	TaxNo           string     `gorm:"column:tax_no"`
	Email           string     `gorm:"column:email"`
	Status          string     `gorm:"column:status"`
	InvoiceNo       string     `gorm:"column:invoice_no"`
	RejectionReason string     `gorm:"column:rejection_reason"`
	Mark            string     `gorm:"column:mark"`
	RequestedAt     time.Time  `gorm:"column:requested_at"`
	IssuedAt        *time.Time `gorm:"column:issued_at"`
}

type recipient struct {
	Name      string `json:"name"`
	Recipient string `json:"recipient"`
	Mobile    string `json:"mobile"`
	Address   string `json:"address"`
	Detail    string `json:"detail"`
	Province  string `json:"province"`
	City      string `json:"city"`
	District  string `json:"district"`
}

type listFilter struct {
	DateFrom    string
	DateTo      string
	MerID       uint64
	Status      string
	OrderType   string
	Keyword     string
	UserType    string
	UserKeyword string
}

func (h *Handler) List(c *gin.Context) {
	page, limit := pagination(c)
	f := parseFilter(c)
	q := h.base(c)
	q = applyFilters(q, f)
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
	return h.business.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order_invoice AS oi").
		Select(`oi.id, oi.order_id, o.order_no, o.merchant_id, o.store_id, o.user_id,
			o.merchant_name_snapshot AS merchant_name, o.store_name_snapshot AS store_name,
			IFNULL(u.nickname,'') AS nickname, IFNULL(u.mobile,'') AS user_mobile,
			o.status AS order_status, o.pay_amount, o.created_at AS order_created_at,
			o.is_system_del, o.recipient_snapshot,
			CASE WHEN EXISTS(
				SELECT 1 FROM qixi_crm_b_refund r WHERE r.order_id=o.id AND r.status='refunded'
			) THEN 1 ELSE 0 END AS has_refunded,
			oi.profile_type, IFNULL(oi.invoice_type,1) AS invoice_type,
			IFNULL(oi.receipt_sn,'') AS receipt_sn, IFNULL(oi.invoice_amount,0) AS invoice_amount,
			oi.title, oi.tax_no, oi.email, oi.status, oi.invoice_no, oi.rejection_reason,
			IFNULL(oi.mark,'') AS mark, oi.requested_at, oi.issued_at`).
		Joins("JOIN qixi_crm_b_order AS o ON o.id=oi.order_id").
		Joins("LEFT JOIN qixi_crm_b_user AS u ON u.id=o.user_id")
}

func parseFilter(c *gin.Context) listFilter {
	f := listFilter{
		DateFrom:    strings.TrimSpace(c.Query("date_from")),
		DateTo:      strings.TrimSpace(c.Query("date_to")),
		Status:      strings.TrimSpace(c.Query("status")),
		OrderType:   strings.TrimSpace(c.Query("order_type")),
		Keyword:     strings.TrimSpace(c.Query("keyword")),
		UserType:    strings.TrimSpace(c.Query("user_type")),
		UserKeyword: strings.TrimSpace(c.Query("user_keyword")),
	}
	if mer := strings.TrimSpace(c.Query("mer_id")); mer != "" {
		if id, err := strconv.ParseUint(mer, 10, 64); err == nil {
			f.MerID = id
		}
	}
	return f
}

func applyFilters(q *gorm.DB, f listFilter) *gorm.DB {
	if f.DateFrom != "" {
		q = q.Where("oi.requested_at >= ?", f.DateFrom+" 00:00:00")
	}
	if f.DateTo != "" {
		q = q.Where("oi.requested_at <= ?", f.DateTo+" 23:59:59")
	}
	if f.MerID > 0 {
		q = q.Where("o.merchant_id = ?", f.MerID)
	}
	switch f.Status {
	case "0", "requested", "unissued":
		q = q.Where("oi.status IN ?", []string{"requested", "rejected"})
	case "1", "issued":
		q = q.Where("oi.status = ?", "issued")
	case "rejected", "voided":
		q = q.Where("oi.status = ?", f.Status)
	}
	q = applyOrderType(q, f.OrderType)
	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		q = q.Where(`(
			o.order_no LIKE ? OR
			JSON_UNQUOTE(JSON_EXTRACT(o.recipient_snapshot,'$.recipient')) LIKE ? OR
			JSON_UNQUOTE(JSON_EXTRACT(o.recipient_snapshot,'$.name')) LIKE ? OR
			JSON_UNQUOTE(JSON_EXTRACT(o.recipient_snapshot,'$.mobile')) LIKE ? OR
			oi.title LIKE ? OR oi.receipt_sn LIKE ? OR oi.invoice_no LIKE ?
		)`, like, like, like, like, like, like, like)
	}
	if f.UserKeyword != "" {
		like := "%" + f.UserKeyword + "%"
		switch f.UserType {
		case "uid":
			if id, err := strconv.ParseUint(f.UserKeyword, 10, 64); err == nil {
				q = q.Where("o.user_id = ?", id)
			} else {
				q = q.Where("1 = 0")
			}
		case "phone":
			q = q.Where("u.mobile LIKE ?", like)
		default:
			q = q.Where("u.nickname LIKE ?", like)
		}
	}
	return q
}

func applyOrderType(q *gorm.DB, orderType string) *gorm.DB {
	switch orderType {
	case "1", "unpaid":
		return q.Where("o.status = ?", "pending_pay").Where("o.is_system_del = 0")
	case "2", "unshipped":
		return q.Where("o.status IN ?", []string{"paid", "fulfilling", "awaiting_final"}).Where("o.is_system_del = 0")
	case "3", "unreceived":
		return q.Where("o.status = ?", "shipped").Where("o.is_system_del = 0")
	case "4", "unevaluated":
		return q.Where(`o.is_system_del = 0 AND o.status = 'completed' AND EXISTS(
			SELECT 1 FROM qixi_crm_b_order_item oi2
			LEFT JOIN qixi_crm_b_product_comment pc ON pc.order_item_id=oi2.id AND pc.deleted_at IS NULL
			WHERE oi2.order_id=o.id AND pc.id IS NULL
		)`)
	case "5", "completed":
		return q.Where("o.status = ?", "completed").Where("o.is_system_del = 0")
	case "6", "refunded":
		return q.Where(`o.is_system_del = 0 AND EXISTS(
			SELECT 1 FROM qixi_crm_b_refund r WHERE r.order_id=o.id AND r.status='refunded'
		)`)
	case "7", "deleted":
		return q.Where("o.is_system_del = 1")
	default:
		return q
	}
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

func views(rows []invoiceRow) []gin.H {
	out := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		recv := parseRecipient(r.RecipientJSON)
		contactName := strings.TrimSpace(recv.Recipient)
		if contactName == "" {
			contactName = strings.TrimSpace(recv.Name)
		}
		contactInfo := strings.TrimSpace(r.Email)
		if contactInfo == "" {
			contactInfo = strings.TrimSpace(recv.Mobile)
		}
		if contactInfo == "" {
			contactInfo = buildAddress(recv)
		}
		amount := r.InvoiceAmount
		if amount <= 0 {
			amount = r.PayAmount
		}
		receiptSN := strings.TrimSpace(r.ReceiptSN)
		if receiptSN == "" {
			receiptSN = strings.TrimSpace(r.InvoiceNo)
		}
		statusLabel := invoiceStatusLabel(r.Status)
		orderStatusLabel := orderStatusText(r)
		out = append(out, gin.H{
			"id":                 r.ID,
			"order_id":           r.OrderID,
			"order_sn":           r.OrderNo,
			"order_no":           r.OrderNo,
			"mer_id":             r.MerchantID,
			"merchant_id":        r.MerchantID,
			"mer_name":           emptyDash(r.MerchantName),
			"merchant_name":      emptyDash(r.MerchantName),
			"store_id":           r.StoreID,
			"store_name":         emptyDash(r.StoreName),
			"uid":                r.UserID,
			"nickname":           emptyDash(r.Nickname),
			"user_phone_mask":    maskPhone(r.UserMobile),
			"pay_price":          r.PayAmount,
			"order_amount":       r.PayAmount,
			"order_status":       r.OrderStatus,
			"order_status_label": orderStatusLabel,
			"invoice_amount":     amount,
			"receipt_price":      amount,
			"receipt_sn":         receiptSN,
			"invoice_no":         strings.TrimSpace(r.InvoiceNo),
			"invoice_type":       r.InvoiceType,
			"invoice_type_label": invoiceTypeLabel(r.InvoiceType),
			"detail_title":       detailTitle(r.InvoiceType, r.ProfileType),
			"profile_type":       r.ProfileType,
			"title_type_label":   titleTypeLabel(r.ProfileType),
			"title":              strings.TrimSpace(r.Title),
			"contact_name":       emptyDash(contactName),
			"contact_info":       emptyDash(contactInfo),
			"tax_no":             strings.TrimSpace(r.TaxNo),
			"tax_no_masked":      maskTax(r.TaxNo),
			"email":              strings.TrimSpace(r.Email),
			"email_masked":       maskEmail(r.Email),
			"status":             r.Status,
			"status_label":       statusLabel,
			"issued":             r.Status == "issued",
			"mark":               strings.TrimSpace(r.Mark),
			"rejection_reason":   strings.TrimSpace(r.RejectionReason),
			"create_time":        r.OrderCreatedAt.Format("2006-01-02 15:04:05"),
			"requested_at":       r.RequestedAt.Format("2006-01-02 15:04:05"),
			"issued_at":          formatTimePtr(r.IssuedAt),
		})
	}
	return out
}

func detailTitle(invoiceType int, profileType string) string {
	if invoiceType == 2 {
		return "企业专用纸质发票"
	}
	if profileType == "enterprise" {
		return "企业电子普通发票"
	}
	return "个人电子普通发票"
}

func parseRecipient(raw []byte) recipient {
	var r recipient
	if len(raw) == 0 {
		return r
	}
	_ = json.Unmarshal(raw, &r)
	return r
}

func buildAddress(r recipient) string {
	if strings.TrimSpace(r.Address) != "" {
		return strings.TrimSpace(r.Address)
	}
	return strings.TrimSpace(r.Province + r.City + r.District + " " + r.Detail)
}

func orderStatusText(r invoiceRow) string {
	if r.HasRefunded == 1 {
		return "已退款"
	}
	if r.IsSystemDel == 1 {
		return "已删除"
	}
	switch r.OrderStatus {
	case "pending_pay":
		return "待付款"
	case "paid", "fulfilling", "awaiting_final":
		return "待发货"
	case "shipped":
		return "待收货"
	case "completed":
		return "已完成"
	case "cancelled":
		return "已取消"
	case "aftersale":
		return "售后中"
	default:
		return "-"
	}
}

func invoiceStatusLabel(status string) string {
	// CRMEB 详情/列表开票态仅为「已开 / 未开」
	if status == "issued" {
		return "已开"
	}
	return "未开"
}

func invoiceTypeLabel(t int) string {
	if t == 2 {
		return "专用发票"
	}
	return "普通发票"
}

func titleTypeLabel(profileType string) string {
	if profileType == "enterprise" {
		return "企业"
	}
	return "个人"
}

func emptyDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "-"
	}
	return v
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
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

func maskPhone(v string) string {
	v = strings.TrimSpace(v)
	runes := []rune(v)
	if len(runes) < 7 {
		return emptyDash(v)
	}
	return string(runes[:3]) + "****" + string(runes[len(runes)-4:])
}
