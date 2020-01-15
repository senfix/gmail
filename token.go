package gmail

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"log"
	"os"
)

type TokenStorage interface {
	SaveToken(token *oauth2.Token)
	GetToken() (err error, token *oauth2.Token)
	GetConfig() (err error, config *oauth2.Config)
}

type tokenStorage struct {
	tokenPath  string
	configPath string
}

func NewTokenStorage(tokenPath string, configPath string) TokenStorage {
	return &tokenStorage{tokenPath: tokenPath, configPath: configPath}
}

func (t tokenStorage) SaveToken(token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", t.tokenPath)
	f, err := os.OpenFile(t.tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func (t tokenStorage) GetToken() (err error, token *oauth2.Token) {
	f, err := os.Open(t.tokenPath)
	if err != nil {
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(token)
	return
}

func (t tokenStorage) GetConfig() (err error, config *oauth2.Config) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err = google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return
}
