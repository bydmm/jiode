package main

import (
	"flag"

	"github.com/bydmm/jiode/client"
	"github.com/bydmm/jiode/server"

	"github.com/joho/godotenv"
)

var cli = flag.Bool("c", false, "client mode")
var addr = flag.String("addr", "localhost:5000", "server host")
var room = flag.String("room", "all", "log channel")
var token = flag.String("token", "", "jiode token")

func main() {
	// 从本地读取环境变量
	godotenv.Load()

	flag.Parse()

	if cli != nil && *cli == true {
		client.RunClient(*addr, *room, *token)
	} else {
		server.RunServer()
	}
}
