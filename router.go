package main

import (
	"fmt"
	"github.com/catinello/base62"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Response gin.H

const (
	S_LINK = "slink"
)
func CreateShort(c *gin.Context) {
	content := c.PostForm("content")
	shortUrl := generateShort(content)

	c.JSON(http.StatusOK, Response{
		"code": "200",
		"msg":  "成功",
		"data": Response{
			"short": viper.GetString("short.prefix") + shortUrl,
		},
	})
}

func generateShort(longUrl string) string {

	var short ShortLink
	short.LongUrl = longUrl
	err := DB.Find(&short, &short).Error
	if err != nil {
		fmt.Println(err)
	}
	if short.ShortUrl != "" {
		return short.ShortUrl
	}
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	sb.WriteString("/")
	sb.WriteString(S_LINK)
	sb.WriteString("/")
	timestamp := time.Now().UnixNano() / 1e6
	sb.WriteString(base62.Encode(int(timestamp)))
	shortUrl := sb.String()
	shortInfo := ShortLink{
		ShortUrl: shortUrl,
		LongUrl:  longUrl,
	}
	DB.Create(&shortInfo)
	return shortUrl
}
