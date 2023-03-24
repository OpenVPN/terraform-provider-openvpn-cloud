package openvpncloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/OpenVPN/terraform-provider-openvpn-cloud/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func TestAccOpenvpncloudRoute_basic(t *testing.T) {
	rn := "openvpncloud_route.test"
	ip, err := acctest.RandIpAddress("10.0.0.0/8")
	require.NoError(t, err)
	route := client.Route{
		Description: "test" + acctest.RandString(10),
		Type:        client.RouteTypeIPV4,
		Value:       ip + "/32",
	}
	routeChanged := route
	routeChanged.Description = acctest.RandStringFromCharSet(10, alphabet)
	networkRandString := "test" + acctest.RandString(10)
	var routeId string

	check := func(r client.Route) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckOpenvpncloudRouteExists(rn, &routeId),
			resource.TestCheckResourceAttr(rn, "description", r.Description),
			resource.TestCheckResourceAttr(rn, "type", r.Type),
			resource.TestCheckResourceAttr(rn, "value", r.Value),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckOpenvpncloudRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenvpncloudRouteConfig(route, networkRandString),
				Check:  check(route),
			},
			{
				Config: testAccOpenvpncloudRouteConfig(routeChanged, networkRandString),
				Check:  check(routeChanged),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: testAccOpenvpncloudRouteImportStateIdFunc(rn),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckOpenvpncloudRouteDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openvpncloud_route" {
			continue
		}
		routeId := rs.Primary.ID
		r, err := client.GetRouteById(routeId)
		if err == nil {
			return err
		}
		if r != nil {
			return errors.New("route still exists")
		}
	}
	return nil
}

func testAccCheckOpenvpncloudRouteExists(n string, routeID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		client := testAccProvider.Meta().(*client.Client)
		_, err := client.GetRouteById(rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccOpenvpncloudRouteImportStateIdFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}
		return rs.Primary.ID, nil
	}
}

func testAccOpenvpncloudRouteConfig(r client.Route, networkRandStr string) string {
	return fmt.Sprintf(`
provider "openvpncloud" {
	base_url = "https://%[1]s.api.openvpn.com"
}
resource "openvpncloud_network" "test" {
	name = "%[5]s"
	default_connector {
	  name          = "%[5]s"
	  vpn_region_id = "fi-hel"
	}
	default_route {
	  value = "10.1.2.0/24"
	  type  = "IP_V4"
	}
}
resource "openvpncloud_route" "test" {
	network_item_id = openvpncloud_network.test.id
	description     = "%[2]s"
	value           = "%[3]s"
	type            = "%[4]s"
}
`, testCloudID, r.Description, r.Value, r.Type, networkRandStr)
}
