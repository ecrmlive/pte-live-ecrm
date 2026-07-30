package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
	"gorm.io/gorm"
)

type CallbackHandler struct {
	db        *gorm.DB
	allowMock bool
}

func NewCallbackHandler(db *gorm.DB, allowMock bool) *CallbackHandler {
	return &CallbackHandler{db: db, allowMock: allowMock}
}
func (h *CallbackHandler) RegisterMock(r gin.IRoutes) { r.POST("/pay/mock", h.Mock) }
func (h *CallbackHandler) Mock(c *gin.Context) {
	if !h.allowMock {
		response.Fail(c, http.StatusNotFound, "接口不存在")
		return
	}
	var in struct {
		GroupOrderID uint64 `json:"group_order_id"`
		UID          uint64 `json:"uid"`
	}
	if err := c.ShouldBindJSON(&in); err != nil || in.GroupOrderID == 0 || in.UID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	created, err := PayMock(c.Request.Context(), h.db, in.UID, in.GroupOrderID)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	response.OK(c, createdResponse(created, true))
}
