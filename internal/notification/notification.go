package notification

import (
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Notifier interface {
	Notify(data interface{}, to string) error
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

func (n SMSNotifier) Notify(msg interface{}, to string) error {
	channel := "whatsapp:"

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(channel + to)
	params.SetFrom(channel + n.from)
	params.SetBody(msg.(string))

	_, err := n.twilio.Api.CreateMessage(params)

	if err != nil {
		return err
	}

	return nil
}
