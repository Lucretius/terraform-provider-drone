package drone

import (
	"fmt"
	"regexp"

	"github.com/Lucretius/terraform-provider-drone/drone/utils"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceCron() *schema.Resource {
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
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"event": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					drone.EventPush,
				}, false),
			},
			"branch": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "master",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"expr": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "@monthly",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"@hourly",
					"@daily",
					"@weekly",
					"@monthly",
					"@yearly",
				}, false),
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceCronCreate,
		Read:   resourceCronRead,
		Update: resourceCronUpdate,
		Delete: resourceCronDelete,
		Exists: resourceCronExists,
	}
}

func resourceCronCreate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return err
	}

	cron, err := client.CronCreate(owner, repo, createCron(data))

	if err != nil {
		return err
	}

	return readCron(data, cron, owner, repo, err)
}

func resourceCronRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return err
	}
	cronName := data.Get("name").(string)
	cron, err := client.Cron(owner, repo, cronName)
	if err != nil {
		return fmt.Errorf("failed to read Drone Cron: %s/%s/%s", owner, repo, cronName)
	}

	return readCron(data, cron, owner, repo, err)
}

func resourceCronUpdate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return err
	}
	cronName := data.Get("name").(string)
	cron, err := client.CronUpdate(owner, repo, cronName, updateCron(data))

	return readCron(data, cron, owner, repo, err)
}

func resourceCronDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return err
	}
	cronName := data.Get("name").(string)
	return client.CronDelete(owner, repo, cronName)
}

func resourceCronExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return false, err
	}
	cronName := data.Get("name").(string)
	cron, err := client.Cron(owner, repo, cronName)
	if err != nil {
		return false, fmt.Errorf("failed to read Drone Cron: %s/%s/%s", owner, repo, cronName)
	}

	exists := cron.Name == cronName
	return exists, err
}

func createCron(data *schema.ResourceData) (repository *drone.Cron) {
	branch := data.Get("branch").(string)
	disabled := data.Get("disabled").(bool)
	expr := data.Get("expr").(string)
	name := data.Get("name").(string)
	event := data.Get("event").(string)

	cron := &drone.Cron{
		Disabled: disabled,
		Branch:   branch,
		Expr:     expr,
		Event:    event,
		Name:     name,
	}

	return cron
}

func updateCron(data *schema.ResourceData) (repository *drone.CronPatch) {
	branch := data.Get("branch").(string)
	disabled := data.Get("disabled").(bool)
	event := data.Get("event").(string)

	cron := &drone.CronPatch{
		Disabled: utils.Bool(disabled),
		Branch:   &branch,
		Event:    &event,
	}
	return cron
}

func readCron(data *schema.ResourceData, cron *drone.Cron, namespace string, repo string, err error) error {
	if err != nil {
		return err
	}

	data.SetId(fmt.Sprintf("%d", cron.ID))

	data.Set("repository", fmt.Sprintf("%s/%s", namespace, repo))
	data.Set("branch", cron.Branch)
	data.Set("disabled", cron.Disabled)
	data.Set("expr", cron.Expr)
	data.Set("name", cron.Name)

	return nil
}
