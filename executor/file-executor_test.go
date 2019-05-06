package executor

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yottta/configbuddy.v2/parser"

	"github.com/yottta/configbuddy.v2/backup"

	"github.com/yottta/configbuddy.v2/model"

	ast "github.com/stretchr/testify/assert"
)

func TestNewFileExecutor(t *testing.T) {
	assert := ast.New(t)

	parser, err := parser.NewParser()
	assert.NoError(err)

	currentDirectory, err := filepath.Abs("./")
	assert.NoError(err)

	testSuite := []struct {
		testName                 string
		fileAction               model.FileAction
		expectedFileName         string
		expectedFileDestination  string
		expectedFileSource       string
		expectedErrorMsgSnapshot string
	}{
		{
			testName: "test common file act",
			fileAction: model.FileAction{
				Command:     "ln -s",
				Destination: ".",
				FileName:    "test_file",
				Source:      "./",
			},
			expectedFileName:        "test_file",
			expectedFileSource:      "./test_file",
			expectedFileDestination: strings.TrimRight(currentDirectory, "/") + "/test_file",
		},
		{
			testName: "without source",
			fileAction: model.FileAction{
				Command:     "ln -s",
				Destination: ".",
				FileName:    "test_file",
			},
			expectedFileName:        "test_file",
			expectedFileSource:      "./test_file",
			expectedFileDestination: strings.TrimRight(currentDirectory, "/") + "/test_file",
		},
		{
			testName: "without destination",
			fileAction: model.FileAction{
				Command:  "ln -s",
				FileName: "test_file",
			},
			expectedErrorMsgSnapshot: "no destination defined for test_file",
		},
		{
			testName: "wrong destination placeholder",
			fileAction: model.FileAction{
				Command:     "ln -s",
				FileName:    "test_file",
				Destination: "$#HOME",
			},
			expectedErrorMsgSnapshot: "unclosed action",
		},
	}

	for _, testCase := range testSuite {
		fileExecutor, err := newFileExecutor(&testCase.fileAction, testCase.fileAction.FileName, &dummyArguments, parser, &mockNoActionBackup{})
		if len(testCase.expectedErrorMsgSnapshot) == 0 {
			assert.NoError(err, fmt.Sprintf("%s -> no error expected", testCase.testName))
			assert.NotNil(fileExecutor, fmt.Sprintf("%s -> file executor expeccted", testCase.testName))

			assert.Equal(testCase.expectedFileName, fileExecutor.fileAction.FileName, fmt.Sprintf("%s -> file name", testCase.testName))
			assert.Equal(testCase.expectedFileDestination, fileExecutor.finalDestination, fmt.Sprintf("%s -> final destination", testCase.testName))
			assert.Equal(testCase.expectedFileSource, fileExecutor.fullPath, fmt.Sprintf("%s -> full path", testCase.testName))
		} else {
			assert.Error(err, fmt.Sprintf("%s -> error expected", testCase.testName))
			assert.Contains(err.Error(), testCase.expectedErrorMsgSnapshot, fmt.Sprintf("%s -> error contains", testCase.testName))
			assert.Nil(fileExecutor, fmt.Sprintf("%s -> nil file executor", testCase.testName))
		}
	}
}

var dummyArguments = model.Arguments{}

type mockNoActionBackup struct{}

func (mb *mockNoActionBackup) Backup(resourcePath string) backup.BackupResult {
	return backup.BackupResult{}
}
