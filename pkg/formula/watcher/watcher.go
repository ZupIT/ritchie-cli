package watcher

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"
	"github.com/radovskyb/watcher"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type WatchManager struct {
	watcher *watcher.Watcher
	formula formula.Builder
	dir     stream.DirListChecker
}

func New(formula formula.Builder, dir stream.DirListChecker) *WatchManager {
	w := watcher.New()

	return &WatchManager{watcher: w, formula: formula, dir: dir}
}

func (w *WatchManager) Watch(workspacePath, formulaPath string) {
	w.watcher.FilterOps(watcher.Write)
	go func() {
		for {
			select {
			case event := <-w.watcher.Event:
				if !event.IsDir() {
					w.build(workspacePath, formulaPath)
					prompt.Info("Waiting for changes...")
				}
			case err := <-w.watcher.Error:
				prompt.Error(err.Error())
			case <-w.watcher.Closed:
				return
			}
		}
	}()

	formulaSrc := fmt.Sprintf("%s/src", formulaPath)
	if err := w.watcher.AddRecursive(formulaSrc); err != nil {
		log.Fatalln(err)
	}

	w.build(workspacePath, formulaPath)

	watchText := fmt.Sprintf("Watching dir %s \n", formulaPath)
	prompt.Info(watchText)

	if err := w.watcher.Start(time.Second * 2); err != nil {
		log.Fatalln(err)
	}
}

func (w WatchManager) build(workspacePath, formulaPath string) {
	buildInfo := fmt.Sprintf(prompt.Teal, "Building formula...")
	s := spinner.StartNew(buildInfo)
	time.Sleep(2 * time.Second)

	if err := w.formula.Build(workspacePath, formulaPath); err != nil {
		errorMsg := fmt.Sprintf(prompt.Red, err)
		s.Error(errors.New(errorMsg))
		return
	}

	success := fmt.Sprintf(prompt.Green, "✔ Build completed!")
	s.Success(success)
	prompt.Info("Now you can run your formula with Ritchie!")
}
