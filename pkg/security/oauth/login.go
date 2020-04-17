package oauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"

	"github.com/ZupIT/ritchie-cli/pkg/security"
)

const (
	htmlLogin = `<!DOCTYPE html>

<html>
  <head>
    <meta charset="utf-8"/>
    <title>Login Ritchie</title>
    <style>
      html, body {
        height: 100%;
        justify-content: center;
        align-items: center; 
      }

      .container {
      	display: flex;
        justify-content: center;
        color: '#111';
        font-size: 30px;
      }
    </style>
  </head>

  <body>
    <div class="container">
      Login Successful
    </div>
    <div class="container">
      <span id="counter">5</span>s to close the browser.
    </div>
  </body>

  <script type="text/javascript"> 
    (function startSetInterval() {
      let count = 5;

      const interval = setInterval(function t() {
        const counter = document.getElementById('counter')
        counter.innerText = count;

        if (count === 0) {
          clearInterval(interval)
          window.close()
        }

        count = count - 1;
        return t;
      }(), 1000); 
    }())
</script>

</html>`
)

type LoginManager struct {
	Resp      chan security.ChanResponse
	serverURL string
}

func NewLoginManager(resp chan security.ChanResponse, serverURL string) *LoginManager {
	return &LoginManager{
		Resp:      resp,
		serverURL: serverURL,
	}
}

// Process login oauth
func (l *LoginManager) Login(org string) error {
	providerConfig, err := providerConfig(org, l.serverURL)
	if err != nil {
		l.Resp <- security.ChanResponse{Error: err}
		return nil
	}
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, providerConfig.Url)
	if err != nil {
		l.Resp <- security.ChanResponse{Error: err}
		return nil
	}
	oauth2Config := oauth2.Config{
		ClientID:     providerConfig.ClientId,
		ClientSecret: "",
		RedirectURL:  callbackURL(),
		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}
	state := "login"
	err = openBrowser(oauth2Config.AuthCodeURL(state))
	if err != nil {
		l.Resp <- security.ChanResponse{Error: err}
		return nil
	}
	http.HandleFunc(callbackPath, l.handlerLogin(provider, state, oauth2Config, ctx))
	return http.ListenAndServe(callbackAddr, nil)
}

func (l LoginManager) handlerLogin(provider *oidc.Provider, state string, oauth2Config oauth2.Config, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := &oidc.Config{
			ClientID: oauth2Config.ClientID,
		}
		verifier := provider.Verifier(cfg)
		if r.URL.Query().Get("state") != state {
			http.Error(w, "state did not match", http.StatusBadRequest)
			l.Resp <- security.ChanResponse{Error: errors.New(fmt.Sprint("state did not match http code ", http.StatusBadRequest))}
			return
		}

		oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			l.Resp <- security.ChanResponse{Error: err}
			return
		}

		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			l.Resp <- security.ChanResponse{Error: errors.New(fmt.Sprint("No id_token field in oauth2 token http code ", http.StatusInternalServerError))}
			return
		}

		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			l.Resp <- security.ChanResponse{Error: err}
			return
		}

		token := oauth2Token.AccessToken
		user := struct {
			Email    string `json:"email"`
			Username string `json:"preferred_username"`
		}{}
		_ = idToken.Claims(&user)
		_, _ = w.Write([]byte(htmlLogin))
		l.Resp <- security.ChanResponse{
			Token:    token,
			Username: user.Username,
		}
	}
}
