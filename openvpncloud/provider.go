package openvpncloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

type Token struct {
	AccessToken string `json:"access_token"`
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OPENVPN_CLOUD_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OPENVPN_CLOUD_CLIENT_SECRET", nil),
			},
			"base_url": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"openvpncloud_network":    resourceNetwork(),
			"openvpncloud_connector":  resourceConnector(),
			"openvpncloud_route":      resourceRoute(),
			"openvpncloud_dns_record": resourceDnsRecord(),
			"openvpncloud_user":       resourceUser(),
			"openvpncloud_host":       resourceHost(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"openvpncloud_network":        dataSourceNetwork(),
			"openvpncloud_connector":      dataSourceConnector(),
			"openvpncloud_user":           dataSourceUser(),
			"openvpncloud_user_group":     dataSourceUserGroup(),
			"openvpncloud_vpn_region":     dataSourceVpnRegion(),
			"openvpncloud_network_routes": dataSourceNetworkRoutes(),
			"openvpncloud_host":           dataSourceHost(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	baseUrl := d.Get("base_url").(string)
	openVPNClient, err := client.NewClient(baseUrl, clientId, clientSecret)
	var diags diag.Diagnostics
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create client",
			Detail:   fmt.Sprintf("Error: %v", err),
		})
		return nil, diags
	}
	return openVPNClient, nil
}
