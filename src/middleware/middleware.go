package middleware

import (
	token "task-tracker-api/src/utils"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := token.TokenValid(ctx)

		if err != nil {
			ctx.JSON(401, gin.H{
				"message": "User is unauthorized",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
