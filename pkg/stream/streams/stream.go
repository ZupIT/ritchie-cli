package streams

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// Unzip wrapper for archive.zip
func Unzip(src string, dest string) error {
	reader, _ := zip.OpenReader(src)
	for _, file := range reader.Reader.File {

		zippedFile, err := file.Open()
		if err != nil {
			return err
		}

		extractedFilePath := filepath.Join(
			dest,
			file.Name,
		)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(extractedFilePath, file.Mode()); err != nil {
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

			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				return err
			}
			outputFile.Close()
		}
		zippedFile.Close()
	}
	defer reader.Close()
	return nil
}
