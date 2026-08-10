package svip

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/queryfilter"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc                 *identity.Service
	adminDB, businessDB *gorm.DB
}

func NewHandler(svc *identity.Service, adminDB, businessDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB, businessDB: businessDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	platformOnly := middleware.RequireAdminRoles("platform")
	manage := middleware.RequireAdminMenu(h.adminDB, "user.svip.manage")
	r.GET("/users", platformOnly, manage, h.ListUsers)
	r.PUT("/users/:id/svip", platformOnly, manage, h.SetSvip)
	operations := middleware.RequireAdminRoles("platform", "operations")
	planRead := middleware.RequireAdminMenu(h.adminDB, "user.svip.plan.read")
	planManage := middleware.RequireAdminMenu(h.adminDB, "user.svip.plan.manage")
	recordRead := middleware.RequireAdminMenu(h.adminDB, "user.svip.record.read")
	r.GET("/svip/plans", operations, planRead, h.ListPlans)
	r.POST("/svip/plans", operations, planManage, h.CreatePlan)
	r.PUT("/svip/plans/:id", operations, planManage, h.UpdatePlan)
	r.PUT("/svip/plans/:id/status", operations, planManage, h.UpdatePlanStatus)
	r.DELETE("/svip/plans/:id", operations, planManage, h.DeletePlan)
	r.GET("/svip/orders", operations, recordRead, h.ListOrders)
	r.GET("/svip/orders/summary", operations, recordRead, h.OrderSummary)
}

// plan and order are supervision projections.  The C-end payment callback is
// the only place allowed to create paid orders or grant membership benefits.
type plan struct {
	ID           uint64  `gorm:"column:id" json:"id"`
	Name         string  `gorm:"column:name" json:"name"`
	CostPrice    float64 `gorm:"column:cost_price" json:"cost_price"`
	Price        float64 `gorm:"column:price" json:"price"`
	PlanType     string  `gorm:"column:plan_type" json:"plan_type"`
	DurationDays *int    `gorm:"column:duration_days" json:"duration_days,omitempty"`
	Benefits     string  `gorm:"column:benefits" json:"benefits"`
	Status       int     `gorm:"column:status" json:"status"`
	Sort         int     `gorm:"column:sort" json:"sort"`
}

type planInput struct {
	Name         string   `json:"name"`
	CostPrice    float64  `json:"cost_price"`
	Price        float64  `json:"price"`
	PlanType     string   `json:"plan_type"`
	DurationDays int      `json:"duration_days"`
	Benefits     []string `json:"benefits"`
	Status       int      `json:"status"`
	Sort         int      `json:"sort"`
}

type order struct {
	ID           uint64  `gorm:"column:id" json:"id"`
	OrderNo      string  `gorm:"column:order_no" json:"order_no"`
	UserID       uint64  `gorm:"column:user_id" json:"user_id"`
	Nickname     string  `gorm:"column:nickname" json:"nickname"`
	Phone        string  `gorm:"column:phone" json:"phone"`
	PlanID       uint64  `gorm:"column:plan_id" json:"plan_id"`
	PlanName     string  `gorm:"column:plan_name" json:"plan_name"`
	PlanType     string  `gorm:"column:plan_type" json:"plan_type"`
	DurationDays *int    `gorm:"column:duration_days" json:"duration_days,omitempty"`
	Amount       float64 `gorm:"column:amount" json:"amount"`
	PayType      string  `gorm:"column:pay_type" json:"pay_type"`
	Status       string  `gorm:"column:status" json:"status"`
	CreatedAt    string  `gorm:"column:created_at" json:"created_at"`
	PaidAt       *string `gorm:"column:paid_at" json:"paid_at,omitempty"`
	EndTime      *string `gorm:"column:end_time" json:"end_time,omitempty"`
}

// orderSummary aligns CRMEB UserOrderRepository::countMemberInfo.
type orderSummary struct {
	PaidMemberCount    int64   `gorm:"column:paid_member_count" json:"paid_member_count"`
	PaidAmount         float64 `gorm:"column:paid_amount" json:"paid_amount"`
	ExpiredMemberCount int64   `gorm:"column:expired_member_count" json:"expired_member_count"`
}

