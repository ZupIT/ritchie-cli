package stream

import (
	"io"
	"io/ioutil"
	"os"
)

type FileExister interface {
	Exists(path string) bool
}

type FileReader interface {
	Read(path string) ([]byte, error)
}

type FileWriter interface {
	Write(path string, content []byte) error
}

type FileCreator interface {
	Create(path string, data io.ReadCloser) error
}

type FileRemover interface {
	Remove(path string) error
}

type FileLister interface {
	List(file string) ([]string, error)
}

type FileCopier interface {
	Copy(src, dst string) error
}

type FileReadExister interface {
	FileReader
	FileExister
}

type FileCopyExistListerWriter interface {
	FileLister
	FileCopier
	FileExister
	FileWriter
}

type FileWriteReadExister interface {
	FileWriter
	FileReader
	FileExister
}

type FileWriteReadExisterLister interface {
	FileWriter
	FileReader
	FileExister
	FileLister
}

type FileWriteReadExistRemover interface {
	FileWriter
	FileReader
	FileExister
	FileRemover
}

type FileWriteCreatorReadExistRemover interface {
	FileWriter
	FileCreator
	FileReader
	FileExister
	FileRemover
}

// FileManager implements FileWriteReadExistRemover
type FileManager struct {
}

// NewFileManager returns a FileManage that writes from w
// reads from r, exists from e and removes from re
func NewFileManager() FileManager {
	return FileManager{}
}

// Read reads the file named by path and returns the contents.
// A successful call returns err == nil
func (f FileManager) Read(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return b, err
}

// Exists returns true if file path exists
func (f FileManager) Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// Write writes content to a file named by path.
// A successful call returns err == nil
func (f FileManager) Write(path string, content []byte) error {
	return ioutil.WriteFile(path, content, os.ModePerm)
}

func (f FileManager) Create(path string, data io.ReadCloser) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()
	if _, err = io.Copy(file, data); err != nil {
		return err
	}

	return nil
}

// Remove removes the named file
func (f FileManager) Remove(path string) error {
	if f.Exists(path) {
		return os.Remove(path)
	}
	return nil
}

// List lists all files in dir path
func (f FileManager) List(dir string) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, f := range fileInfos {
		n := f.Name()
		if !f.IsDir() {
			dirs = append(dirs, n)
		}
	}

	return dirs, nil
}

func (f FileManager) Copy(src, dest string) error {
	input, err := f.Read(src)
	if err != nil {
		return err
	}

	if err := f.Write(dest, input); err != nil {
		return err
	}

	return nil
}
