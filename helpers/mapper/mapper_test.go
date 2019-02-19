package mapper_test

import (
	"github.com/hashicorp/terraform/helper/schema"
	. "github.com/onsi/ginkgo"

	"github.com/Preskton/terraform-provider-twilio/helpers/mapper"
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
	WeaponID           string      `terraform:"id"`
	Name               string      `terraform:"name"`
	Manufacturer       string      `terraform:"manufacturer_name"`
	Stats              WeaponStats `terraform:"stats"`
	PowerUpCosts       []int       `terraform:"power_up_costs"`
	SomethingWithNoTag int         `notthetagyourelookingfor:"lol"`
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
			"power_up_costs": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

var _ = Describe("Preskton's Mappers", func() {
	var (
		weapons = map[string]*Weapon{
			"tkSplatRoller": &Weapon{
				WeaponID:     "TK1337",
				Name:         "Kensa Splat Roller",
				Manufacturer: "Toni Kensa",
				Stats: WeaponStats{
					Power:      100,
					Range:      5,
					RateOfFire: 35,
					Adjective:  "groovy",
					IsOP:       true,
				},
				PowerUpCosts: []int{5, 10, 15, 20, 25},
			},
		}
		expected *Weapon
		tfdata   *schema.ResourceData
		tfschema map[string]*schema.Schema
		err      error
	)

	Describe("Field Mapper", func() {
		Context("When it walks a struct to get the mapped fields for `terraform`", func() {

			expected = weapons["tkSplatRoller"]
			actualMap, mappingErr := mapper.MapStructByTag(expected, "terraform")

			It("should not error", func() {
				Expect(mappingErr).ShouldNot(HaveOccurred())
			})

			It("should provide an entry per field marked with `terraform`", func() {
				Expect(actualMap).ShouldNot(BeNil())
				Expect(len(actualMap)).Should(Equal(5))

				var ok bool

				weaponID, ok := actualMap["id"]
				Expect(ok).To(Equal(true))
				Expect(weaponID).To(Equal(expected.WeaponID))

				name, ok := actualMap["name"]
				Expect(ok).To(Equal(true))
				Expect(name).To(Equal(expected.Name))

				manufacturer, ok := actualMap["manufacturer_name"]
				Expect(ok).To(Equal(true))
				Expect(manufacturer).To(Equal(expected.Manufacturer))
			})

			It("should properly handle list fields", func() {
				powerUpCosts, ok := actualMap["power_up_costs"]
				Expect(ok).To(Equal(true))
				Expect(powerUpCosts).To(Equal(expected.PowerUpCosts))
			})

		})
	})

	Describe("Terraform Marshal", func() {

		BeforeEach(func() {
			expected = weapons["tkSplatRoller"]
			tfdata = resourceTestWidget().TestResourceData()
			tfschema = resourceTestWidget().Schema
			err = mapper.MarshalToTerraform(expected, tfdata, tfschema)
		})

		Context("When it serializes from a struct to a Terraform ResourceData", func() {

			It("should not error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should properly handle the special TF_ID identifier", func() {
				weaponID := tfdata.Id()
				Expect(weaponID).To(Equal(expected.WeaponID))
			})

			It("should populate basic (value-type) ResourceData on the top level per the struct tags", func() {
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

			It("should have properly handled the list fields", func() {
				//powerUpCosts := tfdata.Get("power_up_costs")

				//Expect(powerUpCosts).To(Equal(5))
			})
		})
	})
})
