package runner

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var (
	ErrReadOutputDir  = errors.New(prompt.Red("fail to read output dir"))
	ErrValidOutputDir = errors.New(prompt.Red("Output dir size is different of outputs array in config.json"))
)

type OutputManager struct {
	writer io.Writer
}

func NewOutputManager(
	w io.Writer,
) OutputManager {
	return OutputManager{
		writer: w,
	}
}

func (o OutputManager) Outputs(setup formula.Setup) error {
	msg, err := printAndValidOutputDir(setup)
	if err != nil {
		return err
	}
	_, err = o.writer.Write([]byte(msg))
	return err
}

func printAndValidOutputDir(setup formula.Setup) (string, error) {

	files, err := ioutil.ReadDir(setup.TmpOutputDir)
	if err != nil {
		return "", ErrReadOutputDir
	}
	fOutputs := map[string]string{}

	resolveKey := func(name string) string { return strings.ToUpper(name) }

	if len(files) != len(setup.Config.Outputs) {
		return "", ErrValidOutputDir
	}

	for _, f := range files {
		fName := fmt.Sprintf("%s/%s", setup.TmpOutputDir, f.Name())
		key := resolveKey(f.Name())
		b, err := ioutil.ReadFile(fName)
		if err != nil {
			return "", errors.New(prompt.Red("fail to read file: " + fName))
		}
		fOutputs[key] = string(b)
	}

	var result string
	for _, o := range setup.Config.Outputs {
		key := resolveKey(o.Name)
		v, exist := fOutputs[key]
		if !exist {
			return "", errors.New(prompt.Red("file:" + key + " not found in output dir"))
		}
		if o.Print {
			result += fmt.Sprintf("%s=%s\n", key, v)
		}
	}
	return result, nil
}
