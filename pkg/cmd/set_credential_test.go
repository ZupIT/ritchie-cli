package cmd

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/credential/credsingle"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestNewSingleSetCredentialCmd(t *testing.T) {
	cmd := NewSingleSetCredentialCmd(
		credSetterMock{},
		singleCredSettingsMock{},
		inputSecretMock{},
		inputFalseMock{},
		inputListCredMock{},
		inputPasswordMock{},
		FileManagerMock{},
	)

	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewSingleSetCredentialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewTeamSetCredentialCmd(t *testing.T) {
	cmd := NewTeamSetCredentialCmd(credSetterMock{},
		credSettingsMock{},
		inputSecretMock{},
		inputFalseMock{},
		inputListCredMock{},
		inputPasswordMock{},
		InputMultilineMock{},
	)
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewTeamSetCredentialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func Test_setCredentialCmd_promptResolver(t *testing.T) {
	type fields struct {
		Setter          credential.Setter
		Settings        credential.Settings
		SingleSettings  credential.SingleSettings
		edition         api.Edition
		InputText       prompt.InputText
		InputBool       prompt.InputBool
		InputList       prompt.InputList
		InputPassword   prompt.InputPassword
		InputMultiline  prompt.InputMultiline
		FileReadExister stream.FileReadExister
	}
	tests := []struct {
		name    string
		fields  fields
		want    credential.Detail
		wantErr bool
	}{
		{
			name: "reach default",
			fields: fields{

			},
			want:    credential.Detail{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setCredentialCmd{
				Setter:          tt.fields.Setter,
				Settings:        tt.fields.Settings,
				SingleSettings:  tt.fields.SingleSettings,
				edition:         tt.fields.edition,
				InputText:       tt.fields.InputText,
				InputBool:       tt.fields.InputBool,
				InputList:       tt.fields.InputList,
				InputPassword:   tt.fields.InputPassword,
				InputMultiline:  tt.fields.InputMultiline,
				FileReadExister: tt.fields.FileReadExister,
			}
			got, err := s.promptResolver()
			if (err != nil) != tt.wantErr {
				t.Errorf("promptResolver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("promptResolver() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setCredentialCmd_singlePrompt(t *testing.T) {
	type fields struct {
		Setter          credential.Setter
		Settings        credential.Settings
		SingleSettings  credential.SingleSettings
		edition         api.Edition
		InputText       prompt.InputText
		InputBool       prompt.InputBool
		InputList       prompt.InputList
		InputPassword   prompt.InputPassword
		InputMultiline  prompt.InputMultiline
		FileReadExister stream.FileReadExister
	}
	tests := []struct {
		name    string
		fields  fields
		want    credential.Detail
		wantErr bool
	}{
		{
			name: "error on write default credentials",
			fields: fields{
				SingleSettings: singleCredSettingsCustomMock{
					writeDefaultCredentials: func(path string) error {
						return errors.New("some error")
					},
				},
			},
			want:    credential.Detail{},
			wantErr: true,
		},
		{
			name: "error on read credentials",
			fields: fields{
				SingleSettings: singleCredSettingsCustomMock{
					writeDefaultCredentials: func(path string) error {
						return nil
					},
					readCredentials: func(path string) (credential.Fields, error) {
						return nil, errors.New("some error")
					},
				},
			},
			want:    credential.Detail{},
			wantErr: true,
		},
		{
			name: "error on provider choose",
			fields: fields{
				SingleSettings: singleCredSettingsMock{},
				InputList:      inputListErrorMock{},
			},
			want:    credential.Detail{},
			wantErr: true,
		},
		{
			name: "error on input text",
			fields: fields{
				SingleSettings: singleCredSettingsMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return credsingle.AddNew, nil
					},
				},
				InputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						return "", errors.New("some error")
					},
				},
			},
			want:    credential.Detail{},
			wantErr: true,
		},
		{
			name: "error on input bool",
			fields: fields{
				SingleSettings: singleCredSettingsMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return credsingle.AddNew, nil
					},
				},
				InputText: inputTextMock{},
				InputBool: inputBoolErrorMock{},
			},
			want:    credential.Detail{},
			wantErr: true,
		},
		{
			name: "error on write credentials",
			fields: fields{
				SingleSettings: singleCredSettingsCustomMock{
					writeDefaultCredentials: func(path string) error {
						return nil
					},
					readCredentials: func(path string) (credential.Fields, error) {
						return credential.Fields{}, nil
					},
					writeCredentials: func(fields credential.Fields, path string) error {
						return errors.New("some error")
					},
				},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return credsingle.AddNew, nil
					},
				},
				InputText: inputTextMock{},
				InputBool: inputFalseMock{},
			},
			want:    credential.Detail{},
			wantErr: true,
		},
		{
			name: "success with empty inputs",
			fields: fields{
				SingleSettings: singleCredSettingsCustomMock{
					writeDefaultCredentials: func(path string) error {
						return nil
					},
					readCredentials: func(path string) (credential.Fields, error) {
						return credential.Fields{}, nil
					},
					writeCredentials: func(fields credential.Fields, path string) error {
						return nil
					},
				},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return credsingle.AddNew, nil
					},
				},
				InputText: inputTextMock{},
				InputBool: inputFalseMock{},
			},
			want: credential.Detail{
				Username: "",
				Credential: credential.Credential{
					"mocked text": "mocked text",
				},
				Service: "mocked text",
				Type:    "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setCredentialCmd{
				Setter:          tt.fields.Setter,
				Settings:        tt.fields.Settings,
				SingleSettings:  tt.fields.SingleSettings,
				edition:         tt.fields.edition,
				InputText:       tt.fields.InputText,
				InputBool:       tt.fields.InputBool,
				InputList:       tt.fields.InputList,
				InputPassword:   tt.fields.InputPassword,
				InputMultiline:  tt.fields.InputMultiline,
				FileReadExister: tt.fields.FileReadExister,
			}
			got, err := s.singlePrompt()
			if (err != nil) != tt.wantErr {
				t.Errorf("singlePrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("singlePrompt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setCredentialCmd_inputFile(t *testing.T) {
	type fields struct {
		Setter          credential.Setter
		Settings        credential.Settings
		SingleSettings  credential.SingleSettings
		edition         api.Edition
		InputText       prompt.InputText
		InputBool       prompt.InputBool
		InputList       prompt.InputList
		InputPassword   prompt.InputPassword
		InputMultiline  prompt.InputMultiline
		FileReadExister stream.FileReadExister
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "run with success",
			fields: fields{
				InputText: inputTextMock{},
				FileReadExister: FileManagerMock{},
			},
			want:    "Some response",
			wantErr: false,
		},
		{
			name:    "error on input text",
			fields:  fields{
				InputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						return "", errors.New("some error")
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:    "error on read file",
			fields:  fields{
				InputText: inputTextMock{},
				FileReadExister: FileManagerCustomMock{
					read: func(path string) ([]byte, error) {
						return nil, errors.New("some error")
					},
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setCredentialCmd{
				Setter:          tt.fields.Setter,
				Settings:        tt.fields.Settings,
				SingleSettings:  tt.fields.SingleSettings,
				edition:         tt.fields.edition,
				InputText:       tt.fields.InputText,
				InputBool:       tt.fields.InputBool,
				InputList:       tt.fields.InputList,
				InputPassword:   tt.fields.InputPassword,
				InputMultiline:  tt.fields.InputMultiline,
				FileReadExister: tt.fields.FileReadExister,
			}
			got, err := s.inputFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("inputFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("inputFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setCredentialCmd_stdinResolver(t *testing.T) {
	type fields struct {
		Setter          credential.Setter
		Settings        credential.Settings
		SingleSettings  credential.SingleSettings
		edition         api.Edition
		InputText       prompt.InputText
		InputBool       prompt.InputBool
		InputList       prompt.InputList
		InputPassword   prompt.InputPassword
		InputMultiline  prompt.InputMultiline
		FileReadExister stream.FileReadExister
	}
	tests := []struct {
		name    string
		fields  fields
		want    credential.Detail
		wantErr bool
	}{
		{
			name:    "run with empty edition",
			fields:  fields{
				edition: "",
			},
			want:    credential.Detail{},
			wantErr: true,
		},
		{
			name:    "error on stdin inputs",
			fields:  fields{
				edition: api.Team,
			},
			want:    credential.Detail{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setCredentialCmd{
				Setter:          tt.fields.Setter,
				Settings:        tt.fields.Settings,
				SingleSettings:  tt.fields.SingleSettings,
				edition:         tt.fields.edition,
				InputText:       tt.fields.InputText,
				InputBool:       tt.fields.InputBool,
				InputList:       tt.fields.InputList,
				InputPassword:   tt.fields.InputPassword,
				InputMultiline:  tt.fields.InputMultiline,
				FileReadExister: tt.fields.FileReadExister,
			}
			got, err := s.stdinResolver()
			if (err != nil) != tt.wantErr {
				t.Errorf("stdinResolver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stdinResolver() got = %v, want %v", got, tt.want)
			}
		})
	}
}