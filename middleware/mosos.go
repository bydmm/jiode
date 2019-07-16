package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// JSONLog JSON格式的log
type JSONLog struct {
	Time   string `json:"time"`
	Status int    `json:"status"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Cost   int64  `json:"cost"`
	IP     string `json:"ip"`
}

// Message 消息
type Message struct {
	Client string `json:"c"`
	Msg    string `json:"m"`
}

// MososLogger 日志格式化函数
func MososLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log := JSONLog{
			Time:   end.Format("2006/01/02 - 15:04:05"),
			Status: statusCode,
			Cost:   int64(latency / 1e6),
			Method: method,
			Path:   path,
			IP:     clientIP,
		}
		json, _ := json.Marshal(log)
		go func(json string) {
			SendToMosos(json)
		}(string(json))
	}
}

// SendToMosos 发送给mosos
func SendToMosos(payload string) {
	addr := os.Getenv("MOSOS_ADDR")
	room := os.Getenv("MOSOS_ROOM")
	token := os.Getenv("MOSOS_SECRET_TOKEN")

	u := url.URL{
		Scheme: "http",
		Host:   addr,
		Path:   fmt.Sprintf("api/%s/%s", token, room),
	}

	message := Message{
		Client: "gin server",
		Msg:    payload,
	}
	body, _ := json.Marshal(message)
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(body))

	req.Header.Add("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
}
