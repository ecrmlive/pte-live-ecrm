package serviceportal

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/cs"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/trade"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

// Handler 客服工作台最小竖切：查单 + 快捷回复（阶段 6b / 7 可配置）
type Handler struct {
	ord       *trade.Service
	cs        *cs.Service
	demoToken string
}

func NewHandler(ord *trade.Service, csSvc *cs.Service, demoToken string) *Handler {
	if demoToken == "" {
		demoToken = "service_demo_token"
	}
	return &Handler{ord: ord, cs: csSvc, demoToken: demoToken}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/ping", h.Ping)
	r.GET("/orders/:id", h.OrderDetail)
	r.GET("/quick-replies", h.QuickReplies)
}

func (h *Handler) Ping(c *gin.Context) {
	response.OK(c, gin.H{"prefix": "/api/service/v1", "ok": true})
}

func (h *Handler) authOK(c *gin.Context) bool {
	if c.GetHeader("X-Service-Token") != h.demoToken && c.Query("token") != h.demoToken {
		response.Fail(c, http.StatusUnauthorized, "客服凭证无效")
		return false
	}
	return true
}

func (h *Handler) OrderDetail(c *gin.Context) {
	if !h.authOK(c) {
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	merID64, _ := strconv.ParseUint(c.Query("mer_id"), 10, 64)
	if id == 0 || merID64 == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	o, err := h.ord.GetMerchantOrder(c.Request.Context(), uint(merID64), uint(id))
	if err != nil {
		if errors.Is(err, trade.ErrNotFound) || errors.Is(err, trade.ErrForbidden) {
			response.Fail(c, http.StatusNotFound, "订单不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, o)
}

func (h *Handler) QuickReplies(c *gin.Context) {
	if !h.authOK(c) {
		return
	}
	merID64, _ := strconv.ParseUint(c.DefaultQuery("mer_id", "1"), 10, 64)
	if merID64 == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	rows, err := h.cs.ListEnabled(c.Request.Context(), uint(merID64))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		list = append(list, map[string]any{
			"id":      r.ServiceReplyID,
			"title":   r.Keyword,
			"keyword": r.Keyword,
			"content": r.Content,
			"type":    r.Type,
		})
	}
	response.OK(c, gin.H{"list": list, "mer_id": merID64})
}
