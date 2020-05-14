package fileutil

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// CopyDirectory recursively copies a src directory to a destination.
func CopyDirectory(src, dst string) error {
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dst, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := CreateDirIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := Copy(sourcePath, destPath); err != nil {
				return err
			}
		}

		// `go test` fails on Windows even with this `if` supposedly
		// protecting the `syscall.Stat_t` and `os.Lchown` calls (not
		// available on windows). why?
		/*
			if runtime.GOOS != "windows" {
					stat, ok := fileInfo.Sys().(*syscall.Stat_t)
					if !ok {
						return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
					}
					if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
						return err
					}
			}
		*/

		isSymlink := entry.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, entry.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

// Copy copies a src file to a dst file where src and dst are regular files.
func Copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

// Exists check if file exists
func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// CreateDirIfNotExists creates dir if not exists
func CreateDirIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

// CreateFileIfNotExist creates file if not exists
func CreateFileIfNotExist(file string, content []byte) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		if err := WriteFile(file, content); err != nil {
			return err
		}
	}
	return nil
}

// RemoveDir removes path and any children it contains.
func RemoveDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return nil
}

// CopySymLink copies src symLink to dst symLink
func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

// ReadFile wrapper for ioutil.ReadFile
func ReadFile(path string) ([]byte, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return f, err
}

// ReadAll reads from r until an error or EOF and returns the data it read
func ReadAll(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}

// CopyAll copies from src to dst until either EOF is reached on src or an error occurs
func CopyAll(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

// WriteFile wrapper for ioutil.WriteFile
func WriteFile(path string, content []byte) error {
	return ioutil.WriteFile(path, content, 0644)
}

// WriteFilePerm wrapper for ioutil.WriteFile
func WriteFilePerm(path string, content []byte, perm int32) error {
	return ioutil.WriteFile(path, content, os.FileMode(perm))
}

// RemoveFile wrapper for os.Delete
func RemoveFile(path string) error {
	if exits := Exists(path); exits {
		return os.Remove(path)
	}
	return nil
}

// Unzip wrapper for archive.zip
func Unzip(src string, dest string) error {
	reader, _ := zip.OpenReader(src)
	for _, file := range reader.Reader.File {

		zippedFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zippedFile.Close()

		extractedFilePath := filepath.Join(
			dest,
			file.Name,
		)

		if file.FileInfo().IsDir() {
			err := os.MkdirAll(extractedFilePath, file.Mode())
			if err != nil {
				return err
			}
		} else {
			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				return err
			}
			defer outputFile.Close()

			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				return err
			}
		}
	}
	defer reader.Close()
	return nil
}

// IsNotExistErr
func IsNotExistErr(err error) bool {
	return os.IsNotExist(errors.Cause(err))
}

// Read all files and dir in current directory
func readFilesDir(path string) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fl, err := f.Readdir(-1)
	f.Close()
	return fl, err
}

// Move files from oPath to nPath
func MoveFiles(oPath, nPath string, files []string) error {
	for _, f := range files {
		pwdOF := fmt.Sprintf("%s/%s", oPath, f)
		pwdNF := fmt.Sprintf("%s/%s", nPath, f)
		if err := os.Rename(pwdOF, pwdNF); err != nil {
			return err
		}
	}
	return nil
}

// List new files in nPath differing of oPath
func ListNewFiles(oPath, nPath string) ([]string, error) {
	of, err := readFilesDir(oPath)
	if err != nil {
		return nil, err
	}
	nf, err := readFilesDir(nPath)
	if err != nil {
		return nil, err
	}
	control := make(map[string]int)
	for _, file := range of {
		control[file.Name()]++
	}
	var new []string
	for _, file := range nf {
		if control[file.Name()] == 0 {
			new = append(new, file.Name())
		}
	}
	return new, nil
}
