package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Message 消息
type Message struct {
	Client string `json:"c"`
	Msg    string `json:"m"`
}

// Send 发送给jiode server
func Send(payload string) {
	addr := os.Getenv("JIODE_ADDR")
	room := os.Getenv("JIODE_ROOM")
	token := os.Getenv("JIODE_SECRET_TOKEN")
	name := os.Getenv("JIODE_SERVICE_NAME")

	u := url.URL{
		Scheme: "http",
		Host:   addr,
		Path:   fmt.Sprintf("api/%s/%s", token, room),
	}

	message := Message{
		Client: name,
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
