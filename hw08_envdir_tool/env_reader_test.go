package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestReadDir(t *testing.T) {
	// Place your code here

	const ExisstPath = "./testdata/env"
	const NotExistPath = "./testdata/not exist"
	const EmptyPath = "./testdata/emptydir/emptysubdir"
	const PathWithoutFile = "./testdata/emptydir"

	err := os.Mkdir(PathWithoutFile, os.ModeDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.Mkdir(EmptyPath, os.ModeDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	t.Run("not exist path", func(t *testing.T) {
		res, err := ReadDir(NotExistPath)
		require.ErrorIs(t, err, ErrDirCantOpen)
		require.Nil(t, res)
	})
	t.Run("dir empty", func(t *testing.T) {
		res, err := ReadDir(EmptyPath)
		require.ErrorIs(t, err, ErrDirEmpty)
		require.Nil(t, res)
	})
	t.Run("dir without file", func(t *testing.T) {
		res, err := ReadDir(PathWithoutFile)
		require.ErrorIs(t, err, ErrDirEmpty)
		require.Nil(t, res)
	})
	t.Run("cnt of envs", func(t *testing.T) {
		expectation := Environment{
			"BAR":   EnvValue{"bar", false},
			"EMPTY": EnvValue{"", false},
			"FOO": EnvValue{`   foo
with new line`, false},
			"HELLO": EnvValue{`"hello"`, false},
			"UNSET": EnvValue{"", true},
		}
		res, err := ReadDir(ExisstPath)
		require.NoError(t, err)
		require.Equal(t, res, expectation)
	})

	// clear testenv
	os.Remove(EmptyPath)
	os.Remove(PathWithoutFile)
}
