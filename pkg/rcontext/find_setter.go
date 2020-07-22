package rcontext

import "fmt"

type FindSetterManager struct {
	ctxFile string
	CtxFinder
	Setter
}

func NewFindSetter(homePath string, f CtxFinder, s Setter) FindSetterManager {
	return FindSetterManager{fmt.Sprintf(ContextPath, homePath), f, s}
}
