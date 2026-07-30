package assist

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/assist"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/trade"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
)

type Handler struct {
	svc   *assist.Service
	trade *trade.Service
}

func NewHandler(svc *assist.Service, tradeSvc *trade.Service) *Handler {
	return &Handler{svc: svc, trade: tradeSvc}
}

func (h *Handler) RegisterPublic(r gin.IRoutes) {
	r.GET("/assist/actives", h.List)
	r.GET("/assist/actives/:id", h.Get)
	r.GET("/assist/actives/:id/sets", h.ListSets)
	r.GET("/assist/sets/:id", h.GetSet)
}

func (h *Handler) RegisterAuthed(r gin.IRoutes) {
	r.GET("/assist/actives/:id/mine", h.Mine)
	r.POST("/assist/actives/:id/start", h.Start)
	r.POST("/assist/sets/:id/help", h.Help)
	r.POST("/order/assist/check", h.Check)
	r.POST("/order/assist/create", h.Create)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListApp(c.Request.Context(), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) ListSets(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	rows, err := h.svc.ListSets(c.Request.Context(), uint(id), 20)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) GetSet(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.GetSet(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Mine(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.FindMine(c.Request.Context(), middleware.UID(c), uint(id))
	if err != nil {
		if errors.Is(err, assist.ErrNotFound) {
			response.OK(c, nil)
			return
		}
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Start(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.StartSet(c.Request.Context(), middleware.UID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Help(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.Help(c.Request.Context(), middleware.UID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Check(c *gin.Context) {
	var in trade.AssistInput
	if err := c.ShouldBindJSON(&in); err != nil || in.ProductAssistSetID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	res, err := h.trade.AssistCheck(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Create(c *gin.Context) {
	var in trade.AssistInput
	if err := c.ShouldBindJSON(&in); err != nil || in.ProductAssistSetID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	g, err := h.trade.AssistCreate(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, g)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, assist.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, assist.ErrBadParam), errors.Is(err, assist.ErrInactive),
		errors.Is(err, assist.ErrSoldOut), errors.Is(err, assist.ErrSetNotOpen),
		errors.Is(err, assist.ErrAlreadyHelped), errors.Is(err, assist.ErrSelfHelp),
		errors.Is(err, assist.ErrSetClosed), errors.Is(err, assist.ErrForbidden):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}

func writeTradeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, trade.ErrAddressRequired), errors.Is(err, trade.ErrBadParam),
		errors.Is(err, trade.ErrStockNotEnough),
		errors.Is(err, assist.ErrSetNotOpen), errors.Is(err, assist.ErrInactive),
		errors.Is(err, assist.ErrSoldOut), errors.Is(err, assist.ErrForbidden),
		errors.Is(err, assist.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, trade.ErrNotFound), errors.Is(err, assist.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
