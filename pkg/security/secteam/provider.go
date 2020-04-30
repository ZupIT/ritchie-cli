package secteam

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/security/oauth"
)

func loginChannelProvider(p security.AuthProvider, org, serverURL string) (chan security.ChanResponse, error) {
	cr := make(chan security.ChanResponse)
	switch p {
	case security.OAuthProvider:
		oauthCli := oauth.NewLoginManager(cr, serverURL)
		go func() {
			err := oauthCli.Login(org)
			if err != nil {
				fmt.Sprintln("Error in Login")
				return
			}
		}()
	default:
		return nil, security.ErrUnknownProvider
	}
	return cr, nil
}

func logoutChannelProvider(p security.AuthProvider, org, serverURL string) (chan security.ChanResponse, error) {
	cr := make(chan security.ChanResponse)
	switch p {
	case security.OAuthProvider:
		oauthCli := oauth.NewLogoutManager(org, cr, serverURL)
		go func() {
			err := oauthCli.Logout()
			if err != nil {
				fmt.Sprintln("Error in Logout")
				return
			}
		}()
	default:
		return nil, security.ErrUnknownProvider
	}
	return cr, nil
}
