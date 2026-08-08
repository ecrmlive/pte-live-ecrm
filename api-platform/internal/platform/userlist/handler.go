package userlist

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/queryfilter"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct{ business, admin *gorm.DB }

func New(business, admin *gorm.DB) *Handler { return &Handler{business: business, admin: admin} }

func (h *Handler) Register(r gin.IRoutes) {
	platformRead := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.list.read")}
	r.GET("/user-list", append(platformRead, h.List)...)
	r.GET("/user-list/:id/detail", append(platformRead, h.Detail)...)
	exportUsers := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.list.export")}
	r.POST("/user-list/export", append(exportUsers, h.Export)...)
	userCreate := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.create.execute")}
	r.POST("/user-list", append(userCreate, h.CreateUser)...)
	profileManage := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.profile.manage")}
	r.PUT("/user-list/:id/profile", append(profileManage, h.UpdateProfile)...)
	passwordManage := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.password.reset")}
	r.POST("/user-list/:id/password", append(passwordManage, h.ResetPassword)...)
	platformManage := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.asset.adjust")}
	r.POST("/user-list/:id/assets/adjust", append(platformManage, h.AdjustAsset)...)
	statusManage := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.status.manage")}
	r.POST("/user-list/:id/status", append(statusManage, h.ChangeStatus)...)
	r.GET("/user-list/member-levels", append(platformRead, h.ListMemberLevels)...)
	r.GET("/user-list/coupon-templates", append(platformRead, h.ListCouponTemplates)...)
	r.GET("/user-list/groups", append(platformRead, h.ListGroups)...)
	r.GET("/user-list/labels", append(platformRead, h.ListLabels)...)
	memberManage := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.member.adjust")}
	r.POST("/user-list/:id/member-level", append(memberManage, h.AdjustMemberLevel)...)
	couponManage := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.coupon.manage")}
	r.POST("/user-list/:id/coupons/:coupon_id/issue", append(couponManage, h.IssueCoupon)...)
	r.POST("/user-list/:id/coupons/:coupon_id/revoke", append(couponManage, h.RevokeCoupon)...)
	couponCommandRead := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "marketing.coupon.send.read")}
	r.GET("/user-list/coupon-commands", append(couponCommandRead, h.ListCouponCommands)...)
	couponRecordRead := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "marketing.coupon.record.read")}
	r.GET("/user-list/coupon-records", append(couponRecordRead, h.ListCouponRecords)...)
	referrerManage := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.referrer.manage")}
	r.POST("/user-list/:id/referrer", append(referrerManage, h.ChangeReferrer)...)
	groupManage := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.group.assign")}
	r.POST("/user-list/:id/group", append(groupManage, h.ChangeUserGroup)...)
	r.POST("/user-list/groups/assign", append(groupManage, h.ChangeUserGroups)...)
	labelManage := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.label.assign")}
	r.POST("/user-list/:id/labels", append(labelManage, h.ChangeUserLabels)...)
	r.POST("/user-list/labels/assign", append(labelManage, h.ChangeUsersLabels)...)
	promoterManage := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.promoter.manage")}
	r.POST("/user-list/promoters/assign", append(promoterManage, h.ChangePromoters)...)
	notificationManage := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.notification.send")}
	r.POST("/user-list/:id/notifications", append(notificationManage, h.SendInAppNotification)...)
}

type memberLevel struct {
	ID     uint64 `gorm:"column:id" json:"id"`
	Name   string `gorm:"column:name" json:"name"`
	Rank   int    `gorm:"column:rank" json:"rank"`
	Status int    `gorm:"column:status" json:"status"`
}

type couponTemplate struct {
	CouponID      uint64  `gorm:"column:coupon_id" json:"coupon_id"`
	Name          string  `gorm:"column:name" json:"name"`
	StoreID       uint64  `gorm:"column:store_id" json:"store_id"`
	DiscountType  string  `gorm:"column:discount_type" json:"discount_type"`
	DiscountValue float64 `gorm:"column:discount_value" json:"discount_value"`
	MinAmount     float64 `gorm:"column:min_amount" json:"min_amount"`
}

