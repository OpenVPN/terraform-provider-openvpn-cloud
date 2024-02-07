package cloudconnexa

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudConnexaService_basic(t *testing.T) {
	rn := "cloudconnexa_service.test"
	networkName := acctest.RandStringFromCharSet(10, alphabet)
	service := cloudconnexa.IPService{
		Name: acctest.RandStringFromCharSet(10, alphabet),
	}
	serviceChanged := service
	serviceChanged.Name = fmt.Sprintf("changed-%s", acctest.RandStringFromCharSet(10, alphabet))

	check := func(service cloudconnexa.IPService) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckCloudConnexaServiceExists(rn, networkName),
			resource.TestCheckResourceAttr(rn, "name", service.Name),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaServiceConfig(service, networkName),
				Check:  check(service),
			},
			{
				Config: testAccCloudConnexaServiceConfig(serviceChanged, networkName),
				Check:  check(serviceChanged),
			},
		},
	})
}

func testAccCheckCloudConnexaServiceExists(rn, networkId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		c := testAccProvider.Meta().(*cloudconnexa.Client)
		_, err := c.IPServices.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckCloudConnexaServiceDestroy(state *terraform.State) error {
	c := testAccProvider.Meta().(*cloudconnexa.Client)
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "cloudconnexa_service" {
			continue
		}
		id := rs.Primary.Attributes["id"]
		s, err := c.IPServices.Get(id)
		if err == nil || s != nil {
			return fmt.Errorf("service still exists")
		}
	}
	return nil
}
func testAccCloudConnexaServiceConfig(service cloudconnexa.IPService, networkName string) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
	base_url = "https://%s.api.openvpn.com"
}

resource "cloudconnexa_network" "test" {
	name = "%s"
	description = "test"

	default_connector {
	  name          = "%s"
	  vpn_region_id = "fi-hel"
	}
	default_route {
	  value = "10.1.2.0/24"
	  type  = "IP_V4"
	}
}

resource "cloudconnexa_ip_service" "test" {
	name = "%s"
	type = "SERVICE_DESTINATION"
	description = "test"
	network_item_type = "NETWORK"
	network_item_id = cloudconnexa_network.test.id
	routes = ["test.ua" ]
	config {
		service_types = ["ANY"]
	}
}
`, testCloudID, networkName, fmt.Sprintf("connector_%s", networkName), service.Name)
}
