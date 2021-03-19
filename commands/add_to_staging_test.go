package commands_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/SBanczyk/backup/commands"
	"github.com/SBanczyk/backup/model"
	"github.com/stretchr/testify/assert"
)

func TestBackup(t *testing.T) {
	localDir, _, configDir := createDirs(t)
	someFilePath := createLocalFile(t, localDir, "someFile")
	backupPath := path.Join(configDir, "backupfiles")
	err := model.SaveBackup(backupPath, &model.Backup{})
	assert.NoError(t, err)
	firstFilePath := createLocalFile(t, localDir, "firstFile")
	secondFilePath := createLocalFile(t, localDir, "secondFile")
	stagingPath := path.Join(configDir, "staging")
	err = model.SaveStaging(stagingPath, &model.Staging{
		StagingFiles: []model.StagingPath{{
			Path:   firstFilePath,
			Shadow: true,
		}, {
			Path:   "1qaz",
			Shadow: false,
		}},
		DestroyedFiles: []string{secondFilePath, "654321"},
	})
	assert.NoError(t, err)
	err = commands.AddToStaging(localDir, []string{someFilePath, firstFilePath}, false)
	assert.NoError(t, err)

	staging, err := model.LoadStaging(stagingPath)
	assert.NoError(t, err)
	assert.Len(t, staging.StagingFiles, 3)
	assert.EqualValues(t, &model.Staging{
		StagingFiles: []model.StagingPath{{
			Path:   firstFilePath,
			Shadow: false,
		}, {
			Path:   "1qaz",
			Shadow: false,
		}, {
			Path:   someFilePath,
			Shadow: false,
		},
		}, DestroyedFiles: []string{secondFilePath, "654321"}}, staging)
	err = commands.AddToStaging(localDir, []string{someFilePath, secondFilePath}, true)
	assert.NoError(t, err)
	staging, err = model.LoadStaging(stagingPath)
	assert.NoError(t, err)
	assert.Len(t, staging.StagingFiles, 4)
	assert.EqualValues(t, &model.Staging{
		StagingFiles: []model.StagingPath{{
			Path:   firstFilePath,
			Shadow: false,
		}, {
			Path:   "1qaz",
			Shadow: false,
		}, {
			Path:   someFilePath,
			Shadow: true,
		}, {
			Path:   secondFilePath,
			Shadow: true,
		},
		}, DestroyedFiles: []string{"654321"}}, staging)
}

func createDirs (t *testing.T,) (localdir string, backupdir string, configdir string) {
	t.Helper()
	tempDir := t.TempDir()
	localDir := path.Join(tempDir, "local")
	backupDir := path.Join(tempDir, "backup")
	configDir := path.Join(localDir, ".backup")
	assert.NoError(t, os.Mkdir(localDir, 0777))
	assert.NoError(t, os.Mkdir(backupDir, 0777))
	assert.NoError(t, os.Mkdir(configDir, 0777))
	return localDir, backupDir, configDir
}

func createLocalFile(t *testing.T, localDir string, fileName string) string {
	t.Helper()
	someFilePath := path.Join(localDir, fileName)
	assert.NoError(t, ioutil.WriteFile(someFilePath, []byte(someFilePath), 0666))
	return someFilePath
}