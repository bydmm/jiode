package server

import (
	"sync"

	melody "gopkg.in/olahol/melody.v1"
)

// UserMap 使用hashmap快速查找用户
type UserMap struct {
	lock  *sync.Mutex
	users map[*melody.Session]*User
}

// BuildUserMap 初始化
func BuildUserMap() *UserMap {
	users := UserMap{
		lock:  new(sync.Mutex),
		users: make(map[*melody.Session]*User),
	}
	return &users
}

// GetUsers 获取用户列表
func (userMap *UserMap) GetUsers() map[*melody.Session]*User {
	return userMap.users
}

// GetUser 获取用户
func (userMap *UserMap) GetUser(s *melody.Session) *User {
	return userMap.users[s]
}

// Add 添加用户
func (userMap *UserMap) Add(user *User) {
	userMap.lock.Lock()
	userMap.users[user.Session] = user
	userMap.lock.Unlock()
}

// Delete 移出用户
func (userMap *UserMap) Delete(s *melody.Session) {
	userMap.lock.Lock()
	delete(userMap.users, s)
	userMap.lock.Unlock()
}

// Count 数量
func (userMap *UserMap) Count() int {
	count := 0
	for range userMap.users {
		count++
	}
	return count
}
