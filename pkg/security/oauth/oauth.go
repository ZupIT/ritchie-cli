package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"runtime"
)

const (
	HTTPStatus299      = 299
	callbackAddr       = "localhost:8888"
	callbackPath       = "/ritchie/callback"
	callbackURLPattern = "http://localhost:8888/%s"
	providerURLPattern = "%s/oauth"
)

type ProviderConfig struct {
	Url      string `json:"url"`
	ClientId string `json:"clientId"`
}

func callbackURL() string {
	return fmt.Sprintf(callbackURLPattern, callbackPath)
}

func providerConfig(organization, serverURL string) (ProviderConfig, error) {
	var provideConfig ProviderConfig
	url := fmt.Sprintf(providerURLPattern, serverURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return provideConfig, fmt.Errorf("Failed to providerConfig for org %s. \n%v", organization, err)
	}

	req.Header.Set("x-org", organization)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return provideConfig, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode > HTTPStatus299 {
		return provideConfig, fmt.Errorf("Failed to call url. %v for org %s. Status code: %d\n", url, organization, resp.StatusCode)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return provideConfig, fmt.Errorf("Failed parse response to body: %s\n", string(bodyBytes))
	}

	json.Unmarshal(bodyBytes, &provideConfig)
	return provideConfig, nil
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = nil
	}
	return err
}
