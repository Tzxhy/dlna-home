package devices

import (
	"fmt"
	"sort"

	"gitee.com/tzxhy/dlna-home/soapcalls"
	"github.com/koron/go-ssdp"
	"github.com/pkg/errors"
)

// LoadSSDPServices returns a map with all the devices that support the
// AVTransport service.
func LoadSSDPServices(delay int) (map[string]string, error) {
	// Reset device list every time we call this.
	deviceList := make(map[string]string)
	list, err := ssdp.Search(ssdp.All, delay, "")
	if err != nil {
		return nil, fmt.Errorf("LoadSSDPservices search error: %w", err)
	}

	for _, srv := range list {
		// We only care about the AVTransport services for basic actions
		// (stop,play,pause). If we need support other functionalities
		// like volume control we need to use the RenderingControl service.
		if srv.Type == "urn:schemas-upnp-org:service:AVTransport:1" {
			friendlyName, err := soapcalls.GetFriendlyName(srv.Location)
			if err != nil {
				continue
			}

			deviceList[friendlyName] = srv.Location
		}
	}

	if len(deviceList) > 0 {
		return deviceList, nil
	}

	return nil, errors.New("loadSSDPservices: No available Media Renderers")
}

// DevicePicker will pick the nth device from the devices input map.
func DevicePicker(devices map[string]string, n int) (string, error) {
	if n > len(devices) || len(devices) == 0 || n <= 0 {
		return "", errors.New("devicePicker: Requested device not available")
	}

	keys := make([]string, 0)
	for k := range devices {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for q, k := range keys {
		if n == q+1 {
			return devices[k], nil
		}
	}
	return "", errors.New("devicePicker: Something went terribly wrong")
}
