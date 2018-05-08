package executor

import (
	"testing"

	"github.com/andreic92/configbuddy.v2/model"

	ast "github.com/stretchr/testify/assert"
)

const (
	directoryWithTestingFiles = "../testing"
)

func Test01TestReadFile(t *testing.T) {
	assert := ast.New(t)
	// this should be an error
	res, err := readFile("test.yml")
	assert.Error(err)
	assert.Nil(res)

	// this should be ok
	file := directoryWithTestingFiles + "/test1.yml"
	res, err = readFile(file)
	assert.NoError(err)
	assert.NotNil(res)
	assert.Equal(1, len(res.Config.FileActions))
	assert.Equal(1, len(res.Config.Includes))
	assert.Equal("test_1_included.yml", res.Config.Includes[0])

	invalidFileContent := directoryWithTestingFiles + "/test_invalid_file.yml"
	res, err = readFile(invalidFileContent)
	assert.Error(err)
	assert.Contains(err.Error(), "error unmarshaling JSON")
	assert.Nil(res)
}

func TestAppendConfigToAnother(t *testing.T) {
	assert := ast.New(t)

	// this should be ok
	file1 := directoryWithTestingFiles + "/test1.yml"
	file2 := directoryWithTestingFiles + "/test_1_included.yml"
	conf1, err := readFile(file1)
	assert.NoError(err)
	assert.NotNil(conf1)
	assert.Equal(1, len(conf1.Config.FileActions))
	assert.Equal(1, len(conf1.Config.Includes))
	assert.Equal("test_1_included.yml", conf1.Config.Includes[0])

	conf2, err := readFile(file2)
	assert.NoError(err)
	assert.NotNil(conf2)
	assert.Equal(1, len(conf2.Config.FileActions))
	assert.Equal(0, len(conf2.Config.Includes))

	err = appendActionsToGlobalConfig(conf2, conf1)
	assert.NoError(err)
	assert.Equal(2, len(conf1.Config.FileActions))

	conf1.Config.FileActions = nil
	err = appendActionsToGlobalConfig(conf2, conf1)
	assert.NoError(err)
	assert.Equal(1, len(conf1.Config.FileActions))
}

func TestLoadConfig(t *testing.T) {
	assert := ast.New(t)

	var conf *model.ConfigWrapper
	file1 := directoryWithTestingFiles + "/test1.yml"
	conf, err := loadConfig(conf, file1)
	assert.NoError(err)
	assert.NotNil(conf)
	assert.Equal(2, len(conf.Config.FileActions))
	assert.Equal(1, len(conf.Config.PackageActions))

	conf = nil
	invalidFile := directoryWithTestingFiles + "/test_invalid_file.yml"
	conf, err = loadConfig(conf, invalidFile)
	assert.Error(err)
	assert.Contains(err.Error(), "error unmarshaling JSON")
	assert.Nil(conf)

	conf = nil
	notExistingImportFile := directoryWithTestingFiles + "/test_not_existing_import.yml"
	conf, err = loadConfig(conf, notExistingImportFile)
	assert.Error(err)
	assert.Contains(err.Error(), "no such file or directory")
	assert.Nil(conf)
}

func TestReadConfigs(t *testing.T) {
	assert := ast.New(t)

	// happy scenario
	file1 := directoryWithTestingFiles + "/test1.yml"
	config := &model.Arguments{
		Configs: []string{file1},
	}
	executor := &applicationExecutor{configs: config}
	err := executor.readConfigs()
	assert.NoError(err)
	assert.NotNil(executor.finalConf)
	assert.Equal(2, len(executor.finalConf.Config.FileActions))

	// no configs given - return without any notice
	config = &model.Arguments{
		Configs: []string{},
	}
	executor = &applicationExecutor{configs: config}
	err = executor.readConfigs()
	assert.NoError(err)
	assert.Nil(executor.finalConf)

	// invalid content
	invalidFile := directoryWithTestingFiles + "/test_invalid_file.yml"
	config = &model.Arguments{
		Configs: []string{invalidFile},
	}
	executor = &applicationExecutor{configs: config}
	err = executor.readConfigs()
	assert.Error(err)
	assert.Contains(err.Error(), "error unmarshaling JSON")
	assert.Nil(executor.finalConf)

	// import valid file with a non existent imported file
	notExistingImportFile := directoryWithTestingFiles + "/test_not_existing_import.yml"
	config = &model.Arguments{
		Configs: []string{notExistingImportFile},
	}
	executor = &applicationExecutor{configs: config}
	err = executor.readConfigs()
	assert.Error(err)
	assert.Contains(err.Error(), "no such file or directory")
	assert.Nil(executor.finalConf)
}
