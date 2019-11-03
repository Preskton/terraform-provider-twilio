package twilio

import (
	log "github.com/sirupsen/logrus"

	twilio "github.com/kevinburke/twilio-go"
)

// Config contains our different configuration attributes and instantiates our Twilio client.
type Config struct {
	AccountSID string
	AuthToken  string
	Endpoint   string
}

// TerraformTwilioContext is our Terraform context that will contain both our Twilio client and configuration for access downstream.
type TerraformTwilioContext struct {
	client        *twilio.Client
	configuration Config
}

// Client creates a Twilio client and prepares it for use with Terraform.
func (config *Config) Client() (interface{}, error) {
	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("Initializing Twilio client")

	// TODO Support unique endpoints

	client := twilio.NewClient(config.AccountSID, config.AuthToken, nil)

	context := TerraformTwilioContext{
		client:        client,
		configuration: *config,
	}

	return &context, nil
}
