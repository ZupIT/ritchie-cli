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

package stream

import (
	"os"
	"testing"
)

func TestDirExists(t *testing.T) {
	fileManager := NewFileManager()
	dirManager := NewDirManager(fileManager)

	ans := dirManager.Exists(os.TempDir())
	if ans == false {
		t.Errorf("DirManager Exists got false")
	}
}

func TestIsDir(t *testing.T) {
	fileManager := NewFileManager()
	dirManager := NewDirManager(fileManager)

	ans := dirManager.IsDir(os.TempDir())
	if ans == false {
		t.Errorf("DirManager IsDir got false")
	}
}

func TestDirList(t *testing.T) {
	fileManager := NewFileManager()
	dirManager := NewDirManager(fileManager)

	_, err := dirManager.List(os.TempDir(), false)
	if err != nil {
		t.Errorf("DirManager List got %s", err)
	}
}
