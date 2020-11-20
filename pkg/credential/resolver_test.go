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

import (
	"errors"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/env"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestCredentialResolver(t *testing.T) {
	fileManager := stream.NewFileManager()
	tempDirectory := os.TempDir()
	contextFinder := env.NewFinder(tempDirectory, fileManager)
	credentialSetter := NewSetter(tempDirectory, contextFinder)
	credentialFinder := NewFinder(tempDirectory, contextFinder, fileManager)

	defer os.RemoveAll(File(tempDirectory, "", ""))

	type in struct {
		credSetter      Setter
		credFinder      CredFinder
		credentialField string
		pass            prompt.InputPassword
	}

	type out struct {
		want string
		err  error
	}

	var tests = []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "resolve new provider",
			in: in{
				credSetter:      credentialSetter,
				credFinder:      credentialFinder,
				credentialField: "CREDENTIAL_PROVIDER_KEY",
				pass:            passwordMock{value: "key"},
			},
			out: out{
				want: "key",
			},
		},
		{
			name: "resolve new key",
			in: in{
				credSetter:      credentialSetter,
				credFinder:      credentialFinder,
				credentialField: "CREDENTIAL_PROVIDER_KEY2",
				pass:            passwordMock{value: "key2"},
			},
			out: out{
				want: "key2",
			},
		},
		{
			name: "resolve existing key",
			in: in{
				credSetter:      credentialSetter,
				credFinder:      credentialFinder,
				credentialField: "CREDENTIAL_PROVIDER_KEY2",
				pass:            passwordMock{value: "key2"},
			},
			out: out{
				want: "key2",
			},
		},
		{
			name: "error to read password",
			in: in{
				credSetter:      credentialSetter,
				credFinder:      credentialFinder,
				credentialField: "CREDENTIAL_PROVIDER_NEW_ERROR",
				pass:            passwordMock{err: errors.New("error to read password")},
			},
			out: out{
				err: errors.New("error to read password"),
			},
		},
		{
			name: "error to read password",
			in: in{
				credSetter:      setCredMock{err: errors.New("error to set credential")},
				credFinder:      credentialFinder,
				credentialField: "CREDENTIAL_AWS_USER",
				pass:            passwordMock{value: "username"},
			},
			out: out{
				err: errors.New("error to set credential"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentialResolver := NewResolver(tt.in.credFinder, tt.in.credSetter, tt.in.pass)
			got, err := credentialResolver.Resolve(tt.in.credentialField)

			if tt.out.err != nil && tt.out.err.Error() != err.Error() {
				t.Errorf("Resolve(%s) got %v, want %v", tt.name, err, tt.out.err)
			}
			if got != tt.out.want {
				t.Errorf("Resolve(%s) got %v, want %v", tt.name, got, tt.out.want)
			}
		})
	}
}

type passwordMock struct {
	value string
	err   error
}

func (pass passwordMock) Password(string, ...string) (string, error) {
	return pass.value, pass.err
}

type setCredMock struct {
	err error
}

func (s setCredMock) Set(d Detail) error {
	return s.err
}
