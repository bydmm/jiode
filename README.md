# jiode(脚得)

在你本机的命令行直接查看远端服务器日志的利器

## jiode是做什么的

总感觉前后端有扯不尽的皮，丢不完的锅。每天登陆跳板机，看个日志还要输入验证码真的好累啊！！！

没问题，使用jiode，别你jio得你传了，直接在本机命令行看远端日志吧。

## 效果演示
```bash
$ docker run bydmm/jiode -c -addr 120.x.x.x:7788 -room dev -token password
2019/07/24 05:08:30 connecting to ws://120.x.x.x:7788/ws/password/join
========================== 2019/07/24 - 13:05:01 内网api-dev ==========================
 {
  "cost": 0,
  "ip": "192.168.80.1",
  "method": "GET",
  "path": "/api/v1/blog/my?timestamp=1563944701220",
  "req_body": {},
  "res_body": {
    "code": 0,
    "data": {
      "templates": []
    },
    "result": ""
  },
  "status": 200,
  "time": "2019/07/24 - 05:05:01"
}
```

## Docker启动方法

```shell
# 启动日志总收集服务器，为日志生产者和日志查看者服务
# JIODE_SECRET_TOKEN是密钥，作为接入这个服务的密码，避免滥用
# -p 9999:3000 原服务启动在3000端口，然后我把他暴露在9999端口，记得开防火墙
docker run -e JIODE_SECRET_TOKEN="123456" -p 9999:3000 bydmm/jiode
```

```shell
# 本机cli启动方法
# -addr 设置远端日志服务器的ip和端口
# -token 是上面的密码，要一致
# room 是日志频道，具体看下面日志生产者怎么设置的
docker run bydmm/jiode -c -addr 127.0.0.1:9999 -token 778899 -room room1
```

## 发送日志

```shell
# JIODE需要你配置以下环境变量，具体设置方法请根据你的实际情况决定
# windows用set, linux用export, docker用-e等等
JIODE_ADDR="127.0.0.1:5000" # jiode服务器地址
JIODE_ROOM="dev" # 房间名，自己随便定
JIODE_SECRET_TOKEN="778899" # 密钥
JIODE_SERVICE_NAME="room1" # 发送日志的服务名称，自定义
```

```golang
# jiode实现了gin框架的日志发送中间件，如下直接使用
import jiode "github.com/bydmm/jiode/middleware"

func main() {
    r := gin.Default()
	// 生产环境发日志流量太大
	if gin.Mode() != "release" {
		r.Use(jiode.JSONDump())
	}
	r.GET("/check", func(c *gin.Context){
		c.string(200, "ok")
	})
    r.Run(":3000")
}
```

#### 如果你有定制需求，也可以自行调接口发送

```
POST localhost:5000/api/:token/:room
{
	"c": "server 1", # 发送端名称
	"m": "{"time":"2019/07/16 - 20:30:29","status":200,"method":"GET","path":"/api/v1/videos","cost":0,"ip":"::1"}" # 发送日志
}
```

## 自行安装方法

```shell
git clone https://github.com/bydmm/jiode
go install
export JIODESECRET_TOKEN="abcdef" # 设置接入token, 免得别人乱来。如果不设置也可以，jiode会随机给个otken
jiode # 启动jiode接受服务器，一般在公网云端启动
```

