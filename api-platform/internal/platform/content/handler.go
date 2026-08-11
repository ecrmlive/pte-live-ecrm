package content

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/content"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc     *content.Service
	adminDB *gorm.DB
}

func NewHandler(svc *content.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/notices", h.List)
	operationWrite := middleware.RequireAdminRoles("platform", "operations")
	noticeWrite := middleware.RequireAdminMenu(h.adminDB, "content.notice.manage")
	r.POST("/notices", operationWrite, noticeWrite, h.Create)
	r.PUT("/notices/:id", operationWrite, noticeWrite, h.Update)
	r.DELETE("/notices/:id", operationWrite, noticeWrite, h.Delete)
	r.GET("/agreements", h.ListAgreements)
	r.GET("/agreements/:key", h.GetAgreement)
	r.PUT("/agreements/:key", operationWrite, middleware.RequireAdminMenuAny(h.adminDB,
		"setting.agreement.manage",
		"user.svip.agreement.manage",
		"user.level.description.manage",
		"accounts.invoice.desc.manage",
	), h.SaveAgreement)
	r.GET("/setting/sms", h.GetSMS)
	r.PUT("/setting/sms", middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.adminDB, "setting.sms.manage"), h.SaveSMS)
	shopSetting := middleware.RequireAdminRoles("platform")
	shopManage := middleware.RequireAdminMenu(h.adminDB, "setting.shop.manage")
	r.GET("/setting/shop", shopSetting, shopManage, h.GetShop)
	r.PUT("/setting/shop", shopSetting, shopManage, h.SaveShop)
	marginSetting := middleware.RequireAdminRoles("platform")
	marginManage := middleware.RequireAdminMenu(h.adminDB, "store.margin_config.manage")
	r.GET("/setting/margin", marginSetting, marginManage, h.GetMargin)
	r.PUT("/setting/margin", marginSetting, marginManage, h.SaveMargin)
	merchantApplySetting := middleware.RequireAdminRoles("platform")
	merchantApplyManage := middleware.RequireAdminMenu(h.adminDB, "merchant.mgmt.settings.manage")
	r.GET("/setting/merchant-apply", merchantApplySetting, merchantApplyManage, h.GetMerchantApply)
	r.PUT("/setting/merchant-apply", merchantApplySetting, merchantApplyManage, h.SaveMerchantApply)
	agentZoneSetting := middleware.RequireAdminRoles("platform")
	agentZoneManage := middleware.RequireAdminMenu(h.adminDB, "region.agent_settings.manage")
	r.GET("/setting/agent-zone", agentZoneSetting, agentZoneManage, h.GetAgentZone)
	r.PUT("/setting/agent-zone", agentZoneSetting, agentZoneManage, h.SaveAgentZone)
	paySetting := middleware.RequireAdminRoles("platform")
	payManage := middleware.RequireAdminMenu(h.adminDB, "setting.pay.manage")
	r.GET("/setting/pay", paySetting, payManage, h.GetPay)
	r.PUT("/setting/pay", paySetting, payManage, h.SavePay)
	wechatAppSetting := middleware.RequireAdminRoles("platform")
	wechatAppManage := middleware.RequireAdminMenu(h.adminDB, "app.wechat.manage")
	r.GET("/setting/wechat-app", wechatAppSetting, wechatAppManage, h.GetWechatApp)
	r.PUT("/setting/wechat-app", wechatAppSetting, wechatAppManage, h.SaveWechatApp)
	productSetting := middleware.RequireAdminRoles("platform")
	r.GET("/product/price-descriptions", productSetting, h.GetPriceDescriptions)
	r.PUT("/product/price-descriptions", productSetting, h.SavePriceDescriptions)
	r.GET("/product/activity-labels", productSetting, h.GetActivityLabels)
	r.PUT("/product/activity-labels", productSetting, h.SaveActivityLabels)

	storageSetting := middleware.RequireAdminRoles("platform")
	storageManage := middleware.RequireAdminMenu(h.adminDB, "setting.storage.manage")
	r.GET("/setting/storage", storageSetting, storageManage, h.GetStorage)
	r.PUT("/setting/storage", storageSetting, storageManage, h.SaveStorage)
	userSetupSetting := middleware.RequireAdminRoles("platform")
	userSetupManage := middleware.RequireAdminMenu(h.adminDB, "user.setup.manage")
	r.GET("/setting/user-setup", userSetupSetting, userSetupManage, h.GetUserSetup)
	r.PUT("/setting/user-setup", userSetupSetting, userSetupManage, h.SaveUserSetup)
	transferSetting := middleware.RequireAdminRoles("platform")
	transferManage := middleware.RequireAdminMenu(h.adminDB, "accounts.transfer_settings.manage")
	r.GET("/setting/transfer-settings", transferSetting, transferManage, h.GetTransferSettings)
	r.PUT("/setting/transfer-settings", transferSetting, transferManage, h.SaveTransferSettings)
	distributionSetting := middleware.RequireAdminRoles("platform", "operations")
	distributionManage := middleware.RequireAdminMenu(h.adminDB, "promoter.config.manage")
	r.GET("/setting/distribution", distributionSetting, distributionManage, h.GetDistribution)
	r.PUT("/setting/distribution", distributionSetting, distributionManage, h.SaveDistribution)
	groupBuyingSetting := middleware.RequireAdminRoles("platform", "operations")
	groupBuyingManage := middleware.RequireAdminMenu(h.adminDB, "marketing.combination.manage")
	r.GET("/setting/group-buying", groupBuyingSetting, groupBuyingManage, h.GetGroupBuying)
	r.PUT("/setting/group-buying", groupBuyingSetting, groupBuyingManage, h.SaveGroupBuying)
	integralSetting := middleware.RequireAdminRoles("platform", "operations")
	integralManage := middleware.RequireAdminMenu(h.adminDB, "marketing.integral.config")
	r.GET("/setting/integral", integralSetting, integralManage, h.GetIntegral)
	r.PUT("/setting/integral", integralSetting, integralManage, h.SaveIntegral)
	balanceSetting := middleware.RequireAdminRoles("platform", "operations")
	balanceRead := middleware.RequireAdminMenu(h.adminDB, "marketing.balance.settings.read")
	balanceManage := middleware.RequireAdminMenu(h.adminDB, "marketing.balance.settings.manage")
	r.GET("/setting/balance", balanceSetting, balanceRead, h.GetBalance)
	r.PUT("/setting/balance", balanceSetting, balanceManage, h.SaveBalance)

	appSetting := middleware.RequireAdminRoles("platform")
	routineManage := middleware.RequireAdminMenu(h.adminDB, "app.routine.manage")
	r.GET("/setting/routine-app", appSetting, routineManage, h.GetRoutineApp)
	r.PUT("/setting/routine-app", appSetting, routineManage, h.SaveRoutineApp)
	wechatReplyManage := middleware.RequireAdminMenu(h.adminDB, "app.wechat_reply.manage")
	r.GET("/setting/wechat-reply", appSetting, wechatReplyManage, h.GetWechatReply)
	r.PUT("/setting/wechat-reply", appSetting, wechatReplyManage, h.SaveWechatReply)
	wechatMenusManage := middleware.RequireAdminMenu(h.adminDB, "app.wechat_menus.manage")
	r.GET("/setting/wechat-menus", appSetting, wechatMenusManage, h.GetWechatMenus)
	r.PUT("/setting/wechat-menus", appSetting, wechatMenusManage, h.SaveWechatMenus)
	wechatTemplateManage := middleware.RequireAdminMenu(h.adminDB, "app.wechat_template.manage")
	r.GET("/setting/wechat-template", appSetting, wechatTemplateManage, h.GetWechatTemplate)
	r.PUT("/setting/wechat-template", appSetting, wechatTemplateManage, h.SaveWechatTemplate)
	wechatNewsManage := middleware.RequireAdminMenu(h.adminDB, "app.wechat_news.manage")
	r.GET("/setting/wechat-news", appSetting, wechatNewsManage, h.GetWechatNews)
	r.PUT("/setting/wechat-news", appSetting, wechatNewsManage, h.SaveWechatNews)

	maintainSetting := middleware.RequireAdminRoles("platform")
	maintainManage := middleware.RequireAdminMenu(h.adminDB, "maintain.cache.manage")
	r.POST("/maintain/cache/clear", maintainSetting, maintainManage, h.ClearMaintainCache)
	// 热搜/组合数据/备份/系统表单已迁至 nativeconfigitem（qixi_crm_a_config_item）。
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListAdmin(c.Request.Context(), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Create(c *gin.Context) {
	var in content.NoticeInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	n, err := h.svc.Create(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, n)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in content.NoticeInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	n, err := h.svc.Update(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, n)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListAgreements(c *gin.Context) {
	list, err := h.svc.ListAgreements(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) GetAgreement(c *gin.Context) {
	row, err := h.svc.GetAgreement(c.Request.Context(), c.Param("key"))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) SaveAgreement(c *gin.Context) {
	var in content.AgreeSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.SaveAgreement(c.Request.Context(), c.Param("key"), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

type smsSaveReq struct {
	Config string `json:"config" binding:"required"`
}

func (h *Handler) GetSMS(c *gin.Context) {
	raw, err := h.svc.GetSMSConfig(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"config": raw, "note": "仅允许无密钥 stub 配置；不保存或回显真实短信通道凭据"})
}

func (h *Handler) SaveSMS(c *gin.Context) {
	var req smsSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	raw, err := h.svc.SaveSMSConfig(c.Request.Context(), req.Config)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"config": raw})
}

type jsonConfigSaveReq struct {
	Config string `json:"config" binding:"required"`
}

func (h *Handler) GetShop(c *gin.Context) {
	raw, err := h.svc.GetShopConfig(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"config": raw, "note": "仅保存商城基础开关与超时规则；不含密钥或支付凭据"})
}

func (h *Handler) SaveShop(c *gin.Context) {
	var req jsonConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	raw, err := h.svc.SaveShopConfig(c.Request.Context(), req.Config)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"config": raw})
}

