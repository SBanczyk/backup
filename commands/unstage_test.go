package commands_test

import (
	"testing"

	"github.com/SBanczyk/backup/commands"
	"github.com/SBanczyk/backup/model"
	"github.com/stretchr/testify/assert"
)

func TestUnstage(t *testing.T) {
	localDir, _, stagingDir, _ := createDirs(t)
	err := model.SaveStaging(stagingDir, &model.Staging{
		StagingFiles: []model.StagingPath{{
			Path:   "qwerty",
			Shadow: true,
		}, {
			Path:   "1qaz",
			Shadow: false,
		}},
		DestroyedFiles: []string{"123456", "654321"},
	})
	assert.NoError(t, err)
	err = commands.Unstage(localDir, []string{"qwerty", "123456"})
	assert.NoError(t, err)
	staging, err := model.LoadStaging(stagingDir)
	assert.NoError(t, err)
	assert.EqualValues(t, &model.Staging{
		StagingFiles: []model.StagingPath{{
			Path:   "1qaz",
			Shadow: false,
		},
		}, DestroyedFiles: []string{"654321"}}, staging)
}
