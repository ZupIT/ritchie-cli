package rcontext

import (
	"encoding/json"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type FindManager struct {
	CtxFile string
	File    stream.FileReadExister
}

func NewFinder(homePath string, file stream.FileReadExister) FindManager {
	return FindManager{
		CtxFile: fmt.Sprintf(ContextPath, homePath),
		File:    file,
	}
}

func (f FindManager) Find() (ContextHolder, error) {
	ctxHolder := ContextHolder{}

	if !f.File.Exists(f.CtxFile) {
		return ctxHolder, nil
	}

	file, err := f.File.Read(f.CtxFile)
	if err != nil {
		return ctxHolder, err
	}

	if err := json.Unmarshal(file, &ctxHolder); err != nil {
		return ctxHolder, err
	}

	return ctxHolder, nil
}
