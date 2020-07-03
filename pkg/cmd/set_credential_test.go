package cmd

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

func TestNewSingleSetCredentialCmd(t *testing.T) {
	cmd := NewSingleSetCredentialCmd(credSetterMock{}, inputSecretMock{}, inputFalseMock{}, inputListCredMock{}, inputPasswordMock{}, InputMultilineMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewSingleSetCredentialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewTeamSetCredentialCmd(t *testing.T) {
	cmd := NewTeamSetCredentialCmd(credSetterMock{}, credSettingsMock{}, inputSecretMock{}, inputFalseMock{}, inputListCredMock{}, inputPasswordMock{}, InputMultilineMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewTeamSetCredentialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewSingleSetCredentialCmdWithEntryFile(t *testing.T) {
	errEntry := errors.New("some error of entry")

	tmpfile := createTemporaryFile()
	defer os.Remove(tmpfile.Name())

	type editableFields struct {
		inputText prompt.InputText
		inputList prompt.InputList
	}
	tests := []struct {
		name           string
		editableFields editableFields
		wantErr        bool
		wantedError    error
	}{
		{
			name: "run set_credential with success when prompt entry selected",
			editableFields: editableFields{
				inputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == MsgTypeCredentialInPrompt {
							return "teste=teste", nil
						}
						return "some_input", nil
					},
				},
				inputList: inputListCustomMock{
					list: func(name string, list []string) (string, error) {
						if name == MsgTypeEntry {
							return EntriesTypeCredentialPrompt, nil
						}
						return "some_input", nil
					},
				},
			},
			wantErr:     false,
			wantedError: nil,
		},
		{
			name: "run set_credential with error in entry credential when prompt entry selected",
			editableFields: editableFields{
				inputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == MsgTypeCredentialInPrompt {
							return "", errEntry
						}
						return "some_input", nil
					},
				},
				inputList: inputListCustomMock{
					list: func(name string, list []string) (string, error) {
						if name == MsgTypeEntry {
							return EntriesTypeCredentialPrompt, nil
						}
						return "some_input", nil
					},
				},
			},
			wantErr:     true,
			wantedError: errEntry,
		},
		{
			name: "run set_credential with success when file entry selected",
			editableFields: editableFields{
				inputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == MsgTypeEntryPath {
							return tmpfile.Name(), nil
						}
						return "some_input", nil
					},
				},
				inputList: inputListCustomMock{
					list: func(name string, list []string) (string, error) {
						if name == MsgTypeEntry {
							return EntriesTypeCredentialFile, nil
						}
						return "some_input", nil
					},
				},
			},
			wantErr:     false,
			wantedError: nil,
		},
		{
			name: "run set_credential with error in entry path when file entry selected",
			editableFields: editableFields{
				inputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == MsgTypeEntryPath {
							return "", errEntry
						}
						return "some_input", nil
					},
				},
				inputList: inputListCustomMock{
					list: func(name string, list []string) (string, error) {
						if name == MsgTypeEntry {
							return EntriesTypeCredentialFile, nil
						}
						return "some_input", nil
					},
				},
			},
			wantErr:     true,
			wantedError: errEntry,
		},
		{
			name: "run set_credential with error in select type entry when file entry selected",
			editableFields: editableFields{
				inputText: inputTextMock{},
				inputList: inputListCustomMock{
					list: func(name string, list []string) (string, error) {
						if name == MsgTypeEntry {
							return "", errEntry
						}
						return "some_input", nil
					},
				},
			},
			wantErr:     true,
			wantedError: errEntry,
		},
		{
			name: "run set_credential with error in key of credential when file entry selected",
			editableFields: editableFields{
				inputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == MsgTypeKeyCredential {
							return "", errEntry
						}
						return "some_input", nil
					},
				},
				inputList: inputListCustomMock{
					list: func(name string, list []string) (string, error) {
						if name == MsgTypeEntry {
							return EntriesTypeCredentialFile, nil
						}
						return "some_input", nil
					},
				},
			},
			wantErr:     true,
			wantedError: errEntry,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewSingleSetCredentialCmd(
				credSetterMock{},
				tt.editableFields.inputText,
				inputFalseMock{},
				tt.editableFields.inputList,
				inputPasswordMock{},
			)
			cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			if cmd == nil {
				t.Errorf("NewTeamSetCredentialCmd got %v", cmd)
			}

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("%s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func createTemporaryFile() *os.File {
	content := []byte("temporary file's content")
	tmpfile, err := ioutil.TempFile("", "example")

	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpfile
}
