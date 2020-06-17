package fileutil

type FileUtilService interface {
	ReadFile(path string) ([]byte, error)
	WriteFilePerm(path string, content []byte, perm int32) error
}

type DefaultFileUtilService struct{}

func (s DefaultFileUtilService) ReadFile(path string) ([]byte, error) {
	return ReadFile(path)
}

func (s DefaultFileUtilService) WriteFilePerm(path string, content []byte, perm int32) error {
	return WriteFilePerm(path, content, perm)
}
