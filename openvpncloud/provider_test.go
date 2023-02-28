package openvpncloud

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const alphabet = "abcdefghigklmnopqrstuvwxyz"

var testCloudID = os.Getenv("OPENVPNCLOUD_TEST_ORGANIZATION")
var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"openvpncloud": testAccProvider,
	}
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"openvpncloud": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OPENVPN_CLOUD_CLIENT_ID"); v == "" {
		t.Fatal("OPENVPN_CLOUD_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("OPENVPN_CLOUD_CLIENT_SECRET"); v == "" {
		t.Fatal("OPENVPN_CLOUD_CLIENT_SECRET must be set for acceptance tests")
	}
}
