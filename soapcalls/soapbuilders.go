package soapcalls

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type playEnvelope struct {
	XMLName  xml.Name `xml:"s:Envelope"`
	Schema   string   `xml:"xmlns:s,attr"`
	Encoding string   `xml:"s:encodingStyle,attr"`
	PlayBody playBody `xml:"s:Body"`
}

type playBody struct {
	XMLName    xml.Name   `xml:"s:Body"`
	PlayAction playAction `xml:"u:Play"`
}

type playAction struct {
	XMLName     xml.Name `xml:"u:Play"`
	AVTransport string   `xml:"xmlns:u,attr"`
	InstanceID  string
	Speed       string
}

type pauseEnvelope struct {
	XMLName   xml.Name  `xml:"s:Envelope"`
	Schema    string    `xml:"xmlns:s,attr"`
	Encoding  string    `xml:"s:encodingStyle,attr"`
	PauseBody pauseBody `xml:"s:Body"`
}

type pauseBody struct {
	XMLName     xml.Name    `xml:"s:Body"`
	PauseAction pauseAction `xml:"u:Pause"`
}

type pauseAction struct {
	XMLName     xml.Name `xml:"u:Pause"`
	AVTransport string   `xml:"xmlns:u,attr"`
	InstanceID  string
	Speed       string
}

type stopEnvelope struct {
	XMLName  xml.Name `xml:"s:Envelope"`
	Schema   string   `xml:"xmlns:s,attr"`
	Encoding string   `xml:"s:encodingStyle,attr"`
	StopBody stopBody `xml:"s:Body"`
}

type stopBody struct {
	XMLName    xml.Name   `xml:"s:Body"`
	StopAction stopAction `xml:"u:Stop"`
}

type stopAction struct {
	XMLName     xml.Name `xml:"u:Stop"`
	AVTransport string   `xml:"xmlns:u,attr"`
	InstanceID  string
	Speed       string
}

type setAVTransportEnvelope struct {
	XMLName  xml.Name           `xml:"s:Envelope"`
	Schema   string             `xml:"xmlns:s,attr"`
	Encoding string             `xml:"s:encodingStyle,attr"`
	Body     setAVTransportBody `xml:"s:Body"`
}

type setAVTransportBody struct {
	XMLName           xml.Name          `xml:"s:Body"`
	SetAVTransportURI setAVTransportURI `xml:"u:SetAVTransportURI"`
}

type setAVTransportURI struct {
	XMLName            xml.Name `xml:"u:SetAVTransportURI"`
	AVTransport        string   `xml:"xmlns:u,attr"`
	InstanceID         string
	CurrentURI         string
	CurrentURIMetaData currentURIMetaData `xml:"CurrentURIMetaData"`
}

type currentURIMetaData struct {
	XMLName xml.Name `xml:"CurrentURIMetaData"`
	Value   []byte   `xml:",chardata"`
}

type didLLite struct {
	XMLName      xml.Name     `xml:"DIDL-Lite"`
	SchemaDIDL   string       `xml:"xmlns,attr"`
	DC           string       `xml:"xmlns:dc,attr"`
	Sec          string       `xml:"xmlns:sec,attr"`
	SchemaUPNP   string       `xml:"xmlns:upnp,attr"`
	DIDLLiteItem didLLiteItem `xml:"item"`
}

type didLLiteItem struct {
	SecCaptionInfo   secCaptionInfo   `xml:"sec:CaptionInfo"`
	SecCaptionInfoEx secCaptionInfoEx `xml:"sec:CaptionInfoEx"`
	XMLName          xml.Name         `xml:"item"`
	Restricted       string           `xml:"restricted,attr"`
	UPNPClass        string           `xml:"upnp:class"`
	DCtitle          string           `xml:"dc:title"`
	ID               string           `xml:"id,attr"`
	ParentID         string           `xml:"parentID,attr"`
	ResNode          []resNode        `xml:"res"`
}

type resNode struct {
	XMLName      xml.Name `xml:"res"`
	ProtocolInfo string   `xml:"protocolInfo,attr"`
	Value        string   `xml:",chardata"`
}

type secCaptionInfo struct {
	XMLName xml.Name `xml:"sec:CaptionInfo"`
	Type    string   `xml:"sec:type,attr"`
	Value   string   `xml:",chardata"`
}

type secCaptionInfoEx struct {
	XMLName xml.Name `xml:"sec:CaptionInfoEx"`
	Type    string   `xml:"sec:type,attr"`
	Value   string   `xml:",chardata"`
}

