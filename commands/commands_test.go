package commands_test

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/SBanczyk/backup/commands"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	tempDir := t.TempDir()
	localDir := path.Join(tempDir, "local")
	backupDir := path.Join(tempDir, "backup")
	backupfilesPath := path.Join(backupDir, "backupfiles")
	assert.NoError(t, os.Mkdir(localDir, 0777))
	assert.NoError(t, os.Mkdir(backupDir, 0777))
	assert.NoError(t, ioutil.WriteFile(backupfilesPath, []byte("123456"), 0666))
	err := commands.InitFs(localDir, backupDir)
	assert.NoError(t, err)
	localDirPath := path.Join(localDir, ".backup")
	fileInfo, err := os.Stat(localDirPath)
	assert.True(t, fileInfo.IsDir())
	assert.NoError(t, err)
	configPath := path.Join(localDirPath, "fs_backend.config")
	backupfilesFilePath := path.Join(localDirPath, "backupfiles")
	configPathInfo, err := os.ReadFile(configPath)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(configPathInfo), backupDir))
	backupFilesContent, err := os.ReadFile(backupfilesFilePath)
	assert.NoError(t, err)
	assert.EqualValues(t, "123456", backupFilesContent)
}

func TestInitFsFileExist(t *testing.T) {
	tempDir := t.TempDir()
	localDir := path.Join(tempDir, "local")
	backupDir := path.Join(tempDir, "backup")
	assert.NoError(t, os.Mkdir(localDir, 0777))
	assert.NoError(t, os.Mkdir(backupDir, 0777))
	backupDirPath := path.Join(localDir, ".backup")
	assert.NoError(t, os.Mkdir(backupDirPath, 0777))
	err := commands.InitFs(localDir, backupDir)
	assert.Error(t, err)
}
