package template

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

const (
	languageDir = "languages"
	rootDir     = "root"
)

var (
	templatePath = []string{"repos", "commons", "templates", "create_formula"}
)

type Manager interface {
	Languages() ([]string, error)
	LangTemplateFiles(lang string) ([]File, error)
	ResolverNewPath(oldPath, newDir, lang, workspacePath string) (string, error)
}

type File struct {
	Path  string
	IsDir bool
}

func NewManager(ritchieHome string) Manager {
	return DefaultManager{ritchieHome}
}

type DefaultManager struct {
	ritchieHome string
}

func (tm DefaultManager) templateDir() string {
	return path.Join(tm.ritchieHome, path.Join(templatePath...))
}

func (tm DefaultManager) Languages() ([]string, error) {
	tplD := path.Join(tm.templateDir(), languageDir)

	dirs, err := ioutil.ReadDir(tplD)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, d := range dirs {
		if d.IsDir() {
			result = append(result, d.Name())
		}
	}

	return result, nil
}

func (tm DefaultManager) LangTemplateFiles(lang string) ([]File, error) {
	tplD := tm.templateDir()

	langDir := path.Join(tplD, languageDir, lang)

	languageTpl, err := readDirRecursive(langDir)
	if err != nil {
		return nil, err
	}

	rootTplDir := path.Join(tplD, rootDir)
	rootTpl, err := readDirRecursive(rootTplDir)
	if err != nil {
		return nil, err
	}

	return append(languageTpl, rootTpl...), nil
}

func readDirRecursive(dir string) ([]File, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var fileNames []File
	for _, f := range files {
		if f.IsDir() {
			dirFiles, err := readDirRecursive(path.Join(dir, f.Name()))
			if err != nil {
				return nil, err
			}
			fileNames = append(fileNames, dirFiles...)
		}
		fileNames = append(fileNames, File{
			Path:  path.Join(dir, f.Name()),
			IsDir: f.IsDir(),
		})

	}
	return fileNames, nil
}

func (tm DefaultManager) ResolverNewPath(oldPath, formulaPath, lang, workspacePath string) (string, error) {
	tplD := tm.templateDir()
	langTplPath := path.Join(tplD, languageDir, lang)
	rootTplPath := path.Join(tplD, rootDir)

	if strings.Contains(oldPath, rootTplPath) {
		return strings.Replace(oldPath, rootTplPath, workspacePath, 1), nil
	}

	if strings.Contains(oldPath, langTplPath) {
		return strings.Replace(oldPath, langTplPath, formulaPath, 1), nil
	}

	return "", fmt.Errorf("fail to resolve new Path %s", oldPath)

}
