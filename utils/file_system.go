package utils

import (
	"os"
	"path/filepath"
)

func IsPathExists(pathName string) bool {
	_, err := os.Stat(pathName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func MakeSurePathExists(dirName string) {
	if IsPathExists(dirName) {
		return
	}
	parentDir := filepath.Dir(dirName)
	MakeSurePathExists(parentDir)
	os.Mkdir(dirName, 0777)
}
