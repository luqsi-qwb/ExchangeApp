package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuqsi/exchangeapp/global"
	"github.com/kuqsi/exchangeapp/models"
	"github.com/kuqsi/exchangeapp/utils"
)

func Register(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//对密码进行加密处理
	hashpwd, err := utils.HashPassward(user.Passward)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	user.Passward = hashpwd

	//生成token
	token, err := utils.CreateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//将这个表自动迁移到数据库上面，如果没有该表则自动生成
	if err := global.Db.AutoMigrate(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//将该用户插入表中
	if err := global.Db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func Login(ctx *gin.Context) {
	var input struct {
		User     string `json:"username"`
		Passward string `json:"passward"`
	}

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User

	if err := global.Db.Where("username = ?", input.User).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "找不到该用户",
		})
		return
	}

	//判断密码是否正确
	if err := utils.CheckPassward(input.Passward, user.Passward); !err {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "密码出错",
		})
		return
	}
	//密码正确返回一个token
	token, err := utils.CreateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
