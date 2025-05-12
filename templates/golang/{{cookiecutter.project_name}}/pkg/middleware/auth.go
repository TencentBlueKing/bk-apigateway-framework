package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/account"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/common"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/utils/ginx"
)

// UserAuth 进行用户身份认证，并将用户信息注入到 context 中
func UserAuth(authBackend account.AuthBackend) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 重定向链接（当前访问的链接）
		scheme := lo.Ternary(c.Request.TLS != nil, "https", "http")
		referUrl := fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, c.Request.RequestURI)
		ginH := gin.H{"loginUrl": fmt.Sprintf("%s?c_url=%s", authBackend.GetLoginUrl(), referUrl)}

		userToken, err := c.Request.Cookie(common.UserTokenKey)
		// 没有获取到用户凭证信息 -> 401 -> 让用户通过页面登录
		if err != nil {
			c.HTML(http.StatusUnauthorized, "401.html", ginH)
			c.Abort()
			return
		}

		session := sessions.Default(c)
		if userToken.Value == session.Get(common.UserTokenKey) {
			// 从 session 获取用户信息并注入到 context
			ginx.SetUserID(c, session.Get(common.UserIDKey).(string))
			c.Next()
			return
		}

		userInfo, err := authBackend.GetUserInfo(c.Request.Context(), userToken.Value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		// 获取到用户凭证信息 -> 设置 context & session -> 通过
		ginx.SetUserID(c, userInfo.ID)
		session.Set(common.UserTokenKey, userToken.Value)
		session.Set(common.UserIDKey, userInfo.ID)
		_ = session.Save()
		c.Next()
	}
}
