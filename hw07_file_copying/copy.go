package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOffsetIncorrect       = errors.New("offset provide incorrectry, offset < 0")
	ErrLimitIncorrect        = errors.New("limit provide incorrectly, limit <= 0")
	ErrNoSourceFileExist     = errors.New("source file not exist")
	ErrSourceFilePath        = errors.New("source file path not provide or provide icorrectly")
	ErrDistFilePathEmpty     = errors.New("distination file path not provide or provide incorrectly")
	ErrUnknown               = errors.New("unexpected error")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	switch {
	case fromPath == "":
		return fmt.Errorf("%w: file path - \"%s\"", ErrSourceFilePath, fromPath)
	case toPath == "":
		return fmt.Errorf("%w: file path - \"%s\"", ErrDistFilePathEmpty, toPath)
	case offset < 0:
		return ErrOffsetIncorrect
	case limit < 0:
		return ErrLimitIncorrect
	}

	file, err := os.Open(fromPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoSourceFileExist
		}
		return fmt.Errorf("%w: %w", ErrUnknown, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnknown, err)
	}

	switch {
	case info.IsDir():
		return fmt.Errorf("%w: file is directory", ErrUnsupportedFile)
	case info.Size() <= 0:
		return fmt.Errorf("%w: file size <= 0", ErrUnsupportedFile)
	case info.Size() < offset:
		return fmt.Errorf("%w: offset %v, file size %v", ErrOffsetExceedsFileSize, offset, info.Size())
	}

	if limit == 0 {
		limit = info.Size()
	} else {
		limit = min(info.Size(), limit)
	}

	distFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnknown, err)
	}

	stat := pb.Full.Start64(limit)
	proxyReader := stat.NewProxyReader(file)
	io.CopyN(distFile, proxyReader, limit)
	stat.Finish()

	return nil
}
