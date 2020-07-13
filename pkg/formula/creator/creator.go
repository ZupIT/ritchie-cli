package creator

import (
	"fmt"
	"path"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var (
	ErrRepeatedCommand = prompt.NewError("this command already exists")
)

type CreateManager struct {
	treeManager tree.Manager
	dir         stream.DirCreater
	file        stream.FileWriteReadExister
	tplM        template.Manager
}

func NewCreator(
	tm tree.Manager,
	dir stream.DirCreater,
	file stream.FileWriteReadExister,
	tplM template.Manager,
) CreateManager {
	return CreateManager{treeManager: tm, dir: dir, file: file, tplM: tplM}
}

func (c CreateManager) Create(cf formula.Create) error {
	if err := c.isValidCmd(cf.FormulaCmd); err != nil {
		return err
	}

	if err := c.dir.Create(cf.WorkspacePath); err != nil {
		return err
	}

	fCmdName := cf.FormulaCmdName()

	if err := c.generateFormulaFiles(cf.FormulaPath, cf.Lang, fCmdName, cf.WorkspacePath); err != nil {
		return err
	}

	return nil
}

func (c CreateManager) isValidCmd(fCmd string) error {
	trees, err := c.treeManager.Tree()
	if err != nil {
		return err
	}

	s := strings.Split(fCmd, " ")
	cp := fmt.Sprintf("root_%s", strings.Join(s[1:len(s)-1], "_"))
	u := s[len(s)-1]
	for _, v := range trees {
		for _, j := range v.Commands {
			if j.Parent == cp && j.Usage == u {
				return ErrRepeatedCommand

			}
		}
	}
	return nil
}

func (c CreateManager) generateFormulaFiles(fPath, lang, fCmdName, workSpcPath string) error {

	if err := c.dir.Create(fPath); err != nil {
		return err
	}

	if err := c.createHelpFiles(fCmdName, workSpcPath); err != nil {
		return err
	}
	if err := c.createUmaskFile(fPath); err != nil {
		return err
	}

	if err := c.applyLangTemplate(lang, fPath, workSpcPath); err != nil {
		return err
	}

	return nil
}

func (c CreateManager) applyLangTemplate(lang, formulaPath, workspacePath string) error {

	tFiles, err := c.tplM.LangTemplateFiles(lang)
	if err != nil {
		return err
	}

	for _, f := range tFiles {
		if f.IsDir {
			newPath, err := c.tplM.ResolverNewPath(f.Path, formulaPath, lang, workspacePath)
			if err != nil {
				return err
			}
			err = c.dir.Create(newPath)
			if err != nil {
				return err
			}
		} else {
			newPath, err := c.tplM.ResolverNewPath(f.Path, formulaPath, lang, workspacePath)
			if err != nil {
				return err
			}
			if c.file.Exists(newPath) {
				continue
			}
			tpl, err := c.file.Read(f.Path)
			if err != nil {
				return err
			}
			newDir, _ := path.Split(newPath)
			err = c.dir.Create(newDir)
			if err != nil {
				return err
			}
			err = c.file.Write(newPath, tpl)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c CreateManager) createUmaskFile(fPath string) error {
	return c.file.Write(path.Join(fPath, "set_umask.sh"), []byte(template.Umask))
}

func (c CreateManager) createHelpFiles(formulaCmdName, workSpacePath string) error {
	dirs := strings.Split(formulaCmdName, " ")
	for i := 0; i < len(dirs); i++ {
		d := dirs[0 : i+1]
		tPath := path.Join(workSpacePath, path.Join(d...))
		helpPath := fmt.Sprintf("%s/help.txt", tPath)
		if !c.file.Exists(helpPath) {
			err := c.file.Write(helpPath, []byte(template.Help))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
