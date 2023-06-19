resource "openvpncloud_user" "this" {
  for_each   = var.users
  username   = each.value.username
  email      = each.value.email
  first_name = split("_", each.key)[0]
  last_name  = split("_", each.key)[1]
  group_id   = lookup(var.groups, each.value.group)
  role       = each.value.role
}
