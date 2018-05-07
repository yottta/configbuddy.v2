package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andreic92/configbuddy.v2/model"
	"github.com/andreic92/configbuddy.v2/utils"
)

type BackupService interface {
	Backup(resourcePath string) BackupResult
}

type backupStrategy int

const (
	backupStrategyDisabled        backupStrategy = 0
	backupStrategyBakFile         backupStrategy = 1
	backupStrategyCopyToDirectory backupStrategy = 2
)

type defaultBackupService struct {
	backupStrategy  backupStrategier
	backupDirectory string
	backupActivated bool
}

type BackupResult struct {
	Performed   bool
	InitialPath string
	FinalPath   string
	Error       error
}

func NewBackupService(config *model.Arguments) (BackupService, error) {
	strategy := backupStrategyDisabled
	backupActivated := false
	backupDirectory := config.BackupDirectory
	if len(backupDirectory) > 0 {
		if err := checkDirectory(backupDirectory); err != nil {
			return nil, err
		}
		backupActivated = true
		strategy = backupStrategyCopyToDirectory
	} else if config.BackupActivated {
		backupActivated = true
		strategy = backupStrategyBakFile
	}

	return &defaultBackupService{
		backupStrategy:  getStrategier(strategy, backupDirectory),
		backupDirectory: backupDirectory,
		backupActivated: backupActivated,
	}, nil
}

func (d *defaultBackupService) Backup(resourcePath string) BackupResult {
	sourceAbsPath, err := filepath.Abs(resourcePath)
	if err != nil {
		return BackupResult{
			Performed: false,
			Error:     err,
		}
	}
	// check if the file requested for backup exists
	_, err = os.Stat(sourceAbsPath)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			return BackupResult{
				Performed: false,
			}
		}
		return BackupResult{
			Performed: false,
			Error:     err,
		}
	}
	targetAbsPath := ""
	if d.backupActivated {
		targetAbsPath, err = d.backupStrategy.extractTargetPath(sourceAbsPath)
		if err != nil {
			return BackupResult{
				InitialPath: sourceAbsPath,
				Performed:   false,
				Error:       err,
			}
		}

		err = utils.ExecuteCommand(fmt.Sprintf("cp -RL %s %s", sourceAbsPath, targetAbsPath))
		if err != nil {
			return BackupResult{
				InitialPath: sourceAbsPath,
				FinalPath:   targetAbsPath,
				Performed:   false,
				Error:       err,
			}
		}
	}

	err = utils.RemoveResource(sourceAbsPath)
	return BackupResult{
		InitialPath: sourceAbsPath,
		FinalPath:   targetAbsPath,
		Performed:   d.backupActivated,
		Error:       err,
	}
}

func checkDirectory(directory string) error {
	fileInfo, err := os.Stat(directory)
	if err != nil {
		return os.MkdirAll(directory, os.ModePerm)
	}
	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		return nil
	default:
		return fmt.Errorf("Given path (%s) is not a directory", directory)
	}
}
