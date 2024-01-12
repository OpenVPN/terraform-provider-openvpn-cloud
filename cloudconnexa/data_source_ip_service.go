package cloudconnexa

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

func dataSourceIPService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIPServiceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceServiceConfig(),
			},
			"network_item_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_item_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIPServiceRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	service, err := c.IPServices.Get(
		data.Id(),
	)

	if err != nil {
		return diag.FromErr(err)
	}
	setResourceData(data, service)
	return nil
}
