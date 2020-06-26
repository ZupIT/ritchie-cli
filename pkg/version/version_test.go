package version

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

type StubResolverVersions struct {
	stableVersion func(fromCmd bool) (string, error)
}

func (r StubResolverVersions) StableVersion(fromCmd bool) (string, error) {
	return r.stableVersion(fromCmd)
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

func TestDefaultVersionResolver_StableVersion(t *testing.T) {

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
		fromCmd          bool
		CurrentVersion   string
		StableVersionUrl string
		FileUtilService  fileutil.Service
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
				fromCmd: true,
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
				fromCmd: true,
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
				fromCmd: false,
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
				fromCmd: false,
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
			got, err := r.StableVersion(false)
			if (err != nil) != tt.wantErr {
				t.Errorf("StableVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StableVersion() got = %v, want %v", got, tt.want)
			}
			if tt.fields.fromCmd {
				got, err = r.StableVersion(true)
				if (err != nil) != tt.wantErr {
					t.Errorf("StableVersion() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("StableVersion() got = %v, want %v", got, tt.want)
				}
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
		name string
		args args
		want string
	}{
		{
			name: "Should return empty when current version equals to stableVersion",
			args: args{
				resolve: StubResolverVersions{
					stableVersion: func(fromCmd bool) (string, error) {
						return "1.0.0", nil
					},
				},
				currentVersion: "1.0.0",
			},
			want: "",
		},
		{
			name: "Should return msg when current version it not equals to stableVersion",
			args: args{
				resolve: StubResolverVersions{
					stableVersion: func(fromCmd bool) (string, error) {
						return "1.0.1", nil
					},
				},
				currentVersion: "1.0.0",
			},
			want: MsgRitUpgrade,
		},
		{
			name: "Should return empty on error in StableVersion ",
			args: args{
				resolve: StubResolverVersions{
					stableVersion: func(fromCmd bool) (string, error) {
						return "", errors.New("any error")
					},
				},
				currentVersion: "1.0.0",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyNewVersion(tt.args.resolve, tt.args.currentVersion); got != tt.want {
				t.Errorf("VerifyNewVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
