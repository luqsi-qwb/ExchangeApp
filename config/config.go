package config

import (
	"log"

	"github.com/spf13/viper"
)

// 使用viper读取配置信息
type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Dns          string
		MaxIdleConns int
		MaxOpenCons  int
	}
}

var AppConfig *Config //全局变量

func InitConfig() {
	//初始化viper
	viper.SetConfigName("config")              //配置文件的名字
	viper.SetConfigType("yml")                 //配置文件的类型
	viper.SetConfigFile("./config/config.yml") //配置文件的路径

	//读取配置信息
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error ReadInConfig,err is %v", err)

	}

	AppConfig = &Config{}

	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Error Unmarshal,err is %v", err)
	}

	initDb()
	InitRedis()
}
