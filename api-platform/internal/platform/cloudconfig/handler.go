package cloudconfig

import (
	"context"
	"errors"
	"net/http"
	"strings"

	configdomain "github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/cloudconfig"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
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
	// 高德 Web JS 需在浏览器侧使用 Key/安全密钥；云配置页对 Secret 字段掩码，故单独提供已鉴权只读接口。
	r.GET("/setting/map-client-config", h.MapClientConfig)
}

func (h *Handler) List(c *gin.Context) {
	rows, err := h.svc.List(c.Request.Context())
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"list": rows})
}

// MapClientConfig 返回平台后台地图取点所需的高德 Web JS 凭证（已登录管理员）。
// 仅下发 Web JS Key / 安全密钥；Web 服务 Key 与各端 Key 不下发。
func (h *Handler) MapClientConfig(c *gin.Context) {
	values, err := h.svc.Values(c.Request.Context(), "amap")
	if err != nil {
		writeErr(c, err)
		return
	}
	key := strings.TrimSpace(values["amap_web_js_key"])
	security := strings.TrimSpace(values["amap_web_js_security_code"])
	response.OK(c, gin.H{
		"provider":                  "amap",
		"amap_web_js_key":           key,
		"amap_web_js_security_code": security,
		"configured":                key != "" && security != "",
	})
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
