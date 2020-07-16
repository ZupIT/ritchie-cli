package runner

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestInputManager_Inputs(t *testing.T) {

	inputJson := `
[
  {
	"name" : "sample_text",
	"type" : "text",
	"label" : "Type : ",
	"cache" : {
	  "active": true,
	  "qtd" : 6,
	  "newLabel" : "Type new value. "
	}
  },
  {
	"name" : "sample_list",
	"type" : "text",
	"default" : "in1",
	"items" : ["in_list1", "in_list2", "in_list3", "in_listN"],
	"label" : "Pick your : "
  },
  {
	"name" : "sample_bool",
	"type" : "bool",
	"default" : "false",
	"items" : ["false", "true"],
	"label" : "Pick: "
  },
  {
	"name" : "test_resolver",
	"type" : "CREDENTIAL_TEST"
  }
]
`
	var inputs []formula.Input
	_ = json.Unmarshal([]byte(inputJson), &inputs)

	setup := formula.Setup{
		Config: formula.Config{
			Inputs: inputs,
		},
		FormulaPath: os.TempDir(),
	}

	type in struct {
		iText  inputMock
		iList  inputMock
		iBool  inputMock
		iPass  inputMock
		inType api.TermInputType
		stdin  string
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success stdin",
			in: in{
				iText:  inputMock{text: formula.DefaultCacheNewLabel},
				iList:  inputMock{text: "test"},
				iBool:  inputMock{boolean: false},
				iPass:  inputMock{text: "******"},
				inType: api.Stdin,
				stdin:  `{"sample_text":"test_text","sample_list":"test_list","sample_bool": false}`,
			},
			want: nil,
		},
		{
			name: "success prompt",
			in: in{
				iText:  inputMock{text: formula.DefaultCacheNewLabel},
				iList:  inputMock{text: "test"},
				iBool:  inputMock{boolean: false},
				iPass:  inputMock{text: "******"},
				inType: api.Prompt,
			},
			want: nil,
		},
		{
			name: "error unknown prompt",
			in: in{
				iText:  inputMock{text: formula.DefaultCacheNewLabel},
				iList:  inputMock{text: "test"},
				iBool:  inputMock{boolean: false},
				iPass:  inputMock{text: "******"},
				inType: api.TermInputType(3),
			},
			want: ErrInputNotRecognized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolvers := env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}}
			iText := tt.in.iText
			iList := tt.in.iList
			iBool := tt.in.iBool
			iPass := tt.in.iPass

			inputManager := NewInputManager(resolvers, iList, iText, iBool, iPass)

			cmd := &exec.Cmd{}
			if tt.in.inType == api.Stdin {
				cmd.Stdin = strings.NewReader(tt.in.stdin)
			}

			got := inputManager.Inputs(cmd, setup, tt.in.inType)

			if (tt.want != nil && got == nil) || got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Inputs(%s) got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
