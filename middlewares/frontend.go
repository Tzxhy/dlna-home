package middlewares

import (
	"io/ioutil"
	"net/http"
	"strings"

	"gitee.com/tzxhy/dlna-home/initial"
	"github.com/gin-gonic/gin"
)

// FrontendFileHandler 前端静态文件处理
func FrontendFileHandler() gin.HandlerFunc {

	// 读取index.html
	file, _ := initial.StatikFS.Open("/index.html")

	fileContentBytes, _ := ioutil.ReadAll(file)

	fileContent := string(fileContentBytes)

	fileServer := http.FileServer(initial.StatikFS)
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// API 跳过
		if strings.HasPrefix(path, "/api") {
			c.Next()
			return
		}
		_, err := initial.StatikFS.Open(path)
		// 不存在的路径和index.html均返回index.html
		if (path == "/index.html") || (path == "/") || err != nil {

			c.Header("Content-Type", "text/html")
			c.String(200, fileContent)
			c.Abort()
			return
		}

		// 存在的静态文件
		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
