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
	cmd := NewSingleSetCredentialCmd(credSetterMock{}, inputSecretMock{}, inputFalseMock{}, inputListCredMock{}, inputPasswordMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewSingleSetCredentialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewTeamSetCredentialCmd(t *testing.T) {
	cmd := NewTeamSetCredentialCmd(credSetterMock{}, credSettingsMock{}, inputSecretMock{}, inputFalseMock{}, inputListCredMock{}, inputPasswordMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewTeamSetCredentialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewSingleSetCredentialCmdWithEntryArchive(t *testing.T) {
	errPath := errors.New("some error of path")
	cmd := NewSingleSetCredentialCmd(
		credSetterMock{},
		inputTextCustomMock{
			text: func(name string, required bool) (string, error) {
				if name == MsgTypeEntryPath {
					return "/test", errPath
				}
				return "some_path", nil
			},
		},
		inputFalseMock{},
		inputListCustomMock{
			list: func(name string, list []string) (string, error) {
				if name == MsgTypeEntry {
					return EntriesTypeCredentialFile, nil
				}
				return "some_input", nil
			},
		},
		inputPasswordMock{},
	)
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewTeamSetCredentialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != errPath {
		t.Errorf("%s = %q, want %q", cmd.Use, err, errPath)
	}
}

func TestNewSingleSetCredentialCmdWithEntryArchiveWithError1(t *testing.T) {
	errCredential := errors.New("some error of entry")
	cmd := NewSingleSetCredentialCmd(
		credSetterMock{},
		inputTextCustomMock{
			text: func(name string, required bool) (string, error) {
				if name == MsgTypeCredentialInPrompt {
					return "", errCredential
				}
				return "some_path", nil
			},
		},
		inputFalseMock{},
		inputListCustomMock{
			list: func(name string, list []string) (string, error) {
				if name == MsgTypeEntry {
					return "", errCredential
				}
				return "some_input", nil
			},
		},
		inputPasswordMock{},
	)
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewTeamSetCredentialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != errCredential {
		t.Errorf("%s = %q, want %q", cmd.Use, err, errCredential)
	}
}

func TestNewSingleSetCredentialCmdWithEntryArchiveWithError2(t *testing.T) {
	errEntry := errors.New("some error of entry")
	cmd := NewSingleSetCredentialCmd(
		credSetterMock{},
		inputTextCustomMock{
			text: func(name string, required bool) (string, error) {
				if name == EntriesTypeCredentialPrompt {
					return "", errEntry
				}
				return "some_path", nil
			},
		},
		inputFalseMock{},
		inputListCustomMock{
			list: func(name string, list []string) (string, error) {
				if name == MsgTypeEntry {
					return EntriesTypeCredentialPrompt, nil
				}
				return "some_input", nil
			},
		},
		inputPasswordMock{},
	)
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewTeamSetCredentialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != errEntry {
		t.Errorf("%s = %q, want %q", cmd.Use, err, errEntry)
	}
}

func TestNewSingleSetCredentialCmdWithEntryArchiveWithOptions(t *testing.T) {
	errEntry := errors.New("some error of entry")
	errReader := errors.New("no such file or directory")
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
			name: "run prompt with success",
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
			name: "run archive with success",
			editableFields: editableFields{
				inputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == MsgTypeEntryPath {
							return tmpfile.Name(), nil
						}
						return "some_input=some_input", nil
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
			name: "run archive with error in type msg type entry",
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
			name: "run archive with error in type key file entry",
			editableFields: editableFields{
				inputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == "Type key of your credential: " {
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
			name: "run archive with error in type value file entry",
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
			wantedError: errReader,
		},
		{
			name: "run archive with error in prompt entry credential",
			editableFields: editableFields{
				inputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == MsgTypeCredentialInPrompt {
							return "", errEntry
						}
						return "some_input", nil
					},
				},
				inputList: inputListMock{},
			},
			wantErr:     true,
			wantedError: errEntry,
		},
		{
			name: "run archive with error in prompt format entry credential",
			editableFields: editableFields{
				inputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						if name == MsgTypeCredentialInPrompt {
							return "some_input=", nil
						}
						return "some_input", nil
					},
				},
				inputList: inputListMock{},
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

			if tt.wantErr && err != tt.wantedError {
				t.Errorf("%s error = %v, wantedError %v", tt.name, err, tt.wantedError)
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
