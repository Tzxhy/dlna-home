package routers

import (
	"gitee.com/tzxhy/dlna-home/controllers"
	"gitee.com/tzxhy/dlna-home/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.Use(middlewares.FrontendFileHandler())
	v1 := r.Group("/api/v1")

	v1.GET("device-list", controllers.GetDeviceList)
	v1.GET("playlist", controllers.GetPlayList)
	v1.POST("delete-playlist", controllers.DeletePlayList)
	v1.POST("create-playlist", controllers.CreatePlayList)

	v1.POST("update-playlist", controllers.SetPlayList)

	v1.POST("action", controllers.Action)

	return r
}
