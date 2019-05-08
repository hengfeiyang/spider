package host

import (
	"crypto/md5"
	"encoding/hex"
	"net"
	"strings"

	"github.com/safeie/spider/common/util"
)

// UUID 根据host信息计算设备唯一识别码
// 计算方法：所有网卡的MAC地址HEX
func UUID() string {
	h := md5.New()

	ir, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, i := range ir {
		h.Write(i.HardwareAddr)
	}
	token := h.Sum(nil)
	dst := make([]byte, 36)
	hex.Encode(dst[:], token[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], token[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], token[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], token[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], token[10:])
	return strings.ToUpper(string(dst))
}

// IPAddress 获取主机的IP地址
func IPAddress() string {
	var addr string
	if localAddrs := util.LocalIPAddrs(); localAddrs != nil {
		addr = localAddrs[0]
	} else {
		addr = "127.0.0.1"
	}
	return addr
}
