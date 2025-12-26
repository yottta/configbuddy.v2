package utils

import (
	"fmt"
	log "log/slog"
	"os/exec"
	"strings"
)

func ExecuteCommand(command string) error {
	cmdArr := strings.Fields(command)
	if len(cmdArr) < 2 {
		return fmt.Errorf("invalid command, arguments missing. command: '%s'", command)
	}

	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	log.With("command", command).Debug("prepare to run command")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.With("error", err).
			With("stdout/stderr", fmt.Sprintf("%s", stdoutStderr)).
			With("command", command).Error("command execution error")
		return err
	}
	return nil
}

func RemoveResource(resourcePath string) error {
	log.With("resource path", resourcePath).Debug("remove resource")
	return ExecuteCommand(fmt.Sprintf("rm -Rf %s", resourcePath))
}