type setMuteEnvelope struct {
	XMLName     xml.Name    `xml:"s:Envelope"`
	Schema      string      `xml:"xmlns:s,attr"`
	Encoding    string      `xml:"s:encodingStyle,attr"`
	SetMuteBody setMuteBody `xml:"s:Body"`
}

type setMuteBody struct {
	XMLName       xml.Name      `xml:"s:Body"`
	SetMuteAction setMuteAction `xml:"u:SetMute"`
}
type GetPositionBodyAction struct {
	InstanceID string
	Schema     string `xml:"xmlns:u,attr"`
}
type GetPositionBody struct {
	GetPositionBodyAction GetPositionBodyAction `xml:"u:GetPositionInfo"`
}
type getPositionEnvelope struct {
	XMLName     xml.Name        `xml:"s:Envelope"`
	Schema      string          `xml:"xmlns:s,attr"`
	Encoding    string          `xml:"s:encodingStyle,attr"`
	GetPosition GetPositionBody `xml:"s:Body"`
}

type setMuteAction struct {
	XMLName          xml.Name `xml:"u:SetMute"`
	RenderingControl string   `xml:"xmlns:u,attr"`
	InstanceID       string
	Channel          string
	DesiredMute      string
}

type getMuteEnvelope struct {
	XMLName     xml.Name    `xml:"s:Envelope"`
	Schema      string      `xml:"xmlns:s,attr"`
	Encoding    string      `xml:"s:encodingStyle,attr"`
	GetMuteBody getMuteBody `xml:"s:Body"`
}

type getMuteBody struct {
	XMLName       xml.Name      `xml:"s:Body"`
	GetMuteAction getMuteAction `xml:"u:GetMute"`
}

type getMuteAction struct {
	XMLName          xml.Name `xml:"u:GetMute"`
	RenderingControl string   `xml:"xmlns:u,attr"`
	InstanceID       string
	Channel          string
}

type getVolumeEnvelope struct {
	XMLName       xml.Name      `xml:"s:Envelope"`
	Schema        string        `xml:"xmlns:s,attr"`
	Encoding      string        `xml:"s:encodingStyle,attr"`
	GetVolumeBody getVolumeBody `xml:"s:Body"`
}

type getVolumeBody struct {
	XMLName         xml.Name        `xml:"s:Body"`
	GetVolumeAction getVolumeAction `xml:"u:GetVolume"`
}

type getVolumeAction struct {
	XMLName          xml.Name `xml:"u:GetVolume"`
	RenderingControl string   `xml:"xmlns:u,attr"`
	InstanceID       string
	Channel          string
}

type setVolumeEnvelope struct {
	XMLName       xml.Name      `xml:"s:Envelope"`
	Schema        string        `xml:"xmlns:s,attr"`
	Encoding      string        `xml:"s:encodingStyle,attr"`
	SetVolumeBody setVolumeBody `xml:"s:Body"`
}

type setVolumeBody struct {
	XMLName         xml.Name        `xml:"s:Body"`
	SetVolumeAction setVolumeAction `xml:"u:SetVolume"`
}

type setVolumeAction struct {
	XMLName          xml.Name `xml:"u:SetVolume"`
	RenderingControl string   `xml:"xmlns:u,attr"`
	InstanceID       string
	Channel          string
	DesiredVolume    string
}

