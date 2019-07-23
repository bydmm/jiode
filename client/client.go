package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/bydmm/jiode/model"

	"github.com/gorilla/websocket"
)

// RunClient 客户端
func RunClient(addr string, room string, token string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{
		Scheme: "ws",
		Host:   addr,
		Path:   fmt.Sprintf("ws/%s/%s", token, room),
	}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var msg model.Message
			json.Unmarshal(message, &msg)

			loc := time.FixedZone("UTC-8", +8*60*60)
			time := time.Now().In(loc).Format("2006/01/02 - 15:04:05")
			var v interface{}
			err = json.Unmarshal([]byte(msg.Msg), &v)
			if err == nil {
				body, _ := json.MarshalIndent(v, "", "  ")
				fmt.Printf("========================== %s %s ==========================\n %s\n", time, msg.Client, body)
			} else {
				fmt.Printf("%s %s\n %s\n", time, msg.Client, msg.Msg)
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