func (h *Handler) GetMargin(c *gin.Context) {
	raw, err := h.svc.GetMarginConfig(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{
		"config": raw,
		"note":   "对齐 CRMEB 保证金补缴提醒：开关开启后按天数连续提醒；期满未补足则自动关闭店铺（定时关店任务待接入）",
	})
}

func (h *Handler) SaveMargin(c *gin.Context) {
	var req jsonConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	raw, err := h.svc.SaveMarginConfig(c.Request.Context(), req.Config)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"config": raw})
}

func (h *Handler) GetMerchantApply(c *gin.Context) {
	raw, err := h.svc.GetMerchantApplyConfig(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{
		"config": raw,
		"note":   "商户入驻页背景图与自定义表单字段；系统表单字段固定展示不可删",
	})
}

func (h *Handler) SaveMerchantApply(c *gin.Context) {
	var req jsonConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	raw, err := h.svc.SaveMerchantApplyConfig(c.Request.Context(), req.Config)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"config": raw})
}

func (h *Handler) GetAgentZone(c *gin.Context) {
	raw, err := h.svc.GetAgentZoneConfig(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{
		"config": raw,
		"note":   "区域代理默认三级提成与代理申请自定义表单；系统字段固定展示不可删",
	})
}

func (h *Handler) SaveAgentZone(c *gin.Context) {
	var req jsonConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	raw, err := h.svc.SaveAgentZoneConfig(c.Request.Context(), req.Config)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"config": raw})
}

