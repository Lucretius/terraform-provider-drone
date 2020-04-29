package drone

import (
	"fmt"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestOrgSecret(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testProviders,
		CheckDestroy: testOrgSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testOrgSecretConfigBasic(
					"foo",
					"password",
					"1234567890",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"drone_orgsecret.secret",
						"namespace",
						"foo",
					),
					resource.TestCheckResourceAttr(
						"drone_orgsecret.secret",
						"name",
						"password",
					),
					resource.TestCheckResourceAttr(
						"drone_orgsecret.secret",
						"value",
						"1234567890",
					),
				),
			},
		},
	})
}

func testOrgSecretConfigBasic(namespace, name, value string) string {
	return fmt.Sprintf(`
	resource "drone_orgsecret" "secret" {
		namespace = "%s"
		name      = "%s"
		value     = "%s"
	}
	`, namespace, name, value)
}

func testOrgSecretDestroy(state *terraform.State) error {
	client := testProvider.Meta().(drone.Client)

	for _, resource := range state.RootModule().Resources {
		if resource.Type != "drone_orgsecret" {
			continue
		}

		namespace := resource.Primary.Attributes["namespace"]
		name := resource.Primary.Attributes["name"]

		err := client.OrgSecretDelete(namespace, name)
		if err == nil {
			return fmt.Errorf("namespace Secret still exists: %s/%s", namespace, name)
		}
	}

	return nil
}