// couponCommandRow is an immutable platform manual coupon command. It intentionally
// contains user/operator IDs rather than account, phone, or other personal details.
type couponCommandRow struct {
	ID              uint64    `gorm:"column:id" json:"id"`
	UserID          uint64    `gorm:"column:user_id" json:"user_id"`
	CouponID        uint64    `gorm:"column:coupon_id" json:"coupon_id"`
	CouponUserID    uint64    `gorm:"column:coupon_user_id" json:"coupon_user_id"`
	CouponName      string    `gorm:"column:coupon_name" json:"coupon_name"`
	StoreID         uint64    `gorm:"column:store_id" json:"store_id"`
	Action          string    `gorm:"column:action" json:"action"`
	FromStatus      string    `gorm:"column:from_status" json:"from_status"`
	ToStatus        string    `gorm:"column:to_status" json:"to_status"`
	Reason          string    `gorm:"column:reason" json:"reason"`
	OperatorAdminID uint64    `gorm:"column:operator_admin_id" json:"operator_admin_id"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
}

// couponRecordRow is a user coupon lifecycle projection. It never joins a user
// profile so the platform record page does not expose account or phone data.
type couponRecordRow struct {
	ID          uint64    `gorm:"column:id" json:"id"`
	UserID      uint64    `gorm:"column:user_id" json:"user_id"`
	CouponID    uint64    `gorm:"column:coupon_id" json:"coupon_id"`
	CouponName  string    `gorm:"column:coupon_name" json:"coupon_name"`
	StoreID     uint64    `gorm:"column:store_id" json:"store_id"`
	Source      string    `gorm:"column:source" json:"source"`
	Status      string    `gorm:"column:status" json:"status"`
	ObtainedAt  time.Time `gorm:"column:obtained_at" json:"obtained_at"`
	UsedOrderID *uint64   `gorm:"column:used_order_id" json:"used_order_id"`
}

type userGroup struct {
	GroupID   uint64 `gorm:"column:group_id" json:"group_id"`
	GroupName string `gorm:"column:group_name" json:"group_name"`
	Sort      int    `gorm:"column:sort" json:"sort"`
}

type userLabel struct {
	LabelID   uint64 `gorm:"column:label_id" json:"label_id"`
	LabelName string `gorm:"column:label_name" json:"label_name"`
	Sort      int    `gorm:"column:sort" json:"sort"`
}

func (h *Handler) ListMemberLevels(c *gin.Context) {
	levels := make([]memberLevel, 0)
	if err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_member_level AS l").Select("l.id,l.name,l.rank,l.status").Where("l.status=1 AND l.deleted_at IS NULL").Order("l.rank ASC,l.id ASC").Scan(&levels).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "会员等级查询失败")
		return
	}
	response.OK(c, gin.H{"list": levels})
}

func (h *Handler) ListCouponTemplates(c *gin.Context) {
	rows := make([]couponTemplate, 0)
	now := time.Now()
	if err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_coupon_template_view").
		Select("coupon_id,name,store_id,discount_type,discount_value,min_amount").
		Where("status=1 AND (starts_at IS NULL OR starts_at<=?) AND (ends_at IS NULL OR ends_at>=?)", now, now).
		Order("store_id ASC,coupon_id ASC").Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "优惠券模板查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

// ListCouponCommands exposes the platform's own issue/revoke audit trail. Coupon
// lifecycle writes remain in commandCoupon and checkout/payment; this endpoint is
// deliberately read-only and never returns an idempotency key.
func (h *Handler) ListCouponCommands(c *gin.Context) {
	userID, ok := optionalPositiveQuery(c, "user_id")
	if !ok {
		response.Fail(c, http.StatusBadRequest, "用户 ID 参数错误")
		return
	}
	couponID, ok := optionalPositiveQuery(c, "coupon_id")
	if !ok {
		response.Fail(c, http.StatusBadRequest, "优惠券 ID 参数错误")
		return
	}
	action := strings.TrimSpace(c.Query("action"))
	if action != "" && action != "issue" && action != "revoke" {
		response.Fail(c, http.StatusBadRequest, "优惠券操作类型错误")
		return
	}

	db := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_user_coupon_command_audit AS a").
		Joins("LEFT JOIN qixi_crm_b_coupon_template_view AS ct ON ct.coupon_id=a.coupon_id")
	if userID != 0 {
		db = db.Where("a.user_id=?", userID)
	}
	if couponID != 0 {
		db = db.Where("a.coupon_id=?", couponID)
	}
	if action != "" {
		db = db.Where("a.action=?", action)
	}
	db = queryfilter.ApplyCreatedAtRange(db, c, "a.created_at")
	var total int64
	if err := db.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "优惠券发送记录查询失败")
		return
	}
	page, limit := paging(c)
	rows := make([]couponCommandRow, 0)
	if err := db.Select("a.id,a.user_id,a.coupon_id,a.coupon_user_id,COALESCE(ct.name,'已删除优惠券模板') AS coupon_name,COALESCE(ct.store_id,0) AS store_id,a.action,a.from_status,a.to_status,a.reason,a.operator_admin_id,a.created_at").
		Order("a.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "优惠券发送记录查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total})
}

// ListCouponRecords is read-only supervision of actual obtained/used coupon
// facts. It does not attempt to infer or mutate checkout's locked/used states.
func (h *Handler) ListCouponRecords(c *gin.Context) {
	userID, ok := optionalPositiveQuery(c, "user_id")
	if !ok {
		response.Fail(c, http.StatusBadRequest, "用户 ID 参数错误")
		return
	}
	couponID, ok := optionalPositiveQuery(c, "coupon_id")
	if !ok {
		response.Fail(c, http.StatusBadRequest, "优惠券 ID 参数错误")
		return
	}
	status := strings.TrimSpace(c.Query("status"))
	if status != "" && status != "unused" && status != "locked" && status != "used" && status != "expired" {
		response.Fail(c, http.StatusBadRequest, "优惠券状态错误")
		return
	}

	db := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_coupon_user AS cu").
		Joins("LEFT JOIN qixi_crm_b_coupon_template_view AS ct ON ct.coupon_id=cu.coupon_id")
	if userID != 0 {
		db = db.Where("cu.user_id=?", userID)
	}
	if couponID != 0 {
		db = db.Where("cu.coupon_id=?", couponID)
	}
	if status != "" {
		db = db.Where("cu.status=?", status)
	}
	db = queryfilter.ApplyCreatedAtRange(db, c, "cu.obtained_at")
	var total int64
	if err := db.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "优惠券领取记录查询失败")
		return
	}
	page, limit := paging(c)
	rows := make([]couponRecordRow, 0)
	if err := db.Select("cu.id,cu.user_id,cu.coupon_id,COALESCE(ct.name,'已删除优惠券模板') AS coupon_name,COALESCE(ct.store_id,0) AS store_id,cu.source,cu.status,cu.obtained_at,cu.used_order_id").
		Order("cu.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "优惠券领取记录查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total})
}

func (h *Handler) ListGroups(c *gin.Context) {
	rows := make([]userGroup, 0)
	if err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_user_group").
		Select("group_id,group_name,sort").Where("is_del=0").Order("sort DESC,group_id DESC").Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "用户分组查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) ListLabels(c *gin.Context) {
	rows := make([]userLabel, 0)
	if err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_user_label").
		Select("label_id,label_name,sort").Where("is_del=0").Order("sort DESC,label_id DESC").Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "用户标签查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

type createUserInput struct {
	Account        string `json:"account"`
	Password       string `json:"password"`
	Nickname       string `json:"nickname"`
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

type updateUserProfileInput struct {
	Nickname       string `json:"nickname"`
	AvatarURL      string `json:"avatar_url"`
	Gender         int    `json:"gender"`
	Bio            string `json:"bio"`
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

type resetUserPasswordInput struct {
	Password       string `json:"password"`
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

var (
	errUserAdminNotFound = errors.New("user admin command target not found")
	errUserAdminConflict = errors.New("user admin command conflict")
)

func (h *Handler) CreateUser(c *gin.Context) {
	var in createUserInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "新增用户参数错误")
		return
	}
	in.Account, in.Nickname, in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Account), strings.TrimSpace(in.Nickname), strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if !validUserCommand(in.Reason, in.IdempotencyKey) || len([]rune(in.Account)) < 3 || len([]rune(in.Account)) > 191 || len([]rune(in.Nickname)) < 1 || len([]rune(in.Nickname)) > 64 || len(in.Password) < 12 || len(in.Password) > 72 {
		response.Fail(c, http.StatusBadRequest, "账号、昵称、密码、原因或幂等键错误")
		return
	}
	fingerprint := fingerprint("create", in.Account, in.Nickname, in.Reason)
	out := gin.H{}
	err := h.business.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var audit struct {
			UserID             uint64 `gorm:"column:user_id"`
			RequestFingerprint string `gorm:"column:request_fingerprint"`
			Reason             string `gorm:"column:reason"`
			OperatorID         uint64 `gorm:"column:operator_admin_id"`
		}
		if e := tx.Table("qixi_crm_b_user_admin_command_audit").Where("action='create' AND idempotency_key=?", in.IdempotencyKey).Take(&audit).Error; e == nil {
			// Password input is intentionally excluded from the immutable audit.
			// A duplicate key is therefore never replayed: accepting it could hide
			// a changed initial password behind a successful retry response.
			if audit.RequestFingerprint != fingerprint || audit.Reason != in.Reason || audit.OperatorID != uint64(middleware.AdminID(c)) {
				return errUserAdminConflict
			}
			return errUserAdminConflict
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		var count int64
		if e := tx.Table("qixi_crm_b_user_identity").Where("channel='pc' AND subject=?", in.Account).Count(&count).Error; e != nil {
			return e
		}
		if count != 0 {
			return errUserAdminConflict
		}
		hash, e := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if e != nil {
			return e
		}
		user := struct {
			ID          uint64 `gorm:"column:id"`
			Nickname    string `gorm:"column:nickname"`
			Status      int    `gorm:"column:status"`
			AuthVersion uint64 `gorm:"column:auth_version"`
		}{Nickname: in.Nickname, Status: 1, AuthVersion: 1}
		if e := tx.Table("qixi_crm_b_user").Create(&user).Error; e != nil {
			return e
		}
		userID := user.ID
		if userID == 0 {
			return errors.New("created user id unavailable")
		}
		if e := tx.Table("qixi_crm_b_user_identity").Create(map[string]any{"user_id": userID, "channel": "pc", "subject": in.Account, "credential_hash": string(hash)}).Error; e != nil {
			if isDuplicateKey(e) {
				return errUserAdminConflict
			}
			return e
		}
		if e := tx.Table("qixi_crm_b_user_profile").Create(map[string]any{"user_id": userID, "source_channel": "pc"}).Error; e != nil {
			return e
		}
		if e := tx.Table("qixi_crm_b_user_admin_command_audit").Create(map[string]any{"action": "create", "user_id": userID, "request_fingerprint": fingerprint, "reason": in.Reason, "operator_admin_id": middleware.AdminID(c), "idempotency_key": in.IdempotencyKey}).Error; e != nil {
			return e
		}
		out = gin.H{"user_id": userID, "replayed": false}
		return nil
	})
	writeUserAdminResult(c, err, out, "新增用户失败")
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, ok := positiveID(c.Param("id"))
	var in updateUserProfileInput
	if !ok || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "用户资料参数错误")
		return
	}
	in.Nickname, in.AvatarURL, in.Bio, in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Nickname), strings.TrimSpace(in.AvatarURL), strings.TrimSpace(in.Bio), strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if !validUserCommand(in.Reason, in.IdempotencyKey) || len([]rune(in.Nickname)) < 1 || len([]rune(in.Nickname)) > 64 || len([]rune(in.Bio)) > 500 || in.Gender < 0 || in.Gender > 2 || !validAvatarURL(in.AvatarURL) {
		response.Fail(c, http.StatusBadRequest, "资料、原因或幂等键错误")
		return
	}
	fingerprint := fingerprint("profile", strconv.FormatUint(userID, 10), in.Nickname, in.AvatarURL, strconv.Itoa(in.Gender), in.Bio, in.Reason)
	out := gin.H{}
	err := h.business.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var user struct{ ID uint64 }
		if e := tx.Table("qixi_crm_b_user").Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").Where("id=?", userID).Take(&user).Error; e != nil {
			return e
		}
		var audit struct {
			UserID             uint64 `gorm:"column:user_id"`
			RequestFingerprint string `gorm:"column:request_fingerprint"`
			Reason             string `gorm:"column:reason"`
			OperatorID         uint64 `gorm:"column:operator_admin_id"`
		}
		if e := tx.Table("qixi_crm_b_user_admin_command_audit").Where("action='profile_update' AND idempotency_key=?", in.IdempotencyKey).Take(&audit).Error; e == nil {
			if audit.UserID != userID || audit.RequestFingerprint != fingerprint || audit.Reason != in.Reason || audit.OperatorID != uint64(middleware.AdminID(c)) {
				return errUserAdminConflict
			}
			out = gin.H{"user_id": userID, "replayed": true}
			return nil
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		if e := tx.Table("qixi_crm_b_user").Where("id=?", userID).Update("nickname", in.Nickname).Error; e != nil {
			return e
		}
		if e := tx.Exec("INSERT INTO qixi_crm_b_user_profile (user_id,avatar_url,gender,bio,source_channel) VALUES (?,?,?,?, 'pc') ON DUPLICATE KEY UPDATE avatar_url=VALUES(avatar_url),gender=VALUES(gender),bio=VALUES(bio)", userID, in.AvatarURL, in.Gender, in.Bio).Error; e != nil {
			return e
		}
		if e := tx.Table("qixi_crm_b_user_admin_command_audit").Create(map[string]any{"action": "profile_update", "user_id": userID, "request_fingerprint": fingerprint, "reason": in.Reason, "operator_admin_id": middleware.AdminID(c), "idempotency_key": in.IdempotencyKey}).Error; e != nil {
			return e
		}
		out = gin.H{"user_id": userID, "replayed": false}
		return nil
	})
	writeUserAdminResult(c, err, out, "用户资料保存失败")
}

func (h *Handler) ResetPassword(c *gin.Context) {
	userID, ok := positiveID(c.Param("id"))
	var in resetUserPasswordInput
	if !ok || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "重置密码参数错误")
		return
	}
	in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if !validUserCommand(in.Reason, in.IdempotencyKey) || len(in.Password) < 12 || len(in.Password) > 72 {
		response.Fail(c, http.StatusBadRequest, "密码、原因或幂等键错误")
		return
	}
	out := gin.H{}
	err := h.business.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		// Password material is deliberately neither audited nor fingerprinted.
		// Lock the target before inspecting command history so the credential
		// update and auth-version revocation share one serialization boundary.
		var user struct {
			ID          uint64
			AuthVersion uint64 `gorm:"column:auth_version"`
		}
		if e := tx.Table("qixi_crm_b_user").Clauses(clause.Locking{Strength: "UPDATE"}).Select("id,auth_version").Where("id=?", userID).Take(&user).Error; e != nil {
			return e
		}
		var audit struct {
			UserID     uint64 `gorm:"column:user_id"`
			Reason     string `gorm:"column:reason"`
			OperatorID uint64 `gorm:"column:operator_admin_id"`
		}
		if e := tx.Table("qixi_crm_b_user_admin_command_audit").Where("action='password_reset' AND idempotency_key=?", in.IdempotencyKey).Take(&audit).Error; e == nil {
			// A password is never persisted as an audit fingerprint. Replaying a
			// credential command would therefore make a changed password
			// indistinguishable from a transport retry. Fail closed instead.
			if audit.UserID != userID || audit.Reason != in.Reason || audit.OperatorID != uint64(middleware.AdminID(c)) {
				return errUserAdminConflict
			}
			return errUserAdminConflict
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		hash, e := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if e != nil {
			return e
		}
		res := tx.Table("qixi_crm_b_user_identity").Where("user_id=? AND channel='pc'", userID).Update("credential_hash", string(hash))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errUserAdminNotFound
		}
		if e := tx.Table("qixi_crm_b_user").Where("id=? AND auth_version=?", userID, user.AuthVersion).Update("auth_version", user.AuthVersion+1).Error; e != nil {
			return e
		}
		if e := tx.Table("qixi_crm_b_user_admin_command_audit").Create(map[string]any{"action": "password_reset", "user_id": userID, "request_fingerprint": fingerprint("password_reset", strconv.FormatUint(userID, 10), in.Reason), "reason": in.Reason, "operator_admin_id": middleware.AdminID(c), "idempotency_key": in.IdempotencyKey}).Error; e != nil {
			return e
		}
		out = gin.H{"user_id": userID, "replayed": false}
		return nil
	})
	writeUserAdminResult(c, err, out, "重置密码失败")
}

func validUserCommand(reason, key string) bool {
	return len([]rune(reason)) >= 2 && len([]rune(reason)) <= 500 && len([]rune(key)) >= 8 && len([]rune(key)) <= 128
}
func validAvatarURL(value string) bool {
	return value == "" || strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "/demo/")
}
func fingerprint(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	return hex.EncodeToString(sum[:])
}

func isDuplicateKey(err error) bool {
	var mysqlErr *mysqlDriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
func writeUserAdminResult(c *gin.Context, err error, out gin.H, failure string) {
	switch {
	case err == nil:
		response.OK(c, out)
	case errors.Is(err, gorm.ErrRecordNotFound), errors.Is(err, errUserAdminNotFound):
		response.Fail(c, http.StatusNotFound, "用户或 PC 登录身份不存在")
	case errors.Is(err, errUserAdminConflict):
		response.Fail(c, http.StatusConflict, "幂等键与既有用户管理命令不一致")
	default:
		response.Fail(c, http.StatusInternalServerError, failure)
	}
}

type detailProfile struct {
	ID         uint64    `gorm:"column:id" json:"id"`
	Nickname   string    `gorm:"column:nickname" json:"nickname"`
	Mobile     string    `gorm:"column:mobile" json:"mobile"`
	Status     int       `gorm:"column:status" json:"status"`
	Balance    float64   `gorm:"column:balance" json:"balance"`
	Points     int64     `gorm:"column:points" json:"points"`
	Commission float64   `gorm:"column:commission" json:"commission"`
	LevelName  string    `gorm:"column:level_name" json:"level_name"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

type assetRow struct {
	ID            uint64    `gorm:"column:id" json:"id"`
	AssetType     string    `gorm:"column:asset_type" json:"asset_type"`
	Amount        float64   `gorm:"column:amount" json:"amount"`
	ReferenceType string    `gorm:"column:reference_type" json:"reference_type"`
	ReferenceID   string    `gorm:"column:reference_id" json:"reference_id"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

type membershipLogRow struct {
	ID                uint64    `gorm:"column:id" json:"id"`
	LevelName         string    `gorm:"column:level_name" json:"level_name"`
	PreviousLevelName string    `gorm:"column:previous_level_name" json:"previous_level_name"`
	ChangeType        string    `gorm:"column:change_type" json:"change_type"`
	Note              string    `gorm:"column:note" json:"note"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`
}

type signRow struct {
	ID             uint64    `gorm:"column:id" json:"id"`
	SignDate       string    `gorm:"column:sign_date" json:"sign_date"`
	Points         int64     `gorm:"column:points" json:"points"`
	ContinuousDays int       `gorm:"column:continuous_days" json:"continuous_days"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
}

type browseRow struct {
	ID        uint64    `gorm:"column:id" json:"id"`
	ProductID uint64    `gorm:"column:product_id" json:"product_id"`
	StoreID   uint64    `gorm:"column:store_id" json:"store_id"`
	StoreName string    `gorm:"column:store_name" json:"store_name"`
	Title     string    `gorm:"column:title" json:"title"`
	ViewedAt  time.Time `gorm:"column:viewed_at" json:"viewed_at"`
}

type orderRow struct {
	ID        uint64    `gorm:"column:id" json:"id"`
	OrderNo   string    `gorm:"column:order_no" json:"order_no"`
	StoreName string    `gorm:"column:store_name" json:"store_name"`
	PayAmount float64   `gorm:"column:pay_amount" json:"pay_amount"`
	Quantity  int       `gorm:"column:total_quantity" json:"total_quantity"`
	Status    string    `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

type couponRow struct {
	ID            uint64     `gorm:"column:id" json:"id"`
	CouponID      uint64     `gorm:"column:coupon_id" json:"coupon_id"`
	StoreID       uint64     `gorm:"column:store_id" json:"store_id"`
	Name          string     `gorm:"column:name" json:"name"`
	DiscountType  string     `gorm:"column:discount_type" json:"discount_type"`
	DiscountValue float64    `gorm:"column:discount_value" json:"discount_value"`
	MinAmount     float64    `gorm:"column:min_amount" json:"min_amount"`
	Status        string     `gorm:"column:status" json:"status"`
	ObtainedAt    time.Time  `gorm:"column:obtained_at" json:"obtained_at"`
	EndsAt        *time.Time `gorm:"column:ends_at" json:"ends_at"`
}

type distributionView struct {
	ParentUserID    uint64 `gorm:"column:parent_user_id" json:"parent_user_id"`
	ParentNickname  string `gorm:"column:parent_nickname" json:"parent_nickname"`
	DirectUserCount int64  `gorm:"column:direct_user_count" json:"direct_user_count"`
	PromoterStatus  int8   `gorm:"column:promoter_status" json:"promoter_status"`
}

type row struct {
	ID            uint64    `gorm:"column:id" json:"id"`
	Nickname      string    `gorm:"column:nickname" json:"nickname"`
	AvatarURL     string    `gorm:"column:avatar_url" json:"avatar_url"`
	Mobile        string    `gorm:"column:mobile" json:"mobile"`
	SourceChannel string    `gorm:"column:source_channel" json:"source_channel"`
	Status        int       `gorm:"column:status" json:"status"`
	Balance       float64   `gorm:"column:balance" json:"balance"`
	Points        int64     `gorm:"column:points" json:"points"`
	LevelName     string    `gorm:"column:level_name" json:"level_name"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

type userExportInput struct {
	ID      uint64 `json:"id"`
	Keyword string `json:"keyword"`
	Status  *int   `json:"status"`
	Reason  string `json:"reason"`
}

func (h *Handler) List(c *gin.Context) {
	page, limit := paging(c)
	status := strings.TrimSpace(c.Query("status"))
	if status != "" && status != "0" && status != "1" {
		response.Fail(c, http.StatusBadRequest, "用户状态错误")
		return
	}
	q := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_user AS u").
		Select("u.id,u.nickname,u.mobile,u.status,u.created_at,COALESCE(p.avatar_url,'') AS avatar_url,COALESCE(p.source_channel,'') AS source_channel,COALESCE(a.balance,0) AS balance,COALESCE(a.points,0) AS points,COALESCE(l.name,'') AS level_name").
		Joins("LEFT JOIN qixi_crm_b_user_profile AS p ON p.user_id=u.id").
		Joins("LEFT JOIN qixi_crm_b_member_account AS a ON a.user_id=u.id").
		Joins("LEFT JOIN qixi_crm_b_member_level AS l ON l.id=a.level_id")
	if status != "" {
		// status=1：启用且未注销（注销会把 status 置 0，见 auth.CancelAccount）
		q = q.Where("u.status=?", status)
	}
	if idRaw := strings.TrimSpace(c.Query("id")); idRaw != "" {
		id, err := strconv.ParseUint(idRaw, 10, 64)
		if err != nil || id == 0 {
			response.Fail(c, http.StatusBadRequest, "用户 ID 错误")
			return
		}
		q = q.Where("u.id=?", id)
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		if len([]rune(keyword)) > 64 {
			response.Fail(c, http.StatusBadRequest, "用户搜索词过长")
			return
		}
		like := "%" + keyword + "%"
		q = q.Where("u.nickname LIKE ? OR u.mobile LIKE ? OR CAST(u.id AS CHAR) = ?", like, like, keyword)
	}
	if nickname := strings.TrimSpace(c.Query("nickname")); nickname != "" {
		q = q.Where("u.nickname LIKE ?", "%"+nickname+"%")
	}
	if phone := strings.TrimSpace(c.Query("phone")); phone != "" {
		q = q.Where("u.mobile LIKE ?", "%"+phone+"%")
	}
	q = queryfilter.ApplyCreatedAtRange(q, c, "u.created_at")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, 500, "用户查询失败")
		return
	}
	var rows []row
	if err := q.Order("u.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, 500, "用户查询失败")
		return
	}
	for i := range rows {
		rows[i].Mobile = mask(rows[i].Mobile)
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

// Export produces a bounded, privacy-minimised CSV payload for a manual browser
// download. It deliberately excludes addresses, identity credentials, payment
// and withdrawal data, and records each export request before returning bytes.
func (h *Handler) Export(c *gin.Context) {
	var in userExportInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "用户导出参数错误")
		return
	}
	in.Keyword, in.Reason = strings.TrimSpace(in.Keyword), strings.TrimSpace(in.Reason)
	if (in.Status != nil && *in.Status != 0 && *in.Status != 1) || len([]rune(in.Keyword)) > 64 || len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 {
		response.Fail(c, http.StatusBadRequest, "导出筛选条件或原因错误")
		return
	}
	q := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_user AS u").
		Select("u.id,u.nickname,u.mobile,u.status,u.created_at,COALESCE(a.balance,0) AS balance,COALESCE(a.points,0) AS points,COALESCE(l.name,'') AS level_name").
		Joins("LEFT JOIN qixi_crm_b_member_account AS a ON a.user_id=u.id").Joins("LEFT JOIN qixi_crm_b_member_level AS l ON l.id=a.level_id")
	if in.Status != nil {
		q = q.Where("u.status=?", *in.Status)
	}
	if in.ID != 0 {
		q = q.Where("u.id=?", in.ID)
	}
	if in.Keyword != "" {
		q = q.Where("u.nickname LIKE ?", "%"+in.Keyword+"%")
	}
	// Read one additional row so `truncated` means there were records omitted,
	// rather than merely that the result happened to contain exactly 5000 rows.
	rows := make([]row, 0, 5001)
	if err := q.Order("u.id DESC").Limit(5001).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "用户导出查询失败")
		return
	}
	truncated := len(rows) > 5000
	if truncated {
		rows = rows[:5000]
	}
	var output bytes.Buffer
	output.Write([]byte{0xEF, 0xBB, 0xBF}) // Excel UTF-8 BOM; all content remains utf8mb4-compatible.
	w := csv.NewWriter(&output)
	if err := w.Write([]string{"用户ID", "昵称", "手机号（脱敏）", "状态", "余额", "积分", "会员等级", "注册时间"}); err != nil {
		response.Fail(c, http.StatusInternalServerError, "用户导出生成失败")
		return
	}
	for _, item := range rows {
		status := "停用"
		if item.Status == 1 {
			status = "启用"
		}
		level := item.LevelName
		if level == "" {
			level = "普通会员"
		}
		if err := w.Write([]string{strconv.FormatUint(item.ID, 10), csvCell(item.Nickname), mask(item.Mobile), status, strconv.FormatFloat(item.Balance, 'f', 2, 64), strconv.FormatInt(item.Points, 10), csvCell(level), item.CreatedAt.Format(time.RFC3339)}); err != nil {
			response.Fail(c, http.StatusInternalServerError, "用户导出生成失败")
			return
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		response.Fail(c, http.StatusInternalServerError, "用户导出生成失败")
		return
	}
	fingerprint := fingerprint("user_export", strconv.FormatUint(in.ID, 10), in.Keyword, strconv.Itoa(exportStatusValue(in.Status)), in.Reason)
	if err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_user_export_audit").Create(map[string]any{"query_fingerprint": fingerprint, "row_count": len(rows), "reason": in.Reason, "operator_admin_id": middleware.AdminID(c)}).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "用户导出审计写入失败")
		return
	}
	response.OK(c, gin.H{"file_name": "用户导出_" + time.Now().Format("20060102150405") + ".csv", "content": output.String(), "row_count": len(rows), "truncated": truncated})
}

