package openvpncloud

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "Use a `openvpncloud_network` data source to read an OpenVPN Cloud network.",
		ReadContext: dataSourceNetworkRead,
		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network name.",
			},
			"egress": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Boolean to indicate whether this network provides an egress or not.",
			},
			"internet_access": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of internet access provided. Valid values are `BLOCKED`, `GLOBAL_INTERNET`, or `LOCAL`. Defaults to `LOCAL`.",
			},
			"system_subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The IPV4 and IPV6 subnets automatically assigned to this network.",
			},
			"routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The routes associated with this network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The route id.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of route. Valid values are `IP_V4`, `IP_V6`, and `DOMAIN`.",
						},
						"subnet": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of the route, either an IPV4 address, an IPV6 address, or a DNS hostname.",
						},
					},
				},
			},
			"connectors": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of connectors associated with this network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connector id.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connector name.",
						},
						"network_item_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the network with which the connector is associated.",
						},
						"network_item_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network object type of the connector. This typically will be set to `NETWORK`.",
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
				},
			},
		},
	}
}

func dataSourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	networkName := d.Get("name").(string)
	network, err := c.GetNetworkByName(networkName)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if network == nil {
		return append(diags, diag.Errorf("Network with name %s was not found", networkName)...)
	}
	d.Set("network_id", network.Id)
	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("egress", network.Egress)
	d.Set("internet_access", network.InternetAccess)
	d.Set("system_subnets", network.SystemSubnets)
	d.Set("routes", getRoutesSlice(&network.Routes))
	d.Set("connectors", getConnectorsSlice(&network.Connectors))
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

func getRoutesSlice(networkRoutes *[]client.Route) []interface{} {
	routes := make([]interface{}, len(*networkRoutes), len(*networkRoutes))
	for i, r := range *networkRoutes {
		route := make(map[string]interface{})
		route["id"] = r.Id
		route["subnet"] = r.Subnet
		route["type"] = r.Type
		routes[i] = route
	}
	return routes
}

func getConnectorsSlice(connectors *[]client.Connector) []interface{} {
	conns := make([]interface{}, len(*connectors), len(*connectors))
	for i, c := range *connectors {
		connector := make(map[string]interface{})
		connector["id"] = c.Id
		connector["name"] = c.Name
		connector["network_item_id"] = c.NetworkItemId
		connector["network_item_type"] = c.NetworkItemType
		connector["vpn_region_id"] = c.VpnRegionId
		connector["ip_v4_address"] = c.IPv4Address
		connector["ip_v6_address"] = c.IPv6Address
		conns[i] = connector
	}
	return conns
}
