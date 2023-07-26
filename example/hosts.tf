resource "openvpncloud_host" "this" {
  name = "test-host"
  connector {
    name          = "test-connector"
    vpn_region_id = "eu-central-1"
  }
}
