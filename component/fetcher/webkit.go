package fetcher

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Webkit Webkit下载器
type Webkit struct {
	option *Option
}

// Fetch 执行请求
func (t *Webkit) Fetch(url string) (*Response, error) {
	if url == "" {
		return nil, errors.New("Webkit.Fetch url is empty")
	}
	if strings.Index(url, "://") == -1 {
		return nil, errors.New("Webkit.Fetch url is not begin with http:// or https://")
	}

	jsRes, err := t.PhantomJS(t.option.method, url, t.option.params)
	if err != nil {
		return nil, err
	}

	res := new(Response)
	res.Code = jsRes.Code
	res.Cookie = jsRes.Cookie
	if len(jsRes.Header) > 0 {
		h := make(http.Header)
		for k, v := range jsRes.Header {
			h.Set(k, v)
		}
		res.Header = h
	}
	switch res.Code {
	case http.StatusOK:
		res.Body = []byte(jsRes.Body)
	default:
		return res, fmt.Errorf("Webkit.Fetch.Do Error: %v", res.Code)
	}

	return res, nil
}
