package repo

import (
	"encoding/json"
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	reposDirName  = "repos"
	reposFileName = "repositories.json"
)

var _ formula.RepositoryWriter = Writer{}

type Writer struct {
	reposFilePath string
	file          stream.FileWriter
}

func NewWriter(ritHome string, file stream.FileWriter) Writer {
	reposFilePath := filepath.Join(ritHome, reposDirName, reposFileName)
	return Writer{reposFilePath: reposFilePath, file: file}
}

func (w Writer) Write(repos formula.Repos) error {
	bytes, err := json.MarshalIndent(repos, "", "\t")
	if err != nil {
		return err
	}

	if err := w.file.Write(w.reposFilePath, bytes); err != nil {
		return err
	}

	return nil
}
