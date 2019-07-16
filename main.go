package main

import (
	"mosos/server"

	"github.com/joho/godotenv"
)

func main() {
	// 从本地读取环境变量
	godotenv.Load()

	server.RunServer()
}
