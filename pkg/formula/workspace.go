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

package formula

const (
	DefaultWorkspaceName = "Default"
	DefaultWorkspaceDir  = "ritchie-formulas-local"
	WorkspacesFile       = "formula_workspaces.json"
)

type (
	Workspace struct {
		Name string `json:"name"`
		Dir  string `json:"dir"`
	}

	Workspaces map[string]string
)

type WorkspaceAdder interface {
	Add(workspace Workspace) error
}

type WorkspaceLister interface {
	List() (Workspaces, error)
}

type WorkspaceDeleter interface {
	Delete(workspace Workspace) error
}

type WorkspaceListDeleter interface {
	WorkspaceDeleter
	WorkspaceLister
}

type WorkspaceAddLister interface {
	WorkspaceAdder
	WorkspaceLister
}
