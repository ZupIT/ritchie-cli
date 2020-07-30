/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package watcher

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"
	"github.com/radovskyb/watcher"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type WatchManager struct {
	watcher *watcher.Watcher
	formula formula.LocalBuilder
	dir     stream.DirListChecker
}

func New(formula formula.LocalBuilder, dir stream.DirListChecker) *WatchManager {
	w := watcher.New()

	return &WatchManager{watcher: w, formula: formula, dir: dir}
}

func (w *WatchManager) Watch(workspacePath, formulaPath string) {
	w.watcher.FilterOps(watcher.Write)
	go func() {
		for {
			select {
			case event := <-w.watcher.Event:
				if !event.IsDir() && !strings.Contains(event.Path, "/dist") {
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

	if err := w.watcher.AddRecursive(formulaPath); err != nil {
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
	buildInfo := prompt.Bold("Building formula...")
	s := spinner.StartNew(buildInfo)
	time.Sleep(2 * time.Second)

	if err := w.formula.Build(workspacePath, formulaPath); err != nil {
		errorMsg := prompt.Red(err.Error())
		s.Error(errors.New(errorMsg))
		return
	}

	success := prompt.Green("✔ Build completed!")
	s.Success(success)
	prompt.Info("Now you can run your formula with Ritchie!")
}
