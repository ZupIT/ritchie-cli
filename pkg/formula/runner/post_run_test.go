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

package runner

import (
	"errors"
	"testing"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostRun(t *testing.T) {
	type in struct {
		wrErr  error
		rdErr  error
		apErr  error
		mvErr  error
		rmErr  error
		lnErr  error
		exist  bool
		dir    error
		setup  formula.Setup
		docker bool
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				setup:  formula.Setup{},
				docker: false,
			},
			want: nil,
		},
		{
			name: "success docker",
			in: in{
				setup:  formula.Setup{},
				docker: true,
			},
			want: nil,
		},
		{
			name: "error remove .env file docker",
			in: in{
				rmErr:  errors.New("error to remove .env file"),
				setup:  formula.Setup{},
				docker: true,
			},
			want: errors.New("error to remove .env file"),
		},
		{
			name: "error list new files",
			in: in{
				lnErr:  errors.New("error to list new files"),
				setup:  formula.Setup{},
				docker: false,
			},
			want: errors.New("error to list new files"),
		},
		{
			name: "error move new files",
			in: in{
				mvErr:  errors.New("error to move new files"),
				setup:  formula.Setup{},
				docker: false,
			},
			want: errors.New("error to move new files"),
		},
		{
			name: "error remove work dir",
			in: in{
				dir:    errors.New("error to remove workdir"),
				setup:  formula.Setup{},
				docker: false,
			},
			want: nil,
		},
		{
			name: "input deprecated",
			in: in{
				setup: formula.Setup{
					Config: formula.Config{
						Inputs: formula.Inputs{
							formula.Input{
								Type: "dynamic",
							},
						},
					},
				},
				docker: false,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dm := &mocks.DirManagerMock{}
			dm.On("Remove", mock.Anything).Return(tt.in.dir)
			fm := &mocks.FileManagerMock{}
			fm.On("Write", mock.Anything, mock.Anything).Return(tt.in.wrErr)
			fm.On("Read", mock.Anything).Return([]byte{}, tt.in.rdErr)
			fm.On("Exists", mock.Anything).Return(tt.in.exist)
			fm.On("Append", mock.Anything, mock.Anything).Return(tt.in.apErr)
			fm.On("Move", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.mvErr)
			fm.On("Remove", mock.Anything).Return(tt.in.rmErr)
			fm.On("ListNews", mock.Anything, mock.Anything).Return([]string{}, tt.in.lnErr)
			runner := NewPostRunner(fm, dm)
			got := runner.PostRun(tt.in.setup, tt.in.docker)

			assert.Equal(t, tt.want, got)
		})
	}

}
