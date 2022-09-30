package openvpncloud

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func dataSourceNetworkRoutes() *schema.Resource {
	return &schema.Resource{
		Description: "Use an `openvpncloud_network_routes` data source to read all the routes associated with an OpenVPN Cloud network.",
		ReadContext: dataSourceNetworkRoutesRead,
		Schema: map[string]*schema.Schema{
			"network_item_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the OpenVPN Cloud network of the routes to be discovered.",
			},
			"routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of routes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of route. Valid values are `IP_V4`, `IP_V6`, and `DOMAIN`.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of the route, either an IPV4 address, an IPV6 address, or a DNS hostname.",
						},
					},
				},
			},
		},
	}
}

func dataSourceNetworkRoutesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	routes, err := c.GetRoutes(d.Get("network_item_id").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	configRoutes := make([]map[string]interface{}, len(routes))
	for i, r := range routes {
		route := make(map[string]interface{})
		routeType := r.Type
		route["type"] = routeType
		if routeType == client.RouteTypeIPV4 || routeType == client.RouteTypeIPV6 {
			route["value"] = r.Subnet
		} else if routeType == client.RouteTypeDomain {
			route["value"] = r.Domain
		}
		configRoutes[i] = route
	}
	d.Set("routes", configRoutes)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
