package parser

import (
	"fmt"
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
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

func GetPckManager() (string, error) {
	if runtime.GOOS == DarwinSystem {
		return BrewPackageManager, nil
	}

	if runtime.GOOS == LinuxSystem {
		return getLinuxPckManager()
	}

	return "", fmt.Errorf("The current OS (%s) is not supported", runtime.GOOS)
}

func getLinuxPckManager() (string, error) {
	for fileToCheck, packageManager := range linuxSupportedSystems {
		stat, err := os.Stat(fileToCheck)
		if err != nil {
			log.WithError(err).
				WithField("FileToCheck", fileToCheck).
				WithField("PackageManager", packageManager).
				Debug("Skipped package manager check")
			continue
		}
		if !stat.IsDir() {
			return packageManager, nil
		}
	}
	return "", ErrorNotSupportedDistro
}