func (h *Handler) GetPay(c *gin.Context) {
	raw, err := h.svc.GetPayConfig(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"config": raw, "note": "仅保存支付方式开关；不保存或回显微信/支付宝密钥、证书或令牌"})
}

func (h *Handler) SavePay(c *gin.Context) {
	var req jsonConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	raw, err := h.svc.SavePayConfig(c.Request.Context(), req.Config)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"config": raw})
}

func (h *Handler) GetWechatApp(c *gin.Context) {
	raw, err := h.svc.GetWechatAppConfig(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"config": raw, "note": "仅保存公众号名称与启用开关；不保存或回显 AppSecret、Token、EncodingAESKey 等密钥"})
}

func (h *Handler) SaveWechatApp(c *gin.Context) {
	var req jsonConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	raw, err := h.svc.SaveWechatAppConfig(c.Request.Context(), req.Config)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"config": raw})
}

type cacheListSaveReq struct {
	List []content.CacheListItem `json:"list"`
}

func (h *Handler) GetPriceDescriptions(c *gin.Context) {
	list, err := h.svc.GetCacheList(c.Request.Context(), content.PriceDescriptionCacheKey)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) SavePriceDescriptions(c *gin.Context) {
	var req cacheListSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	list, err := h.svc.SaveCacheList(c.Request.Context(), content.PriceDescriptionCacheKey, req.List)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) GetActivityLabels(c *gin.Context) {
	list, err := h.svc.GetCacheList(c.Request.Context(), content.ActivityLabelCacheKey)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) SaveActivityLabels(c *gin.Context) {
	var req cacheListSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	list, err := h.svc.SaveCacheList(c.Request.Context(), content.ActivityLabelCacheKey, req.List)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) GetStorage(c *gin.Context) {
	h.getJSONSetting(c, h.svc.GetStorageConfig, "仅保存对象存储开关与桶名展示；不保存或回显 SecretId/SecretKey")
}

func (h *Handler) SaveStorage(c *gin.Context) {
	h.saveJSONSetting(c, h.svc.SaveStorageConfig)
}

func (h *Handler) GetUserSetup(c *gin.Context) {
	h.getJSONSetting(c, h.svc.GetUserSetupConfig, "对齐 CRMEB 用户设置：默认头像、扩展信息字段、登录注册与注册有礼；不含短信或第三方登录密钥")
}

func (h *Handler) SaveUserSetup(c *gin.Context) {
	h.saveJSONSetting(c, h.svc.SaveUserSetupConfig)
}

func (h *Handler) GetTransferSettings(c *gin.Context) {
	h.getJSONSetting(c, h.svc.GetTransferSettingsConfig, "仅保存转账监管开关与最低金额；真实打款凭据不在后台保存")
}

