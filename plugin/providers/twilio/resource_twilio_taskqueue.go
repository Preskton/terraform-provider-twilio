package twilio

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceTwilioTaskQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceTwilioTaskQueueCreate,
		Read:   resourceTwilioTaskQueueRead,
		Update: resourceTwilioTaskQueueUpdate,
		Delete: resourceTwilioTaskQueueDelete,
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

func flattentaskQueueForCreate(d *schema.ResourceData) url.Values {
	v := make(url.Values)

	v.Add("FriendlyName", d.Get("friendly_name").(string))
	return v
}

func resourceTwilioTaskQueueCreate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioTaskQueueCreate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	workspaceSid := d.Get("workspace_sid").(string)
	createParams := flattentaskQueueForCreate(d)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.TaskRouter.Workspace.Queues.Create")

	taskQueue, err := client.TaskRouter.Workspace(workspaceSid).Queues.Create(context, createParams)
	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.TaskRouter.Workspace.Queues.Create failed")

		return err
	}
	d.SetId(taskQueue.Sid)
	d.Set("friendly_name", taskQueue.FriendlyName)
	d.Set("date_created", taskQueue.DateCreated)
	d.Set("date_updated", taskQueue.DateUpdated)
	d.Set("target_workers", taskQueue.TargetWorkers)
	d.Set("task_order", taskQueue.TaskOrder)
	d.Set("url", taskQueue.URL)
	return nil
}

func resourceTwilioTaskQueueRead(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioTaskQueueRead")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()
	workspaceSid := d.Get("workspace_sid").(string)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.TaskRouter.Workspace.Queues.Get")

	taskQueue, err := client.TaskRouter.Workspace(workspaceSid).Queues.Get(context, sid)
	if err != nil {
		log.WithFields(
			log.Fields{
				"parent_account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.TaskRouter.Workspace.Queues.Get failed")

		return err
	}
	d.SetId(taskQueue.Sid)
	d.Set("friendly_name", taskQueue.FriendlyName)
	d.Set("date_created", taskQueue.DateCreated)
	d.Set("date_updated", taskQueue.DateUpdated)
	d.Set("target_workers", taskQueue.TargetWorkers)
	d.Set("task_order", taskQueue.TaskOrder)
	d.Set("url", taskQueue.URL)
	return nil
}

func resourceTwilioTaskQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTwilioTaskQueueDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioTaskQueueDelete")

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
	).Debug("START client.TaskRouter.Workspace.Queues.Delete")

	err := client.TaskRouter.Workspace(workspaceSid).Queues.Delete(context, sid)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
			"queue_sid":   sid,
		},
	).Debug("END client.TaskRouter.Workspace.Queues.Delete")
	if err != nil {
		return fmt.Errorf("Failed to delete taskQueue: %s", err.Error())
	}
	return nil
}
