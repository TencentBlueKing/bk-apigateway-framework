package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NOTE: SaaS 开发者可根据需要添加 Code 字段用于标记错误类型

// Response 通用响应体
type Response struct {
	Message   string `json:"message"`
	Data      any    `json:"data"`
	RequestID string `json:"requestID"`
}

// SetResp 为指定的 gin.Context 设置成功响应数据（建议 200 <= statusCode < 300）
func SetResp(c *gin.Context, statusCode int, data any) {
	// 204 状态码特殊处理
	if statusCode == http.StatusNoContent {
		c.Status(statusCode)
		return
	}
	c.JSON(statusCode, Response{Message: "", Data: data, RequestID: GetRequestID(c)})
}

// SetErrResp 为指定的 gin.Context 设置错误响应数据
func SetErrResp(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{Message: message, Data: nil, RequestID: GetRequestID(c)})
}

// PaginatedResp 分页响应数据体
type PaginatedResp struct {
	Count   int64 `json:"count"`
	Results any   `json:"results"`
}

// NewPaginatedRespData 创建分页响应数据体
// 注意：results 类型应该是 Slice / Array
func NewPaginatedRespData(count int64, results any) PaginatedResp {
	return PaginatedResp{Count: count, Results: results}
}
