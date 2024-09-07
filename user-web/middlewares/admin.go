package middlewares

import (
	"LFshop-api/user-web/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 返回一个函数，中间件标准写法
// 加了这个中间件，会检查用户是否有权限调这个方法
func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)
		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
