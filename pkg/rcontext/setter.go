package rcontext

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type SetterManager struct {
	ctxFile string
	finder  Finder
	writer  stream.FileWriter
}

func NewSetter(homePath string, f Finder, w stream.FileWriter) Setter {
	return SetterManager{
		ctxFile: fmt.Sprintf(ContextPath, homePath),
		finder:  f,
		writer:  w,
	}
}

func (s SetterManager) Set(ctx string) (ContextHolder, error) {
	ctxHolder, err := s.finder.Find()
	if err != nil {
		return ContextHolder{}, err
	}

	ctxHolder.Current = strings.ReplaceAll(ctx, DefaultCtx, "")
	if ctx != DefaultCtx {
		if ctxHolder.All == nil {
			ctxHolder.All = make([]string, 0)
		}

		if !sliceutil.Contains(ctxHolder.All, ctx) {
			ctxHolder.All = append(ctxHolder.All, ctx)
		}
	}

	b, err := json.Marshal(&ctxHolder)
	if err != nil {
		return ContextHolder{}, err
	}
	if err := s.writer.Write(s.ctxFile, b); err != nil {
		return ContextHolder{}, err
	}

	return ctxHolder, nil
}