func exportStatusValue(status *int) int {
	if status == nil {
		return -1
	}
	return *status
}

// csvCell protects spreadsheet applications from interpreting user-controlled
// nicknames or level names as formulas when a platform operator opens exports.
func csvCell(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "=") || strings.HasPrefix(value, "+") || strings.HasPrefix(value, "-") || strings.HasPrefix(value, "@") {
		return "'" + value
	}
	return value
}

// Detail exposes a privacy-minimised, read-only supervision view. It does not
// return delivery addresses, recipient snapshots, invoice profiles, payment
// transactions, or withdrawal account snapshots.
func (h *Handler) Detail(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "用户 ID 错误")
		return
	}
	profile := detailProfile{}
	err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_user AS u").
		Select("u.id,u.nickname,u.mobile,u.status,u.created_at,COALESCE(a.balance,0) AS balance,COALESCE(a.points,0) AS points,COALESCE(a.commission,0) AS commission,COALESCE(l.name,'') AS level_name").
		Joins("LEFT JOIN qixi_crm_b_member_account AS a ON a.user_id=u.id").
		Joins("LEFT JOIN qixi_crm_b_member_level AS l ON l.id=a.level_id").
		Where("u.id=?", id).Take(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "用户不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "用户详情查询失败")
		return
	}
	profile.Mobile = mask(profile.Mobile)

	assets := make([]assetRow, 0)
	membershipLogs := make([]membershipLogRow, 0)
	signs := make([]signRow, 0)
	browseHistory := make([]browseRow, 0)
	orders := make([]orderRow, 0)
	coupons := make([]couponRow, 0)
	distribution := distributionView{}
	ctx := c.Request.Context()
	if err := h.business.WithContext(ctx).Table("qixi_crm_b_asset_ledger").Select("id,asset_type,amount,reference_type,reference_id,created_at").Where("user_id=?", id).Order("id DESC").Limit(20).Scan(&assets).Error; err != nil ||
		h.business.WithContext(ctx).Table("qixi_crm_b_member_level_log AS log").
			Select("log.id,COALESCE(current_level.name,'普通会员') AS level_name,COALESCE(previous_level.name,'') AS previous_level_name,log.change_type,log.note,log.created_at").
			Joins("LEFT JOIN qixi_crm_b_member_level AS current_level ON current_level.id=log.level_id").
			Joins("LEFT JOIN qixi_crm_b_member_level AS previous_level ON previous_level.id=log.previous_level_id").
			Where("log.user_id=?", id).Order("log.id DESC").Limit(20).Scan(&membershipLogs).Error != nil ||
		h.business.WithContext(ctx).Table("qixi_crm_b_user_sign").
			Select("id,DATE_FORMAT(sign_date, '%Y-%m-%d') AS sign_date,points,continuous_days,created_at").
			Where("user_id=?", id).Order("sign_date DESC,id DESC").Limit(31).Scan(&signs).Error != nil ||
		h.business.WithContext(ctx).Table("qixi_crm_b_user_browse_history AS h").
			Select("h.id,h.product_id,h.store_id,COALESCE(p.store_name,'') AS store_name,COALESCE(p.title,'商品已下架') AS title,h.viewed_at").
			Joins("LEFT JOIN qixi_crm_b_product_view AS p ON p.product_id=h.product_id").
			Where("h.user_id=?", id).Order("h.viewed_at DESC,h.id DESC").Limit(20).Scan(&browseHistory).Error != nil ||
		h.business.WithContext(ctx).Table("qixi_crm_b_order").
			Select("id,order_no,store_name_snapshot AS store_name,pay_amount,total_quantity,status,created_at").
			Where("user_id=?", id).Order("id DESC").Limit(20).Scan(&orders).Error != nil ||
		h.business.WithContext(ctx).Table("qixi_crm_b_coupon_user AS u").
			Select("u.id,u.coupon_id,c.store_id,c.name,c.discount_type,c.discount_value,c.min_amount,u.status,u.obtained_at,c.ends_at").
			Joins("JOIN qixi_crm_b_coupon_template_view AS c ON c.coupon_id=u.coupon_id").
			Where("u.user_id=?", id).Order("u.id DESC").Limit(20).Scan(&coupons).Error != nil ||
		h.business.WithContext(ctx).Table("qixi_crm_b_user AS u").
			Select("COALESCE(rel.parent_user_id,0) AS parent_user_id,COALESCE(parent.nickname,'') AS parent_nickname,COALESCE(children.direct_user_count,0) AS direct_user_count,COALESCE(promoter.status,-1) AS promoter_status").
			Joins("LEFT JOIN qixi_crm_b_distribution_relation AS rel ON rel.user_id=u.id").
			Joins("LEFT JOIN qixi_crm_b_user AS parent ON parent.id=rel.parent_user_id").
			Joins("LEFT JOIN (SELECT parent_user_id,COUNT(*) AS direct_user_count FROM qixi_crm_b_distribution_relation WHERE parent_user_id IS NOT NULL GROUP BY parent_user_id) AS children ON children.parent_user_id=u.id").
			Joins("LEFT JOIN qixi_crm_b_distribution_promoter AS promoter ON promoter.user_id=u.id").
			Where("u.id=?", id).Scan(&distribution).Error != nil {
		response.Fail(c, http.StatusInternalServerError, "用户审计记录查询失败")
		return
	}
	response.OK(c, gin.H{"profile": profile, "assets": assets, "membership_logs": membershipLogs, "signs": signs, "browse_history": browseHistory, "orders": orders, "coupons": coupons, "distribution": distribution})
}
func paging(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if p < 1 {
		p = 1
	}
	if l < 1 || l > 100 {
		l = 20
	}
	return p, l
}
func mask(v string) string {
	r := []rune(strings.TrimSpace(v))
	if len(r) <= 4 {
		return ""
	}
	return string(r[:3]) + "****" + string(r[len(r)-4:])
}

