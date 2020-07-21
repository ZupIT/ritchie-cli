package rtutorial

import (
	"encoding/json"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

type SetterManager struct {
	tutorialFile string
	finder       Finder
}

func NewSetter(homePath string, f Finder) Setter {
	return SetterManager{tutorialFile: fmt.Sprintf(TutorialPath, homePath), finder: f}
}

func (s SetterManager) Set(tutorial string) (TutorialHolder, error) {
	tutorialHolder, err := s.finder.Find()
	if err != nil {
		return TutorialHolder{Current: DefaultTutorial}, err
	}

	tutorialHolder.Current = tutorial

	b, err := json.Marshal(&tutorialHolder)
	if err != nil {
		return TutorialHolder{}, err
	}
	if err := fileutil.WriteFilePerm(s.tutorialFile, b, 0600); err != nil {
		return TutorialHolder{}, err
	}

	return tutorialHolder, nil
}
