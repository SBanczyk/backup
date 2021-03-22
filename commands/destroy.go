package commands

import (
	"fmt"
	"path"

	"github.com/SBanczyk/backup/model"
)

func Destroy(currentDir string, paths []string) error {
	stagingPath := path.Join(currentDir, ".backup", "staging")
	staging, err := model.LoadStaging(stagingPath)
	if err != nil {
		return err
	}
	backupPath := path.Join(currentDir, ".backup", "backupfiles")
	backup, err := model.LoadBackup(backupPath)
	if err != nil {
		return err
	}
	for i := range paths {
		isFound := false
	Loop:
		for k := range backup.Files {
			for j := range backup.Files[k] {
				if backup.Files[k][j].Path == paths[i] {
					isFound = true
					break Loop
				}
			}
		}
		if !isFound {
			fmt.Printf("File not found in backupfiles")
		} else {
			staging.StagingFiles = removeFromStaging(staging.StagingFiles, paths[i])
			staging.DestroyedFiles = addToDestroyed(staging.DestroyedFiles, paths[i])
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
