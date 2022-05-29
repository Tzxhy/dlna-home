package httphandlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"gitee.com/tzxhy/dlna-home/constants"
	"gitee.com/tzxhy/dlna-home/share"
	"gitee.com/tzxhy/dlna-home/soapcalls"
	"gitee.com/tzxhy/dlna-home/utils"
)

// HTTPserver - new http.Server instance.
type HTTPserver struct {
	http *http.Server
	mux  *http.ServeMux
	// We only need to run one ffmpeg
	// command at a time, per server instance
	ffmpeg *exec.Cmd
}

// Screen interface is used to push message back to the user
// as these are returned by the subscriptions.
type Screen interface {
	EmitMsg(string)
	Fini()
}

// We use this type to be able to test
// the serveContent function without the
// need of os.Open in the tests.
type osFileType struct {
	time time.Time
	file io.ReadSeeker
	path string
}

// StartServer will start a HTTP server to serve the selected media files and
// also handle the subscriptions requests from the DMR devices.
func (s *HTTPserver) StartServer(
	serverStarted chan<- struct{},
	tvpayload *soapcalls.TVPayload,
) error {
	mURL, err := url.Parse(tvpayload.MediaURL)
	if err != nil {
		return fmt.Errorf("failed to parse MediaURL: %w", err)
	}

	// callbackURL, err := url.Parse(tvpayload.CallbackURL)
	// if err != nil {
	// 	return fmt.Errorf("failed to parse CallbackURL: %w", err)
	// }

	s.mux.HandleFunc(mURL.Path, s.serveMediaHandler())
	// s.mux.HandleFunc(callbackURL.Path, s.callbackHandler(tvpayload))

	ln, err := net.Listen("tcp", s.http.Addr)
	if err != nil {
		return fmt.Errorf("server listen error: %w", err)
	}

	serverStarted <- struct{}{}
	_ = s.http.Serve(ln)

	return nil
}

func (s *HTTPserver) serveMediaHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		renderer_url := req.URL.Query()["renderer_url"]
		rendererUrl, decodeErr := url.QueryUnescape(renderer_url[0])
		if decodeErr != nil {
			return
		}
		header := req.Header
		val, has := header["Range"]
		isNewMedia := false
		if has {
			isNewMedia = strings.Contains(val[0], "bytes=0-") // 有的话，就切，没有，就不切
		}

		tv := share.TvDataMap[rendererUrl]

		v := tv.CurrentIdx
		list := tv.PlayListUrls
		if tv.PlayMode == constants.PLAY_MODE_SEQ { // 顺序播放
			if v < 0 { // 无，写0
				tv.CurrentIdx = 0
				v = 0
				utils.WriteLog("first media: " + list[v].Name)
			} else { // 有
				if isNewMedia { // 下一曲？ 加1取模；否则还是当前url
					length := int16(len(list))
					nextIdx := (v + 1) % length
					if v+1 >= length { // 播放完毕，直接停止
						tv.SendtoTV("Stop")
						return
					}
					tv.CurrentIdx = nextIdx
					v = nextIdx
					utils.WriteLog("change media: " + list[v].Name)
				}
			}

		} else if tv.PlayMode == constants.PLAY_MODE_REPEAT_ONE { // 单曲循环

		} else if tv.PlayMode == constants.PLAY_MODE_LIST_REPEAT { // 列表循环
			if v < 0 { // 无，写0
				tv.CurrentIdx = 0
				v = 0
				utils.WriteLog("first media: " + list[v].Name)
			} else { // 有
				if isNewMedia { // 下一曲？ 加1取模；否则还是当前url
					length := int16(len(list))
					nextIdx := (v + 1) % length
					tv.CurrentIdx = nextIdx
					v = nextIdx
					utils.WriteLog("change media: " + list[v].Name)
				}
			}
		} else if tv.PlayMode == constants.PLAY_MODE_RANDOM { // 随机播放
			// if len(tv.PlayListTempUrls) == 0 { // 没有初始化过
			// 	var playListTempUrls = make([]models.AudioItem, len(tv.PlayListUrls))
			// 	copy(playListTempUrls, tv.PlayListUrls)
			// 	utils.Shuffle(&playListTempUrls)
			// 	tv.PlayListTempUrls = playListTempUrls
			// }

			list = tv.PlayListTempUrls
			if v < 0 { // 无，写0
				tv.CurrentIdx = 0
				v = 0
				utils.WriteLog("first media: " + list[v].Name)
			} else { // 有
				if isNewMedia { // 下一曲？ 加1取模；否则还是当前url
					length := int16(len(list))
					nextIdx := (v + 1) % length
					tv.CurrentIdx = nextIdx
					v = nextIdx
					utils.WriteLog("change media: " + list[v].Name)
				}
			}
		}

		mediaURLinfo, _ := utils.StreamURL(context.Background(), list[v].Url)

		serveContent(w, req, tv, mediaURLinfo, s.ffmpeg)
	}
}

// func (s *HTTPserver) callbackHandler(tv *soapcalls.TVPayload) http.HandlerFunc {
// 	return func(w http.ResponseWriter, req *http.Request) {
// 		reqParsed, _ := io.ReadAll(req.Body)
// 		sidVal, sidExists := req.Header["Sid"]
// 		reqParsedUnescape := html.UnescapeString(string(reqParsed))
// 		utils.WriteLog("callback: ")
// 		utils.WriteLog(reqParsedUnescape)
// 		utils.WriteLog(req.Header)

