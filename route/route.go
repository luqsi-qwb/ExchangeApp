package route

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kuqsi/exchangeapp/controllers"
	"github.com/kuqsi/exchangeapp/middlewares"
)

func SetupRouter() *gin.Engine {
	//
	//设置默认路由
	r := gin.Default()

	//允许跨域访问
	r.Use(cors.New(cors.Config{
		//配置开始
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://192.168.189.1:5173"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//认证路由注册
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)

		auth.POST("/register", controllers.Register)
	}
	//获取货币交换律路由
	api := r.Group("/api")
	api.GET("/exchangeRates", controllers.GetExchangeRate)

	api.Use(middlewares.AuthMiddleware())
	{
		api.POST("/exchangeRates", controllers.CreateExchangeRate)

		api.POST("/articles", controllers.CreateArticle)
		api.GET("/articles", controllers.GetArticles)
		api.GET("/articles/:id", controllers.GetArticleById)

		api.POST("/articles/:id/like", controllers.LikeArticle)
		api.GET("articles/:id/like", controllers.GetArticleLikes)
	}

	return r
}
