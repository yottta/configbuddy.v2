package backup

import (
	"strings"
	"testing"

	ast "github.com/stretchr/testify/assert"
)

const (
	copyToBakFileInResource     = "/root/.data"
	copyToBakFileExpectedPrefix = "/root/"
	copyToBakFileExpectedSufix  = "_.data"
	copyToBakFileExpectedLength = len(copyToBakFileInResource) + 10 // approx.

	copyToDirFileInResource        = "/root/.data"
	copyToDirFileInBackupDirectory = "/root/bak/"
	copyToDirFileExpectedPrefix    = copyToDirFileInBackupDirectory
	copyToDirFileExpectedSufix     = "_.data"
	copyToDirFileExpectedLength    = len(copyToDirFileInResource) + 10 // approx.
)

func TestDefaultStrategy(t *testing.T) {
	assert := ast.New(t)
	strategier := getStrategier(backupStrategyDisabled, "")
	assert.Nil(strategier)
}

func TestCopyToBakStrategy(t *testing.T) {
	assert := ast.New(t)

	strategier := getStrategier(backupStrategyBakFile, "")
	assert.NotNil(strategier)

	res, err := strategier.extractTargetPath(copyToBakFileInResource)
	assert.NoError(err)
	assert.True(strings.HasPrefix(res, copyToBakFileExpectedPrefix))
	assert.True(strings.HasSuffix(res, copyToBakFileExpectedSufix))
	assert.True(len(res) > copyToBakFileExpectedLength)

	res, err = strategier.extractTargetPath("")
	assert.Error(err)
	assert.Contains(err.Error(), "invalid resource name")
	assert.Equal(res, "")
}

func TestCopyToDirStrategy(t *testing.T) {
	assert := ast.New(t)

	strategier := getStrategier(backupStrategyCopyToDirectory, copyToDirFileInBackupDirectory)
	assert.NotNil(strategier)

	res, err := strategier.extractTargetPath(copyToDirFileInResource)
	assert.NoError(err)
	assert.True(strings.HasPrefix(res, copyToDirFileExpectedPrefix))
	assert.True(strings.HasSuffix(res, copyToDirFileExpectedSufix))
	assert.True(len(res) > copyToDirFileExpectedLength)

	res, err = strategier.extractTargetPath("")
	assert.Error(err)
	assert.Contains(err.Error(), "invalid resource name")
	assert.Equal(res, "")
}
