package openvpncloud

import (
	"context"

	"github.com/OpenVPN/terraform-provider-openvpn-cloud/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `openvpncloud_user` to create an OpenVPN Cloud user.",
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "A username for the user.",
			},
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "An invitation to OpenVPN cloud account will be sent to this email. It will include an initial password and a VPN setup guide.",
			},
			"first_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 20),
				Description:  "User's first name.",
			},
			"last_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 20),
				Description:  "User's last name.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The UUID of a user's group.",
			},
			"role": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "MEMBER",
				Description: "The type of user role. Valid values are `ADMIN`, `MEMBER`, or `OWNER`.",
			},
			"devices": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "When a user signs in, the device that they use will be added to their account. You can read more at [OpenVPN Cloud Device](https://openvpn.net/cloud-docs/device/).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 32),
							Description:  "A device name.",
						},
						"description": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 120),
							Description:  "A device description.",
						},
						"ipv4_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "An IPv4 address of the device.",
						},
						"ipv6_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "An IPv6 address of the device.",
						},
					},
				},
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	username := d.Get("username").(string)
	email := d.Get("email").(string)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	groupId := d.Get("group_id").(string)
	role := d.Get("role").(string)
	configDevices := d.Get("devices").([]interface{})
	var devices []client.Device
	for _, d := range configDevices {
		device := d.(map[string]interface{})
		devices = append(
			devices,
			client.Device{
				Name:        device["name"].(string),
				Description: device["description"].(string),
				IPv4Address: device["ipv4_address"].(string),
				IPv6Address: device["ipv6_address"].(string),
			},
		)

	}
	u := client.User{
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		GroupId:   groupId,
		Devices:   devices,
		Role:      role,
	}
	user, err := c.CreateUser(u)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(user.Id)
	return append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "The user's role cannot be changed using the code.",
		Detail:   "There is a bug in OpenVPN Cloud API that prevents setting the user's role during the creation. All users are created as Members by default. Once it's fixed, the provider will be updated accordingly.",
	})
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	userId := d.Id()
	u, err := c.GetUserById(userId)

	// If group_id is not set, OpenVPN cloud sets it to the default group.
	var groupId string
	if d.Get("group_id") == "" {
		// The group has not been explicitly set.
		// Set it to an empty string to keep the default group.
		groupId = ""
	} else {
		groupId = u.GroupId
	}

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if u == nil {
		d.SetId("")
	} else {
		d.Set("username", u.Username)
		d.Set("email", u.Email)
		d.Set("first_name", u.FirstName)
		d.Set("last_name", u.LastName)
		d.Set("group_id", groupId)
		d.Set("devices", u.Devices)
		d.Set("role", u.Role)
	}
	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	if !d.HasChanges("first_name", "last_name", "group_id", "email") {
		return diags
	}

	u, err := c.GetUserById(d.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	_, email := d.GetChange("email")
	_, firstName := d.GetChange("first_name")
	_, lastName := d.GetChange("last_name")
	_, role := d.GetChange("role")
	oldGroupId, newGroupId := d.GetChange("group_id")
	// The group has not been set explicitly.
	// The update endpoint requires group_id to be set, so we should set it to the default group.
	groupId := newGroupId.(string)
	if oldGroupId.(string) == "" && groupId == "" {
		g, err := c.GetUserGroup("Default")
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		groupId = g.Id
	}
	status := u.AccountStatus

	err = c.UpdateUser(client.User{
		Id:        d.Id(),
		Email:     email.(string),
		FirstName: firstName.(string),
		LastName:  lastName.(string),
		GroupId:   groupId,
		Role:      role.(string),
		Status:    status,
	})

	return append(diags, diag.FromErr(err)...)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	userId := d.Id()
	err := c.DeleteUser(userId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
