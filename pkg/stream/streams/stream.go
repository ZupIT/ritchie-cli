package streams

import (
	"archive/zip"
	"io"
	"log"
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
		defer zippedFile.Close()

		extractedFilePath := filepath.Join(
			dest,
			file.Name,
		)

		if file.FileInfo().IsDir() {
			log.Println("Directory Created:", extractedFilePath)
			if err := os.MkdirAll(extractedFilePath, file.Mode()); err != nil {
				return err
			}
		} else {
			log.Println("File extracted:", file.Name)

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
