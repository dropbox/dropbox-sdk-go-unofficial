package dropbox

import (
	"net/http"

	"golang.org/x/oauth2"
)

type apiImpl struct {
	client  *http.Client
	verbose bool
}

func Client(token string, verbose bool) Api {
	var conf = &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/1/oauth2/authorize",
			TokenURL: "https://api.dropbox.com/1/oauth2/token",
		},
	}
	tok := &oauth2.Token{AccessToken: token}
	return &apiImpl{conf.Client(oauth2.NoContext, tok), verbose}
}
