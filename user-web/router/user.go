package router

import (
	"LFshop-api/user-web/api"
	"LFshop-api/user-web/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RouterGroup类似一个handler的切片
func InitUserRouter(Router *gin.RouterGroup) {
	zap.S().Info("配置user相关的url")
	UserRouter := Router.Group("user").Use(middlewares.JWTAuth())
	{
		// 加一个参数，设置某些方法需要登录才能访问
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
		UserRouter.POST("register", api.Register)
	}
}
