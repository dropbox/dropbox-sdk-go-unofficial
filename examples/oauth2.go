package examples

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/dropbox/dropbox-sdk-go-unofficial"
)

func config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "<FILL_THIS_IN>",
		ClientSecret: "<FILL_THIS_IN>",
		Endpoint:     dropbox.OAuthEndpoint(""),
		// You will need to update this to a redirect uri that is added in your app settings.
		RedirectURL: "http://localhost:3000/oauth/dropbox",
	}
}

func Oauth2Server() {
	http.HandleFunc("/oauth/dropbox/connect", dropboxRedirect)
	// This needs to match your RedirectURL above in config()
	http.HandleFunc("/oauth/dropbox", dropboxCodeExchange)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func dropboxRedirect(w http.ResponseWriter, r *http.Request) {
	conf := config()
	http.Redirect(w, r, conf.AuthCodeURL("fake-state"), http.StatusFound)
}

func dropboxCodeExchange(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	log.Println("Code is... ", code)

	conf := config()
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("Token: %v\n", tok)))
}
