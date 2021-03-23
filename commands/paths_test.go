package commands

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckBaseDir(t *testing.T) {
	tempDir := t.TempDir()
	localDir := path.Join(tempDir, "xyz", "abc", "qwe")
	assert.NoError(t, os.MkdirAll(localDir, 0777))
	secondDir := path.Join(tempDir, "qaz", "wsx")
	assert.NoError(t, os.MkdirAll(secondDir, 0777))
	backupDir := path.Join(tempDir, "xyz", ".backup")
	assert.NoError(t, os.MkdirAll(backupDir, 0777))
	baseDir := path.Join(tempDir, "xyz")
	baseDirChecked, err := checkBaseDir(localDir)
	assert.NoError(t, err)
	assert.EqualValues(t, baseDir, baseDirChecked)
	_, err = checkBaseDir(secondDir)
	assert.Error(t, err)
	baseDirCheck, err := checkBaseDir(baseDir)
	assert.NoError(t, err)
	assert.EqualValues(t, baseDir, baseDirCheck)
	_, err = checkBaseDir(tempDir)
	assert.Error(t, err)
}
