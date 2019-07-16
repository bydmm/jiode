package server

import (
	"gopkg.in/olahol/melody.v1"
)

// Melody 实例
var Melody *melody.Melody

// clients 客户端列表
var clients *UserMap

// InitClients 初始化在线列表
func InitClients() {
	clients = BuildUserMap()
}

// MelodyInit 初始化Melody
func MelodyInit() {
	Melody = melody.New()
	InitClients()
}
