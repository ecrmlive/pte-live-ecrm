package svip

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/response"
)

type Handler struct{ svc *identity.Service }

func NewHandler(svc *identity.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/users", h.ListUsers)
	r.PUT("/users/:id/svip", middleware.RequirePlatformMenu(h.svc, identity.PlatPermSvipUpdate), h.SetSvip)
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
			"uid": u.UID, "account": u.Account, "nickname": u.Nickname, "phone": u.Phone,
			"is_svip": u.IsSvip, "svip_endtime": u.SvipEndtime,
			"is_svip_active": identity.UserSvipActive(&u),
			"integral":      u.Integral, "now_money": u.NowMoney,
		})
	}
	response.OK(c, gin.H{"list": list, "total": res.Total, "page": res.Page, "limit": res.Limit})
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
