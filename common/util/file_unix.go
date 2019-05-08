// +build darwin linux unix

package util

import (
	"os"

	"golang.org/x/sys/unix"
)

// IsWritable check path is writeable, can return true, can not return false
func IsWritable(path string) bool {
	err := unix.Access(path, unix.O_RDWR)
	if err == nil {
		return true
	}
	// Check if error is "no such file or directory"
	if _, ok := err.(*os.PathError); ok {
		return false
	}
	return false
}
