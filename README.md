# jiode

便利的日志系统

## 启动方法

```
git clone https://github.com/bydmm/jiode
go install
export JIODESECRET_TOKEN="abcdef" # 设置接入token, 免得别人乱来。如果不设置也可以，jiode会随机给个otken
jiode -s # 启动jiode接受服务器，一般在公网云端启动
```

```
jiode -addr localhost:5000 -room dev1 -token abcdef # 在本机启动客户端，开始接受日志
```

## 发送日志

jiode实现了gin框架的日志发送中间件，如下直接使用

```
import jiode "github.com/bydmm/jiode/middleware"

func main() {
    r := gin.Default()
	r.Use(jiode.jiodeLogger())
    r.Run(":3000")
}
```

如果你有定制需求，也可以自行调接口发送

```
POST localhost:5000/api/:token/:room
{
	"c": "server 1", # 发送端名称
	"m": "{"time":"2019/07/16 - 20:30:29","status":200,"method":"GET","path":"/api/v1/videos","cost":0,"ip":"::1"}" # 发送日志
}
```
