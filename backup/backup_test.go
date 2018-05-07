package backup

import (
	"os"
	"testing"

	"github.com/andreic92/configbuddy.v2/model"

	ast "github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	assert := ast.New(t)
	assert.True(true)

	params := &model.Arguments{
		BackupActivated: true,
		BackupDirectory: "",
	}
	bakServ, err := NewBackupService(params)
	assert.NoError(err)
	assert.NotNil(bakServ)

	params.BackupDirectory = "relative_path/bak_dir"
	bakServ, err = NewBackupService(params)
	assert.NoError(err)
	assert.NotNil(bakServ)

	assertDir(assert, params.BackupDirectory)
	deleteResource(assert, "relative_path")

	params.BackupDirectory = "backup.go"
	bakServ, err = NewBackupService(params)
	assert.Error(err)
	assert.Contains(err.Error(), "is not a directory")
	assert.Nil(bakServ)
}

func TestBackupBakFile(t *testing.T) {
	assert := ast.New(t)
	assert.True(true)

	params := &model.Arguments{
		BackupActivated: true,
		BackupDirectory: "",
	}
	bakServ, err := NewBackupService(params)
	assert.NoError(err)
	assert.NotNil(bakServ)

	testFile := "test_file"
	_, err = os.Create(testFile)
	assert.NoError(err)

	res := bakServ.Backup(testFile)
	assert.NoError(res.Error)
	assert.True(res.Performed)
	assertFile(assert, res.FinalPath)
	deleteResource(assert, res.FinalPath)
}

func TestBackupBakFileNonExistentSource(t *testing.T) {
	assert := ast.New(t)
	assert.True(true)

	params := &model.Arguments{
		BackupActivated: true,
		BackupDirectory: "",
	}
	bakServ, err := NewBackupService(params)
	assert.NoError(err)
	assert.NotNil(bakServ)

	testFile := "test_file"

	res := bakServ.Backup(testFile)
	assert.NoError(res.Error)
	assert.False(res.Performed)
	assertNoFile(assert, res.FinalPath)
}

func TestBackupBakFileEmptyResourceName(t *testing.T) {
	assert := ast.New(t)
	assert.True(true)

	params := &model.Arguments{
		BackupActivated: true,
		BackupDirectory: "",
	}
	bakServ, err := NewBackupService(params)
	assert.NoError(err)
	assert.NotNil(bakServ)

	testFile := ""

	res := bakServ.Backup(testFile)
	assert.Error(res.Error)
	assert.False(res.Performed)
}

func assertFile(assert *ast.Assertions, filePath string) {
	fi, err := os.Stat(filePath)
	assert.NoError(err)
	assert.NotNil(fi)

	assert.False(fi.IsDir())
}

func assertNoFile(assert *ast.Assertions, filePath string) {
	fi, err := os.Stat(filePath)
	assert.Error(err)
	assert.Contains(err.Error(), "no such file or directory")
	assert.Nil(fi)
}

func assertDir(assert *ast.Assertions, filePath string) {
	fi, err := os.Stat(filePath)
	assert.NoError(err)
	assert.NotNil(fi)

	assert.True(fi.IsDir())
}

func deleteResource(assert *ast.Assertions, path string) {
	assert.NoError(os.RemoveAll(path))
}
