package openvpncloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `openvpncloud_network` to create an OpenVPN Cloud Network.",
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the network.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "The display description for this resource. Defaults to `Managed by Terraform`.",
			},
			"egress": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Boolean to control whether this network provides an egress or not.",
			},
			"internet_access": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      client.InternetAccessLocal,
				ValidateFunc: validation.StringInSlice([]string{client.InternetAccessBlocked, client.InternetAccessGlobalInternet, client.InternetAccessLocal}, false),
				Description:  "The type of internet access provided. Valid values are `BLOCKED`, `GLOBAL_INTERNET`, or `LOCAL`. Defaults to `LOCAL`.",
			},
			"system_subnets": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The IPV4 and IPV6 subnets automatically assigned to this network.",
			},
			"default_route": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The default route of this network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      client.RouteTypeIPV4,
							ValidateFunc: validation.StringInSlice([]string{client.RouteTypeIPV4, client.RouteTypeIPV6, client.RouteTypeDomain}, false),
							Description:  "The type of route. Valid values are `IP_V4`, `IP_V6`, and `DOMAIN`.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The target value of the default route.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of this resource.",
						},
					},
				},
			},
			"default_connector": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The default connector of this network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of this resource.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the connector automatically created and attached to this network.",
						},
						"vpn_region_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the region where the default connector will be deployed.",
						},
						"network_item_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network object type. This typically will be set to `NETWORK`.",
						},
						"network_item_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parent network id.",
						},
						"ip_v4_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPV4 address of the default connector.",
						},
						"ip_v6_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPV6 address of the default connector.",
						},
					},
				},
			},
		},
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	configConnector := d.Get("default_connector").([]interface{})[0].(map[string]interface{})
	connectors := []client.Connector{
		{
			Name:        configConnector["name"].(string),
			VpnRegionId: configConnector["vpn_region_id"].(string),
		},
	}
	n := client.Network{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Egress:         d.Get("egress").(bool),
		InternetAccess: d.Get("internet_access").(string),
		Connectors:     connectors,
	}
	network, err := c.CreateNetwork(n)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(network.Id)
	configRoute := d.Get("default_route").([]interface{})[0].(map[string]interface{})
	defaultRoute, err := c.CreateRoute(network.Id, client.Route{
		Type:  configRoute["type"].(string),
		Value: configRoute["value"].(string),
	})
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	defaultRouteWithIdSlice := make([]map[string]interface{}, 1)
	defaultRouteWithIdSlice[0] = map[string]interface{}{
		"id": defaultRoute.Id,
	}
	d.Set("default_route", defaultRouteWithIdSlice)
	return append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "The default connector for this network needs to be set up manually",
		Detail:   "Terraform only creates the OpenVPN Cloud default connector object for this network, but additional manual steps are required to associate a host in your infrastructure with this connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.",
	})
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	network, err := c.GetNetworkById(d.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if network == nil {
		d.SetId("")
		return diags
	}
	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("egress", network.Egress)
	d.Set("internet_access", network.InternetAccess)
	d.Set("system_subnets", network.SystemSubnets)
	if len(d.Get("default_connector").([]interface{})) > 0 {
		configConnector := d.Get("default_connector").([]interface{})[0].(map[string]interface{})
		connectorName := configConnector["name"].(string)
		networkConnectors, err := c.GetConnectorsForNetwork(network.Id)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		err = d.Set("default_connector", getConnectorSlice(networkConnectors, network.Id, connectorName))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if len(d.Get("default_route").([]interface{})) > 0 {
		configRoute := d.Get("default_route").([]interface{})[0].(map[string]interface{})
		route, err := c.GetNetworkRoute(d.Id(), configRoute["id"].(string))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		if route == nil {
			d.Set("default_route", []map[string]interface{}{})
		} else {
			defaultRoute := []map[string]interface{}{
				{
					"id":   configRoute["id"].(string),
					"type": route.Type,
				},
			}
			if route.Type == client.RouteTypeIPV4 || route.Type == client.RouteTypeIPV6 {
				defaultRoute[0]["value"] = route.Subnet
			} else if route.Type == client.RouteTypeDomain {
				defaultRoute[0]["value"] = route.Domain
			}
			d.Set("default_route", defaultRoute)
		}
	}
	return diags
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	if d.HasChange("default_connector") {
		old, new := d.GetChange("default_connector")
		oldSlice := old.([]interface{})
		newSlice := new.([]interface{})
		if len(oldSlice) == 0 && len(newSlice) == 1 {
			// This happens when importing the resource
			newConnector := client.Connector{
				Name:            newSlice[0].(map[string]interface{})["name"].(string),
				VpnRegionId:     newSlice[0].(map[string]interface{})["vpn_region_id"].(string),
				NetworkItemType: client.NetworkItemTypeNetwork,
			}
			_, err := c.AddConnector(newConnector, d.Id())
			if err != nil {
				return append(diags, diag.FromErr(err)...)
			}
		} else {
			oldMap := oldSlice[0].(map[string]interface{})
			newMap := newSlice[0].(map[string]interface{})
			if oldMap["name"].(string) != newMap["name"].(string) || oldMap["vpn_region_id"].(string) != newMap["vpn_region_id"].(string) {
				newConnector := client.Connector{
					Name:            newMap["name"].(string),
					VpnRegionId:     newMap["vpn_region_id"].(string),
					NetworkItemType: client.NetworkItemTypeNetwork,
				}
				_, err := c.AddConnector(newConnector, d.Id())
				if err != nil {
					return append(diags, diag.FromErr(err)...)
				}
				if len(oldMap["id"].(string)) > 0 {
					// This can sometimes happen when importing the resource
					err = c.DeleteConnector(oldMap["id"].(string), d.Id(), oldMap["network_item_type"].(string))
					if err != nil {
						return append(diags, diag.FromErr(err)...)
					}
				}
			}
		}
	}
	if d.HasChange("default_route") {
		old, new := d.GetChange("default_route")
		oldSlice := old.([]interface{})
		newSlice := new.([]interface{})
		if len(oldSlice) == 0 && len(newSlice) == 1 {
			// This happens when importing the resource
			newMap := newSlice[0].(map[string]interface{})
			routeType := newMap["type"]
			routeValue := newMap["value"]
			route := client.Route{
				Type:  routeType.(string),
				Value: routeValue.(string),
			}
			defaultRoute, err := c.CreateRoute(d.Id(), route)
			if err != nil {
				return append(diags, diag.FromErr(err)...)
			}
			defaultRouteWithIdSlice := make([]map[string]interface{}, 1)
			defaultRouteWithIdSlice[0] = map[string]interface{}{
				"id": defaultRoute.Id,
			}
			err = d.Set("default_route", defaultRouteWithIdSlice)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		} else {
			newMap := newSlice[0].(map[string]interface{})
			routeId := newMap["id"]
			routeType := newMap["type"]
			routeValue := newMap["value"]
			route := client.Route{
				Id:    routeId.(string),
				Type:  routeType.(string),
				Value: routeValue.(string),
			}
			err := c.UpdateRoute(d.Id(), route)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	}
	if d.HasChanges("name", "description", "internet_access", "egress") {
		_, newName := d.GetChange("name")
		_, newDescription := d.GetChange("description")
		_, newEgress := d.GetChange("egress")
		_, newAccess := d.GetChange("internet_access")
		err := c.UpdateNetwork(client.Network{
			Id:             d.Id(),
			Name:           newName.(string),
			Description:    newDescription.(string),
			Egress:         newEgress.(bool),
			InternetAccess: newAccess.(string),
		})
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	return append(diags, resourceNetworkRead(ctx, d, m)...)
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	networkId := d.Id()
	err := c.DeleteNetwork(networkId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

func getNetworkConnectorSlice(networkConnectors []client.Connector, networkId string, connectorName string) []interface{} {
	if len(networkConnectors) == 0 {
		return nil
	}
	connectorsList := make([]interface{}, 1)
	for _, c := range networkConnectors {
		if c.NetworkItemId == networkId && c.Name == connectorName {
			connector := make(map[string]interface{})
			connector["id"] = c.Id
			connector["name"] = c.Name
			connector["network_item_id"] = c.NetworkItemId
			connector["network_item_type"] = c.NetworkItemType
			connector["vpn_region_id"] = c.VpnRegionId
			connector["ip_v4_address"] = c.IPv4Address
			connector["ip_v6_address"] = c.IPv6Address
			connectorsList[0] = connector
			break
		}
	}
	return connectorsList
}
