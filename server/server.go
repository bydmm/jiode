package server

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bydmm/jiode/model"
	"github.com/bydmm/jiode/util"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

// HandleLogin 登录处理
func HandleLogin(room string, s *melody.Session) {
	user := User{
		Session: s,
		room:    room,
	}
	clients.Add(&user)
}

// RunServer 运行服务器
func RunServer() {
	MelodyInit()

	token := os.Getenv("JIODE_SECRET_TOKEN")
	if token == "" {
		token = util.RandStringRunes(6)
	}

	r := gin.Default()

	Melody.HandleConnect(func(s *melody.Session) {
		fmt.Println("用户连接")
		HandleLogin("", s)
	})

	Melody.HandleDisconnect(func(s *melody.Session) {
		formUser := clients.GetUser(s)
		fmt.Println("用户退出房间" + formUser.Room())
		clients.Delete(s)
	})

	Melody.HandleMessage(func(s *melody.Session, msg []byte) {
		user := clients.GetUser(s)
		user.SetRoom(string(msg))
		fmt.Println("用户进入房间" + user.Room())
	})

	r.GET(fmt.Sprintf("/ws/%s/join", token), func(c *gin.Context) {
		Melody.HandleRequest(c.Writer, c.Request)
	})

	r.POST(fmt.Sprintf("/api/%s/:room", token), func(c *gin.Context) {
		room := c.Param("room")
		if room == "" {
			return
		}
		var jsonBody model.Message
		c.BindJSON(&jsonBody)
		msg, _ := json.Marshal(jsonBody)

		Melody.BroadcastFilter(msg, func(s *melody.Session) bool {
			toUser := clients.GetUser(s)
			return toUser != nil && toUser.Room() == room
		})

		c.JSON(200, map[string]string{
			"res": "ok",
		})
	})

	r.Run()
}
