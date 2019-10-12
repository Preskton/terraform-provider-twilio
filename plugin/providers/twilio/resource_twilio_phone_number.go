package twilio

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/spf13/cast"
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"search": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"area_code": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
						"fallback_http_method": &schema.Schema{
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
						"fallback_http_method": &schema.Schema{
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

func addIfNotEmpty(params url.Values, key string, value interface{}) {
	s := cast.ToString(value)

	if s != "" {
		params.Add(key, s)
	}
}

func makeCreateRequestPayload(d *schema.ResourceData) url.Values {
	createRequestPayload := make(url.Values)

	addIfNotEmpty(createRequestPayload, "FriendlyName", d.Get("friendly_name"))
	addIfNotEmpty(createRequestPayload, "AddressSid", d.Get("address_sid"))
	addIfNotEmpty(createRequestPayload, "TrunkSid", d.Get("trunk_sid"))
	addIfNotEmpty(createRequestPayload, "IdentitySid", d.Get("identity_sid"))

	if sms := d.Get("sms").(*schema.Set); sms.Len() > 0 {
		sms := sms.List()[0].(map[string]interface{})

		addIfNotEmpty(createRequestPayload, "SmsApplicationSid", sms["application_sid"])
		addIfNotEmpty(createRequestPayload, "SmsFallbackUrl", sms["fallback_url"])
		addIfNotEmpty(createRequestPayload, "SmsFallbackMethod", sms["fallback_http_method"])
		addIfNotEmpty(createRequestPayload, "SmsMethod", sms["primary_http_method"])
		addIfNotEmpty(createRequestPayload, "SmsUrl", sms["primary_url"])
	}

	if voice := d.Get("voice").(*schema.Set); voice.Len() > 0 {
		voice := voice.List()[0].(map[string]interface{})

		addIfNotEmpty(createRequestPayload, "VoiceApplicationSid", voice["application_sid"])
		addIfNotEmpty(createRequestPayload, "VoiceFallbackUrl", voice["fallback_url"])
		addIfNotEmpty(createRequestPayload, "VoiceFallbackMethod", voice["fallback_http_method"]) // TODO Map to safe values
		addIfNotEmpty(createRequestPayload, "VoiceMethod", voice["primary_http_method"])          // TODO Map to safe values
		addIfNotEmpty(createRequestPayload, "VoiceUrl", voice["primary_url"])
		addIfNotEmpty(createRequestPayload, "VoiceCallerIdLookup", voice["caller_id_enabled"])
		addIfNotEmpty(createRequestPayload, "VoiceReceiveMode", voice["recieve_mode"]) // TODO Map to Twilio 
	}

	if statusCallback := d.Get("status_callback").(*schema.Set); statusCallback.Len() > 0 {
		statusCallback := statusCallback.List()[0].(map[string]interface{})

		addIfNotEmpty(createRequestPayload, "SmsMethod", statusCallback["primary_http_method"]) // TODO Map to safe values
		addIfNotEmpty(createRequestPayload, "SmsUrl", statusCallback["primary_url"])
	}

	if emergency := d.Get("emergency").(*schema.Set); emergency.Len() > 0 {
		emergency := emergency.List()[0].(map[string]interface{})

		addIfNotEmpty(createRequestPayload, "EmergencyStatus", emergency["enabled"]) // TODO Map to Twilio values
		addIfNotEmpty(createRequestPayload, "EmergencyAddressSid", emergency["address_sid"])
	}

	return createRequestPayload
}

func resourceTwilioPhoneNumberCreate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioPhoneNumberCreate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	//var searchParams url.Values
	searchParams := make(url.Values)

	areaCode := d.Get("area_code").(string)
	if len(areaCode) > 0 {
		searchParams.Set("AreaCode", areaCode)
	}

	search := d.Get("search").(string)	
	if len(search) > 0 {
		if !strings.Contains(search, "*") {
			// Assume they want to search with this as the start of the number
			search = search + "*"
		}

		searchParams.Set("Contains", search)
	}

	countryCode := d.Get("country_code").(string)

	log.WithFields(
		log.Fields{
			"account_sid":  config.AccountSID,
			"country_code": countryCode,
		},
	).Debug("START client.Available.Numbers.Local.GetPage")

	// TODO switch based on the type of number to buy local, mobile, intl
	searchResult, err := client.AvailableNumbers.Local.GetPage(context, countryCode, searchParams)

	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid":  config.AccountSID,
				"country_code": countryCode,
			},
		).Error("Caught an unexpected error when searching for phone numbers")

		return err
	}

	log.WithFields(
		log.Fields{
			"account_sid":  config.AccountSID,
			"country_code": countryCode,
			"search":       search,
			"result_count": len(searchResult.Numbers),
		},
	).Debug("END client.Available.Nubmers.Local.GetPage")

	if searchResult != nil && len(searchResult.Numbers) == 0 {
		log.WithFields(
			log.Fields{
				"account_sid":  config.AccountSID,
				"country_code": countryCode,
			},
		).Error("No phone numbers matched the search patterns")

		return errors.New("No numbers found that match your search")
	}

	// Grab the first number that matches
	number := searchResult.Numbers[0]

	// Per https://www.twilio.com/docs/phone-numbers/api/incoming-phone-numbers#create-an-incomingphonenumber-resource
	// the number must be in E.164 format, aka number with +, country code, number, without any other punctuation
	re := regexp.MustCompile("[ -]")
	e164Number := re.ReplaceAllLiteralString(number.PhoneNumber.Friendly(), "")

	buyParams := makeCreateRequestPayload(d)
	buyParams.Set("PhoneNumber", e164Number)

	log.WithFields(
		log.Fields{
			"account_sid":  config.AccountSID,
			"phone_number": e164Number,
		},
	).Debug("START client.IncomingNumbers.Create")

	buyResult, err := client.IncomingNumbers.Create(context, buyParams)

	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid":  config.AccountSID,
				"phone_number": e164Number,
			},
		).Error("Caught an error when attempting to purchase phone number: " + err.Error())

		return err
	}

	d.SetId(buyResult.Sid)
	d.Set("number", e164Number)

	log.WithFields(
		log.Fields{
			"account_sid":      config.AccountSID,
			"phone_number":     e164Number,
			"phone_number_sid": buyResult.Sid,
		},
	).Debug("END client.IncomingNumbers.Create")

	return nil
}

func resourceTwilioPhoneNumberRead(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioPhoneNumberRead")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	log.Debug("Getting SID")

	sid := d.Id()

	log.Debug("Getting phone_number")

	phoneNumber := d.Get("number").(string)

	log.WithFields(
		log.Fields{
			"account_sid":      config.AccountSID,
			"phone_number":     phoneNumber,
			"phone_number_sid": sid,
		},
	).Debug("START client.IncomingNumbers.Get")

	ph, err := client.IncomingNumbers.Get(context, sid)

	fmt.Printf(ph.APIVersion)

	log.WithFields(
		log.Fields{
			"account_sid":      config.AccountSID,
			"phone_number":     phoneNumber,
			"phone_number_sid": sid,
		},
	).Debug("END client.IncomingNumbers.Get")

	if err != nil {
		return fmt.Errorf("Failed to refresh number: %s", err.Error())
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
	phoneNumber := d.Get("number").(string)

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
