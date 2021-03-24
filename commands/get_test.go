package commands_test

import (
	"os"
	"path"
	"testing"

	"github.com/SBanczyk/backup/commands"
	"github.com/SBanczyk/backup/model"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	localDir, backupDir, stagingPath, backupPath := initFsBackend(t)
	err := model.SaveStaging(stagingPath, &model.Staging{})
	assert.NoError(t, err)
	localFilePath := createLocalFile(t, localDir, "localFile")
	hash, err := calculateHash(localFilePath)
	assert.NoError(t, err)
	err = os.Rename(localFilePath, path.Join(backupDir, hash))
	assert.NoError(t, err)
	err = model.SaveBackup(backupPath, &model.Backup{
		Files: map[string][]model.BackupPath{
			hash: {
				{
					Path:   "secondFile",
					Shadow: true,
				},
			},
		},
	})
	assert.NoError(t, err)
	newPath := path.Join(localDir, "secondFile")
	err = commands.Get(localDir, []string{newPath})
	assert.NoError(t, err)
	newPathInfo, err := os.ReadFile(newPath)
	assert.NoError(t, err)
	assert.EqualValues(t, localFilePath, newPathInfo)
}

func TestGetInStaging(t *testing.T) {
	localDir, backupDir, stagingPath, backupPath := initFsBackend(t)
	localFilePath := createLocalFile(t, localDir, "localFile")
	hash, err := calculateHash(localFilePath)
	assert.NoError(t, err)
	err = os.Rename(localFilePath, path.Join(backupDir, hash))
	assert.NoError(t, err)
	err = model.SaveStaging(stagingPath, &model.Staging{
		StagingFiles: []model.StagingPath{
			{
				Path:   "secondFile",
				Shadow: true,
			},
		},
	})
	assert.NoError(t, err)
	err = model.SaveBackup(backupPath, &model.Backup{
		Files: map[string][]model.BackupPath{
			hash: {
				{
					Path:   "secondFile",
					Shadow: true,
				},
			},
		},
	})
	assert.NoError(t, err)
	newPath := path.Join(localDir, "secondFile")
	err = commands.Get(localDir, []string{newPath})
	assert.NoError(t, err)
	_, err = os.Stat(newPath)
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestGetDifferentHash(t *testing.T) {
	localDir, backupDir, stagingPath, backupPath := initFsBackend(t)
	err := model.SaveStaging(stagingPath, &model.Staging{})
	assert.NoError(t, err)
	localFilePath := createLocalFile(t, localDir, "localFile")
	hash, err := calculateHash(localFilePath)
	assert.NoError(t, err)
	err = os.Rename(localFilePath, path.Join(backupDir, hash))
	assert.NoError(t, err)
	err = model.SaveBackup(backupPath, &model.Backup{
		Files: map[string][]model.BackupPath{
			hash: {
				{
					Path:   "secondFile",
					Shadow: true,
				},
			},
		},
	})
	assert.NoError(t, err)
	newPath := path.Join(localDir, "secondFile")
	err = os.WriteFile(newPath, []byte("123456"), 0777)
	assert.NoError(t, err)
	err = commands.Get(localDir, []string{newPath})
	assert.NoError(t, err)
	newPathInfo, err := os.ReadFile(newPath)
	assert.NoError(t, err)
	assert.EqualValues(t, localFilePath, newPathInfo)
}

func initFsBackend(t *testing.T) (localdir string, backupdir string, stagingpath string, backuppath string) {
	t.Helper()
	tempDir := t.TempDir()
	localDir := path.Join(tempDir, "local")
	backupDir := path.Join(tempDir, "backup")
	configDir := path.Join(localDir, ".backup")
	stagingPath := path.Join(configDir, "staging")
	backupPath := path.Join(configDir, "backupfiles")
	assert.NoError(t, os.Mkdir(localDir, 0777))
	assert.NoError(t, os.Mkdir(backupDir, 0777))
	err := commands.InitFs(localDir, backupDir)
	assert.NoError(t, err)
	return localDir, backupDir, stagingPath, backupPath
}
