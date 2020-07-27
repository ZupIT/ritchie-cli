package rcontext

import "fmt"

type FindSetterManager struct {
	ctxFile string
	Finder
	Setter
}

func NewFindSetter(homePath string, f Finder, s Setter) FindSetterManager {
	return FindSetterManager{fmt.Sprintf(ContextPath, homePath), f, s}
}
