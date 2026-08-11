package cloudconfig

import (
	"context"
	"errors"
	"net/http"
	"strings"

	configdomain "github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/cloudconfig"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc       *configdomain.Service
	adminDB   *gorm.DB
	publisher paymentPublisher
}

type paymentPublisher interface {
	Publish(context.Context, map[string]string) error
}

func NewHandler(svc *configdomain.Service, adminDB *gorm.DB, publisher ...paymentPublisher) *Handler {
	h := &Handler{svc: svc, adminDB: adminDB}
	if len(publisher) > 0 {
		h.publisher = publisher[0]
	}
	return h
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/setting/cloud-configs", h.List)
	r.GET("/setting/cloud-configs/:group", h.Get)
	// 新统一后台的身份、角色和按钮权限均来自 qixi_crm_a_*。
	// 不能使用 legacy identity.Service（其旧模型会访问 qixi_m_* 表）。
	appSetting := middleware.RequireAdminRoles("platform")
	routineManage := middleware.RequireAdminMenu(h.adminDB, "app.routine.manage")
	r.GET("/setting/routine-config", appSetting, routineManage, h.GetRoutine)
	r.PUT("/setting/routine-config", appSetting, routineManage, h.SaveRoutine)
	mobileAppManage := middleware.RequireAdminMenu(h.adminDB, "app.mobile.manage")
	r.GET("/setting/mobile-app-config/:platform", appSetting, mobileAppManage, h.GetMobileApp)
	r.PUT("/setting/mobile-app-config/:platform", appSetting, mobileAppManage, h.SaveMobileApp)
	pushManage := middleware.RequireAdminMenu(h.adminDB, "app.push.manage")
	r.GET("/setting/push-config/:platform", appSetting, pushManage, h.GetPush)
	r.PUT("/setting/push-config/:platform", appSetting, pushManage, h.SavePush)
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
	h.getGroup(c, c.Param("group"))
}

func (h *Handler) GetRoutine(c *gin.Context) {
	h.getGroup(c, "wechat_mini_program")
}

func (h *Handler) getGroup(c *gin.Context, group string) {
	row, err := h.svc.Get(c.Request.Context(), group)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Save(c *gin.Context) {
	h.saveGroup(c, c.Param("group"))
}

func (h *Handler) SaveRoutine(c *gin.Context) {
	h.saveGroup(c, "wechat_mini_program")
}

func (h *Handler) GetMobileApp(c *gin.Context) {
	group, ok := nativeConfigGroup("mobile_app", c.Param("platform"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "不支持的应用平台")
		return
	}
	h.getGroup(c, group)
}

func (h *Handler) SaveMobileApp(c *gin.Context) {
	group, ok := nativeConfigGroup("mobile_app", c.Param("platform"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "不支持的应用平台")
		return
	}
	h.saveGroup(c, group)
}

func (h *Handler) GetPush(c *gin.Context) {
	group, ok := nativeConfigGroup("umeng_push", c.Param("platform"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "不支持的推送平台")
		return
	}
	h.getGroup(c, group)
}

func (h *Handler) SavePush(c *gin.Context) {
	group, ok := nativeConfigGroup("umeng_push", c.Param("platform"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "不支持的推送平台")
		return
	}
	h.saveGroup(c, group)
}

func nativeConfigGroup(prefix, platform string) (string, bool) {
	switch strings.TrimSpace(platform) {
	case "ios", "android", "harmony":
		return prefix + "_" + strings.TrimSpace(platform), true
	default:
		return "", false
	}
}

func (h *Handler) saveGroup(c *gin.Context, group string) {
	var in configdomain.SaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.Save(c.Request.Context(), group, in, middleware.AdminID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	if group == "payment" && h.publisher != nil {
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
