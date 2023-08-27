package openvpncloud

import (
	"context"

	"github.com/OpenVPN/terraform-provider-openvpn-cloud/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDnsRecord() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `openvpncloud_dns_record` to create a DNS record on your VPN.",
		CreateContext: resourceDnsRecordCreate,
		ReadContext:   resourceDnsRecordRead,
		DeleteContext: resourceDnsRecordDelete,
		UpdateContext: resourceDnsRecordUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The DNS record name.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "The description for the UI. Defaults to `Managed by Terraform`.",
			},
			"ip_v4_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv4Address,
				},
				Description: "The list of IPV4 addresses to which this record will resolve.",
			},
			"ip_v6_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv6Address,
				},
				Description: "The list of IPV6 addresses to which this record will resolve.",
			},
		},
	}
}

func resourceDnsRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	domain := d.Get("domain").(string)
	description := d.Get("description").(string)
	ipV4Addresses := d.Get("ip_v4_addresses").([]interface{})
	ipV4AddressesSlice := make([]string, 0)
	for _, a := range ipV4Addresses {
		ipV4AddressesSlice = append(ipV4AddressesSlice, a.(string))
	}
	ipV6Addresses := d.Get("ip_v6_addresses").([]interface{})
	ipV6AddressesSlice := make([]string, 0)
	for _, a := range ipV6Addresses {
		ipV6AddressesSlice = append(ipV6AddressesSlice, a.(string))
	}
	dr := client.DnsRecord{
		Domain:        domain,
		Description:   description,
		IPV4Addresses: ipV4AddressesSlice,
		IPV6Addresses: ipV6AddressesSlice,
	}
	dnsRecord, err := c.CreateDnsRecord(dr)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(dnsRecord.Id)
	return diags
}

func resourceDnsRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	recordId := d.Id()
	r, err := c.GetDnsRecord(recordId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if r == nil {
		d.SetId("")
	} else {
		d.Set("domain", r.Domain)
		d.Set("description", r.Description)
		d.Set("ip_v4_addresses", r.IPV4Addresses)
		d.Set("ip_v6_addresses", r.IPV6Addresses)
	}
	return diags
}

func resourceDnsRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	_, domain := d.GetChange("domain")
	_, description := d.GetChange("description")
	_, ipV4Addresses := d.GetChange("ip_v4_addresses")
	ipV4AddressesSlice := getAddressesSlice(ipV4Addresses.([]interface{}))
	_, ipV6Addresses := d.GetChange("ip_v6_addresses")
	ipV6AddressesSlice := getAddressesSlice(ipV6Addresses.([]interface{}))
	dr := client.DnsRecord{
		Id:            d.Id(),
		Domain:        domain.(string),
		Description:   description.(string),
		IPV4Addresses: ipV4AddressesSlice,
		IPV6Addresses: ipV6AddressesSlice,
	}
	err := c.UpdateDnsRecord(dr)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

func resourceDnsRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	routeId := d.Id()
	err := c.DeleteDnsRecord(routeId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

func getAddressesSlice(addresses []interface{}) []string {
	addressesSlice := make([]string, 0)
	for _, a := range addresses {
		addressesSlice = append(addressesSlice, a.(string))
	}
	return addressesSlice
}
