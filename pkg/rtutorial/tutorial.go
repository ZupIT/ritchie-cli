package rtutorial

const TutorialPath = "%s/tutorial.json"

const DefaultTutorial = "enabled"

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
