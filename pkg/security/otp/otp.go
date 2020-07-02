package otp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/http/headers"
	"github.com/ZupIT/ritchie-cli/pkg/server"
)

const (
	otpUrlPattern    = "%s/otp"
)

func NewOtpResolver(hc *http.Client) DefaultOtpResolver {
	return DefaultOtpResolver{
		httpClient: hc,
	}
}

type DefaultOtpResolver struct {
	httpClient *http.Client
}

func (dor DefaultOtpResolver) RequestOtp(serverUrl, organization string) (Response, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(otpUrlPattern, serverUrl), nil)
	response := Response{}
	if err != nil {
		return response, err
	}
	req.Header.Set(headers.XOrg, organization)
	resp, err := dor.httpClient.Do(req)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf(server.ErrPattern, serverUrl, resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, err
	}
	return response, nil
}
