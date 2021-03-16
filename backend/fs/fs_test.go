package fs_test

import (
	"github.com/SBanczyk/backup/backend/fs"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestFs(t *testing.T) {
	backupDir, localDir := prepareDirs(t)
	backupfilesPath := path.Join(backupDir, "backupfiles")
	assert.NoError(t, ioutil.WriteFile(backupfilesPath, []byte("123456"), 0666))
	testfilePath := path.Join(localDir, "testfile")
	assert.NoError(t, ioutil.WriteFile(testfilePath, []byte("654321"), 0666))
	backend, err := fs.Init(backupDir)
	assert.NoError(t, err)
	assert.NoError(t, backend.UploadFile(testfilePath, "cloudName"))
	downloadedfilePath := path.Join(localDir, "downloaded")
	assert.NoError(t, backend.DownloadFile("cloudName", downloadedfilePath))
	text, err := ioutil.ReadFile(downloadedfilePath)
	assert.NoError(t, err)
	assert.EqualValues(t, "654321", text)
	downloadedbackupfiles, err := backend.DownloadBackupFilesFile()
	assert.NoError(t, err)
	text1, err := ioutil.ReadFile(downloadedbackupfiles)
	assert.NoError(t, err)
	assert.EqualValues(t, "123456", text1)
	assert.NoError(t, backend.RemoveFile("cloudName"))
	_, err = os.Stat(path.Join(backupDir, "cloudName"))
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestUploadNotExist(t *testing.T) {

	backupDir, localDir := prepareDirs(t)
	testfilePath := path.Join(localDir, "testfile")
	backend, err := fs.Init(backupDir)
	assert.NoError(t, err)
	err = backend.UploadFile(testfilePath, "cloudName")
	assert.Error(t, err)
}

func TestDownloadNotExist(t *testing.T) {
	backupDir, localDir := prepareDirs(t)
	testfilePath := path.Join(localDir, "testfile")
	backend, err := fs.Init(backupDir)
	assert.NoError(t, err)
	err = backend.DownloadFile("cloudName", testfilePath)
	assert.Error(t, err)
}

func TestRemoveNotExist(t *testing.T) {
	backupDir, _ := prepareDirs(t)
	backend, err := fs.Init(backupDir)
	assert.NoError(t, err)
	err = backend.RemoveFile("cloudName")
	assert.Error(t, err)
}

func prepareDirs(t *testing.T) (localDir string, backupDir string) {
	tempDir := t.TempDir()
	localDir = path.Join(tempDir, "local")
	backupDir = path.Join(tempDir, "backup")
	assert.NoError(t, os.Mkdir(localDir, 0777))
	assert.NoError(t, os.Mkdir(backupDir, 0777))
	return localDir, backupDir
}
