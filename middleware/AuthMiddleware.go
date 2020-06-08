package middleware

import (
	"micor/ginessential/common"
	"micor/ginessential/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 进行授权认证
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")
		//vcalidate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]                        // 提取Token的有效部分
		token, claims, err := common.ParseToken(tokenString) // 对Token的有效部分进行解析
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		//验证通过后获取Claiim中的userId
		userID := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userID)

		//用户不存在
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}
		//用户存在，将user信息写入上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}
