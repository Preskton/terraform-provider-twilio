package twilio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	twilio "github.com/kevinburke/twilio-go"
	"github.com/spf13/cast"

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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this phone number.",
			},
			"search": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Look for this number sequence anywhere in the phone number.",
			},
			"area_code": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Look for a number within this area code.",
			},
			"country_code": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Two letter ISO country code in which you want to search for a number. See https://support.twilio.com/hc/en-us/articles/223183068-Twilio-international-phone-number-availability-and-their-capabilities for details on available countries.",
			},
			"number": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full phone number, including country and area code.",
			},
			"friendly_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A friendly, human-readable name by which you can refer to this number.",
			},
			"date_created": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the phone number was created.",
			},
			"date_updated": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the phone number was laste updated.",
			},
			"address_requirements": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Address requirements imposed on this number, if any.",
			},
			"is_beta": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not this phone number is new to Twilio (beta status).",
			},
			"is_mms_capable": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not this phone number is MMS-capable.",
			},
			"is_sms_capable": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not this phone number is SMS-capable.",
			},
			"is_voice_capable": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not this phone number is voice-capable..",
			},									
			"sms": &schema.Schema{
				Type:     schema.TypeSet,
				MinItems: 0,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_sid": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SID of the Twilio application to invoke when an SMS is sent to this number.",
						},
						"primary_http_method": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HTTP method for the primary URL. Can be `GET` or `POST`, defaults to `POST`.",
						},
						"primary_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL called when an SMS is sent to this number.",
						},
						"fallback_http_method": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HTTP method for the fallback URL. Can be `GET` or `POST`, defaults to `POST`.",
						},
						"fallback_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL called if the primary URL returns a non-favorable status code.",
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
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL called when a whenever a status change occurs on this number.",
						},
						"http_method": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HTTP method for the status callback URL. Can be `GET` or `POST`, defaults to `POST`.",
						},
					},
				},
			},
			"voice": &schema.Schema{
				Type:     schema.TypeSet,
				MinItems: 0,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_sid": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SID of the Twilio application to invoke when a call is started with this number.",
						},
						"primary_http_method": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HTTP method for the primary URL. Can be `GET` or `POST`, defaults to `POST`.",
						},
						"primary_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL called when a phone call starts on this number.",
						},
						"fallback_http_method": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HTTP method for the fallback URL. Can be `GET` or `POST`, defaults to `POST`.",
						},
						"fallback_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL called if the primary URL returns a non-favorable status code.",
						},
						"caller_id_enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If caller ID is enabled or not for this number. If enabled, incurs additional charge per call (see console for pricing). Can be `true` or `false`, defaults to `false`.",
						},
						"receive_mode": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Determines if the line is set up for voice or fax. Can be `voice` or `fax`, defaults to `voice`.",
						},
					},
				},
			},
			"address_sid": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SID of the address associated with this phone number. May be required for certain countries.",
			},
			"trunk_sid": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SID of the voice trunk that will handle calls to this number. If set, overrides any voice URLs or applications: only the trunk will recieve the incoming call.",
			},
			"identity_sid": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SID of the identity associated with the phone number. May be required in certain countries.",
			},
			"emergency": &schema.Schema{
				Type:     schema.TypeSet,
				MinItems: 0,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status of this phone number. Either `active` or `inactive`.",
						},
						"address_sid": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SID of the address used for emergency calling from this number. The address must be validated before it can be used for emergency purposes.",
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
		addIfNotEmpty(createRequestPayload, "VoiceReceiveMode", voice["receive_mode"]) // TODO Map to Twilio
	}

	if statusCallback := d.Get("status_callback").(*schema.Set); statusCallback.Len() > 0 {
		statusCallback := statusCallback.List()[0].(map[string]interface{})

		addIfNotEmpty(createRequestPayload, "StatusCallbackMethod", statusCallback["http_method"]) // TODO Map to safe values
		addIfNotEmpty(createRequestPayload, "StatusCallback", statusCallback["url"])
	}

	if emergency := d.Get("emergency").(*schema.Set); emergency.Len() > 0 {
		emergency := emergency.List()[0].(map[string]interface{})

		addIfNotEmpty(createRequestPayload, "EmergencyStatus", emergency["status"]) // TODO Map to Twilio values
		addIfNotEmpty(createRequestPayload, "EmergencyAddressSid", emergency["address_sid"])
	}

	return createRequestPayload
}

func hashAnything(item interface{}) int {
	s := cast.ToString(item)
	return hashcode.String(s)
}

