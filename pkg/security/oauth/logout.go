package oauth

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"net/http"
)

const (
	logoutUrlPattern = "%s/protocol/openid-connect/logout?redirect_uri=%s"
	htmlLogout       = `<!DOCTYPE html>

<html>
	<head>
		<meta charset="utf-8"/>
		<title>Logout Ritchie</title>
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
			Logout Successful
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

type LogoutManager struct {
	Resp chan security.ChanResponse
	Org  string
	serverURL string
}

func NewLogoutManager(org string, resp chan security.ChanResponse, serverURL string) *LogoutManager {
	return &LogoutManager{
		Resp: resp,
		Org:  org,
		serverURL: serverURL,
	}
}

func (l *LogoutManager) Logout() error {
	providerConfig, err := providerConfig(l.Org, l.serverURL)
	if err != nil {
		l.Resp <- security.ChanResponse{Error: err}
		return nil
	}

	logoutUrl := fmt.Sprintf(logoutUrlPattern, providerConfig.Url, callbackURL())
	err = openBrowser(logoutUrl)
	if err != nil {
		l.Resp <- security.ChanResponse{Error: err}
		return nil
	}

	http.HandleFunc(callbackPath, l.handlerLogout())
	return http.ListenAndServe(callbackAddr, nil)
}

func (l *LogoutManager) handlerLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(htmlLogout))
		l.Resp <- security.ChanResponse{}
	}
}
