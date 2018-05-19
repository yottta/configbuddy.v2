package parser

import (
	"os/user"
	"testing"

	ast "github.com/stretchr/testify/assert"
)

const (
	noTemplate                     = "just a dummy string"
	notDefaultPlaceholdersTemplate = "this is the users's home directory: ##HOME##"
	invalidTemplate                = "this is the users's home directory: $#HOME"
	okTemplate                     = "this is the users's home directory: $#HOME#$"
)

func TestParser(t *testing.T) {
	assert := ast.New(t)

	parser, err := NewParser()
	assert.NoError(err)
	assert.NotNil(parser)

	res, err := parser.Parse(noTemplate)
	assert.NoError(err)
	assert.Equal(noTemplate, res)

	res, err = parser.Parse(notDefaultPlaceholdersTemplate)
	assert.NoError(err)
	assert.Equal(notDefaultPlaceholdersTemplate, res)

	res, err = parser.Parse(invalidTemplate)
	assert.Error(err)
	assert.Contains(err.Error(), "unclosed action")

	usr, err := user.Current()
	assert.NoError(err, "Couldn't get the user info")
	res, err = parser.Parse(okTemplate)
	assert.NoError(err)
	assert.NotContains(res, defaultParserPrefix)
	assert.Contains(res, usr.HomeDir)
}
