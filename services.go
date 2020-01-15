package gmail

import "github.com/senfix/gmail/parser"

type Services struct {
	Client   Client
	Messages Messages
}

func New(tokenStorage TokenStorage) Services {
	client := GetClient(tokenStorage)
	emailParser := parser.NewEmailParser()
	messages := NewMessages(client, emailParser)

	return Services{
		Client:   client,
		Messages: messages,
	}
}
