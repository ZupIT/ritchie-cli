package cmd

import (
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

func Test_setCredentialCmd_runPrompt(t *testing.T) {
	type in struct {
		Setter        credential.Setter
		credFile      credential.ReaderWriterPather
		file          stream.FileReadExister
		InputText     prompt.InputText
		InputBool     prompt.InputBool
		InputList     prompt.InputList
		InputPassword prompt.InputPassword
	}
	var tests = []struct {
		name    string
		in      in
		wantErr bool
		want    string
	}{
		{
			name: "Run with success",
			in: in{
				Setter:   credSetterMock{},
				credFile: credSettingsMock{},
				file: sMocks.FileReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
					ReadMock: func(path string) ([]byte, error) {
						return nil, nil
					},
				},
				InputText:     inputSecretMock{},
				InputBool:     inputFalseMock{},
				InputList:     inputListCredMock{},
				InputPassword: inputPasswordMock{},
			},
			wantErr: false,
		},
		{
			name: "Run with success AddNew",
			in: in{
				Setter:        credSetterMock{},
				credFile:      credSettingsMock{},
				InputText:     inputSecretMock{},
				InputBool:     inputFalseMock{},
				InputList:     inputListCustomMock{credential.AddNew},
				InputPassword: inputPasswordMock{},
			},
			wantErr: false,
		},
		{
			name: "Fail when list return err",
			in: in{
				Setter:        credSetterMock{},
				credFile:      credSettingsMock{},
				InputText:     inputSecretMock{},
				InputBool:     inputFalseMock{},
				InputList:     inputListErrorMock{},
				InputPassword: inputPasswordMock{},
			},
			wantErr: true,
		},
		{
			name: "Fail when text return err",
			in: in{
				Setter:        credSetterMock{},
				credFile:      credSettingsMock{},
				InputText:     inputTextErrorMock{},
				InputBool:     inputFalseMock{},
				InputList:     inputListCustomMock{credential.AddNew},
				InputPassword: inputPasswordMock{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewSetCredentialCmd(
				tt.in.Setter,
				tt.in.credFile,
				tt.in.file,
				tt.in.InputText,
				tt.in.InputBool,
				tt.in.InputList,
				tt.in.InputPassword,
			)
			o.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := o.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("set credential command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
