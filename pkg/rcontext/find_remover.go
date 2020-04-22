package rcontext

type FindRemoverManager struct {
	Finder
	Remover
}

func NewFindRemover(f Finder, r Remover) FindRemoverManager {
	return FindRemoverManager{f, r}
}
