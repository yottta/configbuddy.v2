package utils

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func ExecuteCommand(command string) error {
	cmdArr := strings.Fields(command)
	if len(cmdArr) < 2 {
		return fmt.Errorf("invalid command, arguments missing. command: '%s'", command)
	}

	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	log.WithField("command", command).Debug("prepare to run command")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.WithField("stdout/stderr", fmt.Sprintf("%s", stdoutStderr)).WithField("command", command).Error("command execution error")
		return err
	}
	return nil
}

func RemoveResource(resourcePath string) error {
	log.WithField("resource path", resourcePath).Debug("remove resource")
	return ExecuteCommand(fmt.Sprintf("rm -Rf %s", resourcePath))
}
