resource "cloudconnexa_user_group" "this" {
  name           = "test-group"
  vpn_region_ids = ["eu-central-1"]
  connect_auth   = "AUTH"
}
