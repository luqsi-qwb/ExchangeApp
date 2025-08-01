package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuqsi/exchangeapp/utils"
)

//中间件：路由拦截器，让他不要直接进入到处理函数当中

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing Token",
			})
			ctx.Abort()
			return
		}

		username, err := utils.ParseJwt(token)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}
		//存储一下键值对，方便中间件的使用
		ctx.Set("username", username)
		ctx.Next() //进入下一个中间件处理
	}
}