func setAVTransportSoapBuild(mediaURL, mediaType, subtitleURL string) ([]byte, error) {
	mediaTypeSlice := strings.Split(mediaType, "/")

	var class string
	switch mediaTypeSlice[0] {
	case "audio":
		class = "object.item.audioItem.musicTrack"
	case "image":
		class = "object.item.imageItem.photo"
	default:
		class = "object.item.videoItem.movie"
	}

	mediaTitle := mediaURL
	mediaTitlefromURL, err := url.Parse(mediaURL)
	if err == nil {
		mediaTitle = strings.TrimLeft(mediaTitlefromURL.Path, "/")
	}

	re, err := regexp.Compile(`[&<>\\]+`)
	if err != nil {
		return nil, fmt.Errorf("setAVTransportSoapBuild regex compile error: %w", err)
	}
	mediaTitle = re.ReplaceAllString(mediaTitle, "")
	log.Println("mediaURL: ", mediaURL)
	l := didLLite{
		XMLName:    xml.Name{},
		SchemaDIDL: "urn:schemas-upnp-org:metadata-1-0/DIDL-Lite/",
		DC:         "http://purl.org/dc/elements/1.1/",
		Sec:        "http://www.sec.co.kr/",
		SchemaUPNP: "urn:schemas-upnp-org:metadata-1-0/upnp/",
		DIDLLiteItem: didLLiteItem{
			XMLName:    xml.Name{},
			ID:         "0",
			ParentID:   "-1",
			Restricted: "false",
			UPNPClass:  class,
			DCtitle:    mediaTitle,
			ResNode: []resNode{
				{
					XMLName:      xml.Name{},
					ProtocolInfo: fmt.Sprintf("http-get:*:%s:*", mediaType),
					Value:        mediaURL,
				}, {
					XMLName:      xml.Name{},
					ProtocolInfo: "http-get:*:text/srt:*",
					Value:        subtitleURL,
				},
			},
			SecCaptionInfo: secCaptionInfo{
				XMLName: xml.Name{},
				Type:    "srt",
				Value:   subtitleURL,
			},
			SecCaptionInfoEx: secCaptionInfoEx{
				XMLName: xml.Name{},
				Type:    "srt",
				Value:   subtitleURL,
			},
		},
	}
	a, err := xml.Marshal(l)
	if err != nil {
		return nil, fmt.Errorf("setAVTransportSoapBuild #1 Marshal error: %w", err)
	}

	d := setAVTransportEnvelope{
		XMLName:  xml.Name{},
		Schema:   "http://schemas.xmlsoap.org/soap/envelope/",
		Encoding: "http://schemas.xmlsoap.org/soap/encoding/",
		Body: setAVTransportBody{
			XMLName: xml.Name{},
			SetAVTransportURI: setAVTransportURI{
				XMLName:     xml.Name{},
				AVTransport: "urn:schemas-upnp-org:service:AVTransport:1",
				InstanceID:  "0",
				CurrentURI:  mediaURL,
				CurrentURIMetaData: currentURIMetaData{
					XMLName: xml.Name{},
					Value:   a,
				},
			},
		},
	}
	xmlStart := []byte("<?xml version='1.0' encoding='utf-8'?>")
	b, err := xml.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("setAVTransportSoapBuild #2 Marshal error: %w", err)
	}

	// Samsung TV hack.
	b = bytes.ReplaceAll(b, []byte("&#34;"), []byte(`"`))
	b = bytes.ReplaceAll(b, []byte("&amp;"), []byte("&"))

	return append(xmlStart, b...), nil
}

func getPositionSoapBuild() ([]byte, error) {
	d := getPositionEnvelope{
		XMLName:  xml.Name{},
		Schema:   "http://schemas.xmlsoap.org/soap/envelope/",
		Encoding: "http://schemas.xmlsoap.org/soap/encoding/",
		GetPosition: GetPositionBody{
			GetPositionBodyAction{
				"0",
				"urn:schemas-upnp-org:service:AVTransport:1",
			},
		},
	}
	xmlStart := []byte("<?xml version='1.0' encoding='utf-8'?>")
	b, err := xml.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("getPositionSoapBuild #2 Marshal error: %w", err)
	}

	// Samsung TV hack.
	b = bytes.ReplaceAll(b, []byte("&#34;"), []byte(`"`))
	b = bytes.ReplaceAll(b, []byte("&amp;"), []byte("&"))

	return append(xmlStart, b...), nil
}

func playSoapBuild() ([]byte, error) {
	d := playEnvelope{
		XMLName:  xml.Name{},
		Schema:   "http://schemas.xmlsoap.org/soap/envelope/",
		Encoding: "http://schemas.xmlsoap.org/soap/encoding/",
		PlayBody: playBody{
			XMLName: xml.Name{},
			PlayAction: playAction{
				XMLName:     xml.Name{},
				AVTransport: "urn:schemas-upnp-org:service:AVTransport:1",
				InstanceID:  "0",
				Speed:       "1",
			},
		},
	}
	xmlStart := []byte("<?xml version='1.0' encoding='utf-8'?>")
	b, err := xml.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("playSoapBuild Marshal error: %w", err)
	}

	return append(xmlStart, b...), nil
}

func stopSoapBuild() ([]byte, error) {
	d := stopEnvelope{
		XMLName:  xml.Name{},
		Schema:   "http://schemas.xmlsoap.org/soap/envelope/",
		Encoding: "http://schemas.xmlsoap.org/soap/encoding/",
		StopBody: stopBody{
			XMLName: xml.Name{},
			StopAction: stopAction{
				XMLName:     xml.Name{},
				AVTransport: "urn:schemas-upnp-org:service:AVTransport:1",
				InstanceID:  "0",
				Speed:       "1",
			},
		},
	}
	xmlStart := []byte("<?xml version='1.0' encoding='utf-8'?>")
	b, err := xml.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("stopSoapBuild Marshal error: %w", err)
	}

	return append(xmlStart, b...), nil
}

