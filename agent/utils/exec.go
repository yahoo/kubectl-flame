package utils

import "os/exec"

func ExecuteCommand(cmd *exec.Cmd) (int, string, error) {
	exitCode := 0
	output, err := cmd.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}

	return exitCode, string(output), err
}
