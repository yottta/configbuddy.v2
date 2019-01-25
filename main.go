package main

import (
	"fmt"
	"os"
	"strings"

	cli "github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"

	"github.com/andreic92/configbuddy.v2/backup"
	"github.com/andreic92/configbuddy.v2/executor"
	"github.com/andreic92/configbuddy.v2/model"
	"github.com/andreic92/configbuddy.v2/parser"
)

const (
	appSystemCode  = "configbuddy"
	appDescription = "App done for installing my dotfiles"
)

func main() {
	app := initApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("App could not start, error=[%s]\n", err)
		os.Exit(1)
	}
}

func initApp() *cli.Cli {
	app := cli.App(appSystemCode, appDescription)

	configs := app.StringsOpt("c", []string{}, "The path for config files")
	backupActivated := app.BoolOpt("b", false, "Boolean saying if the backup should be performed or not. If you want backup to directory, specify -p too")
	backupDirectory := app.StringOpt("p", "", "Path of the folder where the backup will be performed")
	loggingLevel := app.StringOpt("l", "info", getLoggingFlagDescription())

	app.Action = func() {
		initLogging(*loggingLevel)
		log.Infof("Configbuddy started")

		args := &model.Arguments{
			Configs:         *configs,
			BackupDirectory: *backupDirectory,
			BackupActivated: *backupActivated,
		}

		backupService, err := backup.NewBackupService(args)
		if err != nil {
			log.WithError(err).Error("Could not create the backup instance")
			return
		}
		parser, err := parser.NewParser()
		if err != nil {
			log.WithError(err).Error("Could not create the parser instance")
			return
		}

		err = executor.StartConfiguring(args, parser, backupService)
		if err != nil {
			log.WithError(err).Error("Error during configuration process")
			return
		}
	}

	return app
}

func initLogging(loggingLevel string) {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logLevel, err := log.ParseLevel(strings.ToLower(loggingLevel))
	if err != nil {
		panic(err)
	}
	log.SetLevel(logLevel)
	log.Infof("Logging level set to %s", strings.ToLower(loggingLevel))
}

func getLoggingFlagDescription() string {
	var levelsAsString []string
	for _, lvl := range log.AllLevels {
		levelsAsString = append(levelsAsString, lvl.String())
	}
	return fmt.Sprintf("The logging level. Valid values: %s", strings.Join(levelsAsString, ","))
}
