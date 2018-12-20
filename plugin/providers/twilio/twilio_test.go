package twilio_test

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	. "github.com/onsi/ginkgo"

	"github.com/Preskton/terraform-provider-twilio/plugin/providers/twilio"
	. "github.com/onsi/gomega"
)

type helperTestData struct {
	SomeID     int    `terraform:"some_id"`
	SomeString string `terraform:"some_string"`
}

var _ = Describe("Twilio Terraform Provider", func() {
	var (
		testData = helperTestData{SomeID: 1337, SomeString: "Yarn"}
		tfdata   *schema.ResourceData
		err      error
	)

	Describe("Serialization helpers", func() {
		Context("When it serializes from a struct to a Terraform ResourceData", func() {
			BeforeEach(func() {
				tfdata = &schema.ResourceData{}
				err = twilio.MarshalToTerraform(testData, tfdata)

				fmt.Printf("%+v\n", tfdata)
			})

			It("should not error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should populate the ResourceData per the struct tags", func() {
				someID := tfdata.Get("some_id").(int)
				Expect(someID).To(Equal(1337))

				someString := tfdata.Get("some_string").(string)
				Expect(someString).To(Equal("Yarn"))
			})
		})
	})

})
