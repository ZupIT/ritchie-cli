package secteam

import (
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/security/oauth"
)

func loginChannelProvider(p security.AuthProvider, org string) (chan security.ChanResponse, error) {
	cr := make(chan security.ChanResponse)
	switch p {
	case security.OAuthProvider:
		oauthCli := oauth.NewLoginManager(cr)
		go oauthCli.Login(org)
	default:
		return nil, security.ErrUnknownProvider
	}
	return cr, nil
}

func logoutChannelProvider(p security.AuthProvider, org string) (chan security.ChanResponse, error) {
	cr := make(chan security.ChanResponse)
	switch p {
	case security.OAuthProvider:
		oauthCli := oauth.NewLogoutManager(org, cr)
		go oauthCli.Logout()
	default:
		return nil, security.ErrUnknownProvider
	}
	return cr, nil
}
