package notification

import (
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Notifier interface {
	Notify(msg string, to string) error
}

type SMSNotifier struct {
	twilio *twilio.RestClient
	from   string
}

func NewSMSNotifier(t *TwilioClient) (Notifier, error) {
	return &SMSNotifier{
		twilio: t.tc,
		from:   t.from,
	}, nil
}

func (n SMSNotifier) Notify(msg string, to string) error {

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(n.from)
	params.SetBody(msg)

	_, err := n.twilio.Api.CreateMessage(params)

	if err != nil {
		return err
	}

	return nil
}
