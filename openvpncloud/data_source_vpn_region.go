package openvpncloud

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func dataSourceVpnRegion() *schema.Resource {
	return &schema.Resource{
		Description: "Use a `openvpncloud_vpn_region` data source to read an OpenVPN Cloud VPN region.",
		ReadContext: dataSourceVpnRegionRead,
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the region.",
			},
			"continent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The continent of the region.",
			},
			"country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The country of the region.",
			},
			"country_iso": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ISO code of the country of the region.",
			},
			"region_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the region.",
			},
		},
	}
}

func dataSourceVpnRegionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	vpnRegionId := d.Get("region_id").(string)
	vpnRegion, err := c.GetVpnRegion(vpnRegionId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if vpnRegion == nil {
		return append(diags, diag.Errorf("VPN region with id %s was not found", vpnRegionId)...)
	}
	d.Set("region_id", vpnRegion.Id)
	d.Set("continent", vpnRegion.Continent)
	d.Set("country", vpnRegion.Country)
	d.Set("country_iso", vpnRegion.CountryISO)
	d.Set("region_name", vpnRegion.RegionName)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
