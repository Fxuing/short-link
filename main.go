package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func init() {
	viper.SetConfigName("config")
	// 设置配置文件和可执行二进制文件在用一个目录
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("no such config file")
		} else {
			// Config file was found but another error was produced
			log.Println("read config error")
		}
		log.Fatal(err)
	}
}

func main() {

	var db, err = InitDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	router := gin.Default()
	router.Use(Redirect(), gin.Recovery())
	slink :=router.Group("slink")
	{
		slink.POST("/short", CreateShort)
	}
	router.Run(fmt.Sprintf(":%s", viper.GetString("server.port")))

}

func Redirect() gin.HandlerFunc {
	return func(context *gin.Context) {
		url := context.Request.URL
		var short ShortLink
		short.ShortUrl = url.String()
		err := DB.Find(&short, &short).Error
		if err != nil {
			fmt.Println(err)
		}
		if short.LongUrl != "" {
			context.Redirect(http.StatusMovedPermanently, short.LongUrl)
		}
	}

}
