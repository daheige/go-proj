package middleware

import "github.com/gin-gonic/gin"

func NotFoundHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"message": "this page not found!",
		})
	}
}
