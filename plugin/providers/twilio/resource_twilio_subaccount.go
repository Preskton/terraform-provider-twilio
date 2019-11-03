package twilio

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"

	log "github.com/sirupsen/logrus"
)

func resourceTwilioSubaccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTwilioSubaccountCreate,
		Read:   resourceTwilioSubaccountRead,
		Update: resourceTwilioSubaccountUpdate,
		Delete: resourceTwilioSubaccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"parent_account_sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
			},
			"auth_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func flattenSubaccountForCreate(d *schema.ResourceData) url.Values {
	v := make(url.Values)

	v.Add("FriendlyName", d.Get("friendly_name").(string))

	return v
}

func flattenSubaccountForDelete(d *schema.ResourceData) url.Values {
	v := make(url.Values)

	v.Add("status", "closed")

	return v
}

func resourceTwilioSubaccountCreate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioSubaccountCreate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	createParams := flattenSubaccountForCreate(d)

	log.WithFields(
		log.Fields{
			"parent_account_sid": config.AccountSID,
		},
	).Debug("START client.AccountsCreate")

	createResult, err := client.Accounts.Create(context, createParams)

	if err != nil {
		log.WithFields(
			log.Fields{
				"parent_account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.AccountsCreate failed")

		return err
	}

	d.SetId(createResult.Sid)
	d.Set("status", createResult.Status)
	d.Set("auth_token", createResult.AuthToken)
	d.Set("friendly_name", createResult.FriendlyName) // In the event that the name wasn't specified, Twilio generates one for you
	d.Set("date_created", createResult.DateCreated)
	d.Set("date_updated", createResult.DateUpdated)
	d.Set("parent_account_sid", createResult.OwnerAccountSid)

	log.WithFields(
		log.Fields{
			"parent_account_sid": config.AccountSID,
			"subaccount_sid":     createResult.Sid,
		},
	).Debug("END client.AccountsCreate")

	return nil
}

func resourceTwilioSubaccountRead(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioSubaccountRead")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()

	log.WithFields(
		log.Fields{
			"parent_account_sid": config.AccountSID,
			"subaccount_sid":     sid,
		},
	).Debug("START client.Accounts.Get")

	account, err := client.Accounts.Get(context, sid)

	d.Set("status", account.Status)
	d.Set("auth_token", account.AuthToken)
	d.Set("friendly_name", account.FriendlyName) // In the event that the name wasn't specified, Twilio generates one for you
	d.Set("date_created", account.DateCreated)
	d.Set("date_updated", account.DateUpdated)
	d.Set("parent_account_sid", account.OwnerAccountSid)

	log.WithFields(
		log.Fields{
			"parent_account_sid": config.AccountSID,
			"subaccount_sid":     sid,
		},
	).Debug("END client.AccountsGet")

	if err != nil {
		return fmt.Errorf("Failed to refresh account: %s", err.Error())
	}

	return nil
}

func resourceTwilioSubaccountUpdate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("Not implemented")
}

func resourceTwilioSubaccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioSubaccountDelete")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()

	updateData := flattenSubaccountForDelete(d)

	log.WithFields(
		log.Fields{
			"parent_account_sid": config.AccountSID,
			"subaccount_sid":     sid,
		},
	).Debug("START client.Accounts.Delete")

	_, err := client.Accounts.Update(context, sid, updateData)

	log.WithFields(
		log.Fields{
			"parent_account_sid": config.AccountSID,
			"subaccount_sid":     sid,
		},
	).Debug("END client.Accounts.Delete")

	if err != nil {
		return fmt.Errorf("Failed to delete account: %s", err.Error())
	}

	return nil
}
