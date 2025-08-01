package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/kuqsi/exchangeapp/global"
)

func LikeArticle(ctx *gin.Context) { //增加点赞数
	articleid := ctx.Param("id")

	LikeKey := "article:" + articleid + ":likes"

	if err := global.RedisDb.Incr(LikeKey).Err(); err != nil { //自动增加数量，如果没有该文章则自动创建并且赋值为0
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "点赞数成功加一",
	})

}

func GetArticleLikes(ctx *gin.Context) {
	articleid := ctx.Param("id")

	LikeKey := "article:" + articleid + ":likes"

	likes, err := global.RedisDb.Get(LikeKey).Result()

	if err == redis.Nil {
		likes = "0"
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"likes": likes,
	})
}
