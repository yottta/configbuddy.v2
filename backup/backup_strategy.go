package backup

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type backupStrategier interface {
	extractTargetPath(resourcePath string) (string, error)
}

func getStrategier(strategy backupStrategy, backupDirectory string) backupStrategier {
	switch strategy {
	case backupStrategyBakFile:
		return &copyToBakFileStrategy{}
	case backupStrategyCopyToDirectory:
		return &copyToDirectoryStrategy{backupDirectory: backupDirectory}
	default:
		return nil
	}
}

/* *********** COPY TO DIRECTORY STRATEGY *********** */
type copyToDirectoryStrategy struct {
	backupDirectory string
}

func (cd *copyToDirectoryStrategy) extractTargetPath(resourcePath string) (string, error) {
	if len(resourcePath) == 0 {
		return "", fmt.Errorf("invalid resource name - empty resource name")
	}
	return filepath.Abs(
		path.Join(
			cd.backupDirectory,
			getCurrentTimeFileNameStamp()+"_"+strings.Replace(resourcePath, string(filepath.Separator), "_", -1)),
	)
}

/* ************************************************** */

/* ************ COPY TO BAK FILE STRATEGY *********** */
type copyToBakFileStrategy struct {
}

func (cb *copyToBakFileStrategy) extractTargetPath(resourcePath string) (string, error) {
	if len(resourcePath) == 0 {
		return "", fmt.Errorf("invalid resource name - empty resource name")
	}
	return filepath.Abs(path.Join(
		path.Dir(resourcePath),
		getCurrentTimeFileNameStamp()+"_"+path.Base(resourcePath)),
	)
}

/* ************************************************** */

func getCurrentTimeFileNameStamp() string {
	now := time.Now()
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d.%06d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond())
}
