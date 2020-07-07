package repo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/http/headers"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

type TeamLoader struct {
	serverFinder server.Finder
	client       *http.Client
	session      session.Manager
	formula.RepoAdder
}

func NewTeamLoader(serverFinder server.Finder, client *http.Client, session session.Manager, adder formula.RepoAdder) TeamLoader {
	return TeamLoader{
		serverFinder: serverFinder,
		client:       client,
		session:      session,
		RepoAdder:    adder,
	}
}

func (dm TeamLoader) Load() error {
	sess, err := dm.session.Current()
	if err != nil {
		return err
	}

	cfg, err := dm.serverFinder.Find()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(providerPath, cfg.URL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set(headers.XOrg, sess.Organization)
	req.Header.Set(headers.XAuthorization, sess.AccessToken)
	resp, err := dm.client.Do(req)
	if err != nil {
		return err
	}

	body, err := fileutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("%d - %s\n", resp.StatusCode, string(body))
	}

	var dd []formula.Repository
	if err := json.Unmarshal(body, &dd); err != nil {
		return err
	}

	for _, v := range dd {
		if err := dm.Add(v); err != nil {
			return err
		}
	}

	return nil
}