func (h *Handler) SaveTransferSettings(c *gin.Context) {
	h.saveJSONSetting(c, h.svc.SaveTransferSettingsConfig)
}

func (h *Handler) GetDistribution(c *gin.Context) {
	h.getJSONSetting(c, h.svc.GetDistributionConfig, "分销配置")
}

func (h *Handler) SaveDistribution(c *gin.Context) {
	h.saveJSONSetting(c, h.svc.SaveDistributionConfig)
}

func (h *Handler) GetGroupBuying(c *gin.Context) {
	h.getJSONSetting(c, h.svc.GetGroupBuyingConfig, "拼团设置：虚拟成团启用与真实成团最小比例")
}

func (h *Handler) SaveGroupBuying(c *gin.Context) {
	h.saveJSONSetting(c, h.svc.SaveGroupBuyingConfig)
}

func (h *Handler) GetIntegral(c *gin.Context) {
	h.getJSONSetting(c, h.svc.GetIntegralConfig, "积分配置：开关、抵用比例、赠送与清除规则")
}

func (h *Handler) SaveIntegral(c *gin.Context) {
	h.saveJSONSetting(c, h.svc.SaveIntegralConfig)
}

func (h *Handler) GetBalance(c *gin.Context) {
	h.getJSONSetting(c, h.svc.GetBalanceConfig, "余额设置：余额功能、充值开关、最低金额与注意事项")
}

func (h *Handler) SaveBalance(c *gin.Context) {
	h.saveJSONSetting(c, h.svc.SaveBalanceConfig)
}

func (h *Handler) GetRoutineApp(c *gin.Context) {
	h.getAppStubSetting(c, content.RoutineAppConfigKey)
}

func (h *Handler) SaveRoutineApp(c *gin.Context) {
	h.saveAppStubSetting(c, content.RoutineAppConfigKey)
}

func (h *Handler) GetWechatReply(c *gin.Context) {
	h.getAppStubSetting(c, content.WechatReplyConfigKey)
}

func (h *Handler) SaveWechatReply(c *gin.Context) {
	h.saveAppStubSetting(c, content.WechatReplyConfigKey)
}

func (h *Handler) GetWechatMenus(c *gin.Context) {
	buttons, err := h.svc.GetWechatMenus(c.Request.Context())
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{
		"wechat_menus": buttons,
		"note":         "菜单配置保存在平台缓存；未配置公众号服务端凭据时仅本地保存，不向微信推送",
	})
}

func (h *Handler) SaveWechatMenus(c *gin.Context) {
	var req struct {
		Button []content.WechatMenuButton `json:"button"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	saved, err := h.svc.SaveWechatMenus(c.Request.Context(), req.Button)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{
		"wechat_menus": saved,
		"published":    false,
		"note":         "已保存本地菜单配置；向微信公众号发布需配置服务端凭据后另行对接",
	})
}

func (h *Handler) GetWechatTemplate(c *gin.Context) {
	h.getAppStubSetting(c, content.WechatTemplateConfigKey)
}

func (h *Handler) SaveWechatTemplate(c *gin.Context) {
	h.saveAppStubSetting(c, content.WechatTemplateConfigKey)
}

func (h *Handler) GetWechatNews(c *gin.Context) {
	h.getAppStubSetting(c, content.WechatNewsConfigKey)
}

func (h *Handler) SaveWechatNews(c *gin.Context) {
	h.saveAppStubSetting(c, content.WechatNewsConfigKey)
}

func (h *Handler) ClearMaintainCache(c *gin.Context) {
	if err := h.svc.ClearMaintainCache(c.Request.Context()); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true, "note": "已提交缓存清理请求；不含密钥或敏感凭据"})
}

func (h *Handler) getJSONSetting(c *gin.Context, load func(context.Context) (string, error), note string) {
	raw, err := load(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"config": raw, "note": note})
}

func (h *Handler) saveJSONSetting(c *gin.Context, save func(context.Context, string) (string, error)) {
	var req jsonConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	raw, err := save(c.Request.Context(), req.Config)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"config": raw})
}

func (h *Handler) getAppStubSetting(c *gin.Context, key string) {
	raw, err := h.svc.GetAppStubConfig(c.Request.Context(), key)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"config": raw, "note": "仅保存名称与启用开关；不含 AppSecret、Token 或 EncodingAESKey 等密钥"})
}

func (h *Handler) saveAppStubSetting(c *gin.Context, key string) {
	var req jsonConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	raw, err := h.svc.SaveAppStubConfig(c.Request.Context(), key, req.Config)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"config": raw})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, content.ErrNotFound), errors.Is(err, content.ErrAgreeNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, content.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
