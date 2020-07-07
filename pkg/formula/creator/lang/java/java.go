package java

import (
	"fmt"
	"os"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileextensions"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/java/template"
)

type Java struct {
	formula.Lang
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error
}

func New(
	c formula.Creator,
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error,
) Java {
	return Java{
		Lang: formula.Lang{
			Creator:      c,
			FileFormat:   fileextensions.Java,
			StartFile:    template.StartFile,
			Main:         template.Main,
			Makefile:     template.Makefile,
			Run:          template.Run,
			Dockerfile:   template.Dockerfile,
			File:         template.File,
			WindowsBuild: template.WindowsBuild,
			Compiled:     false,
			UpperCase:    true,
		},
		createGenericFiles: createGenericFiles,
	}
}

func (j Java) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := j.createGenericFiles(srcDir, pkg, dir, j.Lang); err != nil {
		return err
	}

	runTemplate := fmt.Sprintf("%s/run_template", srcDir)
	if err := fileutil.WriteFilePerm(runTemplate, []byte(j.Run), 0777); err != nil {
		return err
	}

	if err := fileutil.CreateDirIfNotExists(pkgDir, os.ModePerm); err != nil {
		return err
	}

	templateFileJava := strings.ReplaceAll(j.File, formula.NameBin, pkg)
	firstUpper := strings.Title(strings.ToLower(pkg))
	templateFileJava = strings.ReplaceAll(templateFileJava, formula.NameBinFirstUpper, firstUpper)
	pkgFile := fmt.Sprintf("%s/%s%s", pkgDir, firstUpper, j.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(templateFileJava)); err != nil {
		return err
	}

	return nil
}
