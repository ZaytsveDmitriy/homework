package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	// Place your code here.
	t.Run("limit < 0", func(t *testing.T) {
		err := Copy(".//testdata//input.txt", "out.txt", 0, -1)
		require.ErrorIs(t, err, ErrLimitIncorrect)
	})
	t.Run("offset < 0", func(t *testing.T) {
		err := Copy(".//testdata//input.txt", "out.txt", -1, 0)
		require.ErrorIs(t, err, ErrOffsetIncorrect)
	})
	t.Run("no dist path", func(t *testing.T) {
		err := Copy(".//testdata//input.txt", "", 0, 0)
		require.ErrorIs(t, err, ErrDistFilePathEmpty)
	})
	t.Run("no source path", func(t *testing.T) {
		err := Copy("", "out", 0, 0)
		require.ErrorIs(t, err, ErrSourceFilePath)
	})
	t.Run("no source file exist", func(t *testing.T) {
		err := Copy("no file exist", "out", 0, 0)
		require.ErrorIs(t, err, ErrNoSourceFileExist)
	})
	t.Run("file offset > file size", func(t *testing.T) {
		err := Copy("./testdata/input.txt", "out", 20000, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
	t.Run("file is dir", func(t *testing.T) {
		err := Copy("./testdata/dir", "out", 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})
	t.Run("positive case", func(t *testing.T) {
		err := Copy("./testdata/input.txt", "out.txt", 0, 0)
		require.NoError(t, err)
		os.Remove("out.txt")
	})
}
