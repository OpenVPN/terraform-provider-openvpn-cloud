data "openvpncloud_network" "this" {
  name = "test-net"
}

resource "openvpncloud_service" "this" {
  name              = "test-service"
  type              = "IP_SOURCE"
  description       = "test-description"
  routes            = ["10.0.0.2/32"]
  network_item_type = "NETWORK"
  network_item_id   = data.openvpncloud_network.this.network_id

  config {
    service_types = ["ANY"]
  }
}
