package node

import (
	"fmt"
	"os"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileextensions"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/node/template"
)

type Node struct {
	formula.Lang
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error
}

func New(
	c formula.Creator,
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error,
) Node {
	return Node{
		Lang: formula.Lang{
			Creator:      c,
			FileFormat:   fileextensions.JavaScript,
			StartFile:    template.StartFile,
			Main:         template.Index,
			Makefile:     template.Makefile,
			Run:          template.Run,
			Dockerfile:   template.Dockerfile,
			PackageJson:  template.PackageJson,
			File:         template.File,
			WindowsBuild: template.WindowsBuild,
			Compiled:     false,
			UpperCase:    false,
		},
		createGenericFiles: createGenericFiles,
	}
}

func (n Node) Create(srcDir, pkg, pkgDir, dir string) error {
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

	if err := createPackageJson(srcDir, n.PackageJson); err != nil {
		return err
	}

	templateNode := strings.ReplaceAll(n.File, formula.NameBin, pkg)
	pkgFile := fmt.Sprintf("%s/%s%s", pkgDir, pkg, n.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(templateNode)); err != nil {
		return err
	}

	return nil
}

func createPackageJson(dir, tpl string) error {
	return fileutil.WriteFile(fmt.Sprintf("%s/package.json", dir), []byte(tpl))
}
