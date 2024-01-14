package cloudconnexa

import (
	"context"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
	"hash/fnv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceHost() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_host` to create an Cloud Connexa host.",
		CreateContext: resourceHostCreate,
		ReadContext:   resourceHostRead,
		UpdateContext: resourceHostUpdate,
		DeleteContext: resourceHostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the host.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "The description for the UI. Defaults to `Managed by Terraform`.",
			},
			"internet_access": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "LOCAL",
				ValidateFunc: validation.StringInSlice([]string{"BLOCKED", "GLOBAL_INTERNET", "LOCAL"}, false),
				Description:  "The type of internet access provided. Valid values are `BLOCKED`, `GLOBAL_INTERNET`, or `LOCAL`. Defaults to `LOCAL`.",
			},
			"system_subnets": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The IPV4 and IPV6 subnets automatically assigned to this host.",
			},
			"connector": {
				Type:     schema.TypeSet,
				Required: true,
				Set: func(i interface{}) int {
					n := i.(map[string]interface{})["name"]
					h := fnv.New32a()
					h.Write([]byte(n.(string)))
					return int(h.Sum32())
				},
				Description: "The set of connectors to be associated with this host. Can be defined more than once.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the connector associated with this host.",
						},
						"vpn_region_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the region where the connector will be deployed.",
						},
						"network_item_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network object type. This typically will be set to `HOST`.",
						},
						"network_item_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The host id.",
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
						"profile": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var connectors []cloudconnexa.Connector
	configConnectors := d.Get("connector").(*schema.Set)
	for _, c := range configConnectors.List() {
		connectors = append(connectors, cloudconnexa.Connector{
			Name:        c.(map[string]interface{})["name"].(string),
			VpnRegionId: c.(map[string]interface{})["vpn_region_id"].(string),
		})
	}
	h := cloudconnexa.Host{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		InternetAccess: d.Get("internet_access").(string),
		Connectors:     connectors,
	}
	host, err := c.Hosts.Create(h)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(host.Id)
	diagnostics := setConnectorsList(d, c, host.Connectors)
	if diagnostics != nil {
		return diagnostics
	}

	return append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "The connector for this host needs to be set up manually",
		Detail:   "Terraform only creates the Cloud Connexa connector object for this host, but additional manual steps are required to associate a host in your infrastructure with this connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.",
	})
}

func resourceHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	host, err := c.Hosts.Get(d.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if host == nil {
		d.SetId("")
		return diags
	}
	d.Set("name", host.Name)
	d.Set("description", host.Description)
	d.Set("internet_access", host.InternetAccess)
	d.Set("system_subnets", host.SystemSubnets)

	diagnostics := setConnectorsList(d, c, host.Connectors)
	if diagnostics != nil {
		return diagnostics
	}
	return diags
}

func resourceHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	if d.HasChange("connector") {
		old, new := d.GetChange("connector")
		oldSet := old.(*schema.Set)
		newSet := new.(*schema.Set)
		if oldSet.Len() == 0 && newSet.Len() > 0 {
			// This happens when importing the resource
			newConnector := cloudconnexa.Connector{
				Name:            newSet.List()[0].(map[string]interface{})["name"].(string),
				VpnRegionId:     newSet.List()[0].(map[string]interface{})["vpn_region_id"].(string),
				NetworkItemType: "HOST",
			}
			_, err := c.Connectors.Create(newConnector, d.Id())
			if err != nil {
				return append(diags, diag.FromErr(err)...)
			}
		} else {
			for _, o := range oldSet.List() {
				if !newSet.Contains(o) {
					err := c.Connectors.Delete(o.(map[string]interface{})["id"].(string), d.Id(), "HOST")
					if err != nil {
						diags = append(diags, diag.FromErr(err)...)
					}
				}
			}
			for _, n := range newSet.List() {
				if !oldSet.Contains(n) {
					newConnector := cloudconnexa.Connector{
						Name:            n.(map[string]interface{})["name"].(string),
						VpnRegionId:     n.(map[string]interface{})["vpn_region_id"].(string),
						NetworkItemType: "HOST",
					}
					_, err := c.Connectors.Create(newConnector, d.Id())
					if err != nil {
						diags = append(diags, diag.FromErr(err)...)
					}
				}
			}
		}
	}
	if d.HasChanges("name", "description", "internet_access") {
		_, newName := d.GetChange("name")
		_, newDescription := d.GetChange("description")
		_, newAccess := d.GetChange("internet_access")
		err := c.Hosts.Update(cloudconnexa.Host{
			Id:             d.Id(),
			Name:           newName.(string),
			Description:    newDescription.(string),
			InternetAccess: newAccess.(string),
		})
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	return append(diags, resourceHostRead(ctx, d, m)...)
}

func resourceHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	hostId := d.Id()
	err := c.Hosts.Delete(hostId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

func setConnectorsList(data *schema.ResourceData, c *cloudconnexa.Client, connectors []cloudconnexa.Connector) diag.Diagnostics {
	connectorsList := make([]interface{}, len(connectors))
	for i, connector := range connectors {
		connectorsData, err := getConnectorsListItem(c, connector)
		if err != nil {
			return diag.FromErr(err)
		}
		connectorsList[i] = connectorsData
	}
	err := data.Set("connector", connectorsList)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getConnectorsListItem(c *cloudconnexa.Client, connector cloudconnexa.Connector) (map[string]interface{}, error) {
	connectorsData := map[string]interface{}{
		"id":                connector.Id,
		"name":              connector.Name,
		"vpn_region_id":     connector.VpnRegionId,
		"ip_v4_address":     connector.IPv4Address,
		"ip_v6_address":     connector.IPv6Address,
		"network_item_id":   connector.NetworkItemId,
		"network_item_type": connector.NetworkItemType,
	}

	connectorProfile, err := c.Connectors.GetProfile(connector.Id)
	if err != nil {
		return nil, err
	}
	connectorsData["profile"] = connectorProfile
	return connectorsData, nil
}
