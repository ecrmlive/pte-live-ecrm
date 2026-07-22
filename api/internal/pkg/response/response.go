package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{Status: 200, Message: "success", Data: data})
}

func Fail(c *gin.Context, httpCode int, message string) {
	if httpCode <= 0 {
		httpCode = http.StatusBadRequest
	}
	c.JSON(httpCode, Body{Status: httpCode, Message: message})
}
