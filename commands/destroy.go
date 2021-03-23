package commands

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/SBanczyk/backup/model"
)

func Destroy(currentDir string, paths []string) error {
	baseDir, err := checkBaseDir(currentDir)
	if err != nil {
		return err
	}
	stagingPath := path.Join(baseDir, ".backup", "staging")
	staging, err := model.LoadStaging(stagingPath)
	if err != nil {
		return err
	}
	backupPath := path.Join(baseDir, ".backup", "backupfiles")
	backup, err := model.LoadBackup(backupPath)
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
		isFound := false
	Loop:
		for k := range backup.Files {
			for j := range backup.Files[k] {
				if backup.Files[k][j].Path == pathRel {
					isFound = true
					break Loop
				}
			}
		}
		if !isFound {
			fmt.Printf("File not found in backupfiles")
		} else {
			staging.StagingFiles = removeFromStaging(staging.StagingFiles, pathRel)
			staging.DestroyedFiles = addToDestroyed(staging.DestroyedFiles, pathRel)
		}
	}
	err = model.SaveStaging(stagingPath, staging)
	if err != nil {
		return err
	}
	return nil
}

func addToDestroyed(destroyed []string, path string) (newDestroyed []string) {
	for i := range destroyed {
		if destroyed[i] == path {
			return destroyed
		}
	}
	newDestroyed = append(destroyed, path)
	return newDestroyed
}
