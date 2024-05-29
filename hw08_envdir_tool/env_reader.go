package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var (
	ErrDirCantOpen   = errors.New("can`t open dir")
	ErrDirEmpty      = errors.New("dir is empty")
	ErrOpenFile      = errors.New("can`t open file")
	ErrUnexpected    = errors.New("unexpected error")
	ErrForbiddenChar = errors.New("env contain forbidden char")
)

func normalize(b []byte) string {
	b = bytes.TrimRight(b, " \t")
	b = bytes.ReplaceAll(b, []byte{0x00}, []byte{'\n'})
	return string(b)
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDirCantOpen, err)
	}
	if len(files) < 1 {
		return nil, ErrDirEmpty
	}

	env := make(Environment)

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		filePath := path.Join(dir, f.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("%w: %s - %w", ErrOpenFile, f.Name(), err)
		}

		fInfo, err := f.Info()
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrUnexpected, err)
		}

		if fInfo.Size() == 0 {
			env[fInfo.Name()] = EnvValue{Value: "", NeedRemove: true}
			continue
		}

		reader := bufio.NewReaderSize(file, int(fInfo.Size()))
		envVal, _, err := reader.ReadLine()
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrUnexpected, err)
		}
		if strings.Contains(string(envVal), "=") {
			return nil, ErrForbiddenChar
		}

		env[fInfo.Name()] = EnvValue{Value: normalize(envVal), NeedRemove: false}
	}

	if len(env) == 0 {
		return nil, ErrDirEmpty
	}

	return env, nil
}
