package drone

import (
	"fmt"
	
	"github.com/Lucretius/terraform-provider-drone/drone/utils"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOrgSecret() *schema.Resource {
	return &schema.Resource {
		Schema: map[string]*schema.Schema {
			"namespace": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type: schema.TypeString,
				Required: true,
				Sensitive: true,
			},
			"allow_on_pull_request": {
				Type: schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"allow_push_on_pull_request": {
				Type: schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceOrgSecretCreate,
		Read: resourceOrgSecretRead,
		Update: resourceOrgSecretUpdate,
		Delete: resourceOrgSecretDelete,
		Exists: resourceOrgSecretExists,
	}
}

func resourceOrgSecretCreate(data *schema.ResourceData, meta interface {}) error {
	client := meta.(drone.Client)

	namespace := data.Get("namespace").(string)
	secret, err := client.OrgSecretCreate(namespace, createOrgSecret(data))
	if err != nil {
		return err
	}

	readOrgSecret(data, namespace, secret)
	return nil
}

func resourceOrgSecretRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	namespace, name, err := utils.ParseOrgId(data.Id(), "secret_password")
	if err != nil {
		return err
	}

	secret, err := client.OrgSecret(namespace, name)
	if err != nil {
		return err
	}

	readOrgSecret(data, namespace, secret)
	return nil
}

func resourceOrgSecretUpdate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	namespace, _, err := utils.ParseOrgId(data.Id(), "secret_password")
	if err != nil {
		return err
	}

	secret, err := client.OrgSecretUpdate(namespace, createOrgSecret(data))
	if err != nil {
		return err
	}

	data.Set("value", data.Get("value").(string))
	readOrgSecret(data, namespace, secret)
	return nil
}

func resourceOrgSecretDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	namespace, name, err := utils.ParseOrgId(data.Id(), "secret_password")
	if err != nil {
		return err
	}

	return client.OrgSecretDelete(namespace, name)
}

func resourceOrgSecretExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	namespace, name, err := utils.ParseOrgId(data.Id(), "secret_password")
	if err != nil {
		return false, err
	}

	secret, err := client.OrgSecret(namespace, name)
	if err != nil {
		return false, err
	}

	exists := (secret.Name == name) && (err == nil)
	return exists, err
}

func createOrgSecret(data *schema.ResourceData) (secret *drone.Secret) {
	return &drone.Secret{
		Name: data.Get("name").(string),
		Data: data.Get("value").(string),
		PullRequest: data.Get("allow_on_pull_request").(bool),
		PullRequestPush: data.Get("allow_push_on_pull_request").(bool),
	}
}

func readOrgSecret(data *schema.ResourceData, namespace string, secret *drone.Secret) {
	data.SetId(fmt.Sprintf("%s/%s", namespace, secret.Name))
	data.Set("name", secret.Name)
	data.Set("allow_on_pull_request", secret.PullRequest)
	data.Set("allow_push_on_pull_request", secret.PullRequestPush)
}