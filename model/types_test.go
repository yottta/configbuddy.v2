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
		PackageName: "test_pkg_act",
		Alternatives: map[string][]string{
			"ubuntu": []string{
				"test_alternative",
			},
		},
	}
	assert.Equal("{ (model.PackageAction) Name: test_pkg_act; Alternatives: map[ubuntu:[test_alternative]]; }", fmt.Sprintf("%s", pkgAct))

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
		PackageActions: []PackageAction{
			*pkgAct,
		},
		Includes: []string{"first_import", "second_import"},
		Globals:  globals,
	}
	assert.Equal("{ (model.Config) Globals: { (model.Globals) ExitOnError: true; ConfirmEveryPackage: true }; Includes: [first_import second_import]; FileActions: map[test_file_act:{ (model.FileAction) FileName: test_filename; Hidden: true; Source: test_source; Command: test_command; Destination: test_destination }]; PackageActions: [{ (model.PackageAction) Name: test_pkg_act; Alternatives: map[ubuntu:[test_alternative]]; }] }", fmt.Sprintf("%s", config))

	wrap := ConfigWrapper{
		Config:              config,
		ConfigFileDirectory: "config_dir",
		ConfigFilePath:      "config_file_path",
	}
	assert.Equal("{ (model.ConfigWrapper) ConfigFileDirectory: config_dir; ConfigFilePath: config_file_path; Config: { (model.Config) Globals: { (model.Globals) ExitOnError: true; ConfirmEveryPackage: true }; Includes: [first_import second_import]; FileActions: map[test_file_act:{ (model.FileAction) FileName: test_filename; Hidden: true; Source: test_source; Command: test_command; Destination: test_destination }]; PackageActions: [{ (model.PackageAction) Name: test_pkg_act; Alternatives: map[ubuntu:[test_alternative]]; }] } }", fmt.Sprintf("%s", wrap))
}
