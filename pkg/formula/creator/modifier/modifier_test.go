package modifier

import (
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestNewModifiers(t *testing.T) {
	type args struct {
		create formula.Create
	}

	tests := []struct {
		name string
		args args
		in   string
		want string
	}{
		{
			name: "modify with success",
			args: args{
				create: formula.Create{
					FormulaCmd: "rit testing formula",
				},
			},
			in:   `tags: #rit-replace{formulaTags} cmd: #rit-replace{formulaCmd}`,
			want: `tags: "testing", "formula" cmd: rit testing formula`,
		},
		{
			name: "not modify",
			args: args{
				create: formula.Create{
					FormulaCmd: "rit testing formula",
				},
			},
			in:   `some test`,
			want: `some test`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModifiers(tt.args.create)
			got := Modify([]byte(tt.in), m)
			if !reflect.DeepEqual(got, []byte(tt.want)) {
				t.Errorf("\nModify() =\n%v\nwant:\n%v", string(got), tt.want)
			}
		})
	}
}
