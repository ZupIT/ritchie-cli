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

package rcontext

const (
	ContextPath = "%s/contexts"
	CurrentCtx  = "Current -> "
	DefaultCtx  = "default"
)

type ContextHolder struct {
	Current string   `json:"current_context"`
	All     []string `json:"contexts"`
}

type Setter interface {
	Set(ctx string) (ContextHolder, error)
}

type Finder interface {
	Find() (ContextHolder, error)
}

type Remover interface {
	Remove(ctx string) (ContextHolder, error)
}

type FindRemover interface {
	Finder
	Remover
}

type FindSetter interface {
	Finder
	Setter
}
