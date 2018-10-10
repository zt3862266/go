package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

//used for trace request
func Tracecode() gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawPath
		c.Set("LOGID", c.GetHeader("HTTP_X_RONG_LOGID"))
		c.Set("REQID", c.GetHeader("HTTP_X_RONG_REQID"))

		//before request
		c.Next()
		//after request

		end := time.Now()
		latency := end.Sub(start)
		clientIp := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		logId, _ := c.Get("LOGID")
		reqId, _ := c.Get("REQID")
		if raw != "" {
			path = path + "?" + raw
		}
		fmt.Fprintf(gin.DefaultWriter, "%v | %3d | %13v | %15s | %-7s | %26s | %26s | %s\n",
			end.Format("2006-10-02 15:04:05"),
			statusCode,
			latency,
			clientIp,
			method,
			logId,
			reqId,
			path,
		)

	}

}
