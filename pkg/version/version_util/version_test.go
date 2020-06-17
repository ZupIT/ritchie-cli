package version_util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type StubResolverVersions struct {
	getStableVersion func() (string, error)
}

func (r StubResolverVersions) GetStableVersion() (string, error) {
	return r.getStableVersion()
}

type StubFileUtilService struct {
	readFile      func(string) ([]byte, error)
	writeFilePerm func(string, []byte, int32) error
}

func (s StubFileUtilService) ReadFile(path string) ([]byte, error) {

	return s.readFile(path)
}

func (s StubFileUtilService) WriteFilePerm(path string, content []byte, perm int32) error {
	return s.writeFilePerm(path, content, perm)
}

func TestDefaultVersionResolver_GetStableVersion(t *testing.T) {

	// case 1 Should get stableVersion
	expectedResultCase1 := "1.0.0"

	mockHttpCase1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(expectedResultCase1 + "\n"))
	}))

	// case 4 Should not get stableVersion from expired cache
	expectedResultCase4 := "2.5.6"

	mockHttpCase4 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(expectedResultCase4 + "\n"))
	}))

	type fields struct {
		CurrentVersion   string
		StableVersionUrl string
		FileUtilService  fileutil.FileUtilService
		HttpClient       *http.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Should get stableVersion",
			fields: fields{
				CurrentVersion:   "Any value",
				StableVersionUrl: mockHttpCase1.URL,
				FileUtilService: StubFileUtilService{
					readFile: func(s string) ([]byte, error) {
						return []byte{}, errors.New("some error")
					},
					writeFilePerm: func(_ string, _ []byte, _ int32) error { return nil },
				},
				HttpClient: mockHttpCase1.Client(),
			},
			want:    expectedResultCase1,
			wantErr: false,
		},
		{
			name: "Should return err when http.get fail",
			fields: fields{
				CurrentVersion:   "Any value",
				StableVersionUrl: "any value",
				FileUtilService: StubFileUtilService{
					readFile: func(s string) ([]byte, error) {
						return []byte{}, errors.New("some error")
					},
					writeFilePerm: func(_ string, _ []byte, _ int32) error { return nil },
				},
				HttpClient: &http.Client{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Should get stableVersion from cache",
			fields: fields{
				CurrentVersion:   "Any value",
				StableVersionUrl: "any value",
				FileUtilService: StubFileUtilService{
					readFile: func(_ string) ([]byte, error) {
						cache := stableVersionCache{
							"1.4.5",
							time.Now().Add(time.Hour * 1).Unix(),
						}

						return json.Marshal(cache)
					},
					writeFilePerm: func(_ string, _ []byte, _ int32) error { return nil },
				},
				HttpClient: &http.Client{},
			},
			want:    "1.4.5",
			wantErr: false,
		},
		{
			name: "Should not get stableVersion from expired cache",
			fields: fields{
				CurrentVersion:   "Any value",
				StableVersionUrl: mockHttpCase4.URL,
				FileUtilService: StubFileUtilService{
					readFile: func(_ string) ([]byte, error) {
						cache := stableVersionCache{
							"1.5.0",
							time.Now().Add(time.Hour * 1 * -1).Unix(),
						}

						return json.Marshal(cache)
					},
					writeFilePerm: func(_ string, _ []byte, _ int32) error { return nil },
				},
				HttpClient: mockHttpCase4.Client(),
			},
			want:    expectedResultCase4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := DefaultVersionResolver{
				StableVersionUrl: tt.fields.StableVersionUrl,
				FileUtilService:  tt.fields.FileUtilService,
				HttpClient:       tt.fields.HttpClient,
			}
			got, err := r.GetStableVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStableVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetStableVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyNewVersion(t *testing.T) {
	type args struct {
		resolve        Resolver
		currentVersion string
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
	}{
		{
			name: "Should not print warning",
			args: args{
				resolve: StubResolverVersions{
					getStableVersion: func() (string, error) {
						return "1.0.0", nil
					},
				},
				currentVersion: "1.0.0",
			},
			wantWriter: "",
		},
		{
			name: "Should print warning",
			args: args{
				resolve: StubResolverVersions{
					getStableVersion: func() (string, error) {
						return "1.0.1", nil
					},
				},
				currentVersion: "1.0.0",
			},
			wantWriter: fmt.Sprintf(prompt.Yellow, MsgRitUpgrade),
		},
		{
			name: "Should not print on error in GetStableVersion",
			args: args{
				resolve: StubResolverVersions{
					getStableVersion: func() (string, error) {
						return "", errors.New("any error")
					},
				},
				currentVersion: "1.0.0",
			},
			wantWriter: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			VerifyNewVersion(tt.args.resolve, writer, tt.args.currentVersion)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("VerifyNewVersion() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
