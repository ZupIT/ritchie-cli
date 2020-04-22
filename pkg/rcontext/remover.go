package rcontext

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type RemoveManager struct {
	ctxFile string
	finder  Finder
	file    stream.FileWriter
}

func NewRemover(homePath string, f Finder, w stream.FileWriter) RemoveManager {
	return RemoveManager{
		ctxFile: fmt.Sprintf(ContextPath, homePath),
		finder:  f,
		file:    w,
	}
}

func (r RemoveManager) Remove(ctx string) (ContextHolder, error) {
	ctxHolder, err := r.finder.Find()
	if err != nil {
		return ContextHolder{}, err
	}

	ctx = strings.ReplaceAll(ctx, CurrentCtx, "")
	if ctxHolder.Current == ctx {
		ctxHolder.Current = ""
	}

	for i, context := range ctxHolder.All {
		if ctx == context {
			ctxHolder.All = append(ctxHolder.All[:i], ctxHolder.All[i+1:]...)
			break
		}
	}

	b, err := json.Marshal(&ctxHolder)
	if err != nil {
		return ContextHolder{}, err
	}
	if err := r.file.Write(r.ctxFile, b); err != nil {
		return ContextHolder{}, err
	}

	return ctxHolder, nil
}
