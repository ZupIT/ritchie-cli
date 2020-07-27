package modifier

import (
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestFormulaCmd_modify(t *testing.T) {
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
				b: []byte(`$ #rit-replace{formulaCmd}`),
			},
			want: []byte(`$ rit testing formula`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FormulaCmd{
				cf: tt.fields.cf,
			}
			if got := f.modify(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("modify() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
