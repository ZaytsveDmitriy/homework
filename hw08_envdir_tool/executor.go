package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	ex := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	ex.Stderr = os.Stderr
	ex.Stdin = os.Stdin
	ex.Stdout = os.Stdout

	for name, e := range env {
		if e.NeedRemove {
			os.Unsetenv(name)
			continue
		}
		if _, hasEnv := os.LookupEnv(name); hasEnv {
			os.Unsetenv(name)
			os.Setenv(name, e.Value)
		} else {
			os.Setenv(name, e.Value)
		}
	}
	err := ex.Run()
	if err != nil {
		return -1
	}
	return 0
}
