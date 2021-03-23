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
	localDir, _, stagingDir, _ := createDirs(t)
	qwertyPath := path.Join(localDir, "qwerty")
	oneToSixPath := path.Join(localDir, "123456")
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
	err = commands.Unstage(localDir, []string{qwertyPath, oneToSixPath})
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

func TestUnstageDeep(t *testing.T) {
	localDir, _, stagingPath, _ := createDirs(t)
	firstFilePath := path.Join("qwer", "qaz", "firstFile")
	err := model.SaveStaging(stagingPath, &model.Staging{
		StagingFiles: []model.StagingPath{{
			Path:   firstFilePath,
			Shadow: true,
		},
		}, DestroyedFiles: []string{}})
	assert.NoError(t, err)
	newDir := path.Join(localDir, "qwer", "qaz")
	assert.NoError(t, os.MkdirAll(newDir, 0777))
	firstFilePath = createLocalFile(t, newDir, "firstFile")
	err = commands.Unstage(newDir, []string{firstFilePath})
	assert.NoError(t, err)
	staging, err := model.LoadStaging(stagingPath)
	assert.NoError(t, err)
	assert.EqualValues(t, &model.Staging{
		StagingFiles:   []model.StagingPath{},
		DestroyedFiles: []string{},
	}, staging)
}
