package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		
		key := os.Getenv("ADMIN_SECRET")

		auth := ctx.Request.Header.Get("Authorization");

		if auth == "" {
			ctx.JSON(401, gin.H{"error":"Akses tidak diizinkan"})
			ctx.Abort()
			return
		}

		if auth != key {
			ctx.JSON(401, gin.H{"error":"Akses tidak diizinkan"})
			ctx.Abort()
			return
		}

		ctx.Next()

	}
}