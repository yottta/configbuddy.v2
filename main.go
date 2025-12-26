package main

import (
	log "log/slog"
	"os"

	cli "github.com/jawher/mow.cli"

	"github.com/yottta/configbuddy.v2/backup"
	"github.com/yottta/configbuddy.v2/executor"
	"github.com/yottta/configbuddy.v2/model"
	"github.com/yottta/configbuddy.v2/parser"
	"github.com/yottta/go-core/logging"
)

const (
	appSystemCode  = "configbuddy"
	appDescription = "App done for installing my dotfiles"
)

func main() {
	logging.Setup()
	app := initApp()
	err := app.Run(os.Args)
	if err != nil {
		log.With("error", err).Error("app could not start")
		os.Exit(1)
	}
}

func initApp() *cli.Cli {
	app := cli.App(appSystemCode, appDescription)

	configs := app.Strings(cli.StringsOpt{
		Name:  "configs c",
		Value: []string{},
		Desc:  "The path for config files",
	})

	backupActivated := app.Bool(cli.BoolOpt{
		Name:  "backup b",
		Value: false,
		Desc:  "Boolean saying if the backup should be performed or not. If you want backup to directory, specify -p too",
	})
	backupDirectory := app.String(cli.StringOpt{
		Name:  "backup-path p",
		Value: "",
		Desc:  "Path of the folder where the backup will be performed",
	})

	app.Action = func() {
		log.Info("configbuddy started")

		args := &model.Arguments{
			Configs:         *configs,
			BackupDirectory: *backupDirectory,
			BackupActivated: *backupActivated,
		}

		backupService, err := backup.NewBackupService(args)
		if err != nil {
			log.With("error", err).Error("could not create the backup instance")
			return
		}
		parse, err := parser.NewParser()
		if err != nil {
			log.With("error", err).Error("could not create the parser instance")
			return
		}

		err = executor.StartConfiguring(args, parse, backupService)
		if err != nil {
			log.With("error", err).Error("error during configuration process")
			return
		}
	}

	return app
}
