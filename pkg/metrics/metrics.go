package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/session"
)

const (
	urlPattern = "%s/usage"
)

// CmdUse type that represents a metric use
type CmdUse struct {
	Username string `json:"username"`
	Cmd      string `json:"command"`
}

type Sender struct {
	serverURL      string
	httpClient     *http.Client
	sessionManager session.Manager
}

func NewSender(serverURL string, hc *http.Client, sm session.Manager) Sender {
	return Sender{serverURL: fmt.Sprintf(urlPattern, serverURL), httpClient: hc, sessionManager: sm}
}

func (s Sender) SendCommand() {
	session, err := s.sessionManager.Current()
	if err != nil {
		return
	}

	cmdUse := CmdUse{
		Username: session.Username,
		Cmd:      cmd(),
	}

	b, err := json.Marshal(&cmdUse)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, s.serverURL, bytes.NewBuffer(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-org", session.Organization)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.AccessToken))
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}

func cmd() string {
	var args []string
	for i := 0; i < len(os.Args); i++ {
		if i == len(os.Args)-1 {
			args = append(args, os.Args[i])
			continue
		}
		args = append(args, os.Args[i]+" ")
	}
	return strings.Join(args, "")
}
