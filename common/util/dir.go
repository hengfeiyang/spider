package util

import (
	"os"
	"path/filepath"
)

// GetDir get pwd
func GetDir() string {
	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		return ""
	}
	return filepath.Dir(path)
}
