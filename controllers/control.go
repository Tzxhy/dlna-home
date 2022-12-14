package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
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

// 开始单文件播放
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

var isAndroid = runtime.GOOS == "android"
var loadSSDPServiceDelay = 2

func init() {
	if isAndroid {
		loadSSDPServiceDelay = 5
	}
}

// 获取设备列表
func GetDeviceList(c *gin.Context) {
	deviceList, err := devices.LoadSSDPServices(loadSSDPServiceDelay)
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

// 获取播放列表
func GetPlayList(c *gin.Context) {
	list, _ := models.GetPlayList()
	var newList []PlayListItem
	var allPlayListItem = models.PlayListItem{
		Pid:        models.ALL_ITEMS_PID,
		Name:       "所有",
		CreateDate: 0,
	}
	newList = append(newList, PlayListItem{
		allPlayListItem,
		*models.GetAllPlayListItems(),
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

// 删除一个播放列表
func DeletePlayList(c *gin.Context) {
	var deletePlayListReq DeletePlayListReq
	if c.ShouldBind(&deletePlayListReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	models.DeletePlayList(deletePlayListReq.Pid)
}

type CreatePlayListReq struct {
	Name string `json:"name" binding:"required"`
}

// 创建一个播放列表
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

// 重命名播放列表
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

type DeleteSingleResourceReq struct {
	Aid string `json:"aid" form:"aid" binding:"required"`
}

// 删除播放列表中单个资源
func DeleteSingleResource(c *gin.Context) {
	var deleteSingleResourceReq DeleteSingleResourceReq
	if c.ShouldBind(&deleteSingleResourceReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	models.DeleteSingleResource(deleteSingleResourceReq.Aid)

	c.JSON(http.StatusOK, &gin.H{
		"ok": true,
	})
}

type SetPlayListReq struct {
	Pid  string                 `json:"pid" binding:"required"`
	Name string                 `json:"name" form:"name"`
	List []models.ListItemParam `json:"list" binding:"required"`
}
type AddListReq struct {
	Pid  string                 `json:"pid" binding:"required"`
	List []models.ListItemParam `json:"list" binding:"required"`
}

// 全量更新播放列表
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

// 增量更新播放列表
func AddPartialListForPlay(c *gin.Context) {
	var addPlayListReq AddListReq
	if c.ShouldBind(&addPlayListReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	ok := models.AddListForPlay(addPlayListReq.Pid, addPlayListReq.List)

	c.JSON(http.StatusOK, &gin.H{
		"ok": ok,
	})
}

type GetDeviceVolumeReq struct {
	RendererUrl string `json:"renderer_url" form:"renderer_url" binding:"required"`
}

// 获取音量
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

// 设置音量
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

type StatusRespItem struct {
	CurrentItem models.AudioItem `json:"current_item"`
	Status      string           `json:"status"`       // 播放器状态
	RendererUrl string           `json:"renderer_url"` // 播放器地址
}
type StatusResp struct {
	Data map[string]StatusRespItem `json:"data"`
}

// 获取状态
func GetStatus(c *gin.Context) {
	var ret = &StatusResp{}
	ret.Data = make(map[string]StatusRespItem)
	for rendererUrl, value := range share.TvDataMap {
		var item models.AudioItem
		if value.CurrentIdx < 0 {
		} else if value.PlayMode == constants.PLAY_MODE_RANDOM {
			item = value.PlayListTempUrls[value.CurrentIdx]

		} else {
			item = value.PlayListUrls[value.CurrentIdx]

		}
		ret.Data[rendererUrl] = StatusRespItem{
			item,
			value.Status,
			value.RenderingControlURL,
		}
	}

	c.JSON(http.StatusOK, &ret)
}

type ActionReq struct {
	ActionName  string `json:"action_name" binding:"required"`  // start, play, stop, pause, next, changePlayMode, jump
	Pid         string `json:"pid"`                             // 播放列表id
	PlayMode    uint8  `json:"play_mode"`                       // 默认乱序
	TargetAid   string `json:"aid"`                             // 要播放的文件的aid
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

// 操作
func jump(actionReq ActionReq) { // 跳到指定歌曲
	rendererUrl := actionReq.RendererUrl

	tv, ok := share.TvDataMap[rendererUrl]
	if ok {
		idx := 0
		if tv.PlayMode == constants.PLAY_MODE_RANDOM {
			idx = utils.FindIndex(&tv.PlayListTempUrls, func(item models.AudioItem) bool {
				return item.Aid == actionReq.TargetAid
			})
		} else {
			idx = utils.FindIndex(&tv.PlayListUrls, func(item models.AudioItem) bool {
				return item.Aid == actionReq.TargetAid
			})
		}

		tv.CurrentIdx = int16(idx) - 1
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
	if ok {
		// 当前是随机的话，切换到其他的，那么设置 CurrentIdx 去完成接续
		if tv.PlayMode == constants.PLAY_MODE_RANDOM {
			if tv.CurrentIdx < 0 {
				return
			}
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

func checkServerIsOnline() bool {
	parser, err := url.Parse("http://" + globalWhereToListen)
	if err != nil {
		log.Println("解析 globalWhereToListen 失败")
		log.Print(err)
		return false
	}

	host := parser.Hostname()
	port := parser.Port()

	log.Println("host: ", host)
	log.Println("port: ", port)

	err = utils.TcpGather(host, port)
	if err != nil {
		log.Println("TcpGather err: ", err)
		// 有错误，说明服务不在线
		return false
	}
	return true
}

func startPlayPush(actionReq ActionReq) {
	rendererUrl := actionReq.RendererUrl

	tv, ok := share.TvDataMap[rendererUrl]
	list := models.GetPlayListItems(actionReq.Pid)

	if ok { // 已有，直接操作
		if !checkServerIsOnline() { // 服务不在线
			log.Println("服务不在线，重启服务")
			hasCreateMediaServer = false
			delete(share.TvDataMap, rendererUrl)
			startPlayPush((actionReq))
			return
		}
		tv.PlayListUrls = *list
		changeMode(actionReq)
		tv.CurrentIdx = -1
		tv.SendtoTV("Play1")
	} else {
		log.Println("has tv")
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

// 相关操作 end

type GetPositionReq struct {
	RendererUrl string `form:"renderer_url" binding:"required"`
}

// 获取设备位置
func GetPosition(c *gin.Context) {
	var getPositionReq GetPositionReq
	err := c.ShouldBindQuery(&getPositionReq)
	if err != nil {
		log.Println("err: ", err)
		utils.ReturnParamNotValid(c)
		return
	}

	tv, ok := share.TvDataMap[getPositionReq.RendererUrl]
	if ok {
		position, err := tv.GetPositionSoapCall()
		if err != nil {
			log.Println("err: ", err)
			utils.ReturnParamNotValid(c)
			return
		}
		c.JSON(http.StatusOK, &gin.H{
			"position": position,
		})

	}

}

type SetPositionReq struct {
	RendererUrl string `form:"renderer_url" json:"renderer_url" binding:"required"`
	RelTime     uint16 `form:"rel_time" json:"rel_time"`
}

// 设置设备的播放位置
func SetPosition(c *gin.Context) {
	var setPositionReq SetPositionReq
	err := c.ShouldBind(&setPositionReq)
	if err != nil {
		log.Println("err: ", err)
		utils.ReturnParamNotValid(c)
		return
	}

	tv, ok := share.TvDataMap[setPositionReq.RendererUrl]
	if ok {
		time := utils.GetRelTimeFromSecond(setPositionReq.RelTime)
		tv.SetPositionSoapCall(time)
	}
}
