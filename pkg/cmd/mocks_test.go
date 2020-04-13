package cmd

type inputTextMock struct{}

func (inputTextMock) Text(name string, required bool) (string, error) {
	return "mocked text", nil
}

type inputURLMock struct{}

func (inputURLMock) URL(name, defaultValue string) (string, error) {
	return "http://localhost/mocked", nil
}

type inputIntMock struct{}

func (inputIntMock) Int(name string) (int64, error) {
	return 0, nil
}
