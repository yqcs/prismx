package httpResult

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Resp struct {
	Status  string `json:"status"`  // 状态
	Code    int    `json:"code"`    // 状态码
	Data    any    `json:"data"`    // 数据集
	Message string `json:"message"` // 消息
}

// FAILURE 失败数据处理
func FAILURE(c *gin.Context, code int, message string) {
	c.JSON(code, Resp{
		Code:    code,
		Status:  "failure",
		Message: message,
		Data:    nil,
	})
}

// SUCCESS 通常成功数据处理
func SUCCESS(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Resp{
		Code:    http.StatusOK,
		Data:    data,
		Status:  "success",
		Message: "请求成功",
	})
}
