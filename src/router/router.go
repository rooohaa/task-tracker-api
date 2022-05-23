package router

import (
	auth "task-tracker-api/src/controller"
	"task-tracker-api/src/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin := gin.Default()

	public := gin.Group("/api")

	public.POST("/register", auth.Register)
	public.POST("/login", auth.Login)

	authorized := gin.Group("/api")

	authorized.Use(middleware.JwtAuthMiddleware())
	{
		authorized.POST("/change-password", auth.ChangePassword)
		authorized.POST("/send-code", auth.SendCode)
		authorized.POST("/verify-code", auth.VerifyCode)
		// add here endpoints
	}

	gin.Run(":3030")
}
