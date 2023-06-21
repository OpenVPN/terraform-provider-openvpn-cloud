package openvpncloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/OpenVPN/terraform-provider-openvpn-cloud/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenvpncloudUserGroup_basic(t *testing.T) {
	rn := "openvpncloud_user_group.test"
	userGroup := client.UserGroup{
		Name: acctest.RandStringFromCharSet(10, alphabet),
		VpnRegionIds: []string{
			"us-east-1",
		},
	}
	userGroupChanged := userGroup
	userGroupChanged.Name = fmt.Sprintf("changed-%s", acctest.RandStringFromCharSet(10, alphabet))

	check := func(userGroup client.UserGroup) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckOpenvpncloudUserGroupExists(rn),
			resource.TestCheckResourceAttr(rn, "name", userGroup.Name),
			resource.TestCheckResourceAttr(rn, "vpn_region_ids.0", userGroup.VpnRegionIds[0]),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckOpenvpncloudUserGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenvpncloudUserGroupConfig(userGroup),
				Check:  check(userGroup),
			},
			{
				Config: testAccOpenvpncloudUserGroupConfig(userGroupChanged),
				Check:  check(userGroupChanged),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: testAccOpenvpncloudUserImportStateIdFunc(rn),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckOpenvpncloudUserGroupDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openvpncloud_user_group" {
			continue
		}
		username := rs.Primary.Attributes["username"]
		u, err := c.GetUserGroupById(username)
		if err == nil {
			if u != nil {
				return errors.New("user still exists")
			}
		}
	}
	return nil
}

func testAccCheckOpenvpncloudUserGroupExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		c := testAccProvider.Meta().(*client.Client)
		_, err := c.GetUserGroupById(rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccOpenvpncloudUserGroupConfig(userGroup client.UserGroup) string {
	idsStr, _ := json.Marshal(userGroup.VpnRegionIds)

	return fmt.Sprintf(`
provider "openvpncloud" {
	base_url = "https://%s.api.openvpn.com"
}
resource "openvpncloud_user_group" "test" {
	name = "%s"
	vpn_region_ids = %s

}
`, testCloudID, userGroup.Name, idsStr)
}
