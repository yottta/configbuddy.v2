package executor

import (
	"fmt"
	"os"
	"path"

	"github.com/andreic92/configbuddy.v2/backup"
	"github.com/andreic92/configbuddy.v2/parser"
	"github.com/andreic92/configbuddy.v2/utils"

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

func (f *fileExecutor) Execute(parse parser.Parser, backupService backup.BackupService) error {
	if f.fileAction == nil {
		return fmt.Errorf("No file action provided")
	}

	command := f.fileAction.Command
	fullPath := fmt.Sprintf("%s/%s", f.fileAction.Source, f.fileAction.FileName)

	finalDestination, err := getFinalDestination(parse, f.fileAction)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path.Dir(finalDestination), os.ModePerm); err != nil {
		return err
	}
	if res := backupService.Backup(finalDestination); res.Error != nil {
		return res.Error
	}
	return utils.ExecuteCommand(fmt.Sprintf("%s %s %s", command, fullPath, finalDestination))
}

func getFinalDestination(parse parser.Parser, fileAction *model.FileAction) (string, error) {
	destination, err := parse.Parse(fileAction.Destination)
	if err != nil {
		return "", err
	}
	if destination[len(destination)-1:] != "/" {
		destination = destination + "/"
	}
	finalFileName := fileAction.FileName
	if fileAction.Hidden {
		finalFileName = "." + finalFileName
	}
	finalDestination := destination + finalFileName
	return finalDestination, nil
}
