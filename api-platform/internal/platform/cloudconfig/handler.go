package cloudconfig

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	configdomain "github.com/crmlive/qixi-live-ecrm/api-platform/internal/domain/cloudconfig"
	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/domain/identity"
	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/pkg/response"
)

type Handler struct {
	svc       *configdomain.Service
	id        *identity.Service
	publisher paymentPublisher
}

type paymentPublisher interface {
	Publish(context.Context, map[string]string) error
}

func NewHandler(svc *configdomain.Service, id *identity.Service, publisher ...paymentPublisher) *Handler {
	h := &Handler{svc: svc, id: id}
	if len(publisher) > 0 {
		h.publisher = publisher[0]
	}
	return h
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/setting/cloud-configs", h.List)
	r.GET("/setting/cloud-configs/:group", h.Get)
	r.PUT("/setting/cloud-configs/:group", middleware.RequirePlatformMenu(h.id, identity.PlatPermCloudConfigWrite), h.Save)
}

func (h *Handler) List(c *gin.Context) {
	rows, err := h.svc.List(c.Request.Context())
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"list": rows})
}
func (h *Handler) Get(c *gin.Context) {
	row, err := h.svc.Get(c.Request.Context(), c.Param("group"))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}
func (h *Handler) Save(c *gin.Context) {
	var in configdomain.SaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.Save(c.Request.Context(), c.Param("group"), in, middleware.AdminID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	if c.Param("group") == "payment" && h.publisher != nil {
		values, err := h.svc.Values(c.Request.Context(), "payment")
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "支付配置读取失败")
			return
		}
		if err := h.publisher.Publish(c.Request.Context(), values); err != nil {
			response.Fail(c, http.StatusInternalServerError, "支付配置投影失败，请重试保存")
			return
		}
	}
	response.OK(c, row)
}
func writeErr(c *gin.Context, err error) {
	if errors.Is(err, configdomain.ErrBadGroup) || errors.Is(err, configdomain.ErrBadField) || errors.Is(err, configdomain.ErrBadValue) {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Fail(c, http.StatusInternalServerError, "配置操作失败")
}