func hashStringKeyedMap(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		for _, value := range m {
			buf.WriteString(fmt.Sprintf("%s-", cast.ToString(value)))
		}
	}

	return hashcode.String(buf.String())
}

func mapTwilioPhoneNumberToTerraform(ph *twilio.IncomingPhoneNumber, d *schema.ResourceData) error {
	d.Set("sid", ph.Sid)
	d.Set("number", ph.PhoneNumber.Local())

	d.Set("friendly_name", ph.FriendlyName)
	// d.Set("address_sid", p.AddressSid) -- address SID not in twiliogo
	// d.Set("identity_sid", p.IdentitySid) -- identity SID not in twiliogo
	d.Set("trunk_sid", ph.TrunkSid)

	if ph.DateCreated.Valid {
		d.Set("date_created", ph.DateCreated.Time.Format("2006-01-02T15:04:05-07:00"))
	}

	if ph.DateUpdated.Valid {
		d.Set("date_updated", ph.DateUpdated.Time.Format("2006-01-02T15:04:05-07:00"))
	}
	d.Set("address_requirements", ph.AddressRequirements)
	d.Set("is_beta", ph.Beta)
	d.Set("is_mms_capable", ph.Capabilities.MMS)
	d.Set("is_sms_capable", ph.Capabilities.SMS)
	d.Set("is_voice_capable", ph.Capabilities.Voice)
	// d.Set("is_fax_capable", p.Capabilities.Fax) -- p.Capabilities.Fax not in twiliogo

	// Voice set
	voiceMap := make(map[string]interface{})
	voiceMap["application_sid"] = ph.VoiceApplicationSid
	voiceMap["fallback_url"] = ph.VoiceFallbackURL
	voiceMap["fallback_http_method"] = ph.VoiceFallbackMethod
	voiceMap["primary_url"] = ph.VoiceURL
	voiceMap["primary_http_method"] = ph.VoiceMethod
	voiceMap["caller_id_enabled"] = ph.VoiceCallerIDLookup
	// voiceMap["receive_mode"] = p.ReceiveMode -- receive mode not in twiliogo
	voiceSet := d.Get("voice").(*schema.Set)
	voiceSet.Add(voiceMap)

	// sms set
	smsMap := make(map[string]interface{})
	smsMap["application_sid"] = ph.SMSApplicationSid
	smsMap["fallback_url"] = ph.SMSFallbackURL
	smsMap["fallback_http_method"] = ph.SMSFallbackMethod
	smsMap["primary_url"] = ph.SMSURL
	smsMap["primary_http_method"] = ph.SMSMethod
	smsSet := d.Get("sms").(*schema.Set)
	smsSet.Add(smsMap)

	// status_callback
	statusCallbackMap := make(map[string]interface{})
	statusCallbackMap["url"] = ph.SMSFallbackURL
	statusCallbackMap["http_method"] = ph.SMSFallbackMethod
	statusCallbackSet := d.Get("status_callback").(*schema.Set)
	statusCallbackSet.Add(statusCallbackMap)

	// emergency
	emergencyMap := make(map[string]interface{})

	if ph.EmergencyAddressSid.Valid {
		emergencyMap["address_sid"] = ph.EmergencyAddressSid.String
	}
	emergencyMap["status"] = ph.EmergencyStatus
	emergencySet := d.Get("emergency").(*schema.Set)
	emergencySet.Add(emergencyMap)

	return nil
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

	err = mapTwilioPhoneNumberToTerraform(buyResult, d)

	if err != nil  {
		return fmt.Errorf("Encountered error while reading buy result for phone number SID %s and mapping it to TF: %s", buyResult.Sid, err)
	}

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

	if err != nil {
		return fmt.Errorf("Encountered an error when getting phone number SID %s: %s", sid, err)
	}

	err = mapTwilioPhoneNumberToTerraform(ph, d)

	if err != nil {
		return fmt.Errorf("Encountered an error while mapping Twilio API result to terraform: %s", err)
	}

	return nil
}

func resourceTwilioPhoneNumberUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioPhoneNumberDelete")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()

	updatePayload := makeCreateRequestPayload(d)

	//phoneNumber := d.Get("number").(string)
	//updatePayload.Set("PhoneNumber", e164Number)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
			"phone_sid":   sid,
		},
	).Debug("START client.IncomingNumbers.Update")

	_, err := client.IncomingNumbers.Update(context, sid, updatePayload)

	if err != nil {
		return fmt.Errorf("Failed to update phone number SID %s: %s", sid, err)
	}

	return nil
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
