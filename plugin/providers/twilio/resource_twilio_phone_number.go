package twilio

import (
	"context"
	"errors"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTwilioPhoneNumber() *schema.Resource {
	return &schema.Resource{
		Create: resourceTwilioPhoneNumberCreate,
		Read:   resourceTwilioPhoneNumberRead,
		Update: resourceTwilioPhoneNumberUpdate,
		Delete: resourceTwilioPhoneNumberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sid": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"area_code": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"country_code": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"number": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTwilioPhoneNumberCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TerraformTwilioContext).client
	context := context.TODO()

	areaCode := d.Get("area_code").(string)
	countryCode := d.Get("country_code").(string)

	var searchParams url.Values
	searchParams.Set("AreaCode", areaCode)

	// TODO switch based on the type of number to buy local, mobile, intl
	searchResult, err := client.AvailableNumbers.Local.GetPage(context, countryCode, searchParams)

	if err != nil {
		return err
	}

	if searchResult != nil && len(searchResult.Numbers) == 0 {
		return errors.New("No numbers found that match area code")
	}

	number := searchResult.Numbers[0]

	var buyParams url.Values
	buyParams.Set("PhoneNumber", number.PhoneNumber.Local())

	buyResult, err := client.IncomingNumbers.Create(context, buyParams)

	d.SetId(buyResult.Sid)
	d.Set("phone_number", buyResult.PhoneNumber.Local())

	return nil
}

func resourceTwilioPhoneNumberRead(d *schema.ResourceData, meta interface{}) error {
	return errors.New("Not implemented")
}

func resourceTwilioPhoneNumberUpdate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("Not implemented")
}

func resourceTwilioPhoneNumberDelete(d *schema.ResourceData, meta interface{}) error {
	return errors.New("Not implemented")
}
