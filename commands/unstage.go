package commands

import (
	"path"

	"github.com/SBanczyk/backup/model"
)

func Unstage(currentDir string, paths []string) error {
	stagingPath := path.Join(currentDir, ".backup", "staging")
	staging, err := model.LoadStaging(stagingPath)
	if err != nil {
		return err
	}
	for i := range paths {
		staging.DestroyedFiles = removeFromDestroyed(staging.DestroyedFiles, paths[i])
		staging.StagingFiles = removeFromStaging(staging.StagingFiles, paths[i])
	}
	err = model.SaveStaging(stagingPath, staging)
	if err != nil {
		return err
	}
	return nil
}

func removeFromStaging(slice []model.StagingPath, s string) []model.StagingPath {
	newSlice := make([]model.StagingPath, 0)
	for i := range slice {
		if slice[i].Path != s {
			newSlice = append(newSlice, slice[i])
		}
	}
	return newSlice
}
