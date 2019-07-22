package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"github.com/bydmm/jiode/model"
	"github.com/gin-gonic/gin"
)

// JSONDumpLog JSON格式的log
type JSONDumpLog struct {
	Time    string          `json:"time"`
	Status  int             `json:"status"`
	Method  string          `json:"method"`
	Path    string          `json:"path"`
	Cost    int64           `json:"cost"`
	IP      string          `json:"ip"`
	ReqBody json.RawMessage `json:"req_body"`
	ResBody json.RawMessage `json:"res_body"`
}

// JSONDump 日志格式化函数
func JSONDump() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		reqBody := json.RawMessage(`{}`)
		if c.Request.ContentLength > 0 {
			bodyCopy := new(bytes.Buffer)
			// Read the whole body
			io.Copy(bodyCopy, c.Request.Body)
			bodyData := bodyCopy.Bytes()
			// Replace the body with a reader that reads from the buffer
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(bodyData))
			data, _ := ioutil.ReadAll(bodyCopy)
			// var v interface{}
			// json.Unmarshal(data, &v)
			// body, _ := json.Marshal(v)
			reqBody = json.RawMessage(string(data))
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		resBody := json.RawMessage(string(blw.body.String()))

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log := JSONDumpLog{
			Time:    end.Format("2006/01/02 - 15:04:05"),
			Status:  statusCode,
			Cost:    int64(latency / 1e6),
			Method:  method,
			Path:    path,
			IP:      clientIP,
			ReqBody: reqBody,
			ResBody: resBody,
		}

		json, _ := json.Marshal(log)
		go func(json string) {
			model.Send(json)
		}(string(json))
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
