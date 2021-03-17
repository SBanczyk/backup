package commands

import (
	"fmt"
	"os"
	"path"
)

func initCommon(currentDir string) (configDir string, err error) {
	thatDir := path.Join(currentDir, ".backup")
	_, err = os.Stat(thatDir)
	if err != nil {
		if os.IsNotExist(err) {
			err1 := os.Mkdir(thatDir, 0775)
			if err1 != nil {
				return "", err1
			} else {
				return thatDir, nil
			}
		} else {
			return "", err
		}
	} else {
		return "", fmt.Errorf("direcotory exists")
	}

}
