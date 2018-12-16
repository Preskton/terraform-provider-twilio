package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/Preskton/terraform-provider-twilio/plugin/providers/twilio"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Debug("Loading terraform-provider-twilio")

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: twilio.Provider,
	})
}
