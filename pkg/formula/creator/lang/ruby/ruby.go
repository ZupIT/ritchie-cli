package ruby

import (
	"fmt"
	"os"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileextensions"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/ruby/template"
)

type Ruby struct {
	formula.Lang
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error
	Gemfile string
}

func New(
	c formula.Creator,
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error,
) Ruby {
	return Ruby{
		Lang: formula.Lang{
			Creator:      c,
			FileFormat:   fileextensions.Ruby,
			StartFile:    template.StartFile,
			Main:         template.Index,
			Makefile:     template.Makefile,
			Run:          template.Run,
			Dockerfile:   template.Dockerfile,
			File:         template.File,
			WindowsBuild: template.WindowsBuild,
			Compiled:     false,
			UpperCase:    false,
		},
		createGenericFiles: createGenericFiles,
		Gemfile:      template.Gemfile,
	}
}

func (n Ruby) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := n.createGenericFiles(srcDir, pkg, dir, n.Lang); err != nil {
		return err
	}

	runTemplatePath := fmt.Sprintf("%s/run_template", srcDir)
	if err := fileutil.WriteFilePerm(runTemplatePath, []byte(n.Run), 0777); err != nil {
		return err
	}

	if err := fileutil.CreateDirIfNotExists(pkgDir, os.ModePerm); err != nil {
		return err
	}

	if err := createGemfile(srcDir, n.Gemfile); err != nil {
		return err
	}

	templateNode := strings.ReplaceAll(n.File, formula.NameBin, pkg)
	pkgFile := fmt.Sprintf("%s/%s%s", pkgDir, pkg, n.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(templateNode)); err != nil {
		return err
	}

	return nil
}

func createGemfile(dir, tpl string) error {
	return fileutil.WriteFile(fmt.Sprintf("%s/Gemfile", dir), []byte(tpl))
}
