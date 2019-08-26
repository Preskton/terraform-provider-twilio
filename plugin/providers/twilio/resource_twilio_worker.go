package twilio

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceTwilioWorker() *schema.Resource {
	return &schema.Resource{
		Create: resourceTwilioWorkerCreate,
		Read:   resourceTwilioWorkerRead,
		Update: resourceTwilioWorkerUpdate,
		Delete: resourceTwilioWorkerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_sid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"friendly_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func flattenWorkerForCreate(d *schema.ResourceData) url.Values {
	v := make(url.Values)

	v.Add("FriendlyName", d.Get("friendly_name").(string))
	return v
}

func resourceTwilioWorkerCreate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioWorkerCreate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	workspaceSid := d.Get("workspace_sid").(string)
	createParams := flattenWorkerForCreate(d)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.TaskRouter.Workspace.Workers.Create")

	worker, err := client.TaskRouter.Workspace(workspaceSid).Workers.Create(context, createParams)
	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.TaskRouter.Workspace.Workers.Create failed")

		return err
	}
	d.SetId(worker.Sid)
	d.Set("friendly_name", worker.FriendlyName)
	d.Set("date_created", worker.DateCreated)
	d.Set("date_updated", worker.DateUpdated)
	return nil
}

func resourceTwilioWorkerRead(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioWorkerRead")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()
	workspaceSid := d.Get("workspace_sid").(string)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.TaskRouter.Workspace.Workers.Get")

	worker, err := client.TaskRouter.Workspace(workspaceSid).Workers.Get(context, sid)
	if err != nil {
		log.WithFields(
			log.Fields{
				"parent_account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.TaskRouter.Workspace.Workers.Get failed")

		return err
	}
	d.Set("friendly_name", worker.FriendlyName)
	d.Set("date_created", worker.DateCreated)
	d.Set("date_updated", worker.DateUpdated)
	return nil
}

func resourceTwilioWorkerUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTwilioWorkerDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioWorkerDelete")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()
	workspaceSid := d.Get("workspace_sid").(string)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
			"queue_sid":   sid,
		},
	).Debug("START client.TaskRouter.Workspace.Workers.Delete")

	err := client.TaskRouter.Workspace(workspaceSid).Workers.Delete(context, sid)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
			"queue_sid":   sid,
		},
	).Debug("END client.TaskRouter.Workspace.Workers.Delete")
	if err != nil {
		return fmt.Errorf("Failed to delete worker: %s", err.Error())
	}
	return nil
}
