package cloudconnexa

import (
	"fmt"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloudConnexaConnector_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("test-connector")
	resourceName := "cloudconnexa_connector.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaConnectorConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudConnexaConnectorExists(resourceName),
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

func testAccCheckCloudConnexaConnectorExists(n string) resource.TestCheckFunc {
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

func testAccCheckCloudConnexaConnectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cloudconnexa.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudconnexa_connector" {
			continue
		}

		connectorId := rs.Primary.ID
		connector, err := client.Connectors.GetByID(connectorId)

		if err != nil {
			return err
		}

		if connector != nil {
			return fmt.Errorf("connector with ID '%s' still exists", connectorId)
		}
	}

	return nil
}

func testAccCloudConnexaConnectorConfigBasic(rName string) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
  base_url = "https://%[1]s.api.openvpn.com"
}

resource "cloudconnexa_connector" "test" {
  name              = "%s"
  vpn_region_id     = "us-west-1"
  network_item_type = "HOST"
  network_item_id   = "example_network_item_id"
}
`, testCloudID, rName)
}
