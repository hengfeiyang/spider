package util

import (
	"io/ioutil"
	"net/http"
	"time"
)

func HttpGet(url string, timeout time.Duration) ([]byte, error) {
	client := &http.Client{
		Timeout: timeout * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
