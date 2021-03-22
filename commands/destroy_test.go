package commands_test

import (
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
					Path:   firstFilePath,
					Shadow: false,
				},
			},
		},
	})
	assert.NoError(t, err)
	err = model.SaveStaging(stagingPath, &model.Staging{
		StagingFiles: []model.StagingPath{
			{
				Path:   firstFilePath,
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
		DestroyedFiles: []string{firstFilePath},
	}, staging)

}
