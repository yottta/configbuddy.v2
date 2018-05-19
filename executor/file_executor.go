package executor

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/andreic92/configbuddy.v2/backup"
	"github.com/andreic92/configbuddy.v2/parser"
	"github.com/andreic92/configbuddy.v2/utils"

	"github.com/andreic92/configbuddy.v2/model"
)

type FileExecutor struct {
	fileAction *model.FileAction
	args       *model.Arguments

	backupService backup.BackupService
	parser        parser.Parser

	fullPath         string
	finalDestination string
}

func NewFileExecutor(fileAction *model.FileAction,
	fileName string,
	args *model.Arguments,
	parse parser.Parser,
	backupService backup.BackupService) (*FileExecutor, error) {
	if len(fileAction.FileName) == 0 {
		fileAction.FileName = fileName
	}

	// source path
	var fullPath string
	if len(fileAction.Source) > 0 {
		fullPath = fmt.Sprintf("%s/%s", strings.TrimRight(fileAction.Source, "/"), fileAction.FileName)
	} else {
		fullPath = fmt.Sprintf("./%s", fileAction.FileName)
	}

	// target path
	finalDestination, err := getFinalDestination(parse, fileAction)
	if err != nil {
		return nil, err
	}

	return &FileExecutor{
		fileAction: fileAction,
		args:       args,

		backupService: backupService,
		parser:        parse,

		fullPath:         fullPath,
		finalDestination: finalDestination,
	}, nil
}

func (f *FileExecutor) Execute() (err error) {
	err = f.createDirectoriesStructure()
	if err != nil {
		return err
	}

	err = f.backup()
	if err != nil {
		return err
	}

	return utils.ExecuteCommand(f.getCommand())
}

func (f *FileExecutor) getCommand() string {
	return fmt.Sprintf("%s %s %s", f.fileAction.Command, f.fullPath, f.finalDestination)
}

func (f *FileExecutor) backup() error {
	return f.backupService.Backup(f.finalDestination).Error
}

func (f *FileExecutor) createDirectoriesStructure() error {
	return os.MkdirAll(path.Dir(f.finalDestination), os.ModePerm)
}

func getFinalDestination(parse parser.Parser, fileAction *model.FileAction) (string, error) {
	if len(fileAction.Destination) == 0 {
		return "", fmt.Errorf("No destination defined for %s", fileAction.FileName)
	}
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
	return filepath.Abs(finalDestination)
}
