package cloudconnexa

import (
	"context"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const alphabet = "abcdefghigklmnopqrstuvwxyz"

var testCloudID = os.Getenv("CLOUDCONNEXA_TEST_ORGANIZATION")
var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"cloudconnexa": testAccProvider,
	}
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"cloudconnexa": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	err := Provider().InternalValidate()
	require.NoError(t, err)

	// must have the required error when the credentials are not set
	t.Setenv(clientIDEnvVar, "")
	t.Setenv(clientSecretEnvVar, "")
	rc := terraform.ResourceConfig{}
	diags := Provider().Configure(context.Background(), &rc)
	assert.True(t, diags.HasError())

	for _, d := range diags {
		assert.Truef(t, strings.Contains(d.Detail, cloudconnexa.ErrCredentialsRequired.Error()),
			"error message does not contain the expected error: %s", d.Detail)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(clientIDEnvVar); v == "" {
		t.Fatalf("%s must be set for acceptance tests", clientIDEnvVar)
	}
	if v := os.Getenv(clientSecretEnvVar); v == "" {
		t.Fatalf("%s must be set for acceptance tests", clientSecretEnvVar)
	}
}
