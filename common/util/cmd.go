package util

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Command execute system cmd
func Command(bin string, argv []string, baseDir string) ([]byte, error) {
	cmd := exec.Command(bin, argv...)
	if baseDir != "" {
		cmd.Dir = baseDir
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.Bytes(), fmt.Errorf("%s: %s", err, stderr.Bytes())
	}
	return stdout.Bytes(), nil
}
