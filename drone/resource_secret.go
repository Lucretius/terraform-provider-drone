package drone

import (
	"fmt"
	"regexp"

	"github.com/Lucretius/terraform-provider-drone/drone/utils"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[^/ ]+/[^/ ]+$"),
					"Invalid repository (e.g. octocat/hello-world)",
				),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"allow_on_pull_request": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceSecretCreate,
		Read:   resourceSecretRead,
		Update: resourceSecretUpdate,
		Delete: resourceSecretDelete,
		Exists: resourceSecretExists,
	}
}

func resourceSecretCreate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return err
	}

	secret, err := client.SecretCreate(owner, repo, createSecret(data))

	data.Set("value", data.Get("value").(string))

	return readSecret(data, owner, repo, secret, err)
}

func resourceSecretRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, name, err := utils.ParseId(data.Id(), "secret_password")

	if err != nil {
		return err
	}

	secret, err := client.Secret(owner, repo, name)
	if err != nil {
		return fmt.Errorf("failed to read Drone Secret: %s/%s/%s", owner, repo, name)
	}

	return readSecret(data, owner, repo, secret, err)
}

func resourceSecretUpdate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return err
	}

	secret, err := client.SecretUpdate(owner, repo, createSecret(data))

	data.Set("value", data.Get("value").(string))

	return readSecret(data, owner, repo, secret, err)
}

func resourceSecretDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, name, err := utils.ParseId(data.Id(), "secret_password")

	if err != nil {
		return err
	}

	return client.SecretDelete(owner, repo, name)
}

func resourceSecretExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	owner, repo, name, err := utils.ParseId(data.Id(), "secret_password")

	if err != nil {
		return false, err
	}

	secret, err := client.Secret(owner, repo, name)
	if err != nil {
		return false, fmt.Errorf("failed to read Drone Secret: %s/%s/%s", owner, repo, name)
	}

	exists := secret.Name == name

	return exists, err
}

func createSecret(data *schema.ResourceData) (secret *drone.Secret) {
	secret = &drone.Secret{
		Name:        data.Get("name").(string),
		Data:        data.Get("value").(string),
		PullRequest: data.Get("allow_on_pull_request").(bool),
	}

	return
}

func readSecret(data *schema.ResourceData, owner, repo string, secret *drone.Secret, err error) error {
	if err != nil {
		return err
	}

	data.SetId(fmt.Sprintf("%s/%s/%s", owner, repo, secret.Name))

	data.Set("repository", fmt.Sprintf("%s/%s", owner, repo))
	data.Set("name", secret.Name)
	data.Set("allow_on_pull_request", secret.PullRequest)

	return nil
}
