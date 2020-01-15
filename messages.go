package gmail

import (
	"github.com/senfix/gmail/model"
	"github.com/senfix/gmail/parser"
	"google.golang.org/api/gmail/v1"
)

type Messages interface {
	GetMessage(id string) (err error, message model.Message)
	GetMessages(filters ...string) (err error, message []model.Message)
}

type messages struct {
	*gmail.UsersMessagesService
	parser.Email
}

func NewMessages(c Client, s parser.Email) Messages {
	return &messages{c.Users.Messages, s}
}

func (t *messages) GetMessages(filters ...string) (err error, messages []model.Message) {
	call := t.List("me")
	for _, filter := range filters {
		call = call.Q(filter)
	}

	msgs, err := call.Do()
	if err != nil {
		return
	}

	messages = make([]model.Message, len(msgs.Messages))
	for i, m := range msgs.Messages {
		err, messages[i] = t.GetMessage(m.Id)
		if err != nil {
			return err, nil
		}
	}
	return
}

func (t *messages) GetMessage(id string) (err error, message model.Message) {
	response, err := t.Get("me", id).Format("full").Do()
	if err != nil {
		return
	}
	err, message = t.Parse(response)
	return
}
