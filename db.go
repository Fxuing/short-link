package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

/**
初始化数据库，并建立连接
*/
var DB *gorm.DB

func InitDB() (db *gorm.DB, err error) {
	connect :=fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",viper.Get("db.username"),viper.Get("db.password"),viper.Get("db.host"),viper.GetString("db.port"),viper.Get("db.database"))
	db, err = gorm.Open("mysql", connect)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	db.AutoMigrate(&ShortLink{})
	db.LogMode(true)
	DB = db
	return db, err
}

type ShortLink struct {
	gorm.Model
	ShortUrl string
	LongUrl  string
}
