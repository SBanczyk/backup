package commands_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/SBanczyk/backup/commands"
	"github.com/SBanczyk/backup/model"
	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	testDir, err := filepath.Abs("testdata")
	assert.NoError(t, err)
	os.RemoveAll(testDir)
	err = os.MkdirAll(testDir, 0777)
	assert.NoError(t, err)
	localDir := path.Join(testDir, "local")
	backupDir := path.Join(testDir, "backup")
	configDir := path.Join(localDir, ".backup")
	stagingPath := path.Join(configDir, "staging")
	backupPath := path.Join(configDir, "backupfiles")
	assert.NoError(t, os.Mkdir(localDir, 0777))
	assert.NoError(t, os.Mkdir(backupDir, 0777))
	newLocalPath := path.Join(localDir, "newLocal")
	assert.NoError(t, os.Mkdir(newLocalPath, 0777))
	err = commands.InitFs(localDir, backupDir)
	assert.NoError(t, err)
	createLocalFile(t, localDir, "modifiedFile")
	createLocalFile(t, localDir, "modifiedUnstage")
	createLocalFile(t, localDir, "untrackedFile")
	createLocalFile(t, newLocalPath, "modifiedFile1")
	createLocalFile(t, newLocalPath, "modifiedUnstage1")
	createLocalFile(t, newLocalPath, "untrackedFile1")
	createLocalFile(t, newLocalPath, "newFile1")
	err = model.SaveStaging(stagingPath, &model.Staging{
		StagingFiles: []model.StagingPath{
			{
				Path:   "missingFromStaging",
				Shadow: true,
			},
			{
				Path:   "newFile",
				Shadow: true,
			},
			{
				Path:   "modifiedFile",
				Shadow: true,
			},
			{
				Path:   path.Join("newLocal", "missingFromStaging1"),
				Shadow: true,
			},
			{
				Path:   path.Join("newLocal", "newFile1"),
				Shadow: true,
			},
			{
				Path:   path.Join("newLocal", "modifiedFile1"),
				Shadow: true,
			},
		},
		DestroyedFiles: []string{"destroyedFile", path.Join("newLocal", "destroyedFile1")},
	})
	createLocalFile(t, localDir, "newFile")

	assert.NoError(t, err)
	err = model.SaveBackup(backupPath, &model.Backup{
		Files: map[string][]model.BackupPath{
			"123456": {
				{
					Path:   "modifiedFile",
					Shadow: true,
				},
				{
					Path:   "destroyedFile",
					Shadow: true,
				},
				{
					Path:   "missingFromBack",
					Shadow: false,
				},
				{
					Path:   "modifiedUnstage",
					Shadow: true,
				},
				{
					Path:   path.Join("newLocal", "modifiedFile1"),
					Shadow: true,
				},
				{
					Path:   path.Join("newLocal", "destroyedFile1"),
					Shadow: true,
				},
				{
					Path:   path.Join("newLocal", "missingFromBack1"),
					Shadow: false,
				},
				{
					Path:   path.Join("newLocal", "modifiedUnstage1"),
					Shadow: true,
				},
			},
		},
	})
	assert.NoError(t, err)
}
