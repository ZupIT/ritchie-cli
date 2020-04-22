package stream

import (
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

type FileRemover interface {
	Remove(path string) error
}

type FileReadExister interface {
	FileReader
	FileExister
}

type FileWriteReadExister interface {
	FileWriter
	FileReader
	FileExister
}

type FileWriteReadExistRemover interface {
	FileWriter
	FileReader
	FileExister
	FileRemover
}

// ReadFile implements FileReader
type ReadFile struct {
}

// NewFileReader returns a ReadFile
func NewFileReader() ReadFile {
	return ReadFile{}
}

// Read reads the file named by path and returns the contents.
// A successful call returns err == nil
func (p ReadFile) Read(path string) ([]byte, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return f, err
}

// ExistFile implements FileExister
type ExistFile struct {
}

// NewFileExister returns an ExistFile
func NewFileExister() ExistFile {
	return ExistFile{}
}

// Exists returns true if file path exists
func (p ExistFile) Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// WriteFile implements FileWriter
type WriteFile struct {
}

// NewFileWriter returns a WriteFile
func NewFileWriter() WriteFile {
	return WriteFile{}
}

// Write writes content to a file named by path.
// A successful call returns err == nil
func (w WriteFile) Write(path string, content []byte) error {
	return ioutil.WriteFile(path, content, os.ModePerm)
}

// RemoveFile implements FileRemover
type RemoveFile struct {
	FileExister
}

// NewFileRemover returns a RemoveFile
func NewFileRemover(e FileExister) RemoveFile {
	return RemoveFile{e}
}

// Remove removes the named file
func (r RemoveFile) Remove(path string) error {
	if r.Exists(path) {
		return os.Remove(path)
	}
	return nil
}

// FileManager implements FileWriteReadExistRemover
type FileManager struct {
	FileWriter
	FileReader
	FileExister
	FileRemover
}

// NewFileManager returns a FileManager that writes from w
// reads from r, exists from e and removes from re
func NewFileManager(w FileWriter, r FileReader, e FileExister, re FileRemover) FileManager {
	return FileManager{
		FileWriter:  w,
		FileReader:  r,
		FileExister: e,
		FileRemover: re,
	}
}

type ReadExistFile struct {
	ReadFile
	ExistFile
}

func NewReadExister(r ReadFile, e ExistFile) ReadExistFile {
	return ReadExistFile{r, e}
}
