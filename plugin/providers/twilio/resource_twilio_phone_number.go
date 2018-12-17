package twilio

import (
	"context"
	"errors"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"

	log "github.com/sirupsen/logrus"
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
	log.Debug("ENTER resourceTwilioPhoneNumberCreate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	log.Debug("Getting area_code and country_code")

	areaCode := d.Get("area_code").(string)
	countryCode := d.Get("country_code").(string)

	log.Debug("Setting searchParams url.Values")

	//var searchParams url.Values
	searchParams := make(url.Values)
	searchParams.Set("AreaCode", areaCode)

	log.WithFields(
		log.Fields{
			"account_sid":  config.AccountSID,
			"country_code": countryCode,
			"area_code":    areaCode,
		},
	).Debug("START client.Available.Nubmers.Local.GetPage")

	// TODO switch based on the type of number to buy local, mobile, intl
	searchResult, err := client.AvailableNumbers.Local.GetPage(context, countryCode, searchParams)

	if err != nil {
		return err
	}

	log.WithFields(
		log.Fields{
			"account_sid":  config.AccountSID,
			"country_code": countryCode,
			"area_code":    areaCode,
			"result_count": len(searchResult.Numbers),
		},
	).Debug("END client.Available.Nubmers.Local.GetPage")

	if searchResult != nil && len(searchResult.Numbers) == 0 {
		return errors.New("No numbers found that match area code")
	}

	number := searchResult.Numbers[0]

	buyParams := make(url.Values)
	buyParams.Set("PhoneNumber", number.PhoneNumber.Local())

	log.WithFields(
		log.Fields{
			"account_sid":  config.AccountSID,
			"phone_number": number.PhoneNumber.Local(),
		},
	).Debug("START client.IncomingNumbers.Create")

	buyResult, err := client.IncomingNumbers.Create(context, buyParams)

	d.SetId(buyResult.Sid)
	d.Set("phone_number", buyResult.PhoneNumber.Local())

	log.WithFields(
		log.Fields{
			"account_sid":      config.AccountSID,
			"phone_number":     number.PhoneNumber.Local(),
			"phone_number_sid": buyResult.Sid,
		},
	).Debug("END client.IncomingNumbers.Create")

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
