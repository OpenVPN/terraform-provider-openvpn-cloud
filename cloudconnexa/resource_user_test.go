package cloudconnexa

import (
	"errors"
	"fmt"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloudConnexaUser_basic(t *testing.T) {
	rn := "cloudconnexa_user.test"
	user := cloudconnexa.User{
		Username:  acctest.RandStringFromCharSet(10, alphabet),
		FirstName: acctest.RandStringFromCharSet(10, alphabet),
		LastName:  acctest.RandStringFromCharSet(10, alphabet),
		Email:     fmt.Sprintf("terraform-tests+%s@devopenvpn.in", acctest.RandString(10)),
	}
	userChanged := user
	userChanged.Email = fmt.Sprintf("terraform-tests+changed%s@devopenvpn.in", acctest.RandString(10))
	var userID string

	check := func(user cloudconnexa.User) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckCloudConnexaUserExists(rn, &userID),
			resource.TestCheckResourceAttr(rn, "username", user.Username),
			resource.TestCheckResourceAttr(rn, "email", user.Email),
			resource.TestCheckResourceAttr(rn, "first_name", user.FirstName),
			resource.TestCheckResourceAttr(rn, "last_name", user.LastName),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaUserConfig(user),
				Check:  check(user),
			},
			{
				Config: testAccCloudConnexaUserConfig(userChanged),
				Check:  check(userChanged),
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

func testAccCheckCloudConnexaUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cloudconnexa.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudconnexa_user" {
			continue
		}
		username := rs.Primary.Attributes["username"]
		u, err := client.Users.Get(username)
		if err == nil {
			if u != nil {
				return errors.New("user still exists")
			}
		}
	}
	return nil
}

func testAccCheckCloudConnexaUserExists(n string, teamID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		client := testAccProvider.Meta().(*cloudconnexa.Client)
		_, err := client.Users.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCloudConnexaUserImportStateIdFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}
		return rs.Primary.ID, nil
	}
}

func testAccCloudConnexaUserConfig(user cloudconnexa.User) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
	base_url = "https://%s.api.openvpn.com"
}
resource "cloudconnexa_user" "test" {
	username   = "%s"
	email      = "%s"
	first_name = "%s"
	last_name  = "%s"
}
`, testCloudID, user.Username, user.Email, user.FirstName, user.LastName)
}
