resource "openvpncloud_route" "example-routes" {
  for_each = {
    for key, route in var.example-terraform_ipv4_routes : route.value => route
  }
  network_item_id = var.networks["example-network"]
  type            = "IP_V4"
  value           = each.value.value
  description     = each.value.description
}
