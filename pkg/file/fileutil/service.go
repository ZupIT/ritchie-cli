package fileutil

type Service interface {
	ReadFile(path string) ([]byte, error)
	WriteFilePerm(path string, content []byte, perm int32) error
}

type DefaultService struct{}

func (s DefaultService) ReadFile(path string) ([]byte, error) {
	return ReadFile(path)
}

func (s DefaultService) WriteFilePerm(path string, content []byte, perm int32) error {
	return WriteFilePerm(path, content, perm)
}