func positiveID(raw string) (uint64, bool) {
	id, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 64)
	return id, err == nil && id > 0
}

func optionalPositiveQuery(c *gin.Context, key string) (uint64, bool) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return 0, true
	}
	return positiveID(raw)
}

type userStatusInput struct {
	Status         int    `json:"status"`
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

func (h *Handler) ChangeStatus(c *gin.Context) {
	userID, ok := positiveID(c.Param("id"))
	var in userStatusInput
	if !ok || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "用户状态调整参数错误")
		return
	}
	in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if (in.Status != 0 && in.Status != 1) || len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 || len([]rune(in.IdempotencyKey)) < 8 || len([]rune(in.IdempotencyKey)) > 128 {
		response.Fail(c, http.StatusBadRequest, "状态、原因或幂等键错误")
		return
	}
	out := gin.H{}
	err := h.business.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		// Lock the account before checking a prior command. This makes a replay
		// observe the same account state boundary as a first execution and keeps
		// the auth-version bump atomic with the status transition.
		var user struct {
			ID          uint64
			Status      int
			AuthVersion uint64 `gorm:"column:auth_version"`
		}
		if e := tx.Table("qixi_crm_b_user").Clauses(clause.Locking{Strength: "UPDATE"}).Select("id,status,auth_version").Where("id=?", userID).Take(&user).Error; e != nil {
			return e
		}
		var replay struct {
			UserID     uint64 `gorm:"column:user_id"`
			ToStatus   int    `gorm:"column:to_status"`
			Reason     string `gorm:"column:reason"`
			OperatorID uint64 `gorm:"column:operator_admin_id"`
		}
		if e := tx.Table("qixi_crm_b_user_status_command_audit").Where("idempotency_key=?", in.IdempotencyKey).Take(&replay).Error; e == nil {
			if replay.UserID != userID || replay.ToStatus != in.Status || replay.Reason != in.Reason || replay.OperatorID != uint64(middleware.AdminID(c)) {
				return errUserGroupConflict
			}
			out = gin.H{"user_id": userID, "status": in.Status, "replayed": true}
			return nil
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		if user.Status != in.Status {
			if e := tx.Table("qixi_crm_b_user").Where("id=? AND auth_version=?", userID, user.AuthVersion).Updates(map[string]any{"status": in.Status, "auth_version": user.AuthVersion + 1}).Error; e != nil {
				return e
			}
		}
		if e := tx.Table("qixi_crm_b_user_status_command_audit").Create(map[string]any{"user_id": userID, "from_status": user.Status, "to_status": in.Status, "reason": in.Reason, "operator_admin_id": middleware.AdminID(c), "idempotency_key": in.IdempotencyKey}).Error; e != nil {
			return e
		}
		out = gin.H{"user_id": userID, "previous_status": user.Status, "status": in.Status, "replayed": false}
		return nil
	})
	switch {
	case err == nil:
		response.OK(c, out)
	case errors.Is(err, gorm.ErrRecordNotFound):
		response.Fail(c, http.StatusNotFound, "用户不存在")
	case errors.Is(err, errUserGroupConflict):
		response.Fail(c, http.StatusConflict, "幂等键与既有用户状态调整不一致")
	default:
		response.Fail(c, http.StatusInternalServerError, "用户状态调整失败")
	}
}

