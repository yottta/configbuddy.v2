package parser

import (
	"fmt"
	log "log/slog"
	"os"
	"runtime"
)

const (
	DarwinSystem = `darwin`
	LinuxSystem  = `linux`

	BrewPackageManager = `brew`
)

var (
	linuxSupportedSystems = map[string]string{
		"/etc/redhat-release": "yum",
		"/etc/arch-release":   "pacman",
		"/etc/gentoo-release": "emerge",
		"/etc/SuSE-release":   "zypp",
		"/etc/debian_version": "apt-get",
	}
	ErrorNotSupportedDistro = fmt.Errorf("Unsupported distro")
)

func PckManager() (string, error) {
	if runtime.GOOS == DarwinSystem {
		return BrewPackageManager, nil
	}

	if runtime.GOOS == LinuxSystem {
		return linuxPckManager()
	}

	return "", fmt.Errorf("The current OS (%s) is not supported", runtime.GOOS)
}

func linuxPckManager() (string, error) {
	for fileToCheck, packageManager := range linuxSupportedSystems {
		stat, err := os.Stat(fileToCheck)
		if err != nil {
			log.With("error", err).
				With("file to check", fileToCheck).
				With("package manager", packageManager).
				Debug("skipped package manager check")
			continue
		}
		if !stat.IsDir() {
			return packageManager, nil
		}
	}
	return "", ErrorNotSupportedDistro
}
