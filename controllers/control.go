package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"gitee.com/tzxhy/dlna-home/devices"
	"gitee.com/tzxhy/dlna-home/httphandlers"
	"gitee.com/tzxhy/dlna-home/models"
	"gitee.com/tzxhy/dlna-home/soapcalls"
	"gitee.com/tzxhy/dlna-home/utils"
	"github.com/gin-gonic/gin"
)

func GetDeviceList(c *gin.Context) {
	deviceList, err := devices.LoadSSDPServices(1)
	if err != nil {
		var empty = make(map[string]string)
		c.JSON(http.StatusOK, &gin.H{
			"data": empty,
		})
		return
	}
	c.JSON(http.StatusOK, &gin.H{
		"data": deviceList,
	})
}

type PlayListItem struct {
	models.PlayListItem
	List []models.AudioItem `json:"list"`
}

func GetPlayList(c *gin.Context) {
	list, _ := models.GetPlayList()
	var newList []PlayListItem
	for _, item := range *list {
		items := models.GetPlayListItems(item.Pid)
		newList = append(newList, PlayListItem{
			item,
			*items,
		})
	}
	c.JSON(http.StatusOK, &gin.H{
		"list": newList,
	})
}

type DeletePlayListReq struct {
	Pid string `json:"pid" binding:"required"`
}

