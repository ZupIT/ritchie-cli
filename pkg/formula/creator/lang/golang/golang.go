package golang

import (
	"fmt"
	"os"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileextensions"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/golang/template"
)

type Go struct {
	formula.Lang
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error
}

func New(
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error,
) Go {
	return Go{
		Lang: formula.Lang{
			FileFormat:   fileextensions.Go,
			StartFile:    template.StartFile,
			Main:         template.Main,
			Makefile:     template.Makefile,
			Dockerfile:   template.Dockerfile,
			Pkg:          template.Pkg,
			WindowsBuild: template.WindowsBuild,
			Compiled:     true,
			UpperCase:    false,
		},
		createGenericFiles: createGenericFiles,
	}
}

func (g Go) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := g.createGenericFiles(srcDir, pkg, dir, g.Lang); err != nil {
		return err
	}

	if err := createGoModFile(srcDir, pkg); err != nil {
		return err
	}

	if err := fileutil.CreateDirIfNotExists(pkgDir, os.ModePerm); err != nil {
		return err
	}

	templateGo := strings.ReplaceAll(g.Pkg, formula.NameModule, pkg)
	pkgFile := fmt.Sprintf("%s/%s%s", pkgDir, pkg, g.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(templateGo)); err != nil {
		return err
	}
	return nil
}

func createGoModFile(dir, pkg string) error {
	tplFile := template.GoMod
	tplFile = strings.ReplaceAll(tplFile, formula.NameModule, pkg)
	return fileutil.WriteFile(fmt.Sprintf("%s/go.mod", dir), []byte(tplFile))
}
