package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// CreateFileWithParent creates a file and all prefixed directories first
func CreateFileWithParent(file string) error {
	if file == "" {
		return errors.New("assert: blank values are not allowed for 'file'")
	}

	var err error
	parentDir := filepath.Dir(file)

	if err = os.MkdirAll(parentDir, os.ModePerm); err != nil {
		return errors.New(fmt.Sprintf("cannot create parent directory '%v': %v", parentDir, fmt.Errorf("%w", err)))
	}

	if _, err = os.Stat(file); errors.Is(err, os.ErrNotExist) {
		var f *os.File
		f, err = os.Create(file)
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
		if err != nil {
			return errors.New(fmt.Sprintf("cannot create file '%v': %v", file, fmt.Errorf("%w", err)))
		}
	}

	return err
}

// CreateDirectoryRecursively creates a directory recursively
func CreateDirectoryRecursively(dir string) error {
	if dir == "" {
		return errors.New("assert: blank values are not allowed for 'dir'")
	}

	var err error
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return errors.New(fmt.Sprintf("cannot create parent directory '%v': %v", dir, fmt.Errorf("%w", err)))
	}

	return err
}