type assetAdjustmentInput struct {
	AssetType      string      `json:"asset_type"`
	Amount         json.Number `json:"amount"`
	Reason         string      `json:"reason"`
	IdempotencyKey string      `json:"idempotency_key"`
}

type accountBalance struct {
	Balance string `gorm:"column:balance"`
	Points  int64  `gorm:"column:points"`
}

var (
	errAdjustmentNotFound     = errors.New("adjustment user not found")
	errAdjustmentConflict     = errors.New("adjustment conflict")
	errAdjustmentInsufficient = errors.New("adjustment insufficient balance")
)

// AdjustAsset creates one immutable manual adjustment ledger and matching
// audit row inside the same business-schema transaction. It is intentionally
// limited to balance and points; commission, payment and withdrawal flows use
// their own state machines and cannot be mutated through this endpoint.
func (h *Handler) AdjustAsset(c *gin.Context) {
	userID, ok := positiveID(c.Param("id"))
	var in assetAdjustmentInput
	if !ok || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "调账参数错误")
		return
	}
	in.AssetType, in.Reason, in.IdempotencyKey = strings.TrimSpace(in.AssetType), strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if !validAdjustment(in) {
		response.Fail(c, http.StatusBadRequest, "调账类型、金额、说明或幂等键错误")
		return
	}
	amountMinor, _ := adjustmentMinor(in)
	out := gin.H{}
	err := h.business.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var user struct{ ID uint64 }
		if err := tx.Table("qixi_crm_b_user").Select("id").Where("id=?", userID).Take(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errAdjustmentNotFound
			}
			return err
		}
		var previous struct {
			ID            uint64 `gorm:"column:id"`
			Amount        string `gorm:"column:amount"`
			ReferenceID   string `gorm:"column:reference_id"`
			ReferenceType string `gorm:"column:reference_type"`
		}
		if err := tx.Table("qixi_crm_b_asset_ledger").Where("asset_type=? AND idempotency_key=?", in.AssetType, in.IdempotencyKey).Take(&previous).Error; err == nil {
			previousMinor, validPreviousAmount := parseMinor(previous.Amount)
			if !validPreviousAmount || previous.ReferenceType != "platform_manual_adjustment" || previous.ReferenceID != strconv.FormatUint(userID, 10) || previousMinor != amountMinor {
				return errAdjustmentConflict
			}
			var audit struct {
				UserID        uint64 `gorm:"column:user_id"`
				Reason        string `gorm:"column:reason"`
				OperatorAdmin uint64 `gorm:"column:operator_admin_id"`
			}
			if err := tx.Table("qixi_crm_b_user_asset_adjustment_audit").Select("user_id,reason,operator_admin_id").Where("asset_type=? AND idempotency_key=?", in.AssetType, in.IdempotencyKey).Take(&audit).Error; err != nil {
				return err
			}
			if audit.UserID != userID || audit.Reason != in.Reason || audit.OperatorAdmin != uint64(middleware.AdminID(c)) {
				return errAdjustmentConflict
			}
			out = gin.H{"ledger_id": previous.ID, "replayed": true}
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := tx.Exec("INSERT INTO qixi_crm_b_member_account (user_id) VALUES (?) ON DUPLICATE KEY UPDATE user_id=VALUES(user_id)", userID).Error; err != nil {
			return err
		}
		var account accountBalance
		if err := tx.Table("qixi_crm_b_member_account").Clauses(clause.Locking{Strength: "UPDATE"}).Select("balance,points").Where("user_id=?", userID).Take(&account).Error; err != nil {
			return err
		}
		before, validBalance := parseMinor(account.Balance)
		if !validBalance {
			return errors.New("invalid member balance")
		}
		after := before
		if in.AssetType == "balance" {
			if (amountMinor > 0 && before > int64(^uint64(0)>>1)-amountMinor) || (amountMinor < 0 && before < -int64(^uint64(0)>>1)-amountMinor) {
				return errAdjustmentInsufficient
			}
			after = before + amountMinor
			if after < 0 {
				return errAdjustmentInsufficient
			}
			if err := tx.Table("qixi_crm_b_member_account").Where("user_id=?", userID).Update("balance", formatMinor(after)).Error; err != nil {
				return err
			}
		} else {
			before, after = account.Points, account.Points+amountMinor
			if (amountMinor > 0 && account.Points > int64(^uint64(0)>>1)-amountMinor) || (amountMinor < 0 && account.Points < -int64(^uint64(0)>>1)-amountMinor) || after < 0 {
				return errAdjustmentInsufficient
			}
			if err := tx.Table("qixi_crm_b_member_account").Where("user_id=?", userID).Update("points", after).Error; err != nil {
				return err
			}
		}
		amount := formatMinor(amountMinor)
		if in.AssetType == "points" {
			amount = strconv.FormatInt(amountMinor, 10)
		}
		ledger := map[string]any{"user_id": userID, "asset_type": in.AssetType, "amount": amount, "reference_type": "platform_manual_adjustment", "reference_id": strconv.FormatUint(userID, 10), "idempotency_key": in.IdempotencyKey}
		if err := tx.Table("qixi_crm_b_asset_ledger").Create(ledger).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_user_asset_adjustment_audit").Create(map[string]any{"user_id": userID, "asset_type": in.AssetType, "amount": amount, "balance_before": amountForAsset(in.AssetType, before), "balance_after": amountForAsset(in.AssetType, after), "reason": in.Reason, "operator_admin_id": middleware.AdminID(c), "idempotency_key": in.IdempotencyKey}).Error; err != nil {
			return err
		}
		var created struct{ ID uint64 }
		if err := tx.Table("qixi_crm_b_asset_ledger").Select("id").Where("asset_type=? AND idempotency_key=?", in.AssetType, in.IdempotencyKey).Take(&created).Error; err != nil {
			return err
		}
		out = gin.H{"ledger_id": created.ID, "asset_type": in.AssetType, "amount": amount, "balance_before": amountForAsset(in.AssetType, before), "balance_after": amountForAsset(in.AssetType, after), "replayed": false}
		return nil
	})
	switch {
	case err == nil:
		response.OK(c, out)
	case errors.Is(err, errAdjustmentNotFound):
		response.Fail(c, http.StatusNotFound, "用户不存在")
	case errors.Is(err, errAdjustmentConflict), errors.Is(err, errAdjustmentInsufficient):
		response.Fail(c, http.StatusConflict, map[error]string{errAdjustmentConflict: "幂等键与既有调账不一致", errAdjustmentInsufficient: "账户余额或积分不足"}[err])
	default:
		response.Fail(c, http.StatusInternalServerError, "用户调账失败")
	}
}

func validAdjustment(in assetAdjustmentInput) bool {
	amount, validAmount := adjustmentMinor(in)
	if (in.AssetType != "balance" && in.AssetType != "points") || !validAmount || amount == 0 || amount < -100_000_000 || amount > 100_000_000 || len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 || len([]rune(in.IdempotencyKey)) < 8 || len([]rune(in.IdempotencyKey)) > 128 {
		return false
	}
	return true
}

// adjustmentMinor converts a JSON money value to integer cents. Balance is
// never calculated with float64: decimal values beyond two fraction digits,
// scientific notation and integer overflow are rejected before any mutation.
// Points use the same integer transport with an implicit zero fraction.
func adjustmentMinor(in assetAdjustmentInput) (int64, bool) {
	raw := strings.TrimSpace(in.Amount.String())
	if raw == "" || strings.ContainsAny(raw, "eE") {
		return 0, false
	}
	if in.AssetType == "points" {
		value, err := strconv.ParseInt(raw, 10, 64)
		return value, err == nil
	}
	return parseMinor(raw)
}

func parseMinor(raw string) (int64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, false
	}
	sign := int64(1)
	if raw[0] == '-' || raw[0] == '+' {
		if raw[0] == '-' {
			sign = -1
		}
		raw = raw[1:]
	}
	parts := strings.Split(raw, ".")
	if len(parts) > 2 || parts[0] == "" || len(parts[1:]) == 1 && len(parts[1]) > 2 {
		return 0, false
	}
	whole, fraction := parts[0], ""
	if len(parts) == 2 {
		fraction = parts[1]
	}
	if strings.Trim(whole, "0123456789") != "" || strings.Trim(fraction, "0123456789") != "" {
		return 0, false
	}
	for len(fraction) < 2 {
		fraction += "0"
	}
	value, err := strconv.ParseInt(whole+fraction, 10, 64)
	if err != nil {
		return 0, false
	}
	return sign * value, true
}

func formatMinor(value int64) string {
	sign := ""
	if value < 0 {
		sign, value = "-", -value
	}
	return sign + strconv.FormatInt(value/100, 10) + "." + fmt.Sprintf("%02d", value%100)
}

func amountForAsset(assetType string, value int64) string {
	if assetType == "points" {
		return strconv.FormatInt(value, 10)
	}
	return formatMinor(value)
}

