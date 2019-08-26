package twilio

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceTwilioApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceTwilioApplicationCreate,
		Read:   resourceTwilioApplicationRead,
		Update: resourceTwilioApplicationUpdate,
		Delete: resourceTwilioApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"friendly_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func flattenApplicationForCreate(d *schema.ResourceData) url.Values {
	v := make(url.Values)

	v.Add("FriendlyName", d.Get("friendly_name").(string))
	return v
}

func resourceTwilioApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioApplicationCreate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	createParams := flattenApplicationForCreate(d)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.Applications.Create")

	application, err := client.Applications.Create(context, createParams)
	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.Applications.Create failed")

		return err
	}
	d.SetId(application.Sid)
	d.Set("friendly_name", application.FriendlyName)
	d.Set("date_created", application.DateCreated)
	d.Set("date_updated", application.DateUpdated)
	d.Set("voice_url", application.VoiceURL)
	d.Set("sms_url", application.SMSURL)
	return nil
}

func resourceTwilioApplicationRead(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioApplicationRead")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.Applications.Get")

	application, err := client.Applications.Get(context, sid)
	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.Applications.Get failed")

		return err
	}
	d.Set("friendly_name", application.FriendlyName)
	d.Set("date_created", application.DateCreated)
	d.Set("date_updated", application.DateUpdated)
	d.Set("voice_url", application.VoiceURL)
	d.Set("sms_url", application.SMSURL)
	return nil
}

func resourceTwilioApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTwilioApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioApplicationDelete")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()

	log.WithFields(
		log.Fields{
			"account_sid":     config.AccountSID,
			"application_sid": sid,
		},
	).Debug("START client.Applications.Delete")

	err := client.Applications.Delete(context, sid)

	log.WithFields(
		log.Fields{
			"account_sid":     config.AccountSID,
			"application_sid": sid,
		},
	).Debug("END client.Applications.Delete")
	if err != nil {
		return fmt.Errorf("Failed to delete application: %s", err.Error())
	}
	return nil
}
