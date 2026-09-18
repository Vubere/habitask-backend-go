package middleware

import (
	"fmt"
	"habitask-backend-go/pkg/appjwt"

	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			fmt.Println("fasdfas", "no token")
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		tokenString = tokenString[len("Bearer "):]
		userInfo, err := appjwt.VerifyToken(tokenString)
		if err != nil {
			fmt.Println("fasdfas", err.Error())
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		fmt.Println("here====>", userInfo.UserID, userInfo.Username, "line")
		ctx.Set("userId", userInfo.UserID)
		ctx.Set("username", userInfo.Username)
		ctx.Next()
	}
}
