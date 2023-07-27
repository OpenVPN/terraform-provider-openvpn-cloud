data "openvpncloud_network" "this" {
  name = "test-net"
}

resource "openvpncloud_connector" "this" {
  name              = "test-connector"
  vpn_region_id     = "eu-central-1"
  network_item_type = "NETWORK"
  network_item_id   = data.openvpncloud_network.this.network_id
}
