package commands

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/SBanczyk/backup/backend/fs"
	"github.com/SBanczyk/backup/model"
)

func InitFs(currentDir string, targetDir string) error {
	targetDirAbs, err := filepath.Abs(targetDir)
	if err != nil {
		return err
	}
	configDir, err := initCommon(currentDir)
	if err != nil {
		return err
	}
	err = fs.Init(configDir, targetDirAbs)
	if err != nil {
		return err
	}
	backend, err := fs.Load(configDir)
	if err != nil {
		return err
	}
	configDirPath := path.Join(configDir, "backupfiles")
	downloadedFilePath, err := backend.DownloadBackupFilesFile()
	if err != nil {
		if os.IsNotExist(err) {
			err = model.SaveBackup(configDirPath, &model.Backup{})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		err = copyFile(downloadedFilePath, configDirPath)
		if err != nil {
			return err
		}
	}
	stagingPath := path.Join(configDir, "staging")
	err = model.SaveStaging(stagingPath, &model.Staging{})
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
