package commands

import (
	"fmt"
	"github.com/SBanczyk/backup/backend/fs"
	"github.com/SBanczyk/backup/model"
	"os"
	"path"
	"path/filepath"
)

func Get(currentDir string, paths []string) error {
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
	backendPath := path.Join(baseDir, ".backup")
	backend, err := fs.Load(backendPath)
	if err != nil {
		return err
	}
Loop:
	for i := range paths {
		absPath, err := filepath.Abs(paths[i])
		if err != nil {
			return err
		}
		pathRel, err := filepath.Rel(baseDir, absPath)
		if err != nil {
			return err
		}
		for k := range staging.StagingFiles {
			if staging.StagingFiles[k].Path == pathRel {
				fmt.Print("Already in staging")
				continue Loop
			}
		}
		for j := range staging.DestroyedFiles {
			if staging.DestroyedFiles[j] == pathRel {
				fmt.Print("Already in staging")
				continue Loop
			}
		}
		for z := range backup.Files {
			for o := range backup.Files[z] {
				if backup.Files[z][o].Path == pathRel {
					_, err := os.Stat(absPath)
					if err != nil {
						if os.IsNotExist(err) {
							err1 := backend.DownloadFile(z, absPath)
							if err != nil {
								fmt.Printf("%v", err1)
							}
						} else {
							fmt.Printf("%v", err)
						}
					} else {
						hash, err := calculateHash(absPath)
						if err != nil {
							return err
						}
						if hash != z {
							err1 := backend.DownloadFile(z, absPath)
							if err != nil {
								fmt.Printf("%v", err1)
							}
						}
					}
					continue Loop
				}
			}
		}
		fmt.Print("Element not found")

	}
	return nil
}
