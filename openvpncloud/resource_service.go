package openvpncloud

import (
	"context"
	"github.com/OpenVPN/terraform-provider-openvpn-cloud/client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	validValues = []string{"ANY", "BGP", "CUSTOM", "DHCP", "DNS", "FTP", "HTTP", "HTTPS", "IMAP", "IMAPS", "NTP", "POP3", "POP3S", "SMTP", "SMTPS", "SNMP", "SSH", "TELNET", "TFTP"}
)

func resourceService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		DeleteContext: resourceServiceDelete,
		UpdateContext: resourceServiceUpdate,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Default:      "Created by Terraform OpenVPN Cloud Provider",
				ValidateFunc: validation.StringLenBetween(1, 255),
				Optional:     true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"IP_SOURCE", "SERVICE_DESTINATION"}, false),
			},
			"routes": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     resourceServiceConfig(),
			},
			"network_item_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"NETWORK", "HOST"}, false),
			},
			"network_item_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServiceUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*client.Client)
	networkItemId := data.Get("network_item_id").(string)
	networkItemType := data.Get("network_item_type").(string)

	s, err := c.UpdateService(data.Id(), networkItemType, networkItemId, resourceDataToService(data))
	if err != nil {
		return diag.FromErr(err)
	}

	setResourceData(data, s)
	return nil
}

func resourceServiceConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"custom_service_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"icmp_type": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lower_value": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"upper_value": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"service_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {

						val := i.(string)
						for _, validValue := range validValues {
							if val == validValue {
								return nil
							}
						}
						return diag.Errorf("service type must be one of %s", validValues)
					},
				},
			},
		},
	}
}

func resourceServiceRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*client.Client)
	service, err := c.GetService(
		data.Id(),
		data.Get("network_item_type").(string),
		data.Get("network_item_id").(string),
	)

	if err != nil {
		return diag.FromErr(err)
	}
	setResourceData(data, service)
	return nil
}

func setResourceData(data *schema.ResourceData, service *client.Service) {
	data.SetId(service.Id)
	_ = data.Set("name", service.Name)
	_ = data.Set("description", service.Description)
	_ = data.Set("type", service.Type)
	_ = data.Set("routes", flattenRoutes(service.Routes))
	_ = data.Set("config", flattenServiceConfig(service.Config))
	_ = data.Set("network_item_type", service.NetworkItemType)
	_ = data.Set("network_item_id", service.NetworkItemId)
}

func resourceServiceDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*client.Client)
	err := c.DeleteService(
		data.Id(),
		data.Get("network_item_type").(string),
		data.Get("network_item_id").(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenServiceConfig(config *client.ServiceConfig) interface{} {
	var data = map[string]interface{}{
		"custom_service_types": flattenCustomServiceTypes(config.CustomServiceTypes),
		"service_types":        config.ServiceTypes,
	}
	return []interface{}{data}
}

func flattenCustomServiceTypes(types []*client.CustomServiceType) interface{} {
	var data []interface{}
	for _, t := range types {
		data = append(
			data,
			map[string]interface{}{
				"icmp_type": flattenIcmpType(t.IcmpType),
			},
		)
	}
	return data
}

func flattenIcmpType(icmpType []client.Range) interface{} {
	var data []interface{}
	for _, t := range icmpType {
		data = append(
			data,
			map[string]interface{}{
				"lower_value": t.LowerValue,
				"upper_value": t.UpperValue,
			},
		)
	}
	return data
}

func flattenRoutes(routes []*client.Route) []string {
	var data []string
	for _, route := range routes {
		data = append(
			data,
			route.Domain,
		)
	}
	return data
}

func resourceServiceCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*client.Client)
	var diags diag.Diagnostics

	service, err := c.CreateService(resourceDataToService(data))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	data.SetId(service.Id)
	setResourceData(data, service)
	return diags
}

func resourceDataToService(data *schema.ResourceData) *client.Service {
	routes := data.Get("routes").([]interface{})
	var configRoutes []*client.Route
	for _, r := range routes {
		configRoutes = append(
			configRoutes,
			&client.Route{
				Value: r.(string),
			},
		)
	}

	config := client.ServiceConfig{}
	configList := data.Get("config").([]interface{})
	if len(configList) > 0 && configList[0] != nil {

		config.CustomServiceTypes = []*client.CustomServiceType{}
		config.ServiceTypes = []string{}

		mainConfig := configList[0].(map[string]interface{})
		for _, r := range mainConfig["custom_service_types"].([]interface{}) {
			cst := r.(map[string]interface{})
			var icmpTypes []client.Range
			for _, r := range cst["icmp_type"].([]interface{}) {
				icmpType := r.(map[string]interface{})
				icmpTypes = append(
					icmpTypes,
					client.Range{
						LowerValue: icmpType["lower_value"].(int),
						UpperValue: icmpType["upper_value"].(int),
					},
				)
			}
			config.CustomServiceTypes = append(
				config.CustomServiceTypes,
				&client.CustomServiceType{
					IcmpType: icmpTypes,
				},
			)
		}

		for _, r := range mainConfig["service_types"].([]interface{}) {
			config.ServiceTypes = append(config.ServiceTypes, r.(string))
		}
	}

	s := &client.Service{
		Name:            data.Get("name").(string),
		Description:     data.Get("description").(string),
		NetworkItemId:   data.Get("network_item_id").(string),
		NetworkItemType: data.Get("network_item_type").(string),
		Type:            data.Get("type").(string),
		Routes:          configRoutes,
		Config:          &config,
	}
	return s
}
