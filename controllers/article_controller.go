package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/kuqsi/exchangeapp/global"
	"github.com/kuqsi/exchangeapp/models"
	"gorm.io/gorm"
)

var cacheKey = "articles"

func CreateArticle(ctx *gin.Context) {
	var article models.Article

	if err := ctx.ShouldBindBodyWithJSON(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//绑定数据库
	if err := global.Db.AutoMigrate(&article); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	//添加数据
	if err := global.Db.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		fmt.Println("2出错")
		return
	}

	//每当有新文章加入的时候就删除掉缓存中的数据
	if err := global.RedisDb.Del(cacheKey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		fmt.Println("3出错")
		return
	}

	ctx.JSON(http.StatusCreated, article)
}

func GetArticles(ctx *gin.Context) { //采取旁路缓存模式
	cacheData, err := global.RedisDb.Get(cacheKey).Result()
	if err == redis.Nil { //在redis中没有找到
		var articles []models.Article

		if err := global.Db.Find(&articles).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "找不到文章",
				})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
			}
			return
		}
		//将数据放入缓存当中
		articlesJson, err := json.Marshal(articles)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}

		err = global.RedisDb.Set(cacheKey, articlesJson, 10*time.Minute).Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}
		ctx.JSON(http.StatusOK, articles)
		return

	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	//说明找到数据
	var articles []models.Article

	if err := json.Unmarshal([]byte(cacheData), &articles); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, articles)
}

func GetArticleById(ctx *gin.Context) {
	id := ctx.Param("id")
	fmt.Println("id is ", id)
	var article models.Article
	if err := global.Db.Where("id = ?", id).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "找不到该文章",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, article)
}
