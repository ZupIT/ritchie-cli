package rcontext

import (
	"encoding/json"
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"strings"
)

type RemoveManager struct {
	ctxFile string
	finder  CtxFinder
}

func NewRemover(homePath string, f CtxFinder) RemoveManager {
	return RemoveManager{ctxFile: fmt.Sprintf(ContextPath, homePath), finder: f}
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
	if err := fileutil.WriteFilePerm(r.ctxFile, b, 0600); err != nil {
		return ContextHolder{}, err
	}

	return ctxHolder, nil
}