// 		if !sidExists {
// 			http.NotFound(w, req)
// 			return
// 		}

// 		if sidVal[0] == "" {
// 			http.NotFound(w, req)
// 			return
// 		}

// 		uuid := strings.TrimPrefix(sidVal[0], "uuid:")

// 		// Apparently we should ignore the first message
// 		// On some media renderers we receive a STOPPED message
// 		// even before we start streaming.
// 		seq, err := tv.GetSequence(uuid)
// 		if err != nil {
// 			http.NotFound(w, req)
// 			return
// 		}

// 		if seq == 0 {
// 			tv.IncreaseSequence(uuid)
// 			fmt.Fprintf(w, "OK\n")
// 			return
// 		}

// 		previousstate, newstate, err := soapcalls.EventNotifyParser(reqParsedUnescape)
// 		if err != nil {
// 			http.NotFound(w, req)
// 			return
// 		}

// 		if !tv.UpdateMRstate(previousstate, newstate, uuid) {
// 			http.NotFound(w, req)
// 			return
// 		}
// 	}
// }

// StopServer forcefully closes the HTTP server.
func (s *HTTPserver) StopServer() {
	s.http.Close()
}

// NewServer constractor generates a new HTTPserver type.
func NewServer(a string) *HTTPserver {
	mux := http.NewServeMux()
	srv := HTTPserver{
		http:   &http.Server{Addr: a, Handler: mux},
		mux:    mux,
		ffmpeg: new(exec.Cmd),
	}

	return &srv
}

func serveContent(w http.ResponseWriter, r *http.Request, tv *soapcalls.TVPayload, mf interface{}, ff *exec.Cmd) {
	var isMedia bool
	var transcode bool
	var mediaType string

	if tv != nil {
		isMedia = true
		transcode = tv.Transcode
		mediaType = tv.MediaType
	}

	w.Header()["transferMode.dlna.org"] = []string{"Interactive"}

	if isMedia {
		w.Header()["transferMode.dlna.org"] = []string{"Streaming"}
		w.Header()["realTimeInfo.dlna.org"] = []string{"DLNA.ORG_TLAG=*"}
	}

	switch f := mf.(type) {
	case osFileType:
		serveContentCustomType(r, mediaType, transcode, w, f, ff)
	case []byte:
		serveContentBytes(r, mediaType, w, f)
	case io.ReadCloser:
		serveContentReadClose(r, mediaType, transcode, w, f, ff)
	default:
		http.NotFound(w, r)
		return
	}
}

func serveContentBytes(r *http.Request, mediaType string, w http.ResponseWriter, f []byte) {
	if r.Header.Get("getcontentFeatures.dlna.org") == "1" {
		contentFeatures, err := utils.BuildContentFeatures(mediaType, "01", false)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header()["contentFeatures.dlna.org"] = []string{contentFeatures}
	}

	bReader := bytes.NewReader(f)
	name := strings.TrimLeft(r.URL.Path, "/")
	http.ServeContent(w, r, name, time.Now(), bReader)
}

func serveContentReadClose(r *http.Request, mediaType string, transcode bool, w http.ResponseWriter, f io.ReadCloser, ff *exec.Cmd) {
	if r.Header.Get("getcontentFeatures.dlna.org") == "1" {
		contentFeatures, err := utils.BuildContentFeatures(mediaType, "00", transcode)
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
			return
		}

		w.Header()["contentFeatures.dlna.org"] = []string{contentFeatures}
	}

	// Since we're dealing with an io.Reader we can't
	// allow any HEAD requests that some DMRs trigger.
	if transcode && r.Method == http.MethodGet && strings.Contains(mediaType, "video") {
		_ = utils.ServeTranscodedStream(w, f, ff)
		return
	}

	// No seek support
	if r.Method == http.MethodGet {
		_, _ = io.Copy(w, f)
		f.Close()
		return
	}
}

func serveContentCustomType(r *http.Request, mediaType string, transcode bool, w http.ResponseWriter, f osFileType, ff *exec.Cmd) {
	if r.Header.Get("getcontentFeatures.dlna.org") == "1" {

		seek := "01"
		if strings.Contains(mediaType, "video") && transcode {
			seek = "00"
		}

		contentFeatures, err := utils.BuildContentFeatures(mediaType, seek, transcode)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header()["contentFeatures.dlna.org"] = []string{contentFeatures}
	}

	if transcode && r.Method == http.MethodGet && strings.Contains(mediaType, "video") {
		// Since we're dealing with an io.Reader we can't
		// allow any HEAD requests that some DMRs trigger.
		var input interface{} = f.file
		// The only case where we should expect f.path to be ""
		// is only during our unit tests where we emulate the files.
		if f.path != "" {
			input = f.path
		}
		_ = utils.ServeTranscodedStream(w, input, ff)
		return
	}

	name := strings.TrimLeft(r.URL.Path, "/")

	if r.Method == http.MethodGet {
		http.ServeContent(w, r, name, f.time, f.file)
	}
}
