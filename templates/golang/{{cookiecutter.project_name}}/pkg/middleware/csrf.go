package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
)

// CSRF 中间件用于防止跨站请求伪造
func CSRF(appID string, secret string) gin.HandlerFunc {
	return adapter.Wrap(
		csrf.Protect([]byte(secret), csrf.Secure(false), csrf.Path("/"), csrf.CookieName(appID+"-csrf")),
	)
}

// CSRFToken 中间件用于在 cookie 中设置 csrf token
func CSRFToken(appID string, domain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetSameSite(http.SameSiteLaxMode)
		// 参数依次为： cookie 名称，值，有效期（0 表示会话 cookie）
		// 路径（根），域名（ "" 表示当前域），是否只通过 https 访问，httpOnly 属性
		c.SetCookie(appID+"-csrf-token", csrf.Token(c.Request), 0, "/", domain, false, false)
	}
}
