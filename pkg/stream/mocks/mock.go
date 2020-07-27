package stream

type FileReadExisterCustomMock struct {
	ReadMock   func(path string) ([]byte, error)
	ExistsMock func(path string) bool
}

// Read of FileManagerCustomMock
func (fmc FileReadExisterCustomMock) Read(path string) ([]byte, error) {
	return fmc.ReadMock(path)
}

// Exists of FileManagerCustomMock
func (fmc FileReadExisterCustomMock) Exists(path string) bool {
	return fmc.ExistsMock(path)
}

type FileWriterCustomMock struct {
	WriteMock func(path string, content []byte) error
}

func (wcm FileWriterCustomMock) Write(path string, content []byte) error {
	return wcm.WriteMock(path, content)
}

type FileWriterMock struct{}

func (FileWriterMock) Write(path string, content []byte) error {
	return nil
}


