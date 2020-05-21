package stream

import (
	"fmt"
	"os"
)

type DirCreater interface {
	Create(dir string) error
}

type DirRemover interface {
	Remove(dir string) error
}

// CreateDir implements DirCreater
type CreateDir struct {
}

// NewDirCreater returns a CreateDir
func NewDirCreater() CreateDir {
	return CreateDir{}
}

// Create creates a directory named dir
// A successful call returns err == nil
func (c CreateDir) Create(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

// RemoveDir implements DirRemover
type RemoveDir struct {
}

// NewDirRemover returns a RemoveDir
func NewDirRemover() RemoveDir {
	return RemoveDir{}
}

// Remove removes dir and any children it contains.
func (r RemoveDir) Remove(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return nil
}
