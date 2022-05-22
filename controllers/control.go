package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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
	List []models.ListItemParam `json:"list" binding:"required"`
}

func SetPlayList(c *gin.Context) {
	var setPlayListReq SetPlayListReq
	if c.ShouldBind(&setPlayListReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	ok := models.SetPlayList(setPlayListReq.Pid, setPlayListReq.List)
	c.JSON(http.StatusOK, &gin.H{
		"ok": ok,
	})
}

const ()

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
	list := models.GetPlayListItems(actionReq.Pid)
	utils.Shuffle(list)
	songListMap[rendererUrl] = list
	upnpServicesURLs, err := soapcalls.DMRExtractor(actionReq.RendererUrl)
	check(err)
	whereToListen, err := utils.URLtoListenIPandPort(actionReq.RendererUrl)
	check(err)
	callbackPath, err := utils.RandomString()
	check(err)
	url := (*list)[0].Url
	// mediaURL, err := utils.StreamURL(context.Background(), url)
	// check(err)

	mediaURLinfo, err := utils.StreamURL(context.Background(), url)
	check(err)

	mediaType, err := utils.GetMimeDetailsFromStream(mediaURLinfo)
	check(err)
	tvdata := &soapcalls.TVPayload{
		ControlURL:                  upnpServicesURLs.AVTransportControlURL,
		EventURL:                    upnpServicesURLs.AVTransportEventSubURL,
		RenderingControlURL:         upnpServicesURLs.RenderingControlURL,
		CallbackURL:                 "http://" + whereToListen + "/" + callbackPath,
		MediaURL:                    "http://" + whereToListen + "/" + "media",
		SubtitlesURL:                "http://" + whereToListen + "/",
		MediaType:                   mediaType,
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
		err := s.StartServer(serverStarted, mediaURLinfo, "", tvdata, list, rendererUrl)
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

	if ok { // 已有，直接操作
		tv.SendtoTV("Play1")
	} else {
		tvdata := create(actionReq, rendererUrl)

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
