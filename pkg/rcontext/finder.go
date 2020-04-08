package rcontext

import (
	"encoding/json"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

type FindManager struct {
	ctxFile string
}

func NewFinder(homePath string) FindManager {
	return FindManager{ctxFile: fmt.Sprintf(ContextPath, homePath)}
}

func (f FindManager) Find() (ContextHolder, error) {
	ctxHolder := ContextHolder{}

	if !fileutil.Exists(f.ctxFile) {
		return ctxHolder, nil
	}

	file, err := fileutil.ReadFile(f.ctxFile)
	if err != nil {
		return ctxHolder, err
	}

	if err := json.Unmarshal(file, &ctxHolder); err != nil {
		return ctxHolder, err
	}

	return ctxHolder, nil
}
