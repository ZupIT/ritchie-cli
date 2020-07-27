package rtutorial

import (
	"encoding/json"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type SetterManager struct {
	tutorialFile string
	fw           stream.FileWriter
}

func NewSetter(homePath string, fw stream.FileWriter) Setter {
	return SetterManager{
		tutorialFile: fmt.Sprintf(TutorialPath, homePath),
		fw:           fw,
	}
}

func (s SetterManager) Set(tutorial string) (TutorialHolder, error) {
	tutorialHolder := TutorialHolder{Current: DefaultTutorial}
	tutorialHolderDefault := TutorialHolder{Current: DefaultTutorial}

	tutorialHolder.Current = tutorial

	b, err := json.Marshal(&tutorialHolder)
	if err != nil {
		return tutorialHolderDefault, err
	}

	if err := s.fw.Write(s.tutorialFile, b); err != nil {
		return tutorialHolderDefault, err
	}

	return tutorialHolder, nil
}
