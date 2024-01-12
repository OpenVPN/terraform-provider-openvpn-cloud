package cloudconnexa

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

func dataSourceConnector() *schema.Resource {
	return &schema.Resource{
		Description: "Use an `cloudconnexa_connector` data source to read an existing Cloud Connexa connector.",
		ReadContext: dataSourceConnectorRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the connector.",
			},
			"network_item_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the network or host with which the connector is associated.",
			},
			"network_item_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network object type of the connector. This typically will be set to either `NETWORK` or `HOST`.",
			},
			"vpn_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the region where the connector is deployed.",
			},
			"ip_v4_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPV4 address of the connector.",
			},
			"ip_v6_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPV6 address of the connector.",
			},
		},
	}
}

func dataSourceConnectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	connector, err := c.Connectors.GetByName(d.Get("name").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("name", connector.Name)
	d.Set("network_item_id", connector.NetworkItemId)
	d.Set("network_item_type", connector.NetworkItemType)
	d.Set("vpn_region_id", connector.VpnRegionId)
	d.Set("ip_v4_address", connector.IPv4Address)
	d.Set("ip_v6_address", connector.IPv6Address)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
