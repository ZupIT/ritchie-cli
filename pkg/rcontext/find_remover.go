package rcontext

import "fmt"

type FindRemoverManager struct {
	ctxFile string
	Finder
	Remover
}

func NewFindRemover(homePath string, f Finder, r Remover) FindRemoverManager {
	return FindRemoverManager{fmt.Sprintf(ContextPath, homePath), f, r}
}
