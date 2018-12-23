package twilio_test

import (
	"github.com/hashicorp/terraform/helper/schema"
	. "github.com/onsi/ginkgo"

	"github.com/Preskton/terraform-provider-twilio/plugin/providers/twilio"
	. "github.com/onsi/gomega"
)

type WeaponStats struct {
	Power      int    `terraform:"power_value"`
	Range      int    `terraform:"range_value"`
	RateOfFire int    `terraform:"rof"`
	Adjective  string `terraform:"adj"`
	IsOP       bool   `terraform:"is_op"`
}

type Weapon struct {
	WeaponID     int         `terraform:"weapon_id"`
	Name         string      `terraform:"name"`
	Manufacturer string      `terraform:"manufacturer_name"`
	Stats        WeaponStats `terraform:"stats"`
}

func resourceTestWidget() *schema.Resource {
	return &schema.Resource{
		Create: nil,
		Read:   nil,
		Update: nil,
		Delete: nil,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"weapon_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"manufacturer_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"stats": &schema.Schema{
				Type:     schema.TypeSet,
				MinItems: 0,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"power_value": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"range_value": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"rof": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"adj": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_op": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

var _ = Describe("Twilio Terraform Provider", func() {
	var (
		weapons = map[string]*Weapon{
			"tkSplatRoller": &Weapon{
				WeaponID:     1,
				Name:         "Kensa Splat Roller",
				Manufacturer: "Toni Kensa",
				Stats: WeaponStats{
					Power:      100,
					Range:      5,
					RateOfFire: 35,
					Adjective:  "groovy",
					IsOP:       true,
				},
			},
		}
		expected *Weapon
		tfdata   *schema.ResourceData
		tfschema map[string]*schema.Schema
		err      error
	)

	Describe("Terraform Marshal", func() {

		BeforeSuite(func() {
			expected = weapons["tkSplatRoller"]
			tfdata = resourceTestWidget().TestResourceData()
			tfschema = resourceTestWidget().Schema
			err = twilio.MarshalToTerraform(expected, tfdata, tfschema)
		})

		Context("When it serializes from a struct to a Terraform ResourceData", func() {

			It("should not error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should populate basic (value-type) ResourceData per the struct tags", func() {
				weaponID := tfdata.Get("weapon_id").(int)
				Expect(weaponID).To(Equal(expected.WeaponID))

				weaponName := tfdata.Get("name").(string)
				Expect(weaponName).To(Equal(expected.Name))

				manufacturer := tfdata.Get("manufacturer_name").(string)
				Expect(manufacturer).To(Equal(expected.Manufacturer))
			})

			It("should have at least one entry in the nested `stats` Set", func() {
				stats := tfdata.Get("stats").(*schema.Set)
				Expect(stats.Len()).To(Equal(1))
			})

			It("should have copied the field values into the nested `stats` Set values", func() {
				stats := tfdata.Get("stats").(*schema.Set)
				values := stats.List()[0].(map[string]interface{})

				Expect(values["power_value"]).To(Equal(expected.Stats.Power))
				Expect(values["range_value"]).To(Equal(expected.Stats.Range))
				Expect(values["rof"]).To(Equal(expected.Stats.RateOfFire))
				Expect(values["adj"]).To(Equal(expected.Stats.Adjective))
				Expect(values["is_op"]).To(Equal(expected.Stats.IsOP))
			})
		})
	})
})
