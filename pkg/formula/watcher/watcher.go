package watcher

import (
	"fmt"
	"log"
	"time"

	"github.com/radovskyb/watcher"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/spinner"
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

func (w *WatchManager) build(workspacePath, formulaPath string) {
	buildInfo := fmt.Sprintf(prompt.Teal, "Building formula...")
	s := spinner.New(buildInfo)
	s.Start()
	stderr, err := w.formula.Build(workspacePath, formulaPath)
	if err != nil {
		s.Stop()
		msgFormatted := fmt.Sprintf("Build error: \n%s", string(stderr))
		prompt.Error(msgFormatted)
	} else {
		s.Stop()
		prompt.Success("✔ Build completed!")
	}
}
