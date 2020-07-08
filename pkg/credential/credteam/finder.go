package credteam

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/http/headers"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

const urlGetPattern = "%s/credentials/me/%s"

var ErrNotFoundCredential = errors.New("credential not found")

type Finder struct {
	serverFinder   server.Finder
	httpClient     *http.Client
	sessionManager session.Manager
	ctxFinder      rcontext.Finder
}

func NewFinder(serverFinder server.Finder, hc *http.Client, sm session.Manager, cf rcontext.Finder) Finder {
	return Finder{
		serverFinder:   serverFinder,
		httpClient:     hc,
		sessionManager: sm,
		ctxFinder:      cf,
	}
}

func (f Finder) Find(provider string) (credential.Detail, error) {
	session, err := f.sessionManager.Current()
	if err != nil {
		return credential.Detail{}, err
	}

	ctx, err := f.ctxFinder.Find()
	if err != nil {
		return credential.Detail{}, err
	}

	cfg, err := f.serverFinder.Find()
	if err != nil {
		return credential.Detail{}, err
	}

	url := fmt.Sprintf(urlGetPattern, cfg.URL, provider)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return credential.Detail{}, err
	}

	req.Header.Set(headers.XOrg, session.Organization)
	req.Header.Set(headers.XCtx, ctx.Current)
	req.Header.Set(headers.XAuthorization, session.AccessToken)
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return credential.Detail{}, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		cred := credential.Detail{}
		if err := json.NewDecoder(resp.Body).Decode(&cred); err != nil {
			return credential.Detail{}, err
		}
		cred.Username = session.Username
		return cred, nil
	case http.StatusNotFound:
		return credential.Detail{}, ErrNotFoundCredential
	default:
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return credential.Detail{}, err
		}
		log.Printf("Status code: %v", resp.StatusCode)
		return credential.Detail{}, errors.New(string(b))
	}
}
