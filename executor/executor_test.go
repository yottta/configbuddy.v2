package executor

import (
	"testing"

	ast "github.com/stretchr/testify/assert"
)

func Test01TestReadFile(t *testing.T) {
	assert := ast.New(t)
	// this should be an error
	res, err := readFile("test.yml")
	assert.Error(err)
	assert.Nil(res)
}
