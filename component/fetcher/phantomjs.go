package fetcher

import (
	"encoding/json"
	"fmt"
	neturl "net/url"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/safeie/spider/common/util"
)

// PhantomJSResponse Phantomjs执行返回结果
type PhantomJSResponse struct {
	Header map[string]string
	Code   int
	Cookie string
	Body   string
}

var phantomJSBin string
var phantomJSFiles [2]string

func init() {
	phantomJSBin = "/bin/" + runtime.GOOS + "/phantomjs" // BIN
	phantomJSFiles[0] = "/js/phantomjs_get.js"           // GET
	phantomJSFiles[1] = "/js/phantomjs_post.js"          // POST
}

// PhantomJS 执行phatomjs请求渲染
/*
 *  --proxy=<val>                        Sets the proxy server, e.g. '--proxy=http://proxy.company.com:8080'
 *  --proxy-auth=<val>                   Provides authentication information for the proxy, e.g. ''-proxy-auth=username:password'
 *  --proxy-type=<val>                   Specifies the proxy type, 'http' (default), 'none' (disable completely), or 'socks5'
 */
func (t *Webkit) PhantomJS(method string, url string, params map[string]string) (*PhantomJSResponse, error) {
	var paramsBody string
	if len(params) > 0 {
		if method == "GET" {
			var vals map[string]interface{}
			for k, v := range params {
				vals[k] = v
			}
			url += util.HTTPBuildQuery(vals, util.QUERY_RFC3986)
		} else {
			var vals *neturl.Values
			for key, val := range params {
				vals.Add(key, val)
			}
			paramsBody = vals.Encode()
		}
	}
	var args []string
	if proxyAddr := GetProxyAddr(t.option); proxyAddr != "" {
		args = append(args, "--proxy=http://"+proxyAddr)
	}

	switch method {
	case "GET":
		args = append(args, t.option.configDir+phantomJSFiles[0],
			url,
			t.option.charset,
			GetUserAgent(t.option),
			t.option.GetCookie(),
			strconv.Itoa(t.option.renderDelay),
			strconv.Itoa(t.option.timeout),
		)
	case "POST":
		args = append(args, t.option.configDir+phantomJSFiles[1],
			url,
			t.option.charset,
			GetUserAgent(t.option),
			t.option.GetCookie(),
			strconv.Itoa(t.option.renderDelay),
			strconv.Itoa(t.option.timeout),
			paramsBody,
		)
	}

	var err error
	var out []byte
	var res *PhantomJSResponse
	// 由于JS的不稳定性，增加重试次数
	for i := 0; i < 3; i++ {
		out, err = util.Command(t.option.configDir+phantomJSBin, args, os.TempDir())
		if err != nil {
			time.Sleep(300 * time.Millisecond)
			continue
		}
		err = json.Unmarshal(out, &res)
		if err != nil {
			time.Sleep(300 * time.Millisecond)
			continue
		}
		break
	}
	if err != nil {
		err = fmt.Errorf("Webkit.PhantomJS Error: %v", err)
	}
	return res, err
}
