package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassward(pwd string) (string, error) {
	hashpwd, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	if err != nil {
		return "", err
	}
	return string(hashpwd), nil
}

func CreateJWT(username string) (string, error) { //创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	Token, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return "Bearer " + Token, nil
}

func CheckPassward(pwd string, hashpwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashpwd), []byte(pwd))
	return err == nil
}

func ParseJwt(tokenString string) (string, error) { //解析一下token
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpexted Signing Methed")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("token is un vaild")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("无法解析声明")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("找不到用户名")
	}
	return username, nil
}
