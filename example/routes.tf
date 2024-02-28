resource "cloudconnexa_route" "this" {
  for_each = {
    for key, route in var.routes : route.value => route
  }
  network_item_id = var.networks["example-network"]
  type            = "IP_V4"
  value           = each.value.value
  description     = each.value.description
}
