---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudconnexa_dns_record Resource - terraform-provider-cloudconnexa"
subcategory: ""
description: |-
  Use cloudconnexa_dns_record to create a DNS record on your VPN.
---

# cloudconnexa_dns_record (Resource)

Use `cloudconnexa_dns_record` to create a DNS record on your VPN.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain` (String) The DNS record name.

### Optional

- `ip_v4_addresses` (List of String) The list of IPV4 addresses to which this record will resolve.
- `ip_v6_addresses` (List of String) The list of IPV6 addresses to which this record will resolve.

### Read-Only

- `id` (String) The ID of this resource.

## Import

A connector can be imported using the DNS record ID, which can be fetched directly from the API.

```
terraform import cloudconnexa_dns_record.record <record-uuid>
```