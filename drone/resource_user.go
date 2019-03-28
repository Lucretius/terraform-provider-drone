package drone

import (
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

	return readUser(data, user, err)
}

func resourceUserUpdate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	user, err := client.User(data.Id())

	client.UserUpdate(user.Login, updateUser(data))
}

func resourceUserRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	user, err := client.User(data.Id())

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

	exists := (user.Login == login) && (err == nil)

	return exists, err
}

func createUser(data *schema.ResourceData) (user *drone.User) {
	user = &drone.User{
		Login: data.Get("login").(string),
	}

	return
}

func updateUser(data *schema.ResourceData) (user *drone.UserPatch) {
	userPatch := &drone.UserPatch{
		Active:  Bool(data.Get("active").(bool)),
		Admin:   Bool(data.Get("admin").(bool)),
		Machine: Bool(data.Get("machine").(bool)),
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
	return nil
}
