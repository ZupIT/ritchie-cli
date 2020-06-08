package secteam

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/http/headers"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

const urlLoginPattern = "%s/login"

type LoginManager struct {
	serverFinder   server.Finder
	httpClient     *http.Client
	sessionManager session.Manager
}

func NewLoginManager(
	serverFinder server.Finder,
	hc *http.Client,
	sm session.Manager) LoginManager {
	return LoginManager{
		serverFinder:   serverFinder,
		httpClient:     hc,
		sessionManager: sm,
	}
}

func (l LoginManager) Login(user security.User) error {
	cfg, err := l.serverFinder.Find()
	if err != nil {
		return err
	}
	fmt.Println("Organization:", cfg.Organization)

	b, err := json.Marshal(&user)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(urlLoginPattern, cfg.URL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(headers.XOrg, cfg.Organization)

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type loginResponse struct {
		Token string `json:"token"`
		TTL   int64  `json:"ttl"`
	}

	lr := loginResponse{}

	switch resp.StatusCode {
	case 200:
		if err = json.Unmarshal(b, &lr); err != nil {
			return err
		}
		sess := session.Session{
			AccessToken:  lr.Token,
			Organization: cfg.Organization,
			Username:     user.Username,
			TTL:          lr.TTL,
		}
		err = l.sessionManager.Create(sess)
		if err != nil {
			return err
		}
		return nil
	case 401:
		return errors.New("login failed! Verify your credentials")
	default:
		return errors.New("login failed")
	}

}
