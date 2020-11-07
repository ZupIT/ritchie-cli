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

package credential

// Info represents a credential information of the user.
type Detail struct {
	Username   string     `json:"username"`
	Credential Credential `json:"credential"`
	Service    string     `json:"service"`
	Type       Type       `json:"type"`
}

type Type string

func (t Type) String() string {
	return string(t)
}

// A Credential represents the key-value pairs for the Service (User/Pass, Github, Jenkins, etc).
type Credential map[string]string

// Field represents a credential field associated with your type.
type Field struct {
	Name string `json:"field"`
	Type string `json:"type"`
}

type ListCredDatas []ListCredData

type ListCredData struct {
	Provider   string
	Credential string
	Context    string
}

// Fields are used to represents providers.json
type Fields map[string][]Field

type Setter interface {
	Set(d Detail) error
}

type CredFinder interface {
	Find(service string) (Detail, error)
}

type CredDelete interface {
	Delete(service string) error
}

type Reader interface {
	ReadCredentialsFields(path string) (Fields, error)
	ReadCredentialsValue(path string) ([]ListCredData, error)
	ReadCredentialsValueInContext(path string, context string) ([]ListCredData, error)
}

type Writer interface {
	WriteCredentialsFields(fields Fields, path string) error
	WriteDefaultCredentialsFields(path string) error
}

type Pather interface {
	ProviderPath() string
	CredentialsPath() string
}

type ReaderWriterPather interface {
	Reader
	Writer
	Pather
}

type ReaderPather interface {
	Reader
	Pather
}
