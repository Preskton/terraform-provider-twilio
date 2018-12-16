package netbox

import (
	"fmt"
	"strings"
		
	log "github.com/sirupsen/logrus"

	"net/url"

	twiclient "github.com/kevinburke/twilio-go"
)

// Handles configuration and instantiates our Twilio client.
type Config struct {
	AccountSID string,
	AuthToken string,
	Endpoint string,
}

type TerraformTwilioContext struct {
	client        *twiclient.Client
	configuration Config
}

// Client creates a Twilio client and prepares it for use with Terraform.
func (self *Config) Client() (interface{}, error) {
	log.WithFields(
		log.Fields{
			"account_sid": self.AccountSID,
		},
	).Debug("Initializing Netbox client")

	// TODO Support unique endpoints

	client := twiclient.NewClient(self.AccountSID, self.AuthToken, nil)

	context := TerraformTwilioContext{
		client:        client,
		configuration: &self,
	}

	return &context, nil
}
