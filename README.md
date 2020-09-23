## Go 实现短链接服务

        我们在推广的时候，经常会碰到链接或者二维码的时候，这时候链接太长，看起来不太友好，生成的二维码也是密密麻麻的，所以会用到短链接服务。

### 主要思路

1. 直接生成短链接，短链接通过当前时间戳转换`Base62`

2. 将长链接和短链接数据保存到数据库

3. 生成的时候，先查看下链接有没有生成过，如果生成过，用之前生成的短链即可

4. 生成之后，需要写个中间件做预处理，当前url是否为短链，如果是短链，直接重定向到  长链接去。

5. 服务编写好后，`Nginx`加一段配置 ，针对短链转发到短链接服务。

### 使用到的技术

- [gin](https://github.com/gin-gonic/gin) web框架

- [gorm](https://github.com/jinzhu/gorm) orm框架

- [viper ](https://github.com/spf13/viper) 读取配置https://github.com/catinello/base62

- [base62 ](https://github.com/catinello/base62) 生成base62的工具包

### 核心代码

生成短链：

```go
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
```

请求预处理重定向：

`router.Use(Redirect(), gin.Recovery())`

```go
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
```

Nginx 配置:

```nginx
location "~/slink/([a-z]|[A-Z]|[0-9]){7,8}$" {
¦   proxy_pass http://127.0.0.1:7788;
}   
location /slink/short {
¦   proxy_pass http://127.0.0.1:7788/short;
¦   proxy_redirect off;
¦   proxy_set_header Host $host;
¦   proxy_set_header X-Real-IP $remote_addr;
¦   proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
¦   proxy_set_header X-Forwarded-Proto $scheme;
¦   client_max_body_size 30m;
¦   client_body_buffer_size 256k;
¦   proxy_connect_timeout   90; 
¦   proxy_send_timeout      180;
¦   proxy_read_timeout      180;
¦   proxy_buffer_size       256k;
¦   proxy_buffers           16 256k;
¦   proxy_busy_buffers_size 1024k;
¦   proxy_temp_file_write_size      1024k;
}
```

项目地址：https://github.com/Fxuing/short-link


