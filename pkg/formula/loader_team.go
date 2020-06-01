package formula

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

type TeamLoader struct {
	serverFinder server.Finder
	client       *http.Client
	session      session.Manager
	Adder
}

func NewTeamLoader(serverFinder server.Finder, client *http.Client, session session.Manager, adder Adder) TeamLoader {
	return TeamLoader{
		serverFinder: serverFinder,
		client:       client,
		session:      session,
		Adder:        adder,
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

	req.Header.Set("x-org", sess.Organization)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sess.AccessToken))
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

	var dd []Repository
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
