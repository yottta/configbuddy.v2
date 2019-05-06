package executor

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"

	"github.com/yottta/configbuddy.v2/backup"
	"github.com/yottta/configbuddy.v2/model"
	"github.com/yottta/configbuddy.v2/parser"
)

type applicationExecutor struct {
	configs       *model.Arguments
	parser        parser.Parser
	finalConf     *model.ConfigWrapper
	backupService backup.BackupService
}

func StartConfiguring(config *model.Arguments, parse parser.Parser, backupService backup.BackupService) (err error) {
	executor := &applicationExecutor{configs: config, parser: parse, backupService: backupService}
	err = executor.readConfigs()
	if err != nil {
		return
	}

	err = executor.executePackages()
	if err != nil {
		return
	}

	err = executor.executeFiles()
	if err != nil {
		return
	}
	return nil
}

func (a *applicationExecutor) readConfigs() (err error) {
	if len(a.configs.Configs) == 0 {
		err = fmt.Errorf("no config files provided. nothing to do here")
		return
	}

	var cfg *model.ConfigWrapper
	for _, filePath := range a.configs.Configs {
		cfg, err = loadConfig(cfg, filePath)
		if err != nil {
			log.WithError(err).Errorf("error during validate %s", filePath)
			return
		}
	}

	a.finalConf = cfg
	return
}

func (a *applicationExecutor) executePackages() (err error) {
	return nil
	// for name, act := range a.finalConf.Config.PackageActions {
	// 	packageExecutor, err := newPackageExecutor(&act, name, a.configs, a.parser, a.backupService)
	// 	if err != nil {
	// 		log.WithError(err).WithField("file action", act).Error("Error during processing fileAction")
	// 		continue
	// 	}
	// 	err = packageExecutor.execute()
	// 	if err != nil {
	// 		log.WithError(err).WithField("file action", act).Error("Error during processing fileAction")
	// 	}
	// }
	// return
}

func (a *applicationExecutor) executeFiles() (err error) {
	for name, act := range a.finalConf.Config.FileActions {
		skipAct, err := checkExecutable(a.parser, act.ConditionalAction)
		if err != nil {
			return err
		}
		if skipAct {
			log.WithField("Action", name).Info("File action skipped")
			continue
		}

		fileExecutor, err := newFileExecutor(&act, name, a.configs, a.parser, a.backupService)
		if err != nil {
			log.WithError(err).WithField("file action", act).Error("error during processing fileAction")
			continue
		}
		err = fileExecutor.execute()
		if err != nil {
			log.WithError(err).WithField("file action", act).Error("error during processing fileAction")
		}
	}
	return
}

func loadConfig(appendToThis *model.ConfigWrapper, fileToLoad string) (*model.ConfigWrapper, error) {
	cfg, err := readFile(fileToLoad)
	if err != nil {
		return nil, err
	}
	if appendToThis == nil {
		appendToThis = cfg
		err = appendActionsToGlobalConfig(cfg, appendToThis)
		if err != nil {
			return nil, err
		}
	} else {
		err = appendActionsToGlobalConfig(cfg, appendToThis)
		if err != nil {
			return nil, err
		}
	}
	res := appendToThis

	for _, includeFile := range cfg.Config.Includes {
		log.WithField("file", includeFile).Debug("include config")
		_, err := loadConfig(appendToThis, cfg.ConfigFileDirectory+"/"+includeFile)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func appendActionsToGlobalConfig(cfg *model.ConfigWrapper, appendToThis *model.ConfigWrapper) error {
	// file actions
	if appendToThis.Config.FileActions == nil {
		appendToThis.Config.FileActions = make(map[string]model.FileAction)
	}
	for key, val := range cfg.Config.FileActions {
		abs, err := filepath.Abs(cfg.ConfigFileDirectory + "/" + val.Source)
		if err != nil {
			return err
		}
		val.Source = abs
		if strings.HasPrefix(val.Destination, ".") { // if the destination path is relative
			val.Destination = cfg.ConfigFileDirectory + "/" + val.Destination
		}
		appendToThis.Config.FileActions[key] = val
	}

	// package actions
	if appendToThis.Config.PackageActions == nil {
		appendToThis.Config.PackageActions = []model.PackageAction{}
	}
	for _, val := range cfg.Config.PackageActions {
		appendToThis.Config.PackageActions = append(appendToThis.Config.PackageActions, val)
	}
	return nil
}

func readFile(filePath string) (*model.ConfigWrapper, error) {
	abs, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	log.WithField("file", abs).Debug("reading file")
	bytes, err := ioutil.ReadFile(abs)
	if err != nil {
		return nil, err
	}

	var c model.Config
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}
	return &model.ConfigWrapper{
		Config:              &c,
		ConfigFilePath:      abs,
		ConfigFileDirectory: filepath.Dir(abs),
	}, nil
}

func checkExecutable(parse parser.Parser, conditionalAction model.ConditionalAction) (bool, error) {
	if len(conditionalAction.Condition()) == 0 {
		return false, nil
	}

	val, err := parse.EvaluateCondition(conditionalAction.Condition())
	return !val, err
}
