package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestRunCmd(t *testing.T) {
	// Place your code here
	const ExpectedOut = `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2
`

	envs := Environment{
		"BAR":   EnvValue{"bar", false},
		"EMPTY": EnvValue{"", false},
		"FOO": EnvValue{`   foo
with new line`, false},
		"HELLO": EnvValue{`"hello"`, false},
		"UNSET": EnvValue{"", true},
	}

	cmds := []string{"bin/bash", "./testdata/echo.sh", "arg1=1", "arg2=2"}

	t.Run("not exist path", func(t *testing.T) {
		stdOut := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Println(err)
		}

		os.Stdout = w
		os.Setenv("ADDED", "from original env")

		res := RunCmd(cmds, envs)

		w.Close()
		os.Stdout = stdOut

		var buf bytes.Buffer
		io.Copy(&buf, r)

		require.Zero(t, res)
		require.Equal(t, ExpectedOut, buf.String())
	})
}
