package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

// StreamURL returns the response body for the input media URL.
func StreamURL(ctx context.Context, s string) (io.ReadCloser, error) {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return nil, fmt.Errorf("streamURL failed to parse url: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s, nil)

	if err != nil {
		return nil, fmt.Errorf("streamURL failed to call NewRequest: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("streamURL failed to client.Do: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New("streamURL bad status code: " + resp.Status)
	}

	body := resp.Body

	return body, nil
}

const MAX_CACHE_MEDIA = 10

type cacheItem struct {
	Url  string
	Data []byte
}

var MediaCache = make([]cacheItem, 0, MAX_CACHE_MEDIA)

var lock = sync.RWMutex{}

func getCache(url string) ([]byte, error) {
	lock.RLock()
	defer lock.RUnlock()
	idx := FindIndex(&MediaCache, func(item cacheItem) bool {
		return item.Url == url
	})
	if idx >= 0 {
		return MediaCache[idx].Data, nil
	}
	return nil, errors.New("not found")
}
func setCache(url string, data []byte) {
	lock.Lock()
	defer lock.Unlock()
	if len(MediaCache) == MAX_CACHE_MEDIA { // 删除第一个
		// log.Println("删除第一个，再缓存")
		newMediaCache := make([]cacheItem, 0, MAX_CACHE_MEDIA)
		newMediaCache = append(newMediaCache, MediaCache[1:]...)
		newMediaCache = append(newMediaCache, cacheItem{
			url,
			data,
		})
		MediaCache = newMediaCache
	} else {
		// log.Println("添加缓存: ", url)
		MediaCache = append(MediaCache, cacheItem{
			url,
			data,
		})
	}
}

// StreamURL returns the response body for the input media URL.
func StreamURLToBytes(ctx context.Context, s string) ([]byte, error) {

	bytes, err := getCache(s)

	if err == nil {
		// log.Println("返回缓存")
		return bytes, nil
	}

	_, err = url.ParseRequestURI(s)
	if err != nil {
		return nil, fmt.Errorf("streamURL failed to parse url: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s, nil)

	if err != nil {
		return nil, fmt.Errorf("streamURL failed to call NewRequest: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("streamURL failed to client.Do: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New("streamURL bad status code: " + resp.Status)
	}

	body := resp.Body

	newBytes, _ := ioutil.ReadAll(body)

	setCache(s, newBytes)

	return newBytes, nil
}
