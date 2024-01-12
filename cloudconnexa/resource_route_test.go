package cloudconnexa

import (
	"errors"
	"fmt"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func TestAccCloudConnexaRoute_basic(t *testing.T) {
	rn := "cloudconnexa_route.test"
	ip, err := acctest.RandIpAddress("10.0.0.0/8")
	require.NoError(t, err)
	route := cloudconnexa.Route{
		Description: "test" + acctest.RandString(10),
		Type:        "IP_V4",
		Subnet:      ip + "/32",
	}
	routeChanged := route
	routeChanged.Description = acctest.RandStringFromCharSet(10, alphabet)
	networkRandString := "test" + acctest.RandString(10)
	var routeId string

	check := func(r cloudconnexa.Route) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckCloudConnexaRouteExists(rn, &routeId),
			resource.TestCheckResourceAttr(rn, "description", r.Description),
			resource.TestCheckResourceAttr(rn, "type", r.Type),
			resource.TestCheckResourceAttr(rn, "value", r.Subnet),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaRouteConfig(route, networkRandString),
				Check:  check(route),
			},
			{
				Config: testAccCloudConnexaRouteConfig(routeChanged, networkRandString),
				Check:  check(routeChanged),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: testAccCloudConnexaRouteImportStateIdFunc(rn),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudConnexaRouteDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cloudconnexa.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudconnexa_route" {
			continue
		}
		routeId := rs.Primary.ID
		r, err := client.Routes.Get(routeId)
		if err == nil {
			return err
		}
		if r != nil {
			return errors.New("route still exists")
		}
	}
	return nil
}

func testAccCheckCloudConnexaRouteExists(n string, routeID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		client := testAccProvider.Meta().(*cloudconnexa.Client)
		_, err := client.Routes.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCloudConnexaRouteImportStateIdFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}
		return rs.Primary.ID, nil
	}
}

func testAccCloudConnexaRouteConfig(r cloudconnexa.Route, networkRandStr string) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
	base_url = "https://%[1]s.api.openvpn.com"
}
resource "cloudconnexa_network" "test" {
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
resource "cloudconnexa_route" "test" {
	network_item_id = cloudconnexa_network.test.id
	description     = "%[2]s"
	value           = "%[3]s"
	type            = "%[4]s"
}
`, testCloudID, r.Description, r.Subnet, r.Type, networkRandStr)
}
