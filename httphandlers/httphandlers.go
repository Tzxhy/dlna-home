package httphandlers

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

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
	rendererUrl string,
) error {
	mURL, err := url.Parse(tvpayload.MediaURL)
	if err != nil {
		return fmt.Errorf("failed to parse MediaURL: %w", err)
	}

	callbackURL, err := url.Parse(tvpayload.CallbackURL)
	if err != nil {
		return fmt.Errorf("failed to parse CallbackURL: %w", err)
	}

	s.mux.HandleFunc(mURL.Path, s.serveMediaHandler(tvpayload, rendererUrl))
	s.mux.HandleFunc(callbackURL.Path, s.callbackHandler(tvpayload))

	ln, err := net.Listen("tcp", s.http.Addr)
	if err != nil {
		return fmt.Errorf("server listen error: %w", err)
	}

	serverStarted <- struct{}{}
	_ = s.http.Serve(ln)

	return nil
}

var songIdxMap = make(map[string]uint8)

func (s *HTTPserver) serveMediaHandler(
	tv *soapcalls.TVPayload,
	rendererUrl string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		reqParsed, _ := io.ReadAll(req.Body)
		reqParsedUnescape := html.UnescapeString(string(reqParsed))
		utils.WriteLog("media:")
		utils.WriteLog(reqParsedUnescape)
		utils.WriteLog(req.Header)
		header := req.Header
		val, has := header["Range"]
		isNewMedia := false
		if has {
			log.Println("range value: ", val[0])
			isNewMedia = strings.Contains(val[0], "bytes=0-") // 有的话，就切，没有，就不切
		}

		v, has := songIdxMap[rendererUrl]
		list := tv.PlayListUrls
		if !has { // 无，写0
			songIdxMap[rendererUrl] = 0
			v = 0
			utils.WriteLog("首次播放：")
			utils.WriteLog(list[v].Url)
		} else { // 有
			if isNewMedia { // 下一曲？ 加1取模；否则还是当前url
				utils.WriteLog("切换url：")
				length := uint8(len(list))
				nextIdx := (v + 1) % length
				songIdxMap[rendererUrl] = nextIdx
				v = nextIdx
				utils.WriteLog(list[v].Url)
			}
		}

		mediaURLinfo, _ := utils.StreamURL(context.Background(), list[v].Url)

		serveContent(w, req, tv, mediaURLinfo, s.ffmpeg)
	}
}

func (s *HTTPserver) callbackHandler(tv *soapcalls.TVPayload) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		reqParsed, _ := io.ReadAll(req.Body)
		sidVal, sidExists := req.Header["Sid"]
		reqParsedUnescape := html.UnescapeString(string(reqParsed))
		utils.WriteLog("callback: ")
		utils.WriteLog(reqParsedUnescape)
		utils.WriteLog(req.Header)

		if !sidExists {
			http.NotFound(w, req)
			return
		}

		if sidVal[0] == "" {
			http.NotFound(w, req)
			return
		}

		uuid := strings.TrimPrefix(sidVal[0], "uuid:")

		// Apparently we should ignore the first message
		// On some media renderers we receive a STOPPED message
		// even before we start streaming.
		seq, err := tv.GetSequence(uuid)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		if seq == 0 {
			tv.IncreaseSequence(uuid)
			fmt.Fprintf(w, "OK\n")
			return
		}

		previousstate, newstate, err := soapcalls.EventNotifyParser(reqParsedUnescape)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		if !tv.UpdateMRstate(previousstate, newstate, uuid) {
			http.NotFound(w, req)
			return
		}
	}
}

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
