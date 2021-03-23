package commands

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func checkBaseDir(currentDir string) (baseDir string, err error) {
	for {
		currentPath := path.Join(currentDir, ".backup")
		_, err = os.Stat(currentPath)
		if err != nil {
			if os.IsNotExist(err) {
				if currentDir == string(filepath.Separator) {
					return "", fmt.Errorf("BaseDir not found")
				} else {
					currentDir = filepath.Dir(currentDir)
				}
			} else {
				return "", err
			}
		} else {
			return currentDir, nil
		}
	}
}
