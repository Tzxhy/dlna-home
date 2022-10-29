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

	v1.GET("get-status", controllers.GetStatus)
	v1.GET("get-position", controllers.GetPosition)
	v1.POST("set-position", controllers.SetPosition)
	v1.GET("device-list", controllers.GetDeviceList)
	v1.GET("playlist", controllers.GetPlayList)
	v1.POST("delete-playlist", controllers.DeletePlayList)
	v1.POST("create-playlist", controllers.CreatePlayList)

	v1.POST("update-playlist", controllers.SetPlayList)
	v1.POST("update-partial-list", controllers.AddPartialListForPlay)
	v1.POST("delete-single-resource", controllers.DeleteSingleResource)
	v1.POST("rename-playlist", controllers.RenamePlayList)

	v1.POST("action", controllers.Action)
	v1.POST("start-one", controllers.StartOne)
	v1.GET("volume", controllers.GetDeviceVolume)
	v1.POST("volume", controllers.SetDeviceVolume)

	return r
}
