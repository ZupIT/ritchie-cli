package rtutorial

import "fmt"

type FindSetterManager struct {
	tutorialFile string
	Finder
	Setter
}

func NewFindSetter(homePath string, f Finder, s Setter) FindSetterManager {
	return FindSetterManager{fmt.Sprintf(TutorialPath, homePath), f, s}
}
