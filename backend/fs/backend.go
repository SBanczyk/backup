package fs

import (
	"fmt"
	"io"
	"os"
	"path"
)

type backend struct {
	config fsConfig
}

func Init(targetDir string) (*backend, error) {
	sourceFileStat, err := os.Stat(targetDir)
	if err != nil {
		return nil, err
	}
	if !sourceFileStat.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", targetDir)
	}
	return &backend{
		config: fsConfig{
			targetDir: targetDir,
		},
	}, nil
}

func (b *backend) DownloadBackupFilesFile() (string, error) {
	src := path.Join(b.config.targetDir, "backupfiles")
	dst := path.Join(os.TempDir(), "backupfiles")
	err := b.copyFile(src, dst)
	return dst, err
}

func (b *backend) UploadFile(filePath string, cloudName string) error {
	src := filePath
	dst := path.Join(b.config.targetDir, cloudName)
	err := b.copyFile(src, dst)
	return err
}

func (b *backend) DownloadFile(cloudName string, filePath string) error {
	dst := filePath
	src := path.Join(b.config.targetDir, cloudName)
	err := b.copyFile(src, dst)
	return err
}

func (b *backend) RemoveFile(cloudName string) error {
	src := path.Join(b.config.targetDir, cloudName)
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
