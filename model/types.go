package model

import "fmt"

type ConfigWrapper struct {
	Config              *Config
	ConfigFilePath      string
	ConfigFileDirectory string
}
type Config struct {
	Globals        *Globals              `json:"Globals"`
	Includes       []string              `json:"includes"`
	FileActions    map[string]FileAction `json:"FileAction"`
	PackageActions []PackageAction       `json:"PackageAction"`
}

type Globals struct {
	ExitOnError         bool `json:"exitOnError"`
	ConfirmEveryPackage bool `json:"confirmEveryPackage"`
}

type ConditionalAction struct {
	When string `json:"when"`
}

type FileAction struct {
	ConditionalAction
	FileName    string `json:"name"` // if empty the map key will be used
	Hidden      bool   `json:"hidden"`
	Source      string `json:"source"`
	Command     string `json:"command"`
	Destination string `json:"destination"`
}
type PackageAction struct {
	ConditionalAction
	PackageName  string              `json:"name"`         // if empty the map key will be used
	Alternatives map[string][]string `json:"alternatives"` // map distro name with the the package alternative(s) for that specific distro
	Source       string              `json:"source"`
	URL          string              `json:"url"`
	Destination  string              `json:"destination"`
	Sudo         bool                `json:"sudo"`
}

type Arguments struct {
	Configs         []string
	BackupDirectory string
	BackupActivated bool
}

func (c *ConditionalAction) Condition() string {
	return c.When
}
func (w ConfigWrapper) String() string {
	return fmt.Sprintf("{ (model.ConfigWrapper) ConfigFileDirectory: %s; ConfigFilePath: %s; Config: %s }",
		w.ConfigFileDirectory,
		w.ConfigFilePath,
		w.Config)
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
	return fmt.Sprintf("{ (model.PackageAction) Name: %s; Alternatives: %s; }",
		p.PackageName, p.Alternatives)
}
