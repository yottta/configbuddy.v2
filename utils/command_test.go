package utils

import (
	"fmt"
	"testing"

	ast "github.com/stretchr/testify/assert"
)

const (
	testingFileName = "just_a_dummy_file_name"
)

func TestCommand(t *testing.T) {
	assert := ast.New(t)

	err := ExecuteCommand(fmt.Sprintf("touch %s", testingFileName))
	assert.NoError(err)

	err = RemoveResource(testingFileName)
	assert.NoError(err)

	err = ExecuteCommand("touch")
	assert.Error(err)
	assert.Equal("invalid command, arguments missing. command: 'touch'", err.Error())
}
