package rtutorial

import (
	"encoding/json"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type FindManager struct {
	tutorialFile string
	homePath     string
	fr           stream.FileReadExister
}

func NewFinder(homePath string, fr stream.FileReadExister) FindManager {
	return FindManager{
		tutorialFile: fmt.Sprintf(TutorialPath, homePath),
		homePath:     homePath,
		fr:           fr,
	}
}

func (f FindManager) Find() (TutorialHolder, error) {
	tutorialHolder := TutorialHolder{Current: DefaultTutorial}

	if !f.fr.Exists(f.tutorialFile) {
		return tutorialHolder, nil
	}

	file, err := f.fr.Read(f.tutorialFile)
	if err != nil {
		return tutorialHolder, err
	}

	if err := json.Unmarshal(file, &tutorialHolder); err != nil {
		return tutorialHolder, err
	}

	return tutorialHolder, nil
}
