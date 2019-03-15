package executor

import (
	"github.com/andreic92/configbuddy.v2/backup"
	"github.com/andreic92/configbuddy.v2/model"
	"github.com/andreic92/configbuddy.v2/parser"
)

type packageExecutor struct {
	packageAction *model.PackageAction
	args          *model.Arguments

	backupService backup.BackupService
	parser        parser.Parser
}

func newPackageExecutor(packageAction *model.PackageAction, packageName string, args *model.Arguments, parse parser.Parser, backupService backup.BackupService) (*packageExecutor, error) {
	if len(packageAction.PackageName) == 0 {
		packageAction.PackageName = packageName
	}

	// // source path
	// var fullPath string
	// if len(fileAction.Source) > 0 {
	// 	fullPath = fmt.Sprintf("%s/%s", strings.TrimRight(fileAction.Source, "/"), fileAction.FileName)
	// } else {
	// 	fullPath = fmt.Sprintf("./%s", fileAction.FileName)
	// }

	// // target path
	// finalDestination, err := getFinalDestination(parse, fileAction)
	// if err != nil {
	// 	return nil, err
	// }

	return &packageExecutor{
		packageAction: packageAction,
		args:          args,

		backupService: backupService,
		parser:        parse,
	}, nil
}

func (p *packageExecutor) execute() (err error) {
	return nil
}
