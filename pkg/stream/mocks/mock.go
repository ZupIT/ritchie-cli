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
