package cmd

import (
	"errors"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/security/otp"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/server"
)

func TestNewLoginCmd(t *testing.T) {
	cmd := NewLoginCmd(inputTextMock{}, inputPasswordMock{}, loginManagerMock{}, repoLoaderMock{}, findSetterServerMock{}, otpResolverMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewLoginCmd got %v", cmd)

	}
	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func Test_loginCmd_runPrompt(t *testing.T) {
	type fields struct {
		LoginManager  security.LoginManager
		Loader        formula.RepoLoader
		InputText     prompt.InputText
		InputPassword prompt.InputPassword
		Finder        server.Finder
		Resolver      otp.Resolver
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "run with success",
			fields: fields{
				LoginManager:  loginManagerMock{},
				Loader:        repoLoaderMock{},
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				Finder: findSetterServerMock{},
				Resolver:      otpResolverMock{},
			},
			wantErr: false,
		},
		{
			name: "request otp returns error",
			fields: fields{
				LoginManager:  loginManagerMock{},
				Loader:        repoLoaderMock{},
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				Finder: findSetterServerMock{},
				Resolver:      otpResolverCustomMock{
					requestOtp: func(url, organization string) (otp.Response, error) {
						return otp.Response{}, errors.New("some error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "run with success when ask otp",
			fields: fields{
				LoginManager:  loginManagerMock{},
				Loader:        repoLoaderMock{},
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				Finder: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{}, nil
					},
				},
				Resolver:      otpResolverMock{},
			},
			wantErr: false,
		},
		{
			name: "return err when find return err",
			fields: fields{
				LoginManager:  loginManagerMock{},
				Loader:        repoLoaderMock{},
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				Finder: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{}, errors.New("some error")
					},
				},
				Resolver:      otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "return err when MsgUsername return err",
			fields: fields{
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				InputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == MsgUsername {
							return "", errors.New("some error")
						} else {
							return "some_input", nil
						}
					},
				},
				InputPassword: inputPasswordMock{},
				Finder:        findSetterServerMock{},
				Resolver:      otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "return err when MsgPassword return err",
			fields: fields{
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				InputText:    inputTextMock{},
				InputPassword: inputPasswordCustomMock{
					password: func(label string) (string, error) {
						if label == MsgPassword {
							return "", errors.New("some error")
						} else {
							return "some_input", nil
						}
					},
				},
				Finder: findSetterServerMock{},
				Resolver:      otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "return err when MsgOtp return err",
			fields: fields{
				LoginManager: loginManagerMock{},
				Loader:       repoLoaderMock{},
				InputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == MsgOtp {
							return "", errors.New("some error")
						} else {
							return "some_input", nil
						}
					},
				},
				InputPassword: inputPasswordMock{},
				Finder: findSetterServerCustomMock{
					find: func() (server.Config, error) {
						return server.Config{}, nil
					},
				},
				Resolver:      otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "return err when login return err",
			fields: fields{
				LoginManager: loginManagerCustomMock{
					login: func(user security.User) error {
						return errors.New("some error")
					},
				},
				Loader:        repoLoaderMock{},
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				Finder:        findSetterServerMock{},
				Resolver:      otpResolverMock{},
			},
			wantErr: true,
		},
		{
			name: "return err when load return err",
			fields: fields{
				LoginManager: loginManagerMock{},
				Loader: repoLoaderCustomMock{
					load: func() error {
						return errors.New("some error")
					},
				},
				InputText:     inputTextMock{},
				InputPassword: inputPasswordMock{},
				Finder:        findSetterServerMock{},
				Resolver:      otpResolverMock{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoginCmd(
				tt.fields.InputText,
				tt.fields.InputPassword,
				tt.fields.LoginManager,
				tt.fields.Loader,
				tt.fields.Finder,
				tt.fields.Resolver,
			)
			l.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := l.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("login_runPrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
