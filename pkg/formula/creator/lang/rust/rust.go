package rust

import (
	"fmt"
	"os"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileextensions"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/rust/template"
)

type Rust struct {
	formula.Lang
	CargoToml          string
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error
}

func New(
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error,
) Rust {
	return Rust{
		Lang: formula.Lang{
			FileFormat:   fileextensions.Rust,
			StartFile:    template.StartFile,
			Main:         template.Main,
			Makefile:     template.Makefile,
			Run:          template.Run,
			Dockerfile:   template.Dockerfile,
			File:         template.File,
			WindowsBuild: template.WindowsBuild,
			Compiled:     false,
			UpperCase:    false,
		},
		CargoToml:          template.CargoToml,
		createGenericFiles: createGenericFiles,
	}
}

func (r Rust) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := r.createGenericFiles(srcDir, pkg, dir, r.Lang); err != nil {
		return err
	}

	runTemplatePath := fmt.Sprintf("%s/run_template", srcDir)
	if err := fileutil.WriteFilePerm(runTemplatePath, []byte(r.Run), 0777); err != nil {
		return err
	}

	newPkgDir := srcDir + "/src/" + pkg
	if err := fileutil.CreateDirIfNotExists(newPkgDir, os.ModePerm); err != nil {
		return err
	}

	if err := fileutil.MoveFiles(srcDir, srcDir+"/src", []string{"main.rs"}); err != nil {
		return err
	}

	if err := CreateCargoFile(srcDir, r.CargoToml); err != nil {
		return err
	}

	templateRust := strings.ReplaceAll(r.File, formula.NameBin, pkg)
	pkgFile := fmt.Sprintf("%s/mod.rs", newPkgDir)
	if err := fileutil.WriteFile(pkgFile, []byte(templateRust)); err != nil {
		return err
	}

	return nil
}

func CreateCargoFile(dir, cargoTomlTemplate string) error {
	return fileutil.WriteFile(fmt.Sprintf("%s/Cargo.toml", dir), []byte(cargoTomlTemplate))
}
