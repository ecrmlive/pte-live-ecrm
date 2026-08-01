package broadcast

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/broadcast"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/identity"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/response"
)

type Handler struct {
	svc *broadcast.Service
	id  *identity.Service
}

func NewHandler(svc *broadcast.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/broadcast/rooms", h.List)
	r.GET("/broadcast/rooms/:id", h.Get)
	r.POST("/broadcast/rooms", middleware.RequireMerchantMenu(h.id, identity.MerPermBroadcastCreate), h.Create)
	r.PUT("/broadcast/rooms/:id", h.Update)
	r.PUT("/broadcast/rooms/:id/live", middleware.RequireMerchantMenu(h.id, identity.MerPermBroadcastLive), h.SetLive)
	r.PUT("/broadcast/rooms/:id/goods", middleware.RequireMerchantMenu(h.id, identity.MerPermBroadcastGoods), h.SetGoods)
	r.DELETE("/broadcast/rooms/:id", middleware.RequireMerchantMenu(h.id, identity.MerPermBroadcastDelete), h.Delete)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListMerchant(c.Request.Context(), middleware.MerID(c), page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	room, err := h.svc.Get(c.Request.Context(), uint(id), true)
	if err != nil {
		writeErr(c, err)
		return
	}
	if room.MerID != middleware.MerID(c) {
		response.Fail(c, http.StatusForbidden, broadcast.ErrForbidden.Error())
		return
	}
	response.OK(c, room)
}

func (h *Handler) Create(c *gin.Context) {
	var in broadcast.SaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if len(in.ProductIDs) > 0 {
		if err := h.id.RequireMerchantMenu(c.Request.Context(), middleware.AdminID(c), identity.MerPermBroadcastGoods); err != nil {
			writePerm(c, err)
			return
		}
	}
	room, err := h.svc.Create(c.Request.Context(), middleware.MerID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, room)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in broadcast.SaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	adminID := middleware.AdminID(c)
	if in.LiveStatus != nil {
		if err := h.id.RequireMerchantMenu(c.Request.Context(), adminID, identity.MerPermBroadcastLive); err != nil {
			writePerm(c, err)
			return
		}
	}
	if in.ProductIDs != nil {
		if err := h.id.RequireMerchantMenu(c.Request.Context(), adminID, identity.MerPermBroadcastGoods); err != nil {
			writePerm(c, err)
			return
		}
	}
	room, err := h.svc.Update(c.Request.Context(), middleware.MerID(c), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, room)
}

func (h *Handler) SetLive(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		LiveStatus *int16 `json:"live_status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.LiveStatus == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	room, err := h.svc.Update(c.Request.Context(), middleware.MerID(c), uint(id), broadcast.SaveInput{LiveStatus: body.LiveStatus})
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, room)
}

func (h *Handler) SetGoods(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in broadcast.GoodsInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	room, err := h.svc.SetGoods(c.Request.Context(), middleware.MerID(c), uint(id), in.ProductIDs)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, room)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), middleware.MerID(c), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writePerm(c *gin.Context, err error) {
	if errors.Is(err, identity.ErrNoPerm) {
		response.Fail(c, http.StatusForbidden, err.Error())
		return
	}
	if errors.Is(err, identity.ErrNotFound) {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}
	response.Fail(c, http.StatusInternalServerError, "权限校验失败")
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, broadcast.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, broadcast.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, broadcast.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
