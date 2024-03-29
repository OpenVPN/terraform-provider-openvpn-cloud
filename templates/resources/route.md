---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudconnexa_route Resource - terraform-provider-cloudconnexa"
subcategory: ""
description: |-
  Use cloudconnexa_route to create a route on an Cloud Connexa network.
---

# cloudconnexa_route (Resource)

Use `cloudconnexa_route` to create a route on an Cloud Connexa network.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `network_item_id` (String) The id of the network on which to create the route.
- `type` (String) The type of route. Valid values are `IP_V4`, `IP_V6`, and `DOMAIN`.
- `value` (String) The target value of the default route.

### Read-Only

- `id` (String) The ID of this resource.

## Import

A route can be imported using the route ID, which can be fetched directly from the API.

```
terraform import cloudconnexa_route.route <route-uuid>
```