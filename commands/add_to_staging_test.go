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
	tempDir := t.TempDir()
	localDir := path.Join(tempDir, "local")
	backupDir := path.Join(tempDir, "backup")
	assert.NoError(t, os.Mkdir(localDir, 0777))
	assert.NoError(t, os.Mkdir(backupDir, 0777))
	someFilePath := path.Join(localDir, "someFile")
	assert.NoError(t, ioutil.WriteFile(someFilePath, []byte("123456"), 0666))
	err := commands.InitFs(localDir, backupDir)
	assert.NoError(t, err)
	err = commands.AddToStaging(localDir, []string{someFilePath}, false)
	assert.NoError(t, err)
	stagingPath := path.Join(localDir, ".backup", "staging")
	staging, err := model.LoadStaging(stagingPath)
	assert.NoError(t, err)
	assert.Len(t, staging.StagingFiles, 1)
	assert.Contains(t, staging.StagingFiles, model.StagingPath{
		Path:   someFilePath,
		Shadow: false,
	})
	err = commands.AddToStaging(localDir, []string{someFilePath}, true)
	assert.NoError(t, err)
	staging, err = model.LoadStaging(stagingPath)
	assert.NoError(t, err)
	assert.Len(t, staging.StagingFiles, 1)
	assert.Contains(t, staging.StagingFiles, model.StagingPath{
		Path:   someFilePath,
		Shadow: true,
	})
}
