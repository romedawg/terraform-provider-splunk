package splunk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/splunk/go-splunk-client/pkg/client"
	"github.com/splunk/go-splunk-client/pkg/entry"
)

const legacyClientAdminSAMLGroupsInput = `
resource "splunk_admin_saml_groups" "saml_group" {
  name = "new-saml-group"
  roles = ["admin", "power"]
  # default based on provider configuration default
  # use_legacy_client = true
}
`

const adminSAMLGroupsInput = `
resource "splunk_admin_saml_groups" "saml_group" {
  name = "new-saml-group"
  roles = ["admin", "power"]
  use_legacy_client = false
}
`

const legacyClientUpdateAdminSAMLGroupsInput = `
resource "splunk_admin_saml_groups" "saml_group" {
  name = "new-saml-group"
  roles = ["admin", "power", "user"]
  # default based on provider configuration default
  # use_legacy_client = true
}
`

const updateAdminSAMLGroupsInput = `
resource "splunk_admin_saml_groups" "saml_group" {
  name = "new-saml-group"
  roles = ["admin", "power", "user"]
  use_legacy_client = false
}
`

// noResourcesInput needs to be an empty configuration, but not an empty string,
// otherwise the test case thinks it's unset.
const noResourcesInput = ` `

func TestAccSplunkAdminSAMLGroups(t *testing.T) {
	resourceName := "splunk_admin_saml_groups.saml_group"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkAdminSAMLGroupsInputDestroyResources,
		Steps: []resource.TestStep{
			{
				// initial resource uses legacy client
				Config: legacyClientAdminSAMLGroupsInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "new-saml-group"),
					resource.TestCheckResourceAttr(resourceName, "name", "new-saml-group"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "roles.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.1", "power"),
				),
			},
			{
				// migrate resource to new client
				Config: adminSAMLGroupsInput,
				Check: resource.ComposeTestCheckFunc(
					// only "id" is expected to change, from the simple name to the actual URL ID
					resource.TestCheckResourceAttr(resourceName, "id", "https://localhost:8089/services/admin/SAML-groups/new-saml-group"),
					resource.TestCheckResourceAttr(resourceName, "name", "new-saml-group"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "roles.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.1", "power"),
				),
			},
			{
				// migrate back to legacy client to ensure it's not a lock-in situation
				Config: legacyClientAdminSAMLGroupsInput,
				Check: resource.ComposeTestCheckFunc(
					// "id" should be back to the simple name
					resource.TestCheckResourceAttr(resourceName, "id", "new-saml-group"),
					resource.TestCheckResourceAttr(resourceName, "name", "new-saml-group"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "roles.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.1", "power"),
				),
			},
			{
				// for the legacy client, test re-creation of remotely deleted or missing resources, delete the new saml group before updating it
				Config: legacyClientUpdateAdminSAMLGroupsInput,
				PreConfig: func() {
					c, _ := newTestExternalClient()
					if err := c.Delete(entry.SAMLGroup{
						ID: client.ID{
							Title: "new-saml-group",
						},
					}); err != nil {
						t.Error("PreConfig deletion of new-saml-group failed")
					}
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-saml-group"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "roles.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.1", "power"),
					resource.TestCheckResourceAttr(resourceName, "roles.2", "user"),
				),
			},
			{
				// with existing state using the legacy client, test delete.
				// validation of the success of deletion is done in the step below.
				Config: noResourcesInput,
			},
			{
				// new initial resource uses new client
				Config: adminSAMLGroupsInput,
				// but first validate that the resource was previously deleted.
				// it's odd to effectively "test" the previous step here, but there doesn't seem to be
				// a way to 1) test for resource absense or 2) run a post-config step.
				PreConfig: func() {
					c, _ := newTestExternalClient()
					if err := c.Read(entry.SAMLGroup{
						ID: client.ID{
							Title: "new-saml-group",
						},
					}); err != nil {
						if clientError, ok := err.(client.Error); ok {
							if clientError.Code == client.ErrorNotFound {
								return
							}
						}
						t.Error("PreConfig check of previous deletion of new-saml-group failed")
					}
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "https://localhost:8089/services/admin/SAML-groups/new-saml-group"),
					resource.TestCheckResourceAttr(resourceName, "name", "new-saml-group"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "roles.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.1", "power"),
				),
			},
			{
				// for the new client, test re-creation of remotely deleted or missing resources, delete the new saml group before updating it
				Config: updateAdminSAMLGroupsInput,
				PreConfig: func() {
					c, _ := newTestExternalClient()
					if err := c.Delete(entry.SAMLGroup{
						ID: client.ID{
							Title: "new-saml-group",
						},
					}); err != nil {
						t.Error("PreConfig deletion of new-saml-group failed")
					}
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-saml-group"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "roles.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.1", "power"),
					resource.TestCheckResourceAttr(resourceName, "roles.2", "user"),
				),
			},
			{
				// imports always use the new client, as the splunk resource can't indicate which client should be used
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"use_legacy_client"},
			},
		},
	})
}

func testAccSplunkAdminSAMLGroupsInputDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_authorization_roles":
			endpoint := client.BuildSplunkURL(nil, "services", "admin", "SAML-groups", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
