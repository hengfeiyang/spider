// +build windows

package util

// IsWritable check path is writeable, can return true, can not return false
func IsWritable(path string) bool {
	return true
}
