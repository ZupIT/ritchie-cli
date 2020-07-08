package java

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileextensions"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/java/template"
)

type Java struct {
	fCmdName string
	formula.Lang
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error
}

func New(
	c formula.Creator,
	createGenericFiles func(srcDir, pkg, dir string, l formula.Lang) error,
	fCmdName string,
) Java {
	return Java{
		Lang: formula.Lang{
			Creator:      c,
			FileFormat:   fileextensions.Java,
			StartFile:    template.StartFile,
			Main:         template.Main,
			Makefile:     template.Makefile,
			Dockerfile:   template.Dockerfile,
			File:         template.File,
			WindowsBuild: template.WindowsBuild,
			Compiled:     false,
			UpperCase:    true,
		},
		createGenericFiles: createGenericFiles,
		fCmdName:           fCmdName,
	}
}

func (j Java) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := j.createGenericFiles(srcDir, pkg, dir, j.Lang); err != nil {
		return err
	}

	//Todo change it
	fileutil.RemoveFile(path.Join(srcDir, "Main.java"))

	groupId := strings.ReplaceAll(j.fCmdName, " ", ".")
	baseJavaDir := strings.Split("src/main/java", "/")
	groupIdDir := strings.Split(groupId, ".")
	javaSrcDir := path.Join(srcDir, path.Join(baseJavaDir...), path.Join(groupIdDir...))

	if err := fileutil.CreateDirIfNotExists(javaSrcDir, os.ModePerm); err != nil {
		return err
	}

	firstUpper := strings.Title(strings.ToLower(pkg))

	createMainFile(groupId, firstUpper, pkg, javaSrcDir)

	pom := strings.ReplaceAll(template.Pom, "#rit{{groupId}}", groupId)
	pom = strings.ReplaceAll(pom, "#rit{{artifactId}}", pkg)
	fileutil.WriteFile(path.Join(srcDir, "pom.xml"), []byte(pom))

	err := createPkgFile(j, pkg, groupId, firstUpper, javaSrcDir)
	if err != nil {
		return err
	}

	return nil
}

func createPkgFile(j Java, pkg string, groupId string, firstUpper string, javaSrcDir string) error {
	templateFileJava := strings.ReplaceAll(j.File, formula.NameBin, pkg)
	templateFileJava = strings.ReplaceAll(templateFileJava, "{{final-pkg}}", groupId)
	templateFileJava = strings.ReplaceAll(templateFileJava, formula.NameBinFirstUpper, firstUpper)

	templateFileDir := path.Join(javaSrcDir, pkg)
	if err := fileutil.CreateDirIfNotExists(templateFileDir, os.ModePerm); err != nil {
		return err
	}
	pkgTemplateFile := fmt.Sprintf("%s/%s%s", templateFileDir, firstUpper, j.FileFormat)
	if err := fileutil.WriteFile(pkgTemplateFile, []byte(templateFileJava)); err != nil {
		return err
	}
	return nil
}

func createMainFile(groupId string, firstUpper string, pkg string, javaSrcDir string) {
	mainFile := strings.ReplaceAll(template.Main, "{{final-pkg}}", groupId)
	mainFile = strings.ReplaceAll(mainFile, formula.NameBinFirstUpper, firstUpper)
	mainFile = strings.ReplaceAll(mainFile, formula.NameBin, pkg)

	fileutil.WriteFile(path.Join(javaSrcDir, "Main.java"), []byte(mainFile))
}
