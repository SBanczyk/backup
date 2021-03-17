package commands

import (
	"io"
	"os"
	"path"

	"github.com/SBanczyk/backup/backend/fs"
)

func InitFs(currentDir string, targetDir string) error {
	configDir, err := initCommon(currentDir)
	if err != nil {
		return err
	}
	err = fs.Init(configDir, targetDir)
	if err != nil {
		return err
	}
	backend, err := fs.Load(configDir)
	if err != nil {
		return err
	}
	downloadedFilePath, err := backend.DownloadBackupFilesFile()
	if err != nil {
		return err
	}
	configDirPath := path.Join(configDir, "backupfiles")
	err = copyFile(downloadedFilePath, configDirPath)
	if err != nil {
		return err
	}

	return nil
}

func copyFile(src string, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
