package rtutorial

import (
	"encoding/json"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

type FindManager struct {
	tutorialFile string
	homePath     string
}

func NewFinder(homePath string) FindManager {
	return FindManager{
		tutorialFile: fmt.Sprintf(TutorialPath, homePath),
		homePath:     homePath,
	}
}

func (f FindManager) Find() (TutorialHolder, error) {
	tutorialHolder := TutorialHolder{Current: DefaultTutorial}

	if !fileutil.Exists(f.tutorialFile) {
		setter := NewSetter(f.homePath)

		tutorialHolder, err := setter.Set(DefaultTutorial)
		if err != nil {
			return tutorialHolder, err
		}

		return tutorialHolder, nil
	}

	file, err := fileutil.ReadFile(f.tutorialFile)
	if err != nil {
		return tutorialHolder, err
	}

	if err := json.Unmarshal(file, &tutorialHolder); err != nil {
		return tutorialHolder, err
	}

	return tutorialHolder, nil
}
