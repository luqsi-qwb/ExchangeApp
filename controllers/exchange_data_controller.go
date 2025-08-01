package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuqsi/exchangeapp/global"
	"github.com/kuqsi/exchangeapp/models"
)

func CreateExchangeRate(ctx *gin.Context) {
	var rate models.ExchangeRate
	//提取客户端发来的数据
	if err := ctx.ShouldBindBodyWithJSON(&rate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//填充一下时间信息
	rate.Date = time.Now()

	//初始化数据库表
	if err := global.Db.AutoMigrate(&rate); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	//插入数据
	if err := global.Db.Create(&rate).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	ctx.JSON(http.StatusCreated, rate)

}

func GetExchangeRate(ctx *gin.Context) {
	var rates []models.ExchangeRate

	if err := global.Db.Find(&rates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, rates)
}
