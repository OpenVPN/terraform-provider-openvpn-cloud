package openvpncloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/OpenVPN/terraform-provider-openvpn-cloud/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenvpncloudUser_basic(t *testing.T) {
	rn := "openvpncloud_user.test"
	user := client.User{
		Username:  acctest.RandStringFromCharSet(10, alphabet),
		FirstName: acctest.RandStringFromCharSet(10, alphabet),
		LastName:  acctest.RandStringFromCharSet(10, alphabet),
		Email:     fmt.Sprintf("terraform-tests+%s@devopenvpn.in", acctest.RandString(10)),
	}
	userChanged := user
	userChanged.Email = fmt.Sprintf("terraform-tests+changed%s@devopenvpn.in", acctest.RandString(10))
	var userID string

	check := func(user client.User) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckOpenvpncloudUserExists(rn, &userID),
			resource.TestCheckResourceAttr(rn, "username", user.Username),
			resource.TestCheckResourceAttr(rn, "email", user.Email),
			resource.TestCheckResourceAttr(rn, "first_name", user.FirstName),
			resource.TestCheckResourceAttr(rn, "last_name", user.LastName),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckOpenvpncloudUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenvpncloudUserConfig(user),
				Check:  check(user),
			},
			{
				Config: testAccOpenvpncloudUserConfig(userChanged),
				Check:  check(userChanged),
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

func testAccCheckOpenvpncloudUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openvpncloud_user" {
			continue
		}
		username := rs.Primary.Attributes["username"]
		u, err := client.GetUserById(username)
		if err == nil {
			if u != nil {
				return errors.New("user still exists")
			}
		}
	}
	return nil
}

func testAccCheckOpenvpncloudUserExists(n string, teamID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		client := testAccProvider.Meta().(*client.Client)
		_, err := client.GetUserById(rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccOpenvpncloudUserImportStateIdFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}
		return rs.Primary.ID, nil
	}
}

func testAccOpenvpncloudUserConfig(user client.User) string {
	return fmt.Sprintf(`
provider "openvpncloud" {
	base_url = "https://%s.api.openvpn.com"
}
resource "openvpncloud_user" "test" {
	username   = "%s"
	email      = "%s"
	first_name = "%s"
	last_name  = "%s"
}
`, testCloudID, user.Username, user.Email, user.FirstName, user.LastName)
}
