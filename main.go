package main

import (
	"flag"

	"github.com/bydmm/jiode/client"
	"github.com/bydmm/jiode/server"

	"github.com/joho/godotenv"
)

var serv = flag.Bool("s", false, "server mode")
var addr = flag.String("addr", "localhost:5000", "server host")
var room = flag.String("room", "all", "log channel")
var token = flag.String("token", "", "jiode token")

func main() {
	// 从本地读取环境变量
	godotenv.Load()

	flag.Parse()

	if serv != nil && *serv == true {
		server.RunServer()
	} else {
		client.RunClient(*addr, *room, *token)
	}
}
