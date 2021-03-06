package commands

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/SBanczyk/backup/model"
)

func AddToStaging(currentDir string, paths []string, shadow bool) (err error) {
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
		if !strings.HasPrefix(absPath, baseDir) {
			fmt.Printf("File not in baseDir")
			continue
		}
		pathsInfo, err := os.Stat(paths[i])
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}
		if !pathsInfo.Mode().IsRegular() {
			fmt.Printf("%v is not a file", paths[i])
			continue
		}
		isFound := false

		pathRel, err := filepath.Rel(baseDir, absPath)
		if err != nil {
			return err
		}
		for k := range backup.Files {
			for j := range backup.Files[k] {
				if backup.Files[k][j].Path == pathRel {
					isFound = true
					hash, err := calculateHash(paths[i])
					if err != nil {
						return err
					}
					if hash == k {
						staging.DestroyedFiles = removeFromDestroyed(staging.DestroyedFiles, pathRel)
					} else {
						staging.DestroyedFiles = removeFromDestroyed(staging.DestroyedFiles, pathRel)
						staging.StagingFiles = replaceInStaging(staging.StagingFiles, model.StagingPath{
							Path:   pathRel,
							Shadow: shadow,
						})
					}
				}
			}
		}
		if !isFound {
			staging.DestroyedFiles = removeFromDestroyed(staging.DestroyedFiles, pathRel)
			staging.StagingFiles = replaceInStaging(staging.StagingFiles, model.StagingPath{
				Path:   pathRel,
				Shadow: shadow,
			})
		}
	}
	err = model.SaveStaging(stagingPath, staging)
	if err != nil {
		return err
	}
	return nil
}

func removeFromDestroyed(slice []string, s string) []string {
	newSlice := make([]string, 0)
	for i := range slice {
		if slice[i] != s {
			newSlice = append(newSlice, slice[i])
		}
	}
	return newSlice
}

func calculateHash(path string) (calculatedHash string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func replaceInStaging(slice []model.StagingPath, s model.StagingPath) []model.StagingPath {
	newSlice := make([]model.StagingPath, 0)
	isFound := false
	for i := range slice {
		if slice[i].Path == s.Path {
			newSlice = append(newSlice, s)
			isFound = true
		} else {
			newSlice = append(newSlice, slice[i])
		}
	}
	if !isFound {
		newSlice = append(newSlice, s)
	}
	return newSlice
}
