package filedetect

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// FileIsExist 文件是否存在
func FileIsExist(fileName string) (bool, error) {
	fi, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if fi.IsDir() {
		return false, errors.New("this file is dir")
	}

	return true, nil
}

var abolishChars = []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|", "^"}

// CheckFileName 文件名是否合法
func CheckFileName(fileName string) error {
	_, basef := path.Split(fileName)
	for _, char := range abolishChars {
		if strings.Contains(basef, char) {
			return fmt.Errorf("file name has illegal char:%s", char)
		}
	}
	return nil
}
