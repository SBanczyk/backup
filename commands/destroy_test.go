package commands_test

import (
	"os"
	"path"
	"testing"

	"github.com/SBanczyk/backup/commands"
	"github.com/SBanczyk/backup/model"
	"github.com/stretchr/testify/assert"
)

func TestDestroy(t *testing.T) {
	localDir, _, stagingPath, backupPath := createDirs(t)
	firstFilePath := path.Join(localDir, "firstFile")
	err := model.SaveBackup(backupPath, &model.Backup{
		Files: map[string][]model.BackupPath{
			"abcdef": {
				{
					Path:   "firstFile",
					Shadow: false,
				},
			},
		},
	})
	assert.NoError(t, err)
	err = model.SaveStaging(stagingPath, &model.Staging{
		StagingFiles: []model.StagingPath{
			{
				Path:   "firstFile",
				Shadow: false,
			},
		},
	})
	assert.NoError(t, err)
	secondFilePath := path.Join(localDir, "secondFile")
	err = commands.Destroy(localDir, []string{firstFilePath, secondFilePath})
	assert.NoError(t, err)
	staging, err := model.LoadStaging(stagingPath)
	assert.NoError(t, err)
	assert.EqualValues(t, &model.Staging{
		StagingFiles:   []model.StagingPath{},
		DestroyedFiles: []string{"firstFile"},
	}, staging)

}

func TestDestroyDeep(t *testing.T) {
	localDir, _, stagingPath, backupPath := createDirs(t)
	newDir := path.Join(localDir, "qwer", "qaz")
	assert.NoError(t, os.MkdirAll(newDir, 0777))
	firstFilePath := path.Join("qwer", "qaz", "firstFile")
	err := model.SaveBackup(backupPath, &model.Backup{
		Files: map[string][]model.BackupPath{
			"abcdef": {
				{
					Path:   firstFilePath,
					Shadow: false,
				},
			},
		},
	})
	assert.NoError(t, err)
	err = model.SaveStaging(stagingPath, &model.Staging{})
	assert.NoError(t, err)
	firstFileDir := createLocalFile(t, newDir, "firstFile")
	err = commands.Destroy(newDir, []string{firstFileDir})
	assert.NoError(t, err)
	staging, err := model.LoadStaging(stagingPath)
	assert.NoError(t, err)
	assert.EqualValues(t, &model.Staging{
		StagingFiles:   []model.StagingPath{},
		DestroyedFiles: []string{firstFilePath},
	}, staging)
}
