package model

import "fmt"

type ConfigWrapper struct {
	Config              *Config
	ConfigFilePath      string
	ConfigFileDirectory string
}
type Config struct {
	Globals        *Globals                 `json:"Globals"`
	Includes       []string                 `json:"includes"`
	FileActions    map[string]FileAction    `json:"FileAction"`
	PackageActions map[string]PackageAction `json:"PackageAction"`
}

type Globals struct {
	ExitOnError         bool `json:"exitOnError"`
	ConfirmEveryPackage bool `json:"confirmEveryPackage"`
}

type FileAction struct {
	FileName    string `json:"name"` // if empty the map key will be used
	Hidden      bool   `json:"hidden"`
	Source      string `json:"source"`
	Command     string `json:"command"`
	Destination string `json:"destination"`
}
type PackageAction struct {
	Alternatives string `json:"alternatives"`
}

func (c Config) String() string {
	return fmt.Sprintf("{ (model.Config) Globals: %s; Includes: %s; FileActions: %s; PackageActions: %s }",
		c.Globals,
		c.Includes,
		c.FileActions,
		c.PackageActions)
}

func (g Globals) String() string {
	return fmt.Sprintf("{ (model.Globals) ExitOnError: %t; ConfirmEveryPackage: %t }",
		g.ExitOnError,
		g.ConfirmEveryPackage)
}

func (f FileAction) String() string {
	return fmt.Sprintf("{ (model.FileAction) FileName: %s; Hidden: %t; Source: %s; Command: %s; Destination: %s }",
		f.FileName,
		f.Hidden,
		f.Source,
		f.Command,
		f.Destination)
}

func (p PackageAction) String() string {
	return fmt.Sprintf("{ (model.PackageAction) Alternatives: %s; }",
		p.Alternatives)
}
