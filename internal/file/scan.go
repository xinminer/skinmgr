package file

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"skinmgr/internal/log"
)

func ClearTempFile(path string) error {
	list, err := gfile.ScanDirFile(path, "*.fmv", false)
	if err != nil {
		return err
	}
	for _, f := range list {
		originalName := gstr.Replace(f, ".fmv", "")
		if err = gfile.Move(f, originalName); err != nil {
			log.Log.Warnf("Clear history *.fmv error: %v", err)
		}
	}
	return nil
}

func ScanFile(path string, suffix string) (string, error) {

	list, err := gfile.ScanDirFile(path, suffix, false)
	if err != nil {
		return "", err
	}
	if len(list) == 0 {
		return "", errors.New("not found file")
	}

	selected := list[0]
	transferPath := fmt.Sprintf("%s/%s", path, "transfer")
	if err = gfile.Mkdir(transferPath); err != nil {
		log.Log.Errorf("mkdir transfer path error: %v", err)
	}

	selectedInfo, err := gfile.Stat(selected)
	if err != nil {
		log.Log.Errorf("selected info error: %v", err)
	}

	transfer := fmt.Sprintf("%s/%s", transferPath, selectedInfo.Name())

	if err = gfile.Move(selected, transfer); err != nil {
		return "", err
	}

	return transfer, nil
}
