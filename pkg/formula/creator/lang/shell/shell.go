package shell

import (
	"fmt"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileextensions"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/shell/template"
)

type Shell struct {
	formula.Lang
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error
}

func New(
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error,
) Shell {
	return Shell{
		Lang: formula.Lang{
			FileFormat: fileextensions.Shell,
			StartFile:  template.StartFile,
			Main:       template.Main,
			Makefile:   template.Makefile,
			Dockerfile: template.Dockerfile,
			File:       template.File,
			Compiled:   false,
			UpperCase:  false,
		},
		createGenericFiles: createGenericFiles,
	}
}

func (s Shell) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := s.createGenericFiles(srcDir, pkg, dir, s.Lang); err != nil {
		return err
	}

	if err := fileutil.CreateDirIfNotExists(pkgDir, os.ModePerm); err != nil {
		return err
	}

	pkgFile := fmt.Sprintf("%s/%s%s", pkgDir, pkg, s.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(s.File)); err != nil {
		return err
	}

	return nil
}