type memberLevelAdjustmentInput struct {
	LevelID        uint64 `json:"level_id"`
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

var (
	errMemberLevelNotFound = errors.New("member level user or target not found")
	errMemberLevelConflict = errors.New("member level idempotency conflict")
)

// AdjustMemberLevel changes only the current membership projection. It never
// grants SVIP, coupons, balance or order discounts directly; those domains
// keep their own eligibility and pricing state machines.
func (h *Handler) AdjustMemberLevel(c *gin.Context) {
	userID, ok := positiveID(c.Param("id"))
	var in memberLevelAdjustmentInput
	if !ok || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "会员等级调整参数错误")
		return
	}
	in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 || len([]rune(in.IdempotencyKey)) < 8 || len([]rune(in.IdempotencyKey)) > 128 {
		response.Fail(c, http.StatusBadRequest, "调整原因或幂等键错误")
		return
	}
	out := gin.H{}
	err := h.business.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var user struct{ ID uint64 }
		if err := tx.Table("qixi_crm_b_user").Select("id").Where("id=?", userID).Take(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errMemberLevelNotFound
			}
			return err
		}
		if in.LevelID != 0 {
			var target memberLevel
			if err := tx.Table("qixi_crm_b_member_level AS l").Select("l.id,l.name,l.rank,l.status").Where("l.id=? AND l.status=1 AND l.deleted_at IS NULL", in.LevelID).Take(&target).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errMemberLevelNotFound
				}
				return err
			}
		}
		var replay struct {
			ID      uint64
			LevelID sql.NullInt64 `gorm:"column:level_id"`
		}
		if err := tx.Table("qixi_crm_b_member_level_log").Where("user_id=? AND idempotency_key=?", userID, in.IdempotencyKey).Take(&replay).Error; err == nil {
			if uint64(replay.LevelID.Int64) != in.LevelID || (in.LevelID == 0 && replay.LevelID.Valid) {
				return errMemberLevelConflict
			}
			out = gin.H{"log_id": replay.ID, "user_id": userID, "level_id": in.LevelID, "replayed": true}
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := tx.Exec("INSERT INTO qixi_crm_b_member_account (user_id) VALUES (?) ON DUPLICATE KEY UPDATE user_id=VALUES(user_id)", userID).Error; err != nil {
			return err
		}
		var account struct {
			LevelID sql.NullInt64 `gorm:"column:level_id"`
		}
		if err := tx.Table("qixi_crm_b_member_account").Clauses(clause.Locking{Strength: "UPDATE"}).Select("level_id").Where("user_id=?", userID).Take(&account).Error; err != nil {
			return err
		}
		previous := account.LevelID
		if err := tx.Table("qixi_crm_b_member_account").Where("user_id=?", userID).Update("level_id", nullableLevel(in.LevelID)).Error; err != nil {
			return err
		}
		log := map[string]any{"user_id": userID, "level_id": nullableLevel(in.LevelID), "previous_level_id": nullableLevelFromSQL(previous), "change_type": "manual", "note": in.Reason, "idempotency_key": in.IdempotencyKey, "operator_admin_id": middleware.AdminID(c)}
		if err := tx.Table("qixi_crm_b_member_level_log").Create(log).Error; err != nil {
			return err
		}
		var created struct{ ID uint64 }
		if err := tx.Table("qixi_crm_b_member_level_log").Select("id").Where("user_id=? AND idempotency_key=?", userID, in.IdempotencyKey).Take(&created).Error; err != nil {
			return err
		}
		out = gin.H{"log_id": created.ID, "user_id": userID, "level_id": in.LevelID, "previous_level_id": nullableLevelFromSQL(previous), "replayed": false}
		return nil
	})
	switch {
	case err == nil:
		response.OK(c, out)
	case errors.Is(err, errMemberLevelNotFound):
		response.Fail(c, http.StatusNotFound, "用户或会员等级不存在")
	case errors.Is(err, errMemberLevelConflict):
		response.Fail(c, http.StatusConflict, "幂等键与既有会员等级调整不一致")
	default:
		response.Fail(c, http.StatusInternalServerError, "会员等级调整失败")
	}
}

func nullableLevel(levelID uint64) any {
	if levelID == 0 {
		return nil
	}
	return levelID
}
func nullableLevelFromSQL(level sql.NullInt64) any {
	if !level.Valid {
		return nil
	}
	return level.Int64
}

type userGroupCommandInput struct {
	GroupID        uint64   `json:"group_id"`
	UserIDs        []uint64 `json:"user_ids"`
	Reason         string   `json:"reason"`
	IdempotencyKey string   `json:"idempotency_key"`
}

var (
	errUserGroupNotFound = errors.New("user or group not found")
	errUserGroupConflict = errors.New("user group command conflict")
)

// ChangeUserGroup changes one user's CRMEB-compatible operating group. A zero
// group_id means ungrouped; it never changes the user's membership, orders,
// assets, commission or login state.
func (h *Handler) ChangeUserGroup(c *gin.Context) {
	userID, ok := positiveID(c.Param("id"))
	var in userGroupCommandInput
	if !ok || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "用户分组调整参数错误")
		return
	}
	in.UserIDs = []uint64{userID}
	h.changeUserGroups(c, in)
}

func (h *Handler) ChangeUserGroups(c *gin.Context) {
	var in userGroupCommandInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "批量用户分组调整参数错误")
		return
	}
	h.changeUserGroups(c, in)
}

func (h *Handler) changeUserGroups(c *gin.Context, in userGroupCommandInput) {
	userIDs, ok := normalizeUserIDs(in.UserIDs)
	in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if !ok || len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 || len([]rune(in.IdempotencyKey)) < 8 || len([]rune(in.IdempotencyKey)) > 128 {
		response.Fail(c, http.StatusBadRequest, "用户、原因或幂等键错误")
		return
	}
	userIDsJSON, _ := json.Marshal(userIDs)
	userIDsText := string(userIDsJSON)
	out := gin.H{}
	err := h.business.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		// User IDs are normalized in ascending order. Lock them before checking the
		// audit row so overlapping batch commands and identical retries observe a
		// stable user set and cannot race into the idempotency unique index.
		var lockedUsers []struct{ ID uint64 }
		if e := tx.Table("qixi_crm_b_user").Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").Where("id IN ?", userIDs).Order("id ASC").Find(&lockedUsers).Error; e != nil {
			return e
		}
		if len(lockedUsers) != len(userIDs) {
			return errUserGroupNotFound
		}
		if in.GroupID != 0 {
			var group struct {
				ID uint64 `gorm:"column:group_id"`
			}
			if e := tx.Table("qixi_crm_b_user_group").Clauses(clause.Locking{Strength: "UPDATE"}).Select("group_id").Where("group_id=? AND is_del=0", in.GroupID).Take(&group).Error; e != nil {
				if errors.Is(e, gorm.ErrRecordNotFound) {
					return errUserGroupNotFound
				}
				return e
			}
		}
		var replay struct {
			UserIDsJSON string `gorm:"column:user_ids_json"`
			GroupID     uint64 `gorm:"column:group_id"`
			Reason      string `gorm:"column:reason"`
			OperatorID  uint64 `gorm:"column:operator_admin_id"`
		}
		if e := tx.Table("qixi_crm_b_user_group_command_audit").Where("idempotency_key=?", in.IdempotencyKey).Take(&replay).Error; e == nil {
			if !sameUserIDsJSON(replay.UserIDsJSON, userIDs) || replay.GroupID != in.GroupID || replay.Reason != in.Reason || replay.OperatorID != uint64(middleware.AdminID(c)) {
				return errUserGroupConflict
			}
			out = gin.H{"user_ids": userIDs, "group_id": in.GroupID, "replayed": true}
			return nil
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		if e := tx.Table("qixi_crm_b_user").Where("id IN ?", userIDs).Update("group_id", in.GroupID).Error; e != nil {
			return e
		}
		if e := tx.Table("qixi_crm_b_user_group_command_audit").Create(map[string]any{
			"user_ids_json": userIDsText, "group_id": in.GroupID, "reason": in.Reason,
			"operator_admin_id": middleware.AdminID(c), "idempotency_key": in.IdempotencyKey,
		}).Error; e != nil {
			return e
		}
		out = gin.H{"user_ids": userIDs, "group_id": in.GroupID, "replayed": false}
		return nil
	})
	switch {
	case err == nil:
		response.OK(c, out)
	case errors.Is(err, errUserGroupNotFound):
		response.Fail(c, http.StatusNotFound, "用户或用户分组不存在")
	case errors.Is(err, errUserGroupConflict):
		response.Fail(c, http.StatusConflict, "幂等键与既有用户分组调整不一致")
	default:
		response.Fail(c, http.StatusInternalServerError, "用户分组调整失败")
	}
}

func normalizeUserIDs(ids []uint64) ([]uint64, bool) {
	if len(ids) == 0 || len(ids) > 100 {
		return nil, false
	}
	seen := make(map[uint64]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			return nil, false
		}
		seen[id] = struct{}{}
	}
	out := make([]uint64, 0, len(seen))
	for id := range seen {
		out = append(out, id)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out, true
}

func sameUserIDsJSON(raw string, want []uint64) bool {
	var got []uint64
	if json.Unmarshal([]byte(raw), &got) != nil {
		return false
	}
	got, ok := normalizeUserIDs(got)
	if !ok || len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

type userLabelCommandInput struct {
	LabelIDs       []uint64 `json:"label_ids"`
	UserIDs        []uint64 `json:"user_ids"`
	Reason         string   `json:"reason"`
	IdempotencyKey string   `json:"idempotency_key"`
}

var (
	errUserLabelNotFound = errors.New("user or label not found")
	errUserLabelConflict = errors.New("user label command conflict")
)

func (h *Handler) ChangeUserLabels(c *gin.Context) {
	userID, ok := positiveID(c.Param("id"))
	var in userLabelCommandInput
	if !ok || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "用户标签调整参数错误")
		return
	}
	in.UserIDs = []uint64{userID}
	h.changeUsersLabels(c, in)
}

func (h *Handler) ChangeUsersLabels(c *gin.Context) {
	var in userLabelCommandInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "批量用户标签调整参数错误")
		return
	}
	h.changeUsersLabels(c, in)
}

