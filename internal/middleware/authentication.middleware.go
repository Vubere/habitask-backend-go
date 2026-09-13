package middleware

import (
	"habitask-backend-go/pkg/appjwt"

	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		tokenString = tokenString[len("Bearer "):]
		userInfo, err := appjwt.VerifyToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		ctx.Set("userId", userInfo.UserID)
		ctx.Set("username", userInfo.Username)
		ctx.Next()
	}
}