func DeletePlayList(c *gin.Context) {
	var deletePlayListReq DeletePlayListReq
	if c.ShouldBind(&deletePlayListReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	models.DeletePlayList(deletePlayListReq.Pid)
}

type GetDeviceVolumeReq struct {
	RendererUrl string `json:"renderer_url" form:"renderer_url" binding:"required"`
}

func GetDeviceVolume(c *gin.Context) {
	var getDeviceVolumeReq GetDeviceVolumeReq
	if c.ShouldBind(&getDeviceVolumeReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	rendererUrl := getDeviceVolumeReq.RendererUrl

	tv, ok := serverMap[rendererUrl]
	if ok {
		val, err := tv.GetVolumeSoapCall()
		if err != nil {
			log.Println("err: ", err)
		}
		c.JSON(http.StatusOK, &gin.H{
			"ok":    err == nil,
			"level": val,
		})
	} else {
		c.JSON(http.StatusOK, &gin.H{
			"ok":    false,
			"level": 0,
		})
	}
}

type SetDeviceVolumeReq struct {
	RendererUrl string `json:"renderer_url" binding:"required"`
	Level       uint8  `json:"level" form:"level"`
}

func SetDeviceVolume(c *gin.Context) {
	var setDeviceVolumeReq SetDeviceVolumeReq
	if c.ShouldBind(&setDeviceVolumeReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	rendererUrl := setDeviceVolumeReq.RendererUrl

	tv, ok := serverMap[rendererUrl]
	if ok {
		l := strconv.Itoa(int(setDeviceVolumeReq.Level))
		err := tv.SetVolumeSoapCall(l)
		if err != nil {
			log.Println("err: ", err)
		}
		c.JSON(http.StatusOK, &gin.H{
			"ok": err == nil,
		})
	} else {
		c.JSON(http.StatusOK, &gin.H{
			"ok": false,
		})
	}
}

type CreatePlayListReq struct {
	Name string `json:"name" binding:"required"`
}

func CreatePlayList(c *gin.Context) {
	var createPlayList CreatePlayListReq
	if c.ShouldBind(&createPlayList) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	pid, err := models.CreatePlayList(createPlayList.Name)
	if err == nil {
		c.JSON(http.StatusOK, &gin.H{
			"pid": pid,
		})
		return
	}
	c.JSON(http.StatusOK, &gin.H{
		"err": err.Error(),
	})
}

type SetPlayListReq struct {
	Pid  string                 `json:"pid" binding:"required"`
	Name string                 `json:"name" form:"name"`
	List []models.ListItemParam `json:"list" binding:"required"`
}

func SetPlayList(c *gin.Context) {
	var setPlayListReq SetPlayListReq
	if c.ShouldBind(&setPlayListReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	ok := models.SetPlayList(setPlayListReq.Pid, setPlayListReq.List)
	if setPlayListReq.Name != "" {
		models.RenamePlayList(setPlayListReq.Pid, setPlayListReq.Name)
	}
	c.JSON(http.StatusOK, &gin.H{
		"ok": ok,
	})
}

type ActionReq struct {
	// start, play, stop, pause
	ActionName string `json:"action_name" binding:"required"`
	// 播放列表id
	Pid string `json:"pid"`
	// 循环模式
	CycleMode   uint8  `json:"cycle_mode"`
	RendererUrl string `json:"renderer_url" binding:"required"`
}

// 执行操作
func Action(c *gin.Context) {
	var actionReq ActionReq
	if c.ShouldBind(&actionReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	if actionReq.ActionName == "start" {
		startPlayPush(actionReq)
	} else if actionReq.ActionName == "stop" {
		stopPlay(actionReq)
	} else if actionReq.ActionName == "pause" {
		pause(actionReq)
	} else if actionReq.ActionName == "play" {
		play(actionReq)
	}
}
func check(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Encountered error(s): %s\n", err)
		os.Exit(1)
	}
}

var serverMap = make(map[string]*(soapcalls.TVPayload), 1)
var songListMap = make(map[string]*[]models.AudioItem, 0)
var create = func(actionReq ActionReq, rendererUrl string) *soapcalls.TVPayload {

	upnpServicesURLs, err := soapcalls.DMRExtractor(actionReq.RendererUrl)
	check(err)
	whereToListen, err := utils.URLtoListenIPandPort(actionReq.RendererUrl)
	check(err)
	callbackPath, err := utils.RandomString()
	check(err)
	tvdata := &soapcalls.TVPayload{
		ControlURL:                  upnpServicesURLs.AVTransportControlURL,
		EventURL:                    upnpServicesURLs.AVTransportEventSubURL,
		RenderingControlURL:         upnpServicesURLs.RenderingControlURL,
		CallbackURL:                 "http://" + whereToListen + "/" + callbackPath,
		MediaURL:                    "http://" + whereToListen + "/" + "media",
		SubtitlesURL:                "http://" + whereToListen + "/",
		MediaType:                   "",
		CurrentTimers:               make(map[string]*time.Timer),
		MediaRenderersStates:        make(map[string]*soapcalls.States),
		InitialMediaRenderersStates: make(map[string]bool),
		RWMutex:                     &sync.RWMutex{},
		Transcode:                   false,
	}
	serverMap[rendererUrl] = tvdata

	s := httphandlers.NewServer(whereToListen)
	serverStarted := make(chan struct{})

	// We pass the tvdata here as we need the callback handlers to be able to react
	// to the different media renderer states.
	go func() {
		err := s.StartServer(serverStarted, tvdata, rendererUrl)
		check(err)
	}()
	// // Wait for HTTP server to properly initialize
	<-serverStarted
	log.Print("after")
	return tvdata
}

func startPlayPush(actionReq ActionReq) {
	rendererUrl := actionReq.RendererUrl

	tv, ok := serverMap[rendererUrl]
	list := models.GetPlayListItems(actionReq.Pid)
	utils.Shuffle(list)
	songListMap[rendererUrl] = list

	if ok { // 已有，直接操作
		tv.PlayListUrls = *list
		tv.SendtoTV("Play1")
	} else {
		tvdata := create(actionReq, rendererUrl)
		tvdata.PlayListUrls = *list
		if err := tvdata.SendtoTV("Play1"); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}

}

func stopPlay(actionReq ActionReq) {
	rendererUrl := actionReq.RendererUrl

	tv, ok := serverMap[rendererUrl]
	if ok {
		tv.SendtoTV("Stop")
	}
}

func play(actionReq ActionReq) {
	rendererUrl := actionReq.RendererUrl

	tv, ok := serverMap[rendererUrl]
	if ok {
		tv.AVTransportActionSoapCall("Play")
	}
}

func pause(actionReq ActionReq) {
	rendererUrl := actionReq.RendererUrl

	tv, ok := serverMap[rendererUrl]
	if ok {
		tv.AVTransportActionSoapCall("Pause")
	}
}
