package backup

import (
	"fmt"
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
	assert.Contains(res.Error.Error(), "path cannot be empty")
	assert.False(res.Performed)
}

func TestBackupErrorFromStrategy(t *testing.T) {
	assert := ast.New(t)
	assert.True(true)

	bakServ := defaultBackupService{
		backupActivated: true,
		backupDirectory: "",
		backupStrategy:  &mockExtractErrorFileStrategy{},
	}

	testFile := "test_file"
	_, err := os.Create(testFile)
	assert.NoError(err)

	res := bakServ.Backup(testFile)
	assert.Error(res.Error)
	assert.Contains(res.Error.Error(), "mock error")
	assert.False(res.Performed)
}

func TestBackupOverAlreadyExistingFile(t *testing.T) {
	assert := ast.New(t)
	assert.True(true)

	testFile := "test_file"
	bakServ := defaultBackupService{
		backupActivated: true,
		backupDirectory: "",
		backupStrategy:  &mockExtractAlreadyExistingFileStrategy{alreadyExistingFilePath: testFile},
	}

	_, err := os.Create(testFile)
	assert.NoError(err)

	res := bakServ.Backup(testFile)
	assert.Error(res.Error)
	assert.Contains(res.Error.Error(), "exit status 1")
	assert.False(res.Performed)

	deleteResource(assert, testFile)
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

type mockExtractErrorFileStrategy struct {
}

func (m *mockExtractErrorFileStrategy) extractTargetPath(resourcePath string) (string, error) {
	return "", fmt.Errorf("mock error")
}

type mockExtractAlreadyExistingFileStrategy struct {
	alreadyExistingFilePath string
}

func (m *mockExtractAlreadyExistingFileStrategy) extractTargetPath(resourcePath string) (string, error) {
	return m.alreadyExistingFilePath, nil
}
