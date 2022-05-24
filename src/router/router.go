package router

import (
	controller "task-tracker-api/src/controller"
	"task-tracker-api/src/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin := gin.Default()

	public := gin.Group("/api")

	public.POST("/register", controller.Register)
	public.POST("/login", controller.Login)

	authorized := gin.Group("/api")

	authorized.Use(middleware.JwtAuthMiddleware())
	{
		authorized.POST("/change-password", controller.ChangePassword)
		authorized.POST("/send-code", controller.SendCode)
		authorized.POST("/verify-code", controller.VerifyCode)
		authorized.GET("/tasks", controller.GetTasks)
		authorized.GET("/tasks/:id", controller.GetTaskById)
		authorized.POST("/tasks", controller.CreateTask)
		authorized.DELETE("/tasks/:id", controller.DeleteTask)
		authorized.PUT("/tasks/:id", controller.UpdateTask)
	}

	gin.Run(":3030")
}