// changeUsersLabels replaces only operating labels. An empty label_ids list
// explicitly clears labels; it never changes any identity, order or asset fact.
func (h *Handler) changeUsersLabels(c *gin.Context, in userLabelCommandInput) {
	userIDs, userOK := normalizeUserIDs(in.UserIDs)
	labelIDs, labelOK := normalizeLabelIDs(in.LabelIDs)
	in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if !userOK || !labelOK || len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 || len([]rune(in.IdempotencyKey)) < 8 || len([]rune(in.IdempotencyKey)) > 128 {
		response.Fail(c, http.StatusBadRequest, "用户、标签、原因或幂等键错误")
		return
	}
	userIDsJSON, _ := json.Marshal(userIDs)
	labelIDsJSON, _ := json.Marshal(labelIDs)
	out := gin.H{}
	err := h.business.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		// Both sets are normalized in ascending order. Lock the affected users and
		// active labels before reading the audit row so replacement, clearing and
		// retries cannot observe stale membership or race into the unique audit key.
		var lockedUsers []struct{ ID uint64 }
		if e := tx.Table("qixi_crm_b_user").Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").Where("id IN ?", userIDs).Order("id ASC").Find(&lockedUsers).Error; e != nil {
			return e
		}
		if len(lockedUsers) != len(userIDs) {
			return errUserLabelNotFound
		}
		if len(labelIDs) > 0 {
			var lockedLabels []struct {
				ID uint64 `gorm:"column:label_id"`
			}
			if e := tx.Table("qixi_crm_b_user_label").Clauses(clause.Locking{Strength: "UPDATE"}).Select("label_id").Where("label_id IN ? AND is_del=0", labelIDs).Order("label_id ASC").Find(&lockedLabels).Error; e != nil {
				return e
			}
			if len(lockedLabels) != len(labelIDs) {
				return errUserLabelNotFound
			}
		}
		var replay struct {
			UserIDsJSON  string `gorm:"column:user_ids_json"`
			LabelIDsJSON string `gorm:"column:label_ids_json"`
			Reason       string `gorm:"column:reason"`
			OperatorID   uint64 `gorm:"column:operator_admin_id"`
		}
		if e := tx.Table("qixi_crm_b_user_label_command_audit").Where("idempotency_key=?", in.IdempotencyKey).Take(&replay).Error; e == nil {
			if !sameUserIDsJSON(replay.UserIDsJSON, userIDs) || !sameIDsJSON(replay.LabelIDsJSON, labelIDs) || replay.Reason != in.Reason || replay.OperatorID != uint64(middleware.AdminID(c)) {
				return errUserLabelConflict
			}
			out = gin.H{"user_ids": userIDs, "label_ids": labelIDs, "replayed": true}
			return nil
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		if e := tx.Exec("DELETE FROM qixi_crm_b_user_label_relation WHERE uid IN ?", userIDs).Error; e != nil {
			return e
		}
		for _, userID := range userIDs {
			for _, labelID := range labelIDs {
				if e := tx.Table("qixi_crm_b_user_label_relation").Create(map[string]any{"uid": userID, "label_id": labelID}).Error; e != nil {
					return e
				}
			}
		}
		if e := tx.Table("qixi_crm_b_user_label_command_audit").Create(map[string]any{
			"user_ids_json": string(userIDsJSON), "label_ids_json": string(labelIDsJSON), "reason": in.Reason,
			"operator_admin_id": middleware.AdminID(c), "idempotency_key": in.IdempotencyKey,
		}).Error; e != nil {
			return e
		}
		out = gin.H{"user_ids": userIDs, "label_ids": labelIDs, "replayed": false}
		return nil
	})
	switch {
	case err == nil:
		response.OK(c, out)
	case errors.Is(err, errUserLabelNotFound):
		response.Fail(c, http.StatusNotFound, "用户或用户标签不存在")
	case errors.Is(err, errUserLabelConflict):
		response.Fail(c, http.StatusConflict, "幂等键与既有用户标签调整不一致")
	default:
		response.Fail(c, http.StatusInternalServerError, "用户标签调整失败")
	}
}

func normalizeLabelIDs(ids []uint64) ([]uint64, bool) {
	if len(ids) > 50 {
		return nil, false
	}
	if len(ids) == 0 {
		return []uint64{}, true
	}
	seen := make(map[uint64]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			return nil, false
		}
		seen[id] = struct{}{}
	}
	out := make([]uint64, 0, len(seen))
	for id := range seen {
		out = append(out, id)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out, true
}

func sameIDsJSON(raw string, want []uint64) bool {
	var got []uint64
	if json.Unmarshal([]byte(raw), &got) != nil {
		return false
	}
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

type promoterCommandInput struct {
	UserIDs        []uint64 `json:"user_ids"`
	Status         int      `json:"status"`
	Reason         string   `json:"reason"`
	IdempotencyKey string   `json:"idempotency_key"`
}

func (h *Handler) ChangePromoters(c *gin.Context) {
	var in promoterCommandInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "批量推广员设置参数错误")
		return
	}
	userIDs, ok := normalizeUserIDs(in.UserIDs)
	in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if !ok || (in.Status != 0 && in.Status != 1) || !validUserCommand(in.Reason, in.IdempotencyKey) {
		response.Fail(c, http.StatusBadRequest, "用户、资格状态、原因或幂等键错误")
		return
	}
	userIDsJSON, _ := json.Marshal(userIDs)
	out := gin.H{}
	err := h.business.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var lockedUsers []struct{ ID uint64 }
		if e := tx.Table("qixi_crm_b_user").Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id").Where("id IN ?", userIDs).Order("id ASC").Find(&lockedUsers).Error; e != nil {
			return e
		}
		if len(lockedUsers) != len(userIDs) {
			return errUserAdminNotFound
		}
		var replay struct {
			UserIDsJSON string `gorm:"column:user_ids_json"`
			Status      int    `gorm:"column:status"`
			Reason      string `gorm:"column:reason"`
			OperatorID  uint64 `gorm:"column:operator_admin_id"`
		}
		if e := tx.Table("qixi_crm_b_distribution_promoter_command_audit").Where("idempotency_key=?", in.IdempotencyKey).Take(&replay).Error; e == nil {
			if !sameUserIDsJSON(replay.UserIDsJSON, userIDs) || replay.Status != in.Status || replay.Reason != in.Reason || replay.OperatorID != uint64(middleware.AdminID(c)) {
				return errUserAdminConflict
			}
			out = gin.H{"user_ids": userIDs, "status": in.Status, "replayed": true}
			return nil
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		for _, userID := range userIDs {
			if e := tx.Exec("INSERT INTO qixi_crm_b_distribution_promoter (user_id,status) VALUES (?,?) ON DUPLICATE KEY UPDATE status=VALUES(status)", userID, in.Status).Error; e != nil {
				return e
			}
		}
		if e := tx.Table("qixi_crm_b_distribution_promoter_command_audit").Create(map[string]any{"user_ids_json": string(userIDsJSON), "status": in.Status, "reason": in.Reason, "operator_admin_id": middleware.AdminID(c), "idempotency_key": in.IdempotencyKey}).Error; e != nil {
			return e
		}
		out = gin.H{"user_ids": userIDs, "status": in.Status, "replayed": false}
		return nil
	})
	switch {
	case err == nil:
		response.OK(c, out)
	case errors.Is(err, errUserAdminNotFound):
		response.Fail(c, http.StatusNotFound, "用户不存在")
	case errors.Is(err, errUserAdminConflict):
		response.Fail(c, http.StatusConflict, "幂等键与既有推广员设置不一致")
	default:
		response.Fail(c, http.StatusInternalServerError, "批量推广员设置失败")
	}
}

type inAppNotificationInput struct {
	Title          string `json:"title"`
	Body           string `json:"body"`
	CoverURL       string `json:"cover_url"`
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

// SendInAppNotification is the credential-free replacement for CRMEB's
// user-news action. It creates a C-end in-app image-text notification only;
// it never impersonates or calls WeChat/SMS providers.
func (h *Handler) SendInAppNotification(c *gin.Context) {
	userID, ok := positiveID(c.Param("id"))
	var in inAppNotificationInput
	if !ok || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "站内图文通知参数错误")
		return
	}
	in.Title, in.Body, in.CoverURL, in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Title), strings.TrimSpace(in.Body), strings.TrimSpace(in.CoverURL), strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if len([]rune(in.Title)) < 1 || len([]rune(in.Title)) > 255 || len([]rune(in.Body)) < 1 || len([]rune(in.Body)) > 2000 || !validAvatarURL(in.CoverURL) || !validUserCommand(in.Reason, in.IdempotencyKey) {
		response.Fail(c, http.StatusBadRequest, "通知内容、封面、原因或幂等键错误")
		return
	}
	requestFingerprint := fingerprint("in_app_notification", strconv.FormatUint(userID, 10), in.Title, in.Body, in.CoverURL, in.Reason)
	out := gin.H{}
	err := h.business.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var user struct{ ID uint64 }
		if e := tx.Table("qixi_crm_b_user").Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").Where("id=?", userID).Take(&user).Error; e != nil {
			return e
		}
		var replay struct {
			NotificationID     uint64 `gorm:"column:notification_id"`
			RequestFingerprint string `gorm:"column:request_fingerprint"`
			OperatorID         uint64 `gorm:"column:operator_admin_id"`
		}
		if e := tx.Table("qixi_crm_b_user_notification_command_audit").Where("user_id=? AND idempotency_key=?", userID, in.IdempotencyKey).Take(&replay).Error; e == nil {
			if replay.RequestFingerprint != requestFingerprint || replay.OperatorID != uint64(middleware.AdminID(c)) {
				return errUserAdminConflict
			}
			out = gin.H{"user_id": userID, "notification_id": replay.NotificationID, "replayed": true}
			return nil
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		payload, _ := json.Marshal(gin.H{"kind": "image_text", "cover_url": in.CoverURL, "source": "platform_manual"})
		note := struct {
			ID       uint64 `gorm:"column:id"`
			UserID   uint64 `gorm:"column:user_id"`
			Category string `gorm:"column:category"`
			Title    string `gorm:"column:title"`
			Body     string `gorm:"column:body"`
			Payload  string `gorm:"column:payload"`
		}{UserID: userID, Category: "system", Title: in.Title, Body: in.Body, Payload: string(payload)}
		if e := tx.Table("qixi_crm_b_notification").Create(&note).Error; e != nil {
			return e
		}
		if note.ID == 0 {
			return errors.New("notification id unavailable")
		}
		if e := tx.Table("qixi_crm_b_user_notification_command_audit").Create(map[string]any{"user_id": userID, "notification_id": note.ID, "title": in.Title, "request_fingerprint": requestFingerprint, "reason": in.Reason, "operator_admin_id": middleware.AdminID(c), "idempotency_key": in.IdempotencyKey}).Error; e != nil {
			return e
		}
		out = gin.H{"user_id": userID, "notification_id": note.ID, "replayed": false}
		return nil
	})
	writeUserAdminResult(c, err, out, "站内图文通知发送失败")
}

