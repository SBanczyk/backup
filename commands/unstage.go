package commands

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/SBanczyk/backup/model"
)

func Unstage(currentDir string, paths []string) error {
	baseDir, err := checkBaseDir(currentDir)
	if err != nil {
		return err
	}
	stagingPath := path.Join(baseDir, ".backup", "staging")
	staging, err := model.LoadStaging(stagingPath)
	if err != nil {
		return err
	}
	for i := range paths {
		absPath, err := filepath.Abs(paths[i])
		if err != nil {
			return err
		}
		pathRel, err := filepath.Rel(baseDir, absPath)
		if err != nil {
			return err
		}
		fmt.Printf("pathRel: %v\n", pathRel)
		staging.DestroyedFiles = removeFromDestroyed(staging.DestroyedFiles, pathRel)
		staging.StagingFiles = removeFromStaging(staging.StagingFiles, pathRel)
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
