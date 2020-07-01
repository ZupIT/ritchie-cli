package runner

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
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

func (o OutputManager) ValidAndPrint(setup formula.Setup) error {
	_, err := fmt.Fprintf(o.writer, printAndValidOutputDir(setup))
	return err
}

func printAndValidOutputDir(setup formula.Setup) string {

	files, err := ioutil.ReadDir(setup.TmpOutputDir)
	if err != nil {
		return prompt.Red("Fail to read output dir")
	}
	fOutputs := map[string]string{}

	resolveKey := func(name string) string { return strings.ToUpper(name) }

	if len(files) != len(setup.Config.Outputs) {
		return prompt.Red("Output dir size is different of outputs array in config.json")
	}

	for _, f := range files {
		fName := fmt.Sprintf("%s/%s", setup.TmpOutputDir, f.Name())
		key := resolveKey(f.Name())
		b, err := ioutil.ReadFile(fName)
		if err != nil {
			return prompt.Red("fail to read file: " + fName)
		}
		fOutputs[key] = string(b)
	}

	var result string
	for _, o := range setup.Config.Outputs {
		key := resolveKey(o.Name)
		v, exist := fOutputs[key]
		if !exist {
			return prompt.Red("file:" + key + " not found in output dir")
		}
		if o.Print {
			result += fmt.Sprintf("%s=%s\n", key, v)
		}
	}
	return result
}
