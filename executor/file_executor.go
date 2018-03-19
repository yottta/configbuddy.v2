package executor

import (
	"fmt"

	"github.com/andreic92/configbuddy.v2/model"
)

type fileExecutor struct {
	fileAction *model.FileAction
}

func NewFileExecutor(fileAction *model.FileAction, fileName string) *fileExecutor {
	if len(fileAction.FileName) == 0 {
		fileAction.FileName = fileName
	}
	return &fileExecutor{
		fileAction: fileAction,
	}
}

func (f *fileExecutor) Execute() error {
	if f.fileAction == nil {
		return fmt.Errorf("No file action provided")
	}

	// source := f.fileAction.Source
	// fileName := f.fileAction.FileName

	return nil
}
