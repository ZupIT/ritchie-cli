package cmd

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/security/otp"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/server"
)

func TestNewSingleInitCmd(t *testing.T) {
	cmd := NewSingleInitCmd(inputPasswordMock{}, passphraseManagerMock{}, repoLoaderMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")

	if cmd == nil {
		t.Errorf("NewSingleInitCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func Test_initTeamCmd_runPrompt(t *testing.T) {
	type fields struct {
		InputText     prompt.InputText
		InputPassword prompt.InputPassword
		InputURL      prompt.InputURL
		InputBool     prompt.InputBool
		FindSetter    server.FindSetter
		LoginManager  security.LoginManager
		Loader        formula.RepoLoader
		Resolver      otp.Resolver
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Run With Success with empty config",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputFalseMock{},
				FindSetter:    findSetterServerMock{},
				LoginManager:  loginManagerMock{},
				Loader:        repoLoaderMock{},
				Resolver:      otpResolverMock{},
			},
			wantErr: false,
		},
		{
			name: "otp request returns error",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputTrueMock{},
				FindSetter:    findSetterServerMock{},
				LoginManager:  loginManagerMock{},
				Loader:        repoLoaderMock{},
				Resolver:      otpResolverCustomMock{
					requestOtp: func(url, organization string) (otp.Response, error) {
						return otp.Response{}, errors.New("some error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Return err when finder return err",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputFalseMock{},
				FindSetter: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{}, errors.New("some error")
					},
				},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "Return err when setter return err",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputFalseMock{},
				FindSetter: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{}, nil
					},
					set: func(config *server.Config) error {
						return errors.New("some error")
					},
				},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "Run with success when OrganizationAlreadyExists and override false",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputFalseMock{},
				FindSetter: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{
							Organization: "any_value",
						}, nil
					},
					set: func(config *server.Config) error {
						return nil
					},
				},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: false,
		},
		{
			name: "Return err when OrganizationAlreadyExists and MsgOrganization return err",
			fields: fields{
				InputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == MsgOrganization {
							return "", errors.New("some error")
						} else {
							return "some_input", nil
						}
					},
					textWithValidate: func(name string, validate func(interface{}) error) (string, error) {
						return "", nil
					},
				},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputTrueMock{},
				FindSetter: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{
							Organization: "any_value",
						}, nil
					},
					set: func(config *server.Config) error {
						return nil
					},
				},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "Return err when OrganizationAlreadyExists return err",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool: inputBoolCustomMock{
					func(name string, items []string) (bool, error) {
						if name == fmt.Sprintf(msgOrganizationAlreadyExists, "any_value") {
							return false, errors.New("some error")
						} else {
							return false, nil
						}
					},
				},
				FindSetter: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{
							Organization: "any_value",
						}, nil
					},
					set: func(config *server.Config) error {
						return nil
					},
				},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "Run with success when OrganizationAlreadyExists and override true",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputTrueMock{},
				FindSetter: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{
							Organization: "any_value",
						}, nil
					},
					set: func(config *server.Config) error {
						return nil
					},
				},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: false,
		},
		{
			name: "Return err when MsgOrganization return err",
			fields: fields{
				InputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == MsgOrganization {
							return "", errors.New("some error")
						} else {
							return "some_input", nil
						}
					},
				},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputFalseMock{},
				FindSetter:    findSetterServerMock{},
				LoginManager:  loginManagerMock{},
				Loader:        repoLoaderMock{},
				Resolver:      otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "Run with success when ServerURLAlreadyExists and override false",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputFalseMock{},
				FindSetter: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{
							Organization: "any_value",
							URL:          "http://someurl.com.br",
						}, nil
					},
					set: func(config *server.Config) error {
						return nil
					},
				},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: false,
		},
		{
			name: "Return err when ServerURLAlreadyExists return err",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool: inputBoolCustomMock{
					bool: func(name string, items []string) (bool, error) {
						if name == fmt.Sprintf(msgServerURLAlreadyExists, "http://someurl.com.br") {
							return false, errors.New("some error")
						} else {
							return false, nil
						}
					},
				},
				FindSetter: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{
							Organization: "any_value",
							URL:          "http://someurl.com.br",
						}, nil
					},
					set: func(config *server.Config) error {
						return nil
					},
				},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "Run with success when ServerURLAlreadyExists and override true",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputTrueMock{},
				FindSetter: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{
							Organization: "any_value",
							URL:          "http://someurl.com.br",
						}, nil
					},
					set: func(config *server.Config) error {
						return nil
					},
				},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: false,
		},
		{
			name: "Run with success when get otp value",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputTrueMock{},
				FindSetter: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{
							Organization: "any_value",
							URL:          "http://someurl.com.br",
						}, nil
					},
					set: func(config *server.Config) error {
						return nil
					},
				},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: false,
		},
		{
			name: "Return err when MsgLogin return err",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool: inputBoolCustomMock{
					bool: func(name string, items []string) (bool, error) {
						if name == MsgLogin {
							return false, errors.New("some error")
						} else {
							return false, nil
						}
					},
				},
				FindSetter:   findSetterServerMock{},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "Return err when MsgUsername return err",
			fields: fields{
				InputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == MsgUsername {
							return "", errors.New("some error")
						} else {
							return "any_input", nil
						}
					},
				},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputTrueMock{},
				FindSetter:    findSetterServerMock{},
				LoginManager:  loginManagerMock{},
				Loader:        repoLoaderMock{},
				Resolver:      otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "Return err when MsgPassword return err",
			fields: fields{
				InputText: inputTextMock{},
				InputPassword: inputPasswordCustomMock{
					password: func(label string) (string, error) {
						if label == MsgPassword {
							return "", errors.New("some error")
						} else {
							return "any_input", nil
						}
					},
				},
				InputURL:     inputURLMock{},
				InputBool:    inputTrueMock{},
				FindSetter:   findSetterServerMock{},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "Return err when MsgOtp return err",
			fields: fields{
				InputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == MsgOtp {
							return "", errors.New("some error")
						} else {
							return "any_input", nil
						}
					},
				},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputTrueMock{},
				FindSetter: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{}, nil
					},
					set: func(config *server.Config) error {
						return nil
					},
				},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "Return err when o.Login return err",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputTrueMock{},
				FindSetter:    findSetterServerMock{},
				LoginManager: loginManagerCustomMock{
					login: func(user security.User) error {
						return errors.New("some error")
					},
				},
				Loader:   repoLoaderMock{},
				Resolver: otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "Return err when o.Load return err",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL:      inputURLMock{},
				InputBool:     inputTrueMock{},
				FindSetter:    findSetterServerMock{},
				LoginManager:  loginManagerMock{},
				Loader: repoLoaderCustomMock{
					load: func() error {
						return errors.New("some errors")
					},
				},
				Resolver: otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "Return err when MsgServerURL return err and IsValidURL return err",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL: inputURLCustomMock{
					url: func(name, defaultValue string) (string, error) {
						if name == MsgServerURL {
							return "", errors.New("some error")
						} else {
							return "http://localhost/mocked", nil
						}
					},
				},
				InputBool:    inputTrueMock{},
				FindSetter:   findSetterServerMock{},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "Return err when MsgServerURL return err and IsValidURL not return err",
			fields: fields{
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				InputURL: inputURLCustomMock{
					url: func(name, defaultValue string) (string, error) {
						if name == MsgServerURL {
							return "", errors.New("some error")
						} else {
							return "http://localhost/mocked", nil
						}
					},
				},
				InputBool: inputTrueMock{},
				FindSetter: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{
							URL: "http://localhost/mocked",
						}, nil
					},
					set: func(config *server.Config) error {
						return errors.New("some error")
					},
				},
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				Resolver:     otpResolverMock{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewTeamInitCmd(
				tt.fields.InputText,
				tt.fields.InputPassword,
				tt.fields.InputURL,
				tt.fields.InputBool,
				tt.fields.FindSetter,
				tt.fields.LoginManager,
				tt.fields.Loader,
				tt.fields.Resolver,
			)
			o.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := o.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("init_runPrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
