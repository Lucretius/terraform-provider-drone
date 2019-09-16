package drone

import (
	"fmt"
	"regexp"

	"github.com/Lucretius/terraform-provider-drone/drone/utils"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceRepo() *schema.Resource {
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
			"trusted": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"protected": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "private",
			},
			"configuration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  ".drone.yml",
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceRepoCreate,
		Read:   resourceRepoRead,
		Update: resourceRepoUpdate,
		Delete: resourceRepoDelete,
		Exists: resourceRepoExists,
	}
}

func resourceRepoCreate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	// Refresh repository list
	if _, err := client.RepoListSync(); err != nil {
		return err
	}

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return err
	}

	resp, err := client.Repo(owner, repo)

	if err != nil {
		return err
	}
	repository, err := client.RepoUpdate(owner, repo, createRepo(data))

	if err != nil {
		return err
	}
	if !resp.Active {
		_, err = client.RepoEnable(owner, repo)

		if err != nil {
			return err
		}
	}

	return readRepo(data, repository, err)
}

func resourceRepoRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Id())

	if err != nil {
		return err
	}

	repository, err := client.Repo(owner, repo)

	return readRepo(data, repository, err)
}

func resourceRepoUpdate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return err
	}

	repository, err := client.RepoUpdate(owner, repo, createRepo(data))

	return readRepo(data, repository, err)
}

func resourceRepoDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Id())

	if err != nil {
		return err
	}

	return client.RepoDisable(owner, repo)
}

func resourceRepoExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Id())

	if err != nil {
		return false, err
	}

	repository, err := client.Repo(owner, repo)

	exists := (repository.Namespace == owner) && (repository.Name == repo) && (err == nil)

	return exists, err
}

func createRepo(data *schema.ResourceData) (repository *drone.RepoPatch) {
	trusted := data.Get("trusted").(bool)
	protected := data.Get("protected").(bool)
	timeout := int64(data.Get("timeout").(int))
	visibility := data.Get("visibility").(string)
	config := data.Get("configuration").(string)

	repository = &drone.RepoPatch{
		Config:     &config,
		Protected:  &protected,
		Trusted:    &trusted,
		Timeout:    &timeout,
		Visibility: &visibility,
	}

	return
}

func readRepo(data *schema.ResourceData, repository *drone.Repo, err error) error {
	if err != nil {
		return err
	}

	data.SetId(fmt.Sprintf("%s/%s", repository.Namespace, repository.Name))

	data.Set("repository", fmt.Sprintf("%s/%s", repository.Namespace, repository.Name))
	data.Set("trusted", repository.Trusted)
	data.Set("protected", repository.Protected)
	data.Set("timeout", repository.Timeout)
	data.Set("visibility", repository.Visibility)
	data.Set("configuration", repository.Config)

	return nil
}