func pauseSoapBuild() ([]byte, error) {
	d := pauseEnvelope{
		XMLName:  xml.Name{},
		Schema:   "http://schemas.xmlsoap.org/soap/envelope/",
		Encoding: "http://schemas.xmlsoap.org/soap/encoding/",
		PauseBody: pauseBody{
			XMLName: xml.Name{},
			PauseAction: pauseAction{
				XMLName:     xml.Name{},
				AVTransport: "urn:schemas-upnp-org:service:AVTransport:1",
				InstanceID:  "0",
				Speed:       "1",
			},
		},
	}
	xmlStart := []byte("<?xml version='1.0' encoding='utf-8'?>")
	b, err := xml.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("pauseSoapBuild Marshal error: %w", err)
	}

	return append(xmlStart, b...), nil
}

func setMuteSoapBuild(m string) ([]byte, error) {
	if m != "0" && m != "1" {
		return nil, errors.New("setMuteSoapBuild input error. Was expecting 0 or 1.")
	}

	d := setMuteEnvelope{
		XMLName:  xml.Name{},
		Schema:   "http://schemas.xmlsoap.org/soap/envelope/",
		Encoding: "http://schemas.xmlsoap.org/soap/encoding/",
		SetMuteBody: setMuteBody{
			XMLName: xml.Name{},
			SetMuteAction: setMuteAction{
				XMLName:          xml.Name{},
				RenderingControl: "urn:schemas-upnp-org:service:RenderingControl:1",
				InstanceID:       "0",
				Channel:          "Master",
				DesiredMute:      m,
			},
		},
	}
	xmlStart := []byte("<?xml version='1.0' encoding='utf-8'?>")
	b, err := xml.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("setMuteSoapBuild Marshal error: %w", err)
	}

	return append(xmlStart, b...), nil
}

func getMuteSoapBuild() ([]byte, error) {
	d := getMuteEnvelope{
		XMLName:  xml.Name{},
		Schema:   "http://schemas.xmlsoap.org/soap/envelope/",
		Encoding: "http://schemas.xmlsoap.org/soap/encoding/",
		GetMuteBody: getMuteBody{
			XMLName: xml.Name{},
			GetMuteAction: getMuteAction{
				XMLName:          xml.Name{},
				RenderingControl: "urn:schemas-upnp-org:service:RenderingControl:1",
				InstanceID:       "0",
				Channel:          "Master",
			},
		},
	}
	xmlStart := []byte("<?xml version='1.0' encoding='utf-8'?>")
	b, err := xml.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("getMuteSoapBuild Marshal error: %w", err)
	}

	return append(xmlStart, b...), nil
}

func getVolumeSoapBuild() ([]byte, error) {
	d := getVolumeEnvelope{
		XMLName:  xml.Name{},
		Schema:   "http://schemas.xmlsoap.org/soap/envelope/",
		Encoding: "http://schemas.xmlsoap.org/soap/encoding/",
		GetVolumeBody: getVolumeBody{
			XMLName: xml.Name{},
			GetVolumeAction: getVolumeAction{
				XMLName:          xml.Name{},
				RenderingControl: "urn:schemas-upnp-org:service:RenderingControl:1",
				InstanceID:       "0",
				Channel:          "Master",
			},
		},
	}
	xmlStart := []byte("<?xml version='1.0' encoding='utf-8'?>")
	b, err := xml.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("getVolumeSoapBuild Marshal error: %w", err)
	}

	return append(xmlStart, b...), nil
}

func setVolumeSoapBuild(v string) ([]byte, error) {
	d := setVolumeEnvelope{
		XMLName:  xml.Name{},
		Schema:   "http://schemas.xmlsoap.org/soap/envelope/",
		Encoding: "http://schemas.xmlsoap.org/soap/encoding/",
		SetVolumeBody: setVolumeBody{
			XMLName: xml.Name{},
			SetVolumeAction: setVolumeAction{
				XMLName:          xml.Name{},
				RenderingControl: "urn:schemas-upnp-org:service:RenderingControl:1",
				InstanceID:       "0",
				Channel:          "Master",
				DesiredVolume:    v,
			},
		},
	}
	xmlStart := []byte("<?xml version='1.0' encoding='utf-8'?>")
	b, err := xml.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("setVolumeSoapBuild Marshal error: %w", err)
	}

	return append(xmlStart, b...), nil
}
