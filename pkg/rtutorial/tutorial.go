package rtutorial

const TutorialPath = "%s/Tutorial"

const DefaultTutorial = "on"

type TutorialHolder struct {
	Current string `json:"tutorial"`
}

type Setter interface {
	Set(tutorial string) (TutorialHolder, error)
}

type Finder interface {
	Find() (TutorialHolder, error)
}

type FindSetter interface {
	Finder
	Setter
}
