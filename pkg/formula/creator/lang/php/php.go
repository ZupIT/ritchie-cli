package php

import (
	"fmt"
	"os"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileextensions"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/php/template"
)

type Php struct {
	formula.Lang
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error
}

func New(
	c formula.Creator,
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error,
) Php {
	return Php{
		Lang: formula.Lang{
			Creator:      c,
			FileFormat:   fileextensions.Php,
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
	}
}

func (p Php) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := p.createGenericFiles(srcDir, pkg, dir, p.Lang); err != nil {
		return err
	}

	runTemplatePath := fmt.Sprintf("%s/run_template", srcDir)
	if err := fileutil.WriteFilePerm(runTemplatePath, []byte(p.Run), 0777); err != nil {
		return err
	}

	if err := fileutil.CreateDirIfNotExists(pkgDir, os.ModePerm); err != nil {
		return err
	}

	templatePHP := strings.ReplaceAll(p.File, formula.NameBin, pkg)
	pkgFile := fmt.Sprintf("%s/%s%s", pkgDir, pkg, p.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(templatePHP)); err != nil {
		return err
	}

	return nil
}
