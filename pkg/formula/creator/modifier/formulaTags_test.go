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

package modifier

import (
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestFormulaTags_modify(t *testing.T) {
	type fields struct {
		cf formula.Create
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "modify with success",
			fields: fields{
				cf: formula.Create{
					FormulaCmd: "rit testing formula",
				},
			},
			args: args{
				b: []byte(`{"tags": ["#rit-replace{formulaTags}"]}`),
			},
			want: []byte(`{"tags": ["testing", "formula"]}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FormulaTags{
				cf: tt.fields.cf,
			}
			if got := f.modify(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("modify() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
