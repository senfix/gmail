package parser

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/senfix/gmail/model"
	"google.golang.org/api/gmail/v1"
)

type Email interface {
	Parse(input *gmail.Message) (err error, message model.Message)
}

func NewEmailParser() Email {
	mailRegExp, _ := regexp.Compile("<([^>]+)>")
	return &email{mailRegExp}
}

type email struct {
	mailRegExp *regexp.Regexp
}

func (t *email) Parse(input *gmail.Message) (err error, message model.Message) {
	message = model.Message{
		GoogleId: input.Id,
		Message:  t.decodeBody(input.Payload.Body),
	}
	fmt.Printf("%+v\n", input.Payload)
	message.GoogleId = input.Id

	for _, h := range input.Payload.Headers {
		switch h.Name {
		case "Subject":
			message.Subject = h.Value
		case "From":
			message.From = t.parseAddress(h.Value)
		case "To":
			message.To = t.parseAddress(h.Value)
		case "Date":
			sec := input.InternalDate / 1000
			nsec := input.InternalDate - 1000*sec
			message.Date = time.Unix(sec, nsec)
		}
	}

	return
}

func (t *email) decodeBody(input *gmail.MessagePartBody) string {
	fmt.Printf("%+v\n", input)
	data, _ := base64.URLEncoding.DecodeString(input.Data)
	return string(data)
}

func (t *email) parseAddress(input string) []model.Address {
	inputs := strings.Split(input, ",")
	addresses := make([]model.Address, len(inputs))
	for i, v := range inputs {
		matches := t.mailRegExp.FindStringSubmatch(v)
		a := model.Address{
			Name:  "",
			Email: v,
		}
		if len(matches) == 2 {
			a.Email = matches[1]
			a.Name = strings.Trim(v[0:len(v)-len(a.Email)-3], " \"")
			if a.Name == a.Email {
				a.Name = ""
			}
		}
		addresses[i] = a
	}
	return addresses
}
