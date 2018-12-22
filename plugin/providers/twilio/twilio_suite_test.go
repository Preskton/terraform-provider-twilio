package twilio_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	log "github.com/sirupsen/logrus"
)

func TestTwilio(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Twilio Suite")
}
