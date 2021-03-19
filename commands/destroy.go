package commands

import (
	"fmt"
	"path"

	"github.com/SBanczyk/backup/model"
)

func Destroy (currentDir string, paths []string) error {
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
	for i := range paths{
		for k := range backup.Files {
			for j := range backup.Files[k] {
				if backup.Files[k][j].Path != paths[i] {
					fmt.Printf("some error")
					continue
				}else {
					staging.StagingFiles = removeFromStaging(staging.StagingFiles, paths[i])
					
				}
			}
		}
	}
	return nil
}