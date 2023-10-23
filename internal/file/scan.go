package file

import (
	"errors"
	"github.com/gogf/gf/v2/os/gfile"
)

func ScanFile(path string, suffix string) (string, error) {
	list, err := gfile.ScanDirFile(path, suffix, false)
	if err != nil {
		return "", err
	}
	if len(list) == 0 {
		return "", errors.New("not found file")
	}
	return list[0], nil
}
