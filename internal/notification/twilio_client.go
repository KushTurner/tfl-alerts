package notification

import "github.com/twilio/twilio-go"

type TwilioConfig struct {
	From       string
	AccountSID string
	AuthToken  string
}

type TwilioClient struct {
	tc   *twilio.RestClient
	from string
}

func NewTwilioClient(cfg *TwilioConfig) *TwilioClient {
	c := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.AccountSID,
		Password: cfg.AuthToken,
	})

	return &TwilioClient{
		tc:   c,
		from: cfg.From,
	}
}
