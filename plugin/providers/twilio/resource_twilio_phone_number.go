package twilio

import (
	"context"
	"errors"
	"fmt"
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
			"friendly_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"sms": &schema.Schema{
				Type:     schema.TypeSet,
				MinItems: 0,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_sid": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"primary_http_method": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"primary_url": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"fallback_method": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"fallback_url": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"status_callback": &schema.Schema{
				Type:     schema.TypeSet,
				MinItems: 0,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"http_method": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"voice": &schema.Schema{
				Type:     schema.TypeSet,
				MinItems: 0,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_sid": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"primary_http_method": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"primary_url": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"fallback_method": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"fallback_url": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"caller_id_enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"receive_mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"address_sid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"trunk_sid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"identity_sid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"emergency": &schema.Schema{
				Type:     schema.TypeSet,
				MinItems: 0,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"address_sid": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func flattenPhoneNumber(d *schema.ResourceData) url.Values {
	v := make(url.Values)

	v.Add("FriendlyName", d.Get("friendly_name").(string))
	v.Add("AddressSid", d.Get("address_sid").(string))
	v.Add("TrunkSid", d.Get("trunk_sid").(string))
	v.Add("IdentitySid", d.Get("identity_sid").(string))

	// TODO SMS

	if sms := d.Get("sms").(*schema.Set); sms.Len() > 0 {
		sms := sms.List()[0].(map[string]interface{})

		v.Add("SmsApplicationSid", sms["application_sid"].(string))
		v.Add("SmsFallbackUrl", sms["fallback_url"].(string))
		v.Add("SmsFallbackMethod", sms["fallback_http_method"].(string)) // TODO Map to safe values
		v.Add("SmsMethod", sms["primary_http_method"].(string))          // TODO Map to safe values
		v.Add("SmsUrl", sms["primary_url"].(string))
	}

	// TODO Voice

	if voice := d.Get("voice").(*schema.Set); voice.Len() > 0 {
		voice := voice.List()[0].(map[string]interface{})

		v.Add("VoiceApplicationSid", voice["application_sid"].(string))
		v.Add("VoiceFallbackUrl", voice["fallback_url"].(string))
		v.Add("VoiceFallbackMethod", voice["fallback_http_method"].(string)) // TODO Map to safe values
		v.Add("VoiceMethod", voice["primary_http_method"].(string))          // TODO Map to safe values
		v.Add("VoiceUrl", voice["primary_url"].(string))
		v.Add("VoiceCallerIdLookup", voice["caller_id_enabled"].(string))
		v.Add("VoiceReceiveMode", voice["recieve_mode"].(string)) // TODO Map to Twilio values
	}

	// TODO Status Callback

	if statusCallback := d.Get("status_callback").(*schema.Set); statusCallback.Len() > 0 {
		statusCallback := statusCallback.List()[0].(map[string]interface{})

		v.Add("SmsMethod", statusCallback["primary_http_method"].(string)) // TODO Map to safe values
		v.Add("SmsUrl", statusCallback["primary_url"].(string))
	}

	// TODO Emergency

	if emergency := d.Get("emergency").(*schema.Set); emergency.Len() > 0 {
		emergency := emergency.List()[0].(map[string]interface{})

		v.Add("EmergencyStatus", emergency["enabled"].(string)) // TODO Map to Twilio values
		v.Add("EmergencyAddressSid", emergency["address_sid"].(string))
	}

	return v
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
	log.Debug("ENTER resourceTwilioPhoneNumberDelete")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()
	phoneNumber := d.Get("phone_number").(string)

	log.WithFields(
		log.Fields{
			"account_sid":      config.AccountSID,
			"phone_number":     phoneNumber,
			"phone_number_sid": sid,
		},
	).Debug("START client.IncomingNumbers.Release")

	ph, err := client.IncomingNumbers.Get(context, sid)

	fmt.Printf(ph.APIVersion)

	log.WithFields(
		log.Fields{
			"account_sid":      config.AccountSID,
			"phone_number":     phoneNumber,
			"phone_number_sid": sid,
		},
	).Debug("END client.IncomingNumbers.Release")

	if err != nil {
		return fmt.Errorf("Failed to delete/release number: %s", err.Error())
	}

	return nil
}

func resourceTwilioPhoneNumberUpdate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("Not implemented")
}

func resourceTwilioPhoneNumberDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioPhoneNumberDelete")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()
	phoneNumber := d.Get("phone_number").(string)

	log.WithFields(
		log.Fields{
			"account_sid":      config.AccountSID,
			"phone_number":     phoneNumber,
			"phone_number_sid": sid,
		},
	).Debug("START client.IncomingNumbers.Release")

	err := client.IncomingNumbers.Release(context, sid)

	log.WithFields(
		log.Fields{
			"account_sid":      config.AccountSID,
			"phone_number":     phoneNumber,
			"phone_number_sid": sid,
		},
	).Debug("END client.IncomingNumbers.Release")

	if err != nil {
		return fmt.Errorf("Failed to delete/release number: %s", err.Error())
	}

	return nil
}
