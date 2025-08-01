package config

import (
	"log"
	"time"

	"github.com/kuqsi/exchangeapp/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDb() {
	//初始化数据库函数
	dns := AppConfig.Database.Dns
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败,err is %v", err)
	}

	//获取一下数据库连接池
	sqldb, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库连接池失败,err is %v", err)
	}
	//配置一下数据库连接池
	sqldb.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)
	sqldb.SetMaxOpenConns(AppConfig.Database.MaxOpenCons)
	sqldb.SetConnMaxLifetime(time.Hour)


	global.Db = db

}
