package rcontext

import "fmt"

type FindRemoverManager struct {
	ctxFile string
	CtxFinder
	Remover
}

func NewFindRemover(homePath string, f CtxFinder, r Remover) FindRemoverManager {
	return FindRemoverManager{fmt.Sprintf(ContextPath, homePath), f, r}
}
