package rtutorial

import (
	"encoding/json"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

type SetterManager struct {
	tutorialFile string
}

func NewSetter(homePath string) Setter {
	return SetterManager{tutorialFile: fmt.Sprintf(TutorialPath, homePath)}
}

func (s SetterManager) Set(tutorial string) (TutorialHolder, error) {
	tutorialHolder := TutorialHolder{Current: DefaultTutorial}
	tutorialHolderDefault := TutorialHolder{Current: DefaultTutorial}

	tutorialHolder.Current = tutorial

	b, err := json.Marshal(&tutorialHolder)
	if err != nil {
		return tutorialHolderDefault, err
	}

	exists := fileutil.Exists(s.tutorialFile)
	if exists {
		if err := fileutil.WriteFilePerm(s.tutorialFile, b, 0600); err != nil {
			return tutorialHolderDefault, err
		}
	} else {
		err = fileutil.CreateFileIfNotExist(s.tutorialFile, b)
		if err != nil {
			return tutorialHolderDefault, err
		}
	}

	return tutorialHolder, nil
}
