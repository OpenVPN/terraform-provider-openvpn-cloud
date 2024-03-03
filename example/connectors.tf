data "cloudconnexa_network" "test-net" {
  name = "test-net"
}

resource "cloudconnexa_connector" "test-connector" {
  name              = "test-connector"
  vpn_region_id     = "eu-central-1"
  network_item_type = "NETWORK"
  network_item_id   = data.cloudconnexa_network.test-net.network_id
}
