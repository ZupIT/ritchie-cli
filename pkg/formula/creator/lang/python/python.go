package python

import (
	"fmt"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileextensions"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/python/template"
)

type Python struct {
	formula.Lang
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error
}

func New(
	c formula.Creator,
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error,
) Python {
	return Python{
		Lang: formula.Lang{
			Creator:      c,
			FileFormat:   fileextensions.Python,
			StartFile:    template.StartFile,
			Main:         template.Main,
			Makefile:     template.Makefile,
			Dockerfile:   template.Dockerfile,
			File:         template.File,
			WindowsBuild: template.WindowsBuild,
			Compiled:     false,
			UpperCase:    false,
		},
		createGenericFiles: createGenericFiles,
	}
}

func (p Python) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := p.createGenericFiles(srcDir, pkg, dir, p.Lang); err != nil {
		return err
	}

	if err := fileutil.CreateDirIfNotExists(pkgDir, os.ModePerm); err != nil {
		return err
	}

	pkgFile := fmt.Sprintf("%s/%s%s", pkgDir, pkg, p.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(p.File)); err != nil {
		return err
	}

	return nil
}
