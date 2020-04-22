package rcontext

import (
	"encoding/json"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type FindManager struct {
	ctxFile string
	file    stream.FileReadExister
}

func NewFinder(homePath string, file stream.FileReadExister) FindManager {
	return FindManager{
		ctxFile: fmt.Sprintf(ContextPath, homePath),
		file:    file,
	}
}

func (f FindManager) Find() (ContextHolder, error) {
	ctxHolder := ContextHolder{}

	if !f.file.Exists(f.ctxFile) {
		return ctxHolder, nil
	}

	read, err := f.file.Read(f.ctxFile)
	if err != nil {
		return ctxHolder, err
	}

	if err := json.Unmarshal(read, &ctxHolder); err != nil {
		return ctxHolder, err
	}

	return ctxHolder, nil
}
