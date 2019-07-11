package fetcher

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"github.com/safeie/spider/common/util"
)

// Gokit Go下载器
type Gokit struct {
	option *Option
}

// Fetch 执行请求
func (t *Gokit) Fetch(url string, params, headers map[string]string) (*Response, error) {
	if url == "" {
		return nil, errors.New("Gokit.Fetch url is empty")
	}
	if strings.Index(url, "://") == -1 {
		return nil, errors.New("Gokit.Fetch url is not begin with http:// or https://")
	}

	var resp *http.Response
	var req *http.Request
	var reqBody io.Reader
	var err error
	var ps = make(map[string]string)
	for k, v := range t.option.params {
		ps[k] = v
	}
	for k, v := range params {
		ps[k] = v
	}
	if len(ps) > 0 {
		if t.option.method == "GET" {
			var vals = make(map[string]interface{})
			for k, v := range ps {
				vals[k] = v
			}
			url += util.HTTPBuildQuery(vals, util.QUERY_RFC3986)
		} else {
			var vals = make(neturl.Values)
			for key, val := range ps {
				vals.Add(key, val)
			}
			reqBody = strings.NewReader(vals.Encode())
		}
	}
	req, err = http.NewRequest(t.option.method, url, reqBody)
	req.Header.Set("User-Agent", GetUserAgent(t.option))
	for k, v := range t.option.headers {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if t.option.method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	cookie := t.option.GetCookie()
	if len(cookie) > 0 {
		req.Header.Set("Cookie", cookie)
	}

	client := &http.Client{}
	if t.option.timeout > 0 {
		client.Timeout = time.Second * time.Duration(t.option.timeout)
	}

	if proxyAddr := GetProxyAddr(t.option); proxyAddr != "" {
		if proxy, err := neturl.Parse(proxyAddr); err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxy),
			}
		}
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Gokit.Fetch.Do Error: %s", err)
	}

	res := new(Response)
	res.Code = resp.StatusCode
	res.Header = resp.Header
	switch resp.StatusCode {
	case http.StatusOK:
		res.Body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("Gokit.Fetch.ReadAll Error: %s", err)
		}
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect:
		if location := FixURL(url, resp.Header.Get("location")); location != "" {
			return t.Fetch(location, params, headers)
		}
		return res, fmt.Errorf("Gokit.Fetch.Do Redirect to [%s] %s error", resp.Status, resp.Header.Get("location"))
	default:
		return res, fmt.Errorf("Gokit.Fetch.Do Error: %v", resp.Status)
	}

	// 转编码
	if t.option.charset != "UTF-8" {
		nBody, err := ConvertCharset(res, t.option.charset)
		if err != nil {
			return nil, fmt.Errorf("Fetcher.ConvertCharset Error: %s", err)
		}
		res.Body = nBody
	}

	// cookie
	var cookies []string
	for _, c := range resp.Cookies() {
		cookies = append(cookies, c.Name+"="+c.Value)
	}
	res.Cookie = strings.Join(cookies, "; ")

	return res, nil
}
