package gmail

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"log"
)

type Client struct {
	*gmail.Service
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func GetClient(tokenStorage TokenStorage) Client {
	err, config := tokenStorage.GetConfig()
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	err, token := tokenStorage.GetToken()
	if err != nil {
		token = getTokenFromWeb(config)
		tokenStorage.SaveToken(token)
	}

	srv, err := gmail.New(config.Client(context.Background(), token))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	return Client{srv}
}
