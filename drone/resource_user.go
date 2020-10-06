package drone

import (
	"fmt"
	"github.com/Lucretius/terraform-provider-drone/drone/utils"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"login": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"admin": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"machine": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Exists: resourceUserExists,
	}
}

func resourceUserCreate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	user, err := client.UserCreate(createUser(data))

	if err != nil {
		return fmt.Errorf("Unable to create user %s", user.Login)
	}

	return readUser(data, user, err)
}

func resourceUserUpdate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	user, err := client.User(data.Id())

	if err != nil {
		return err
	}

	user, err = client.UserUpdate(user.Login, updateUser(data))

	return readUser(data, user, err)
}

func resourceUserRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	user, err := client.User(data.Id())
	if err != nil {
		return fmt.Errorf("failed to read Drone user with id: %s", data.Id())
	}

	return readUser(data, user, err)
}

func resourceUserDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	return client.UserDelete(data.Id())
}

func resourceUserExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	login := data.Id()

	user, err := client.User(login)
	if err != nil {
		return false, fmt.Errorf("failed to read Drone user with id: %s", data.Id())
	}

	exists := user.Login == login

	return exists, err
}

func createUser(data *schema.ResourceData) (user *drone.User) {
	user = &drone.User{
		Login:   data.Get("login").(string),
		Active:  data.Get("active").(bool),
		Admin:   data.Get("admin").(bool),
		Machine: data.Get("machine").(bool),
	}

	return user
}

func updateUser(data *schema.ResourceData) (user *drone.UserPatch) {
	userPatch := &drone.UserPatch{
		Active:  utils.Bool(data.Get("active").(bool)),
		Admin:   utils.Bool(data.Get("admin").(bool)),
		Machine: utils.Bool(data.Get("machine").(bool)),
	}
	return userPatch
}

func readUser(data *schema.ResourceData, user *drone.User, err error) error {
	if err != nil {
		return err
	}

	data.SetId(user.Login)

	data.Set("login", user.Login)
	data.Set("active", user.Active)
	data.Set("machine", user.Machine)
	data.Set("admin", user.Admin)
	data.Set("token", user.Token)
	return nil
}
