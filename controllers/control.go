package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"gitee.com/tzxhy/dlna-home/constants"
	"gitee.com/tzxhy/dlna-home/devices"
	"gitee.com/tzxhy/dlna-home/httphandlers"
	"gitee.com/tzxhy/dlna-home/models"
	"gitee.com/tzxhy/dlna-home/share"
	"gitee.com/tzxhy/dlna-home/soapcalls"
	"gitee.com/tzxhy/dlna-home/utils"
	"github.com/gin-gonic/gin"
)

type StartOneReq struct {
	Url         string `json:"url" form:"url" binding:"required"`
	RendererUrl string `json:"renderer_url" form:"renderer_url" binding:"required"`
}

func StartOne(c *gin.Context) {
	var actionReq StartOneReq
	if c.ShouldBind(&actionReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	upnpServicesURLs, err := soapcalls.DMRExtractor(actionReq.RendererUrl)
	check(err)

	tvdata := &soapcalls.TVPayload{
		ControlURL:                  upnpServicesURLs.AVTransportControlURL,
		EventURL:                    upnpServicesURLs.AVTransportEventSubURL,
		RenderingControlURL:         upnpServicesURLs.RenderingControlURL,
		CallbackURL:                 "",
		MediaURL:                    actionReq.Url,
		SubtitlesURL:                "",
		MediaType:                   "",
		CurrentTimers:               make(map[string]*time.Timer),
		MediaRenderersStates:        make(map[string]*soapcalls.States),
		InitialMediaRenderersStates: make(map[string]bool),
		RWMutex:                     &sync.RWMutex{},
		Transcode:                   false,
		CurrentIdx:                  -1,
	}
	err = tvdata.SendtoTV("Play1")
	if err != nil {
		log.Println("err: ", err)
	}
	c.JSON(http.StatusOK, &gin.H{
		"ok": err == nil,
	})
}

func GetDeviceList(c *gin.Context) {
	deviceList, err := devices.LoadSSDPServices(2)
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
	var allPlayListItem = models.PlayListItem{
		Pid:        "ALL",
		Name:       "所有",
		CreateDate: 0,
	}
	newList = append(newList, PlayListItem{
		allPlayListItem,
		make([]models.AudioItem, 0),
	})
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

	tv, ok := share.TvDataMap[rendererUrl]
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

	tv, ok := share.TvDataMap[rendererUrl]
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

type RenamePlayListReq struct {
	Pid  string `json:"pid" binding:"required"`
	Name string `json:"new_name" form:"new_name" binding:"required"`
}

func RenamePlayList(c *gin.Context) {
	var renamePlayListReq RenamePlayListReq
	if c.ShouldBind(&renamePlayListReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	models.RenamePlayList(renamePlayListReq.Pid, renamePlayListReq.Name)
	c.JSON(http.StatusOK, &gin.H{
		"ok": true,
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

type StatusRespItem struct {
	Status      string `json:"status"`
	RendererUrl string `json:"renderer_url"`
	CurrentUrl  string `json:"current_url"`
	Name        string `json:"name"`
}
type StatusResp struct {
	Data map[string]StatusRespItem `json:"data"`
}

func GetStatus(c *gin.Context) {
	var ret = &StatusResp{}
	ret.Data = make(map[string]StatusRespItem)
	for rendererUrl, value := range share.TvDataMap {
		var url = ""
		var name = ""
		if value.CurrentIdx < 0 {
		} else if value.PlayMode == constants.PLAY_MODE_RANDOM {
			url = value.PlayListTempUrls[value.CurrentIdx].Url
			name = value.PlayListTempUrls[value.CurrentIdx].Name
		} else {
			url = value.PlayListUrls[value.CurrentIdx].Url
			name = value.PlayListUrls[value.CurrentIdx].Name
		}

		ret.Data[rendererUrl] = StatusRespItem{
			value.Status,
			value.RenderingControlURL,
			url,
			name,
		}
	}

	c.JSON(http.StatusOK, &ret)
}

type ActionReq struct {
	ActionName  string `json:"action_name" binding:"required"`  // start, play, stop, pause, next, changePlayMode, jump
	Pid         string `json:"pid"`                             // 播放列表id
	PlayMode    uint8  `json:"play_mode"`                       // 默认乱序
	TargetIdx   int16  `json:"target_idx"`                      // 要播放的文件的索引
	RendererUrl string `json:"renderer_url" binding:"required"` // 播放器地址
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
	} else if actionReq.ActionName == "next" {
		next(actionReq)
	} else if actionReq.ActionName == "previous" {
		previous(actionReq)
	} else if actionReq.ActionName == "changePlayMode" {
		changeMode(actionReq)
	} else if actionReq.ActionName == "jump" {
		jump(actionReq)
	}
}
func jump(actionReq ActionReq) { // 跳到指定歌曲
	rendererUrl := actionReq.RendererUrl

	tv, ok := share.TvDataMap[rendererUrl]
	if ok {
		tv.CurrentIdx = actionReq.TargetIdx - 1
		if tv.CurrentIdx < -1 {
			tv.CurrentIdx = -1
		}
		tv.SendtoTV("Play1")
	}
}

func next(actionReq ActionReq) {
	rendererUrl := actionReq.RendererUrl

	tv, ok := share.TvDataMap[rendererUrl]
	if ok {
		tv.SendtoTV("Play1")
	}
}

func previous(actionReq ActionReq) {
	rendererUrl := actionReq.RendererUrl

	tv, ok := share.TvDataMap[rendererUrl]
	if ok {
		tv.CurrentIdx -= 2
		if tv.CurrentIdx < -1 {
			tv.CurrentIdx = -1
		}
		tv.SendtoTV("Play1")
	}
}
func changeMode(actionReq ActionReq) {
	rendererUrl := actionReq.RendererUrl

	tv, ok := share.TvDataMap[rendererUrl]
	if ok && tv.PlayMode != actionReq.PlayMode {
		// 当前是随机的话，切换到其他的，那么设置CurrentIdx 去完成接续
		if tv.PlayMode == constants.PLAY_MODE_RANDOM {
			currentItem := tv.PlayListTempUrls[tv.CurrentIdx]
			newIdx := utils.FindIndex(&tv.PlayListUrls, func(item models.AudioItem) bool {
				return item.Aid == currentItem.Aid
			})
			tv.CurrentIdx = int16(newIdx)
		}
		log.Println("change play mode: ", actionReq.PlayMode)
		tv.PlayMode = actionReq.PlayMode
		if actionReq.PlayMode == constants.PLAY_MODE_RANDOM {
			var playListTempUrls = make([]models.AudioItem, len(tv.PlayListUrls))
			copy(playListTempUrls, tv.PlayListUrls)
			utils.Shuffle(&playListTempUrls)
			tv.PlayListTempUrls = playListTempUrls
		}
	}
}
func check(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Encountered error(s): %s\n", err)
		os.Exit(1)
	}
}

var hasCreateMediaServer = false
var globalWhereToListen = ""
var create = func(actionReq ActionReq, rendererUrl string) *soapcalls.TVPayload {

	upnpServicesURLs, err := soapcalls.DMRExtractor(actionReq.RendererUrl)
	check(err)
	if globalWhereToListen == "" {
		globalWhereToListen, _ = utils.URLtoListenIPandPort(actionReq.RendererUrl)
		if !hasCreateMediaServer {
			check(err)
		}
	}

	callbackPath, err := utils.RandomString()
	check(err)
	tvdata := &soapcalls.TVPayload{
		ControlURL:                  upnpServicesURLs.AVTransportControlURL,
		EventURL:                    upnpServicesURLs.AVTransportEventSubURL,
		RenderingControlURL:         upnpServicesURLs.RenderingControlURL,
		CallbackURL:                 "http://" + globalWhereToListen + "/" + callbackPath,
		MediaURL:                    "http://" + globalWhereToListen + "/" + "media?renderer_url=" + url.QueryEscape(actionReq.RendererUrl),
		SubtitlesURL:                "http://" + globalWhereToListen + "/",
		MediaType:                   "",
		CurrentTimers:               make(map[string]*time.Timer),
		MediaRenderersStates:        make(map[string]*soapcalls.States),
		InitialMediaRenderersStates: make(map[string]bool),
		RWMutex:                     &sync.RWMutex{},
		Transcode:                   false,
		CurrentIdx:                  -1,
	}
	share.TvDataMap[rendererUrl] = tvdata

	if !hasCreateMediaServer {
		hasCreateMediaServer = true
		s := httphandlers.NewServer(globalWhereToListen)
		serverStarted := make(chan struct{})

		go func() {
			err := s.StartServer(serverStarted, tvdata)
			check(err)
		}()
		<-serverStarted
	}

	return tvdata
}

func startPlayPush(actionReq ActionReq) {
	rendererUrl := actionReq.RendererUrl

	tv, ok := share.TvDataMap[rendererUrl]
	list := models.GetPlayListItems(actionReq.Pid)

	if ok { // 已有，直接操作
		tv.PlayListUrls = *list
		tv.SendtoTV("Play1")
	} else {
		tvdata := create(actionReq, rendererUrl)
		tvdata.PlayListUrls = *list
		changeMode(actionReq)
		if err := tvdata.SendtoTV("Play1"); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		tvdata.Status = "play"
	}

}

func stopPlay(actionReq ActionReq) {
	rendererUrl := actionReq.RendererUrl

	tv, ok := share.TvDataMap[rendererUrl]
	if ok {
		tv.SendtoTV("Stop")
		tv.Status = "stop"
	}
}

func play(actionReq ActionReq) {
	rendererUrl := actionReq.RendererUrl

	tv, ok := share.TvDataMap[rendererUrl]
	if ok {
		tv.AVTransportActionSoapCall("Play")
	}
}

func pause(actionReq ActionReq) {
	rendererUrl := actionReq.RendererUrl

	tv, ok := share.TvDataMap[rendererUrl]
	if ok {
		tv.AVTransportActionSoapCall("Pause")
	}
}
