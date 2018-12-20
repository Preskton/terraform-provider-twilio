package twilio_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTwilio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Twilio Suite")
}
