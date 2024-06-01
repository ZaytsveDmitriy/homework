package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	ex := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	ex.Stderr = os.Stderr
	ex.Stdin = os.Stdin
	ex.Stdout = os.Stdout
	ex.Env = ex.Environ()
	for name, e := range env {
		if e.NeedRemove {
			ex.Env = append(ex.Env, name+"=")
		} else {
			os.Setenv(name, e.Value)
			ex.Env = append(ex.Env, fmt.Sprint(name, "=", e.Value))
		}
	}
	ex.Run()

	return ex.ProcessState.ExitCode()
}
