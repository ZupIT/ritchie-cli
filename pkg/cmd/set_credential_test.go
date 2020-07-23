package cmd

import (
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/credential/credsingle"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

func Test_setCredentialCmd_runPrompt(t *testing.T) {
	type fields struct {
		Setter         credential.Setter
		SingleSettings credential.SingleSettings
		InputText      prompt.InputText
		InputBool      prompt.InputBool
		InputList      prompt.InputList
		InputPassword  prompt.InputPassword
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Run with success",
			fields: fields{
				Setter:         credSetterMock{},
				SingleSettings: singleCredSettingsMock{},
				InputText:      inputSecretMock{},
				InputBool:      inputFalseMock{},
				InputList:      inputListCredMock{},
				InputPassword:  inputPasswordMock{},
			},
			wantErr: false,
		},
		{
			name: "Run with success AddNew",
			fields: fields{
				Setter:         credSetterMock{},
				SingleSettings: singleCredSettingsMock{},
				InputText:      inputSecretMock{},
				InputBool:      inputFalseMock{},
				InputList:      inputListCustomMock{credsingle.AddNew},
				InputPassword:  inputPasswordMock{},
			},
			wantErr: false,
		},
		{
			name: "Fail when list return err",
			fields: fields{
				Setter:         credSetterMock{},
				SingleSettings: singleCredSettingsMock{},
				InputText:      inputSecretMock{},
				InputBool:      inputFalseMock{},
				InputList:      inputListErrorMock{},
				InputPassword:  inputPasswordMock{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewSetCredentialCmd(
				tt.fields.Setter,
				tt.fields.SingleSettings,
				tt.fields.InputText,
				tt.fields.InputBool,
				tt.fields.InputList,
				tt.fields.InputPassword,
				TutorialFinderMock{},
			)
			o.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := o.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("setCredentialCmd_runPrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
