package credteam

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

const (
	urlConfigPattern = "%s/credentials/config"
)

var (
	ErrFieldsNotFound = errors.New("fields not found")
)

type Settings struct {
	serverFinder   server.Finder
	httpClient     *http.Client
	sessionManager session.Manager
	ctxFinder      rcontext.Finder
}

func NewSettings(serverFinder server.Finder, hc *http.Client, sm session.Manager, cf rcontext.Finder) Settings {
	return Settings{
		serverFinder:   serverFinder,
		httpClient:     hc,
		sessionManager: sm,
		ctxFinder:      cf,
	}
}

func (s Settings) Fields() (credential.Fields, error) {
	session, err := s.sessionManager.Current()
	if err != nil {
		return nil, err
	}

	ctx, err := s.ctxFinder.Find()
	if err != nil {
		return nil, err
	}

	cfg, err := s.serverFinder.Find()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(urlConfigPattern, cfg.URL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-org", session.Organization)
	req.Header.Set("x-ctx", ctx.Current)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.AccessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var cfg credential.Fields
		if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	case http.StatusNotFound:
		return nil, ErrFieldsNotFound
	default:
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(b))
	}
}
