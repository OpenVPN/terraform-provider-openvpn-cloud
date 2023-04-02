package openvpncloud

import (
	"fmt"
	"testing"

	"github.com/OpenVPN/terraform-provider-openvpn-cloud/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenvpncloudDnsRecord_basic(t *testing.T) {
	resourceName := "openvpncloud_dns_record.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckOpenvpncloudDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenvpncloudDnsRecordConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "domain", "test.example.com"),
					resource.TestCheckResourceAttr(resourceName, "ip_v4_addresses.0", "192.168.1.1"),
					resource.TestCheckResourceAttr(resourceName, "ip_v4_addresses.1", "192.168.1.2"),
					resource.TestCheckResourceAttr(resourceName, "ip_v6_addresses.0", "2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
					resource.TestCheckResourceAttr(resourceName, "ip_v6_addresses.1", "2001:0db8:85a3:0000:0000:8a2e:0370:7335"),
				),
			},
		},
	})
}

func testAccCheckOpenvpncloudDnsRecordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openvpncloud_dns_record" {
			continue
		}

		recordId := rs.Primary.ID
		r, err := client.GetDnsRecord(recordId)

		if err != nil {
			return err
		}

		if r != nil {
			return fmt.Errorf("DNS record with ID '%s' still exists", recordId)
		}
	}

	return nil
}

const testAccOpenvpncloudDnsRecordConfigBasic = `
provider "openvpncloud" {
  base_url = "https://%[1]s.api.openvpn.com"
}

resource "openvpncloud_dns_record" "test" {
  domain          = "test.example.com"
  ip_v4_addresses = ["192.168.1.1", "192.168.1.2"]
  ip_v6_addresses = ["2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:0db8:85a3:0000:0000:8a2e:0370:7335"]
}
`
