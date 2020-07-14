package otp

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewOtpResolver(t *testing.T) {
	type args struct {
		hc *http.Client
	}
	tests := []struct {
		name string
		args args
		want DefaultOtpResolver
	}{
		{
			name: "success run",
			args: args{
				hc: http.DefaultClient,
			},
			want: DefaultOtpResolver{httpClient: http.DefaultClient},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOtpResolver(tt.args.hc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOtpResolver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultOtpResolver_RequestOtp(t *testing.T) {

	mockSuccessUrl := "http://localhost:8882"
	mockFailUrl := "http://localhost:8882/request-fail"

	type fields struct {
		httpClient *http.Client
	}
	type args struct {
		serverUrl    string
		organization string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Response
		wantErr bool
	}{
		{
			name:    "success run",
			fields:  fields{
				httpClient: http.DefaultClient,
			},
			args:    args{
				serverUrl: mockSuccessUrl,
				organization: "anyone",
			},
			want:    Response{
				Otp: true,
			},
			wantErr: false,
		},
		{
			name:    "error while doing request",
			fields:  fields{
				httpClient: http.DefaultClient,
			},
			args:    args{
				serverUrl: "any url",
				organization: "anyone",
			},
			want:    Response{},
			wantErr: true,
		},
		{
			name:    "request responds 500 status code",
			fields:  fields{
				httpClient: http.DefaultClient,
			},
			args:    args{
				serverUrl: mockFailUrl,
				organization: "anyone",
			},
			want:    Response{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dor := DefaultOtpResolver{
				httpClient: tt.fields.httpClient,
			}
			got, err := dor.RequestOtp(tt.args.serverUrl, tt.args.organization)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestOtp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RequestOtp() got = %v, want %v", got, tt.want)
			}
		})
	}
}