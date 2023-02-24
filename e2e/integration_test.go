package e2e

import (
	"fmt"
	api "github.com/OpenVPN/terraform-provider-openvpn-cloud/client"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

const (
	OpenvpnHostKey              = "OVPN_HOST"
	OpenvpnCloudClientIdKey     = "OPENVPN_CLOUD_CLIENT_ID"
	OpenvpnCloudClientSecretKey = "OPENVPN_CLOUD_CLIENT_SECRET"
)

func TestCreationDeletion(t *testing.T) {
	validateEnvVars(t)

	terraformOptions := &terraform.Options{

		NoColor: os.Getenv("NO_COLOR") == "1",

		// The path to where our Terraform code is located
		TerraformDir: "../example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	t.Cleanup(func() {
		terraform.Destroy(t, terraformOptions)
	})

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	hostID := terraform.Output(t, terraformOptions, "host_id")
	connectorID := terraform.Output(t, terraformOptions, "connector_id")

	assert.NotEmpty(t, hostID)
	assert.NotEmpty(t, connectorID)

	client, err := api.NewClient(
		os.Getenv(OpenvpnHostKey),
		os.Getenv(OpenvpnCloudClientIdKey),
		os.Getenv(OpenvpnCloudClientSecretKey),
	)
	require.NoError(t, err)

	// Total waiting time: 1min
	totalAttempts := 10
	attemptWaitingTime := 6 * time.Second

	connectorWasOnline := false
	for i := 0; i < totalAttempts; i++ {
		t.Logf("Waiting for connector to be online (%d/%d)", i+1, totalAttempts)
		connector, err := client.GetConnectorById(connectorID)
		require.NoError(t, err, "Invalid connector ID in output")
		if connector.ConnectionStatus == api.ConnectionStatusOnline {
			connectorWasOnline = true
			break
		}
		time.Sleep(attemptWaitingTime)
	}
	assert.True(t, connectorWasOnline)
}

func validateEnvVars(t *testing.T) {
	validateEnvVar(t, OpenvpnHostKey)
	validateEnvVar(t, OpenvpnCloudClientIdKey)
	validateEnvVar(t, OpenvpnCloudClientSecretKey)
}

func validateEnvVar(t *testing.T, envVar string) {
	fmt.Println(os.Getenv(envVar))
	require.NotEmptyf(t, os.Getenv(envVar), "%s must be set for acceptance tests", envVar)
}
