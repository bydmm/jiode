package server

import (
	"fmt"
	"os"

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
	token := os.Getenv("SECRET_TOKEN")

	r := gin.Default()
	r.GET(fmt.Sprintf("/ws/%s/:room", token), func(c *gin.Context) {
		Melody.HandleRequest(c.Writer, c.Request)

		Melody.HandleConnect(func(s *melody.Session) {
			HandleLogin(c.Param("room"), s)
		})

		Melody.HandleMessage(func(s *melody.Session, msg []byte) {
			formUser := clients.GetUser(s)
			Melody.BroadcastFilter(msg, func(s *melody.Session) bool {
				toUser := clients.GetUser(s)
				return toUser != nil && toUser.room == formUser.room
			})
		})
	})

	r.Run(":5000")
}
