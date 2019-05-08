package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// MD5 get hash of string
func MD5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

// MD5Bytes get hash of bytes
func MD5Bytes(s []byte) string {
	m := md5.New()
	m.Write(s)
	return hex.EncodeToString(m.Sum(nil))
}

// MD5File get hash of file content
func MD5File(filename string) string {
	m := md5.New()
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	io.Copy(m, f)
	return hex.EncodeToString(m.Sum(nil))
}
