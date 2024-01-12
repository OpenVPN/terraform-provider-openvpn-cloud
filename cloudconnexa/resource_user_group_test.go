package cloudconnexa

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloudConnexaUserGroup_basic(t *testing.T) {
	rn := "cloudconnexa_user_group.test"
	userGroup := cloudconnexa.UserGroup{
		Name: acctest.RandStringFromCharSet(10, alphabet),
		VpnRegionIds: []string{
			"us-east-1",
		},
	}
	userGroupChanged := userGroup
	userGroupChanged.Name = fmt.Sprintf("changed-%s", acctest.RandStringFromCharSet(10, alphabet))

	check := func(userGroup cloudconnexa.UserGroup) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckCloudConnexaUserGroupExists(rn),
			resource.TestCheckResourceAttr(rn, "name", userGroup.Name),
			resource.TestCheckResourceAttr(rn, "vpn_region_ids.0", userGroup.VpnRegionIds[0]),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaUserGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaUserGroupConfig(userGroup),
				Check:  check(userGroup),
			},
			{
				Config: testAccCloudConnexaUserGroupConfig(userGroupChanged),
				Check:  check(userGroupChanged),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: testAccCloudConnexaUserImportStateIdFunc(rn),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudConnexaUserGroupDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*cloudconnexa.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudconnexa_user_group" {
			continue
		}
		username := rs.Primary.Attributes["username"]
		u, err := c.UserGroups.GetByName(username)
		if err == nil {
			if u != nil {
				return errors.New("user still exists")
			}
		}
	}
	return nil
}

func testAccCheckCloudConnexaUserGroupExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		c := testAccProvider.Meta().(*cloudconnexa.Client)
		_, err := c.UserGroups.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCloudConnexaUserGroupConfig(userGroup cloudconnexa.UserGroup) string {
	idsStr, _ := json.Marshal(userGroup.VpnRegionIds)

	return fmt.Sprintf(`
provider "cloudconnexa" {
	base_url = "https://%s.api.openvpn.com"
}
resource "cloudconnexa_user_group" "test" {
	name = "%s"
	vpn_region_ids = %s

}
`, testCloudID, userGroup.Name, idsStr)
}
