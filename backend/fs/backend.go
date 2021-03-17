package fs

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type backend struct {
	config fsConfig
}

func Init(localDir string, targetDir string,) (error) {
	sourceFileStat, err := os.Stat(targetDir)
	if err != nil {
		return  err
	}
	if !sourceFileStat.IsDir() {
		return fmt.Errorf("%s is not a directory", targetDir)
	}
	b, err := json.Marshal(fsConfig{
		TargetDir: targetDir,
	})
	if err != nil {
		return err
	}
	localConfig := path.Join(localDir, "fs_backend.config")
	err = ioutil.WriteFile(localConfig, b, 0644)
	if err != nil {
		return err
	}
	return nil
	
}

func Load (localDir string) (*backend, error) {
	localConfig := path.Join(localDir, "fs_backend.config")
	file,err := ioutil.ReadFile(localConfig)
	if err != nil {
		return nil, err
	}
	var m fsConfig
	err = json.Unmarshal(file, &m)
	if err != nil {
		return nil, err
	}
	return &backend{
		config: m,
	}, nil
}

func (b *backend) DownloadBackupFilesFile() (string, error) {
	src := path.Join(b.config.TargetDir, "backupfiles")
	dst := path.Join(os.TempDir(), "backupfiles")
	err := b.copyFile(src, dst)
	return dst, err
}

func (b *backend) UploadFile(filePath string, cloudName string) error {
	src := filePath
	dst := path.Join(b.config.TargetDir, cloudName)
	err := b.copyFile(src, dst)
	return err
}

func (b *backend) DownloadFile(cloudName string, filePath string) error {
	dst := filePath
	src := path.Join(b.config.TargetDir, cloudName)
	err := b.copyFile(src, dst)
	return err
}

func (b *backend) RemoveFile(cloudName string) error {
	src := path.Join(b.config.TargetDir, cloudName)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	err = os.Remove(src)
	return err
}

func (b *backend) copyFile(src string, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}
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
