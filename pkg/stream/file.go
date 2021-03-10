/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package stream

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

// FileNewLister list all new files compared between old and new path
type FileNewLister interface {
	ListNews(oldPath, newPath string) ([]string, error)
}

type FileCopier interface {
	Copy(src, dst string) error
}

type FileAppender interface {
	Append(path string, content []byte) error
}

type FileMover interface {
	Move(oldPath, newPath string, files []string) error
}

type FileReadWriter interface {
	FileReader
	FileWriter
}

type FileReadExister interface {
	FileReader
	FileExister
}

type FileNewListMoveRemover interface {
	FileMover
	FileRemover
	FileNewLister
}

type FileListMover interface {
	FileLister
	FileMover
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

type FileWriteReadExistLister interface {
	FileWriter
	FileReader
	FileExister
	FileLister
}

type FileWriteExistAppendRemover interface {
	FileWriter
	FileExister
	FileAppender
	FileRemover
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

type FileWriteRemover interface {
	FileWriter
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

func (f FileManager) Append(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(content); err != nil {
		return err
	}

	return nil
}

func (f FileManager) Move(oldPath, newPath string, files []string) error {
	for _, f := range files {
		pwdOF := filepath.Join(oldPath, f)
		pwdNF := filepath.Join(newPath, f)
		if err := os.Rename(pwdOF, pwdNF); err != nil {
			return err
		}
	}
	return nil
}

func (f FileManager) ListNews(oldPath, newPath string) ([]string, error) {
	of, err := readFilesDir(oldPath)
	if err != nil {
		return nil, err
	}
	nf, err := readFilesDir(newPath)
	if err != nil {
		return nil, err
	}
	control := make(map[string]int)
	for _, file := range of {
		control[file.Name()]++
	}
	var ff []string
	for _, file := range nf {
		if control[file.Name()] == 0 {
			ff = append(ff, file.Name())
		}
	}
	return ff, nil
}

func readFilesDir(path string) ([]os.FileInfo, error) {
	filesInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return filesInfo, err
}
