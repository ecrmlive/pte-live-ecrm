package svip

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/domain/merchant"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/pkg/response"
)

type Handler struct {
	svc *merchant.Service
	id  *identity.Service
}

func NewHandler(svc *merchant.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/setting/svip", h.Get)
	r.PUT("/setting/svip", middleware.RequireMerchantMenu(h.id, identity.MerPermSvipUpdate), h.Update)
}

func (h *Handler) Get(c *gin.Context) {
	cfg, err := h.svc.GetSvipConfig(c.Request.Context(), middleware.MerID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, cfg)
}

func (h *Handler) Update(c *gin.Context) {
	var body struct {
		SvipCouponMerge *int8 `json:"svip_coupon_merge"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.SvipCouponMerge == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	cfg, err := h.svc.UpdateSvipConfig(c.Request.Context(), middleware.MerID(c), *body.SvipCouponMerge)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, cfg)
}
