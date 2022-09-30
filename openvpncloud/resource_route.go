package openvpncloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func resourceRoute() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `openvpncloud_route` to create a route on an OpenVPN Cloud network.",
		CreateContext: resourceRouteCreate,
		ReadContext:   resourceRouteRead,
		DeleteContext: resourceRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{client.RouteTypeIPV4, client.RouteTypeIPV6, client.RouteTypeDomain}, false),
				Description:  "The type of route. Valid values are `IP_V4`, `IP_V6`, and `DOMAIN`.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target value of the default route.",
			},
			"network_item_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the network on which to create the route.",
			},
		},
	}
}

func resourceRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	networkItemId := d.Get("network_item_id").(string)
	routeType := d.Get("type").(string)
	routeValue := d.Get("value").(string)
	r := client.Route{
		Type:  routeType,
		Value: routeValue,
	}
	route, err := c.CreateRoute(networkItemId, r)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(route.Id)
	if routeType == client.RouteTypeIPV4 || routeType == client.RouteTypeIPV6 {
		d.Set("value", route.Subnet)
	} else if routeType == client.RouteTypeDomain {
		d.Set("value", route.Domain)
	}
	return diags
}

func resourceRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	routeId := d.Id()
	r, err := c.GetRouteById(routeId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if r == nil {
		d.SetId("")
	} else {
		d.Set("type", r.Type)
		if r.Type == client.RouteTypeIPV4 || r.Type == client.RouteTypeIPV6 {
			d.Set("value", r.Subnet)
		} else if r.Type == client.RouteTypeDomain {
			d.Set("resourceRouteRead", r.Domain)
		}
		d.Set("network_item_id", r.NetworkItemId)
	}
	return diags
}

func resourceRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	routeId := d.Id()
	networkItemId := d.Get("network_item_id").(string)
	err := c.DeleteRoute(networkItemId, routeId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
