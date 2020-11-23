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
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"
	"github.com/radovskyb/watcher"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const stoppedText = "Press CTRL+C to stop"

type WatchManager struct {
	watcher    *watcher.Watcher
	formula    formula.Builder
	dir        stream.DirListChecker
	sendMetric func(cmd metric.SendCommandDataParams)
}

func New(
	formula formula.Builder,
	dir stream.DirListChecker,
	sendMetric func(cmd metric.SendCommandDataParams),
) *WatchManager {

	w := watcher.New()

	return &WatchManager{
		watcher:    w,
		formula:    formula,
		dir:        dir,
		sendMetric: sendMetric,
	}
}

func (w *WatchManager) closeWatch() {
	fmt.Println("\nStopping...")

	w.watcher.Wait()
	w.watcher.Close()
}

func (w *WatchManager) Watch(formulaPath string, workspace formula.Workspace) {
	w.watcher.FilterOps(watcher.Write)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case event := <-w.watcher.Event:
				if !event.IsDir() && !strings.Contains(event.Path, "/dist") {
					w.build(formulaPath, workspace)
					fmt.Println(prompt.Bold("Waiting for changes...") + "\n" + stoppedText + "\n")
				}
			case err := <-w.watcher.Error:
				prompt.Error(err.Error())
			case <-sigs:
				w.closeWatch()
			case <-w.watcher.Closed:
				return
			}
		}
	}()

	if err := w.watcher.AddRecursive(formulaPath); err != nil {
		log.Fatalln(err)
	}

	w.build(formulaPath, workspace)

	watchText := fmt.Sprintf("Watching dir %s", formulaPath)
	fmt.Println(prompt.Bold(watchText) + "\n" + stoppedText + "\n")

	if err := w.watcher.Start(time.Second * 2); err != nil {
		log.Fatalln(err)
	}
}

func (w WatchManager) build(formulaPath string, workspace formula.Workspace) {
	buildInfo := prompt.Bold("Building formula...")
	s := spinner.StartNew(buildInfo)
	time.Sleep(2 * time.Second)

	info := formula.BuildInfo{FormulaPath: formulaPath, Workspace: workspace}
	if err := w.formula.Build(info); err != nil {
		errorMsg := prompt.Red(err.Error())
		s.Error(errors.New(errorMsg))
		return
	}

	success := prompt.Green("Build completed!")
	s.Success(success)
	prompt.Info("Now you can run your formula with Ritchie!")
}
