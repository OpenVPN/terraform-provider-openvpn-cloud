package openvpncloud

import (
	"context"
	"github.com/OpenVPN/terraform-provider-openvpn-cloud/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `openvpncloud_user_group` to create an OpenVPN Cloud user group.",
		CreateContext: resourceUserGroupCreate,
		ReadContext:   resourceUserGroupRead,
		UpdateContext: resourceUserGroupUpdate,
		DeleteContext: resourceUserGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the user group.",
			},
			"connect_auth": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "AUTO",
				ValidateFunc: validation.StringInSlice([]string{"AUTH", "AUTO", "STRICT_AUTH"}, false),
			},
			"internet_access": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "LOCAL",
				ValidateFunc: validation.StringInSlice([]string{"LOCAL", "BLOCKED", "GLOBAL_INTERNET"}, false),
			},
			"max_device": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "The maximum number of devices that can be connected to the user group.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 40),
				Description:  "The name of the user group.",
			},
			"system_subnets": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Default:     nil,
				Description: "A list of subnets that are accessible to the user group.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vpn_region_ids": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "A list of VPN regions that are accessible to the user group.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceUserGroupUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*client.Client)
	var diags diag.Diagnostics
	ug := resourceDataToUserGroup(data)

	userGroup, err := c.UpdateUserGroup(data.Id(), ug)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if userGroup == nil {
		data.SetId("")
	} else {
		updateUserGroupData(data, userGroup)
	}

	return diags
}

func resourceDataToUserGroup(data *schema.ResourceData) *client.UserGroup {
	name := data.Get("name").(string)
	connectAuth := data.Get("connect_auth").(string)
	maxDevice := data.Get("max_device").(int)
	internetAccess := data.Get("internet_access").(string)
	configSystemSubnets := data.Get("system_subnets").([]interface{})
	var systemSubnets []string
	for _, s := range configSystemSubnets {
		systemSubnets = append(systemSubnets, s.(string))
	}
	configVpnRegionIds := data.Get("vpn_region_ids").([]interface{})
	var vpnRegionIds []string
	for _, r := range configVpnRegionIds {
		vpnRegionIds = append(vpnRegionIds, r.(string))
	}

	ug := &client.UserGroup{
		Name:           name,
		ConnectAuth:    connectAuth,
		MaxDevice:      maxDevice,
		SystemSubnets:  systemSubnets,
		VpnRegionIds:   vpnRegionIds,
		InternetAccess: internetAccess,
	}
	return ug
}

func updateUserGroupData(data *schema.ResourceData, userGroup *client.UserGroup) {
	data.SetId(userGroup.Id)
	_ = data.Set("connect_auth", userGroup.ConnectAuth)
	_ = data.Set("max_device", userGroup.MaxDevice)
	_ = data.Set("name", userGroup.Name)
	_ = data.Set("system_subnets", userGroup.SystemSubnets)
	_ = data.Set("vpn_region_ids", userGroup.VpnRegionIds)
	_ = data.Set("internet_access", userGroup.InternetAccess)
}

func resourceUserGroupDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*client.Client)
	var diags diag.Diagnostics
	err := c.DeleteUserGroup(data.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	data.SetId("")
	return diags
}

func resourceUserGroupRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*client.Client)
	var diags diag.Diagnostics
	userGroup, err := c.GetUserGroupById(data.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if userGroup == nil {
		data.SetId("")
	} else {
		updateUserGroupData(data, userGroup)
	}
	return diags
}

func resourceUserGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	ug := resourceDataToUserGroup(d)

	userGroup, err := c.CreateUserGroup(ug)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	updateUserGroupData(d, userGroup)
	return diags
}
