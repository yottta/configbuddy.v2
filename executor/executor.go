package executor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/andreic92/configbuddy.v2/model"
	"github.com/andreic92/configbuddy.v2/parser"
	"github.com/ghodss/yaml"

	log "github.com/sirupsen/logrus"
)

type applicationExecutor struct {
	configs   *model.Arguments
	parser    parser.Parser
	finalConf *model.ConfigWrapper
}

func StartConfiguring(config *model.Arguments, parse parser.Parser) (err error) {
	executor := &applicationExecutor{configs: config, parser: parse}
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

func (a *applicationExecutor) executePackages() (err error) {
	dat, err := json.Marshal(a.finalConf)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", dat)
	return
}

func (a *applicationExecutor) executeFiles() (err error) {
	for name, act := range a.finalConf.Config.FileActions {
		fileExecutor := NewFileExecutor(&act, name, a.configs)
		fileExecutor.Execute(a.parser)
	}
	return
}

func (a *applicationExecutor) readConfigs() (err error) {
	if len(a.configs.Configs) == 0 {
		log.Infof("No config files provided. Nothing to do here. Exit...")
		return
	}

	var cfg *model.ConfigWrapper
	for _, filePath := range a.configs.Configs {
		cfg, err = loadConfig(cfg, filePath)
		if err != nil {
			log.WithError(err).Errorf("Error during validate %s", filePath)
			return
		}
	}

	a.finalConf = cfg
	return
}

func loadConfig(appendToThis *model.ConfigWrapper, fileToLoad string) (*model.ConfigWrapper, error) {
	cfg, err := readFile(fileToLoad)
	if err != nil {
		return nil, err
	}
	if appendToThis == nil {
		appendToThis = cfg
	} else {
		for key, val := range cfg.Config.FileActions {
			abs, err := filepath.Abs(cfg.ConfigFileDirectory + "/" + val.Source)
			if err != nil {
				return nil, err
			}
			val.Source = abs
			appendToThis.Config.FileActions[key] = val
		}
		for key, val := range cfg.Config.PackageActions {
			appendToThis.Config.PackageActions[key] = val
		}
	}

	for _, includeFile := range cfg.Config.Includes {
		_, err := loadConfig(appendToThis, cfg.ConfigFileDirectory+"/"+includeFile)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func readFile(filePath string) (*model.ConfigWrapper, error) {
	abs, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

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
