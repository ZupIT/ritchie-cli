package credteam

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

const urlCreatePattern = "%s/credentials/%s"

type Setter struct {
	serverURL      string
	httpClient     *http.Client
	sessionManager session.Manager
	ctxFinder      rcontext.Finder
}

func NewSetter(serverURL string, hc *http.Client, sm session.Manager, cf rcontext.Finder) Setter {
	return Setter{
		serverURL:      serverURL,
		httpClient:     hc,
		sessionManager: sm,
		ctxFinder:      cf,
	}
}

func (s Setter) Set(cred credential.Detail) error {
	session, err := s.sessionManager.Current()
	if err != nil {
		return err
	}

	ctx, err := s.ctxFinder.Find()
	if err != nil {
		return err
	}

	b, err := json.Marshal(&cred)
	if err != nil {
		return err
	}

	path := credential.Me
	if cred.Username != credential.Me {
		path = credential.Admin
	} else {
		cred.Username = ""
	}

	url := fmt.Sprintf(urlCreatePattern, s.serverURL, path)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-org", session.Organization)
	req.Header.Set("x-ctx", ctx.Current)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.AccessToken))
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 201:
		return nil
	default:
		log.Printf("Status code: %v", resp.StatusCode)
		return errors.New(string(b))
	}
}
