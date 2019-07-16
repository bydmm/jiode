package server

import (
	melody "gopkg.in/olahol/melody.v1"
)

// User 用户
type User struct {
	login   bool
	room    string
	Session *melody.Session
}

// SetRoom 设置用户所在房间
func (user *User) SetRoom(room string) {
	user.room = room
}

// Room 用户所在地点
func (user User) Room() string {
	return user.room
}
