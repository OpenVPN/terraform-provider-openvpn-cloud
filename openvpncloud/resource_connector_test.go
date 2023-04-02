package openvpncloud

import (
	"fmt"
	"testing"

	"github.com/OpenVPN/terraform-provider-openvpn-cloud/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenvpncloudConnector_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("test-connector")
	resourceName := "openvpncloud_connector.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckOpenvpncloudConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenvpncloudConnectorConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenvpncloudConnectorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpn_region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "network_item_type"),
					resource.TestCheckResourceAttrSet(resourceName, "network_item_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_v4_address"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_v6_address"),
				),
			},
		},
	})
}

func testAccCheckOpenvpncloudConnectorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No connector ID is set")
		}
		return nil
	}
}

func testAccCheckOpenvpncloudConnectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openvpncloud_connector" {
			continue
		}

		connectorId := rs.Primary.ID
		connector, err := client.GetConnectorById(connectorId)

		if err != nil {
			return err
		}

		if connector != nil {
			return fmt.Errorf("Connector with ID '%s' still exists", connectorId)
		}
	}

	return nil
}

func testAccOpenvpncloudConnectorConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "openvpncloud_connector" "test" {
  name             = "%s"
  vpn_region_id    = "us_east"
  network_item_type = "HOST"
  network_item_id   = "example_network_item_id"
}
`, rName)
}
