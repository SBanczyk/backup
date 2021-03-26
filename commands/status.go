package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/SBanczyk/backup/model"
)

func Status(currentDir string) error {
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
	displayMissingInStaging(baseDir, *staging)
	displayNewFile(*staging, *backup)
	displayModifiedFile(baseDir, *staging, *backup)
	displayDestroyedFile(*staging, *backup)
	displayMissingFileFromBackup(baseDir, *staging, *backup)
	displayModifiedUnstageFile(baseDir, *staging, *backup)
	displayUntrackedFiles(baseDir, *staging, *backup)
	return nil
}

func displayMissingInStaging(baseDir string, staging model.Staging) {
	for i := range staging.StagingFiles {
		filePath := path.Join(baseDir, staging.StagingFiles[i].Path)
		_, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("File in staging, but not on disk: %v\n", staging.StagingFiles[i].Path)
			} else {
				fmt.Printf("%v\n", err)
			}
		}
	}
}

func displayNewFile(staging model.Staging, backup model.Backup) {
Loop:
	for i := range staging.StagingFiles {
		for k := range backup.Files {
			for j := range backup.Files[k] {
				if staging.StagingFiles[i].Path == backup.Files[k][j].Path {
					continue Loop
				}
			}
		}
		_, err := os.Stat(staging.StagingFiles[i].Path)
		if err == nil {
			fmt.Printf("New file: %v\n", staging.StagingFiles[i].Path)
		}

	}
}

func displayModifiedFile(baseDir string, staging model.Staging, backup model.Backup) {
Loop:
	for i := range staging.StagingFiles {
		for k := range backup.Files {
			for j := range backup.Files[k] {
				if staging.StagingFiles[i].Path == backup.Files[k][j].Path {
					fileHash, err := calculateHash(path.Join(baseDir, staging.StagingFiles[i].Path))
					if err != nil {
						fmt.Printf("%v\n", err)
					} else if fileHash != k {
						fmt.Printf("Modified file: %v\n", staging.StagingFiles[i].Path)
					}
					continue Loop
				}
			}
		}
	}
}

func displayDestroyedFile(staging model.Staging, backup model.Backup) {
Loop:
	for i := range staging.DestroyedFiles {
		for k := range backup.Files {
			for j := range backup.Files[k] {
				if staging.DestroyedFiles[i] == backup.Files[k][j].Path {
					fmt.Printf("Destroyed file: %v\n", staging.DestroyedFiles[i])
					continue Loop
				}
			}
		}
	}
}

func displayMissingFileFromBackup(baseDir string, staging model.Staging, backup model.Backup) {
	for i := range backup.Files {
	Loop:
		for k := range backup.Files[i] {
			if !backup.Files[i][k].Shadow {
				for j := range staging.StagingFiles {
					if staging.StagingFiles[j].Path == backup.Files[i][k].Path {
						continue Loop
					}
				}
				for j := range staging.DestroyedFiles {
					if staging.DestroyedFiles[j] == backup.Files[i][k].Path {
						continue Loop
					}
				}
				path := path.Join(baseDir, backup.Files[i][k].Path)
				_, err := os.Stat(path)
				if err != nil {
					if os.IsNotExist(err) {
						fmt.Printf("File does not exist: %v\n", backup.Files[i][k].Path)
					} else {
						fmt.Printf("%v", err)
					}
				}
			}
		}
	}
}

func displayModifiedUnstageFile(baseDir string, staging model.Staging, backup model.Backup) {
	for i := range backup.Files {
	Loop:
		for k := range backup.Files[i] {
			for j := range staging.StagingFiles {
				if backup.Files[i][k].Path == staging.StagingFiles[j].Path {
					continue Loop
				}
			}
			for j := range staging.DestroyedFiles {
				if staging.DestroyedFiles[j] == backup.Files[i][k].Path {
					continue Loop
				}
			}
			baseDirPath := path.Join(baseDir, backup.Files[i][k].Path)
			baseDirHash, err := calculateHash(baseDirPath)
			if err != nil {
				if !os.IsNotExist(err) {
					fmt.Printf("%v\n", err)
				}
			} else if baseDirHash != i {
				fmt.Printf("Different hash: %v\n", backup.Files[i][k].Path)
			}
		}
	}
}

func displayUntrackedFiles(baseDir string, staging model.Staging, backup model.Backup) {
	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err1 error) error {
		if err1 != nil {
			return err1
		}
		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}
		if relPath == ".backup" {
			return fs.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		for i := range staging.StagingFiles {
			if staging.StagingFiles[i].Path == relPath {
				return nil
			}
		}
		for i := range staging.DestroyedFiles {
			if staging.DestroyedFiles[i] == relPath {
				return nil
			}
		}
		for i := range backup.Files {
			for k := range backup.Files[i] {
				if backup.Files[i][k].Path == relPath {
					return nil
				}
			}
		}
		fmt.Printf("Untracked file: %v\n", relPath)
		return nil
	})
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
