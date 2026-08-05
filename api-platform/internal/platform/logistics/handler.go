package logistics

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/logistics"
	platformcityevent "github.com/crmlive/pte-live-ecrm/api-platform/internal/event/platformcity"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc     *logistics.Service
	adminDB *gorm.DB
}

func NewHandler(svc *logistics.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	manage := middleware.RequireAdminMenu(h.adminDB, "freight.express.manage")
	platform := middleware.RequireAdminRoles("platform")
	r.GET("/express", platform, manage, h.ListExpress)
	r.POST("/express", platform, manage, h.CreateExpress)
	r.PUT("/express/:id", platform, manage, h.UpdateExpress)
	r.DELETE("/express/:id", platform, manage, h.DeleteExpress)
	r.GET("/city", platform, manage, h.ListCity)
	r.POST("/city/resync", platform, manage, h.ResyncCity)
}

func (h *Handler) ListExpress(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListExpress(c.Request.Context(), page, limit, false)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) CreateExpress(c *gin.Context) {
	var in logistics.ExpressInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateExpress(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateExpress(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in logistics.ExpressInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateExpress(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteExpress(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteExpress(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListCity(c *gin.Context) {
	var parent *uint
	if raw := c.Query("parent_id"); raw != "" {
		v, err := strconv.ParseUint(raw, 10, 64)
		if err == nil {
			u := uint(v)
			parent = &u
		}
	}
	list, err := h.svc.ListCity(c.Request.Context(), parent)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) ResyncCity(c *gin.Context) {
	list, err := h.svc.ListCity(c.Request.Context(), nil)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取城市数据失败")
		return
	}
	err = h.adminDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		for _, city := range list {
			payload, err := json.Marshal(city)
			if err != nil {
				return err
			}
			if err = tx.Table("qixi_crm_a_outbox").Create(map[string]any{"event_type": platformcityevent.Upserted, "aggregate_type": "city", "aggregate_id": strconv.FormatUint(uint64(city.CityID), 10), "payload": payload, "status": "pending"}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "发布城市投影失败")
		return
	}
	response.OK(c, gin.H{"queued": len(list)})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, logistics.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, logistics.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
