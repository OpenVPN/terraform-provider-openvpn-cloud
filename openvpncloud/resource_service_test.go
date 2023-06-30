package openvpncloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"

	"github.com/OpenVPN/terraform-provider-openvpn-cloud/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOpenvpncloudService_basic(t *testing.T) {
	rn := "openvpncloud_service.test"
	networkName := acctest.RandStringFromCharSet(10, alphabet)
	service := client.Service{
		Name: acctest.RandStringFromCharSet(10, alphabet),
	}
	serviceChanged := service
	serviceChanged.Name = fmt.Sprintf("changed-%s", acctest.RandStringFromCharSet(10, alphabet))

	check := func(service client.Service) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckOpenvpncloudServiceExists(rn, networkName),
			resource.TestCheckResourceAttr(rn, "name", service.Name),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckOpenvpncloudServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenvpncloudServiceConfig(service, networkName),
				Check:  check(service),
			},
			{
				Config: testAccOpenvpncloudServiceConfig(serviceChanged, networkName),
				Check:  check(serviceChanged),
			},
		},
	})
}

func testAccCheckOpenvpncloudServiceExists(rn, networkId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		c := testAccProvider.Meta().(*client.Client)
		_, err := c.GetService(rs.Primary.ID, rs.Primary.Attributes["network_item_type"], rs.Primary.Attributes["network_item_id"])
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckOpenvpncloudServiceDestroy(state *terraform.State) error {
	c := testAccProvider.Meta().(*client.Client)
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "openvpncloud_service" {
			continue
		}
		id := rs.Primary.Attributes["id"]
		s, err := c.GetService(id, "c63acae0-b569-4116-9b39-921c1dee62d2", "NETWORK")
		if err == nil || s != nil {
			return fmt.Errorf("service still exists")
		}
	}
	return nil
}

func testAccOpenvpncloudServiceConfig(service client.Service, networkName string) string {
	return fmt.Sprintf(`
provider "openvpncloud" {
	base_url = "https://%s.api.openvpn.com"
}

resource "openvpncloud_network" "test" {
	name = "%s"

	default_connector {
	  name          = "%s"
	  vpn_region_id = "fi-hel"
	}
	default_route {
	  value = "10.1.2.0/24"
	  type  = "IP_V4"
	}
}

resource "openvpncloud_service" "test" {
	name = "%s"
	type = "SERVICE_DESTINATION"
	description = "test"
	network_item_type = "NETWORK"
	network_item_id = openvpncloud_network.test.id
	routes = ["test.ua" ]
	config {
		service_types = ["ANY"]
	}
}
`, testCloudID, networkName, fmt.Sprintf("connector_%s", networkName), service.Name)
}
