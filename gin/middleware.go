package gin

import (
	"github.com/akshaykumar12527/yaag/middleware"
	"github.com/akshaykumar12527/yaag/yaag"
	"github.com/akshaykumar12527/yaag/yaag/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http/httptest"
	"strings"
)

type ResponseWriter struct {
	gin.ResponseWriter
	Data []byte
}

func (w *ResponseWriter) Write(buf []byte) (int, error) {
	w.Data = buf
	return len(buf), nil
}

func Document() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !yaag.IsOn() {
			return
		}
		writer := httptest.NewRecorder()
		apiCall := models.ApiCall{}
		middleware.Before(&apiCall, c.Request)
		w := c.Writer
		resp := &ResponseWriter{c.Writer, []byte{}}
		c.Writer = resp
		c.Next()
		w.Write(resp.Data)
		if writer.Code != 404 {
			headers := map[string]string{}
			for k, v := range c.Writer.Header() {
				log.Println(k, v)
				headers[k] = strings.Join(v, " ")
			}
			if strings.Contains(headers["Content-Type"], "application/json") {
				apiCall.MethodType = c.Request.Method
				apiCall.CurrentPath = strings.Split(c.Request.RequestURI, "?")[0]
				apiCall.ResponseBody = string(resp.Data)
				apiCall.ResponseCode = c.Writer.Status()
				for k, v := range c.Writer.Header() {
					log.Println(k, v)
					headers[k] = strings.Join(v, " ")
				}
				apiCall.ResponseHeader = headers
				go yaag.GenerateHtml(&apiCall)
			}
		}
	}
}
