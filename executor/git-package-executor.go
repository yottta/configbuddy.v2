package executor

import (
	"fmt"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/yottta/configbuddy.v2/model"
	"github.com/yottta/configbuddy.v2/parser"
)

type gitPackageExecutor struct {
	sysPackageExecutor

	packageDestination string
}

func (p *gitPackageExecutor) execute() (err error) {
	log.WithField("PackageName", p.packageAction.PackageName).Info("package action executed (unimplemented)")
	return nil
}

func getPackageDestination(parse parser.Parser, packageAction *model.PackageAction) (string, error) {
	if len(packageAction.Destination) == 0 {
		return "", fmt.Errorf("no destination defined for %s", packageAction.PackageName)
	}
	destination, err := parse.Parse(packageAction.Destination)
	if err != nil {
		return "", err
	}

	if destination[len(destination)-1:] != "/" {
		destination = destination + "/"
	}
	return filepath.Abs(destination)
}
