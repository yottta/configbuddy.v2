package model

import (
	"fmt"
	"testing"

	ast "github.com/stretchr/testify/assert"
)

// test required for being sure about all fields are included in String()
func Test01TestStringFunc(t *testing.T) {
	assert := ast.New(t)

	pkgAct := &PackageAction{
		Alternatives: "test_alternatives",
	}
	assert.Equal("{ (model.PackageAction) Alternatives: test_alternatives; }", fmt.Sprintf("%s", pkgAct))

	fileAct := &FileAction{
		Command:     "test_command",
		Destination: "test_destination",
		FileName:    "test_filename",
		Hidden:      true,
		Source:      "test_source",
	}
	assert.Equal("{ (model.FileAction) FileName: test_filename; Hidden: true; Source: test_source; Command: test_command; Destination: test_destination }", fmt.Sprintf("%s", fileAct))

	globals := &Globals{
		ConfirmEveryPackage: true,
		ExitOnError:         true,
	}
	assert.Equal("{ (model.Globals) ExitOnError: true; ConfirmEveryPackage: true }", fmt.Sprintf("%s", globals))

	config := &Config{
		FileActions: map[string]FileAction{
			"test_file_act": *fileAct,
		},
		PackageActions: map[string]PackageAction{
			"test_pkg_act": *pkgAct,
		},
		Includes: []string{"first_import", "second_import"},
		Globals:  globals,
	}
	assert.Equal("{ (model.Config) Globals: { (model.Globals) ExitOnError: true; ConfirmEveryPackage: true }; Includes: [first_import second_import]; FileActions: map[test_file_act:{ (model.FileAction) FileName: test_filename; Hidden: true; Source: test_source; Command: test_command; Destination: test_destination }]; PackageActions: map[test_pkg_act:{ (model.PackageAction) Alternatives: test_alternatives; }] }", fmt.Sprintf("%s", config))

	wrap := ConfigWrapper{
		Config:              config,
		ConfigFileDirectory: "config_dir",
		ConfigFilePath:      "config_file_path",
	}
	assert.Equal("{ (model.ConfigWrapper) ConfigFileDirectory: config_dir; ConfigFilePath: config_file_path; Config: { (model.Config) Globals: { (model.Globals) ExitOnError: true; ConfirmEveryPackage: true }; Includes: [first_import second_import]; FileActions: map[test_file_act:{ (model.FileAction) FileName: test_filename; Hidden: true; Source: test_source; Command: test_command; Destination: test_destination }]; PackageActions: map[test_pkg_act:{ (model.PackageAction) Alternatives: test_alternatives; }] } }", fmt.Sprintf("%s", wrap))
}
