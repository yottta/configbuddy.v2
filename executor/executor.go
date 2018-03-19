package executor

import (
	"io/ioutil"
	"path/filepath"

	"github.com/andreic92/configbuddy.v2/model"
	"github.com/ghodss/yaml"

	log "github.com/sirupsen/logrus"
)

type applicationExecutor struct {
	configs *model.Arguments
}

func StartConfiguring(config *model.Arguments) {
	executor := &applicationExecutor{config}
	executor.readConfigs()
}

func (a *applicationExecutor) readConfigs() {
	if len(a.configs.Configs) == 0 {
		log.Infof("No config files provided. Nothing to do here. Exit...")
		return
	}

	var cfg *model.ConfigWrapper
	var err error
	for _, filePath := range a.configs.Configs {
		cfg, err = loadConfig(cfg, filePath)
		if err != nil {
			log.WithError(err).Errorf("Error during validate %s", filePath)
			return
		}
	}
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
