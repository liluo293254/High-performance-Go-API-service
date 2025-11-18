package router

import (
	"bgame/internal/handler/user"
	"bgame/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupUserRoutes(r *gin.Engine) {
	userHandler := user.NewUserHandler()
	userGroup := r.Group("/api/user")
	{
		// 公开接口
		userGroup.POST("/regAndLogin", userHandler.RegAndLogin)

		// 需要认证的接口
		userGroup.Use(middleware.AuthUser())
		{
			userGroup.GET("/info", userHandler.GetUserInfo)
		}
	}
}
