package profitsharing

import (
	"errors"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db} }
func (h *Handler) Register(r gin.IRoutes) {
	p := middleware.RequireAdminRoles("platform")
	m := middleware.RequireAdminMenu(h.db, "merchant.profitsharing.review")
	r.GET("/merchant-profitsharing-applications", p, m, h.List)
	r.GET("/merchant-profitsharing-applications/:id", p, m, h.Get)
	r.POST("/merchant-profitsharing-applications/:id/review", p, m, h.Review)
	r.PUT("/merchant-profitsharing-applications/:id/note", p, m, h.Note)
}

type row struct {
	ID            uint      `json:"id"`
	MerchantID    uint      `gorm:"column:merchant_id" json:"merchant_id"`
	ApplicationNo string    `gorm:"column:application_no" json:"application_no"`
	Status        string    `json:"status"`
	Description   string    `json:"description"`
	ReviewNote    string    `gorm:"column:review_note" json:"review_note"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

const tab = "qixi_crm_a_merchant_profitsharing_application"

func (h *Handler) List(c *gin.Context) {
	var rows []row
	q := h.db.WithContext(c.Request.Context()).Table(tab).Order("id DESC")
	if s := strings.TrimSpace(c.Query("status")); s != "" {
		q = q.Where("status=?", s)
	}
	if e := q.Find(&rows).Error; e != nil {
		fail(c, e)
		return
	}
	response.OK(c, gin.H{"list": rows})
}
func (h *Handler) Get(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		return
	}
	var v row
	if e := h.db.WithContext(c.Request.Context()).Table(tab).Where("id=?", id).Take(&v).Error; e != nil {
		fail(c, e)
		return
	}
	response.OK(c, v)
}
func (h *Handler) Review(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		return
	}
	var q struct {
		Approved bool   `json:"approved"`
		Note     string `json:"note"`
	}
	if c.ShouldBindJSON(&q) != nil || !validReviewNote(q.Note) {
		response.Fail(c, 400, "审核参数错误")
		return
	}
	to := "rejected"
	if q.Approved {
		to = "approved"
	}
	e := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var v row
		if e := tx.Table(tab).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id=?", id).Take(&v).Error; e != nil {
			return e
		}
		if v.Status != "applied" {
			return errors.New("分账申请状态已变化")
		}
		res := tx.Table(tab).Where("id=? AND status='applied'", id).Updates(map[string]any{"status": to, "review_note": strings.TrimSpace(q.Note), "reviewed_by": middleware.AdminID(c), "reviewed_at": time.Now()})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errors.New("分账申请状态已变化")
		}
		return tx.Table("qixi_crm_a_merchant_profitsharing_audit").Create(map[string]any{"application_id": id, "from_status": "applied", "to_status": to, "note": strings.TrimSpace(q.Note), "operator_admin_id": middleware.AdminID(c)}).Error
	})
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			response.Fail(c, 404, "分账申请不存在")
		} else if strings.Contains(e.Error(), "状态") {
			response.Fail(c, 409, e.Error())
		} else {
			fail(c, e)
		}
		return
	}
	response.OK(c, gin.H{"ok": true, "status": to})
}
func (h *Handler) Note(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		return
	}
	var q struct {
		Note string `json:"note"`
	}
	if c.ShouldBindJSON(&q) != nil || !validReviewNote(q.Note) {
		response.Fail(c, 400, "备注参数错误")
		return
	}
	res := h.db.WithContext(c.Request.Context()).Table(tab).Where("id=?", id).Update("review_note", strings.TrimSpace(q.Note))
	if res.Error != nil {
		fail(c, res.Error)
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, 404, "分账申请不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
}
func id(c *gin.Context) (uint, bool) {
	n, e := strconv.ParseUint(c.Param("id"), 10, 64)
	if e != nil || n == 0 {
		response.Fail(c, 400, "申请 ID 错误")
		return 0, false
	}
	return uint(n), true
}

func validReviewNote(note string) bool {
	return len([]rune(strings.TrimSpace(note))) <= 500
}

func fail(c *gin.Context, e error) {
	if errors.Is(e, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "分账申请不存在")
		return
	}
	response.Fail(c, 500, "分账申请操作失败")
}
