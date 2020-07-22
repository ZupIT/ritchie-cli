package mock

import (
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

type TutorialSetterMock struct{}

func (TutorialSetterMock) Set(tutorial string) (rtutorial.TutorialHolder, error) {
	return rtutorial.TutorialHolder{}, nil
}

type TutorialFinderMock struct{}

func (TutorialFinderMock) Find() (rtutorial.TutorialHolder, error) {
	return rtutorial.TutorialHolder{}, nil
}

type TutorialFindSetterMock struct{}

func (TutorialFindSetterMock) Find() (rtutorial.TutorialHolder, error) {
	f := TutorialFinderMock{}
	return f.Find()
}

func (TutorialFindSetterMock) Set(tutorial string) (rtutorial.TutorialHolder, error) {
	s := TutorialSetterMock{}
	return s.Set(tutorial)
}
