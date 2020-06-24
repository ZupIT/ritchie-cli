package stream

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type DirCreater interface {
	Create(dir string) error
}

type DirRemover interface {
	Remove(dir string) error
}

type DirLister interface {
	List(dir string, hiddenDir bool) ([]string, error)
}

type DirChecker interface {
	Exists(dir string) bool
	IsDir(dir string) bool
}

type DirCopier interface {
	Copy(src, dst string) error
}

type DirCreateListCopier interface {
	DirCreater
	DirLister
	DirCopier
}

type DirListChecker interface {
	DirLister
	DirChecker
}

type DirManager struct {
	file FileCopier
}

func NewDirManager(file FileCopier) DirManager {
	return DirManager{file: file}
}

// Create creates a directory named dir
// A successful call returns err == nil
func (m DirManager) Create(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

// Remove removes dir and any children it contains.
func (m DirManager) Remove(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return nil
}

func (m DirManager) Exists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}

	return true
}

func (m DirManager) IsDir(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return info.IsDir()
}

// List list all directories associate with dir name, case hiddenDir is
// true, we will add hidden directories to your slice of directories.
// If success, List returns a slice of directories and a nil error.
// If error, List returns an empty slice and a non-nil error.
func (m DirManager) List(dir string, hiddenDir bool) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, f := range fileInfos {
		n := f.Name()
		if f.IsDir() {
			if hiddenDir {
				dirs = append(dirs, n)
				continue
			}

			if !strings.ContainsAny(n, ".") {
				dirs = append(dirs, n)
			}
		}
	}

	return dirs, nil
}

func (m DirManager) Copy(src, dest string) error {
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := m.Create(destPath); err != nil {
				return err
			}
			if err := m.Copy(sourcePath, destPath); err != nil {
				return err
			}
		} else {
			if err := m.file.Copy(sourcePath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}
