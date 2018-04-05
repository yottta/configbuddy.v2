package executor

import (
	"fmt"

	"github.com/andreic92/configbuddy.v2/parser"

	"github.com/andreic92/configbuddy.v2/model"
)

type fileExecutor struct {
	fileAction *model.FileAction
	args       *model.Arguments
}

func NewFileExecutor(fileAction *model.FileAction, fileName string, args *model.Arguments) *fileExecutor {
	if len(fileAction.FileName) == 0 {
		fileAction.FileName = fileName
	}
	return &fileExecutor{
		fileAction: fileAction,
		args:       args,
	}
}

func (f *fileExecutor) Execute(parse parser.Parser) error {
	if f.fileAction == nil {
		return fmt.Errorf("No file action provided")
	}

	// fullPath := fmt.Sprintf("%s/%s", f.fileAction.Source, f.fileAction.FileName)
	// command := f.fileAction.Command
	// destination := f.fileAction.Destination
	destination, err := parse.Parse(f.fileAction.Destination)
	if err != nil {
		return err
	}
	if destination[len(destination)-1:] != "/" {
		destination = destination + "/"
	}
	fmt.Println(destination)

	return nil
}
