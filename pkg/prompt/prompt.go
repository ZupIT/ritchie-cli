package prompt

type InputText interface {
	Text(name string, required bool, helper ...string) (string, error)
}

type InputTextValidator interface {
	Text(name string, validate func(interface{}) error, helper ...string) (string, error)
}

type InputBool interface {
	Bool(name string, items []string) (bool, error)
}

type InputPassword interface {
	Password(label string) (string, error)
}

type InputMultiline interface {
	MultiLineText(name string, required bool) (string, error)
}

type InputList interface {
	List(name string, items []string) (string, error)
}

type InputInt interface {
	Int(name string, helper ...string) (int64, error)
}

type InputEmail interface {
	Email(name string) (string, error)
}

type InputURL interface {
	URL(name, defaultValue string) (string, error)
}
