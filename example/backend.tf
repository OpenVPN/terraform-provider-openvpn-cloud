terraform {
  backend "local" {}
  required_providers {
    openvpncloud = {
      source  = "OpenVPN/openvpn-cloud"
      version = "0.0.7"
    }
  }
}
