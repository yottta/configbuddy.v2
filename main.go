package main

import (
	"os"

	"github.com/andreic92/configbuddy.v2/executor"

	"github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
)

const (
	appSystemCode  = "configbuddy"
	appDescription = "Lightweight configuration system made for simple tasks"
)

func main() {
	app := initApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("App could not start, error=[%s]\n", err)
		return
	}
}

func initApp() *cli.Cli {
	app := cli.App(appSystemCode, appDescription)

	configs := app.StringsOpt("c", []string{}, "The path for config files")

	initLogging()

	app.Action = func() {
		log.Infof("Configbuddy started")

		executor.StartConfiguring(*configs)
	}

	return app
}

func initLogging() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}
