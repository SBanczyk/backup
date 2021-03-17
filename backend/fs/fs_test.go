package fs_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"github.com/SBanczyk/backup/backend"
	"github.com/SBanczyk/backup/backend/fs"
	"github.com/stretchr/testify/assert"
)

func TestFs(t *testing.T) {
	backupDir, localDir, backend := prepareBackend(t)
	backupfilesPath := path.Join(backupDir, "backupfiles")
	assert.NoError(t, ioutil.WriteFile(backupfilesPath, []byte("123456"), 0666))
	testfilePath := path.Join(localDir, "testfile")
	assert.NoError(t, ioutil.WriteFile(testfilePath, []byte("654321"), 0666))
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

	_, localDir, backend := prepareBackend(t)
	testfilePath := path.Join(localDir, "testfile")
	err := backend.UploadFile(testfilePath, "cloudName")
	assert.Error(t, err)
}

func TestDownloadNotExist(t *testing.T) {
	_, localDir, backend := prepareBackend(t)
	testfilePath := path.Join(localDir, "testfile")
	err := backend.DownloadFile("cloudName", testfilePath)
	assert.Error(t, err)
}

func TestRemoveNotExist(t *testing.T) {
	_,_, backend := prepareBackend(t)
	err := backend.RemoveFile("cloudName")
	assert.Error(t, err)
}

func prepareBackend(t *testing.T) (backupDir string, localDir string, b backend.Backend) {
	t.Helper()
	tempDir := t.TempDir()
	localDir = path.Join(tempDir, "local")
	backupDir = path.Join(tempDir, "backup")
	assert.NoError(t, os.Mkdir(localDir, 0777))
	assert.NoError(t, os.Mkdir(backupDir, 0777))
	err := fs.Init(localDir, backupDir)
	assert.NoError(t, err)
	backend, err := fs.Load(localDir)
	assert.NoError(t, err)
	return backupDir, localDir, backend
}
