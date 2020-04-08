package secteam

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"io/ioutil"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/session"
)

const (
	urlPattern = "%s/users"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserManager struct {
	serverURL      string
	httpClient     *http.Client
	sessionManager session.Manager
}

func NewUserManager(serverURL string, hc *http.Client, sm session.Manager) UserManager {
	return UserManager{serverURL: fmt.Sprintf(urlPattern, serverURL), httpClient: hc, sessionManager: sm}
}

func (u UserManager) Create(user security.User) error {
	s, err := u.sessionManager.Current()
	if err != nil {
		return err
	}

	b, err := json.Marshal(&user)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u.serverURL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-org", s.Organization)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))
	resp, err := u.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	default:
		return errors.New(string(b))
	}
}

func (u UserManager) Delete(user security.User) error {
	s, err := u.sessionManager.Current()
	if err != nil {
		return err
	}

	b, err := json.Marshal(&user)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, u.serverURL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-org", s.Organization)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))
	res, err := u.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return ErrUserNotFound
	default:
		return errors.New(string(b))
	}
}