type couponCommandInput struct {
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

type userCouponState struct {
	ID          uint64  `gorm:"column:id"`
	CouponID    uint64  `gorm:"column:coupon_id"`
	UserID      uint64  `gorm:"column:user_id"`
	Status      string  `gorm:"column:status"`
	UsedOrderID *uint64 `gorm:"column:used_order_id"`
}

var (
	errCouponCommandNotFound = errors.New("coupon command user or coupon not found")
	errCouponCommandConflict = errors.New("coupon command conflict")
)

func (h *Handler) IssueCoupon(c *gin.Context)  { h.commandCoupon(c, "issue") }
func (h *Handler) RevokeCoupon(c *gin.Context) { h.commandCoupon(c, "revoke") }

// commandCoupon only creates a user coupon once, or revokes an unused coupon.
// It must not resurrect, delete, lock, or consume coupons: checkout and payment
// own those state transitions.
func (h *Handler) commandCoupon(c *gin.Context, action string) {
	userID, userOK := positiveID(c.Param("id"))
	couponID, couponOK := positiveID(c.Param("coupon_id"))
	var in couponCommandInput
	if !userOK || !couponOK || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "用户券操作参数错误")
		return
	}
	in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 || len([]rune(in.IdempotencyKey)) < 8 || len([]rune(in.IdempotencyKey)) > 128 {
		response.Fail(c, http.StatusBadRequest, "操作原因或幂等键错误")
		return
	}
	out := gin.H{}
	err := h.business.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		// Serialise user-scoped manual coupon commands before checking an existing
		// audit row. A second identical concurrent request observes the first
		// transaction's audit after acquiring this lock and becomes a replay rather
		// than racing into the coupon unique index.
		var user struct{ ID uint64 }
		if err := tx.Table("qixi_crm_b_user").Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").Where("id=?", userID).Take(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errCouponCommandNotFound
			}
			return err
		}
		var replay struct {
			ID           uint64 `gorm:"column:id"`
			UserID       uint64 `gorm:"column:user_id"`
			CouponID     uint64 `gorm:"column:coupon_id"`
			CouponUserID uint64 `gorm:"column:coupon_user_id"`
			Reason       string `gorm:"column:reason"`
			OperatorID   uint64 `gorm:"column:operator_admin_id"`
		}
		if err := tx.Table("qixi_crm_b_user_coupon_command_audit").Where("action=? AND idempotency_key=?", action, in.IdempotencyKey).Take(&replay).Error; err == nil {
			if replay.UserID != userID || replay.CouponID != couponID || replay.Reason != in.Reason || replay.OperatorID != uint64(middleware.AdminID(c)) {
				return errCouponCommandConflict
			}
			out = gin.H{"coupon_user_id": replay.CouponUserID, "user_id": userID, "coupon_id": couponID, "action": action, "replayed": true}
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		var coupon userCouponState
		if action == "issue" {
			now := time.Now()
			var count int64
			if err := tx.Table("qixi_crm_b_coupon_template_view").Where("coupon_id=? AND status=1 AND (starts_at IS NULL OR starts_at<=?) AND (ends_at IS NULL OR ends_at>=?)", couponID, now, now).Count(&count).Error; err != nil {
				return err
			}
			if count != 1 {
				return errCouponCommandNotFound
			}
			if err := tx.Table("qixi_crm_b_coupon_user").Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id=? AND coupon_id=?", userID, couponID).Take(&coupon).Error; err == nil {
				return errCouponCommandConflict
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if err := tx.Table("qixi_crm_b_coupon_user").Create(map[string]any{"user_id": userID, "coupon_id": couponID, "source": "platform_manual", "status": "unused"}).Error; err != nil {
				return err
			}
			if err := tx.Table("qixi_crm_b_coupon_user").Where("user_id=? AND coupon_id=?", userID, couponID).Take(&coupon).Error; err != nil {
				return err
			}
			if err := tx.Table("qixi_crm_b_user_coupon_command_audit").Create(map[string]any{"user_id": userID, "coupon_id": couponID, "coupon_user_id": coupon.ID, "action": action, "from_status": "", "to_status": "unused", "reason": in.Reason, "operator_admin_id": middleware.AdminID(c), "idempotency_key": in.IdempotencyKey}).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Table("qixi_crm_b_coupon_user").Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id=? AND coupon_id=?", userID, couponID).Take(&coupon).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errCouponCommandNotFound
				}
				return err
			}
			if coupon.Status != "unused" || coupon.UsedOrderID != nil {
				return errCouponCommandConflict
			}
			res := tx.Table("qixi_crm_b_coupon_user").Where("id=? AND status='unused' AND used_order_id IS NULL", coupon.ID).Update("status", "expired")
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected != 1 {
				return errCouponCommandConflict
			}
			if err := tx.Table("qixi_crm_b_user_coupon_command_audit").Create(map[string]any{"user_id": userID, "coupon_id": couponID, "coupon_user_id": coupon.ID, "action": action, "from_status": "unused", "to_status": "expired", "reason": in.Reason, "operator_admin_id": middleware.AdminID(c), "idempotency_key": in.IdempotencyKey}).Error; err != nil {
				return err
			}
		}
		out = gin.H{"coupon_user_id": coupon.ID, "user_id": userID, "coupon_id": couponID, "action": action, "replayed": false}
		return nil
	})
	switch {
	case err == nil:
		response.OK(c, out)
	case errors.Is(err, errCouponCommandNotFound):
		response.Fail(c, http.StatusNotFound, "用户或可用优惠券不存在")
	case errors.Is(err, errCouponCommandConflict):
		response.Fail(c, http.StatusConflict, "用户券状态不允许操作或幂等键冲突")
	default:
		response.Fail(c, http.StatusInternalServerError, "用户券操作失败")
	}
}

type referrerInput struct {
	ParentUserID   uint64 `json:"parent_user_id"`
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

var (
	errReferrerNotFound = errors.New("referrer user not found")
	errReferrerConflict = errors.New("referrer change conflict")
)

func (h *Handler) ChangeReferrer(c *gin.Context) {
	userID, ok := positiveID(c.Param("id"))
	var in referrerInput
	if !ok || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "推荐关系参数错误")
		return
	}
	in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 || len([]rune(in.IdempotencyKey)) < 8 || len([]rune(in.IdempotencyKey)) > 128 || in.ParentUserID == userID {
		response.Fail(c, http.StatusBadRequest, "上级用户、原因或幂等键错误")
		return
	}
	out := gin.H{}
	err := h.business.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		// Lock both ends in primary-key order before reading a replay or walking the
		// referral chain. This serialises identical commands and opposing A→B/B→A
		// commands so a concurrent cycle cannot pass two stale chain checks.
		lockedIDs := []uint64{userID}
		if in.ParentUserID != 0 {
			lockedIDs = append(lockedIDs, in.ParentUserID)
		}
		var lockedUsers []struct{ ID uint64 }
		if e := tx.Table("qixi_crm_b_user").Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").Where("id IN ?", lockedIDs).Order("id ASC").Find(&lockedUsers).Error; e != nil {
			return e
		}
		if len(lockedUsers) != len(lockedIDs) {
			return errReferrerNotFound
		}
		var replay struct {
			ID           uint64        `gorm:"column:id"`
			ParentUserID sql.NullInt64 `gorm:"column:parent_user_id"`
			Reason       string        `gorm:"column:reason"`
			OperatorID   uint64        `gorm:"column:operator_admin_id"`
		}
		if e := tx.Table("qixi_crm_b_distribution_relation_audit").Where("user_id=? AND idempotency_key=?", userID, in.IdempotencyKey).Take(&replay).Error; e == nil {
			if uint64(replay.ParentUserID.Int64) != in.ParentUserID || (in.ParentUserID == 0 && replay.ParentUserID.Valid) || replay.Reason != in.Reason || replay.OperatorID != uint64(middleware.AdminID(c)) {
				return errReferrerConflict
			}
			out = gin.H{"user_id": userID, "parent_user_id": in.ParentUserID, "replayed": true}
			return nil
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		if in.ParentUserID != 0 {
			if e := h.ensureNoReferralCycle(tx, userID, in.ParentUserID); e != nil {
				return e
			}
		}
		if e := tx.Exec("INSERT INTO qixi_crm_b_distribution_relation (user_id,parent_user_id,bound_at) VALUES (?,NULL,NOW()) ON DUPLICATE KEY UPDATE user_id=VALUES(user_id)", userID).Error; e != nil {
			return e
		}
		var relation struct {
			ParentUserID sql.NullInt64 `gorm:"column:parent_user_id"`
		}
		if e := tx.Table("qixi_crm_b_distribution_relation").Clauses(clause.Locking{Strength: "UPDATE"}).Select("parent_user_id").Where("user_id=?", userID).Take(&relation).Error; e != nil {
			return e
		}
		previous := relation.ParentUserID
		if e := tx.Table("qixi_crm_b_distribution_relation").Where("user_id=?", userID).Updates(map[string]any{"parent_user_id": nullableLevel(in.ParentUserID), "bound_at": time.Now()}).Error; e != nil {
			return e
		}
		if e := tx.Table("qixi_crm_b_distribution_relation_audit").Create(map[string]any{"user_id": userID, "previous_parent_user_id": nullableLevelFromSQL(previous), "parent_user_id": nullableLevel(in.ParentUserID), "reason": in.Reason, "operator_admin_id": middleware.AdminID(c), "idempotency_key": in.IdempotencyKey}).Error; e != nil {
			return e
		}
		out = gin.H{"user_id": userID, "previous_parent_user_id": nullableLevelFromSQL(previous), "parent_user_id": in.ParentUserID, "replayed": false}
		return nil
	})
	switch {
	case err == nil:
		response.OK(c, out)
	case errors.Is(err, errReferrerNotFound):
		response.Fail(c, http.StatusNotFound, "用户或上级用户不存在")
	case errors.Is(err, errReferrerConflict):
		response.Fail(c, http.StatusConflict, "推荐关系形成循环或幂等键冲突")
	default:
		response.Fail(c, http.StatusInternalServerError, "推荐关系调整失败")
	}
}

func (h *Handler) ensureNoReferralCycle(tx *gorm.DB, userID, parentID uint64) error {
	current := parentID
	for depth := 0; depth < 64 && current != 0; depth++ {
		if current == userID {
			return errReferrerConflict
		}
		var r struct {
			ParentUserID sql.NullInt64 `gorm:"column:parent_user_id"`
		}
		e := tx.Table("qixi_crm_b_distribution_relation").Clauses(clause.Locking{Strength: "UPDATE"}).Select("parent_user_id").Where("user_id=?", current).Take(&r).Error
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return nil
		}
		if e != nil {
			return e
		}
		if !r.ParentUserID.Valid {
			return nil
		}
		current = uint64(r.ParentUserID.Int64)
	}
	if current != 0 {
		return errReferrerConflict
	}
	return nil
}
