package rtutorial

import (
	"encoding/json"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

type FindManager struct {
	tutorialFile string
}

func NewFinder(homePath string) FindManager {
	return FindManager{
		tutorialFile: fmt.Sprintf(TutorialPath, homePath),
	}
}

func (f FindManager) Find() (TutorialHolder, error) {
	tutorialHolder := TutorialHolder{Current: DefaultTutorial}

	if !fileutil.Exists(f.tutorialFile) {
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
