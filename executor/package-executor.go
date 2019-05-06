package executor

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/yottta/configbuddy.v2/backup"
	"github.com/yottta/configbuddy.v2/model"
	"github.com/yottta/configbuddy.v2/parser"
)

const (
	gitPackageSource = "git"
)

type packageExecutor interface {
	execute() (err error)
}
type sysPackageExecutor struct {
	packageAction *model.PackageAction
	args          *model.Arguments

	backupService backup.BackupService
	parser        parser.Parser
}

func newPackageExecutor(packageAction *model.PackageAction, args *model.Arguments, parse parser.Parser, backupService backup.BackupService) (packageExecutor, error) {
	if len(packageAction.PackageName) == 0 {
		return nil, fmt.Errorf("empty package action name")
	}
	if packageAction.Source == gitPackageSource {
		gitPackageDestination, err := getPackageDestination(parse, packageAction)
		if err != nil {
			return nil, err
		}

		if len(packageAction.URL) == 0 {
			return nil, fmt.Errorf("url empty for %s package action", packageAction.PackageName)
		}

		return &gitPackageExecutor{
			sysPackageExecutor: sysPackageExecutor{
				packageAction: packageAction,
				args:          args,

				backupService: backupService,
				parser:        parse,
			},

			packageDestination: gitPackageDestination,
		}, nil
	}
	return &sysPackageExecutor{
		packageAction: packageAction,
		args:          args,

		backupService: backupService,
		parser:        parse,
	}, nil
}

func (p *sysPackageExecutor) execute() (err error) {
	log.WithField("PackageName", p.packageAction.PackageName).Info("package action executed (unimplemented)")
	return nil
}

func (p *sysPackageExecutor) command() string {
	var buff strings.Builder
	if p.packageAction.Sudo {
		buff.WriteString("sudo ")
	}
	buff.WriteString(parser.PackageManagerPlaceholder)

	return "" //fmt.Sprintf("%s %s %s", f.fileAction.Command, f.fullPath, f.finalDestination)
}
