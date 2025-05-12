package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

const (
	// NOTE: SaaS 开发者可根据需求自定义 Limit 上下限
	// Limit 下限
	minLimit = 5
	// Limit 上限
	maxLimit = 100
)

// GetPage 获取分页参数 Page
func GetPage(c *gin.Context) int {
	page := cast.ToInt(c.Query("page"))
	// page 必须为非负整数
	return max(1, page)
}

// GetLimit 获取分页参数 Limit
func GetLimit(c *gin.Context) int {
	limit := cast.ToInt(c.Query("limit"))
	limit = min(maxLimit, limit)
	limit = max(minLimit, limit)
	return limit
}

// GetOffset 获取分页参数 Offset
func GetOffset(c *gin.Context) int {
	page, limit := GetPage(c), GetLimit(c)
	return (page - 1) * limit
}
