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

package prompt

type InputText interface {
	Text(name string, required bool, helper ...string) (string, error)
}

type InputTextValidator interface {
	Text(name string, validate func(interface{}) error, helper ...string) (string, error)
}

type InputBool interface {
	Bool(name string, items []string, helper ...string) (bool, error)
}

type InputPassword interface {
	Password(label string, helper ...string) (string, error)
}

type InputMultiline interface {
	MultiLineText(name string, required bool) (string, error)
}

type InputList interface {
	List(name string, items []string, helper ...string) (string, error)
}

type InputMultiSelect interface {
	MultiSelect(name string, items []string, helper ...string) ([]string, error)
}

type InputInt interface {
	Int(name string, helper ...string) (int64, error)
}

type InputEmail interface {
	Email(name string) (string, error)
}

type InputURL interface {
	URL(name, defaultValue string) (string, error)
}
