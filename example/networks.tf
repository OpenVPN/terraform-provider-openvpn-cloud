resource "cloudconnexa_network" "test-network" {
  name   = "test-network"
  egress = false
  default_route {
    value = "192.168.0.0/24"
  }
  default_connector {
    name          = "test-connector"
    vpn_region_id = "eu-central-1"
  }
}