func (h *Handler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListUsers(c.Request.Context(), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	list := make([]gin.H, 0, len(res.List))
	for _, u := range res.List {
		list = append(list, gin.H{
			"uid": u.UID, "nickname": u.Nickname, "phone_masked": maskPhone(u.Phone),
			"is_svip": u.IsSvip, "svip_endtime": u.SvipEndtime,
			"is_svip_active": identity.UserSvipActive(&u),
			"integral":       u.Integral, "now_money": u.NowMoney,
		})
	}
	response.OK(c, gin.H{"list": list, "total": res.Total, "page": res.Page, "limit": res.Limit})
}

func maskPhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if len(phone) < 7 {
		return ""
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

func (h *Handler) SetSvip(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		IsSvip      int8   `json:"is_svip"`
		SvipEndtime string `json:"svip_endtime"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	in := identity.SvipInput{IsSvip: body.IsSvip}
	if body.SvipEndtime != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", body.SvipEndtime, time.Local); err == nil {
			in.SvipEndtime = &t
		} else if t, err := time.ParseInLocation("2006-01-02", body.SvipEndtime, time.Local); err == nil {
			in.SvipEndtime = &t
		}
	}
	u, err := h.svc.SetUserSvip(c.Request.Context(), uint(id), in)
	if err != nil {
		if errors.Is(err, identity.ErrNotFound) {
			response.Fail(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, identity.ErrBadParam) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "操作失败")
		return
	}
	response.OK(c, gin.H{
		"uid": u.UID, "is_svip": u.IsSvip, "svip_endtime": u.SvipEndtime,
		"is_svip_active": identity.UserSvipActive(u),
	})
}

func (h *Handler) ListPlans(c *gin.Context) {
	page, limit := svipPagination(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_svip_plan")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		svipFail(c, "会员类型查询失败")
		return
	}
	var rows []plan
	if err := q.Select("id,name,cost_price,price,plan_type,duration_days,benefits,status,sort").
		Order("sort ASC,id ASC").
		Offset((page - 1) * limit).
		Limit(limit).
		Scan(&rows).Error; err != nil {
		svipFail(c, "会员类型查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) CreatePlan(c *gin.Context) {
	var in planInput
	if c.ShouldBindJSON(&in) != nil || !validPlan(&in) {
		response.Fail(c, http.StatusBadRequest, "会员类型参数错误")
		return
	}
	if !h.activeInterestNames(c, in.Benefits) {
		return
	}
	row := planFromInput(&in)
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_svip_plan").Create(&row).Error; err != nil {
		svipFail(c, "会员类型创建失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdatePlan(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in planInput
	if id == 0 || c.ShouldBindJSON(&in) != nil || !validPlan(&in) {
		response.Fail(c, http.StatusBadRequest, "会员类型参数错误")
		return
	}
	if !h.activeInterestNames(c, in.Benefits) {
		return
	}
	row := planFromInput(&in)
	values := map[string]any{
		"name": row.Name, "cost_price": row.CostPrice, "price": row.Price, "plan_type": row.PlanType,
		"duration_days": row.DurationDays, "benefits": row.Benefits, "status": row.Status, "sort": row.Sort,
	}
	res := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_svip_plan").Where("id = ?", id).Updates(values)
	if res.Error != nil {
		svipFail(c, "会员类型更新失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "会员类型不存在")
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *Handler) UpdatePlanStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		Status int `json:"status"`
	}
	if id == 0 || c.ShouldBindJSON(&body) != nil || (body.Status != 0 && body.Status != 1) {
		response.Fail(c, http.StatusBadRequest, "会员类型状态参数错误")
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_svip_plan").Where("id = ?", id).Update("status", body.Status)
	if res.Error != nil {
		svipFail(c, "会员类型状态更新失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "会员类型不存在")
		return
	}
	response.OK(c, gin.H{"id": id, "status": body.Status})
}

func (h *Handler) DeletePlan(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "会员类型参数错误")
		return
	}
	var orderCount int64
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_svip_order").Where("plan_id = ?", id).Count(&orderCount).Error; err != nil {
		svipFail(c, "会员类型删除校验失败")
		return
	}
	if orderCount > 0 {
		response.Fail(c, http.StatusBadRequest, "该会员类型已有购买记录，无法删除")
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Exec("DELETE FROM qixi_crm_b_svip_plan WHERE id = ?", id)
	if res.Error != nil {
		svipFail(c, "会员类型删除失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "会员类型不存在")
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *Handler) activeInterestNames(c *gin.Context, benefits []string) bool {
	names := normalizedBenefits(benefits)
	if len(names) == 0 {
		return true
	}
	var total int64
	err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_svip_interest").Where("deleted_at IS NULL AND status=1 AND name IN ?", names).Count(&total).Error
	if err != nil {
		svipFail(c, "会员权益校验失败")
		return false
	}
	if total != int64(len(names)) {
		response.Fail(c, http.StatusBadRequest, "会员权益不存在或已停用")
		return false
	}
	return true
}

func (h *Handler) ListOrders(c *gin.Context) {
	page, limit := svipPagination(c)
	q := applySvipOrderFilters(h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_svip_order AS o").
		Joins("LEFT JOIN qixi_crm_b_user AS u ON u.id = o.user_id"), c)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		svipFail(c, "会员记录查询失败")
		return
	}
	var rows []order
	if err := applySvipOrderFilters(h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_svip_order AS o").
		Joins("LEFT JOIN qixi_crm_b_user AS u ON u.id = o.user_id").
		Joins("LEFT JOIN qixi_crm_b_user_svip AS sv ON sv.user_id = o.user_id"), c).
		Select(`o.id,o.order_no,o.user_id,COALESCE(u.nickname,'') AS nickname,COALESCE(u.mobile,'') AS phone,
o.plan_id,o.plan_name,o.plan_type,o.duration_days,o.amount,COALESCE(o.pay_type,'') AS pay_type,o.status,
o.created_at,o.paid_at,COALESCE(o.end_time, sv.expires_at) AS end_time`).
		Order("o.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		svipFail(c, "会员记录查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) OrderSummary(c *gin.Context) {
	var out orderSummary
	paidQ := applySvipOrderFilters(h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_svip_order AS o").
		Joins("LEFT JOIN qixi_crm_b_user AS u ON u.id = o.user_id").
		Where("o.status = ?", "paid"), c)
	if err := paidQ.Select(`COUNT(DISTINCT o.user_id) AS paid_member_count,
COALESCE(SUM(o.amount), 0) AS paid_amount`).Scan(&out).Error; err != nil {
		svipFail(c, "会员记录统计失败")
		return
	}
	// CRMEB：累计已过期人数按当前会员状态统计（非订单行）。
	var expired int64
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user_svip").
		Where("status <> ? AND expires_at IS NOT NULL AND expires_at < NOW()", "lifetime").
		Count(&expired).Error; err != nil {
		svipFail(c, "会员记录统计失败")
		return
	}
	out.ExpiredMemberCount = expired
	response.OK(c, out)
}

func applySvipOrderFilters(q *gorm.DB, c *gin.Context) *gorm.DB {
	if payType := strings.TrimSpace(c.Query("pay_type")); payType != "" {
		switch payType {
		case "weixin", "wechat":
			q = q.Where("o.pay_type IN ?", []string{"weixin", "wechat"})
		case "alipay", "routine", "sys", "free":
			q = q.Where("o.pay_type = ?", payType)
		default:
			// ignore unknown pay_type to avoid hard-failing list filters
		}
	}
	if title := strings.TrimSpace(c.Query("title")); title != "" {
		q = q.Where("o.plan_name LIKE ?", "%"+title+"%")
	}
	if nickname := strings.TrimSpace(c.Query("nickname")); nickname != "" {
		q = q.Where("u.nickname LIKE ?", "%"+nickname+"%")
	}
	// 兼容旧 keyword：订单号 / 用户 ID / 昵称
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("o.order_no LIKE ? OR CAST(o.user_id AS CHAR) LIKE ? OR u.nickname LIKE ?", like, like, like)
	}
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		if status == "pending" || status == "paid" || status == "closed" {
			q = q.Where("o.status = ?", status)
		}
	}
	return queryfilter.ApplyCreatedAtRange(q, c, "o.created_at")
}

func planFromInput(in *planInput) plan {
	benefits, _ := json.Marshal(normalizedBenefits(in.Benefits))
	row := plan{
		Name: strings.TrimSpace(in.Name), CostPrice: in.CostPrice, Price: in.Price,
		PlanType: in.PlanType, Benefits: string(benefits), Status: in.Status, Sort: in.Sort,
	}
	if in.PlanType != "lifetime" {
		days := in.DurationDays
		row.DurationDays = &days
	}
	return row
}

func validPlan(in *planInput) bool {
	if in == nil || strings.TrimSpace(in.Name) == "" || len([]rune(strings.TrimSpace(in.Name))) > 64 ||
		in.CostPrice < 0 || in.Price < 0 || (in.Status != 0 && in.Status != 1) || in.Sort < 0 || in.Sort > 999999 {
		return false
	}
	benefits := normalizedBenefits(in.Benefits)
	if len(benefits) > 10 {
		return false
	}
	rawNonEmpty := 0
	for _, v := range in.Benefits {
		if strings.TrimSpace(v) != "" {
			rawNonEmpty++
		}
	}
	if rawNonEmpty != len(benefits) {
		return false
	}
	switch in.PlanType {
	case "trial":
		return in.Price == 0 && in.DurationDays > 0 && in.DurationDays <= 31
	case "period":
		return in.Price > 0 && in.DurationDays > 0 && in.DurationDays <= 3660
	case "lifetime":
		return in.Price > 0 && in.DurationDays == 0
	default:
		return false
	}
}

func normalizedBenefits(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || len([]rune(value)) > 32 {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func svipPagination(c *gin.Context) (int, int) {
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

func svipFail(c *gin.Context, message string) {
	response.Fail(c, http.StatusInternalServerError, message)
}
