package rcontext

type FindSetterManager struct {
	Finder
	Setter
}

func NewFindSetter(f Finder, s Setter) FindSetterManager {
	return FindSetterManager{f, s}
}
