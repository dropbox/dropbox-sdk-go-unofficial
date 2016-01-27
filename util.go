package dropbox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type apiImpl struct {
	client *http.Client
}

func Client(token string) Api {
	var conf = &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/1/oauth2/authorize",
			TokenURL: "https://api.dropbox.com/1/oauth2/token",
		},
	}
	tok := &oauth2.Token{AccessToken: token}
	return &apiImpl{conf.Client(oauth2.NoContext, tok)}
}

func OauthClient(appId string, appSecret string, filePath string) (Api, error) {
	var conf = &oauth2.Config{
		ClientID:     appId,
		ClientSecret: appSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/1/oauth2/authorize",
			TokenURL: "https://api.dropbox.com/1/oauth2/token",
		},
	}
	tok, err := readToken(filePath)
	if err != nil {
		url := conf.AuthCodeURL("state")
		fmt.Printf("1. Go to %v\n", url)
		fmt.Println("2. Click \"Allow\" (you might have to log in first).")
		fmt.Println("3. Copy the authorization code.")
		fmt.Printf("Enter the authorization code here: ")

		var code string
		if _, err := fmt.Scan(&code); err != nil {
			return nil, err
		}
		tok, err = conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			return nil, err
		}
		saveToken(filePath, tok)
	}
	return &apiImpl{conf.Client(oauth2.NoContext, tok)}, nil
}

func readToken(filePath string) (*oauth2.Token, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var tok oauth2.Token
	if json.Unmarshal(b, &tok) != nil {
		return nil, err
	}

	if !tok.Valid() {
		return nil, fmt.Errorf("Token %v is no longer valid", tok)
	}

	return &tok, nil
}

func saveToken(filePath string, token *oauth2.Token) {
	if _, err := os.Stat(filePath); err != nil {
		if !os.IsNotExist(err) {
			return
		}
		// create file
		b, err := json.Marshal(token)
		if err != nil {
			return
		}
		if err = ioutil.WriteFile(filePath, b, 0644); err != nil {
			return
		}
	}
}
