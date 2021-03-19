package commands_test

import (
	"os"
	"path"
	"testing"

	"github.com/SBanczyk/backup/commands"
	"github.com/SBanczyk/backup/model"
	"github.com/stretchr/testify/assert"
)

func TestUnstage(t *testing.T) {
	tempDir := t.TempDir()
	stagingDir := path.Join(tempDir, ".backup")
	assert.NoError(t, os.Mkdir(stagingDir, 0777))
	stagingDir = path.Join(stagingDir, "staging")
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
	err = commands.Unstage(tempDir, []string{"qwerty", "123456"})
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
