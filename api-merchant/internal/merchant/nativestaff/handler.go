// Package nativestaff owns store employee accounts after the qixi_crm_m_
// identity migration. It must not use the legacy qixi_m_* staff tables.
package nativestaff

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/setting/staff", h.list)
	r.POST("/setting/staff", middleware.RequireStorePermission(h.db, "staff.create"), h.create)
	r.PUT("/setting/staff/:id", middleware.RequireStorePermission(h.db, "staff.update"), h.update)
	r.DELETE("/setting/staff/:id", middleware.RequireStorePermission(h.db, "staff.delete"), h.remove)
}

type account struct {
	ID              uint64    `gorm:"column:id;primaryKey"`
	StoreID         uint64    `gorm:"column:store_id"`
	Username        string    `gorm:"column:username"`
	PasswordHash    string    `gorm:"column:password_hash"`
	RoleCode        string    `gorm:"column:role_code"`
	DisplayName     string    `gorm:"column:display_name"`
	Phone           string    `gorm:"column:phone"`
	CanAcceptOrders int8      `gorm:"column:can_accept_orders"`
	CanVerifyOrders int8      `gorm:"column:can_verify_orders"`
	CanShipOrders   int8      `gorm:"column:can_ship_orders"`
	Status          int8      `gorm:"column:status"`
	AuthVersion     uint64    `gorm:"column:auth_version"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}

func (account) TableName() string { return "qixi_crm_m_account" }

type saveRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	IsGoods  *int8  `json:"is_goods"`
	IsOpen   *int8  `json:"is_open"`
	IsVerify *int8  `json:"is_verify"`
	Status   *int8  `json:"status"`
}

func (h *Handler) list(c *gin.Context) {
	page, limit := pagination(c)
	q := h.db.WithContext(c.Request.Context()).Model(&account{}).
		Where("store_id = ? AND role_code <> ?", middleware.StoreID(c), "owner")
	if statusRaw := strings.TrimSpace(c.Query("status")); statusRaw != "" {
		if status, err := strconv.Atoi(statusRaw); err == nil && (status == 0 || status == 1) {
			q = q.Where("status = ?", status)
		}
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("(display_name LIKE ? OR phone LIKE ? OR username LIKE ?)", like, like, like)
	}
	switch strings.TrimSpace(c.Query("staff_scope")) {
	case "delivery":
		q = q.Where("role_code = ? OR can_ship_orders = 1", "delivery")
	case "service":
		q = q.Where("role_code = ? OR can_accept_orders = 1", "service")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询员工失败")
		return
	}
	var rows []account
	if err := q.Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询员工失败")
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, staffJSON(row))
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}

func (h *Handler) create(c *gin.Context) {
	var req saveRequest
	if err := c.ShouldBindJSON(&req); err != nil || !validCreate(req) {
		response.Fail(c, http.StatusBadRequest, "员工账号、昵称或密码不符合要求")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "员工密码处理失败")
		return
	}
	row := account{
		StoreID: uint64(middleware.StoreID(c)), Username: strings.TrimSpace(req.Account), PasswordHash: string(hash),
		RoleCode: "clerk", DisplayName: strings.TrimSpace(req.Nickname), Phone: strings.TrimSpace(req.Phone),
		CanAcceptOrders: boolInt(req.IsOpen, 1), CanVerifyOrders: boolInt(req.IsVerify, 1), CanShipOrders: boolInt(req.IsGoods, 1),
		Status: boolInt(req.Status, 1), AuthVersion: 1,
	}
	if err := h.db.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		if isDuplicate(err) {
			response.Fail(c, http.StatusConflict, "员工登录账号已存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "创建员工失败")
		return
	}
	h.audit(c, row.ID, "staff.create", strconv.FormatUint(row.ID, 10))
	response.OK(c, staffJSON(row))
}

func (h *Handler) update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "员工编号不正确")
		return
	}
	var req saveRequest
	if err := c.ShouldBindJSON(&req); err != nil || !validUpdate(req) {
		response.Fail(c, http.StatusBadRequest, "员工资料不符合要求")
		return
	}
	row, err := h.owned(c, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "员工不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询员工失败")
		return
	}
	updates := map[string]any{}
	if name := strings.TrimSpace(req.Nickname); name != "" && name != row.DisplayName {
		updates["display_name"] = name
	}
	if req.Phone != "" && strings.TrimSpace(req.Phone) != row.Phone {
		updates["phone"] = strings.TrimSpace(req.Phone)
	}
	if req.IsOpen != nil && *req.IsOpen != row.CanAcceptOrders {
		updates["can_accept_orders"] = boolInt(req.IsOpen, 0)
	}
	if req.IsVerify != nil && *req.IsVerify != row.CanVerifyOrders {
		updates["can_verify_orders"] = boolInt(req.IsVerify, 0)
	}
	if req.IsGoods != nil && *req.IsGoods != row.CanShipOrders {
		updates["can_ship_orders"] = boolInt(req.IsGoods, 0)
	}
	if req.Status != nil && *req.Status != row.Status {
		updates["status"] = boolInt(req.Status, 0)
	}
	if password := req.Password; password != "" {
		if len(password) < 8 {
			response.Fail(c, http.StatusBadRequest, "员工密码至少需要 8 位")
			return
		}
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			response.Fail(c, http.StatusInternalServerError, "员工密码处理失败")
			return
		}
		updates["password_hash"] = string(hash)
	}
	if len(updates) > 0 {
		if _, sessionChanged := updates["status"]; sessionChanged || updates["password_hash"] != nil {
			updates["auth_version"] = row.AuthVersion + 1
		}
		if err := h.db.WithContext(c.Request.Context()).Model(&account{}).Where("id = ? AND store_id = ?", row.ID, row.StoreID).Updates(updates).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "更新员工失败")
			return
		}
	}
	updated, err := h.owned(c, id)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取员工失败")
		return
	}
	h.audit(c, updated.ID, "staff.update", strconv.FormatUint(updated.ID, 10))
	response.OK(c, staffJSON(updated))
}

func (h *Handler) remove(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "员工编号不正确")
		return
	}
	row, err := h.owned(c, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "员工不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询员工失败")
		return
	}
	if row.ID == uint64(middleware.AdminID(c)) {
		response.Fail(c, http.StatusBadRequest, "不能移除当前登录账号")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Model(&account{}).Where("id = ? AND store_id = ?", row.ID, row.StoreID).Updates(map[string]any{"status": 0, "auth_version": row.AuthVersion + 1}).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "移除员工失败")
		return
	}
	h.audit(c, row.ID, "staff.remove", strconv.FormatUint(row.ID, 10))
	response.OK(c, gin.H{"service_id": row.ID, "removed": true})
}

func (h *Handler) owned(c *gin.Context, id uint64) (account, error) {
	var row account
	err := h.db.WithContext(c.Request.Context()).Where("id = ? AND store_id = ? AND role_code <> ?", id, middleware.StoreID(c), "owner").Take(&row).Error
	return row, err
}

func (h *Handler) audit(c *gin.Context, accountID uint64, action, resourceID string) {
	requestID := strings.TrimSpace(c.GetHeader("X-Request-Id"))
	if requestID == "" {
		requestID = fmt.Sprintf("staff-%d-%d", accountID, time.Now().UnixNano())
	}
	_ = h.db.WithContext(c.Request.Context()).Table("qixi_crm_m_operation_log").Create(map[string]any{
		"store_id": middleware.StoreID(c), "account_id": middleware.AdminID(c), "action": action, "resource_type": "staff", "resource_id": resourceID, "request_id": requestID,
	}).Error
}

func staffJSON(row account) gin.H {
	return gin.H{"service_id": row.ID, "account": row.Username, "nickname": row.DisplayName, "phone": row.Phone, "role_code": row.RoleCode, "status": row.Status, "is_open": row.CanAcceptOrders, "is_verify": row.CanVerifyOrders, "is_goods": row.CanShipOrders, "create_time": row.CreatedAt.Format("2006-01-02 15:04:05")}
}

func pagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}
func parseID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return id, err == nil && id > 0
}
func boolInt(value *int8, fallback int8) int8 {
	if value == nil {
		return fallback
	}
	if *value == 1 {
		return 1
	}
	return 0
}
func validCreate(req saveRequest) bool {
	return len(strings.TrimSpace(req.Account)) >= 3 && len(strings.TrimSpace(req.Account)) <= 64 && len(strings.TrimSpace(req.Nickname)) > 0 && len(strings.TrimSpace(req.Nickname)) <= 64 && len(req.Password) >= 8
}
func validUpdate(req saveRequest) bool {
	return len(strings.TrimSpace(req.Nickname)) <= 64 && len(req.Password) == 0 || (len(strings.TrimSpace(req.Nickname)) <= 64 && len(req.Password) >= 8)
}
func isDuplicate(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "duplicate")
}
